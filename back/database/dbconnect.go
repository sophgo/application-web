package database

import (
	"application-web/global"
	"application-web/logger"
	"application-web/pkg/model"
	"application-web/pkg/repo"
	"database/sql"
	"fmt"
	"regexp"

	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // init only
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
)

const sqlite3Go = "sqlite3_with_go_func"
const admin = "admin"
const recordNotFound = "record not found"

func init() {
	sql.Register(sqlite3Go, &sqlite3.SQLiteDriver{
		ConnectHook: func(c *sqlite3.SQLiteConn) error {
			return c.RegisterFunc("regexp", regexp.MatchString, true)
		},
	})
}

var DB *gorm.DB

type DBUtil struct {
	db *gorm.DB
}

func GetDBUtil(db *gorm.DB) *DBUtil {
	return &DBUtil{db: db}
}

func InitDB() {
	var err error

	sqlDb, err := sql.Open(sqlite3Go, global.System.Db.Path)
	if err != nil {
		panic(err)
	}
	DB, err = gorm.Open("sqlite3", sqlDb)
	if err != nil {
		logger.Error("open %s failed, error: %s", global.System.Db.Path, err)
		panic(err)
	} else {
		DB.DB().SetMaxOpenConns(1)
	}
	global.DB = DB

	DB.SingularTable(true)

	_ = GetDBUtil(DB).CreateTableIfNotExist(&model.User{}, "user_id", "id")
	_ = GetDBUtil(DB).CreateTableIfNotExist(&model.Record{}, "record_id", "id")
	_ = GetDBUtil(DB).CreateTableIfNotExist(&model.AlgoTaskSql{}, "task_id", "id")
	_ = GetDBUtil(DB).CreateTableIfNotExist(&model.FaceAlgoTaskSql{}, "facetask_id", "id")
	_ = GetDBUtil(DB).CreateTableIfNotExist(&model.FaceFeature{}, "facefeature_id", "id")
	_ = GetDBUtil(DB).CreateTableIfNotExist(&model.FaceRecord{}, "facerecord_id", "id")
	_ = GetDBUtil(DB).CreateTableIfNotExist(&model.SearchTaskSql{}, "searchtask_id", "id")
	_ = GetDBUtil(DB).CreateTableIfNotExist(&model.SearchResultRecord{}, "searchresult_id", "id")
	_ = GetDBUtil(DB).CreateTableIfNotExist(&model.PersonLibrary{}, "personlibrary_id", "id")

	_, err = repo.QueryUserWithName(admin)
	if err != nil && strings.EqualFold(recordNotFound, err.Error()) {
		repo.SaveUser(&model.User{
			Model: gorm.Model{
				ID: 1,
			},
			UserID:     "admin",
			Status:     "",
			UserName:   admin,
			Password:   global.System.ServerInfo.Password,
			Token:      "",
			Address:    "",
			Role:       "",
			LoginTime:  time.Time{},
			LockedTime: time.Time{},
			ExpireTime: time.Time{},
			Label:      "",
		})
	}

	c := cron.New(cron.WithSeconds())
	_, err = c.AddFunc("0 0 0 * * ?", func() {
		date := time.Now().Add(-time.Hour * 24 * time.Duration(global.System.Db.SaveDays))
		logger.Info("清理数据：%s %v", global.System.Db.SaveDays, date)
	})
	if err != nil {
		fmt.Println("cron init err:", err)
	}

	c.Start()
}

// 创建表, 支持索引
func (d *DBUtil) CreateTableIfNotExist(schema interface{}, indexName string, columns ...string) error {
	var db = d.db
	if !db.HasTable(schema) {
		if err := db.Debug().CreateTable(schema).Error; err != nil {
			return err
		}
		// 添加唯一索引
		if indexName != "" || len(columns) != 0 {
			db.Model(schema).AddIndex("idx_"+indexName, columns...)
		}
	} else {
		db.AutoMigrate(schema)
	}
	return nil
}
