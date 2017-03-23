package main

import (
	_ "webfs/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	//"strings"
	"log"
)

func main() {
	//添加过滤器
	beego.InsertFilter("/*", beego.BeforeRouter, AccessAllowFilter)

	beego.Run()
}

var AccessAllowFilter = func (ctx *context.Context){
	log.Println("enter AccessAllowFilter ...")
	////获取客户端IP地址
	//origin:=ctx.Request.RemoteAddr
	//
	////获取配置文件中配置的地址信息
	//allowOrigin:=beego.AppConfig.String("allowOrigin")
	//
	////验证客户端IP是否是合法访问地址
	//if strings.Contains(allowOrigin,origin){
		//如果合法则写响应头
		ctx.Output.Header("Access-Control-Allow-Origin","*")
	//}

}