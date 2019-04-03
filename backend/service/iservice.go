package service

import (
	"github.com/irisnet/explorer/backend/lcd"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/orm/document"
	"github.com/irisnet/explorer/backend/types"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type IService struct {
	BaseService
}

func (service *IService) GetModule() Module {
	return Iservice
}

func (service *IService) QueryList(page, size int) interface{} {
	var result []document.SvcDef
	sort := desc(document.Tx_Field_Time)
	pageInfo := queryPage(document.CollectionNmSvcDef, &result, nil, sort, page, size)
	var resVo []model.SvcDef
	if pageInfo.Count > 0 {
		for _, svcDef := range result {
			resVo = append(resVo, buildSvcDef(svcDef))
		}
		pageInfo.Data = resVo
	}

	return pageInfo
}

func (service *IService) Query(name, defChainId string) interface{} {
	var svcInfo IServiceInfo
	if svcDef, err := lcd.QuerySvcDef(name, defChainId); err == nil {
		svcInfo.SvcDef = model.SvcDef{
			Name:              svcDef.Definition.Name,
			ChainId:           svcDef.Definition.ChainID,
			Author:            svcDef.Definition.Author,
			AuthorDescription: svcDef.Definition.AuthorDescription,
			Description:       svcDef.Definition.Description,
			IDLContent:        svcDef.Definition.IdlContent,
		}
	}
	dbm := getDb()
	defer dbm.Session.Close()
	svcBindDoc := dbm.C(document.CollectionNmSvcBind)
	svcInvocationDoc := dbm.C(document.CollectionNmSvcInvocation)

	svcInfo.SvcBindList = service.QuerySvcBinding(name, defChainId, 1, 10, svcBindDoc)

	svcInfo.SvcTxList = service.QuerySvcInvocation(name, defChainId, 1, 10, svcInvocationDoc)
	return svcInfo
}

func (service *IService) QuerySvcBinding(name, defChainId string, page, size int, collection *mgo.Collection) interface{} {
	var svcBindList []model.SvcBind
	var bindingList []document.SvcBind
	query := bson.M{}
	query["def_name"] = name
	query["def_chain_id"] = defChainId

	if collection == nil {
		dbm := getDb()
		defer dbm.Session.Close()
		collection = dbm.C(document.CollectionNmSvcBind)
	}

	sort := desc(document.Tx_Field_Time)
	result := findByPage(collection, &bindingList, query, sort, page, size)
	if result.Count > 0 {
		for _, binding := range bindingList {
			svcBindList = append(svcBindList, model.SvcBind{
				Hash:        binding.Hash,
				DefName:     binding.DefName,
				DefChainID:  binding.DefChainID,
				BindChainID: binding.BindChainID,
				Provider:    binding.Provider,
				BindingType: binding.BindingType,
				Prices:      binding.Prices,
				Available:   binding.Available,
				Level: model.Level{
					AvgRspTime: binding.Level.AvgRspTime,
					UsableTime: binding.Level.UsableTime,
				},
			})
		}

	}
	result.Data = svcBindList
	return result
}

func (service *IService) QuerySvcInvocation(name, defChainId string, page, size int, collection *mgo.Collection) interface{} {
	var svcTxList []model.SvcTx
	var invocationList []document.SvcInvocation
	query := bson.M{}
	query["def_name"] = name
	query["def_chain_id"] = defChainId

	if collection == nil {
		dbm := getDb()
		defer dbm.Session.Close()
		collection = dbm.C(document.CollectionNmSvcInvocation)
	}

	sort := desc(document.Tx_Field_Time)
	result := findByPage(collection, &invocationList, query, sort, page, size)

	if result.Count > 0 {
		for _, invocation := range invocationList {
			sendAddr := invocation.Consumer
			receiveAddr := invocation.Provider
			if invocation.TxType == types.TypeServiceRespond {
				sendAddr = invocation.Provider
				receiveAddr = invocation.Consumer
			}
			svcTxList = append(svcTxList, model.SvcTx{
				Hash:        invocation.Hash,
				ReqId:       invocation.ReqId,
				TxType:      invocation.TxType,
				SendAddr:    sendAddr,
				ReceiveAddr: receiveAddr,
				Height:      invocation.Height,
				Data:        invocation.Data,
				Time:        invocation.Time,
			})
		}
	}
	result.Data = svcTxList
	return result
}

type IServiceInfo struct {
	model.SvcDef
	SvcBindList interface{} `json:"svc_bind_list"`
	SvcTxList   interface{} `json:"svc_tx_list"`
}

func buildSvcDef(svcDef document.SvcDef) model.SvcDef {
	return model.SvcDef{
		Hash:              svcDef.Hash,
		Name:              svcDef.Code,
		ChainId:           svcDef.ChainId,
		Author:            svcDef.Author,
		AuthorDescription: svcDef.AuthorDescription,
		Description:       svcDef.Description,
		Status:            svcDef.Status,
	}
}
