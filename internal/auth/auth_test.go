package auth

import (
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
