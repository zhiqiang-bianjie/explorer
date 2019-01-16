package filter

import (
	"fmt"
	"github.com/irisnet/explorer/backend/cache"
	"github.com/irisnet/explorer/backend/logger"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/types"
	"strings"
	"time"
)

type RateLimitPreFilter struct{}

func (RateLimitPreFilter) Do(request *model.IrisReq, data interface{}) (interface{}, types.BizCode) {
	traceId := logger.Int64("traceId", request.TraceId)
	logger.Info("RateLimitPreFilter", traceId)
	remoteIp := strings.Split(request.RemoteAddr, ":")[0]
	key := fmt.Sprintf("%s:%s%s", types.RedisKeyRateLimitPrefix, remoteIp, request.RequestURI)
	rate, err := cache.Instance().GetInt(key)
	if err != nil {
		cache.Instance().Set(key, int64(1), 1*time.Second)
		return nil, types.CodeSuccess
	} else {
		cache.Instance().Incr(key)
	}
	//TODO
	if rate > 10 {
		logger.Warn("rateLimit filter starts working", logger.String("IP", remoteIp))
		return nil, types.CodeRateLimit
	}
	return nil, types.CodeSuccess
}

func (RateLimitPreFilter) Paths() string {
	return GlobalFilterPath
}

func (RateLimitPreFilter) Type() Type {
	return Pre
}
