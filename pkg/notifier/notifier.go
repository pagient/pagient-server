package notifier

import "github.com/pagient/pagient-server/pkg/model"

type Notifier interface {
	NotifyNewPatient(*model.Patient)
	NotifyUpdatedPatient(*model.Patient)
	NotifyDeletedPatient(*model.Patient)
}
