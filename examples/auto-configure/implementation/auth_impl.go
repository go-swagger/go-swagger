package implementation

import "errors"

type AuthImpl struct{}

func (i *AuthImpl) KeyAuth(token string) (any, error) {
	if token != "example token" {
		return nil, errors.New("wrong token")
	}

	// if return nil, nil, will cause 401 error
	return true, nil
}
