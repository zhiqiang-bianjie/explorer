package controller

import (
	"github.com/gorilla/mux"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/service"
	"github.com/irisnet/explorer/backend/types"
	"github.com/irisnet/explorer/backend/utils"
)

func RegisterIService(r *mux.Router) error {
	funs := []func(*mux.Router) error{
		registerQueryServiceList,
		registerQueryService,
		registerQuerySvcBinding,
		registerQuerySvcInvocation,
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
		page := int(utils.ParseIntWithDefault(QueryParam(request, "page"), 1))
		size := int(utils.ParseIntWithDefault(QueryParam(request, "size"), 10))
		return iService.QueryList(page, size)
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

func registerQuerySvcBinding(r *mux.Router) error {
	doApi(r, types.UrlRegisterQuerySvcBinding, "GET", func(request model.IrisReq) interface{} {
		svcName := Var(request, "svcName")
		defChainId := Var(request, "defChainId")
		page := int(utils.ParseIntWithDefault(QueryParam(request, "page"), 1))
		size := int(utils.ParseIntWithDefault(QueryParam(request, "size"), 10))
		return iService.QuerySvcBinding(svcName, defChainId, page, size, nil)
	})
	return nil
}

func registerQuerySvcInvocation(r *mux.Router) error {
	doApi(r, types.UrlRegisterQuerySvcInvocation, "GET", func(request model.IrisReq) interface{} {
		svcName := Var(request, "svcName")
		defChainId := Var(request, "defChainId")
		page := int(utils.ParseIntWithDefault(QueryParam(request, "page"), 1))
		size := int(utils.ParseIntWithDefault(QueryParam(request, "size"), 10))
		return iService.QuerySvcInvocation(svcName, defChainId, page, size, nil)
	})
	return nil
}
