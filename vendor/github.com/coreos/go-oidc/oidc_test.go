package oidc

import (
	"fmt"
	"testing"
)

const (
	// at_hash value and access_token returned by Google.
	googleAccessTokenHash = "piwt8oCH-K2D9pXlaS1Y-w"
	googleAccessToken     = "ya29.CjHSA1l5WUn8xZ6HanHFzzdHdbXm-14rxnC7JHch9eFIsZkQEGoWzaYG4o7k5f6BnPLj"
	googleSigningAlg      = RS256
	// following values computed by own algo for regression testing
	computed384TokenHash = "_ILKVQjbEzFKNJjUKC2kz9eReYi0A9Of"
	computed512TokenHash = "Spa_APgwBrarSeQbxI-rbragXho6dqFpH5x9PqaPfUI"
)

type accessTokenTest struct {
	name        string
	tok         *IDToken
	accessToken string
	verifier    func(err error) error
}

func (a accessTokenTest) run(t *testing.T) {
	err := a.tok.VerifyAccessToken(a.accessToken)
	result := a.verifier(err)
	if result != nil {
		t.Error(result)
	}
}

func TestAccessTokenVerification(t *testing.T) {
	newToken := func(alg, atHash string) *IDToken {
		return &IDToken{sigAlgorithm: alg, AccessTokenHash: atHash}
	}
	assertNil := func(err error) error {
		if err != nil {
			return fmt.Errorf("want nil error, got %v", err)
		}
		return nil
	}
	assertMsg := func(msg string) func(err error) error {
		return func(err error) error {
			if err == nil {
				return fmt.Errorf("expected error, got success")
			}
			if err.Error() != msg {
				return fmt.Errorf("bad error message, %q, (want %q)", err.Error(), msg)
			}
			return nil
		}
	}
	tests := []accessTokenTest{
		{
			"goodRS256",
			newToken(googleSigningAlg, googleAccessTokenHash),
			googleAccessToken,
			assertNil,
		},
		{
			"goodES384",
			newToken("ES384", computed384TokenHash),
			googleAccessToken,
			assertNil,
		},
		{
			"goodPS512",
			newToken("PS512", computed512TokenHash),
			googleAccessToken,
			assertNil,
		},
		{
			"badRS256",
			newToken("RS256", computed512TokenHash),
			googleAccessToken,
			assertMsg("access token hash does not match value in ID token"),
		},
		{
			"nohash",
			newToken("RS256", ""),
			googleAccessToken,
			assertMsg("id token did not have an access token hash"),
		},
		{
			"badSignAlgo",
			newToken("none", "xxx"),
			googleAccessToken,
			assertMsg(`oidc: unsupported signing algorithm "none"`),
		},
	}
	for _, test := range tests {
		t.Run(test.name, test.run)
	}
}
