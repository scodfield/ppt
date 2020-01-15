package controllers
import (
	"os"
	"io"
	"fmt"
	"github.com/astaxie/beego"
	_"github.com/astaxie/beego/validation"
)

type UploadController struct {
	beego.Controller
}


func (this *UploadController) Get() {
	this.TplName = "loadPhoto.html"
}

func (this *UploadController) Post() {
	result := "succes to upload!!"
	photoFile,photoHeader, err := this.GetFile("Image")
	if err != nil {
		fmt.Println("err: ", err)
	} else {
		defer photoFile.Close()
		// create destination file making sure the path is writeable.
		dst, err := os.Create("photos/" + photoHeader.Filename)
		if err != nil {
			result = "faile to create file"
		}
		defer dst.Close()
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, photoFile); err != nil {
			result = "faile to write file"
		}
	}
	this.Ctx.WriteString(result)
}
