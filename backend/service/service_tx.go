package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/irisnet/explorer/backend/model"
	"github.com/irisnet/explorer/backend/orm/document"
	"github.com/irisnet/explorer/backend/types"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type TxService struct {
	BaseService
}

func (service *TxService) GetModule() Module {
	return Tx
}

func (service *TxService) QueryList(query bson.M, page, pageSize int) model.PageVo {
	var data []document.CommonTx
	pageInfo := queryPage(document.CollectionNmCommonTx, &data, query, desc(document.Tx_Field_Time), page, pageSize)
	pageInfo.Data = buildData(data)
	return pageInfo
}

func (service *TxService) QueryLatest(query bson.M, page, pageSize int) model.PageVo {
	var data []document.CommonTx
	pageInfo := queryPage(document.CollectionNmCommonTx, &data, query, desc(document.Tx_Field_Time), page, pageSize)
	return pageInfo
}

func (service *TxService) Query(hash string) interface{} {
	dbm := getDb()
	defer dbm.Session.Close()

	var result document.CommonTx
	query := bson.M{}
	query[document.Tx_Field_Hash] = hash
	err := dbm.C(document.CollectionNmCommonTx).Find(query).Sort(desc(document.Tx_Field_Time)).One(&result)
	if err != nil {
		panic(types.CodeNotFound)
	}

	tx := service.buildTx(result)

	switch tx.(type) {
	case model.GovTx:
		govTx := tx.(model.GovTx)
		return govTx
	case model.StakeTx:
		stakeTx := tx.(model.StakeTx)
		if stakeTx.Type == types.TypeBeginRedelegation {
			var res document.TxMsg
			err := dbm.C(document.CollectionNmTxMsg).Find(bson.M{document.TxMsg_Field_Hash: stakeTx.Hash}).One(&res)
			if err != nil {
				break
			}
			var msg model.MsgBeginRedelegate
			if err = json.Unmarshal([]byte(res.Content), &msg); err == nil {
				stakeTx.From = msg.DelegatorAddr
				stakeTx.To = msg.ValidatorDstAddr
				stakeTx.Source = msg.ValidatorSrcAddr
			}
		}
		return stakeTx
	}
	return tx
}

func (service *TxService) QueryByAcc(address string, page, size int) model.PageVo {
	var data []document.CommonTx
	query := bson.M{}
	query["$or"] = []bson.M{{document.Tx_Field_From: address}, {document.Tx_Field_To: address}}
	var typeArr []string
	typeArr = append(typeArr, types.TypeTransfer)
	typeArr = append(typeArr, types.DeclarationList...)
	typeArr = append(typeArr, types.StakeList...)
	typeArr = append(typeArr, types.GovernanceList...)
	query[document.Tx_Field_Type] = bson.M{
		"$in": typeArr,
	}
	return queryPage(document.CollectionNmCommonTx, &data, query, desc(document.Tx_Field_Time), page, size)
}

func (service *TxService) CountByType(query bson.M) model.TxStatisticsVo {
	var typeArr []string
	typeArr = append(typeArr, types.TypeTransfer)
	typeArr = append(typeArr, types.DeclarationList...)
	typeArr = append(typeArr, types.StakeList...)
	typeArr = append(typeArr, types.GovernanceList...)
	query[document.Tx_Field_Type] = bson.M{
		"$in": typeArr,
	}

	var counter []struct {
		Type  string `bson:"_id,omitempty"`
		Count int
	}

	c := getDb().C(document.CollectionNmCommonTx)
	defer c.Database.Session.Close()

	pipe := c.Pipe(
		[]bson.M{
			{"$match": query},
			{"$group": bson.M{
				"_id":   "$type",
				"count": bson.M{"$sum": 1},
			}},
		},
	)

	pipe.All(&counter)

	var result model.TxStatisticsVo
	for _, cnt := range counter {
		switch types.Convert(cnt.Type) {
		case types.Trans:
			result.TransCnt = result.TransCnt + cnt.Count
		case types.Declaration:
			result.DeclarationCnt = result.DeclarationCnt + cnt.Count
		case types.Stake:
			result.StakeCnt = result.StakeCnt + cnt.Count
		case types.Gov:
			result.GovCnt = result.GovCnt + cnt.Count
		}
	}
	return result
}

func (service *TxService) CountByDay() []model.TxDayVo {
	c := getDb().C(document.CollectionNmCommonTx)
	defer c.Database.Session.Close()

	now := time.Now()
	d, _ := time.ParseDuration("-336h") //14 days ago
	start := now.Add(d)

	fromDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, start.Location())

	log.Println(fmt.Sprintf("from:%s,to:%s", fromDate.String(), endDate.String()))

	pipe := c.Pipe(
		[]bson.M{
			{"$match": bson.M{
				"time": bson.M{
					"$gte": fromDate,
					"$lt":  endDate,
				},
			}},
			{"$project": bson.M{
				"day": bson.M{"$substr": []interface{}{"$time", 0, 10}},
			}},
			{"$group": bson.M{
				"_id":   "$day",
				"count": bson.M{"$sum": 1},
			}},
			{"$sort": bson.M{
				"_id": 1,
			}},
		},
	)
	var txDays []model.TxDayVo
	var result []model.TxDayVo
	pipe.All(&txDays)
	var i time.Duration
	var j int
	day := start
	for day.Unix() < endDate.Unix() {
		key := day.Format("2006-01-02")
		if len(txDays) > j && txDays[j].Time == key {
			result = append(result, txDays[j])
			j++
		} else {
			var txDay model.TxDayVo
			txDay.Time = key
			txDay.Count = 0
			result = append(result, txDay)
		}
		i++
		day = start.Add(i * 24 * time.Hour)
	}
	return result
}

func buildData(txs []document.CommonTx) []interface{} {
	var txList []interface{}

	if len(txs) == 0 {
		return txList
	}
	for _, tx := range txs {
		txResp := txService.buildTx(tx)
		txList = append(txList, txResp)
	}

	return txList
}

func (service *TxService) buildTx(tx document.CommonTx) interface{} {
	db := getDb()
	defer db.Session.Close()

	switch types.Convert(tx.Type) {
	case types.Trans:
		return model.TransTx{
			BaseTx: buildBaseTx(tx),
			From:   tx.From,
			To:     tx.To,
			Amount: tx.Amount,
		}
	case types.Declaration:
		dtx := model.DeclarationTx{
			BaseTx:   buildBaseTx(tx),
			SelfBond: tx.Amount,
			Owner:    tx.From,
			Pubkey:   tx.StakeCreateValidator.PubKey,
		}
		var blackList = service.QueryBlackList(db)
		if tx.Type == types.TypeCreateValidator {
			var moniker = tx.StakeCreateValidator.Description.Moniker
			var identity = tx.StakeCreateValidator.Description.Identity
			var website = tx.StakeCreateValidator.Description.Website
			var details = tx.StakeCreateValidator.Description.Details
			if desc, ok := blackList[tx.To]; ok {
				moniker = desc.Moniker
				identity = desc.Identity
				website = desc.Website
				details = desc.Details
			}
			dtx.Moniker = moniker
			dtx.Details = details
			dtx.Website = website
			dtx.Identity = identity
		} else if tx.Type == types.TypeEditValidator {
			var moniker = tx.StakeEditValidator.Description.Moniker
			var identity = tx.StakeEditValidator.Description.Identity
			var website = tx.StakeEditValidator.Description.Website
			var details = tx.StakeEditValidator.Description.Details
			if desc, ok := blackList[tx.From]; ok {
				moniker = desc.Moniker
				identity = desc.Identity
				website = desc.Website
				details = desc.Details
			}
			dtx.Moniker = moniker
			dtx.Details = details
			dtx.Website = website
			dtx.Identity = identity
		} else if tx.Type == types.TypeUnjail {
			candidateDb := db.C(document.CollectionNmStakeRoleCandidate)
			var can document.Candidate
			candidateDb.Find(bson.M{document.Candidate_Field_Address: dtx.Owner}).One(&can)
			var moniker = can.Description.Moniker
			var identity = can.Description.Identity
			var website = can.Description.Website
			var details = can.Description.Details
			if desc, ok := blackList[tx.From]; ok {
				moniker = desc.Moniker
				identity = desc.Identity
				website = desc.Website
				details = desc.Details
			}
			dtx.Moniker = moniker
			dtx.Details = identity
			dtx.Website = website
			dtx.Identity = details
		}
		return dtx
	case types.Stake:
		return model.StakeTx{
			TransTx: model.TransTx{
				BaseTx: buildBaseTx(tx),
				From:   tx.From,
				To:     tx.To,
				Amount: tx.Amount,
			},
		}
	case types.Gov:
		govTx := model.GovTx{
			BaseTx:     buildBaseTx(tx),
			Amount:     tx.Amount,
			From:       tx.From,
			ProposalId: tx.ProposalId,
		}

		txMsgDb := db.C(document.CollectionNmTxMsg)
		var res document.TxMsg
		err := txMsgDb.Find(bson.M{document.TxMsg_Field_Hash: govTx.Hash}).One(&res)
		if err != nil {
			return govTx
		}

		if govTx.Type == types.TypeSubmitProposal {
			var msg model.MsgSubmitProposal
			json.Unmarshal([]byte(res.Content), &msg)
			govTx.Title = msg.Title
			govTx.Description = msg.Description
			govTx.ProposalType = msg.ProposalType
		} else if govTx.Type == types.TypeDeposit {
			var msg model.MsgDeposit
			json.Unmarshal([]byte(res.Content), &msg)
			govTx.ProposalId = msg.ProposalID
			govTx.Amount = msg.Amount
		} else if govTx.Type == types.TypeVote {
			var msg model.MsgVote
			json.Unmarshal([]byte(res.Content), &msg)
			govTx.ProposalId = msg.ProposalID
			govTx.Option = msg.Option
		}

		return govTx
	case types.Service:
		baseTx := buildBaseTx(tx)
		switch tx.Type {
		case types.TypeServiceDefine:
			var svcDef model.SvcDef
			if err := model.UnMarshalJSON([]byte(tx.MsgContent), &svcDef); err == nil {
				return model.ServiceTx{
					BaseTx: baseTx,
					Msg:    svcDef,
				}
			}
		case types.TypeServiceBind:
			var svcBind model.SvcBind
			if err := model.UnMarshalJSON([]byte(tx.MsgContent), &svcBind); err == nil {
				return model.ServiceTx{
					BaseTx: baseTx,
					Msg:    svcBind,
				}
			}
		case types.TypeServiceCall:
			var svcReq model.SvcRequest
			if err := model.UnMarshalJSON([]byte(tx.MsgContent), &svcReq); err == nil {
				decodeBytes, err := base64.StdEncoding.DecodeString(svcReq.Input)
				if err == nil {
					svcReq.Input = string(decodeBytes)
				}
				return model.ServiceTx{
					BaseTx: baseTx,
					Msg:    svcReq,
				}
			}
		case types.TypeServiceRespond:
			var svcResp model.SvcResponse
			if err := model.UnMarshalJSON([]byte(tx.MsgContent), &svcResp); err == nil {
				decodeBytes, err := base64.StdEncoding.DecodeString(svcResp.Output)
				if err == nil {
					svcResp.Output = string(decodeBytes)
				}
				return model.ServiceTx{
					BaseTx: baseTx,
					Msg:    svcResp,
				}
			}
		}
		return baseTx

	}
	return tx
}

func buildBaseTx(tx document.CommonTx) model.BaseTx {
	return model.BaseTx{
		Hash:        tx.TxHash,
		BlockHeight: tx.Height,
		Type:        tx.Type,
		Fee:         tx.ActualFee,
		Status:      tx.Status,
		GasLimit:    tx.Fee.Gas,
		GasUsed:     tx.GasUsed,
		GasPrice:    tx.GasPrice,
		Memo:        tx.Memo,
		Timestamp:   tx.Time,
		From:        tx.From,
	}
}
