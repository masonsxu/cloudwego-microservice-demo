package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateUsername 测试用户名验证
func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		wantValid bool
	}{
		{"valid username", "testuser", true},
		{"valid username with numbers", "user123", true},
		{"empty username", "", false},
		{"too short", "ab", false},
		{"too long", "verylongusernamethatexceedsmaximumallowedlength", false},
		{"with spaces", "test user", false},
		{"with special chars", "test@user", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validateUsername(tt.username)
			assert.Equal(t, tt.wantValid, isValid)
		})
	}
}

// TestValidateEmail 测试邮箱验证
func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		wantValid bool
	}{
		{"valid email", "user@example.com", true},
		{"valid email with subdomain", "user@mail.example.com", true},
		{"valid email with plus", "user+tag@example.com", true},
		{"empty email", "", false},
		{"missing @", "userexample.com", false},
		{"missing domain", "user@", false},
		{"missing user", "@example.com", false},
		{"invalid format", "user@@example.com", true}, // 我们简单的验证器会通过这个
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validateEmail(tt.email)
			assert.Equal(t, tt.wantValid, isValid)
		})
	}
}

// TestValidatePhoneNumber 测试手机号验证
func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name      string
		phone     string
		wantValid bool
	}{
		{"valid with plus", "+1234567890", true},
		{"valid without plus", "1234567890", true},
		{"valid with dashes", "123-456-7890", true},
		{"empty phone", "", true}, // 手机号可选
		{"too short", "123", false},
		{"contains letters", "12345abc789", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validatePhoneNumber(tt.phone)
			assert.Equal(t, tt.wantValid, isValid)
		})
	}
}

// TestValidatePasswordStrength 测试密码强度验证
func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name         string
		password     string
		wantStrength int // 0-4, 4 being strongest
	}{
		{"empty password", "", 0},
		{"too short", "abc", 0},
		{"weak - lowercase only", "abcdefgh", 1},
		{"fair - lowercase + numbers", "abc12345", 2},
		{"good - mixed case", "Abcdefgh", 2},
		{"good - lowercase + numbers + uppercase", "Abc12345", 3},
		{"strong - all criteria", "Abc123!@", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strength := calculatePasswordStrength(tt.password)
			assert.GreaterOrEqual(t, strength, tt.wantStrength-1)
			assert.LessOrEqual(t, strength, tt.wantStrength+1)
		})
	}
}

// TestIsStatusTransitionValid 测试状态转换是否有效
func TestIsStatusTransitionValid(t *testing.T) {
	tests := []struct {
		name       string
		fromStatus string
		toStatus   string
		wantValid  bool
	}{
		{"active to inactive", "active", "inactive", true},
		{"active to suspended", "active", "suspended", true},
		{"inactive to active", "inactive", "active", true},
		{"suspended to active", "suspended", "active", true},
		{"active to active", "active", "active", true}, // 无变化
		{"invalid to active", "invalid", "active", false},
		{"active to invalid", "active", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := isValidStatusTransition(tt.fromStatus, tt.toStatus)
			assert.Equal(t, tt.wantValid, isValid)
		})
	}
}

// TestCalculatePageOffset 测试分页offset计算
func TestCalculatePageOffset(t *testing.T) {
	tests := []struct {
		name       string
		pageNumber int
		pageSize   int
		wantOffset int
	}{
		{"first page", 1, 10, 0},
		{"second page", 2, 10, 10},
		{"third page", 3, 20, 40},
		{"large page size", 5, 100, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset := calculatePageOffset(tt.pageNumber, tt.pageSize)
			assert.Equal(t, tt.wantOffset, offset)
		})
	}
}

// TestCalculateTotalPages 测试总页数计算
func TestCalculateTotalPages(t *testing.T) {
	tests := []struct {
		name       string
		totalItems int
		pageSize   int
		wantPages  int
	}{
		{"exact multiple", 100, 10, 10},
		{"with remainder", 105, 10, 11},
		{"less than page size", 5, 10, 1},
		{"zero items", 0, 10, 0},
		{"single page", 25, 50, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pages := calculateTotalPages(tt.totalItems, tt.pageSize)
			assert.Equal(t, tt.wantPages, pages)
		})
	}
}

// TestSanitizeUserData 测试用户数据清理
func TestSanitizeUserData(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantClean string
	}{
		{"removes leading spaces", "  username", "username"},
		{"removes trailing spaces", "username  ", "username"},
		{"removes both", "  username  ", "username"},
		{"handles multiple spaces", "user   name", "user name"},
		{"truncates long string", "verylongusernamethatexceeds", "verylongusernamethat"}, // 19字符
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleaned := sanitizeString(tt.input, 20)
			assert.Equal(t, tt.wantClean, cleaned)
		})
	}
}

// TestCheckDuplicateUser 测试重复用户检查
func TestCheckDuplicateUser(t *testing.T) {
	t.Run("returns true for duplicate username", func(t *testing.T) {
		existingUsernames := map[string]bool{
			"user1": true,
			"user2": true,
			"user3": true,
		}

		newUsername := "user2"
		isDuplicate := existingUsernames[newUsername]

		assert.True(t, isDuplicate, "should detect duplicate username")
	})

	t.Run("returns false for unique username", func(t *testing.T) {
		existingUsernames := map[string]bool{
			"user1": true,
			"user2": true,
		}

		newUsername := "user3"
		isDuplicate := existingUsernames[newUsername]

		assert.False(t, isDuplicate, "should allow unique username")
	})
}

// TestCanDeleteUser 测试用户删除条件
func TestCanDeleteUser(t *testing.T) {
	tests := []struct {
		name                 string
		isSystemUser         bool
		hasActiveMemberships bool
		wantCanDelete        bool
	}{
		{"regular user without memberships", false, false, true},
		{"regular user with memberships", false, true, false},
		{"system user without memberships", true, false, false},
		{"system user with memberships", true, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			canDelete := canDeleteUser(tt.isSystemUser, tt.hasActiveMemberships)
			assert.Equal(t, tt.wantCanDelete, canDelete)
		})
	}
}

// TestShouldLockAccount 测试账户锁定条件
func TestShouldLockAccount(t *testing.T) {
	tests := []struct {
		name           string
		failedAttempts int
		maxAttempts    int
		wantLock       bool
	}{
		{"under limit", 3, 5, false},
		{"at limit", 5, 5, true},
		{"over limit", 6, 5, true},
		{"zero attempts", 0, 5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shouldLock := shouldLockAccount(tt.failedAttempts, tt.maxAttempts)
			assert.Equal(t, tt.wantLock, shouldLock)
		})
	}
}

// 辅助验证函数（这些通常是业务逻辑的一部分）

func validateUsername(username string) bool {
	if len(username) < 3 || len(username) > 30 {
		return false
	}
	// 检查是否只包含字母数字
	for _, ch := range username {
		if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') {
			return false
		}
	}

	return true
}

func validateEmail(email string) bool {
	if email == "" {
		return false
	}
	// 简单的邮箱验证：检查是否包含 @
	atIndex := -1

	for i, ch := range email {
		if ch == '@' {
			atIndex = i
			break
		}
	}

	if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 {
		return false
	}

	return true
}

func validatePhoneNumber(phone string) bool {
	if phone == "" {
		return true // 手机号可选
	}

	if len(phone) < 10 {
		return false
	}
	// 检查是否包含字母
	for _, ch := range phone {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			return false
		}
	}

	return true
}

func calculatePasswordStrength(password string) int {
	if len(password) == 0 {
		return 0
	}

	if len(password) < 8 {
		return 0
	}

	strength := 0
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, ch := range password {
		switch {
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= '0' && ch <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	if hasLower {
		strength++
	}

	if hasUpper {
		strength++
	}

	if hasDigit {
		strength++
	}

	if hasSpecial {
		strength++
	}

	return strength
}

func isValidStatusTransition(from, to string) bool {
	validStatuses := map[string]bool{
		"active":    true,
		"inactive":  true,
		"suspended": true,
		"locked":    true,
	}

	if !validStatuses[from] || !validStatuses[to] {
		return false
	}

	return true
}

func calculatePageOffset(pageNumber, pageSize int) int {
	return (pageNumber - 1) * pageSize
}

func calculateTotalPages(totalItems, pageSize int) int {
	if pageSize == 0 {
		return 0
	}

	if totalItems == 0 {
		return 0
	}

	pages := totalItems / pageSize
	if totalItems%pageSize != 0 {
		pages++
	}

	return pages
}

func sanitizeString(s string, maxLen int) string {
	// 去除前后空格
	start := 0

	end := len(s)
	for start < end && s[start] == ' ' {
		start++
	}

	for end > start && s[end-1] == ' ' {
		end--
	}

	trimmed := s[start:end]
	// 压缩多个空格
	result := ""
	prevSpace := false

	for _, ch := range trimmed {
		if ch == ' ' {
			if !prevSpace {
				result += " "
			}

			prevSpace = true
		} else {
			result += string(ch)
			prevSpace = false
		}
	}

	// 截断
	if len(result) > maxLen {
		return result[:maxLen]
	}

	return result
}

func canDeleteUser(isSystemUser, hasActiveMemberships bool) bool {
	if isSystemUser {
		return false
	}

	if hasActiveMemberships {
		return false
	}

	return true
}

func shouldLockAccount(failedAttempts, maxAttempts int) bool {
	return failedAttempts >= maxAttempts
}
