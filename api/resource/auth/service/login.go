package authservice

import (
	authdto "trackpro/api/resource/auth/dto"
	"trackpro/util/ctx"
	"trackpro/util/jwt"
)

type LoginService interface {
	Login(app *ctx.Application, authLogin authdto.AuthLoginTemp) (*jwt.TokenDetail, error)
}

type loginService struct {
}

func NewLoginService() LoginService {
	return &loginService{}
}

func (service *loginService) Login(app *ctx.Application, authLogin authdto.AuthLoginTemp) (*jwt.TokenDetail, error) {
	// Get Public Key from Redis

	// Access Token
	tokenDetail, err := jwt.ParseToken(app.Logger, authLogin.Token, app.Config.AccessTokenPublicKey)
	if err != nil {
		return nil, err
	}

	return tokenDetail, nil
}
