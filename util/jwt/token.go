package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/rs/zerolog"
	"time"
)

type TokenDetail struct {
	Role   string
	UserId string
}

func CreateToken(logger *zerolog.Logger, userId string, expireDuration time.Duration, role string) (token string, publicKey string, err error) {
	// Generate a new RSA Private Key
	rawPrivateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	realPrivateKey, err := jwk.FromRaw(rawPrivateKey)
	if err != nil {
		logger.Error().Err(fmt.Errorf(`failed to create JWK: %s`, err.Error()))
		return "", "", err
	}

	// Get publicKey and convert to string
	rawPublicKey := &rawPrivateKey.PublicKey
	bytesPublicKey := x509.MarshalPKCS1PublicKey(rawPublicKey)
	// ==> EncodeToMemory returns the PEM encoding
	publicKey = string(pem.EncodeToMemory(
		// ==> A Block represents a PEM encoded structure.
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: bytesPublicKey,
		}),
	)

	// Build jwt claims
	claims, err := jwt.NewBuilder().
		Issuer("trackpro").
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(expireDuration)).
		Claim("userId", userId).
		Claim("role", role).
		Build()

	if err != nil {
		logger.Error().Err(fmt.Errorf(`failed to build claims: %s`, err.Error()))
		return "", "", err
	}

	// Sign the Claims with PrivateKey
	bytesToken, err := jwt.Sign(claims, jwt.WithKey(jwa.RS256, realPrivateKey))
	if err != nil {
		logger.Error().Err(fmt.Errorf(`failed to sign claims: %s`, err.Error()))
		return "", "", err
	}

	token = string(bytesToken)

	return token, publicKey, nil
}

func ParseToken(token string, publicKey string) (TokenDetail, error) {
	panic(1)
}
