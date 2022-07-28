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
}

type CheckInterfaceImp struct{}

var CheckImp CheckInterface = &CheckInterfaceImp{}

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
