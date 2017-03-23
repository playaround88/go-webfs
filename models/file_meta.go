package models

import (
	"fmt"
	"log"
	"strconv"
)

const (
	KEY_PREFIX  = "file_meta:"
	SEQ = "file_meta:seq"
)

type FileMeta struct {
	Id int64	`json:"id""`
	//文件名
	Name string	`json:"name"`
	//文件大小
	Size int64	`json:"size"`
	//文件访问地址
	Url string	`json:"url"`
	//文件删除的URL地址
	DeleteURL string	`json:"deleteUrl"`
	//删除类型，DEL/
	DeleteType string	`json:"deleteType"`
	//服务器保存的文件名
	SaveName string	`json:"saveName"`
	//文件保存的分类，共有情况下的区分
	Catalog string	`json:"catalog"`
	//可能的错误信息
	Error string	`json:"error"`
}

func NewFileMeta() *FileMeta{
	return &FileMeta{
		DeleteType:"DEL",
		DeleteURL:"/del/",
	}
}



func (fm *FileMeta) Load(){
	//判断id值异常
	if fm.Id <= 0 {
		log.Println("加载FileMeta，但是ID值异常id：",fm.Id)
		return
	}
	//从redis查询hash
	ssmcd:=client.HGetAll(KEY_PREFIX+fmt.Sprintf("%d",fm.Id))
	result:=ssmcd.Val()

	// 映射结果
	fm.Name=result["Name"]

	s,_:=strconv.ParseInt(result["Size"],10, 64)
	fm.Size=s

	fm.Url=result["Url"]
	fm.DeleteURL=result["DeleteURL"]
	fm.DeleteType=result["DeleteType"]
	fm.SaveName=result["SaveName"]
	fm.Catalog=result["Catalog"]
}

func (fm *FileMeta) Del(){
	client.Del(KEY_PREFIX+fmt.Sprintf("%d",fm.Id))
}

func (fm *FileMeta) Save(){
	//生成ID值，并回设
	ic:=client.Incr(SEQ)
	fm.Id=ic.Val()
	//XXX 设置Hash字段，通道方式节省交互网络开销
	client.HSet(KEY_PREFIX+fmt.Sprintf("%d",fm.Id),"Name", fm.Name)
	client.HSet(KEY_PREFIX+fmt.Sprintf("%d",fm.Id),"Size", fm.Size)
	client.HSet(KEY_PREFIX+fmt.Sprintf("%d",fm.Id),"Url", fm.Url)
	client.HSet(KEY_PREFIX+fmt.Sprintf("%d",fm.Id),"DeleteURL", fm.DeleteURL)
	client.HSet(KEY_PREFIX+fmt.Sprintf("%d",fm.Id),"DeleteType", fm.DeleteType)
	client.HSet(KEY_PREFIX+fmt.Sprintf("%d",fm.Id),"SaveName", fm.SaveName)
	client.HSet(KEY_PREFIX+fmt.Sprintf("%d",fm.Id),"Catalog", fm.Catalog)
}