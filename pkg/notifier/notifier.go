package notifier

import "github.com/pagient/pagient-server/pkg/model"

// Notifier interface for async view updates
type Notifier interface {
	NotifyNewPatient(*model.Patient)
	NotifyUpdatedPatient(*model.Patient)
	NotifyDeletedPatient(*model.Patient)
}
