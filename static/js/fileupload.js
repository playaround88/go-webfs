/**
 * 强化jquery-file-upload插件功能。
 * 
 * @Author wutb
 * @date 2016-04-12
 * 
 * @dependence artTemplate(模板引擎)
 * 
 * 1.上传成功后表格显示。
 * 2.更新页面渲染。 
 * 3.查看页面渲染。
 */

/*
 * 1.上传成功后表格显示。 
 * 使用方法 
 * 初始化文件上传组件时，把fileupload_options，放到第一个调用参数。
 * $('#fileupload').fileupload(fileupload.options,{...自定义配置项});
 */
//上传文件的统一路径
if(!webfs){
	alert("文件上传服务器url配置异常，文件上传功能不可用！");
}
if(!webfs_catalog){
	var webfs_catalog='default';
	alert("没有设置文件存放服务器的分类目录，使用default");
}
var fileupload={};
fileupload['uploadPath'] = webfs+'/upload';
fileupload['delPath'] = webfs+'/del'
fileupload['queryPath'] = webfs+'/query';

var fileTabTmp='<table class="uploadFiles">'
			+'<thead><tr>'
				+'<th>ID</th>'
				+'<th>文件名称</th>'
				+'<th>文件大小(KB)</th>'
				+'<th>上传状态</th>'
				+'<th>操作</th>'
			+'</tr></thead>'
			+'<tbody></tbody>'
			+'</table>';

var fileTrTmp='<tr>'
				+'<td>{{id}}</td>'
				+'<td><a target="_blank" href="/{{saveName}}">{{name}}</a></td>'
				+'<td>{{size | sizeTransfer:"K"}}K</td>'
				+'<td>{{error?error:"上传成功"}}</td>'
				+'<td><a class="delbtn" delurl="'+fileupload.delPath+'?id={{id}}">删除</a></td>'
			+'</tr>';

//模板引擎的配置
/*
 * 文件大小的转换，只支持K和M两种转换
 */
template.helper('sizeTransfer',function(size, unit){
	unit=unit.toLowerCase();
	var result=0;
	switch (unit) {
	case 'k':
		result=size/1024;
		break;
	case 'm':
		result=size/(1024*1024);
		break;
	default:
		throw 'not support unit';
		break;
	}
	return result;
});

fileupload.options = {
	dataType : 'json',
	url : fileupload.uploadPath+'?catalog='+webfs_catalog, 
	valueId:'', //上传组件所属与模型对应字段
	tabRenderTo:'', //上传表格，要渲染的目的地
	done : function(e, options) {
		//请求异常
		if(options.textStatus != 'success'){
			alert('服务器异常，请求失败！' + options.textStatus);
			return false;
		}
		//返回正常，但json中包含异常
		if(options.result.error){
			alert(options.result.error);
			return false;
		}
		//正常情况的处理
		var fileTab=$(options.tabRenderTo).children('table');
		if(fileTab.length==0){
			$(options.tabRenderTo).html(fileTabTmp);
			//重新绑定参数
			fileTab=$(options.tabRenderTo).children('table');
		}
		var file_tbody=fileTab.children('tbody');
		var tr_render=template.compile(fileTrTmp);

		var tr=$(tr_render(options.result));
		tr.children('td:last').children('a').click({"options":options}, fileupload_del);
		file_tbody.append(tr);
		//回调函数
		options.afterDone(options.result);
	},
	fail : function(e, options){
		alert('上传异常，状态码：'+options.textStatus+',error:'+options.errorThrown);
	},
	afterDone:function(file){//默认将file的ID值，附加到valueId组件的value中
		var v=$(this.valueId);
		if(v.val()){
			v.val(v.val()+','+file.id);
		}else{
			v.val(file.id);
		}
	},
	afterDel:function(){
		var val=$(this.valueId).val();
		
		val=val.replace(","+this.options.id,"");
		val=val.replace(this.options.id,"");
		
		$(this.valueId).val(val)
	}
};

function fileupload_del(e){
	var obj=this;
	var url=$(obj).attr('delurl');
    function delTr(){
    	$(obj).parent().parent().remove();
    }
	$.get(url,{},function(data,status){
		if(status=='success'){
			var fileName='';
			for(var n in data){
				fileName=n;
				break;//每次只能删除一个文件
			}
			if(!fileName || ! data[fileName] ){
				alert("文件删除失败");
				return ;
			}
			delTr();
			alert('文件删除成功！');
			//回调接口函数
			e.data.options.afterDel();
		}else{
			alert('服务器异常!');
		}
	},'json');
}

/*
 * 2.更新页面渲染。 
 * jsp页面渲染好三个标签
 * <input id="fileupload" type="file" name="files[]" multiple />
 * <input id="record" type="hidden" value="" />
 * <div id="files" class="files"> </div>
 * js在dom加载完毕后，调用下面的方法，完成页面的渲染
 * conf={
 * 		'valueId':'#importFile',//初始化的字段值存放组件ID
 * 		'fileId':'#fileupload',//文件组件的ID
 * 		'tabRenderTo':'#files' //表格渲染到的容器ID
 * }
 */
fileupload.initUpdate=function(conf){
	if(!conf["valueId"] || !conf["fileId"] || !conf["tabRenderTo"]){
		alert('配置错误，请查证！');
		return ;
	}
	//初始化上传组件
	var file=$(conf.fieldId).fileupload(fileupload.options,{
		valueId:conf.valueId,
		tabRenderTo:conf.tabRenderTo
	})
	
	//渲染表格
	$(conf.tabRenderTo).html(fileTabTmp);
	var fileTab=$(conf.tabRenderTo).children('table');
	var ids=$(conf.valueId).val().split(',');
	for(var i in ids){
		$.get(fileupload.queryPath+"?id="+ids[i],{},function(data,status){
			if(status!='success'){
				alert("文件信息查询异常");
				return ;
			}
			
			var tr_render=template.compile(fileTrTmp);
			var tr=$(tr_render(data));
			tr.children('td:last').children('a').click({"options":file.options,"file":data}, fileupload_del);
			fileTab.children('tbody').append(tr);
		},'json');
	}
}

/*
 * 3.查看页面渲染。
 * jsp页面渲染好两个个标签
 * <input id="record" type="hidden" value="" />
 * <div id="files" class="files"> </div>
 * js在dom加载完毕后，调用下面的方法，完成页面的渲染
 * conf={
 * 		'valueId':'#importFile',//初始化的字段值存放组件ID
 * 		'tabRenderTo':'#files' //表格渲染到的容器ID
 * }
 */
fileupload.initShow=function(conf){
	//渲染表格
	$(conf.tabRenderTo).html(fileTabTmp);
	var fileTab=$(conf.tabRenderTo).children('table');
	var ids=$(conf.valueId).val().split(',');
	for(var i in ids){
		$.get(fileupload.queryPath+"?id="+ids[i],{},function(data,status){
			if(status!='success'){
				alert("文件信息查询异常");
				return ;
			}
			
			var tr_render=template.compile(fileTrTmp);
			var tr=$(tr_render(data));
			tr.children('td:last').html('<a class="downloadbtn" target="_blank" href="/'+data.saveName+'">下载</a>');
			fileTab.children('tbody').append(tr);
		},'json');
	}
}