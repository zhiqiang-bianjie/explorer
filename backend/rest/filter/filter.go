package filter

import (
	"github.com/irisnet/explorer/backend/logger"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/types"
)

const (
	Pre  Type = 0
	Post Type = 1

	GlobalFilterPath = "*"
)

type Filter interface {
	Do(request *model.IrisReq, data interface{}) (bool, interface{}, types.BizCode)
	GetPath() string
}
type FChain []Filter

func NewFChain() FChain {
	return make(FChain, 0)
}
func (chain FChain) Append(f Filter) FChain {
	c := append(chain, f)
	return c
}

type Router map[string]FChain
type Type int

var preRouter = make(Router, 0)
var postRouter = make(Router, 0)

func Register(path string, typ Type, fs FChain) {
	var router = getRouter(typ)
	if _, ok := router[path]; ok {
		logger.Panic("duplicate registration filter", logger.String("path", path))
	}
	router[path] = fs
}

func DoFilters(req *model.IrisReq, data interface{}, typ Type) (bool, interface{}, types.BizCode) {
	var router = getRouter(typ)
	//check global filters
	globalFilters, ok := router[GlobalFilterPath]
	if ok {
		for _, f := range globalFilters {
			ok, data, err := f.Do(req, data)
			if !ok {
				return false, data, err
			}
		}
	}

	//check custom filters
	customFilters, ok := router[req.RequestURI]
	if ok {
		for _, f := range customFilters {
			ok, data, err := f.Do(req, data)
			if !ok {
				return false, data, err
			}
		}
	}
	return true, nil, types.CodeSuccess
}

func getRouter(typ Type) Router {
	var router Router
	switch typ {
	case Pre:
		router = preRouter
		break
	case Post:
		router = postRouter
		break
	default:
		logger.Panic("not existed filter type", logger.Any("type", typ))
	}
	return router
}
