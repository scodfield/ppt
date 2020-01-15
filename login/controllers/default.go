package controllers
import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Data["Website"] = "beego.who"
	this.Data["Email"] = "who@gmail.com"
	this.TplName = "index.html"
}
