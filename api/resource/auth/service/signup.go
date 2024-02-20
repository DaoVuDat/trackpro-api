package authservice

import (
	"context"
	"github.com/DaoVuDat/trackpro-api/api/model/project-management/public/model"
	accountdto "github.com/DaoVuDat/trackpro-api/api/resource/account/dto"
	accountrepo "github.com/DaoVuDat/trackpro-api/api/resource/account/repo"
	authdto "github.com/DaoVuDat/trackpro-api/api/resource/auth/dto"
	passworddto "github.com/DaoVuDat/trackpro-api/api/resource/password/dto"
	passwordrepo "github.com/DaoVuDat/trackpro-api/api/resource/password/repo"
	profiledto "github.com/DaoVuDat/trackpro-api/api/resource/profile/dto"
	profilerepo "github.com/DaoVuDat/trackpro-api/api/resource/profile/repo"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/DaoVuDat/trackpro-api/util/password"
	"strings"
)

type SignUpService interface {
	SignUp(app *ctx.Application, authSignUp authdto.AuthSignUp) (accessToken, refreshToken, role, userId string, err error)
}

type signupService struct {
	ctx                context.Context
	accountCreateRepo  accountrepo.CreateAccountRepo
	profileCreateRepo  profilerepo.CreateProfileRepo
	createPasswordRepo passwordrepo.CreatePasswordRepo
	tokenService       TokenService
}

func NewSignUpService(
	ctx context.Context,
	createAccountRepo accountrepo.CreateAccountRepo,
	createProfileRepo profilerepo.CreateProfileRepo,
	createPasswordRepo passwordrepo.CreatePasswordRepo,
	tokenService TokenService,
) SignUpService {
	return &signupService{
		ctx:                ctx,
		accountCreateRepo:  createAccountRepo,
		profileCreateRepo:  createProfileRepo,
		createPasswordRepo: createPasswordRepo,
		tokenService:       tokenService,
	}
}

func (service *signupService) SignUp(app *ctx.Application, authSignUp authdto.AuthSignUp) (accessToken, refreshToken, role, userid string, err error) {

	// create account
	_ = &model.Account{}

	curCtx := service.ctx
	// open transaction for creating account then profile
	tx, err := app.Db.BeginTx(curCtx, nil)
	if err != nil {
		return "", "", "", "", err
	}
	defer tx.Rollback()

	// Create Account
	accountCreate := accountdto.AccountCreate{
		UserName: strings.ToLower(authSignUp.UserName),
	}

	account, err := service.accountCreateRepo.CreateTX(app, curCtx, tx, accountCreate)
	if err != nil {
		app.Logger.Error().Err(err)
		return "", "", "", "", err
	}

	// Create Password
	hashedPassword, err := password.HashPassword(authSignUp.Password)
	if err != nil {
		return "", "", "", "", err
	}

	passwordCreate := passworddto.PasswordCreate{
		UserId:         account.ID,
		HashedPassword: hashedPassword,
	}

	ok, err := service.createPasswordRepo.CreateTX(app, curCtx, tx, passwordCreate)

	if err != nil {
		return "", "", "", "", err
	}

	if !ok {
		return "", "", "", "", common.FailCreateError
	}

	// Create Profile
	profileCreate := profiledto.ProfileCreate{
		UserID:    account.ID,
		FirstName: authSignUp.FirstName,
		LastName:  authSignUp.LastName,
	}

	app.Logger.Debug().Msgf("%v", profileCreate)

	_, err = service.profileCreateRepo.CreateTX(app, curCtx, tx, profileCreate)
	if err != nil {
		return "", "", "", "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", "", "", "", err
	}

	// Access Token
	privateClaims := authdto.PrivateClaimsForToken{
		UserId: account.ID.String(),
		Role:   account.Type.String(),
	}

	accessToken, refreshToken, role, err = service.tokenService.CreateAccessTokenAndRefreshToken(app, curCtx, privateClaims)

	return accessToken, refreshToken, role, account.ID.String(), nil
}
