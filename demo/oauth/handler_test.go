package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	mockMain "oauth/mock"
	"testing"
)

func TestHandler_GenerateCodeChallenges(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOauth := mockMain.NewMockOauth(ctrl)
	h := NewHandler(mockOauth)

	testState := "test_state"
	testCodeVerifier := "test_code_verifier"
	testCodeChallenge := "test_code_challenge"

	tests := []struct {
		name      string
		mockCall  func()
		want      *GenerateCodeChallengesResponse
		wantError error
	}{
		{
			name: "正常系/GenerateCodeChallengesResponseが正しく返される",
			mockCall: func() {
				mockOauth.EXPECT().GenerateState().Return(testState, nil)
				mockOauth.EXPECT().GenerateCodeVerifier().Return(testCodeVerifier, nil)
				mockOauth.EXPECT().GenerateCodeChallenge(gomock.Any()).Return(testCodeChallenge, nil)
			},
			want: &GenerateCodeChallengesResponse{
				State:         testState,
				CodeVerifier:  testCodeVerifier,
				CodeChallenge: testCodeChallenge,
			},
		},
		{
			name: "異常系/GenerateStateがエラー",
			mockCall: func() {
				mockOauth.EXPECT().GenerateState().Return("", errors.New("GenerateState error"))
			},
			want:      nil,
			wantError: errors.New("GenerateState error"),
		},
		{
			name: "異常系/GenerateCodeVerifierがエラー",
			mockCall: func() {
				mockOauth.EXPECT().GenerateState().Return(testState, nil)
				mockOauth.EXPECT().GenerateCodeVerifier().Return("", errors.New("GenerateCodeVerifier error"))
			},
			want:      nil,
			wantError: errors.New("GenerateCodeVerifier error"),
		},
		{
			name: "異常系/GenerateCodeChallengeがエラー",
			mockCall: func() {
				mockOauth.EXPECT().GenerateState().Return(testState, nil)
				mockOauth.EXPECT().GenerateCodeVerifier().Return(testCodeVerifier, nil)
				mockOauth.EXPECT().GenerateCodeChallenge(gomock.Any()).Return("", errors.New("GenerateCodeChallenge error"))
			},
			want:      nil,
			wantError: errors.New("GenerateCodeChallenge error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			tt.mockCall()

			res, err := h.GenerateCodeChallenges()
			if tt.wantError != nil {
				a.Error(err)
				a.Equal(tt.wantError.Error(), err.Error())
			} else {
				a.NoError(err)
				a.Equal(tt.want, res)
			}
		})
	}
}
