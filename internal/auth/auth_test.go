package auth

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Sanity tests.
func TestHashPassword(t *testing.T) {
	var hashTests = []struct {
		password string
	}{
		{
			"",
		},
		{
			"antidisestablishmenterianism",
		},
	}

	for _, tt := range hashTests {
		t.Run(tt.password, func(t *testing.T) {
			result, err := HashPassword(tt.password)
			assert.NoError(t, err)

			passwordMatchesHash := CheckPasswordHash(tt.password, result)

			assert.True(t, passwordMatchesHash)
		})
	}
}

func TestValidateJWT(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		secret := "test_secret_key"
		userID, _ := rand.Int(rand.Reader, big.NewInt(64))
		convertedUserId := int32(userID.Int64())

		tokenString, err := MakeJWT(convertedUserId, secret)
		if err != nil {
			t.Errorf("Failed to create JWT: %v", err)
		}

		decodedID, err := ValidateJWT(tokenString, secret)
		if err != nil || convertedUserId != int32(decodedID) {
			t.Errorf("Expected %v, got %v (error: %v)", userID, decodedID, err)
		}
	})

	t.Run("Invalid secret", func(t *testing.T) {
		secret := "correct_test_key"
		userID, _ := rand.Int(rand.Reader, big.NewInt(64))
		convertedUserId := int32(userID.Int64())

		tokenString, err := MakeJWT(convertedUserId, secret)
		if err != nil {
			t.Errorf("failed to create jwt: %v", err)
		}

		_, err = ValidateJWT(tokenString, "fakeSecret")
		if err == nil {
			t.Errorf("expected error due to incorrect secret")
		}
	})
}
