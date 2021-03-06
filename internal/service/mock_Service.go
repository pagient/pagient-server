// Code generated by mockery v1.0.0. DO NOT EDIT.

package service

import mock "github.com/stretchr/testify/mock"
import model "github.com/pagient/pagient-server/internal/model"

// MockService is an autogenerated mock type for the Service type
type MockService struct {
	mock.Mock
}

// CallPatient provides a mock function with given fields: _a0
func (_m *MockService) CallPatient(_a0 *model.Patient) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Patient) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ChangeUserPassword provides a mock function with given fields: _a0
func (_m *MockService) ChangeUserPassword(_a0 *model.User) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.User) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateClient provides a mock function with given fields: _a0
func (_m *MockService) CreateClient(_a0 *model.Client) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Client) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreatePatient provides a mock function with given fields: _a0
func (_m *MockService) CreatePatient(_a0 *model.Patient) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Patient) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateToken provides a mock function with given fields: _a0
func (_m *MockService) CreateToken(_a0 *model.Token) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Token) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateUser provides a mock function with given fields: _a0
func (_m *MockService) CreateUser(_a0 *model.User) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.User) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeletePatient provides a mock function with given fields: _a0
func (_m *MockService) DeletePatient(_a0 *model.Patient) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Patient) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteToken provides a mock function with given fields: _a0
func (_m *MockService) DeleteToken(_a0 *model.Token) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Token) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListClients provides a mock function with given fields:
func (_m *MockService) ListClients() ([]*model.Client, error) {
	ret := _m.Called()

	var r0 []*model.Client
	if rf, ok := ret.Get(0).(func() []*model.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPagerPatientsByStatus provides a mock function with given fields: _a0
func (_m *MockService) ListPagerPatientsByStatus(_a0 ...model.PatientStatus) ([]*model.Patient, error) {
	_va := make([]interface{}, len(_a0))
	for _i := range _a0 {
		_va[_i] = _a0[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*model.Patient
	if rf, ok := ret.Get(0).(func(...model.PatientStatus) []*model.Patient); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Patient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...model.PatientStatus) error); ok {
		r1 = rf(_a0...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPagers provides a mock function with given fields:
func (_m *MockService) ListPagers() ([]*model.Pager, error) {
	ret := _m.Called()

	var r0 []*model.Pager
	if rf, ok := ret.Get(0).(func() []*model.Pager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Pager)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPatients provides a mock function with given fields:
func (_m *MockService) ListPatients() ([]*model.Patient, error) {
	ret := _m.Called()

	var r0 []*model.Patient
	if rf, ok := ret.Get(0).(func() []*model.Patient); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Patient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListTokensByUser provides a mock function with given fields: _a0
func (_m *MockService) ListTokensByUser(_a0 string) ([]*model.Token, error) {
	ret := _m.Called(_a0)

	var r0 []*model.Token
	if rf, ok := ret.Get(0).(func(string) []*model.Token); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUsers provides a mock function with given fields:
func (_m *MockService) ListUsers() ([]*model.User, error) {
	ret := _m.Called()

	var r0 []*model.User
	if rf, ok := ret.Get(0).(func() []*model.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Login provides a mock function with given fields: _a0, _a1
func (_m *MockService) Login(_a0 string, _a1 string) (*model.User, bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string, string) *model.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string, string) bool); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Get(1).(bool)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ShowClient provides a mock function with given fields: _a0
func (_m *MockService) ShowClient(_a0 uint) (*model.Client, error) {
	ret := _m.Called(_a0)

	var r0 *model.Client
	if rf, ok := ret.Get(0).(func(uint) *model.Client); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ShowClientByUser provides a mock function with given fields: _a0
func (_m *MockService) ShowClientByUser(_a0 string) (*model.Client, error) {
	ret := _m.Called(_a0)

	var r0 *model.Client
	if rf, ok := ret.Get(0).(func(string) *model.Client); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Client)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ShowPager provides a mock function with given fields: _a0
func (_m *MockService) ShowPager(_a0 uint) (*model.Pager, error) {
	ret := _m.Called(_a0)

	var r0 *model.Pager
	if rf, ok := ret.Get(0).(func(uint) *model.Pager); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Pager)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ShowPatient provides a mock function with given fields: _a0
func (_m *MockService) ShowPatient(_a0 uint) (*model.Patient, error) {
	ret := _m.Called(_a0)

	var r0 *model.Patient
	if rf, ok := ret.Get(0).(func(uint) *model.Patient); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Patient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ShowToken provides a mock function with given fields: _a0
func (_m *MockService) ShowToken(_a0 string) (*model.Token, error) {
	ret := _m.Called(_a0)

	var r0 *model.Token
	if rf, ok := ret.Get(0).(func(string) *model.Token); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ShowUser provides a mock function with given fields: _a0
func (_m *MockService) ShowUser(_a0 string) (*model.User, error) {
	ret := _m.Called(_a0)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string) *model.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ShowUserByToken provides a mock function with given fields: _a0
func (_m *MockService) ShowUserByToken(_a0 string) (*model.User, error) {
	ret := _m.Called(_a0)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string) *model.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePatient provides a mock function with given fields: _a0
func (_m *MockService) UpdatePatient(_a0 *model.Patient) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Patient) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
