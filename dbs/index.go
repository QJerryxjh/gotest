package dbs

import (
	"errors"
	"os"

	"github.com/Unknwon/goconfig"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

// 获取当前执行栈的绝对路径
func GetCallerPath() (path string) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return
}

func InitDB(env string) (err error) {
	path := GetCallerPath()
	config, err := goconfig.LoadConfigFile(path + "/conf.ini")
	if err != nil {
		return
	}
	username, err := config.GetValue(env, "username")
	if err != nil {
		return
	}
	database, err := config.GetValue(env, "database")
	if err != nil {
		return
	}
	password, err := config.GetValue(env, "password")
	if err != nil {
		return
	}
	host, err := config.GetValue(env, "host")
	if err != nil {
		return
	}
	port, err := config.GetValue(env, "port")
	if err != nil {
		return
	}
	charset, err := config.GetValue(env, "charset")
	if err != nil {
		return
	}
	parseTime, err := config.GetValue(env, "parseTime")
	if err != nil {
		return
	}
	loc, err := config.GetValue(env, "loc")
	if err != nil {
		return
	}
	if username == "" || password == "" || host == "" || charset == "" || parseTime == "" || loc == "" || port == "" || database == "" {
		err = errors.New("缺少参数")
		return
	}

	db, err := gorm.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/"+database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		err = errors.New(err.Error())
		return
	}
	db.DB().SetMaxIdleConns(1024)
	db.DB().SetMaxOpenConns(256)
	DB = db
	return
}

func InitEnvironment(env string) error {
	path := GetCallerPath()
	flag := false
	if _, err := os.Stat(path + "/conf.ini"); err != nil {
		file, err := os.Create("conf.ini")
		if err != nil {
			return err
		}
		defer file.Close()
		flag = true
	}

	if flag {
		con, err := goconfig.LoadConfigFile(path + "/conf.ini")
		if err != nil {
			return err
		}
		con.SetValue("root", "database", "mysql")
		con.SetValue("root", "username", "root")
		con.SetValue("root", "password", "12345678")
		con.SetValue("root", "host", "127.0.0.1")
		con.SetValue("root", "port", "3306")
		con.SetValue("root", "charset", "utf8")
		con.SetValue("root", "parseTime", "True")
		con.SetValue("root", "loc", "Local")
		con.SetValue(env, "database", "blog")
		con.SetValue(env, "username", "admin")
		con.SetValue(env, "password", "12345678")
		con.SetValue(env, "host", "127.0.0.1")
		con.SetValue(env, "port", "3306")
		con.SetValue(env, "charset", "utf8")
		con.SetValue(env, "parseTime", "True")
		con.SetValue(env, "loc", "Local")
		err = goconfig.SaveConfigFile(con, "conf.ini")
	}
	err := InitDB(env)
	if err != nil {
		return err
	}
	return nil
}

func Close() (err error) {
	err = DB.Close()
	if err != nil {
		return
	}
	return
}
