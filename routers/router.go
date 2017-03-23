package routers

import (
	"webfs/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.AutoRouter(&controllers.FileController{})
}
