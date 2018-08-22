# FirestoreRestore

[![Docker Automated build](https://img.shields.io/docker/automated/christianalexander/firestorerestore.svg?style=popout-square)](https://hub.docker.com/r/christianalexander/firestorerestore/)

An _experimental_ utility to export and import Google Cloud Firestore databases.

See the blog post at https://christianalexander.com/2018/07/18/firestore-backup-restore/

## Requirements

1. Create a [Google Cloud Storage bucket](https://console.cloud.google.com/storage/) in the same project as the Firestore Database.
2. Create a [service account](https://console.cloud.google.com/iam-admin/serviceaccounts/project) with the "Cloud Datastore Import Export Admin" role.
3. Store the JSON key somewhere it can be retrieved later.

## Building and Running

### From Release
1. Grab a release from the [releases](https://github.com/ChristianAlexander/FirestoreRestore/releases) page
2. Extract the release
3. Run `FirestoreRestore`

### From Source
1. Clone this repo
2. Run `make`
3. Run `./FirestoreRestore`

## Usage

*All commands must be run with the `GOOGLE_APPLICATION_CREDENTIALS` environment variable set to the path where the JSON key is located. See [Google documentation](https://cloud.google.com/docs/authentication/production/#setting_the_environment_variable) for more details.*

### Backup

`FirestoreRestore -backup -wait -p <GOOGLE PROJECT ID> -b <BUCKET NAME> -n <BACKUP NAME>`

If `-n` is not specified, the backup will be named after the current time.

### Restore

`FirestoreRestore -wait -p <GOOGLE PROJECT ID> -b <BUCKET NAME> -n <BACKUP NAME>`

The `-n` value should be a path relative to the root of the bucket, such as `backups/abcd`
