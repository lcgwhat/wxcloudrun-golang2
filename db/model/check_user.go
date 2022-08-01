package model

import (
	"database/sql"
	"time"
)

type CheckUser struct {
	Id       int            `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Name     string         `gorm:"column:name;type:varchar(50);NOT NULL" json:"name"`
	CreateAt time.Time      `gorm:"column:create_at;type:datetime;NOT NULL" json:"create_at"`
	OpenId   sql.NullString `gorm:"column:open_id;type:varchar(255);comment:微信openid" json:"open_id"`
}

func (m *CheckUser) TableName() string {
	return "check_user"
}
