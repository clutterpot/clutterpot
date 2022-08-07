package auth

import (
	"time"

	"github.com/clutterpot/clutterpot/model"
	"github.com/lestrrat-go/jwx/jwt"
)

func (a *Auth) NewAccessToken(privateClaims *Claims) (jwt.Token, string, error) {
	now := time.Now().UTC()
	claims := map[string]any{
		jwt.IssuedAtKey:   now,
		jwt.ExpirationKey: now.Add(time.Minute * 10),
		"uid":             privateClaims.UserID,
		"usr":             privateClaims.Username,
		"knd":             privateClaims.Kind,
	}

	return a.jwt.Encode(claims)
}

func (a *Auth) NewRefreshToken(session *model.Session) (jwt.Token, string, error) {
	claims := map[string]any{
		jwt.IssuedAtKey:   session.CreatedAt,
		jwt.ExpirationKey: session.ExpiresAt,
		"sid":             session.ID,
	}

	return a.jwt.Encode(claims)
}
