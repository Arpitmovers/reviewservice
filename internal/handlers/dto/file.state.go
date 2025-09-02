package handlers

type FileState string

const (
	FileStateSuccess    FileState = "Success"    // processing completed successfully
	FileStateFailed     FileState = "Failed"     // processing failed, retry possible
	FileStateError      FileState = "Error"      // unrecoverable error, skip/reject
	FileStateInProgress FileState = "InProgress" // currently being processed
)
