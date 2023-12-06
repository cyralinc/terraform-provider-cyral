package operationtype

type OperationType string

const (
	Create = OperationType("create")
	Read   = OperationType("read")
	Update = OperationType("update")
	Delete = OperationType("delete")
)
