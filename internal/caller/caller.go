package caller

import (
	"sort"
	"time"

	"github.com/pagient/pagient-server/internal/model"
	"github.com/pagient/pagient-server/internal/notifier"
	"github.com/pagient/pagient-server/internal/service"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// SoftwareBridge provides abstraction for different practitioner software
type SoftwareBridge interface {
	GetToBeExaminedPatients() ([]*model.Patient, error)
	GetExaminedPatients() ([]*model.Patient, error)
}

// Caller struct encapsulates the surgery software bridge
type Caller struct {
	service  service.PatientService
	notifier notifier.Notifier
	bridge   SoftwareBridge
}

// NewCaller returns a surgery software bridge struct
func NewCaller(s service.PatientService, bridge SoftwareBridge) *Caller {
	return &Caller{
		service: s,
		bridge:  bridge,
	}
}

// Run runs the bridge functionality in a new goroutine repeated by given every every
func (c *Caller) Run(every time.Duration, stop <-chan struct{}) error {
	ticker := time.NewTicker(every)
	go func() {
		for {
			select {
			case <-ticker.C:
				patients, err := c.service.ListPagerPatientsByStatus(model.PatientStatusPending)
				if err != nil {
					log.Error().
						Err(err).
						Msg("get not yet alerted patients having pagers failed")

					continue
				}

				queuedPatients, err := c.bridge.GetToBeExaminedPatients()
				if err != nil {
					log.Error().
						Err(err).
						Msg("get to be examined patients from software bridge failed")

					continue
				}

				toBeCalledPatients := intersect(patients, queuedPatients)
				if err := c.callPatients(toBeCalledPatients); err != nil {
					log.Error().
						Err(err).
						Msg("call patients failed")

					continue
				}

				patients, err = c.service.ListPagerPatientsByStatus(model.PatientStatusPending, model.PatientStatusCall, model.PatientStatusCalled)
				if err != nil {
					log.Error().
						Err(err).
						Msg("get examined/finished patients having pagers failed")

					continue
				}

				finishedPatients, err := c.bridge.GetExaminedPatients()
				if err != nil {
					log.Error().
						Err(err).
						Msg("get examined patients from software bridge failed")

					continue
				}

				notReturnedPagerPatients := intersect(patients, finishedPatients)
				if err := c.markExaminedPatientsFinished(notReturnedPagerPatients); err != nil {
					log.Error().
						Err(err).
						Msg("set patients finished failed")

					continue
				}
			case <-stop:
				// close goroutine
				ticker.Stop()
				return
			}
		}
	}()

	<-stop

	return nil
}

func (c *Caller) callPatients(patients []*model.Patient) error {
	for _, patient := range patients {
		if err := c.service.CallPatient(patient); err != nil {
			return errors.Wrap(err, "call patient failed")
		}
	}

	return nil
}

func (c *Caller) markExaminedPatientsFinished(patients []*model.Patient) error {
	for _, patient := range patients {
		patient.Status = model.PatientStatusFinished
		if _, err := c.service.UpdatePatient(patient); err != nil {
			return errors.Wrap(err, "update patient failed")
		}
	}

	return nil
}

func intersect(patientsA, patientsB []*model.Patient) []*model.Patient {
	sortPatientsByID(patientsB)

	intersection := make([]*model.Patient, min(len(patientsA), len(patientsB)))
	for _, patientA := range patientsA {
		for _, patientB := range patientsB {
			if patientA.ID < patientB.ID {
				break
			}

			if patientA.ID == patientB.ID {
				intersection = append(intersection, patientA)
			}

		}
	}

	return intersection
}

func sortPatientsByID(patients []*model.Patient) {
	sort.Slice(patients, func(i, j int) bool {
		return patients[i].ID < patients[j].ID
	})
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
