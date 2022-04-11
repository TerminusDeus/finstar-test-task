package v1

import (
	"finstar-test-task/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *controller) IncreaseBalance(request proto.IncreaseBalanceRequest) (*proto.IncreaseBalanceResponse, error) {
	if err := request.Validate(); err != nil {
		c.l.Error(err, "http - v1 - IncreaseBalance")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := c.r.IncreaseBalance(request.GetUserId(), request.GetReceipt()); err != nil {
		c.l.Error(err, "http - v1 - IncreaseBalance")

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.IncreaseBalanceResponse{}, nil
}
