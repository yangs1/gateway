package loadBalance

import (
	"gateway/model"
	"gorm.io/gorm"
)

type ServerInfo struct {
	gorm.Model
	ServerType int    `json:"server_type" ` //接入类型,0是前缀
	RoundType  int    `json:"round_type"`   // 轮询方式
	ServerName string `json:"service_name"`
	ServerDesc string `json:"service_desc"`

	OpenAuth  int    `json:"open_auth"`
	BlackList string `json:"black_list"`
	WhiteList string `json:"white_list"`
}

func (_ *ServerInfo) TableName() string {
	return "server_info"
}

// ============================= service ================================================

func (s *ServerInfo) PageList(db *gorm.DB, page model.PageInput) ([]ServerInfo, int64, error) {
	total := int64(0)
	list := []ServerInfo{} //切片
	offset := (page.PageNum - 1) * page.PageSize

	query := db.Table(s.TableName())

	tx := query.Limit(page.PageSize).Offset(offset).Order("id desc").Find(&list)

	if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound { //数据库查找不到数据也是一种错误(匹配不到相应的数据)
		return nil, 0, tx.Error
	}

	// todo 校验
	query.Limit(page.PageSize).Offset(offset).Count(&total)

	return list, total, nil
}
