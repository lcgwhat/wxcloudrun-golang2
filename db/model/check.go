package model

import "time"

type Check struct {
	Id        int       `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT" json:"id"`
	Username  string    `gorm:"column:username;type:varchar(50);NOT NULL" json:"username"`
	CheckDate time.Time `gorm:"column:check_date;type:datetime;comment:签到日期;NOT NULL" json:"check_date"`
	UserId    int       `gorm:"column:user_id;type:int(11);NOT NULL" json:"user_id"`
	Note      string    `gorm:"column:note;type:varchar(255);comment:留言;NOT NULL" json:"note"`
	Status    int       `gorm:"column:status;type:smallint(6);default:1;comment:状态（1.正常状态2.补签3.失效）;NOT NULL" json:"status"`
}

func (m *Check) TableName() string {
	return "check"
}
