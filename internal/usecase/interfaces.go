// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"finstar-test-task/proto"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	Controller interface {
		IncreaseBalance(proto.IncreaseBalanceRequest) (*proto.IncreaseBalanceResponse, error)
		TransferBalance(proto.TransferBalanceRequest) (*proto.TransferBalanceResponse, error)
	}

	// UserRepo -.
	UserRepo interface {
		IncreaseBalance(userId uint64, receipt float64) error
		TransferBalance(userIdFrom, userIdTo uint64, writeOff float64) error
	}
)
