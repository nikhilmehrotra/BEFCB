var ScorecardFullyearRanking = {
	Processing:ko.observable(false),
	Data:ko.observableArray([]),
	// Filter
	Country:ko.observable(""),
	CountryList:ko.observableArray([]), //If you need this, please get this data from ScorecardFullyearRanking.Data
	RegionCountryScorecard:ko.observable(false),
	RegionOne:ko.observable(""),
	CountryOne:ko.observable(""),
}
// $(document).ready(function(){
// 	ajaxPost("/web-cb/masterregion/getdata",{},function(res){
// 		ScorecardFullyearRanking.CountryList(res.Data);
// 	})
// })
ScorecardFullyearRanking.Get = function(DataSource){
	// console.log(DataSource)
	if (typeof DataSource !== "undefined" && DataSource.length > 0){		
		if(DataSource.length > 0){
			totalMonth = DataSource.length
			total = DataSource[0].length
			var newArr = []
			for(v=0; v<total; v++){
				var newArrDetail = []
				for(vv=0; vv<totalMonth; vv++){
					// console.log("zzz", DataSource[vv][v] )
					if(DataSource[vv][v].NAGap){
						DataSource[vv][v].CountryCode = "";
						DataSource[vv][v].Name = "";
						DataSource[vv][v].RegionName = "";
					}
					newArrDetail.push(DataSource[vv][v])
		        }
				newArr.push({"Count": v,"Data":newArrDetail})
		    }
			DataSource = newArr
		}

		ScorecardFullyearRanking.Data(DataSource)
		return false;
	}
	// ScorecardFullyearRanking.Render();
}
// ScorecardFullyearRanking.Render = function(){
// 	var DataSource = ko.mapping.toJS(ScorecardFullyearRanking.Data());
// 	// Start Render
// }

ScorecardFullyearRanking.CountryValue = ko.observable("")
ScorecardFullyearRanking.CountryList = ko.observableArray(["INA","SG"])
ScorecardFullyearRanking.CountryValue.subscribe(function(d){
	if(d != ""){
		$('.tdmodified').removeClass('mark')
		$('.tdmodified.'+d).addClass('mark')
	} else{
		$('.tdmodified').removeClass('mark')
	}
})