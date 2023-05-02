package loadBalance

import (
	"gorm.io/gorm"
)

type ServerHttp struct {
	gorm.Model
	ServerId uint `json:"server_id"`

	Ip     string `json:"ip"`
	Weight int    `json:"weight"`
}

func (_ *ServerHttp) TableName() string {
	return "server_http"
}

// ============================= service ================================================

func (s *ServerHttp) PageList(db *gorm.DB, search ServerInfo) ([]ServerHttp, error) {

	httpSearch := &ServerHttp{
		ServerId: search.ID,
	}

	var httpResult []ServerHttp

	tx := db.Table(s.TableName()).Where(httpSearch).Find(&httpResult)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound { //数据库查找不到数据也是一种错误(匹配不到相应的数据)
		return nil, tx.Error
	}

	return httpResult, nil
}
