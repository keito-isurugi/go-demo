//go:generate mockgen -source=oauth.go -destination=./mock/oauth_mock.go
package main

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"golang.org/x/oauth2"
)

type Oauth interface {
	GenerateRandomString(length int, charSet string) (string, error)
	GenerateState() (string, error)
	GenerateCodeVerifier() (string, error)
	GenerateCodeChallenge(verifier string) (string, error)
}

type oauth struct {
	readRandomBytes func([]byte) (int, error)
}

func NewOauth(readRandomBytes func([]byte) (int, error)) Oauth {
	return &oauth{
		readRandomBytes: readRandomBytes,
	}
}

const (
	stateLength    = 32
	verifierLength = 128
	charSet        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~"
)

func (o *oauth) GenerateRandomString(length int, charSet string) (string, error) {
	result := make([]byte, length)

	for i := range result {
		b := make([]byte, 1)
		_, err := o.readRandomBytes(b)
		if err != nil {
			return "", err
		}
		result[i] = charSet[b[0]%byte(len(charSet))]
	}

	return string(result), nil
}

func (o *oauth) GenerateState() (string, error) {
	return o.GenerateRandomString(stateLength, charSet)
}

func (o *oauth) GenerateCodeVerifier() (string, error) {
	return o.GenerateRandomString(verifierLength, charSet)
}

func (o *oauth) GenerateCodeChallenge(verifier string) (string, error) {
	if verifier == "" {
		return "", errors.New("code_verifier cannot be empty")
	}

	hash := sha256.Sum256([]byte(verifier))

	codeChallenge := base64.RawURLEncoding.EncodeToString(hash[:])
	return codeChallenge, nil
}

func (o *oauth) GenerateCodeVerifierWithOauth2() string {
	return oauth2.GenerateVerifier()
}

func (o *oauth) GenerateCodeChallengeOauth2(verifier string) string {
	return oauth2.S256ChallengeFromVerifier(verifier)
}
