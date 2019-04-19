package service

import (
	"fmt"
	"github.com/irisnet/explorer/backend/logger"
	"github.com/irisnet/explorer/backend/orm"
	"github.com/irisnet/explorer/backend/orm/document"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	accountService = &AccountService{
		BaseService{
			collection: document.CollectionNmAccount,
		},
	}

	blockService = &BlockService{
		BaseService{
			collection: document.CollectionNmBlock,
		},
	}

	commonService = &CommonService{}

	proposalService = &ProposalService{
		BaseService{
			collection: document.CollectionNmProposal,
		},
	}

	stakeService = &CandidateService{
		BaseService{
			collection: document.CollectionNmStakeRoleCandidate,
		},
	}

	txService = &TxService{
		BaseService{
			collection: document.CollectionNmCommonTx,
		},
	}
	delegatorService = &DelegatorService{
		BaseService{
			collection: document.CollectionNmStakeRoleDelegator,
		},
	}
)

func GetAccountService() *AccountService {
	return accountService
}

func GetBlockService() *BlockService {
	return blockService
}

func GetCommonService() *CommonService {
	return commonService
}

func GetProposalService() *ProposalService {
	return proposalService
}

func GetCandidateService() *CandidateService {
	return stakeService
}

func GetDelegatorService() *DelegatorService {
	return delegatorService
}

func GetTxService() *TxService {
	return txService
}

type BaseService struct {
	tid        string
	collection string
}

func (base *BaseService) Collection() string {
	return base.collection
}

func (base *BaseService) SetTid(traceId string) {
	base.tid = traceId
}

func (base *BaseService) GetTid() string {
	return base.tid
}

func (base *BaseService) GetTraceLog() zap.Field {
	return logger.String("traceId", base.GetTid())
}

func (base *BaseService) QueryBlackList(database *mgo.Database) map[string]document.BlackList {
	var blackListStore = database.C(document.CollectionNmBlackList)
	var blackList []document.BlackList
	var blackListMap = make(map[string]document.BlackList)
	if err := blackListStore.Find(nil).All(&blackList); err == nil {
		for _, v := range blackList {
			blackListMap[v.OperatorAddr] = v
		}
	}
	return blackListMap
}

func (base *BaseService) QueryAll(selector, condition bson.M, sort string, size int, result interface{}) error {
	var query = orm.NewQuery()
	defer query.Release()
	query.SetCollection(base.Collection()).
		SetCondition(condition).
		SetSelector(selector).
		SetSort(sort).
		SetSize(size).
		SetResult(result)

	err := query.Exec()
	if err != nil {
		logger.Error("queryAll error", logger.Any("query", condition), logger.String("err", err.Error()))
	}
	return err
}

func (base *BaseService) QueryOne(selector, condition bson.M, result interface{}) error {
	var query = orm.NewQuery()
	defer query.Release()
	query.SetCollection(base.Collection()).
		SetCondition(condition).
		SetSelector(selector).
		SetResult(result)

	err := query.Exec()
	if err != nil {
		logger.Error("queryOne", logger.Any("query", condition), logger.String("err", err.Error()))
	}
	return err
}

func (base *BaseService) PageQuery(selector, condition bson.M, sort string, page, size int, result interface{}) (int, error) {
	var query = orm.NewQuery()
	defer query.Release()
	query.SetCollection(base.Collection()).
		SetCondition(condition).
		SetSelector(selector).
		SetSort(sort).
		SetPage(page).
		SetSize(size).
		SetResult(result)

	cnt, err := query.ExecPage()
	if err != nil {
		logger.Error("pageQuery", logger.Any("query", condition), logger.String("err", err.Error()))
	}

	return cnt, err
}

func getDb() *mgo.Database {
	return orm.GetDatabase()
}

func desc(field string) string {
	return fmt.Sprintf("-%s", field)
}

func asc(field string) string {
	return fmt.Sprintf("%s", field)
}
