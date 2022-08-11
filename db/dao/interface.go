package dao

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
)

// CounterInterface 计数器数据模型接口
type CounterInterface interface {
	GetCounter(id int32) (*model.CounterModel, error)
	UpsertCounter(counter *model.CounterModel) error
	ClearCounter(id int32) error
}

// CounterInterfaceImp 计数器数据模型实现
type CounterInterfaceImp struct{}

// Imp 实现实例
var Imp CounterInterface = &CounterInterfaceImp{}

type CheckInterface interface {
	UpsertCheck(check *model.Check) error
	GetCheck(userId int, checkDate time.Time) (*model.Check, error)
	ExistCheck(userId int, checkDate time.Time) bool
	GetCheckFail(openId string, month time.Time) int64
}

type CheckInterfaceImp struct{}

var CheckImp CheckInterface = &CheckInterfaceImp{}

func (imp *CheckInterfaceImp) GetCheckFail(openId string, month time.Time) int64 {
	var check = new(model.Check)
	count := int64(0)
	cli := db.Get()
	date := month.Format("2006-01")
	cli.Table(check.TableName()).Where("username = ?", openId).Where("DATE_FORMAT(check_date,'%Y-%m')=?", date).Where("status=?", model.CheckStatus_3).Count(&count)
	return count
}
func (imp *CheckInterfaceImp) UpsertCheck(check *model.Check) error {
	cli := db.Get()
	return cli.Table(check.TableName()).Save(check).Error
}
func (imp *CheckInterfaceImp) ExistCheck(userId int, checkDate time.Time) bool {
	var check = new(model.Check)
	count := int64(0)
	cli := db.Get()
	date := checkDate.Format("2006-01-02")
	err := cli.Table(check.TableName()).Where("user_id = ?", userId).Where("check_date = ?", date).Count(&count).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
	}
	return count > 0
}
func (imp *CheckInterfaceImp) GetCheck(userId int, checkDate time.Time) (*model.Check, error) {
	var err error
	var check = new(model.Check)

	cli := db.Get()
	err = cli.Table(check.TableName()).Where("user_id = ?", userId).Where("check_date = ?", checkDate).First(check).Error

	return check, err
}

type CheckUserInterface interface {
	GetCheckUserByOpenId(openId string) (*model.CheckUser, error)
	GetCheckUserById(id int32) (*model.CheckUser, error)
	SetCheckDate(userId int, ContinuousDay int, CheckDate time.Time) error
	UpsertCherUser(user *model.CheckUser) error
}
type CheckUserInterfaceImp struct{}

var CheckUserImp = &CheckUserInterfaceImp{}

func (c *CheckUserInterfaceImp) UpsertCherUser(u *model.CheckUser) error {
	cli := db.Get()
	return cli.Model(u).Save(u).Error
}
func (c *CheckUserInterfaceImp) GetCheckUserByOpenId(openId string) (*model.CheckUser, error) {
	cli := db.Get()
	var checkUser = new(model.CheckUser)
	err := cli.Table(checkUser.TableName()).Where("open_id = ?", openId).First(checkUser).Error
	return checkUser, err
}

func (c *CheckUserInterfaceImp) GetCheckUserById(id int32) (*model.CheckUser, error) {
	cli := db.Get()
	var checkUser = new(model.CheckUser)
	err := cli.Table(checkUser.TableName()).Where("id = ?", id).First(checkUser).Error
	return checkUser, err
}

func (c *CheckUserInterfaceImp) SetCheckDate(userId int, ContinuousDay int, CheckDate time.Time) error {
	cli := db.Get()
	return cli.Exec("Update check_user set continuous_day=? , last_day=? where id = ?", ContinuousDay, CheckDate, userId).Error
}
