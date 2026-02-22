package errno

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestErrNo_Error(t *testing.T) {
	err := NewErrNo(201001, "用户不存在")

	expected := "ErrorCode: 201001, ErrorMsg: 用户不存在"
	assert.Equal(t, expected, err.Error())
}

func TestNewErrNo(t *testing.T) {
	err := NewErrNo(200100, "参数错误")

	assert.Equal(t, int32(200100), err.ErrCode)
	assert.Equal(t, "参数错误", err.ErrMsg)
}

func TestErrNo_WithMessage(t *testing.T) {
	originalErr := NewErrNo(201001, "用户不存在")
	newErr := originalErr.WithMessage("找不到指定用户")

	assert.Equal(t, int32(201001), newErr.ErrCode)
	assert.Equal(t, "找不到指定用户", newErr.ErrMsg)

	// 原始错误不应被修改
	assert.Equal(t, "用户不存在", originalErr.ErrMsg)
}

func TestErrNo_Code(t *testing.T) {
	err := NewErrNo(201002, "用户已存在")

	assert.Equal(t, int32(201002), err.Code())
}

func TestErrNo_Message(t *testing.T) {
	err := NewErrNo(201002, "用户已存在")

	assert.Equal(t, "用户已存在", err.Message())
}

func TestToKitexError(t *testing.T) {
	t.Run("converts ErrNo to BizStatusError", func(t *testing.T) {
		errNo := NewErrNo(201001, "用户不存在")
		kitexErr := ToKitexError(errNo)

		assert.NotNil(t, kitexErr)
		assert.Contains(t, kitexErr.Error(), "201001")
		assert.Contains(t, kitexErr.Error(), "用户不存在")
	})

	t.Run("returns nil for nil error", func(t *testing.T) {
		result := ToKitexError(nil)

		assert.Nil(t, result)
	})

	t.Run("wraps non-ErrNo errors", func(t *testing.T) {
		normalErr := errors.New("普通错误")
		kitexErr := ToKitexError(normalErr)

		assert.NotNil(t, kitexErr)
		assert.Contains(t, kitexErr.Error(), "Operation failed")
	})
}

func TestIsRecordNotFound(t *testing.T) {
	t.Run("returns true for GORM ErrRecordNotFound", func(t *testing.T) {
		err := gorm.ErrRecordNotFound
		assert.True(t, IsRecordNotFound(err))
	})

	t.Run("returns false for other errors", func(t *testing.T) {
		err := errors.New("其他错误")
		assert.False(t, IsRecordNotFound(err))
	})

	t.Run("returns false for nil error", func(t *testing.T) {
		assert.False(t, IsRecordNotFound(nil))
	})

	t.Run("returns false for wrapped GORM error", func(t *testing.T) {
		err := errors.New("wrapped error")
		assert.False(t, IsRecordNotFound(err))
	})
}

func TestWrapDatabaseError(t *testing.T) {
	t.Run("wraps error with message", func(t *testing.T) {
		originalErr := errors.New("数据库连接失败")
		wrappedErr := WrapDatabaseError(originalErr, "创建用户失败")

		assert.Equal(t, int32(ErrorCodeOperationFailed), wrappedErr.Code())
		assert.Contains(t, wrappedErr.Message(), "创建用户失败")
		assert.Contains(t, wrappedErr.Message(), "数据库连接失败")
	})

	t.Run("returns empty error for nil input", func(t *testing.T) {
		result := WrapDatabaseError(nil, "操作失败")

		assert.Equal(t, int32(0), result.Code())
		assert.Empty(t, result.Message())
	})
}

func TestErrorCodes(t *testing.T) {
	t.Run("user related error codes", func(t *testing.T) {
		assert.Equal(t, int32(201001), int32(ErrorCodeUserNotFound))
		assert.Equal(t, int32(201002), int32(ErrorCodeUserAlreadyExists))
		assert.Equal(t, int32(201003), int32(ErrorCodeUsernameAlreadyExists))
		assert.Equal(t, int32(201004), int32(ErrorCodeEmailAlreadyExists))
		assert.Equal(t, int32(201016), int32(ErrorCodeInvalidCredentials))
		assert.Equal(t, int32(201018), int32(ErrorCodeMustChangePassword))
	})

	t.Run("organization related error codes", func(t *testing.T) {
		assert.Equal(t, int32(202001), int32(ErrorCodeOrganizationNotFound))
		assert.Equal(t, int32(202002), int32(ErrorCodeOrganizationAlreadyExists))
	})

	t.Run("department related error codes", func(t *testing.T) {
		assert.Equal(t, int32(203001), int32(ErrorCodeDepartmentNotFound))
		assert.Equal(t, int32(203002), int32(ErrorCodeDepartmentAlreadyExists))
	})

	t.Run("role related error codes", func(t *testing.T) {
		assert.Equal(t, int32(207001), int32(ErrorCodeRoleDefinitionNotFound))
		assert.Equal(t, int32(207002), int32(ErrorCodeRoleNameAlreadyExists))
	})
}
