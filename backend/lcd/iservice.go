package lcd

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/explorer/backend/conf"
	"github.com/irisnet/explorer/backend/logger"
	"github.com/irisnet/explorer/backend/utils"
)

func QuerySvcDef(name, defChainId string) (result SvcDefVo, err error) {
	url := fmt.Sprintf(UrlServiceDefine, conf.Get().Hub.LcdUrl, defChainId, name)
	resBytes, err := utils.Get(url)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(resBytes, &result); err != nil {
		logger.Error("get svcDef error", logger.String("err", err.Error()))
		return result, err
	}
	return result, nil
}

func QuerySvcBindingsList(name, defChainId string) (result []SvcBindVo, err error) {
	url := fmt.Sprintf(UrlServiceBindingList, conf.Get().Hub.LcdUrl, defChainId, name)
	resBytes, err := utils.Get(url)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(resBytes, &result); err != nil {
		logger.Error("get svcDef error", logger.String("err", err.Error()))
		return result, err
	}
	return result, nil
}
