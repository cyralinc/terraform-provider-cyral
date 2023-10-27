package core

type OperationType string

const (
	OperationTypeCreate = OperationType("create")
	OperationTypeRead   = OperationType("read")
	OperationTypeUpdate = OperationType("update")
	OperationTypeDelete = OperationType("delete")
)
