package implementation

import "errors"

type AuthImpl struct{}

func (i *AuthImpl) KeyAuth(token string) (interface{}, error) {
	if token != "example token" {
		return nil, errors.New("Wrong token")
	}
	return nil, nil
}
