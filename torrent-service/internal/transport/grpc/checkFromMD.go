package grpc

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/pkg/contracts/grpcmetakey"
	"github.com/monstrong/gracker2/pkg/domain"
	"github.com/monstrong/gracker2/pkg/errs"
	"google.golang.org/grpc/metadata"
)

func checkIDfromMD(MD metadata.MD) (userID uuid.UUID, err error) {
	vals := MD.Get(grpcmetakey.GRPCMetaKeyUserID)

	if len(vals) > 0 {
		if userID, err = uuid.Parse(vals[0]); err != nil {
			return uuid.Nil, &errs.GatewayWrongFieldError{Field: grpcmetakey.GRPCMetaKeyUserID, Reason: "invalid user id format"}
		}
	} else {
		return uuid.Nil, &errs.GatewayWrongFieldError{Field: grpcmetakey.GRPCMetaKeyUserID, Reason: "empty user id"}
	}

	return
}

func checkSenderRolefromMD(MD metadata.MD) (SenderRole domain.Role, err error) {
	vals := MD.Get(grpcmetakey.GRPCMetaKeyRole)

	if len(vals) > 0 {
		roleValue, err := strconv.ParseInt(vals[0], 10, 8)
		if err != nil {
			return domain.RoleUnspecified, &errs.GatewayWrongFieldError{Field: grpcmetakey.GRPCMetaKeyRole, Reason: "invalid role format"}
		}
		SenderRole = domain.Role(roleValue)
	} else {
		return domain.RoleUnspecified, &errs.GatewayWrongFieldError{Field: grpcmetakey.GRPCMetaKeyRole, Reason: "empty role"}
	}

	return
}
