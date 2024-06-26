// Code generated by MockGen. DO NOT EDIT.
// Source: internal/application/handlers/handler.go

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	reflect "reflect"

	links "github.com/cfif1982/urlshtr.git/internal/domain/links"
	gomock "github.com/golang/mock/gomock"
)

// MockRepositoryInterface is a mock of RepositoryInterface interface.
type MockRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryInterfaceMockRecorder
}

// MockRepositoryInterfaceMockRecorder is the mock recorder for MockRepositoryInterface.
type MockRepositoryInterfaceMockRecorder struct {
	mock *MockRepositoryInterface
}

// NewMockRepositoryInterface creates a new mock instance.
func NewMockRepositoryInterface(ctrl *gomock.Controller) *MockRepositoryInterface {
	mock := &MockRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryInterface) EXPECT() *MockRepositoryInterfaceMockRecorder {
	return m.recorder
}

// AddLink mocks base method.
func (m *MockRepositoryInterface) AddLink(link *links.Link) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLink", link)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddLink indicates an expected call of AddLink.
func (mr *MockRepositoryInterfaceMockRecorder) AddLink(link interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLink", reflect.TypeOf((*MockRepositoryInterface)(nil).AddLink), link)
}

// GetLinkByKey mocks base method.
func (m *MockRepositoryInterface) GetLinkByKey(key string) (*links.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinkByKey", key)
	ret0, _ := ret[0].(*links.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinkByKey indicates an expected call of GetLinkByKey.
func (mr *MockRepositoryInterfaceMockRecorder) GetLinkByKey(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinkByKey", reflect.TypeOf((*MockRepositoryInterface)(nil).GetLinkByKey), key)
}
