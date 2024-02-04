package authservice

import (
	"context"
	"trackpro/api/model/project-management/public/model"
	accountdto "trackpro/api/resource/account/dto"
	accountrepo "trackpro/api/resource/account/repo"
	authdto "trackpro/api/resource/auth/dto"
	profiledto "trackpro/api/resource/profile/dto"
	profilerepo "trackpro/api/resource/profile/repo"
	"trackpro/util/ctx"
	"trackpro/util/jwt"
)

type SignUpService interface {
	SignUp(app *ctx.Application, authSignUp authdto.AuthSignUp) (string, error)
}

type signupService struct {
	ctx               context.Context
	accountCreateRepo accountrepo.CreateAccountRepo
	profileCreateRepo profilerepo.CreateProfileRepo
}

func NewSignUpService(
	ctx context.Context,
	accountCreateRepo accountrepo.CreateAccountRepo,
	profileCreateRepo profilerepo.CreateProfileRepo,

) SignUpService {
	return &signupService{
		ctx:               ctx,
		accountCreateRepo: accountCreateRepo,
		profileCreateRepo: profileCreateRepo,
	}
}

func (service *signupService) SignUp(app *ctx.Application, authSignUp authdto.AuthSignUp) (string, error) {
	// create account
	_ = &model.Account{}

	curCtx := service.ctx
	// open transaction for creating account then profile
	tx, err := app.Db.BeginTx(curCtx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	accountCreate := accountdto.AccountCreate{
		UserName: authSignUp.UserName,
	}

	account, err := service.accountCreateRepo.CreateTX(app, curCtx, tx, accountCreate)
	if err != nil {
		return "", err
	}

	profileCreate := profiledto.ProfileCreate{
		UserID:    account.ID,
		FirstName: authSignUp.FirstName,
		LastName:  authSignUp.LastName,
	}

	app.Logger.Debug().Msgf("%v", profileCreate)

	_, err = service.profileCreateRepo.CreateTX(app, curCtx, tx, profileCreate)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	// Access Token
	tokenDetail, err := jwt.CreateToken(
		app.Logger,
		account.ID.String(),
		account.Type.String(),
		app.Config.AccessTokenPrivateKey,
		app.Config.RefreshTokenExpiredIn,
	)
	if err != nil {
		return "", err
	}

	// Store Public Key into Redis to cache
	return *tokenDetail.Token, nil
}
