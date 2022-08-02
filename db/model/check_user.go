package model

import (
	"time"
)

//CheckUser 用户
type CheckUser struct {
	Id            int       `gorm:"column:id;type:int(11);primary_key" json:"id"`
	OpenId        string    `gorm:"column:open_id;type:varchar(255);comment:微信openid" json:"open_id"`
	Username      string    `gorm:"column:username;type:varchar(50)" json:"username"`
	ContinuousDay int       `gorm:"column:continuous_day;type:int(11);default:0;comment:连续签到" json:"continuous_day"`
	LastDay       time.Time `gorm:"column:last_day;type:date;comment:上次签到日期" json:"last_day"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime;default:CURRENT_TIMESTAMP" json:"create_time"`
}

func (m *CheckUser) TableName() string {
	return "check_user"
}
