var AdoptionModuleAnalysis = {
	Processing:ko.observable(false),
	InitiativeData:ko.observable(),
	InitiativeId:ko.observable(""),
	Country:ko.observable(""),
	Year:ko.observable(Now.getFullYear()),
	ActiveMetric:ko.observable(""),
	// Data Source
	GoLiveDate:ko.observable(),
	YearList:ko.observableArray([]),
	InitiativeList:ko.observableArray([]),
	MetricData:ko.observableArray([]),
	IsResultAvailable:ko.observable(false),
	Data:ko.observable(),
	Lock:ko.observable(false),
}

AdoptionModuleAnalysis.CheckOwnedAccess = function(){
    var id = AdoptionModuleAnalysis.InitiativeId();
    var sources = AdoptionModule.OwnedInitiative();
    
    if(sources.indexOf(id) >= 0 ){
        return true;
    }else{
        return false;
    }
}
AdoptionModuleAnalysis.ExportPDF = function(){
	kendo.pdf.defineFont({
	  "Helvetica Neue": "/web-cb/static/fonts/HelveticaNeue.ttf"
	});
	kendo.drawing.drawDOM($('#pdf-man')).then(function (group) {
	  var title = "AdoptionModuleAnalysis.pdf"
	  kendo.drawing.pdf.saveAs(group, title);
	//       $('.btn-export-png').show()
	// $('.btn-export-pdf').show()
	// $('#ovheader').hide()
	})
}
AdoptionModuleAnalysis.InitiativeId.subscribe(function(val){
	var InitiativeData = Enumerable.From(AdoptionModuleAnalysis.InitiativeList()).Where("$.InitiativeID == '"+val+"'").FirstOrDefault();
	AdoptionModuleAnalysis.InitiativeData(InitiativeData);
	if(!AdoptionModuleAnalysis.Lock()){
		AdoptionModuleAnalysis.GetData();
	}
})
AdoptionModuleAnalysis.Year.subscribe(function(){
	if(!AdoptionModuleAnalysis.Lock()){
		AdoptionModuleAnalysis.GetData();
	}
})
AdoptionModuleAnalysis.Back = function(){
	AdoptionModule.Mode("MAIN");
	AdoptionModule.GetData();
	AdoptionModuleAnalysis.ActiveMetric("");
}
AdoptionModuleAnalysis.Get = function(InitiativeData,InitiativeList){
	AdoptionModuleAnalysis.Lock(true);
	AdoptionModuleAnalysis.Processing(true);
	AdoptionModule.Mode("ANALYSIS");
	AdoptionModuleAnalysis.Year(Now.getFullYear());
	AdoptionModuleAnalysis.InitiativeId(InitiativeData.InitiativeID);
	AdoptionModuleAnalysis.InitiativeData(InitiativeData);
	for(var x in InitiativeList){
		delete InitiativeList[x]["MetricData"]
	}
	AdoptionModuleAnalysis.InitiativeList(InitiativeList);
	AdoptionModuleAnalysis.Lock(false);
	AdoptionModuleAnalysis.GetData();
}
AdoptionModuleAnalysis.ResetFilterCountry = function(){
	AdoptionModuleAnalysis.Country("");
	AdoptionModuleAnalysis.GetDataCountry();
}
AdoptionModuleAnalysis.FilterCountry = function(d){
	AdoptionModuleAnalysis.Country(d.Country);
	AdoptionModuleAnalysis.GetDataCountry();
}
AdoptionModuleAnalysis.GetRAGValue = function(val){
	return val;
}
AdoptionModuleAnalysis.GetRAGValueFromNumber = function(val){
	val = parseInt(val)
	color = (val == 1) ? 'red ' : ( (val==2)?'amber ': ((val==3)?'green ':'') );
	return color;
}
AdoptionModuleAnalysis.GetValue = function(isNA,val){
	if(isNA){
		return "-";
	}
	return kendo.toString(val,"N1");
}
AdoptionModuleAnalysis.GetDetail = function(){
	AdoptionModuleDetail.Get(ko.mapping.toJS(AdoptionModuleAnalysis));
}
AdoptionModuleAnalysis.Render = function(dataSource){
	var d = dataSource.AnalyticData;
	var NASumData = dataSource.NASumData;
	var AnalyticDataSource = []


	var xdata = {Period:"Jan"}
	if(!NASumData.NASumJan){
		xdata.MoM = d.SumJan;
		xdata.Cumulative = d.CumulativeJan;
	}
	if(!NASumData.NASumTotalJan){
		xdata.DisplayTotal = d.SumTotalJan > d.SumJan ? (d.SumTotalJan-d.SumJan) : d.SumTotalJan;
		xdata.Total = d.SumTotalJan;
		xdata.CumulativeTotal = d.CumulativeTotalJan;
		// console.log(xdata)
	}
	AnalyticDataSource.push(xdata)


	// if(!NASumData.NASumJan){		
	// 	AnalyticDataSource.push({Period:"Jan",MoM:d.SumJan,DisplayTotal:(d.SumTotalJan-d.SumJan),Total:d.SumTotalJan,Cumulative:d.CumulativeJan,CumulativeTotal:d.CumulativeTotalJan})
	// }else{
	// 	AnalyticDataSource.push({Period:"Jan"})
	// }
	var xdata = {Period:"Feb"}
	if(!NASumData.NASumFeb){
		xdata.MoM = d.SumFeb;
		xdata.Cumulative = d.CumulativeFeb;
	}
	if(!NASumData.NASumTotalFeb){
		xdata.DisplayTotal = d.SumTotalFeb > d.SumFeb ? (d.SumTotalFeb-d.SumFeb) : d.SumTotalFeb;
		xdata.Total = d.SumTotalFeb;
		xdata.CumulativeTotal = d.CumulativeTotalFeb;
		// console.log(xdata)
	}
	AnalyticDataSource.push(xdata)


	// if(!NASumData.NASumFeb){	
	// 	AnalyticDataSource.push({Period:"Feb",MoM:d.SumFeb,DisplayTotal:(d.SumTotalFeb-d.SumFeb),Total:d.SumTotalFeb,Cumulative:d.CumulativeFeb,CumulativeTotal:d.CumulativeTotalFeb})
	// }else{
	// 	AnalyticDataSource.push({Period:"Feb"})
	// }
	var xdata = {Period:"Mar"}
	if(!NASumData.NASumMar){
		xdata.MoM = d.SumMar;
		xdata.Cumulative = d.CumulativeMar;
	}
	if(!NASumData.NASumTotalMar){
		xdata.DisplayTotal = d.SumTotalMar > d.SumMar ? (d.SumTotalMar-d.SumMar) : d.SumTotalMar;
		xdata.Total = d.SumTotalMar;
		xdata.CumulativeTotal = d.CumulativeTotalMar;
		// console.log(xdata)
	}
	AnalyticDataSource.push(xdata)


	// if(!NASumData.NASumMar){	
	// 	AnalyticDataSource.push({Period:"Mar",MoM:d.SumMar,DisplayTotal:(d.SumTotalMar-d.SumMar),Total:d.SumTotalMar,Cumulative:d.CumulativeMar,CumulativeTotal:d.CumulativeTotalMar})
	// }else{
	// 	AnalyticDataSource.push({Period:"Mar"})
	// }
	var xdata = {Period:"Apr"}
	if(!NASumData.NASumApr){
		xdata.MoM = d.SumApr;
		xdata.Cumulative = d.CumulativeApr;
	}
	if(!NASumData.NASumTotalApr){
		xdata.DisplayTotal = d.SumTotalApr > d.SumApr ? (d.SumTotalApr-d.SumApr) : d.SumTotalApr;
		xdata.Total = d.SumTotalApr;
		xdata.CumulativeTotal = d.CumulativeTotalApr;
		// console.log(xdata)
	}
	AnalyticDataSource.push(xdata)


	// if(!NASumData.NASumApr){	
	// 	AnalyticDataSource.push({Period:"Apr",MoM:d.SumApr,DisplayTotal:(d.SumTotalApr-d.SumApr),Total:d.SumTotalApr,Cumulative:d.CumulativeApr,CumulativeTotal:d.CumulativeTotalApr})
	// }else{
	// 	AnalyticDataSource.push({Period:"Apr"})
	// }
	var xdata = {Period:"May"}
	if(!NASumData.NASumMay){
		xdata.MoM = d.SumMay;
		xdata.Cumulative = d.CumulativeMay;
	}
	if(!NASumData.NASumTotalMay){
		xdata.DisplayTotal = d.SumTotalMay > d.SumMay ? (d.SumTotalMay-d.SumMay) : d.SumTotalMay;
		xdata.Total = d.SumTotalMay;
		xdata.CumulativeTotal = d.CumulativeTotalMay;
		// console.log(xdata)
	}
	AnalyticDataSource.push(xdata)


	// if(!NASumData.NASumMay){	
	// 	AnalyticDataSource.push({Period:"May",MoM:d.SumMay,DisplayTotal:(d.SumTotalMay-d.SumMay),Total:d.SumTotalMay,Cumulative:d.CumulativeMay,CumulativeTotal:d.CumulativeTotalMay})
	// }else{
	// 	AnalyticDataSource.push({Period:"May"})
	// }
	var xdata = {Period:"Jun"}
	if(!NASumData.NASumJun){
		xdata.MoM = d.SumJun;
		xdata.Cumulative = d.CumulativeJun;
	}
	if(!NASumData.NASumTotalJun){
		xdata.DisplayTotal = d.SumTotalJun > d.SumJun ? (d.SumTotalJun-d.SumJun) : d.SumTotalJun;
		xdata.Total = d.SumTotalJun;
		xdata.CumulativeTotal = d.CumulativeTotalJun;
		// console.log(xdata)
	}
	AnalyticDataSource.push(xdata)


	// if(!NASumData.NASumJun){	
	// 	AnalyticDataSource.push({Period:"Jun",MoM:d.SumJun,DisplayTotal:(d.SumTotalJun-d.SumJun),Total:d.SumTotalJun,Cumulative:d.CumulativeJun,CumulativeTotal:d.CumulativeTotalJun})
	// }else{
	// 	AnalyticDataSource.push({Period:"Jun"})
	// }
	var xdata = {Period:"Jul"}
	if(!NASumData.NASumJul){
		xdata.MoM = d.SumJul;
		xdata.Cumulative = d.CumulativeJul;
	}
	if(!NASumData.NASumTotalJul){
		xdata.DisplayTotal = d.SumTotalJul > d.SumJul ? (d.SumTotalJul-d.SumJul) : d.SumTotalJul;
		xdata.Total = d.SumTotalJul;
		xdata.CumulativeTotal = d.CumulativeTotalJul;
		// console.log(xdata)
	}
	AnalyticDataSource.push(xdata)


	// if(!NASumData.NASumJul){	
	// 	AnalyticDataSource.push({Period:"Jul",MoM:d.SumJul,DisplayTotal:(d.SumTotalJul-d.SumJul),Total:d.SumTotalJul,Cumulative:d.CumulativeJul,CumulativeTotal:d.CumulativeTotalJul})
	// }else{
	// 	AnalyticDataSource.push({Period:"Jul"})
	// }
	var xdata = {Period:"Aug"}
	if(!NASumData.NASumAug){
		xdata.MoM = d.SumAug;
		xdata.Cumulative = d.CumulativeAug;
	}
	if(!NASumData.NASumTotalAug){
		xdata.DisplayTotal = d.SumTotalAug > d.SumAug ? (d.SumTotalAug-d.SumAug) : d.SumTotalAug;
		xdata.Total = d.SumTotalAug;
		xdata.CumulativeTotal = d.CumulativeTotalAug;
		// console.log(xdata)
	}
	AnalyticDataSource.push(xdata)


	// if(!NASumData.NASumAug){	
	// 	AnalyticDataSource.push({Period:"Aug",MoM:d.SumAug,DisplayTotal:(d.SumTotalAug-d.SumAug),Total:d.SumTotalAug,Cumulative:d.CumulativeAug,CumulativeTotal:d.CumulativeTotalAug})
	// }else{
	// 	AnalyticDataSource.push({Period:"Aug"})
	// }
	var xdata = {Period:"Sep"}
	if(!NASumData.NASumSep){
		xdata.MoM = d.SumSep;
		xdata.Cumulative = d.CumulativeSep;
	}
	if(!NASumData.NASumTotalSep){
		xdata.DisplayTotal = d.SumTotalSep > d.SumSep ? (d.SumTotalSep-d.SumSep) : d.SumTotalSep;
		xdata.Total = d.SumTotalSep;
		xdata.CumulativeTotal = d.CumulativeTotalSep;
		// console.log(xdata)
	}
	AnalyticDataSource.push(xdata)


	// if(!NASumData.NASumSep){	
	// 	AnalyticDataSource.push({Period:"Sep",MoM:d.SumSep,DisplayTotal:(d.SumTotalSep-d.SumSep),Total:d.SumTotalSep,Cumulative:d.CumulativeSep,CumulativeTotal:d.CumulativeTotalSep})
	// }else{
	// 	AnalyticDataSource.push({Period:"Sep"})
	// }
	var xdata = {Period:"Oct"}
	if(!NASumData.NASumOct){
		xdata.MoM = d.SumOct;
		xdata.Cumulative = d.CumulativeOct;
	}
	if(!NASumData.NASumTotalOct){
		xdata.DisplayTotal = d.SumTotalOct > d.SumOct ? (d.SumTotalOct-d.SumOct) : d.SumTotalOct;
		xdata.Total = d.SumTotalOct;
		xdata.CumulativeTotal = d.CumulativeTotalOct;
		// console.log(xdata)
	}
	AnalyticDataSource.push(xdata)


	// if(!NASumData.NASumOct){	
	// 	AnalyticDataSource.push({Period:"Oct",MoM:d.SumOct,DisplayTotal:(d.SumTotalOct-d.SumOct),Total:d.SumTotalOct,Cumulative:d.CumulativeOct,CumulativeTotal:d.CumulativeTotalOct})
	// }else{
	// 	AnalyticDataSource.push({Period:"Oct"})tricNam
	// }
	var xdata = {Period:"Nov"}
	if(!NASumData.NASumNov){
		xdata.MoM = d.SumNov;
		xdata.Cumulative = d.CumulativeNov;
	}
	if(!NASumData.NASumTotalNov){
		xdata.DisplayTotal = d.SumTotalNov > d.SumNov ? (d.SumTotalNov-d.SumNov) : d.SumTotalNov;
		xdata.Total = d.SumTotalNov;
		xdata.CumulativeTotal = d.CumulativeTotalNov;
		// console.log(xdata)
	}
	AnalyticDataSource.push(xdata)


	// if(!NASumData.NASumNov){	
	// 	AnalyticDataSource.push({Period:"Nov",MoM:d.SumNov,DisplayTotal:(d.SumTotalNov-d.SumNov),Total:d.SumTotalNov,Cumulative:d.CumulativeNov,CumulativeTotal:d.CumulativeTotalNov})
	// }else{
	// 	AnalyticDataSource.push({Period:"Nov"})
	// }
	var xdata = {Period:"Dec"}
	if(!NASumData.NASumDec){
		xdata.MoM = d.SumDec;
		xdata.Cumulative = d.CumulativeDec;
	}
	if(!NASumData.NASumTotalDec){
		xdata.DisplayTotal = d.SumTotalDec > d.SumDec ? (d.SumTotalDec-d.SumDec) : d.SumTotalDec;
		xdata.Total = d.SumTotalDec;
		xdata.CumulativeTotal = d.CumulativeTotalDec;
		// console.log(xdata)
	}
	AnalyticDataSource.push(xdata)


	// if(!NASumData.NASumDec){	
	// 	AnalyticDataSource.push({Period:"Dec",MoM:d.SumDec,DisplayTotal:(d.SumTotalDec-d.SumDec),Total:d.SumTotalDec,Cumulative:d.CumulativeDec,CumulativeTotal:d.CumulativeTotalDec})
	// }else{
	// 	AnalyticDataSource.push({Period:"Dec"})
	// }
	// var chartwidth = $("#analysis-chart-width").width();
	chartwidth = 1140;
	var fontStyle = "11px Helvetica Neue";
	var fontStyleHeading = "14px Arial,Helvetica,sans-serif";
	// $("#analysis-table").width((chartwidth-53));
	// console.log(AnalyticDataSource)
	var SelectedMetric = "";

	if(AdoptionModuleAnalysis.ActiveMetric()!==""){
		var tmp = Enumerable.From(dataSource.MetricData).Where("$.MetricId === '"+AdoptionModuleAnalysis.ActiveMetric()+"'").FirstOrDefault()
		
		if(typeof tmp !== "undefined"){
			SelectedMetric = tmp.MetricName;	
		}
		
	}
	$("#analysis-chart").kendoChart({
		dataSource: {
            data: AnalyticDataSource,
        },
        title:{
        	text:AdoptionModuleAnalysis.InitiativeData().Initiatives+" - "+SelectedMetric+" "+(AdoptionModuleAnalysis.Country()!==""?" ( "+AdoptionModuleAnalysis.Country()+" )":""),
        	font:fontStyleHeading
        },
        legend: {
            visible: false
        },
        chartArea: {
		    width: chartwidth,
		    margin:{
		    	left:120,
		    	right:0
		    }
		},
        seriesDefaults: {
        	overlay: {
		      gradient: "none"
		    },
		    tooltip: {
		        visible: true,
		        format:"{0:N2}"
		    },
            labels: {
                visible: false,
            },
        },
        series: [
        	{field: "MoM",type: "column",axis: "MoM",color:"#2890C0",stack:true},
        	{field: "DisplayTotal",type: "column",axis: "MoM",color:"#a1c5e0",stack:true,tooltip:{template:'#:kendo.toString(dataItem.Total,"N2")#'}},
        	{field: "Cumulative",type: "line",axis: "Cumulative",color:"#3F9C35",style:"smooth",markers:{visible:false}},
        	{field: "CumulativeTotal",type: "line",axis: "Cumulative",color:"#9fd18b",style:"smooth",markers:{visible:false}},
        ],
        categoryAxis: {
            field: "Period",
            axisCrossingValues: [0, 12],
            majorGridLines: {visible:false},
            minorGridLines: {visible:false},
            labels:{
            	font: fontStyle
            },
            line:{
            	visible:false
            }
        },
        valueAxes: [
         	{
	            name: "MoM",
	            title: { text: "Month on Month",font: fontStyleHeading },
	            majorGridLines: {visible:false},
            	minorGridLines: {visible:false},
            	labels:{
	            	font: fontStyle
	            },
	            line:{
	            	visible:false
	            }
	        }, {
	            name: "Cumulative",
	            title: { text: "Cumulative",font: fontStyleHeading},
	            majorGridLines: {visible:false},
            	minorGridLines: {visible:false},
            	labels:{
	            	font: fontStyle
	            },
	            line:{
	            	visible:false
	            }
	        },
        ],
	});
}
AdoptionModuleAnalysis.GetDataCountry = function(){
	AdoptionModuleAnalysis.Processing(true);
	var parm = {
		InitiativeId:AdoptionModuleAnalysis.InitiativeId(),
		Year:AdoptionModuleAnalysis.Year(),
		Country:AdoptionModuleAnalysis.Country(),
		MetricId:AdoptionModuleAnalysis.ActiveMetric(),
	}
	ajaxPost("/web-cb/adoptionmodule/getanalysisdata",parm,function(res){
		// console.log(res.Data)
		if(res.IsError){
			swal("Error",res.Message,"error")
			return false;
		}
		AdoptionModuleAnalysis.ActiveMetric(res.Data.ActiveMetric);
		AdoptionModuleAnalysis.IsResultAvailable(res.Data.IsResultAvailable);
		if(res.Data.IsResultAvailable){
			AdoptionModuleAnalysis.Render(res.Data);
		}
		AdoptionModuleAnalysis.Processing(false);
		$("body").scrollTop(0);
	})
}
AdoptionModuleAnalysis.ChangeMetric = function(d){
	AdoptionModuleAnalysis.ActiveMetric(d.MetricId);
	AdoptionModuleAnalysis.GetData();
}
AdoptionModuleAnalysis.GetData = function(){
	AdoptionModuleAnalysis.Country("");
	AdoptionModuleAnalysis.Processing(true);
	var parm = {
		InitiativeId:AdoptionModuleAnalysis.InitiativeId(),
		Year:parseInt(AdoptionModuleAnalysis.Year()),
		Country:AdoptionModuleAnalysis.Country(),
		MetricId:AdoptionModuleAnalysis.ActiveMetric(),
	}
	ajaxPost("/web-cb/adoptionmodule/getanalysisdata",parm,function(res){
		// console.log(res.Data)
		if(res.IsError){
			swal("Error",res.Message,"error")
			return false;
		}

		AdoptionModuleAnalysis.ActiveMetric(res.Data.ActiveMetric);
		AdoptionModuleAnalysis.IsResultAvailable(res.Data.IsResultAvailable);
		AdoptionModuleAnalysis.MetricData(res.Data.MetricData);
		if(res.Data.IsResultAvailable){
			AdoptionModuleAnalysis.Data(res.Data);
			AdoptionModuleAnalysis.Render(res.Data);
		}
		var GoLiveDate = res.Data.GoLiveDate;
		var SelectedYear = parseInt(AdoptionModuleAnalysis.Year())
		for(var x in GoLiveDate){
			var LiveDate = getUTCDate(GoLiveDate[x]);
			var Year = LiveDate.getFullYear();
			if(SelectedYear==Year){
				var z = $("#analysis-table tr[country='"+x+"'] td[month='"+kendo.toString(LiveDate,"MMM")+"']")
				z.css("border","2px dashed #4d94ff")
				z.css("padding","4px;")
				var z = $("#analysis-table tr[country='"+x+"'] td[monthtotal='"+kendo.toString(LiveDate,"MMM")+"']")
				z.css("border","2px dashed #4d94ff")
				z.css("padding","4px;")
			}
		}
		AdoptionModuleAnalysis.Processing(false);
	})
}
AdoptionModuleAnalysis.Init = function(){
	// Get Year List
	var Last10Years = AdoptionModuleAnalysis.Year()-10;
	var NextYears = AdoptionModuleAnalysis.Year()+1;
	for (var i=Last10Years;i<=NextYears;i++){
		AdoptionModuleAnalysis.YearList.push(i);
	}

}
$(document).ready(function(){
	AdoptionModuleAnalysis.Init();
})