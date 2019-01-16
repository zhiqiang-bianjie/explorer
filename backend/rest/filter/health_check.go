package filter

import (
	"encoding/json"
	"github.com/irisnet/explorer/backend/logger"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/types"
	"github.com/irisnet/explorer/backend/utils"
)

type health struct {
	Status  string `json:"status"`
	ErrCode string `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    struct {
		IsMaintained bool `json:"isMaintained"`
	} `json:"data"`
}

type HealthCheckFilter struct {
}

func (HealthCheckFilter) Type() Type {
	return Pre
}

func (HealthCheckFilter) Paths() []string {
	return []string{GlobalFilterPath}
}

func (HealthCheckFilter) Do(request *model.IrisReq, data interface{}) (interface{}, types.BizCode) {
	traceId := logger.Int64("traceId", request.TraceId)
	logger.Info("AclPreFilter", traceId)
	//TODO
	bz, err := utils.Get("http://192.168.150.7:9180/ops_ctl/latest")
	if err != nil {
		logger.Error("filer/health_check request api error", logger.String("err", err.Error()))
	}

	var h health
	if err := json.Unmarshal(bz, &h); err != nil {
		logger.Error("filer/health_check unmarshal error", logger.String("err", err.Error()))
	}

	if h.Status == "success" && !h.Data.IsMaintained {
		return nil, types.CodeSuccess
	}
	return nil, types.CodeSysmaintenance
}
