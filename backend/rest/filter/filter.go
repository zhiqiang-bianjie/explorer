package filter

import (
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/types"
	"strings"
)

const (
	Pre  Type = 0
	Post Type = 1

	GlobalFilterPath = "*"
)

var filterChain = make(FChain, 0)

type Filter interface {
	Do(request *model.IrisReq, data interface{}) (interface{}, types.BizCode)
	Paths() string
	Type() Type
}

type FChain []Filter
type Type int

func Register(filter Filter) {
	filterChain = append(filterChain, filter)
}

func DoFilters(req *model.IrisReq, data interface{}, typ Type) (interface{}, types.BizCode) {
	var reqUrl = req.RequestURI
	for _, f := range filterChain {
		if typ != f.Type() {
			continue
		}
		paths := f.Paths()
		if paths == GlobalFilterPath {
			data, err := f.Do(req, data)
			if !err.Success() {
				return data, err
			}
		} else {
			urls := strings.Split(paths, ",")
			for _, url := range urls {
				if reqUrl == url {
					data, err := f.Do(req, data)
					if !err.Success() {
						return data, err
					}
				}
			}
		}

	}
	return nil, types.CodeSuccess
}
