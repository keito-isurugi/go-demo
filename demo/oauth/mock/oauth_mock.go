// Code generated by MockGen. DO NOT EDIT.
// Source: oauth.go
//
// Generated by this command:
//
//	mockgen -source=oauth.go -destination=./mock/oauth_mock.go
//

// Package mock_main is a generated GoMock package.
package mock_main

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockOauth is a mock of Oauth interface.
type MockOauth struct {
	ctrl     *gomock.Controller
	recorder *MockOauthMockRecorder
	isgomock struct{}
}

// MockOauthMockRecorder is the mock recorder for MockOauth.
type MockOauthMockRecorder struct {
	mock *MockOauth
}

// NewMockOauth creates a new mock instance.
func NewMockOauth(ctrl *gomock.Controller) *MockOauth {
	mock := &MockOauth{ctrl: ctrl}
	mock.recorder = &MockOauthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOauth) EXPECT() *MockOauthMockRecorder {
	return m.recorder
}

// GenerateCodeChallenge mocks base method.
func (m *MockOauth) GenerateCodeChallenge(verifier string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateCodeChallenge", verifier)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateCodeChallenge indicates an expected call of GenerateCodeChallenge.
func (mr *MockOauthMockRecorder) GenerateCodeChallenge(verifier any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateCodeChallenge", reflect.TypeOf((*MockOauth)(nil).GenerateCodeChallenge), verifier)
}

// GenerateCodeVerifier mocks base method.
func (m *MockOauth) GenerateCodeVerifier() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateCodeVerifier")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateCodeVerifier indicates an expected call of GenerateCodeVerifier.
func (mr *MockOauthMockRecorder) GenerateCodeVerifier() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateCodeVerifier", reflect.TypeOf((*MockOauth)(nil).GenerateCodeVerifier))
}

// GenerateRandomString mocks base method.
func (m *MockOauth) GenerateRandomString(length int, charSet string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateRandomString", length, charSet)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateRandomString indicates an expected call of GenerateRandomString.
func (mr *MockOauthMockRecorder) GenerateRandomString(length, charSet any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRandomString", reflect.TypeOf((*MockOauth)(nil).GenerateRandomString), length, charSet)
}

// GenerateState mocks base method.
func (m *MockOauth) GenerateState() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateState")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateState indicates an expected call of GenerateState.
func (mr *MockOauthMockRecorder) GenerateState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateState", reflect.TypeOf((*MockOauth)(nil).GenerateState))
}
