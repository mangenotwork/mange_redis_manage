//var Host = "http://192.168.0.9:8334"
var Host = window.location.protocol + "//" + window.location.host;

function errorTip(xhr, textStatus){
	//alert("请求失败");
    //console.log("请求失败");
	console.log(xhr);
    console.log(textStatus)
    window.parent.location.href="/notlogin";
}

/**  
 * ajax post提交  
 * @param url  
 * @param param  
 * @param datat 为html,json,text  
 * @param callback回调函数  
 * @return  
 */
function sendAjax_post(url, param, callback) {
    console.log(url);
    $.ajax({
        type: "post",
        url: Host+url,
        data: JSON.stringify(param),
        dataType: 'json',
        success: callback,
        error: function(xhr,textStatus) {
            errorTip(xhr, textStatus);
        }
    });
}


/**  
 * ajax get提交  
 * @param url  
 * @param param 
 * @param datat 为html,json,text  
 * @param callback回调函数  
 * @return  
 */
function sendAjax_get(url, callback) {
    $.ajax({
        type: "get",
        url: url,
        success: callback,
        error: function(xhr, textStatus) {
            errorTip(xhr, textStatus);
        }
    });
}

function sendAjax_get_tb(url,callback){
    $.ajax({
        type: "get",
        url: url,
        async: false,
        success: callback,
        error: function(xhr, textStatus) {
            errorTip(xhr, textStatus);
        }
    });
}


function sendAjax_get_test(url) {
    $.ajax({
        type: "get",
        url: url,
        success: function(data){
            console.log(data);
        },
        error: function(xhr, textStatus) {
            errorTip(xhr, textStatus);
        }
    });
}


//公用成功弹窗
function PostTip(url, param){
	sendAjax_post(url,param,function(data){
		console.log(data);
		if (data.code == 0){
			alert(data.data);
		}else{
			alert(data.mag);
		}
	})
}


var URL = function (){}
//获取key信息
URL.prototype.GetKeyInfo = function(rid,dbid,key){
    return "/api/redis/auth/keys/info/"+rid+"?db="+dbid+"&key="+key;
}
//获取db树型结构数据
URL.prototype.DBTreeInfo = function(rid,dbid,match){
    return "/api/redis/auth/keytree/"+rid+"?db="+dbid+"&match="+match;
}
//key搜索
URL.prototype.SearchKey = function(rid,dbid,match){
    return "/api/redis/auth/keysearch/"+rid+"?db="+dbid+"&match="+match;
}
//新增key
URL.prototype.AddKey = function(rid,dbid){
    return "/api/redis/auth/db/info/"+rid+"?db="+dbid;
}