package jwt

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog"
	"time"
)

type TokenDetail struct {
	Role    string
	UserId  string
	Token   *string
	TokenId *string
}

func CreateToken(logger *zerolog.Logger, userId string, role string, privateKey string, expireDuration time.Duration) (tokenDetail *TokenDetail, err error) {
	// Decode Private Key
	decodePrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		logger.Error().Err(fmt.Errorf(`failed to decode PrivateKey: %s`, err.Error()))
		return nil, err
	}

	realPrivateKey, err := jwk.ParseKey(decodePrivateKey, jwk.WithPEM(true))
	if err != nil {
		logger.Error().Err(fmt.Errorf(`failed to parse Key: %s`, err.Error()))
		return nil, err
	}

	tokenDetail = &TokenDetail{
		Role:    role,
		UserId:  userId,
		Token:   nil,
		TokenId: nil,
	}

	tokenId := uuid.NewString()
	tokenDetail.TokenId = &tokenId

	// Build jwt claims
	claims, err := jwt.NewBuilder().
		JwtID(tokenId).
		Issuer("trackpro").
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(expireDuration)).
		Claim("userId", userId).
		Claim("role", role).
		Build()

	if err != nil {
		logger.Error().Err(fmt.Errorf(`failed to build claims: %s`, err.Error()))
		return nil, err
	}

	// Sign the Claims with PrivateKey
	bytesToken, err := jwt.Sign(claims, jwt.WithKey(jwa.RS256, realPrivateKey))
	if err != nil {
		logger.Error().Err(fmt.Errorf(`failed to sign claims: %s`, err.Error()))
		return nil, err
	}

	token := string(bytesToken)
	tokenDetail.Token = &token

	return tokenDetail, nil
}

func ParseToken(logger *zerolog.Logger, token string, publicKey string) (*TokenDetail, error) {
	// Decode Public Key
	decodePublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		logger.Error().Err(fmt.Errorf(`failed to decode PublicKey: %s`, err.Error()))
		return nil, err
	}

	realPubKey, _ := jwk.ParseKey(decodePublicKey, jwk.WithPEM(true))
	if err != nil {
		logger.Error().Err(fmt.Errorf(`failed to parse Key: %s`, err.Error()))
		return nil, err
	}

	// Parse token for claims
	claims, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.RS256, realPubKey), jwt.WithValidate(false))
	if err != nil {
		fmt.Printf("jwt.Parse failed: %s\n", err)
		logger.Error().Err(fmt.Errorf(`failed to parse token to claims: %s`, err.Error()))
		return nil, err
	}

	// Create TokenDetail
	var role string
	var userId string
	if roleInf, ok := claims.Get("role"); !ok {
		err = fmt.Errorf(`failed to get private claim "role"`)
		logger.Error().Err(err)
		return nil, err
	} else {
		if role, ok = roleInf.(string); !ok {
			err = fmt.Errorf(`"role" expected to be string, but got %T`, roleInf)
			logger.Error().Err(err)
			return nil, err
		}
	}

	if userIdInterface, ok := claims.Get("userId"); !ok {
		err = fmt.Errorf(`failed to get private claim "userid"`)
		logger.Error().Err(err)
		return nil, err
	} else {
		if userId, ok = userIdInterface.(string); !ok {
			err = fmt.Errorf(`"userId" expected to be string, but got %T`, userIdInterface)
			logger.Error().Err(err)
			return nil, err
		}
	}

	tokenDetail := &TokenDetail{
		Role:   role,
		UserId: userId,
	}
	return tokenDetail, nil
}
