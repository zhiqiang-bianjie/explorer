package filter

import (
	"github.com/irisnet/explorer/backend/logger"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/types"
	"strings"
)

//ACL TODO
const AclStr = "127.0.0.1"

type AclFilter struct {
}

func (AclFilter) GetPath() string {
	return GlobalFilterPath
}

func (AclFilter) Do(request *model.IrisReq, data interface{}) (bool, interface{}, types.BizCode) {
	if "*" == AclStr {
		return true, nil, types.CodeSuccess
	}

	aclList := strings.Split(AclStr, ",")
	if len(aclList) == 0 {
		return true, nil, types.CodeSuccess
	}

	remoteIp := strings.Split(request.RemoteAddr, ":")[0]
	for _, acc := range aclList {
		if acc == remoteIp {
			return true, nil, types.CodeSuccess
		}
	}

	logger.Warn("access unauthorized", logger.String("IP", remoteIp))
	return false, nil, types.CodeUnauthorized
}
