package conf

import (
	"database/sql"
	firebase "firebase.google.com/go/v4"
	"git.in.codoon.com/Overseas/runbox/first-test/cloudsql"
)

var Firebaer_app *firebase.App
var Fb_mysql *sql.DB

func InitBase() {

	//初始化数据库
	db := cloudsql.GetDB()
	Fb_mysql = db
	//初始化firebaseAPP
	/*	app := firebase_conf.InitializeAppWithRefreshToken()

		Firebaer_app = app*/

}
