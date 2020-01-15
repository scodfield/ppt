package controllers
import (
	_"os"
	_"io"
	"io/ioutil"
	_"fmt"
	"github.com/astaxie/beego"
	_"github.com/astaxie/beego/validation"
)

type ShowController struct {
	beego.Controller
}


func (this *ShowController) Get() {
	fileAddr, err := ioutil.ReadDir("./photos")
	if err != nil {
		this.Data["Result"] = false
		this.Data["Welcome"] = "sorry, wrong url"
	} else {
		this.Data["Result"] = true
		this.Data["Welcome"] = "Welcome to see"
	}
	this.Data["Photos"] = fileAddr
	// var listHtml string 
	// for _, fileInfo := range fileAddr {
	// 	imgID := fileInfo.Name 
	// 	listHtml += "<li><a href=\"/view?id="+imgID+"\">imgID</a></li>"
	// }
	// listStr = "<ol>"+listHtml+"</ol>"
	this.TplName = "showPhoto.html"
}

