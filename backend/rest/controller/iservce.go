package controller

import (
	"github.com/gorilla/mux"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/service"
	"github.com/irisnet/explorer/backend/types"
)

func RegisterIService(r *mux.Router) error {
	funs := []func(*mux.Router) error{
		registerQueryServiceList,
		registerQueryService,
	}

	for _, fn := range funs {
		if err := fn(r); err != nil {
			return err
		}
	}
	return nil
}

type IService struct {
	*service.IService
}

var iService = IService{
	service.Get(service.Iservice).(*service.IService),
}

func registerQueryServiceList(r *mux.Router) error {
	doApi(r, types.UrlRegisterQueryServiceList, "GET", func(request model.IrisReq) interface{} {
		return iService.QueryList(1, 10)
	})
	return nil
}

func registerQueryService(r *mux.Router) error {
	doApi(r, types.UrlRegisterQueryService, "GET", func(request model.IrisReq) interface{} {
		svcName := Var(request, "svcName")
		defChainId := Var(request, "defChainId")
		return iService.Query(svcName, defChainId)
	})
	return nil
}