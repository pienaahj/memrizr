package service

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pienaahj/memrizr/account/model"
	"github.com/pienaahj/memrizr/account/model/mocks"
)

func TestNewPairFromUser(t *testing.T) {
	var idExp int64 = 15 * 60
	var refreshExp int64 = 3 * 24 * 2600

	mockTokenRepository := new(mocks.MockTokenRepository)

	priv, _ := os.ReadFile("../rsa_private_test.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(priv)
	pub, _ := os.ReadFile("../rsa_public_test.pem")
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(pub)
	secret := "anotsorandomtestsecret"

	// instantiate a common token service to be used by all tests
	tokenService := NewTokenService(&TSConfig{
		TokenRepository:       mockTokenRepository,
		PrivKey:               privKey,
		PubKey:                pubKey,
		RefreshSecret:         secret,
		IDExpirationSecs:      idExp,
		RefreshExpirationSecs: refreshExp,
	})

	// include password to make sure it is not serialized
	// since json tag is "-"
	uid, _ := uuid.NewRandom()
	u := &model.User{
		UID:      uid,
		Email:    "bob@bob.com",
		Password: "blarghedymcblarghface",
	}

	// Setup mock call responses in setup before t.Run statements
	uidErrorCase, _ := uuid.NewRandom()
	uErrorCase := &model.User{
		UID:      uidErrorCase,
		Email:    "failure@failure.com",
		Password: "blarghedymcblarghface",
	}
	prevID := "a_previous_tokenID"

	setSuccessArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		u.UID.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	setErrorArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		uErrorCase.UID.String(),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
	}

	deleteWithPrevIDArguments := mock.Arguments{
		mock.AnythingOfType("*context.emptyCtx"),
		u.UID.String(),
		prevID,
	}

	// mock call argument/responses
	mockTokenRepository.On("SetRefreshToken", setSuccessArguments...).Return(nil)
	mockTokenRepository.On("SetRefreshToken", setErrorArguments...).Return(fmt.Errorf("Error setting refresh token"))
	mockTokenRepository.On("DeleteRefreshToken", deleteWithPrevIDArguments...).Return(nil)
	t.Run("Returns a token pair with values", func(t *testing.T) {
		ctx := context.Background()                                    // updated from context.TODO()
		tokenPair, err := tokenService.NewPairFromUser(ctx, u, prevID) // replaced "" with prevID from setup
		assert.NoError(t, err)

		// SetRefreshToken should be called with setSuccessArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since prevID is ""
		mockTokenRepository.AssertCalled(t, "DeleteRefreshToken", deleteWithPrevIDArguments...)

		var s string
		assert.IsType(t, s, tokenPair.IDToken)

		// decode the Base64URL encoded string
		// simpler to use jwt library which is already imported
		idTokenClaims := &IDTokenCustomClaims{}

		_, err = jwt.ParseWithClaims(tokenPair.IDToken, idTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return pubKey, nil
		})

		assert.NoError(t, err)

		// assert claims on idToken
		expectedClaims := []interface{}{
			u.UID,
			u.Email,
			u.Name,
			u.ImageURL,
			u.Website,
		}
		actualIDClaims := []interface{}{
			idTokenClaims.User.UID,
			idTokenClaims.User.Email,
			idTokenClaims.User.Name,
			idTokenClaims.User.ImageURL,
			idTokenClaims.User.Website,
		}

		assert.ElementsMatch(t, expectedClaims, actualIDClaims)
		assert.Empty(t, idTokenClaims.User.Password) // password should never be encoded to json

		/*
			// unixTime := time.Now().Unix() - type int64
			unixTime := jwt.NewNumericDate(time.Now()) // type jwt.NumericDate - type struct{ time.Time}
			// tokenExp := unixTime + 60*15 // 15 minutes from current time
			// tokenExp := jwt.NewNumericDate(time.Now().Add(15 * time.Minute)) // 15 minutes from current time
			// tokenExp := unixTime + exp
			tokenExp := jwt.NewNumericDate(time.Now().Add(time.Duration(exp) * time.Second)) // 15 minutes from current time



		*/

		// for idToken
		// expiresAt := time.Unix(idTokenClaims.RegisteredClaims.ExpiresAt.Time.Unix(), 0)
		expiresAt := idTokenClaims.RegisteredClaims.ExpiresAt.Time
		// fmt.Println("Expires at:", expiresAt)
		expectedExpiresAt := time.Now().Add(time.Duration(idExp) * time.Second)
		// fmt.Println("Should Expire at:", expectedExpiresAt)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)

		refreshTokenClaims := &RefreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.RefreshToken, refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		assert.IsType(t, s, tokenPair.RefreshToken)

		// assert claims on refresh token
		assert.NoError(t, err)
		assert.Equal(t, u.UID, refreshTokenClaims.UID)

		expiresAt = time.Unix(refreshTokenClaims.RegisteredClaims.ExpiresAt.Time.Unix(), 0)
		expectedExpiresAt = time.Now().Add(time.Duration(refreshExp) * time.Second)
		assert.WithinDuration(t, expectedExpiresAt, expiresAt, 5*time.Second)
	})

	t.Run("Error setting refresh token", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, uErrorCase, "")
		assert.Error(t, err) // should return an error

		// SetRefreshToken should be called with setErrorArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setErrorArguments...)
		// DeleteRefreshToken should not be since SetRefreshToken causes method to return
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})

	t.Run("Empty string provided for prevID", func(t *testing.T) {
		ctx := context.Background()
		_, err := tokenService.NewPairFromUser(ctx, u, "")
		assert.NoError(t, err)

		// SetRefreshToken should be called with setSuccessArguments
		mockTokenRepository.AssertCalled(t, "SetRefreshToken", setSuccessArguments...)
		// DeleteRefreshToken should not be called since prevID is ""
		mockTokenRepository.AssertNotCalled(t, "DeleteRefreshToken")
	})
}
