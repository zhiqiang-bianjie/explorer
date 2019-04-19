package controller

import (
	"github.com/gorilla/mux"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/service"
	"github.com/irisnet/explorer/backend/types"
	"github.com/irisnet/explorer/backend/utils"
)

// mux.Router registrars
func RegisterBlock(r *mux.Router) error {
	funs := []func(*mux.Router) error{
		registerQueryBlock,
		registerQueryBlocks,
		registerQueryRecentBlocks,
		registerQueryBlocksPrecommits,
	}

	for _, fn := range funs {
		if err := fn(r); err != nil {
			return err
		}
	}
	return nil
}

func registerQueryBlock(r *mux.Router) error {
	doApi(r, types.UrlRegisterQueryBlock, "GET", func(request model.IrisReq) interface{} {
		h := Var(request, "height")
		height, ok := utils.ParseInt(h)
		if !ok {
			panic(types.CodeInValidParam)
			return nil
		}
		result := service.GetBlockService().Query(height)
		return result
	})
	return nil
}

func registerQueryBlocks(r *mux.Router) error {
	doApi(r, types.UrlRegisterQueryBlocks, "GET", func(request model.IrisReq) interface{} {
		page := int(utils.ParseIntWithDefault(QueryParam(request, "page"), 1))
		size := int(utils.ParseIntWithDefault(QueryParam(request, "size"), 100))
		result := service.GetBlockService().QueryList(page, size)
		return result
	})

	return nil
}

func registerQueryRecentBlocks(r *mux.Router) error {
	doApi(r, types.UrlRegisterQueryRecentBlocks, "GET", func(request model.IrisReq) interface{} {
		return service.GetBlockService().QueryRecent()
	})

	return nil
}

func registerQueryBlocksPrecommits(r *mux.Router) error {
	doApi(r, types.UrlRegisterQueryBlocksPrecommits, "GET", func(request model.IrisReq) interface{} {

		address := Var(request, "address")
		page, size := GetPage(request)

		result := service.GetBlockService().QueryPrecommits(address, page, size)
		return result
	})

	return nil
}
