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
	dbStats := db.DB().Stats()
	if !db.RecordNotFound() && result.Error != nil {
		fmt.Printf("dbStats ERROR MaxOpenConnections:%v OpenConnections:%v InUse:%v Idle:%v WaitCount:%v WaitDuration:%v MaxIdleClosed:%v MaxLifetimeClosed:%v \n", dbStats.MaxOpenConnections, dbStats.OpenConnections, dbStats.InUse, dbStats.Idle, dbStats.WaitCount, dbStats.WaitDuration, dbStats.MaxIdleClosed, dbStats.MaxLifetimeClosed)
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("dbStats OK MaxOpenConnections:%v OpenConnections:%v InUse:%v Idle:%v WaitCount:%v WaitDuration:%v MaxIdleClosed:%v MaxLifetimeClosed:%v \n", dbStats.MaxOpenConnections, dbStats.OpenConnections, dbStats.InUse, dbStats.Idle, dbStats.WaitCount, dbStats.WaitDuration, dbStats.MaxIdleClosed, dbStats.MaxLifetimeClosed)

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

	id, err := strconv.Atoi(c.URLParams["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := getUserInfo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, fmt.Sprintf("user not found ID:%v", id), http.StatusNotFound)
		return
	}

	tpl = template.Must(template.ParseFiles("view/user/edit.html"))
	tpl.Execute(w, FormData{*user, ""})
}

// UserUpdate ユーザ更新
func UserUpdate(c web.C, w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(c.URLParams["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := getUserInfo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, fmt.Sprintf("user not found ID:%v", id), http.StatusNotFound)
		return
	}

	user.Name = r.FormValue("Name")
	if err := model.UserValidate(*user); err != nil {
		var Mess string
		errs := valval.Errors(err)
		for _, errInfo := range errs {
			Mess += fmt.Sprint(errInfo.Error)
		}
		tpl = template.Must(template.ParseFiles("view/user/edit.html"))
		tpl.Execute(w, FormData{*user, Mess})
		return
	}

	db.Save(&user)
	http.Redirect(w, r, "/user/index", 301)
}

// getUserInfo ユーザ情報を取得する
func getUserInfo(id int) (*model.User, error) {

	u := model.User{}
	db = db.First(&u, "id = ?", id)

	if db.Error != nil {
		if db.RecordNotFound() {
			return nil, nil
		}
		return nil, db.Error
	}

	return &u, nil
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
		ConnectionMaxLifetimeSeconds int `yaml:"ConnectionMaxLifetimeSeconds"`
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
	db.DB().SetConnMaxLifetime(time.Duration(config.Connection.ConnectionMaxLifetimeSeconds) * time.Second)
}
