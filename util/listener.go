package util

import (
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// CosListener is a struct representing a listener for COS progress events.
type CosListener struct {
	fo      *FileOperations
	counter *Counter
}

// ProgressChangedCallback is a callback function that is triggered when the progress of the COS operation changes.
// It updates the transfer size and deal size based on the event type and consumed bytes.
func (l *CosListener) ProgressChangedCallback(event *cos.ProgressEvent) {
	switch event.EventType {
	case cos.ProgressStartedEvent:
	case cos.ProgressDataEvent:
		l.fo.Monitor.updateTransferSize(event.RWBytes)
		l.fo.Monitor.updateDealSize(event.RWBytes)
		l.counter.TransferSize += event.RWBytes
	case cos.ProgressCompletedEvent:
	case cos.ProgressFailedEvent:
		l.fo.Monitor.updateDealSize(-event.ConsumedBytes)
	default:
		fmt.Printf("Progress Changed Error: unknown progress event type\n")
	}
	freshProgress()
}
