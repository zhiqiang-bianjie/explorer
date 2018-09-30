package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/irisnet/explorer/backend/types"
	"github.com/irisnet/explorer/backend/utils"
	"github.com/irisnet/irishub-sync/store/document"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func registerQueryGovTx(r *mux.Router) error {
	r.HandleFunc("/api/tx/gov/{page}/{size}", queryGovTx).Methods("GET")
	return nil
}

func queryGovTx(w http.ResponseWriter, r *http.Request) {
	var data []document.CommonTx
	query := bson.M{}
	query["type"] = bson.M{
		"$in": types.GovernanceList,
	}
	var result types.Page
	count := utils.QueryByPage("tx_common", &data, query, "-time", r)
	if count > 0 {
		result = types.Page{
			Count: count,
			Data:  buildGovResp(data),
		}
	}
	resp, err := json.Marshal(result)
	if err == nil {
		w.Write(resp)
	}
}

func buildGovResp(txs []document.CommonTx) []types.GovTx {
	var txList []types.GovTx
	if len(txs) == 0 {
		return txList
	}
	for _, tx := range txs {
		txResp := buildTx(tx)
		txList = append(txList, txResp.(types.GovTx))
	}

	return txList
}
