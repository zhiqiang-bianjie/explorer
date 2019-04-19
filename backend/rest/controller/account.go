package controller

import (
	"github.com/gorilla/mux"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/service"
	"github.com/irisnet/explorer/backend/types"
)

// mux.Router registrars
func RegisterAccount(r *mux.Router) error {
	funs := []func(*mux.Router) error{
		registerQueryAccount,
		registerQueryAllAccount,
	}

	for _, fn := range funs {
		if err := fn(r); err != nil {
			return err
		}
	}
	return nil
}

func registerQueryAccount(r *mux.Router) error {

	doApi(r, types.UrlRegisterQueryAccount, "GET", func(request model.IrisReq) interface{} {
		address := Var(request, "address")
		result := service.GetAccountService().Query(address)
		return result
	})

	return nil
}

func registerQueryAllAccount(r *mux.Router) error {
	doApi(r, types.UrlRegisterQueryAllAccount, "GET", func(request model.IrisReq) interface{} {
		page, size := GetPage(request)
		result := service.GetAccountService().QueryAll(page, size)
		return result
	})
	return nil
}
