package bridge

import (
	"testing"

	bridgeModel "github.com/pagient/pagient-server/internal/bridge/model"
	"github.com/pagient/pagient-server/internal/model"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDefaultBridge_GetToBeExaminedPatients(t *testing.T) {
	tests := map[string]struct {
		RoomAssignments []*bridgeModel.RoomAssignment
		DBError         error
		Patients        []*model.Patient
	}{
		"no room assignments": {
			RoomAssignments: nil,
			DBError: nil,
			Patients: []*model.Patient{},
		},
		"empty room assignments": {
			RoomAssignments: []*bridgeModel.RoomAssignment{},
			DBError: nil,
			Patients: []*model.Patient{},
		},
		"some room assignments": {
			RoomAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 1,
				},
				{
					PID: 2,
				},
				{
					PID: 3,
				},
			},
			DBError: nil,
			Patients: []*model.Patient{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
				{
					ID: 3,
				},
			},
		},
		"database error": {
			RoomAssignments: nil,
			DBError: errors.New("sample test error"),
			Patients: nil,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		tx := &MockTx{}
		tx.On("Commit").Return(nil).Once()
		tx.On("Rollback").Return(nil).Once()
		tx.On("GetRoomAssignments", mock.AnythingOfType("string"), mock.AnythingOfType("uint")).Return(test.RoomAssignments, test.DBError).Once()

		db:= &MockDB{}
		db.On("Begin").Return(tx, nil).Once()

		bridge := NewBridge(db)

		patientsExaminedNext, err := bridge.GetToBeExaminedPatients()
		assert.ElementsMatch(t, test.Patients, patientsExaminedNext)
		if test.DBError != nil {
			assert.Error(t, err)
			assert.EqualError(t, test.DBError, errors.Cause(err).Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestDefaultBridge_GetExaminedPatients(t *testing.T) {
	tests := map[string]struct {
		LastAssignments []*bridgeModel.RoomAssignment
		RoomAssignments []*bridgeModel.RoomAssignment
		DBError         error
		Patients        []*model.Patient
	}{
		"no last assignments": {
			LastAssignments: nil,
			RoomAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 1,
				},
			},
			DBError: nil,
			Patients: nil,
		},
		"no examined patients": {
			LastAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 1,
				},
				{
					PID: 3,
				},
			},
			RoomAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 3,
				},
				{
					PID: 1,
				},
			},
			DBError: nil,
			Patients: nil,
		},
		"no more room assignments": {
			LastAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 1,
				},
				{
					PID: 2,
				},
			},
			RoomAssignments: nil,
			DBError: nil,
			Patients: []*model.Patient{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
		},
		"all last assignments are now finished patients": {
			LastAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 1,
				},
				{
					PID: 2,
				},
			},
			RoomAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 3,
				},
				{
					PID: 4,
				},
			},
			DBError: nil,
			Patients: []*model.Patient{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
		},
		"a few examined patients": {
			LastAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 1,
				},
				{
					PID: 2,
				},
				{
					PID: 3,
				},
				{
					PID: 4,
				},
			},
			RoomAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 3,
				},
				{
					PID: 2,
				},
			},
			DBError: nil,
			Patients: []*model.Patient{
				{
					ID: 1,
				},
				{
					ID: 4,
				},
			},
		},
		"database error": {
			LastAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 1,
				},
				{
					PID: 2,
				},
			},
			RoomAssignments: nil,
			DBError: errors.New("sample test error"),
			Patients: nil,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		tx := &MockTx{}
		tx.On("Commit").Return(nil).Times(2)
		tx.On("Rollback").Return(nil).Times(2)

		callCount := 0
		tx.On("GetRoomAssignments", mock.AnythingOfType("string"), mock.AnythingOfType("uint")).
			Return(func (s string, u ...uint) []*bridgeModel.RoomAssignment {
				callCount++
				if callCount == 1 {
					return test.LastAssignments
				}
				return test.RoomAssignments
			}, test.DBError).
			Times(2)

		db:= &MockDB{}
		db.On("Begin").Return(tx, nil).Times(2)

		bridge := NewBridge(db)

		patientsExamined, err := bridge.GetExaminedPatients()
		assert.ElementsMatch(t, nil, patientsExamined)

		patientsExamined, err = bridge.GetExaminedPatients()
		assert.ElementsMatch(t, test.Patients, patientsExamined)
		if test.DBError != nil {
			assert.Error(t, err)
			assert.EqualError(t, test.DBError, errors.Cause(err).Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
