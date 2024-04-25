// Code generated by MockGen. DO NOT EDIT.
// Source: app.go
//
// Generated by this command:
//
//	mockgen -source=app.go -destination=mock_dependencies_test.go -package=application_test
//

// Package application_test is a generated GoMock package.
package application_test

import (
	context "context"
	reflect "reflect"

	domain "github.com/alcortesm/demo-mongodb-transactions/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockGroupRepo is a mock of GroupRepo interface.
type MockGroupRepo struct {
	ctrl     *gomock.Controller
	recorder *MockGroupRepoMockRecorder
}

// MockGroupRepoMockRecorder is the mock recorder for MockGroupRepo.
type MockGroupRepoMockRecorder struct {
	mock *MockGroupRepo
}

// NewMockGroupRepo creates a new mock instance.
func NewMockGroupRepo(ctrl *gomock.Controller) *MockGroupRepo {
	mock := &MockGroupRepo{ctrl: ctrl}
	mock.recorder = &MockGroupRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGroupRepo) EXPECT() *MockGroupRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockGroupRepo) Create(ctx context.Context, group *domain.Group) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, group)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockGroupRepoMockRecorder) Create(ctx, group any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockGroupRepo)(nil).Create), ctx, group)
}

// Load mocks base method.
func (m *MockGroupRepo) Load(ctx context.Context, id string) (*domain.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", ctx, id)
	ret0, _ := ret[0].(*domain.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load.
func (mr *MockGroupRepoMockRecorder) Load(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockGroupRepo)(nil).Load), ctx, id)
}

// Update mocks base method.
func (m *MockGroupRepo) Update(ctx context.Context, group *domain.Group) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, group)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockGroupRepoMockRecorder) Update(ctx, group any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockGroupRepo)(nil).Update), ctx, group)
}

// MockUuider is a mock of Uuider interface.
type MockUuider struct {
	ctrl     *gomock.Controller
	recorder *MockUuiderMockRecorder
}

// MockUuiderMockRecorder is the mock recorder for MockUuider.
type MockUuiderMockRecorder struct {
	mock *MockUuider
}

// NewMockUuider creates a new mock instance.
func NewMockUuider(ctrl *gomock.Controller) *MockUuider {
	mock := &MockUuider{ctrl: ctrl}
	mock.recorder = &MockUuiderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUuider) EXPECT() *MockUuiderMockRecorder {
	return m.recorder
}

// NewString mocks base method.
func (m *MockUuider) NewString() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewString")
	ret0, _ := ret[0].(string)
	return ret0
}

// NewString indicates an expected call of NewString.
func (mr *MockUuiderMockRecorder) NewString() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewString", reflect.TypeOf((*MockUuider)(nil).NewString))
}