package routers
import (
	"login/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/register", &controllers.RegisterController{})
	beego.Router("/upload", &controllers.UploadController{})
	beego.Router("/show", &controllers.ShowController{})
}