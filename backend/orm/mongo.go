package orm

import (
	"github.com/irisnet/explorer/backend/conf"
	"github.com/irisnet/explorer/backend/logger"
	"github.com/irisnet/explorer/backend/model"
	"time"

	"gopkg.in/mgo.v2"
)

func init() {

	dialInfo := &mgo.DialInfo{
		Addrs:     conf.Get().Db.Addrs,
		Database:  conf.Get().Db.Database,
		Username:  conf.Get().Db.UserName,
		Password:  conf.Get().Db.Password,
		Direct:    false,
		Timeout:   time.Second * 10,
		PoolLimit: conf.Get().Db.PoolLimit,
	}

	var err error
	session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		logger.Error("start mongo client failed", logger.String("err", err.Error()))
	}
	session.SetMode(mgo.Monotonic, true)
}

var session *mgo.Session

func GetDatabase() *mgo.Database {
	return session.Clone().DB(conf.Get().Db.Database)
}

func QueryRows(collation string, data interface{}, m map[string]interface{}, sort string, page, size int) model.PageVo {
	c := GetDatabase().C(collation)
	defer c.Database.Session.Close()
	count, err := c.Find(m).Count()
	if err != nil {
		logger.Error("QueryRows Count failed", logger.String("err", err.Error()))
		return model.PageVo{Count: 0, Data: nil}
	}
	err = c.Find(m).Skip((page - 1) * size).Limit(size).Sort(sort).All(data)
	if err != nil {
		logger.Error("QueryRows Find failed", logger.String("err", err.Error()))
		return model.PageVo{Count: count, Data: nil}
	} else {
		return model.PageVo{Count: count, Data: data}
	}
}

func QueryRowsField(query MQuery) (int, []map[string]interface{}, error) {
	var result []map[string]interface{}
	c := GetDatabase().C(query.C)
	defer c.Database.Session.Close()
	count, err := c.Find(query.Q).Count()
	if err != nil {
		return count, result, err
	}
	err = c.Find(query.Q).Select(query.Selector).Skip((query.Page - 1) * query.Size).Limit(query.Size).Sort(query.Sort).All(&result)
	return count, result, err
}

func LimitQuery(query MQuery) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	c := GetDatabase().C(query.C)
	defer c.Database.Session.Close()
	err := c.Find(query.Q).Select(query.Selector).Sort(query.Sort).Limit(query.Size).All(&result)
	if err != nil {
		logger.Error("Limit error", logger.String("err", err.Error()))
	}
	return result, err
}

func QueryRow(collation string, data interface{}, m map[string]interface{}) error {
	c := GetDatabase().C(collation)
	defer c.Database.Session.Close()
	return c.Find(m).One(data)
}

func QueryRowField(query MQuery) map[string]interface{} {
	var result map[string]interface{}
	c := GetDatabase().C(query.C)
	defer c.Database.Session.Close()
	err := c.Find(query.Q).One(&result)
	if err != nil {
		logger.Error("data not found", logger.String("err", err.Error()))
	}
	return result
}

type MQuery struct {
	C        string
	Result   interface{}
	Q        map[string]interface{}
	Sort     string
	Page     int
	Size     int
	Selector interface{}
}
