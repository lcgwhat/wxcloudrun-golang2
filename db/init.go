package db

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"sync"
	"time"
	"wxcloudrun-golang/db/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type dbInstance2 struct {
	db *gorm.DB
	sync.Mutex
}

var dbInstance *gorm.DB

type dbConf struct {
	user     string
	pwd      string
	addr     string
	dataBase string
}

var conf dbConf = dbConf{
	user:     "root",
	pwd:      "",
	addr:     "127.0.0.1:3306",
	dataBase: "test",
}

// Init 初始化数据库
func Init() error {
	port := os.Getenv("MY_PORT")

	source := "%s:%s@tcp(%s)/%s?readTimeout=1500ms&writeTimeout=1500ms&charset=utf8&loc=Local&&parseTime=true"
	user := os.Getenv("MYSQL_USERNAME")
	pwd := os.Getenv("MYSQL_PASSWORD")
	addr := os.Getenv("MYSQL_ADDRESS")
	dataBase := os.Getenv("MYSQL_DATABASE")
	if dataBase == "" {
		dataBase = "golang_demo"
	}
	if port == "" {
		user = conf.user
		pwd = conf.pwd
		addr = conf.addr
		dataBase = conf.dataBase
	}
	source = fmt.Sprintf(source, user, pwd, addr, dataBase)
	fmt.Println("start init mysql with ", source)

	db, err := gorm.Open(mysql.Open(source), &gorm.Config{
		Logger: getGormLogger(), // 使用自定义 Logger
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		}})
	if err != nil {
		fmt.Println("DB Open error,err=", err.Error())
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("DB Init error,err=", err.Error())
		return err
	}

	// 用于设置连接池中空闲连接的最大数量
	sqlDB.SetMaxIdleConns(100)
	// 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(200)
	// 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	dbInstance = db
	initMySqlTables(dbInstance)
	fmt.Println("finish init mysql with ", source)
	return nil
}

// Get ...
func Get() *gorm.DB {
	var try = 3
	if dbInstance == nil {
		err := Init()
		if err != nil {
			try--
			time.Sleep(time.Millisecond * 200)
			if try == 0 {
				panic(err)
			}
		} else {
			return dbInstance
		}

	}
	db, err := dbInstance.DB()
	if err != nil {
		panic(err)
	}
	err = db.Ping()

	if err != nil {
		for try > 0 {
			err = Init()
			if err != nil {
				try--
				time.Sleep(time.Millisecond * 200)
				if try == 0 {
					panic(err)
				}
			} else {
				return dbInstance
			}
		}
	}
	return dbInstance
}

// 数据库表初始化
func initMySqlTables(db *gorm.DB) {
	err := db.AutoMigrate(
		model.Check{},
	)
	if err != nil {
		fmt.Println("migrate table failed", err.Error())
	}
}

// 自定义 Logger（使用文件记录日志）
func getGormLogger() logger.Interface {
	var logMode logger.LogLevel

	switch "info" {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	return logger.New(getGormLogWriter(), logger.Config{
		SlowThreshold:             200 * time.Millisecond, // 慢 SQL 阈值
		LogLevel:                  logMode,                // 日志级别
		IgnoreRecordNotFoundError: false,                  // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  false,                  // 禁用彩色打印
	})
}
func getGormLogWriter() logger.Writer {
	var writer io.Writer
	port := os.Getenv("MY_PORT")

	// 是否启动日志文件
	if port != "" {
		// 自定义writer
		writer = &lumberjack.Logger{
			Filename:   "./sql.log",
			MaxSize:    500,
			MaxBackups: 3,
			MaxAge:     20,
			Compress:   true,
		}
	} else {
		// 默认 Writer
		writer = os.Stdout
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}
