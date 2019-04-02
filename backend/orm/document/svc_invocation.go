package document

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const CollectionNmSvcInvocation = "nj_scf_svc_invocation"

type SvcInvocation struct {
	Hash        string    `bson:"hash"`
	ReqId       string    `bson:"req_id"`
	TxType      string    `bson:"tx_type"`
	DefChainID  string    `bson:"def_chain_id"`
	DefName     string    `bson:"def_name"`
	BindChainID string    `bson:"bind_chain_id"`
	ReqChainID  string    `bson:"req_chain_id"`
	MethodID    int16     `bson:"method_id"`
	Consumer    string    `bson:"consumer"`
	Provider    string    `bson:"provider"`
	Height      int64     `bson:"height"`
	Data        string    `bson:"data"`
	Time        time.Time `bson:"time"`
}

func (m SvcInvocation) Name() string {
	return CollectionNmSvcInvocation
}

func (m SvcInvocation) PkKvPair() map[string]interface{} {
	return bson.M{"hash": m.Hash}
}
