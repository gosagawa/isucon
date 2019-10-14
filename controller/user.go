package controller

import (
	"io/ioutil"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/wcl48/valval"
	"github.com/zenazn/goji/web"
	"gopkg.in/yaml.v2"

	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gosagawa/isucon/model"
)

var tpl *template.Template
var db *gorm.DB

// Password パスワード
const Password = "user:user"

// FormData  フォームデータ
type FormData struct {
	User model.User
	Mess string
}

func init() {
	connect()
}

// UserIndex ユーザ情報
func UserIndex(c web.C, w http.ResponseWriter, r *http.Request) {
	Users := []model.User{}
	result := db.Find(&Users)
	if !db.RecordNotFound() && result.Error != nil {
		fmt.Println(db.DB().Stats())
		panic(result.Error)
	}

	tpl = template.Must(template.ParseFiles("view/user/index.html"))
	tpl.Execute(w, Users)
}

// UserNew ユーザ新規作成
func UserNew(c web.C, w http.ResponseWriter, r *http.Request) {
	tpl = template.Must(template.ParseFiles("view/user/new.html"))
	tpl.Execute(w, FormData{model.User{}, ""})
}

// UserCreate ユーザ作成
func UserCreate(c web.C, w http.ResponseWriter, r *http.Request) {
	User := model.User{Name: r.FormValue("Name")}
	if err := model.UserValidate(User); err != nil {
		var Mess string
		errs := valval.Errors(err)
		for _, errInfo := range errs {
			Mess += fmt.Sprint(errInfo.Error)
		}
		tpl = template.Must(template.ParseFiles("view/user/new.html"))
		tpl.Execute(w, FormData{User, Mess})
	} else {
		db.Create(&User)
		http.Redirect(w, r, "/user/index", 301)
	}
}

// UserEdit ユーザ編集
func UserEdit(c web.C, w http.ResponseWriter, r *http.Request) {
	User := model.User{}
	User.ID, _ = strconv.Atoi(c.URLParams["id"])
	db.Find(&User)
	tpl = template.Must(template.ParseFiles("view/user/edit.html"))
	tpl.Execute(w, FormData{User, ""})
}

// UserUpdate ユーザ更新
func UserUpdate(c web.C, w http.ResponseWriter, r *http.Request) {
	User := model.User{}
	User.ID, _ = strconv.Atoi(c.URLParams["id"])
	db.Find(&User)
	User.Name = r.FormValue("Name")
	if err := model.UserValidate(User); err != nil {
		var Mess string
		errs := valval.Errors(err)
		for _, errInfo := range errs {
			Mess += fmt.Sprint(errInfo.Error)
		}
		tpl = template.Must(template.ParseFiles("view/user/edit.html"))
		tpl.Execute(w, FormData{User, Mess})
	} else {
		db.Save(&User)
		http.Redirect(w, r, "/user/index", 301)
	}
}

// UserDelete ユーザ削除
func UserDelete(c web.C, w http.ResponseWriter, r *http.Request) {
	User := model.User{}
	User.ID, _ = strconv.Atoi(c.URLParams["id"])
	db.Delete(&User)
	http.Redirect(w, r, "/user/index", 301)
}

func connect() {
	yml, err := ioutil.ReadFile("conf/db.yml")
	if err != nil {
		panic(err)
	}

	type configConnection struct {
		MaxOpenConnections           int `yaml:"MaxOpenConnections"`
		MaxIdleConnections           int `yaml:"MaxIdleConnections"`
		ConnectionMaxLifetimeMinutes int `yaml:"ConnectionMaxLifetimeMinutes"`
	}
	type configYml struct {
		Host       string
		Port       string
		User       string
		Password   string
		Db         string
		LogMode    bool
		Connection configConnection `yaml:"Connection"`
	}
	t := make(map[string]configYml)

	err = yaml.Unmarshal([]byte(yml), &t)
	if err != nil {
		panic(err)
	}

	config := t["development"]

	db, err = gorm.Open("mysql", config.User+":"+config.Password+"@tcp("+config.Host+":"+config.Port+")/"+config.Db+"?loc=Asia%2FTokyo&parseTime=true&charset=utf8mb4")
	if err != nil {
		panic(err)
	}
	db.LogMode(config.LogMode)
	db.DB().SetMaxOpenConns(config.Connection.MaxOpenConnections)
	db.DB().SetMaxIdleConns(config.Connection.MaxIdleConnections)
	db.DB().SetConnMaxLifetime(time.Duration(config.Connection.ConnectionMaxLifetimeMinutes) * time.Minute)
	if err != nil {
		panic(err)
	}
}
