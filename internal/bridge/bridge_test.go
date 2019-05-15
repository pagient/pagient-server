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
		roomAssignments []*bridgeModel.RoomAssignment
		dbError         error
		patients        []*model.Patient
	}{
		"no room assignments": {
			roomAssignments: nil,
			dbError:         nil,
			patients:        []*model.Patient{},
		},
		"empty room assignments": {
			roomAssignments: []*bridgeModel.RoomAssignment{},
			dbError:         nil,
			patients:        []*model.Patient{},
		},
		"some room assignments": {
			roomAssignments: []*bridgeModel.RoomAssignment{
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
			dbError: nil,
			patients: []*model.Patient{
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
			roomAssignments: nil,
			dbError:         errors.New("sample test error"),
			patients:        nil,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		db := &MockDB{}
		db.On("GetRoomAssignments", mock.AnythingOfType("string"), mock.AnythingOfType("uint")).Return(test.roomAssignments, test.dbError).Once()

		bridge := NewBridge(db)

		patientsExaminedNext, err := bridge.GetToBeExaminedPatients()
		assert.ElementsMatch(t, test.patients, patientsExaminedNext)
		if test.dbError != nil {
			assert.Error(t, err)
			assert.EqualError(t, test.dbError, errors.Cause(err).Error())
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestDefaultBridge_GetExaminedPatients(t *testing.T) {
	tests := map[string]struct {
		lastAssignments []*bridgeModel.RoomAssignment
		roomAssignments []*bridgeModel.RoomAssignment
		dbError         error
		patients        []*model.Patient
	}{
		"no last assignments": {
			lastAssignments: nil,
			roomAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 1,
				},
			},
			dbError:  nil,
			patients: nil,
		},
		"no examined patients": {
			lastAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 1,
				},
				{
					PID: 3,
				},
			},
			roomAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 3,
				},
				{
					PID: 1,
				},
			},
			dbError:  nil,
			patients: nil,
		},
		"no more room assignments": {
			lastAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 1,
				},
				{
					PID: 2,
				},
			},
			roomAssignments: nil,
			dbError:         nil,
			patients: []*model.Patient{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
		},
		"all last assignments are now finished patients": {
			lastAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 1,
				},
				{
					PID: 2,
				},
			},
			roomAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 3,
				},
				{
					PID: 4,
				},
			},
			dbError: nil,
			patients: []*model.Patient{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
		},
		"a few examined patients": {
			lastAssignments: []*bridgeModel.RoomAssignment{
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
			roomAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 3,
				},
				{
					PID: 2,
				},
			},
			dbError: nil,
			patients: []*model.Patient{
				{
					ID: 1,
				},
				{
					ID: 4,
				},
			},
		},
		"database error": {
			lastAssignments: []*bridgeModel.RoomAssignment{
				{
					PID: 1,
				},
				{
					PID: 2,
				},
			},
			roomAssignments: nil,
			dbError:         errors.New("sample test error"),
			patients:        nil,
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		db := &MockDB{}
		callCount := 0
		db.On("GetRoomAssignments", mock.AnythingOfType("string"), mock.AnythingOfType("uint")).
			Return(func(s string, u ...uint) []*bridgeModel.RoomAssignment {
				callCount++
				if callCount == 1 {
					return test.lastAssignments
				}
				return test.roomAssignments
			}, test.dbError).
			Times(2)

		bridge := NewBridge(db)

		patientsExamined, err := bridge.GetExaminedPatients()
		assert.ElementsMatch(t, nil, patientsExamined)

		patientsExamined, err = bridge.GetExaminedPatients()
		assert.ElementsMatch(t, test.patients, patientsExamined)
		if test.dbError != nil {
			assert.Error(t, err)
			assert.EqualError(t, test.dbError, errors.Cause(err).Error())
		} else {
			assert.NoError(t, err)
		}
	}
}
