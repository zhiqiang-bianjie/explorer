package filter

import (
	"github.com/irisnet/explorer/backend/logger"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/types"
	"strings"
)

//ACL TODO
const AclList = "127.0.0.1"

type AclPreFilter struct {
}

func (AclPreFilter) Paths() []string {
	return []string{GlobalFilterPath}
}

func (AclPreFilter) Type() Type {
	return Pre
}

func (AclPreFilter) Do(request *model.IrisReq, data interface{}) (interface{}, types.BizCode) {
	traceId := logger.Int64("traceId", request.TraceId)
	logger.Info("AclPreFilter", traceId)
	if "*" == AclList {
		return nil, types.CodeSuccess
	}

	aclList := strings.Split(AclList, ",")
	if len(aclList) == 0 {
		return nil, types.CodeSuccess
	}

	remoteIp := strings.Split(request.RemoteAddr, ":")[0]
	for _, acc := range aclList {
		if acc == remoteIp {
			return nil, types.CodeSuccess
		}
	}

	logger.Warn("access unauthorized", logger.String("IP", remoteIp))
	return nil, types.CodeUnauthorized
}
