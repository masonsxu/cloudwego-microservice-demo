package errno

import (
	"errors"
	"fmt"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"gorm.io/gorm"
)

// ErrNo 错误码结构定义
type ErrNo struct {
	ErrCode int32
	ErrMsg  string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("ErrorCode: %d, ErrorMsg: %s", e.ErrCode, e.ErrMsg)
}

func NewErrNo(code int32, msg string) ErrNo {
	return ErrNo{
		ErrCode: code,
		ErrMsg:  msg,
	}
}

// WithMessage 快速设置
func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

// Code / Message 快速获取，避免直接访问字段
func (e ErrNo) Code() int32     { return e.ErrCode }
func (e ErrNo) Message() string { return e.ErrMsg }

// ToKitexError 转换为Kitex错误
func ToKitexError(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(ErrNo); ok {
		// 使用Kitex的NewBizStatusError创建业务错误
		return kerrors.NewBizStatusError(e.Code(), e.Message())
	}

	// 其他类型的错误包装为通用操作失败错误
	return kerrors.NewBizStatusError(
		int32(ErrorCodeOperationFailed),
		"Operation failed",
	)
}

// IsResearcherNotFound 检查是否为研究人员未找到错误
func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// WrapDatabaseError 包装数据库错误为业务错误
func WrapDatabaseError(err error, message string) ErrNo {
	if err == nil {
		return ErrNo{}
	}

	return NewErrNo(ErrorCodeOperationFailed, fmt.Sprintf("%s: %v", message, err))
}
