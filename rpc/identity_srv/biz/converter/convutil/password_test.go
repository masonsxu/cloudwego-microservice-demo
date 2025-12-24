package convutil

import (
	"testing"
)

func TestHashPassword_Success(t *testing.T) {
	password := "mySecretPassword123"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	if hash == "" {
		t.Error("HashPassword() returned empty hash")
	}

	if hash == password {
		t.Error("HashPassword() returned plaintext password")
	}
}

func TestHashPassword_DifferentHashesForSamePassword(t *testing.T) {
	password := "mySecretPassword123"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// bcrypt 每次生成不同的哈希（因为有随机盐）
	if hash1 == hash2 {
		t.Error("HashPassword() should generate different hashes due to random salt")
	}
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("")

	if err != nil {
		t.Fatalf("HashPassword(\"\") error = %v", err)
	}

	// 空密码也应该能被哈希
	if hash == "" {
		t.Error("HashPassword(\"\") should return non-empty hash")
	}
}

func TestVerifyPassword_Correct(t *testing.T) {
	password := "mySecretPassword123"
	hash, _ := HashPassword(password)

	if !VerifyPassword(password, hash) {
		t.Error("VerifyPassword() = false for correct password")
	}
}

func TestVerifyPassword_Incorrect(t *testing.T) {
	password := "mySecretPassword123"
	wrongPassword := "wrongPassword"
	hash, _ := HashPassword(password)

	if VerifyPassword(wrongPassword, hash) {
		t.Error("VerifyPassword() = true for incorrect password")
	}
}

func TestVerifyPassword_EmptyPassword(t *testing.T) {
	hash, _ := HashPassword("")

	if !VerifyPassword("", hash) {
		t.Error("VerifyPassword() = false for correct empty password")
	}

	if VerifyPassword("anyPassword", hash) {
		t.Error("VerifyPassword() = true for incorrect password against empty hash")
	}
}

func TestVerifyPassword_InvalidHash(t *testing.T) {
	if VerifyPassword("password", "invalid-hash") {
		t.Error("VerifyPassword() = true for invalid hash format")
	}
}

func TestVerifyPassword_EmptyHash(t *testing.T) {
	if VerifyPassword("password", "") {
		t.Error("VerifyPassword() = true for empty hash")
	}
}

func TestHashPassword_SpecialCharacters(t *testing.T) {
	passwords := []string{
		"P@ssw0rd!#$%",
		"密码123",
		"パスワード",
		"emoji🔐password",
		"  spaces  ",
		"newline\npassword",
	}

	for _, password := range passwords {
		hash, err := HashPassword(password)
		if err != nil {
			t.Errorf("HashPassword(%q) error = %v", password, err)
			continue
		}

		if !VerifyPassword(password, hash) {
			t.Errorf("VerifyPassword(%q) = false after hashing", password)
		}
	}
}

func TestHashPassword_LongPassword(t *testing.T) {
	// bcrypt 有 72 字节的限制，超过会返回错误
	longPassword := ""
	for i := 0; i < 100; i++ {
		longPassword += "a"
	}

	_, err := HashPassword(longPassword)

	// golang.org/x/crypto/bcrypt 对超过 72 字节的密码返回错误
	if err == nil {
		t.Error("HashPassword(>72 bytes) should return error")
	}
}

func TestHashPassword_MaxLengthPassword(t *testing.T) {
	// 测试 72 字节边界（bcrypt 最大支持长度）
	maxPassword := ""
	for i := 0; i < 72; i++ {
		maxPassword += "a"
	}

	hash, err := HashPassword(maxPassword)
	if err != nil {
		t.Fatalf("HashPassword(72 bytes) error = %v", err)
	}

	if !VerifyPassword(maxPassword, hash) {
		t.Error("VerifyPassword(72 bytes) = false")
	}
}

// Benchmark 测试
func BenchmarkHashPassword(b *testing.B) {
	password := "benchmarkPassword123"
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(password)
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "benchmarkPassword123"
	hash, _ := HashPassword(password)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = VerifyPassword(password, hash)
	}
}
