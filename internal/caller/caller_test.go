package caller

import (
	"testing"
	"time"

	"github.com/pagient/pagient-server/internal/model"
	"github.com/pagient/pagient-server/internal/service"

	"github.com/stretchr/testify/assert"
)

func TestCaller_Run(t *testing.T) {
	patientPool := map[int]*model.Patient{
		1: {
			ID: 1,
		},
		2: {
			ID: 2,
		},
		3: {
			ID: 3,
		},
		4: {
			ID: 4,
		},

		5: {
			ID: 5,
		},
		6: {
			ID: 6,
		},
	}

	tests := map[string]struct {
		every            time.Duration
		repeats          int
		patients         []*model.Patient
		toBeExamined     []*model.Patient
		calledPatients   []*model.Patient
		haveBeenExamined []*model.Patient
		finishedPatients []*model.Patient
	}{
		"should run every given every": {
			every:            time.Duration(50) * time.Millisecond,
			repeats:          2,
			patients:         nil,
			toBeExamined:     nil,
			calledPatients:   nil,
			haveBeenExamined: nil,
			finishedPatients: nil,
		},
		"should call \"pending\" patients that are examined next": {
			every:   time.Duration(1) * time.Millisecond,
			repeats: 1,
			patients: []*model.Patient{
				patientPool[1],
				patientPool[2],
				patientPool[3],
				patientPool[4],
			},
			toBeExamined: []*model.Patient{
				patientPool[3],
				patientPool[6],
				patientPool[1],
				patientPool[5],
			},
			calledPatients: []*model.Patient{
				patientPool[3],
				patientPool[1],
			},
			haveBeenExamined: nil,
			finishedPatients: nil,
		},
		"should set status \"finished\" for examined patients": {
			every:   time.Duration(1) * time.Millisecond,
			repeats: 1,
			patients: []*model.Patient{
				patientPool[1],
				patientPool[2],
				patientPool[3],
				patientPool[4],
			},
			toBeExamined:   nil,
			calledPatients: nil,
			haveBeenExamined: []*model.Patient{
				patientPool[3],
				patientPool[6],
				patientPool[1],
				patientPool[5],
			},
			finishedPatients: []*model.Patient{
				patientPool[3],
				patientPool[1],
			},
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		s := &service.MockService{}
		s.On("ListPagerPatientsByStatus", model.PatientStatusPending).Return(test.patients, nil)
		s.On("ListPagerPatientsByStatus", model.PatientStatusPending, model.PatientStatusCall, model.PatientStatusCalled).Return(test.patients, nil)

		for _, patient := range test.calledPatients {
			s.
				On("CallPatient", patient).
				Return(nil).
				Times(test.repeats)
		}

		for _, patient := range test.finishedPatients {
			s.
				On("UpdatePatient", patient).
				Return(nil).
				Times(test.repeats)
		}

		b := &MockSoftwareBridge{}
		b.On("GetToBeExaminedPatients").
			Return(test.toBeExamined, nil)
		b.On("GetExaminedPatients").
			Return(test.haveBeenExamined, nil)

		caller := NewCaller(s, b)
		stop := make(chan struct{}, 1)

		go func() {
			err := caller.Run(test.every, stop)
			assert.NoError(t, err)
		}()

		// wait til all repeats should have been executed plus half the every duration for safety
		maxDur := test.every*time.Duration(test.repeats) + test.every/2
		<-time.After(maxDur)
		close(stop)
	}
}
