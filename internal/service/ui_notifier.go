package service

import "github.com/pagient/pagient-server/internal/model"

// UINotifier interface for async view updates
type UINotifier interface {
	NotifyNewPatient(*model.Patient)
	NotifyUpdatedPatient(*model.Patient)
	NotifyDeletedPatient(*model.Patient)
}
