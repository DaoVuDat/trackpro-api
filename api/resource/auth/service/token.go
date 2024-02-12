package authservice

import (
	"context"
	"errors"
	"fmt"
	authdto "github.com/DaoVuDat/trackpro-api/api/resource/auth/dto"
	"github.com/DaoVuDat/trackpro-api/api/router/common"
	"github.com/DaoVuDat/trackpro-api/util/ctx"
	"github.com/DaoVuDat/trackpro-api/util/jwt"
	"github.com/redis/go-redis/v9"
	"strings"
)

type TokenService interface {
	CreateAccessTokenAndRefreshToken(app *ctx.Application, curCtx context.Context, privateClaims authdto.PrivateClaimsForToken) (accessToken string, refreshToken string, role string, err error)
	GetToken(app *ctx.Application, curCtx context.Context, tokenDetail authdto.TokenDetailForRedis) (string, error)
	RevokeToken(app *ctx.Application, curCtx context.Context, revokeAll bool, tokenDetail authdto.TokenDetailForRedis) (bool, error)
	RefreshAccessToken(app *ctx.Application, curCtx context.Context,
		privateClaims authdto.PrivateClaimsForToken,
		oldAccessTokenDetail authdto.TokenDetailForRedis,
		refreshTokenDetail authdto.TokenDetailForRedis) (accessToken string, refreshToken string, role string, err error)
}

type tokenService struct {
	redisClient *redis.Client
}

func NewTokenService(client *redis.Client) TokenService {
	return &tokenService{
		redisClient: client,
	}
}

func (service *tokenService) CreateAccessTokenAndRefreshToken(app *ctx.Application, curCtx context.Context, privateClaims authdto.PrivateClaimsForToken) (accessToken string, refreshToken string, role string, err error) {
	accessTokenDetail, err := jwt.CreateToken(
		app.Logger,
		privateClaims.UserId,
		privateClaims.Role,
		app.Config.AccessTokenPrivateKey,
		app.Config.AccessTokenExpiredIn,
	)
	if err != nil {
		return "", "", role, err
	}

	// Refresh Token
	refreshTokenDetail, err := jwt.CreateToken(
		app.Logger,
		privateClaims.UserId,
		privateClaims.Role,
		app.Config.RefreshTokenPrivateKey,
		app.Config.RefreshTokenExpiredIn,
	)
	if err != nil {
		return "", "", role, err
	}

	// Store to redis
	accessKeyRedis := fmt.Sprintf("%s:%s:%s", privateClaims.UserId, authdto.Access, *accessTokenDetail.Token)
	refreshKeyRedis := fmt.Sprintf("%s:%s:%s", privateClaims.UserId, authdto.Refresh, *refreshTokenDetail.Token)

	if err = service.redisClient.Set(curCtx, accessKeyRedis, authdto.Active, app.Config.AccessTokenExpiredIn).Err(); err != nil {
		return "", "", "", errors.New("redis caching failure")
	}

	if err = service.redisClient.Set(curCtx, refreshKeyRedis, *accessTokenDetail.Token, app.Config.RefreshTokenExpiredIn).Err(); err != nil {
		return "", "", "", errors.New("redis caching failure")
	}

	return *accessTokenDetail.Token, *refreshTokenDetail.Token, privateClaims.Role, nil
}

func (service *tokenService) GetToken(app *ctx.Application, curCtx context.Context, tokenDetail authdto.TokenDetailForRedis) (string, error) {
	key := fmt.Sprintf("%s:%s:%s", tokenDetail.UserId, tokenDetail.TokenType, tokenDetail.Token)
	if value, err := service.redisClient.Get(curCtx, key).Result(); err != nil {
		if errors.Is(err, redis.Nil) {
			return "", common.TokenExpired
		}
		return "", err
	} else {
		return value, nil
	}
}

func (service *tokenService) RevokeToken(app *ctx.Application, curCtx context.Context, revokeAll bool, tokenDetail authdto.TokenDetailForRedis) (bool, error) {
	if !revokeAll {
		key := fmt.Sprintf("%s:%s:%s", tokenDetail.UserId, tokenDetail.TokenType, tokenDetail.Token)

		if err := service.redisClient.Del(curCtx, key).Err(); err != nil {
			return false, errors.New("redis remove cache failure")
		}
	} else {
		allKeys := fmt.Sprintf("%s:*", tokenDetail.UserId)
		var startCursor uint64

		_, err := service.redisClient.Pipelined(curCtx, func(pipe redis.Pipeliner) error {

			iter := service.redisClient.Scan(curCtx, startCursor, allKeys, 0).Iterator()
			for iter.Next(curCtx) {
				key := iter.Val()
				pipe.Del(curCtx, key)
			}
			return nil
		})

		if err != nil {
			return false, errors.New("redis remove all cache failure")
		}
	}
	return true, nil
}

func (service *tokenService) RefreshAccessToken(
	app *ctx.Application,
	curCtx context.Context,
	privateClaims authdto.PrivateClaimsForToken,
	oldAccessTokenDetail authdto.TokenDetailForRedis,
	refreshTokenDetail authdto.TokenDetailForRedis,
) (accessToken string, refreshToken string, role string, err error) {
	// check refresh token in
	accessTokenValueOfRefreshToken, err := service.GetToken(app, curCtx, refreshTokenDetail)
	if err != nil {

		return "", "", "", err
	}

	// Check oldAccessToken equal accessTokenValueOfRefreshToken
	if strings.Compare(oldAccessTokenDetail.Token, accessTokenValueOfRefreshToken) != 0 {
		// warning -> revoke all sessions
		_, err := service.RevokeToken(app, curCtx, true, authdto.TokenDetailForRedis{
			UserId:    refreshTokenDetail.UserId,
			TokenType: "",
			Token:     "",
		})

		if err != nil {
			return "", "", "", err
		}

		return "", "", "", common.OldAccessTokenAndRefreshTokenNotMatch
	}

	// create new pair token and remove the old refresh token, only use one time
	accessToken, refreshToken, role, err = service.CreateAccessTokenAndRefreshToken(app, curCtx, privateClaims)
	if err != nil {
		return "", "", "", err
	}

	// revoke current refreshToken
	_, err = service.RevokeToken(app, curCtx, false, refreshTokenDetail)
	if err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, role, nil
}
