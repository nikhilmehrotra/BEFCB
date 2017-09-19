$(document).ready(function(){
	ajaxPost("/web-cb/region/getdata", {}, function(result){})
	ajaxPost("/web-cb/region/get", {}, function(result){})
	ajaxPost("/web-cb/region/save", {}, function(data){})
	ajaxPost("/web-cb/region/delete", {RegionId:"59030edf8d5f772ad0b77693"}, function(data){})
})