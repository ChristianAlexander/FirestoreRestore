package main

// JobMetadata contains the status of an import/export job.
type JobMetadata struct {
	Name     string `json:"name"`
	Metadata struct {
		Type            string         `json:"@type"`
		StartTime       string         `json:"startTime"`
		OperationState  OperationState `json:"operationState"`
		OutputURIPrefix string         `json:"outputUriPrefix"`
	}
}

// OperationState is the state of an import/export job
type OperationState string

var (
	// OperationStateUnspecified - Unspecified.
	OperationStateUnspecified OperationState = "STATE_UNSPECIFIED"

	// OperationStateInitializing - Request is being prepared for processing.
	OperationStateInitializing OperationState = "INITIALIZING"

	// OperationStateProcessing - Request is actively being processed.
	OperationStateProcessing OperationState = "PROCESSING"

	// OperationStateCancelling - Request is in the process of being cancelled after user called cancel on the operation.
	OperationStateCancelling OperationState = "CANCELLING"

	// OperationStateFinalizing - Request has been processed and is in its finalization stage.
	OperationStateFinalizing OperationState = "FINALIZING"

	// OperationStateSuccessful - Request has completed successfully.
	OperationStateSuccessful OperationState = "SUCCESSFUL"

	// OperationStateFailed - Request has finished being processed, but encountered an error.
	OperationStateFailed OperationState = "FAILED"

	// OperationStateCancelled - Request has finished being cancelled after user called cancel on the operation."
	OperationStateCancelled OperationState = "CANCELLED"
)
