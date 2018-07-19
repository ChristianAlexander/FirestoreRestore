package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	firestore "google.golang.org/api/firestore/v1beta1"
)

var (
	projectID, bucket, backupName string
	wait, shouldBackup            bool
)

func init() {
	flag.StringVar(&projectID, "p", "", "[required] the Google Cloud Project in which the Firestore database and the Google Cloud Storage bucket exist")
	flag.StringVar(&bucket, "b", "", "[required] the Google Cloud Storage bucket for backup storage")
	flag.StringVar(&backupName, "n", time.Now().Format(time.RFC3339), "the name of the backup, defaults to the current datetime")
	flag.BoolVar(&wait, "wait", false, "wait for the operation to complete")
	flag.BoolVar(&shouldBackup, "backup", false, "true to backup, false to restore")

	flag.Parse()

	if projectID == "" || bucket == "" {
		flag.Usage()
		os.Exit(1)
	}

	if !shouldBackup && backupName == "" {
		logrus.Fatalf("Must specify a backup name with the -n flag")
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logrus.Warnf("Received signal %s", sig)
		cancel()
	}()

	hc, err := google.DefaultClient(ctx, firestore.CloudPlatformScope, firestore.DatastoreScope)
	if err != nil {
		logrus.Fatalf("Failed to create Google HTTP client: %v", err)
	}

	f, err := firestore.New(hc)
	if err != nil {
		logrus.Fatalf("Failed to create Firestore client: %v", err)
	}

	var jobName string

	if shouldBackup {
		n, err := backup(ctx, f)
		if err != nil {
			logrus.Fatalf("Failed to back up: %v", err)
		}

		logrus.Infof("Created backup operation '%s'", n)
		jobName = n
	} else {
		n, err := restore(ctx, f)
		if err != nil {
			logrus.Fatalf("Failed to restore: %v", err)
		}
		logrus.Infof("Created restore operation '%s'", n)
		jobName = n
	}

	if wait {
		err := waitOnJob(ctx, hc, jobName)
		if err != nil {
			logrus.Fatalf("Failed while waiting on job: %v", err)
		}
	}

	logrus.Infoln("Process completed")
}

func backup(ctx context.Context, f *firestore.Service) (string, error) {
	backupURI := fmt.Sprintf("gs://%s/%s", bucket, backupName)
	logrus.Infof("Requesting backup of '%s' to '%s'", projectID, backupURI)

	res, err := f.Projects.Databases.ExportDocuments(fmt.Sprintf("projects/%s/databases/(default)", projectID), &firestore.GoogleFirestoreAdminV1beta1ExportDocumentsRequest{
		OutputUriPrefix: backupURI,
	}).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create backup operation: %v", err)
	}

	return res.Name, nil
}

func restore(ctx context.Context, f *firestore.Service) (string, error) {
	backupURI := fmt.Sprintf("gs://%s/%s", bucket, backupName)
	logrus.Infof("Restoring backup '%s' to '%s'", backupURI, projectID)

	res, err := f.Projects.Databases.ImportDocuments(fmt.Sprintf("projects/%s/databases/(default)", projectID), &firestore.GoogleFirestoreAdminV1beta1ImportDocumentsRequest{
		InputUriPrefix: backupURI,
	}).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create restore operation: %v", err)
	}

	return res.Name, nil
}

func waitOnJob(ctx context.Context, c *http.Client, jobName string) error {
	t := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			{
				return nil
			}
		case <-t.C:
			{
				res, err := c.Get(fmt.Sprintf("https://firestore.googleapis.com/v1/%s", jobName))
				if err != nil {
					return fmt.Errorf("failed to get job status: %v", err)
				}

				b, err := ioutil.ReadAll(res.Body)
				if err != nil {
					return fmt.Errorf("failed to read response: %v", err)
				}

				var md JobMetadata
				err = json.Unmarshal(b, &md)
				if err != nil {
					return fmt.Errorf("failed to unmarshal response: %v", err)
				}

				if md.Metadata.OperationState == OperationStateCancelled ||
					md.Metadata.OperationState == OperationStateFailed {
					return fmt.Errorf("operation did not complete: %s", md.Metadata.OperationState)
				}

				if md.Metadata.OperationState == OperationStateSuccessful {
					return nil
				}
			}
		}
	}
}
