package authservice

import (
	"context"
	accountrepo "github.com/DaoVuDat/trackpro-api/api/resource/account/repo"
	authdto "github.com/DaoVuDat/trackpro-api/api/resource/auth/dto"
	passwordrepo "github.com/DaoVuDat/trackpro-api/api/resource/password/repo"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/DaoVuDat/trackpro-api/util/password"
	"strings"
)

type LoginService interface {
	Login(app *ctx.Application, curCtx context.Context, authLogin authdto.AuthLogin) (accessToken, refreshToken, role string, err error)
}

type loginService struct {
	findAccountRepo  accountrepo.FindAccountRepo
	findPasswordRepo passwordrepo.FindPasswordRepo
	tokenService     TokenService
}

func NewLoginService(
	findAccountRepo accountrepo.FindAccountRepo,
	findPasswordRepo passwordrepo.FindPasswordRepo,
	tokenService TokenService,
) LoginService {
	return &loginService{
		findAccountRepo:  findAccountRepo,
		findPasswordRepo: findPasswordRepo,
		tokenService:     tokenService,
	}
}

func (service *loginService) Login(app *ctx.Application, curCtx context.Context, authLogin authdto.AuthLogin) (
	accessToken, refreshToken, role string,
	err error,
) {
	// find account
	account, err := service.findAccountRepo.FindByUserName(app, strings.ToLower(authLogin.UserName))
	if err != nil {
		return "", "", "", err
	}

	// get password
	pwd, err := service.findPasswordRepo.Find(app, account.ID)
	if err != nil {
		return "", "", "", err
	}

	// compare password
	valid, err := password.ComparedPassword(pwd, authLogin.Password)
	if err != nil {
		return "", "", "", err
	}

	if !valid {
		return "", "", "", err
	}

	// generate pair token
	accessToken, refreshToken, role, err = service.tokenService.CreateAccessTokenAndRefreshToken(app, curCtx, authdto.PrivateClaimsForToken{
		UserId: account.ID.String(),
		Role:   account.Type.String(),
	})
	if err != nil {
		return "", "", "", err
	}
	return accessToken, refreshToken, role, nil
}
