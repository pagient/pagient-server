package bridge

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/pagient/pagient-server/pkg/config"
	"github.com/pagient/pagient-server/pkg/model"
	"github.com/pagient/pagient-server/pkg/notifier"
	"github.com/pagient/pagient-server/pkg/service"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Bridge struct {
	cfg            *config.Config
	patientService service.PatientService
	notifier       notifier.Notifier
}

func NewBridge(cfg *config.Config, patientService service.PatientService, notifier notifier.Notifier) *Bridge {
	return &Bridge{
		cfg:            cfg,
		patientService: patientService,
		notifier:       notifier,
	}
}

func (bridge *Bridge) Run(stop <-chan struct{}) error {
	connectionString := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&encrypt=disable",
		bridge.cfg.Bridge.DbUser,
		bridge.cfg.Bridge.DbPassword,
		bridge.cfg.Bridge.DbUrl,
		bridge.cfg.Bridge.DbName)

	db, err := gorm.Open("mssql", connectionString)
	if err != nil {
		return errors.Wrap(err, "failed to connect database")
	}
	db.LogMode(zerolog.GlobalLevel() <= zerolog.DebugLevel)
	db.SetLogger(&log.Logger)
	defer db.Close()

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				patients, err := bridge.patientService.GetAll()
				if err != nil {
					log.Error().
						Err(err).
						Msg("get all patients failed")

					continue
				}

				patients = bridge.filterPatients(patients)

				if len(patients) == 0 {
					continue
				}

				var assignments []*patientRoomAssignment
				db.Raw("SELECT TOP(?) PDS6_WZ.* FROM PDS6_WZ JOIN PDS6_STWZ ON PDS6_WZ.WZID = PDS6_STWZ.WZID "+
					"WHERE PDS6_STWZ.CODE = ? ORDER BY PDS6_WZ.FLGNR ASC",
					bridge.cfg.Bridge.CallActionQueuePosition, bridge.cfg.Bridge.CallActionWZ).
					Scan(&assignments)

				if len(assignments) == 0 {
					continue
				}

				// loop over patients to check if there's a room assignment
				for _, patient := range patients {
					for _, assignment := range assignments {
						// there's an assignment queued in the first X items
						if assignment.PatientID == patient.ID {
							// call patient
							patient.Status = model.PatientStateCall
							if _, err := bridge.patientService.Update(patient); err != nil {
								log.Error().
									Err(err).
									Msg("update patient failed")
							}
							break
						}
					}
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

func (bridge *Bridge) filterPatients(patients []*model.Patient) []*model.Patient {
	var filteredPatients []*model.Patient

	for _, patient := range patients {
		if patient.PagerID != 0 &&
			(patient.Status != model.PatientStateCall &&
				patient.Status != model.PatientStateCalled) {
			filteredPatients = append(filteredPatients, patient)
		}
	}

	return filteredPatients
}
