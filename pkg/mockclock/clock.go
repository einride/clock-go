// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/einride/clock-go/pkg/clock (interfaces: Clock,Ticker)

// Package mockclock is a generated GoMock package.
package mockclock

import (
	clock "github.com/einride/clock-go/pkg/clock"
	gomock "github.com/golang/mock/gomock"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	time "time"
)

// MockClock is a mock of Clock interface
type MockClock struct {
	ctrl     *gomock.Controller
	recorder *MockClockMockRecorder
}

// MockClockMockRecorder is the mock recorder for MockClock
type MockClockMockRecorder struct {
	mock *MockClock
}

// NewMockClock creates a new mock instance
func NewMockClock(ctrl *gomock.Controller) *MockClock {
	mock := &MockClock{ctrl: ctrl}
	mock.recorder = &MockClockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClock) EXPECT() *MockClockMockRecorder {
	return m.recorder
}

// After mocks base method
func (m *MockClock) After(arg0 time.Duration) <-chan time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "After", arg0)
	ret0, _ := ret[0].(<-chan time.Time)
	return ret0
}

// After indicates an expected call of After
func (mr *MockClockMockRecorder) After(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "After", reflect.TypeOf((*MockClock)(nil).After), arg0)
}

// NewTicker mocks base method
func (m *MockClock) NewTicker(arg0 time.Duration) clock.Ticker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewTicker", arg0)
	ret0, _ := ret[0].(clock.Ticker)
	return ret0
}

// NewTicker indicates an expected call of NewTicker
func (mr *MockClockMockRecorder) NewTicker(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewTicker", reflect.TypeOf((*MockClock)(nil).NewTicker), arg0)
}

// Now mocks base method
func (m *MockClock) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now
func (mr *MockClockMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockClock)(nil).Now))
}

// NowProto mocks base method
func (m *MockClock) NowProto() *timestamppb.Timestamp {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NowProto")
	ret0, _ := ret[0].(*timestamppb.Timestamp)
	return ret0
}

// NowProto indicates an expected call of NowProto
func (mr *MockClockMockRecorder) NowProto() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NowProto", reflect.TypeOf((*MockClock)(nil).NowProto))
}

// Since mocks base method
func (m *MockClock) Since(arg0 time.Time) time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Since", arg0)
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// Since indicates an expected call of Since
func (mr *MockClockMockRecorder) Since(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Since", reflect.TypeOf((*MockClock)(nil).Since), arg0)
}

// MockTicker is a mock of Ticker interface
type MockTicker struct {
	ctrl     *gomock.Controller
	recorder *MockTickerMockRecorder
}

// MockTickerMockRecorder is the mock recorder for MockTicker
type MockTickerMockRecorder struct {
	mock *MockTicker
}

// NewMockTicker creates a new mock instance
func NewMockTicker(ctrl *gomock.Controller) *MockTicker {
	mock := &MockTicker{ctrl: ctrl}
	mock.recorder = &MockTickerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTicker) EXPECT() *MockTickerMockRecorder {
	return m.recorder
}

// C mocks base method
func (m *MockTicker) C() <-chan time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "C")
	ret0, _ := ret[0].(<-chan time.Time)
	return ret0
}

// C indicates an expected call of C
func (mr *MockTickerMockRecorder) C() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "C", reflect.TypeOf((*MockTicker)(nil).C))
}

// Stop mocks base method
func (m *MockTicker) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop
func (mr *MockTickerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockTicker)(nil).Stop))
}
