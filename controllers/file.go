package controllers

import (
	"github.com/astaxie/beego"
	"webfs/models"
	"os"
)

type FileController struct {
	beego.Controller
}

// 获取文件大小的接口
type Size interface {
	Size() int64
}

func (c *FileController) Test() {
	c.TplName = "fileupload.html"
}

// Upload 上传一个文件
// 返回一个json格式的fileMeta
// jquery-fileupload组件一次上传多个的时候，底层也是一个个的上传
func (c *FileController) Upload() {
	f, h, err := c.GetFile("file")
	if err != nil {
		c.Ctx.WriteString("{'error':true,'message':'获取上传列表异常'}")
		return
	}
	defer f.Close()

	//初始化
	fileMeta := models.NewFileMeta()

	fileMeta.Catalog = c.GetString("catalog")
	fileMeta.Name = h.Filename
	//文件大小
	if sizeInterface, ok := f.(Size); ok {
		fileMeta.Size=sizeInterface.Size()
	}
	// XXX 保存文件
	fileMeta.SaveName = "static/upload/" + h.Filename
	c.SaveToFile("file", fileMeta.SaveName)
	//存储文件元信息
	fileMeta.Save()

	//返回结果
	c.Data["json"] = fileMeta
	c.ServeJSON()
}

//Download 下载文件，直接在静态文件目录，直接返回saveName即可下载

//Del 删除一个已上传的
func (c *FileController) Del() {
	var result map[string]interface{}=make(map[string]interface{})
	id,err:=c.GetInt64("id")
	if err!=nil {
		result["error"]="访问参数异常：没有获取到待删除文件的编码"
		c.Data["json"]=result
		c.ServeJSON()
		return
	}
	//load filemeta
	fm:=models.NewFileMeta()
	fm.Id=id
	fm.Load()

	//如果加载后的名称或者保存位置信息为空，即为没有查询到文件元信息
	if fm.Name=="" || fm.SaveName==""{
		result["error"]="没有加载到文件元信息"
		c.Data["json"]=result
		c.ServeJSON()
		return
	}

	//根据SaveName删除文件
	os.Remove(fm.SaveName)
	fm.Del()

	//构造返回结果
	result[fm.Name]=true
	c.Data["json"]=result
	c.ServeJSON()
}

//Query 获取一个已上传文件的信息
func (c *FileController) Query() {
	id,err:=c.GetInt64("id")
	if err!=nil {
		c.Data["json"]=nil
		c.ServeJSON()
		return
	}
	//load filemeta
	fm:=models.NewFileMeta()
	fm.Id=id
	fm.Load()

	if fm.Name=="" || fm.SaveName==""{
		c.Data["json"]=nil
	}else{
		c.Data["json"]=fm
	}

	c.ServeJSON()
}
