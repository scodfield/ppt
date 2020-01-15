package controllers
import (
	_"fmt"
	"time"
	_"encoding/json"
	_"unsafe"
	"login/db"
	"github.com/astaxie/beego"
	_"github.com/astaxie/beego/validation"
	_"github.com/mitchellh/mapstructure"
)

type LoginController struct {
	beego.Controller
}


func (this *LoginController) Get() {
	this.TplName = "login.html"
}

func (this *LoginController) Post() {
	var info db.User
	uname := this.GetString("username")
	upass := this.GetString("password")
	if db.WhetherRegistered(uname) {
		infoIn := db.GetAccountInfo(uname)
		info, _ = infoIn.(db.User)
		// check the password
		if upass == info.Password {
			var loginLog db.LoginLog
			loginLog = db.LoginLog{AccId:info.AccId,Name:info.Name,LoginTime:time.Now()}
			db.SetLoginCache(info)
			db.SetLoginLog(loginLog) 
			this.Data["Name"] = uname
			this.TplName = "myInfo.html"
		} else {
			this.TplName = "login.html"
		}
	} else {
		this.TplName = "register.html"
	}
}