package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/types"
	"github.com/irisnet/explorer/backend/utils"
	"io/ioutil"
	"net/http"
)

func RegisterNodes(r *mux.Router) error {
	apis := []func(*mux.Router) error{
		RegisterQueryNodes,
		RegisterQueryNodeLocation,
		RegisterQueryFaucet,
		RegisterApply,
	}

	for _, fn := range apis {
		if err := fn(r); err != nil {
			return err
		}
	}
	return nil
}

func RegisterQueryNodes(r *mux.Router) error {
	doApi(r, types.UrlRegisterQueryNodes, "GET", func(request model.IrisReq) interface{} {
		bz := utils.GetNodes()
		return bz
	})

	return nil
}

func RegisterQueryNodeLocation(r *mux.Router) error {
	doApi(r, types.UrlRegisterQueryIp, "POST", func(request model.IrisReq) interface{} {
		body, _ := ioutil.ReadAll(request.Body)
		var params map[string][]string
		json.Unmarshal(body, &params)
		ips := params["ipdata"]
		var ipMap = make([]string, len(ips))
		for i, ip := range ips {
			url := fmt.Sprintf(types.UrlNodeLocation, ip)
			resp, err := http.Get(url)
			if err != nil {
				panic(types.CodeExtSysFailed)
				return nil
			}
			body, _ := ioutil.ReadAll(resp.Body)
			ipMap[i] = string(body)
		}
		return ipMap
	})

	return nil
}

func RegisterQueryFaucet(r *mux.Router) error {
	doApi(r, types.UrlRegisterQueryFaucet, "GET", func(request model.IrisReq) interface{} {
		return utils.GetFaucetAccount(request.Request)
	})
	return nil
}

func RegisterApply(r *mux.Router) error {
	doApi(r, types.UrlRegisterApply, "POST", func(request model.IrisReq) interface{} {
		return utils.Apply(request.Request)
	})
	return nil
}