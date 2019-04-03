package model

import (
	"github.com/irisnet/explorer/backend/orm/document"
	"time"
)

type SvcDef struct {
	Hash              string   `json:"hash"`
	Name              string   `json:"name"`
	ChainId           string   `json:"chain_id"`
	Description       string   `json:"description"`
	Tags              []string `json:"tags"`
	Author            string   `json:"author"`
	AuthorDescription string   `json:"author_description"`
	IDLContent        string   `json:"idl_content"`
	Status            string   `json:"status"`
}

type SvcBind struct {
	Hash        string         `json:"hash"`
	DefName     string         `json:"def_name"`
	DefChainID  string         `json:"def_chain_id"`
	BindChainID string         `json:"bind_chain_id"`
	Provider    string         `json:"provider"`
	BindingType string         `json:"binding_type"`
	Deposit     document.Coins `json:"deposit"`
	Prices      document.Coins `json:"price"`
	Level       Level          `json:"level"`
	Available   bool           `json:"available"`
}

type Level struct {
	AvgRspTime int64 `json:"avg_rsp_time"`
	UsableTime int64 `json:"usable_time"`
}

type SvcRequest struct {
	DefChainID  string `json:"def_chain_id"`
	DefName     string `json:"def_name"`
	BindChainID string `json:"bind_chain_id"`
	ReqChainID  string `json:"req_chain_id"`
	MethodID    int16  `json:"method_id"`
	Provider    string `json:"provider"`
	Consumer    string `json:"consumer"`
	Input       string `json:"input"`
	ServiceFee  Coins  `json:"service_fee"`
	Profiling   bool   `json:"profiling"`
}

type SvcResponse struct {
	ReqChainID string `json:"req_chain_id"`
	RequestID  string `json:"request_id"`
	Provider   string `json:"provider"`
	Output     string `json:"output"`
	ErrorMsg   string `json:"error_msg"`
}

type SvcTx struct {
	Hash        string    `json:"hash"`
	ReqId       string    `json:"req_id"`
	TxType      string    `json:"tx_type"`
	SendAddr    string    `json:"send_addr"`
	ReceiveAddr string    `json:"receive_addr"`
	Height      int64     `json:"height"`
	Data        string    `json:"data"`
	Time        time.Time `json:"time"`
}
