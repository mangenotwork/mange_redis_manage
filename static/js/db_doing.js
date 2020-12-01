var db_id = 0;
var global_key = "";
var hash_value_id = 1;
var list_value_id = 1;
var set_value_id = 1;
var zset_value_id = 1;
var nowe_key = "";


//获取key的信息并更新页面
function KeysInfo(dbid,key){
	var url = new URL();
	sendAjax_get_tb(url.GetKeyInfo(rid,dbid,key),function(data){
		db_id = dbid;
		global_key = key;
		nowe_key = data.data.key_name;

		$("#keyinfo").show()
		$("#addkey_div").hide();
		$("#stringkey_btn").hide();
		$("#key_btn").hide();
		$("#keyname").empty();
		$("#keyname").append(data.data.key_name);
		$("#key_db").empty();
		$("#key_db").append(data.data.key_db);
		$("#ttl").empty();
		$("#ttl").append(data.data.ttl);
		$("#key_type").empty();
		$("#key_type").append(data.data.key_type);
		$("#key_size").empty();
		$("#key_size").append(data.data.size);
		$("#value_code").empty();

		if (data.data.key_type == "string"){
			var value_data = JSON.stringify(data.data.value)
			$("#value_code").append(value_data);
			$("#txtdata").show();
			$(".datagrid").hide();
			$("#stringkey_btn").show();
		}else{
			//console.log(data.data.value);
			//console.log("typeof = ", typeof(data.data.value));
			var columns_val = [];
			var data_val = Object;
			var data_show = [];

			//hash
			if (data.data.key_type == "hash"){
				columns_val = [
					{field:'k',title:'Field',sortable:true,width:300},
					{field:'v',title:'Value',sortable:true,width:300},
				];
				data_val = JSON.parse(data.data.value);
					//console.log("data_val = ",data_val, typeof(data_val));		
				for(var a in data_val){
					//console.log(a, data_val[a]);
					data_show.push({"k":a,"v":data_val[a]});
				}
			}

			//list, set
			if (data.data.key_type == "list" || data.data.key_type == "set"){
				data_val = data.data.value;
				columns_val = [
					{field:'v',title:'Value',sortable:true,width:600},
				];
				for(var a in data_val){
					data_show.push({"v":data_val[a]});
				}
			}

			//zset
			if (data.data.key_type == "zset"){
				data_val = data.data.value;
				columns_val = [
					{field:'k',title:'Value',sortable:true,width:300},
					{field:'v',title:'Score',sortable:true,width:300},
				];
				var show_1 = {};
				for(var a in data_val){
					//console.log(data_val[a])
					if (a%2 === 0){
						show_1.k = data_val[a];
					}else{
						show_1.v = data_val[a];
						//console.log(show_1)
						data_show.push(show_1);
						show_1 = {};
					}
				}
			}

			$('#listdata').datagrid({
				singleSelect: true,
				columns:[columns_val],
					data: data_show
			});

			$("#txtdata").hide();
			$(".datagrid").show();
			$("#key_btn").show();
		}

		//修改key名称
		$(".now_dbid").empty();
		$(".now_dbid").append(dbid)
		$(".now_key").empty();
		$(".now_key").append(data.data.key_name)
	});
}

//获取key树的信息并更新页面
function DBTree(dbid,match){
	var url = new URL();
    sendAjax_get(url.DBTreeInfo(rid,dbid,match),function(data){
        var init_data = data.data;
        var newdata = []

        $('#tt'+dbid).tree({
	    	data: init_data,
	    	onBeforeExpand: function(node){
	    		//展开前
	    		// 加载数据，加载key
	    		//console.log(node);
	    		//console.log(node.target);
				//console.log(node.text); 
				var reg = `\\(.*?\\)`
				text = node.text.replace(new RegExp(reg, `ig`), ``);
				text = text+":";
				//console.log(text);

				//var urls = "/api/redis/auth/keytree/"+rid+"?db="+dbid+"&match="+text;
				sendAjax_get_tb(url.DBTreeInfo(rid,dbid,text),function(data){
				    //console.log(data.data);
				    newdata = data.data;
				    //console.log("in newdata = ",newdata);
				    // console.log(node.target);
				    if  (node.children === undefined){
					    $('#tt'+dbid).tree('append', {
							parent: node.target,
							data: newdata
						});
				    }
				});
			},
			onExpand: function(node){
			},
			onBeforeCollapse: function(node){
			},
			onClick: function(node){
				//单击展开
				//console.log(node);
				//console.log(node.children);
				if  (node.children === undefined && node.state === "open"){
					KeysInfo(dbid,node.text);
				}else{
					$('#tt'+dbid).tree('expand',node.target);
				}
	    	},
	    	onDblClick:function(node){
	    		//双击收起
	    		if  (node.children === undefined && node.state === "open"){
	    			KeysInfo(dbid,node.text);
	    		}else{
	    			$('#tt'+dbid).tree('collapse',node.target);
	    		}
	    	}
		});
    });
}

//点击db 事件
$(".opendb").click(function(){
	var dbid = $(this).attr("db");
	if ($("#"+dbid).html() == ""){
		$("#"+dbid).append('<ul id="tt'+$(this).attr("db")+'" style="width: 300px;"></ul>');
		DBTree($(this).attr("db"),"");
	}else{
		$("#"+dbid).empty();
	}
	
});

//搜索结果界面渲染
function SearchTree(dbid,match){
	var url = new URL();
    sendAjax_get(url.SearchKey(rid,dbid,match),function(data){
        var init_data = data.data;
        console.log(init_data)
        if (init_data === null){
        	init_data = [{
        		"text":"未找到含有"+match+"的key"
        	}];
        }
        $('#tt'+dbid).tree({
	    	data: init_data,
	    	onClick: function(node){
	    		if  (node.children === undefined && node.state === "open"){
					KeysInfo(dbid,node.text);
				}
	    	}
		});
    });
}

//搜索key 的输入事件
$(".search").keyup(function(e){
	var dbid = $(this).attr("db");
	$("#"+dbid).empty();
	$("#"+dbid).append('<ul id="tt'+$(this).attr("db")+'" style="width: 300px;"></ul>');
	SearchTree($(this).attr("db"),$(this).val());

})

//db树的刷新事件
$(".reload").click(function(){
	var dbid = $(this).attr("db");
	$("#"+dbid).empty();
	$("#"+dbid).append('<ul id="tt'+$(this).attr("db")+'" style="width: 300px;"></ul>');
	DBTree($(this).attr("db"),"");
})

//更新key的值界面显示点击事件
$("#update").click(function(){
	var v_code_txt = $("#value_code").html();
	$("#value_code").hide();
	$("#value_text").val(v_code_txt);
	$("#value_text").show();
	$("#submit").show();
})

//取消更新key的值界面隐藏点击事件
$("#notupdate").click(function(){
	$("#value_text").val();
	$("#value_text").hide();
	$("#value_code").show();
	$("#submit").hide();
})

//key重命名弹窗
$(".rename_key").click(function(){
	$("#rename_win").window({
		title:"key重命名",
		width:600,
    	height:400,
   	 	modal:true
	});
});

//key删除点击事件
$(".delete_key").click(function(){
	$("#delete_win").window({
		title:"删除再次确认",
		width:420,
    	height:210,
   	 	modal:true
	});
});

//刷新当前key
$("#refresh_key").click(function(){
	KeysInfo(db_id,global_key);
});

//重置当前key的ttl弹窗
$(".set_ttl").click(function(){
	$("#set_ttl_win").window({
		title:"重置TTL",
		width:600,
    	height:400,
   	 	modal:true
	});
});

//addkey dbinfo html
function addkey_dbinfo_html(dbid,keys,avgttl,rxp){
	return '<span style="font-size: 21px;">DBID: <span style="color: red;">'+dbid+'</span></span>'
		   +'<span style="font-size: 21px;margin-left: 50px">Yey总数: '+keys+'</span></span>'
		   +'<span style="font-size: 21px;margin-left: 50px">平均TTL: '+avgttl+'</span></span>'
		   +'<span style="font-size: 21px;margin-left: 50px">Expires: '+rxp+'</span></span>'
}

//addkey
$(".addkey").click(function(){
	var dbid = $(this).attr("db");
	//console.log(dbid)
	var url = new URL();
    sendAjax_get(url.AddKey(rid,dbid),function(data){
    	console.log(data);
    	$("#addkey_div").show();
    	$("#keyinfo").hide();
    	$("#addkey_dbinfo").empty();
    	$("#addkey_dbinfo").append(addkey_dbinfo_html(data.data.db, data.data.keys_count, data.data.avgttl, data.data.expires));

    });
});



//新增redis key value 交互，单行输入框增加
//Class OneInput(setid, names, divid)
var OneInput = function(setid, names, divid){
	this.names = names;
	this.value_id = setid;
	this.div = $("#"+divid);
}
OneInput.prototype.init_show = function(){
	this.div.empty();
	this.div.append('<input class="'+this.names+'_input addinput" type="text">'
					+'<a onclick="'+this.names+'_add()" class="addinputbtn"><span>+</span></a><br>');
	this.value_id = 1;
}
OneInput.prototype.html = function(){
	var id_str = this.names+"_id_"+this.value_id;
	console.log(this.input);
	return '<div id="'+id_str+'">'
			+'<input class="'+this.names+'_input addinput" type="text">'
			+'<a onclick="'+this.names+'_add()" class="addinputbtn"><span>+</span></a>'
			+'<a onclick="'+this.names+'_del('+id_str+')" class="addinputbtn"><span>-</span></a><br></div>'
}
OneInput.prototype.add = function(){
	this.value_id++
	this.div.append(this.html());
}
OneInput.prototype.del = function(o){
	$("#"+o.id).empty();
}
OneInput.prototype.reset = function(){
	this.init_show();
	this.resetskeyinput();
}
OneInput.prototype.html2input = function(){
	var id_str = this.names+"_id_"+this.value_id;
	return '<div id="'+id_str+'"><input class="'+this.names+'_input addinput" type="text">'
			+'<input class="'+this.names+'_input addinput" type="text">'
			+'<a onclick="'+this.names+'_add()" class="addinputbtn"><span>+</span></a>'
			+'<a onclick="'+this.names+'_del('+id_str+')" class="addinputbtn"><span>-</span></a><br></div>'
}
OneInput.prototype.init_show2 = function(){
	this.div.empty();
	this.div.append('<input  class="'+this.names+'_input addinput" type="text">'
					+'<input class="'+this.names+'_input addinput" type="text">'
					+'<a onclick="'+this.names+'_add()" class="addinputbtn"><span>+</span></a><br>');
	this.value_id = 1;
}
OneInput.prototype.add2 = function(){
	this.value_id++
	this.div.append(this.html2input());
}
OneInput.prototype.reset2 = function(){
	this.init_show2();
	this.resetskeyinput();
}
OneInput.prototype.post = function(result,keytype){
	
	var keyname = $("#"+this.names+"_key").val();
	var ttl = $("#"+this.names+"_ttl").val();

	if (result.length === 0){
		alert("Value不能为空");
		return
	}

	if (keyname === ""){
		alert("Key Name不能为空");
		return
	}

	if (ttl === ""){
		ttl = "-1"
	}

	var dbid = $("#addkey_dbid").html();

	var param = {
		"db_id":parseInt(dbid),
		"key":keyname,
		"key_type":keytype,
		"value":result,
		"ttl":parseInt(ttl)
	}

	console.log(param);

	var url = "/api/redis/auth/keys/create/{{.rid}}";
	sendAjax_post(url,param,function(data){
	    console.log(data);
	    if (data.code == 1){
	    	alert(data.mag);
	    	return
	    }
	    alert("key添加成功");
	    location.reload();

	})

}
OneInput.prototype.subimt2 = function(keytype){
	var result = new Array();
	$('input[class="'+this.names+'_input"]').each(function(j,k){
		if (j%2 === 0){
			key = k.value
		}else{
			val = k.value
			if (key != "" && val != ""){
				var map=new Map();
				map[key]=val;
				result.push(map);
			}
		}
	});
	this.post(result,keytype);
}
OneInput.prototype.subimt = function(keytype){
	var result = new Array();
	$('input[class="'+this.names+'_input"]').each(function(j,k){
		if (k.value != ""){
			result.push(k.value)
		}
	});
	this.post(result,keytype);
}
OneInput.prototype.stringsubimt = function(keytype){
	var result = $("#"+this.names+"_input").val();
	this.post(result,keytype);
}
OneInput.prototype.resetstring = function(){
	$("#"+this.names+"_input").val("");
	this.resetskeyinput();
}
OneInput.prototype.resetskeyinput = function(){
	$("#"+this.names+"_key").val("");
	$("#"+this.names+"_ttl").val("");
}


var string_add = new OneInput(1, "string_value", "");
function string_subimt(){
	string_add.stringsubimt("string");
}
function string_reset(){
	string_add.resetstring();
}

var set_add = new OneInput(set_value_id, "set_value", "add_set_list");
function set_value_add(){
	set_add.add();
}
function set_value_del(a){
	set_add.del(a);
}
function set_subimt(){
	set_add.subimt("set");
}
function set_reset(){
	set_add.reset();
}

var list_add = new OneInput(list_value_id, "list_value", "add_list_list");
function list_value_add(){
	list_add.add();
}
function list_value_del(a){
	list_add.del(a);
}
function list_subimt(){
	list_add.subimt("list");
}
function list_reset(){
	list_add.reset();
}

var hash_add = new OneInput(hash_value_id, "hash_value", "add_hash_list");
function hash_value_add(){
	hash_add.add2();
}
function hash_value_del(a){
	hash_add.del(a);
}
function hash_subimt(){
	hash_add.subimt2("hash");
}
function hash_reset(){
	hash_add.reset2();
}

var zset_add = new OneInput(zset_value_id, "zset_value", "add_zset_list");
function zset_value_add(){
	zset_add.add2();
}
function zset_value_del(a){
	zset_add.del(a);
}
function zset_subimt(){
	zset_add.subimt2("zset");
}
function zset_reset(){
	zset_add.reset2();
}


//更新当前value的提交按钮事件
function switchadddiv(divtype){

	$("#add_string_div").hide();
	$("#add_hash_div").hide();
	$("#add_list_div").hide();
	$("#add_set_div").hide();
	$("#add_zset_div").hide();
	$("#add_channel_div").hide();
	$("#add_import_div").hide();

	switch(divtype){
		case "string":
	    	$("#add_string_div").show();
	    	break;
		case "hash":
	    	$("#add_hash_div").show();
	    	hash_add.init_show2();
	    	break;
		case "list":
	    	$("#add_list_div").show();
	    	list_add.init_show();
	    	break;
		case "set":
	    	$("#add_set_div").show();
	    	set_add.init_show();
	    	break;
		case "zset":
	    	$("#add_zset_div").show();
	    	zset_add.init_show2();
	    	break;
	    case "channel":
	    	$("#add_channel_div").show();
	    	break;
	    case "import":
	    	$("#add_import_div").show();
	    	break;
	   default:
	    $("#add_string_div").show();
	  };
}


$('#listdata').datagrid({
	singleSelect: true,
	rownumbers: true,
	fitColumns: false,//设置为 true，则会自动扩大或缩小列的尺寸以适应网格的宽度并且防止水平滚动。
	striped: true,//设置为 true，则把行条纹化
	pagination:true,//设置为 true，则在数据网格（datagrid）底部显示分页工具栏。
	pageNumber: 1,//初始化页码。
	pageSize: 20,//初始化页面尺寸。
	pageList: [10,20,30,50],//初始化页面尺寸的选择列表。
	onDblClickRow: function(rowIndex, rowData){
		console.log(rowIndex);
		console.log(rowData);

	},
	onSelect: function(rowIndex, rowData){
		console.log(rowIndex);
		console.log(rowData);
	},
	toolbar: [{
		iconCls: 'icon-edit',
		handler: function(){
			var d = $("#listdata").datagrid('getSelections');
			console.log(d);


			$('#edit').window({
			    width:600,
				height:350,
				modal:true,
				collapsible:false,
				title: "客户端操作",
				minimizable:false,
				maximizable:false
			}); 

			for (var i in d){
				console.log(i);
				console.log(d[i]);
				$('#edit').append('k:<input class="edit_input" type="text" style="width:368px" value="'+d[i]["k"]+'"><br>');
				$('#edit').append('v:<input class="edit_input" type="text" style="width:368px" value="'+d[i]["v"]+'"><br>');
			}


		}
	},'-',{
		iconCls: 'icon-help',
		handler: function(){alert('help')}
	},'-',{
		iconCls: 'icon-cancel',
		handler: function(){
			var select = $("#listdata").datagrid('getSelections');
			console.log(select);
			if (select.length == 0) {
				alert("请选中要删除的行")
			}
		}
	}
	]
});


$('#ttl_dt').datetimebox({
    value: '20/9/2020 12:8:01',
    required: true,
    showSeconds: true
});

//修改key的ttl
function setttl(typeid) {
	var ttl = $("#ttl_dt").datetimebox('getValue');
	if (typeid == 0){
		ttl = $("#ttl_ss").val();
	} 

	console.log(nowe_key,db_id,typeid,ttl);
	var url = "/api/redis/auth/keys/modify/ttl/"+rid+"?db="+db_id+"&key="+nowe_key+"&type="+typeid+"&ttl="+ttl;
    sendAjax_get(url,function(data){
    	console.log(data);
    	if (data.code == 0){
    		alert("设置成功");
    	}else{
    		alert("设置失败:"+data.mag);
    	}

    });
}

//修改key的名称
function renamekey(){
	var keynewname = $("#key_new_name").val();
	console.log(nowe_key,db_id,keynewname);
	var url = "/api/redis/auth/keys/modify/name/"+rid+"?db="+db_id+"&key="+nowe_key+"&newname="+keynewname;
    sendAjax_get(url,function(data){
    	console.log(data);
    	if (data.code == 0){
    		alert("设置成功");
    	}else{
    		alert("设置失败:"+data.mag);
    	}

    });

}

//删除key
function delkey(){
	var url = "/api/redis/auth/keys/delete/"+rid+"?db="+db_id+"&key="+nowe_key;
    sendAjax_get(url,function(data){
    	console.log(data);
    	if (data.code == 0){
    		alert("删除成功");
    	}else{
    		alert("删除失败:"+data.mag);
    	}

    });
}


