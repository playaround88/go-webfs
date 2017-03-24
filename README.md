# go-webfs
基于go beego的文件上传服务器

## 依赖
| 封装jquery-fileupload插件  
| 基于go语言的web框架beego开发  
| 持久层选用了redis

## 特性
1. 对于多个关联多个文件的情况，简化文件字段，只需要一个字段存储。
2. 文件上传后，显示已上传文件列表。
3. 可以根据存储的文件ids（comma seperate），渲染更新页面和显示页面。
4. redis存储，高性能，可扩展。


## 跨域说明
跨域借助html5的一个特性，即在response header中设置access-control-allow-origin参数。  
目前，仅仅是直接设置了*（允许所有跨域访问），可以根据具体的使用修改。

## 快速上手
基本的使用示例请参考/file/test  
注意，跨域情况下，应该把所有的js引用都带上webfs的访问地址

