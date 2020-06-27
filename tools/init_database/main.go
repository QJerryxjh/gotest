package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Unknwon/goconfig"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DbManger struct {
	ID        uint      `gorm:"primary_key"`
	Username  string    `gorm:"type:varchar(16);not null;unique_index"`
	Password  string    `gorm:"type:varchar(32);not null"`
	Name      string    `gorm:"type:varchar(64);not null;"`
	Email     string    `gorm:"type:varchar(32)"`
	Gender    string    `gorm:"type:varchar(16)"`
	Birthday  time.Time `gorm:"type:datetime;default:null"`
	CreatedAt time.Time `gorm:"type:datetime"`
}

type InitDataType struct {
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

// 获取当前执行栈的绝对路径
func GetCallerPath() (path string) {
	_, file, _, ok := runtime.Caller(0)
	res := filepath.Dir(file)
	if !ok {
		panic(errors.New("Can not get current file info"))
	}
	return res
}

func cryptoMd5(str string) string {
	md5Init := md5.New()
	md5Init.Write([]byte(str))
	return hex.EncodeToString(md5Init.Sum(nil))
}

// 数据库配置
func Config() (err error) {
	filePath := GetCallerPath()
	iniFilePath := filePath + "/conf.ini"
	flag := false
	// 查看当前兄弟文件有没有conf.ini
	if _, err := os.Stat(iniFilePath); err != nil {
		// 没有就创建一个
		file, err1 := os.Create("conf.ini")
		if err1 != nil {
			return err1
		}
		err2 := file.Close()
		if err2 != nil {
			return err2
		}
		flag = true
	}

	conf, err := goconfig.LoadConfigFile(iniFilePath)
	if err != nil {
		return
	}

	if flag {
		// 没有ini文件，是自己创建的新文件，为文件添加值
		conf.SetValue("root", "database", "mysql")
		conf.SetValue("root", "username", "root")
		conf.SetValue("root", "password", "12345678")
		conf.SetValue("root", "host", "127.0.0.1")
		conf.SetValue("root", "port", "3306")
		conf.SetValue("root", "charset", "utf8")
		conf.SetValue("root", "parseTime", "True")
		conf.SetValue("root", "loc", "Local")

		conf.SetValue("pro", "database", "blog")
		conf.SetValue("pro", "username", "admin")
		conf.SetValue("pro", "password", "12345678")
		conf.SetValue("pro", "host", "127.0.0.1")
		conf.SetValue("pro", "port", "3306")
		conf.SetValue("pro", "charset", "utf8")
		conf.SetValue("pro", "parseTime", "True")
		conf.SetValue("pro", "loc", "Local")
		err = goconfig.SaveConfigFile(conf, "conf.ini")
		if err != nil {
			return
		}
	}

	username, err := conf.GetValue("pro", "username")
	if err != nil {
		return
	}
	database, err := conf.GetValue("pro", "database")
	if err != nil {
		return
	}
	password, err := conf.GetValue("pro", "password")
	if err != nil {
		return
	}
	host, err := conf.GetValue("pro", "host")
	if err != nil {
		return
	}
	port, err := conf.GetValue("pro", "port")
	if err != nil {
		return
	}
	rootPassword, err := conf.GetValue("root", "password")
	if err != nil {
		return
	}

	var initDatabase InitDataType
	initDatabase.Username = username
	initDatabase.Password = password
	initDatabase.Port = port
	initDatabase.Host = host
	initDatabase.Database = database
	err = initDatabase.InitDatabase(rootPassword)
	if err != nil {
		return
	}

	return
}

func (self *InitDataType) InitDatabase(rootPassword string) (err error) {
	db, err := gorm.Open("mysql", "root:"+rootPassword+"@(127.0.0.1)/mysql?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return
	}
	err = db.Exec("CREATE USER if not exists '" + self.Username + "'@'%' IDENTIFIED BY'" + self.Password + "';").Error

	defer func() {
		if err != nil {
			err = db.Exec("drop user '" + self.Username + "'@'%';").Error
			err = db.Exec("drop user '" + self.Username + "'@'localhost';").Error
		}
	}()
	if err != nil {
		return
	}

	err = db.Exec("alter user " + self.Username + "@'%' identified with mysql_native_password by '" + self.Password + "';").Error

	if err != nil {
		return
	}

	err = db.Exec("create database if not exists " + self.Database + " charset utf8;").Error
	defer func() {
		if err != nil {
			db.Exec("drop database " + self.Database + ";")
		}
	}()
	if err != nil {
		return
	}

	err = db.Exec("grant all privileges on " + self.Database + ".* to '" + self.Username + "'@'%';").Error
	if err != nil {
		return
	}
	err = db.Exec("flush privileges;").Error
	if err != nil {
		return
	}
	err = db.Exec("use " + self.Database + ";").Error
	if err != nil {
		return
	}

	if !db.HasTable(&DbManger{}) {
		err = db.CreateTable(&DbManger{}).Error
		if err != nil {
			return
		}
	}

	var user DbManger
	var count int
	err = db.Where("username = ?", "admin").First(&user).Count(&count).Error
	if count > 0 {
		// 存在管理员了
		return
	}

	password := "123123"
	password = cryptoMd5(password)
	user1_birthday, err := time.Parse("01/02/2006", "06/27/2020")
	if err != nil {
		return
	}
	created_at, err := time.Parse("01/02/2006", "06/27/2020")
	if err != nil {
		return
	}
	user1 := DbManger{Username: "admin", Name: "Manager1", Email: "qjerryxjh@163.com", Gender: "男", Birthday: user1_birthday, CreatedAt: created_at}
	user1.Password = password

	err = db.Create(&user1).Error
	if err != nil {
		return
	}

	return
}

func main() {
	err := Config()

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	if err != nil {
		panic(err)
	} else {
		log.Println("初始化管理员成功")
	}
}
