package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/pagient/pagient-server/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strconv"
	"sync"
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

func TestPatientFileRepository_Get(t *testing.T) {
	patientID := 1

	tests := map[string]struct {
		patientID       int
		patientRet      *model.Patient
		patientRetErr   error
		resultingErrMsg string
	}{
		"successful": {
			patientID: 1,
			patientRet: &model.Patient{
				ID:   patientID,
				Name: "test1",
			},
		},
		"with db error": {
			patientID:       patientID,
			patientRetErr:   assert.AnError,
			resultingErrMsg: "read patient failed: " + assert.AnError.Error(),
		},
		"successful but patient or repository not found": {
			patientID:  patientID,
			patientRet: &model.Patient{},
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		driver := &mockDriver{}
		driver.
			On("Read", patientCollection, strconv.Itoa(test.patientID), mock.AnythingOfType("*model.Patient")).
			Return(test.patientRetErr).
			Run(func(args mock.Arguments) {
				if test.patientRetErr == nil && test.patientRet != nil {
					arg := args.Get(2).(*model.Patient)
					arg.ID = test.patientRet.ID
					arg.Name = test.patientRet.Name
				}
			}).
			Once()

		repository := &patientRepository{
			lock: &sync.Mutex{},
			db:   driver,
		}

		patient, err := repository.Get(test.patientID)
		assert.Equal(t, test.patientRet, patient)
		if test.resultingErrMsg != "" {
			assert.EqualError(t, err, test.resultingErrMsg)
		} else {
			assert.Equal(t, nil, err)
		}
		driver.AssertExpectations(t)
	}
}

func TestPatientFileRepository_GetAll(t *testing.T) {
	tests := map[string]struct {
		patientsStr     []string
		patients        []*model.Patient
		readErr         error
		resultingErrMsg string
	}{
		"successful": {
			patientsStr: []string{
				fmt.Sprintf("{ \"ID\": %d, \"Name\": \"%s\" }", 1, "test1"),
				fmt.Sprintf("{ \"ID\": %d, \"Name\": \"%s\" }", 2, "test2"),
				fmt.Sprintf("{ \"ID\": %d, \"Name\": \"%s\" }", 3, "test3"),
				fmt.Sprintf("{ \"ID\": %d, \"Name\": \"%s\" }", 4, "test4"),
			},
			patients: []*model.Patient{
				{
					ID:   1,
					Name: "test1",
				}, {
					ID:   2,
					Name: "test2",
				}, {
					ID:   3,
					Name: "test3",
				}, {
					ID:   4,
					Name: "test4",
				},
			},
		},
		"successful with empty repository": {
			readErr: os.ErrNotExist,
		},
		"with db error": {
			readErr:         assert.AnError,
			resultingErrMsg: "read patients failed: " + assert.AnError.Error(),
		},
		"with unmarshal error": {
			patientsStr: []string{
				fmt.Sprintf("{ \"ID\": %d, \"Name\": \"%s\"", 1, "test1"),
				fmt.Sprintf("{ \"ID\": %d, \"Name\": \"%s\", \"NonExistentField\": \"empty\" }", 2, "test2"),
			},
			resultingErrMsg: "json unmarshal failed: unexpected end of JSON input",
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		driver := &mockDriver{}
		driver.
			On("ReadAll", patientCollection).
			Return(test.patientsStr, test.readErr).
			Once()

		repository := &patientRepository{
			lock: &sync.Mutex{},
			db:   driver,
		}

		patients, err := repository.GetAll()
		assert.Equal(t, test.patients, patients)
		if test.resultingErrMsg != "" {
			assert.EqualError(t, err, test.resultingErrMsg)
		} else {
			assert.Equal(t, nil, err)
		}
		driver.AssertExpectations(t)
	}
}

func TestPatientFileRepository_Add(t *testing.T) {
	tests := map[string]struct {
		patient         *model.Patient
		patientRet      *model.Patient
		readPatErr      error
		writePatErr     error
		resultingErrMsg string
	}{
		"successful": {
			patient: &model.Patient{
				ID:       1,
				Name:     "test",
				ClientID: 1,
			},
		},
		"with write patient error": {
			patient: &model.Patient{
				ID:       2,
				Name:     "test4",
				ClientID: 2,
			},
			writePatErr:     assert.AnError,
			resultingErrMsg: "write patient failed: " + assert.AnError.Error(),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		driver := &mockDriver{}

		driver.
			On("Read", patientCollection, strconv.Itoa(test.patient.ID), mock.AnythingOfType("*model.Patient")).
			Return(test.readPatErr).
			Run(func(args mock.Arguments) {
				if test.readPatErr == nil && test.patientRet != nil {
					arg := args.Get(2).(*model.Patient)
					arg.ID = test.patientRet.ID
					arg.Name = test.patientRet.Name
				}
			}).
			Once()

		if test.patientRet == nil && test.readPatErr == nil {
			driver.
				On("Write", patientCollection, strconv.Itoa(test.patient.ID), test.patient).
				Return(test.writePatErr).
				Once()
		}

		repository := &patientRepository{
			lock: &sync.Mutex{},
			db:   driver,
		}

		patient := new(model.Patient)
		*patient = *test.patient
		_, err := repository.Add(patient)
		if test.resultingErrMsg != "" {
			assert.EqualError(t, err, test.resultingErrMsg)
		} else {
			assert.Equal(t, nil, err)
		}
		driver.AssertExpectations(t)
	}
}

func TestPatientFileRepository_Update(t *testing.T) {
	tests := map[string]struct {
		patient         *model.Patient
		patientRet      *model.Patient
		readPatErr      error
		writePatErr     error
		resultingErrMsg string
	}{
		"successful": {
			patient: &model.Patient{
				ID:   1,
				Name: "test1",
			},
		},
		"with db error on read": {
			patient: &model.Patient{
				ID:   2,
				Name: "test2",
			},
			readPatErr:      assert.AnError,
			resultingErrMsg: "read patient failed: " + assert.AnError.Error(),
		},
		"with patient or repository not found": {
			patient: &model.Patient{
				ID:   3,
				Name: "test3",
			},
			readPatErr:      os.ErrNotExist,
			resultingErrMsg: "patient not found",
		},
		"with db error on write": {
			patient: &model.Patient{
				ID:   4,
				Name: "test4",
			},
			writePatErr:     assert.AnError,
			resultingErrMsg: "write patient failed: " + assert.AnError.Error(),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		driver := &mockDriver{}

		driver.
			On("Read", patientCollection, strconv.Itoa(test.patient.ID), mock.AnythingOfType("*model.Patient")).
			Return(test.readPatErr).
			Run(func(args mock.Arguments) {
				if test.readPatErr == nil && test.patientRet != nil {
					arg := args.Get(2).(*model.Patient)
					arg.ID = test.patientRet.ID
					arg.Name = test.patientRet.Name
				}
			}).
			Once()

		if test.patientRet == nil && test.readPatErr == nil {
			driver.
				On("Write", patientCollection, strconv.Itoa(test.patient.ID), test.patient).
				Return(test.writePatErr).
				Once()
		}

		repository := &patientRepository{
			lock: &sync.Mutex{},
			db:   driver,
		}

		_, err := repository.Update(test.patient)
		if test.readPatErr != nil || test.writePatErr != nil {
			assert.EqualError(t, err, test.resultingErrMsg)
		} else {
			assert.Equal(t, nil, err)
		}
		driver.AssertExpectations(t)
	}
}

func TestPatientFileRepository_Remove(t *testing.T) {
	tests := map[string]struct {
		patient         *model.Patient
		deleteErr       error
		resultingErrMsg string
	}{
		"successful": {
			patient: &model.Patient{
				ID:   1,
				Name: "test1",
			},
		},
		"with patient or repository not found": {
			patient: &model.Patient{
				ID:   2,
				Name: "test3",
			},
			deleteErr:       os.ErrNotExist,
			resultingErrMsg: "patient not found",
		},
		"with db error on delete": {
			patient: &model.Patient{
				ID:   3,
				Name: "test2",
			},
			deleteErr:       assert.AnError,
			resultingErrMsg: "delete patient failed: " + assert.AnError.Error(),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		driver := &mockDriver{}

		driver.
			On("Delete", patientCollection, strconv.Itoa(test.patient.ID)).
			Return(test.deleteErr).
			Once()

		repository := &patientRepository{
			lock: &sync.Mutex{},
			db:   driver,
		}

		_, err := repository.Remove(test.patient)
		if test.deleteErr != nil {
			assert.EqualError(t, err, test.resultingErrMsg)
		} else {
			assert.Equal(t, nil, err)
		}
		driver.AssertExpectations(t)
	}
}
