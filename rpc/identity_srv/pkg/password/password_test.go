package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	t.Run("successfully hashes password", func(t *testing.T) {
		password := "MySecurePassword123!"
		hash, err := HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash)
		assert.Contains(t, hash, "$2a$") // bcrypt hash prefix
	})

	t.Run("generates different hashes for same password", func(t *testing.T) {
		password := "SamePassword123!"
		hash1, err1 := HashPassword(password)
		hash2, err2 := HashPassword(password)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2) // Different salt
	})

	t.Run("handles empty password", func(t *testing.T) {
		password := ""
		hash, err := HashPassword(password)

		require.NoError(t, err)
		assert.NotEmpty(t, hash)
	})
}

func TestVerifyPassword(t *testing.T) {
	t.Run("verifies correct password", func(t *testing.T) {
		password := "CorrectPassword123!"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		isValid := VerifyPassword(password, hash)
		assert.True(t, isValid)
	})

	t.Run("rejects incorrect password", func(t *testing.T) {
		password := "CorrectPassword123!"
		wrongPassword := "WrongPassword123!"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		isValid := VerifyPassword(wrongPassword, hash)
		assert.False(t, isValid)
	})

	t.Run("rejects empty password against non-empty hash", func(t *testing.T) {
		password := "SomePassword123!"
		hash, err := HashPassword(password)
		require.NoError(t, err)

		isValid := VerifyPassword("", hash)
		assert.False(t, isValid)
	})

	t.Run("handles empty hash", func(t *testing.T) {
		password := "SomePassword123!"
		isValid := VerifyPassword(password, "")
		assert.False(t, isValid)
	})

	t.Run("handles invalid hash format", func(t *testing.T) {
		password := "SomePassword123!"
		invalidHash := "not-a-valid-bcrypt-hash"

		isValid := VerifyPassword(password, invalidHash)
		assert.False(t, isValid)
	})
}

func TestPasswordHashingIntegration(t *testing.T) {
	t.Run("common password patterns", func(t *testing.T) {
		passwords := []string{
			"simple",
			"Password123",
			"Complex!@#123",
			"中文密码测试",
			"very long password with multiple words and numbers 123456789",
		}

		for _, password := range passwords {
			hash, err := HashPassword(password)
			require.NoError(t, err, "Failed to hash password: %s", password)
			assert.True(t, VerifyPassword(password, hash), "Failed to verify password: %s", password)
		}
	})
}

func BenchmarkHashPassword(b *testing.B) {
	password := "BenchmarkPassword123!"

	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "BenchmarkPassword123!"
	hash, _ := HashPassword(password)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		VerifyPassword(password, hash)
	}
}
