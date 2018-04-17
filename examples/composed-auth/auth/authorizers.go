package auth

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/davecgh/go-spew/spew"

	jwt "github.com/dgrijalva/jwt-go"
	errors "github.com/go-openapi/errors"
	models "github.com/go-swagger/go-swagger/examples/composed-auth/models"
	logging "github.com/op/go-logging"
)

const (
	privateKeyPath = "keys/apiKey.prv"
	publicKeyPath  = "keys/apiKey.pem"
)

var (
	logger *logging.Logger
	userDb map[string]string

	// Keys used to sign and verify our tokens
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

// roleClaims describes the format of our JWT token's claims
type roleClaims struct {
	Roles []string `json:"roles"`
	jwt.StandardClaims
}

func init() {
	logger = logging.MustGetLogger("auth")
	logging.SetLevel(logging.DEBUG, "auth")

	// emulates the loading of a local users database
	userDb = map[string]string{
		"fred": "scrum",
		"ivan": "terrible",
	}

	// loads public keys to verify our tokens
	verifyKeyBuf, err := ioutil.ReadFile(publicKeyPath)
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
	logger.Debugf("Credentials: %q:%q", user, pass)
	if password, ok := userDb[user]; ok {
		if pass == password {
			return &models.Principal{
				Name: user,
			}, nil
		}
	}
	logger.Debug("Bad credentials")
	return nil, errors.New(401, "Unauthorized: not a registered user")
}

// IsReseller tells if the API key is a JWT signed by us with a claim to be a reseller
func IsReseller(token string) (*models.Principal, error) {
	logger.Debug("Parsing token")
	claims, err := parseAndCheckToken(token)
	if err == nil {
		logger.Debugf("Token claims: %s", spew.Sdump(claims))
		if claims.Issuer == "example.com" && claims.Id != "" {
			isReseller := false
			for _, role := range claims.Roles {
				if role == "reseller" {
					isReseller = true
					break
				}
			}
			if isReseller {
				return &models.Principal{
					Name: claims.Id,
				}, nil
			}
			logger.Debug("Bad claims")
			return nil, errors.New(403, "Forbidden: insufficient API key privileges")
		}
	}
	logger.Debug("Bad credentials")
	return nil, errors.New(401, "Unauthorized: invalid API key token: %v", err)
}

// HasRole tells if the Bearer token is a JWT signed by us with a claim to be
// member of an authorization scope.
// We verify that the claimed role is one of the passed scopes
func HasRole(token string, scopes []string) (*models.Principal, error) {
	logger.Debugf("Runtime passed scopes: %v", scopes)
	claims, err := parseAndCheckToken(token)
	if err == nil {
		logger.Debugf("Token claims: %s", spew.Sdump(claims))
		if claims.Issuer == "example.com" {
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
			if isInScopes {
				return &models.Principal{
					Name:  claims.Id,
					Roles: claimedRoles,
				}, nil
			}
			logger.Debug("Bad claims")
			return nil, errors.New(403, "Forbidden: insufficient privileges")
		}
	}
	logger.Debug("Bad credentials")
	return nil, errors.New(401, "Unauthorized: invalid Bearer token: %v", err)
}

func parseAndCheckToken(token string) (*roleClaims, error) {
	// the API key is a JWT signed by us with a claim to be a reseller
	parsedToken, err := jwt.ParseWithClaims(token, &roleClaims{}, func(parsedToken *jwt.Token) (interface{}, error) {
		// the key used to validate tokens
		return verifyKey, nil
	})

	if err == nil {
		if claims, ok := parsedToken.Claims.(*roleClaims); ok && parsedToken.Valid {
			return claims, nil
		}
	}
	return nil, err

}
