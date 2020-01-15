package controllers
import (
	_"fmt"
	"login/db"
	"github.com/astaxie/beego"
	_"github.com/astaxie/beego/validation"
)

type RegisterController struct {
	beego.Controller
}

// type User struct {
// 	AccId  int
// 	Name string 
// 	DevId string 
// 	SdkType string 
// }

func (this *RegisterController) Get() {
	this.TplName = "register.html"
}

func (this *RegisterController) Post() {
	uname := this.GetString("username")
	if db.WhetherRegistered(uname) {
		this.Data["IsAlready"] = true
		this.TplName = "regResult.html"
	} else {
		user := db.User{
			Name : uname,
			Password : this.GetString("password"),
		}
		setErr := db.SetAccountInfo(uname,user)
		if setErr != nil {
			this.Data["IsFail"] = true
			this.TplName = "regResult.html"
		} else {
			this.Data["IsSucc"] = true
			this.TplName = "login.html"
		}
	}
	this.Data["Name"] = uname
}