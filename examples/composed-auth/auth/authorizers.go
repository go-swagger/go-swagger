package auth

import (
	"crypto/rsa"
	"os"

	errors "github.com/go-openapi/errors"
	models "github.com/go-swagger/go-swagger/examples/composed-auth/models"
	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	// currently unused: privateKeyPath = "keys/apiKey.prv"
	publicKeyPath = "keys/apiKey.pem"
	issuerName    = "example.com"
)

var (
	userDb map[string]string

	// Keys used to sign and verify our tokens
	verifyKey *rsa.PublicKey
	// currently unused: signKey   *rsa.PrivateKey
)

// roleClaims describes the format of our JWT token's claims
type roleClaims struct {
	Roles []string `json:"roles"`
	jwt.MapClaims
}

func init() {
	// emulates the loading of a local users database
	userDb = map[string]string{
		"fred": "scrum",
		"ivan": "terrible",
	}

	// loads public keys to verify our tokens
	verifyKeyBuf, err := os.ReadFile(publicKeyPath)
	if err != nil {
		panic("Cannot load public key for tokens")
	}
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyKeyBuf)
	if err != nil {
		panic("Invalid public key for tokens")
	}
}

// Customized authorizer methods for our sample API

// IsRegistered determines if the user is properly registered,
// i.e if a valid username:password pair has been provided
func IsRegistered(user, pass string) (*models.Principal, error) {
	password, ok := userDb[user]
	if !ok || pass != password {
		return nil, errors.New(401, "Unauthorized: not a registered user")
	}

	return &models.Principal{
		Name: user,
	}, nil
}

// IsReseller tells if the API key is a JWT signed by us with a claim to be a reseller
func IsReseller(token string) (*models.Principal, error) {
	claims, err := parseAndCheckToken(token)
	if err != nil {
		return nil, errors.New(401, "Unauthorized: invalid API key token: %v", err)
	}

	issuer, err := claims.GetIssuer()
	if issuer != issuerName {
		return nil, errors.New(401, "Unauthorized: invalid Bearer token: invalid issuer %v", err)
	}

	id, err := claims.GetSubject()
	if err != nil {
		return nil, errors.New(401, "Unauthorized: invalid Bearer token: missing subject: %v", err)
	}

	if issuer != issuerName || id == "" {
		return nil, errors.New(403, "Forbidden: insufficient API key privileges")
	}

	isReseller := false
	for _, role := range claims.Roles {
		if role == "reseller" {
			isReseller = true
			break
		}
	}

	if !isReseller {
		return nil, errors.New(403, "Forbidden: insufficient API key privileges")
	}

	return &models.Principal{
		Name: id,
	}, nil
}

// HasRole tells if the Bearer token is a JWT signed by us with a claim to be
// member of an authorization scope.
// We verify that the claimed role is one of the passed scopes
func HasRole(token string, scopes []string) (*models.Principal, error) {
	claims, err := parseAndCheckToken(token)
	if err != nil {
		return nil, errors.New(401, "Unauthorized: invalid Bearer token: %v", err)
	}
	issuer, err := claims.GetIssuer()
	if err != nil {
		return nil, errors.New(401, "Unauthorized: invalid Bearer token: invalid issuer %v", err)
	}

	if issuer != issuerName {
		return nil, errors.New(401, "Unauthorized: invalid Bearer token: invalid issuer %s", issuer)
	}

	id, err := claims.GetSubject()
	if err != nil {
		return nil, errors.New(401, "Unauthorized: invalid Bearer token: missing subject: %v", err)
	}

	isInScopes := false
	claimedRoles := []string{}
	for _, scope := range scopes {
		for _, role := range claims.Roles {
			if role == scope {
				isInScopes = true
				// we enrich the principal with all claimed roles within scope (hence: not breaking here)
				claimedRoles = append(claimedRoles, role)
			}
		}
	}

	if !isInScopes {
		return nil, errors.New(403, "Forbidden: insufficient privileges")
	}

	return &models.Principal{
		Name:  id,
		Roles: claimedRoles,
	}, nil
}

func parseAndCheckToken(token string) (*roleClaims, error) {
	// the API key is a JWT signed by us with a claim to be a reseller
	parsedToken, err := jwt.ParseWithClaims(token, &roleClaims{}, func(parsedToken *jwt.Token) (interface{}, error) {
		// the key used to validate tokens
		return verifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*roleClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New(401, "invalid token missing expected claims")
	}

	return claims, nil
}
