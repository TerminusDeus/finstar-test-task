package v1

import (
	"finstar-test-task/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *controller) TransferBalance(request proto.TransferBalanceRequest) (*proto.TransferBalanceResponse, error) {
	if err := request.Validate(); err != nil {
		c.l.Error(err, "http - v1 - TransferBalance")

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := c.r.TransferBalance(request.GetUserIdFrom(), request.GetUserIdTo(), request.GetWriteOff()); err != nil {
		c.l.Error(err.Error(), "http - v1 - TransferBalance")

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.TransferBalanceResponse{}, nil
}
