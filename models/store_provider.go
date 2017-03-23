package models

import "os"

type StoreProvider interface {
	//保存文件
	Store(file os.File) (*FileMeta,error)
	//获取文件
	GetFile(meta *FileMeta) (os.File,error)
	//删除文件
	DelFile(meta *FileMeta) error
}
