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

type RateLimitFilter struct{}

func (RateLimitFilter) Do(request *model.IrisReq, data interface{}) (bool, interface{}, types.BizCode) {

	remoteIp := strings.Split(request.RemoteAddr, ":")[0]
	key := fmt.Sprintf("%s:%s%s", types.RedisKeyRateLimitPrefix, remoteIp, request.RequestURI)
	rate, err := cache.Instance().GetInt(key)
	if err != nil {
		cache.Instance().Set(key, int64(1), 1*time.Second)
		return true, nil, types.CodeSuccess
	} else {
		cache.Instance().Incr(key)
	}
	//TODO
	if rate > 10 {
		logger.Warn("rateLimit filter starts working", logger.String("IP", remoteIp))
		return false, nil, types.CodeRateLimit
	}
	return true, nil, types.CodeSuccess
}

func (RateLimitFilter) GetPath() string {
	return GlobalFilterPath
}
