package v1

import (
	"finstar-test-task/internal/usecase"
	"finstar-test-task/pkg/logger"
)

type controller struct {
	r usecase.UserRepo
	l logger.Interface
}

func New(r usecase.UserRepo, l logger.Interface) *controller {
	return &controller{r, l}
}
