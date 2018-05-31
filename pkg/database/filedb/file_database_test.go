package filedb

import (
	"fmt"
	"os"
	"testing"

	"github.com/pagient/pagient-api/pkg/model"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDriver struct {
	mock.Mock
}

func (d *mockDriver) Write(c, r string, v interface{}) error {
	args := d.Called(c, r, v)
	return args.Error(0)
}

func (d *mockDriver) Read(c, r string, v interface{}) error {
	args := d.Called(c, r, v)
	return args.Error(0)
}

func (d *mockDriver) ReadAll(c string) ([]string, error) {
	args := d.Called(c)
	return args.Get(0).([]string), args.Error(1)
}

func (d *mockDriver) Delete(c, r string) error {
	args := d.Called(c, r)
	return args.Error(0)
}

func TestFileDatabase_GetPatient(t *testing.T) {
	patientID := uuid.NewV4()

	tests := map[string]struct {
		patientID       uuid.UUID
		patientRet      *model.Patient
		patientRetErr   error
		resultingErrMsg string
	}{
		"successful": {
			patientID: patientID,
			patientRet: &model.Patient{
				ID:   patientID,
				Name: "test1",
			},
		},
		"with db error": {
			patientID:       patientID,
			patientRetErr:   assert.AnError,
			resultingErrMsg: assert.AnError.Error(),
		},
		"successful but patient or database not found": {
			patientID:       patientID,
			patientRetErr:   os.ErrNotExist,
			resultingErrMsg: "patient not found",
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		driver := &mockDriver{}
		driver.
			On("Read", patientCollection, test.patientID.String(), mock.AnythingOfType("*model.Patient")).
			Return(test.patientRetErr).
			Run(func(args mock.Arguments) {
				if test.patientRetErr == nil {
					arg := args.Get(2).(*model.Patient)
					arg.ID = test.patientRet.ID
					arg.Name = test.patientRet.Name
				}
			}).
			Once()

		db := &FileDatabase{
			driver: driver,
		}

		patient, err := db.GetPatient(test.patientID)
		assert.Equal(t, test.patientRet, patient)
		if test.resultingErrMsg != "" {
			assert.EqualError(t, err, test.resultingErrMsg)
		} else {
			assert.Equal(t, nil, err)
		}
		driver.AssertExpectations(t)
	}
}

func TestFileDatabase_GetPatients(t *testing.T) {
	patientID := uuid.NewV4()

	tests := map[string]struct {
		patientsStr     []string
		patients        []*model.Patient
		readErr         error
		resultingErrMsg string
	}{
		"successful": {
			patientsStr: []string{
				fmt.Sprintf("{ \"ID\":\"%s\", \"Name\":\"test1\" }", patientID),
				fmt.Sprintf("{ \"ID\":\"%s\", \"Name\":\"test2\" }", patientID),
				fmt.Sprintf("{ \"ID\":\"%s\", \"Name\":\"test3\" }", patientID),
				fmt.Sprintf("{ \"ID\":\"%s\", \"Name\":\"test4\" }", patientID),
			},
			patients: []*model.Patient{
				{
					ID:   patientID,
					Name: "test1",
				}, {
					ID:   patientID,
					Name: "test2",
				}, {
					ID:   patientID,
					Name: "test3",
				}, {
					ID:   patientID,
					Name: "test4",
				},
			},
			readErr:         nil,
			resultingErrMsg: "",
		},
		"successful with empty database": {
			readErr: os.ErrNotExist,
		},
		"with db error": {
			readErr:         assert.AnError,
			resultingErrMsg: assert.AnError.Error(),
		},
		"with unmarshal error": {
			patientsStr: []string{
				"{ \"ID\":1, \"Name\":\"test1\"",
				"{ \"ID\":2, \"Name\":\"test2\", \"NonExistentField\": \"empty\" }",
			},
			resultingErrMsg: "unexpected end of JSON input",
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		driver := &mockDriver{}
		driver.
			On("ReadAll", patientCollection).
			Return(test.patientsStr, test.readErr).
			Once()

		db := &FileDatabase{
			driver: driver,
		}

		patients, err := db.GetPatients()
		assert.Equal(t, test.patients, patients)
		if test.resultingErrMsg != "" {
			assert.EqualError(t, err, test.resultingErrMsg)
		} else {
			assert.Equal(t, nil, err)
		}
		driver.AssertExpectations(t)
	}
}

func TestFileDatabase_AddPatient(t *testing.T) {
	tests := map[string]struct {
		patient         *model.Patient
		patientID       uuid.UUID
		writePatErr     error
		resultingErrMsg string
	}{
		"successful": {
			patient: &model.Patient{
				Name:     "test",
				ClientID: 1,
			},
		},
		"with write patient error": {
			patient: &model.Patient{
				Name:     "test4",
				ClientID: 2,
			},
			writePatErr:     assert.AnError,
			resultingErrMsg: assert.AnError.Error(),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		driver := &mockDriver{}

		driver.
			On("Write", patientCollection, mock.AnythingOfType("string"), mock.AnythingOfType("*model.Patient")).
			Return(test.writePatErr).
			Once()

		db := &FileDatabase{
			driver: driver,
		}

		patient := new(model.Patient)
		*patient = *test.patient
		err := db.AddPatient(patient)
		if test.resultingErrMsg != "" {
			assert.EqualError(t, err, test.resultingErrMsg)
			assert.Equal(t, patient.ID, uuid.Nil)
		} else {
			assert.Equal(t, nil, err)
			assert.NotEqual(t, patient.ID, uuid.Nil)
		}
		driver.AssertExpectations(t)
	}
}

func TestFileDatabase_UpdatePatient(t *testing.T) {
	tests := map[string]struct {
		patient         *model.Patient
		deleteErr       error
		writeErr        error
		resultingErrMsg string
	}{
		"successful": {
			patient: &model.Patient{
				ID:   uuid.NewV4(),
				Name: "test1",
			},
		},
		"with db error on delete": {
			patient: &model.Patient{
				ID:   uuid.NewV4(),
				Name: "test2",
			},
			deleteErr:       assert.AnError,
			resultingErrMsg: assert.AnError.Error(),
		},
		"with patient or database not found": {
			patient: &model.Patient{
				ID:   uuid.NewV4(),
				Name: "test3",
			},
			deleteErr:       os.ErrNotExist,
			resultingErrMsg: "patient not found",
		},
		"with db error on write": {
			patient: &model.Patient{
				ID:   uuid.NewV4(),
				Name: "test4",
			},
			writeErr:        assert.AnError,
			resultingErrMsg: assert.AnError.Error(),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		driver := &mockDriver{}

		driver.
			On("Delete", patientCollection, test.patient.ID.String()).
			Return(test.deleteErr).
			Once()

		if test.deleteErr == nil {
			driver.
				On("Write", patientCollection, test.patient.ID.String(), test.patient).
				Return(test.writeErr).
				Once()
		}

		db := &FileDatabase{
			driver: driver,
		}

		err := db.UpdatePatient(test.patient)
		if test.deleteErr != nil || test.writeErr != nil {
			assert.EqualError(t, err, test.resultingErrMsg)
		} else {
			assert.Equal(t, nil, err)
		}
		driver.AssertExpectations(t)
	}
}

func TestFileDatabase_RemovePatient(t *testing.T) {
	tests := map[string]struct {
		patient         *model.Patient
		deleteErr       error
		resultingErrMsg string
	}{
		"successful": {
			patient: &model.Patient{
				ID:   uuid.NewV4(),
				Name: "test1",
			},
		},
		"with patient or database not found": {
			patient: &model.Patient{
				ID:   uuid.NewV4(),
				Name: "test3",
			},
			deleteErr:       os.ErrNotExist,
			resultingErrMsg: "patient not found",
		},
		"with db error on delete": {
			patient: &model.Patient{
				ID:   uuid.NewV4(),
				Name: "test2",
			},
			deleteErr:       assert.AnError,
			resultingErrMsg: assert.AnError.Error(),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		driver := &mockDriver{}

		driver.
			On("Delete", patientCollection, test.patient.ID.String()).
			Return(test.deleteErr).
			Once()

		db := &FileDatabase{
			driver: driver,
		}

		err := db.RemovePatient(test.patient)
		if test.deleteErr != nil {
			assert.EqualError(t, err, test.resultingErrMsg)
		} else {
			assert.Equal(t, nil, err)
		}
		driver.AssertExpectations(t)
	}
}
