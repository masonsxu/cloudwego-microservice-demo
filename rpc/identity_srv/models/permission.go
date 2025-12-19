package models

// Permission 定义了一个具体的操作权限，由“资源+动作”构成，并可附加约束条件。
// 这个结构体用于数据库持久化，与 aip/rpc 层的 thrift 定义解耦。
type Permission struct {
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description,omitempty"`
}
