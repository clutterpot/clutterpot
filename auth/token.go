package auth

import (
	"github.com/clutterpot/clutterpot/model"
)

func (a *Auth) NewToken(user *model.User) (string, error) {
	claims := map[string]any{
		"id":  user.ID,
		"usr": user.Username,
	}

	_, tokenString, err := a.JWT.Encode(claims)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
