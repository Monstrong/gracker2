package interceptors

import (
	"context"
	"strconv"

	"github.com/monstrong/gracker2/pkg/contracts/grpcmetakey"
	"github.com/monstrong/gracker2/pkg/domain"
	"github.com/monstrong/gracker2/pkg/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func BlockForBannedAndBillingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		MD, _ := metadata.FromIncomingContext(ctx)
		vals := MD.Get(grpcmetakey.GRPCMetaKeyRole)

		if len(vals) > 0 {
			roleValue, err := strconv.ParseInt(vals[0], 10, 8)
			if err != nil {
				return nil, &errs.GatewayWrongFieldError{Field: grpcmetakey.GRPCMetaKeyRole, Reason: "invalid role format"}
			}
			SenderRole := domain.Role(roleValue)
			if SenderRole < domain.UserRole || SenderRole == domain.BillingRole {
				return nil, errs.ErrPermissionDenied
			}
		}

		return handler(ctx, req)
	}
}
