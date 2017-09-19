var ScorecardMetricCountryRanking = {
	Processing:ko.observable(false),
	Data:ko.observableArray([]),
	// Filter
	Country:ko.observable(""),
	ListMetric:ko.observableArray([]),
	CountryList:ko.observableArray([]), //If you need this, please get this data from ScorecardMetricCountryRanking.Data
	RegionCountryScorecard:ko.observable(false),
	RegionOne:ko.observable(""),
	CountryOne:ko.observable(""),
	SortBy:ko.observable("gaptotarget"),
}
// $(document).ready(function(){
// 	ajaxPost("/web-cb/masterregion/getdata",{},function(res){
// 		ScorecardMetricCountryRanking.CountryList(res.Data);
// 	})
// })

ScorecardMetricCountryRanking.SortBy.subscribe(function(v){
	// console.log(v)
	ScorecardAnalysis.GetMetricsRankingData()
	// parm:{
	// 	SortBy: ScorecardMetricCountryRanking.SortBy()
	// }
	// ajaxPost("/web-cb/masterregion/getdata",parm,function(res){
	// })
})

ScorecardMetricCountryRanking.Get = function(){

	ScorecardMetricCountryRanking.ListMetric([])
	ScorecardMetricCountryRanking.ListMetric.push({name: "Rank", id: "Rank", Color: "#A1AFC2"})
	_.each(ScorecardAnalysis.MetricsCountryRankingData(), function(v,i){
		if(v.Type == "cumulative"){
			// v.Color = "rgb(106,193,123)"
			v.Color = "rgb(106,193,123)"
		}else if(v.Type == "relative"){
			// v.Color = "rgb(105,190,40)"
			v.Color = "rgb(40,144,192)"
		}else if(v.Type == "spot"){
			// v.Color = "rgb(0,92,132)"
			v.Color = "rgb(109,110,103)"
		}else{
			v.Color = "#A1AFC2"
		}
	    ScorecardMetricCountryRanking.ListMetric.push({name: v.Name, Color: v.Color, DataPoint: v.Name, Type: v.Type, Id: v.Id, LastPeriod: toUTC(kendo.parseDate(v.Period,"yyyyMMdd"))})
	})

	DataSource = ScorecardAnalysis.MetricsCountryRankingData();
	var NewData = []

	if(ScorecardMetricCountryRanking.ListMetric().length > 1){
		var Higher = 0;
		_.each(DataSource, function(vv,ii){
			var MaxLength = vv.DataList.length
			if(MaxLength > Higher){
				Higher = MaxLength
			}
		})
		for(i = 1; i <= Higher; i++){
			var DataInside = []
			DataInside.push({CountryName: "",CountryCode: "",Region: "" , Rank: i})
			_.each(DataSource, function(vv,ii){	
				data = _.find(vv.DataList, function(vvv){return vvv.Rank == i});
				if(data != undefined){
	                DataInside.push(data)
	            } else{
	                DataInside.push({CountryName: "",CountryCode: "",Region: "" , Rank: i})
	            }
			})
			NewData.push({Data :DataInside})
		}
		
	}
	DataSource = NewData

	ScorecardMetricCountryRanking.Data(DataSource)
}

ScorecardMetricCountryRanking.CountryValue = ko.observable("")
ScorecardMetricCountryRanking.CountryList = ko.observableArray(["INA","SG"])
ScorecardMetricCountryRanking.CountryValue.subscribe(function(d){
	if(d != ""){
		$('.tdmodified').removeClass('mark')
		$('.tdmodified.'+d).addClass('mark')
	} else{
		$('.tdmodified').removeClass('mark')
	}
})