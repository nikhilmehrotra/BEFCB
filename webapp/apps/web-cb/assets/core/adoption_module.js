TabMenuValue = ko.observable('AdoptionModule')
var AdoptionModule = {
	Mode:ko.observable("MAIN"),
	Processing:ko.observable(false),	
	Geography:ko.observable("ALL"),
	GeographyList:ko.observableArray([
		{value:"ALL",name:"All"},
		{value:"GLOBAL",name:"Global"},
		{value:"COUNTRY",name:"Country"},
	]),
	Country:ko.observable(""),
	CountryList:ko.observableArray([]),
	Data:ko.observableArray([]),
	OwnedInitiative:ko.observableArray([]),
	Lock:ko.observable(false),
}
AdoptionModule.Get = function(){
	AdoptionModule.Mode("MAIN");
	AdoptionModule.GetData();
}
AdoptionModule.Country.subscribe(function(val){
	if(!AdoptionModule.Lock()){
		AdoptionModule.GetData();
	}
});
AdoptionModule.Geography.subscribe(function(val){
	AdoptionModule.Lock(true);
	AdoptionModule.Country("");
	AdoptionModule.Lock(false);
	AdoptionModule.GetData();
});
AdoptionModule.GetAnalysis = function(d){
	var InitiativeList = ko.mapping.toJS(AdoptionModule.Data());
	AdoptionModuleAnalysis.Get(d.parentData,InitiativeList);
}
AdoptionModule.GetGeographyData = function(a,b,c){
	if(a){
		return "Global";
	}else{
		if(c.length>0){
			if(typeof c == "string"){
				return c;
			}
			return c.join(", ");
		}else if(b.length>0){
			return b.join(", ");
		}else{
			return "Global";
		}
	}
}


AdoptionModule.GetData = function(){
	AdoptionModule.Processing(true);
	AdoptionModule.Data([]);
	var parm = {
		Geography:AdoptionModule.Geography(),
		Country:(AdoptionModule.Country()==='Country'?'':AdoptionModule.Country())
	}
	ajaxPost("/web-cb/adoptionmodule/getdata",parm,function(res){
		if(res.IsError){
			swal("Error",res.Message,"error")
			return false;
		}
		AdoptionModule.Processing(false);
		var sources = res.Data.Sources;
		AdoptionModule.OwnedInitiative(res.Data.OwnedInitiative);
		for(var x in sources){
			for(var z in sources[x].MetricData){
				sources[x].MetricData[z]["parentData"] = sources[x];
			}
		}
		AdoptionModule.Data(sources);
	})
}
AdoptionModule.Init = function(){
	AdoptionModule.GetCountryData();
	AdoptionModule.GetData();
}
AdoptionModule.GetCountryData = function(){
	ajaxPost("/web-cb/region/getdata",{},function(res){
		if(res.IsError){
			return false;
		}
		AdoptionModule.CountryList(res.Data);
	});
}
$(document).ready(function(){
	AdoptionModule.Init();
})