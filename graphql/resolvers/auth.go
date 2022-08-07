package resolvers

import (
	"context"

	"github.com/clutterpot/clutterpot/auth"
	"github.com/clutterpot/clutterpot/model"
)

func (r *mutationResolver) Login(ctx context.Context, email, password string) (*model.AuthPayload, error) {
	user, err := r.Store.User.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if err := model.CompareHashAndPassword(user.Password, password); err != nil {
		return nil, err
	}

	accessToken, accessTokenString, err := r.Auth.NewAccessToken(&auth.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Kind:     user.Kind,
	})
	if err != nil {
		return nil, err
	}

	session, err := r.Store.Session.Create(model.SessionInput{UserID: user.ID})
	if err != nil {
		return nil, err
	}

	_, refreshTokenString, err := r.Auth.NewRefreshToken(session)
	if err != nil {
		return nil, err
	}

	return &model.AuthPayload{
		AccessToken:  accessTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    accessToken.Expiration(),
		RefreshToken: &refreshTokenString,
	}, nil
}

func (r *mutationResolver) RefreshAccessToken(ctx context.Context, refreshToken string) (*model.AuthPayload, error) {
	_, sessionClaims, err := r.Auth.Decode(refreshToken)
	if err != nil {
		return nil, err
	}

	sessionUser, err := r.Store.Session.GetByID(sessionClaims["sid"].(string))
	if err != nil {
		return nil, err
	}

	accessToken, accessTokenString, err := r.Auth.NewAccessToken(&auth.Claims{
		UserID:   sessionUser.User.ID,
		Username: sessionUser.User.Username,
		Kind:     sessionUser.User.Kind,
	})
	if err != nil {
		return nil, err
	}

	return &model.AuthPayload{
		AccessToken: accessTokenString,
		TokenType:   "Bearer",
		ExpiresIn:   accessToken.Expiration(),
	}, nil
}
