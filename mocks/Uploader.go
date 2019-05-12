// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import io "io"
import mock "github.com/stretchr/testify/mock"

import uuid "github.com/gofrs/uuid"

// Uploader is an autogenerated mock type for the Uploader type
type Uploader struct {
	mock.Mock
}

// Store provides a mock function with given fields: src
func (_m *Uploader) Store(src io.Reader) (uuid.UUID, error) {
	ret := _m.Called(src)

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(io.Reader) uuid.UUID); ok {
		r0 = rf(src)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(io.Reader) error); ok {
		r1 = rf(src)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}