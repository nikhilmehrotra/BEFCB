var ScorecardCountryAnalysis = {
	Processing:ko.observable(false),
	Data:ko.observableArray([]),
	// Filter
    Period:ko.observable(new Date(Now.getFullYear(),Now.getMonth(),1)),
	Country:ko.observable(""),
	CountryList:ko.observableArray([]), //If you need this, please get this data from ScorecardCountryAnalysis.Data
	TrendChartData:ko.observableArray([]),
	ActiveCountry:ko.observable(''),
    SortByValue:ko.observable(''),
    SortByValueType:ko.observable(''),
    UserCountry:ko.observable(''),
    OwnedCountry:ko.observable(false),
}

ScorecardCountryAnalysis.Period.subscribe(function(v){
    var DataReference = ko.mapping.toJS(ScorecardAnalysis.DataReference);
    ScorecardAnalysis.Processing(true)
    var url = "/web-cb/scorecardanalysis/getdata";
    v.setDate(1)
    // console.log(v,getUTCDate(v),kendo.toString(getUTCDate(v),"yyyyMMdd"))
    var parm = {
        BusinessMetrics:DataReference.Id,
        Region:ScorecardAnalysis.Region(),
        Country:ScorecardAnalysis.Country(),
        Period:kendo.toString(v,"yyyyMMdd"),
    };
    ajaxPost(url,parm,function(res){
        ScorecardAnalysis.Processing(false)
        if (res.IsError){
            swal("Error!",res.Message,"error");
        }

        if(typeof res.Data !== "undefined" && res.Data !== null){
            ScorecardCountryAnalysis.Get(res.Data.CountryAnalysis);
            console.log(res.Data.CountryAnalysis)
        }
    });
})

// ScorecardCountryAnalysis.SortByValue.subscribe(function(v){
//     var z = ko.mapping.toJS(ScorecardCountryAnalysis.SortByValue())
//     console.log(v,z)
//     DataSource = ScorecardCountryAnalysis.Data()
//     DataSource = _.sortBy(DataSource,[v], ['asc'])
//     ScorecardCountryAnalysis.Data(DataSource)
// })

ScorecardCountryAnalysis.ChangeSort = function(v){
    var z = ko.mapping.toJS(ScorecardCountryAnalysis.SortByValue())
    var typ = ko.mapping.toJS(ScorecardCountryAnalysis.SortByValueType())
    DataSource = ScorecardCountryAnalysis.Data()
    DataSource = _.sortBy(DataSource,[v], ["asc"])
    if(z == v && typ == 'asc'){
        ScorecardCountryAnalysis.SortByValueType('desc')
        DataSource = DataSource.reverse()
    } else if(typ == 'desc'){
        ScorecardCountryAnalysis.SortByValueType('asc')
    }

    // console.log(ScorecardCountryAnalysis.SortByValue(), ScorecardCountryAnalysis.SortByValueType(), v, z, typ)
    ScorecardCountryAnalysis.Data(DataSource)
    ScorecardCountryAnalysis.SortByValue(v)
}

ScorecardCountryAnalysis.Get = function(DataSource){
	// dummy data ----------------------------------
	// var DataSource = [
	// 	{CountryName: "Indonesia", Actual: 1, Target: 1, Gap: 1},
	// 	{CountryName: "Singapore", Actual: 2, Target: 2, Gap: 0.25},
	// 	{CountryName: "Japan", Actual: 3, Target: 3, Gap: 0.35},
	// 	{CountryName: "United State of America", Actual: 4, Target: 4, Gap: 0.75},
	// 	{CountryName: "Rusia", Actual: 4, Target: 4, Gap: 0.25},
	// 	{CountryName: "China", Actual: 5, Target: 5, Gap: 0.65},
	// ];
	// ---------------------------------------------
    // console.log("--", DataSource)
	if (typeof DataSource !== "undefined"){
        ScorecardCountryAnalysis.SortByValue('Gap')
        ScorecardCountryAnalysis.SortByValueType('desc')
        DataSource = _.sortBy(DataSource,["Gap"], ["desc"])
        DataSource = DataSource.reverse()
        _.each(DataSource, function(v,i){
            DataSource[i].GaptoTarget = 100 - (v.ActualPercent*100);
            if(v.ActualPercent*100 < 0){
                DataSource[i].GaptoTarget = 100;
            } else if(v.ActualPercent*100 > 100){
                DataSource[i].GaptoTarget = 0;   
            }

            DataSource[i].LastYearGaptoTarget = 100 - (v.LastYearActualPercent*100);
            if(v.LastYearActualPercent*100 < 0){
                DataSource[i].LastYearGaptoTarget = 100;
            } else if(v.LastYearActualPercent*100 > 100){
                DataSource[i].LastYearGaptoTarget = 0;   
            }

            // DataSource[i].LastYearActual = (DataSource[i].LastYearActual !== undefined) ? DataSource[i].LastYearActual : 0;
            // DataSource[i].LastYearNAActual = (DataSource[i].LastYearNAActual !== undefined) ? DataSource[i].LastYearNAActual : true;
            // DataSource[i].LastYearActualPercent = (DataSource[i].LastYearActualPercent !== undefined) ? DataSource[i].LastYearActualPercent : 0;
            // DataSource[i].LastYearActualPercentage = (DataSource[i].LastYearActualPercentage !== undefined) ? DataSource[i].LastYearActualPercentage : 0;
            // DataSource[i].LastYearBudgetProjection = (DataSource[i].LastYearBudgetProjection !== undefined) ? DataSource[i].LastYearBudgetProjection : 0;
            // DataSource[i].LastYearGap = (DataSource[i].LastYearGap !== undefined) ? DataSource[i].LastYearGap : 0;
            // DataSource[i].LastYearNAActualPercent = (DataSource[i].LastYearNAActualPercent !== undefined) ? DataSource[i].LastYearNAActualPercent : true;
            // DataSource[i].LastYearNABudgetProjection = (DataSource[i].LastYearNABudgetProjection !== undefined) ? DataSource[i].LastYearNABudgetProjection : true;
            // DataSource[i].LastYearNAGap = (DataSource[i].LastYearNAGap !== undefined) ? DataSource[i].LastYearNAGap : true;
            // DataSource[i].LastYearNAProjection = (DataSource[i].LastYearNAProjection !== undefined) ? DataSource[i].LastYearNAProjection : true;
            // DataSource[i].LastYearNATarget = (DataSource[i].LastYearNATarget !== undefined) ? DataSource[i].LastYearNATarget : true;
            // DataSource[i].LastYearName = (DataSource[i].LastYearName !== undefined) ? DataSource[i].LastYearName : "";
            // DataSource[i].LastYearOverTarget = (DataSource[i].LastYearOverTarget !== undefined) ? DataSource[i].LastYearOverTarget : 0;
            // DataSource[i].LastYearProjection = (DataSource[i].LastYearProjection !== undefined) ? DataSource[i].LastYearProjection : 0;
            // DataSource[i].LastYearTarget = (DataSource[i].LastYearTarget !== undefined) ? DataSource[i].LastYearTarget : 0;
            // DataSource[i].LastYearRAG = (DataSource[i].LastYearRAG !== undefined) ? DataSource[i].LastYearRAG : "";
            if(DataSource[i].LastYearActual == 0){
                DataSource[i].YoY = 100;
            } else{
                DataSource[i].YoY = ((DataSource[i].Actual - DataSource[i].LastYearActual) / DataSource[i].LastYearActual )*100;
            }
        })
		ScorecardCountryAnalysis.Data(DataSource)
		return false;
	}
	// ScorecardCountryAnalysis.Render();
}

// ScorecardCountryAnalysis.Render = function(){
// 	var DataSource = ko.mapping.toJS(ScorecardCountryAnalysis.Data());
// 	// Start Render
// }

ScorecardCountryAnalysis.GetTrendChart = function(CountryName, index){
    if( ( ScorecardCountryAnalysis.OwnedCountry() && CountryName == ScorecardCountryAnalysis.UserCountry() ) || !ScorecardCountryAnalysis.OwnedCountry() ){
    	if(CountryName == ScorecardCountryAnalysis.ActiveCountry() ){
        	ScorecardCountryAnalysis.ActiveCountry('');
        	$("#TrendChart").html('');
        } else{
        	ScorecardCountryAnalysis.Processing(false)
            var url = "/web-cb/scorecardanalysis/getcountryanalysistrendlinedata";
            var parm = {
        		Country: CountryName,
        		BusinessMetrics: ScorecardAnalysis.DataReference().Id,
                Period:kendo.toString(getUTCDate(ScorecardAnalysis.DataReference().Period),"yyyyMMdd"),
        	};
            ajaxPost(url,parm,function(res){
                // console.log(res)
                ScorecardCountryAnalysis.Processing(false)
                if (res.IsError){
                    swal("Error!",res.Message,"error");
                }

                resEdited = []
                _.each(res.Data, function(v,i){
                    if(v.length > 0){
                        resEdited.push(v[0])
                    }
                })

         		ScorecardCountryAnalysis.ActiveCountry(CountryName)
                ScorecardCountryAnalysis.TrendChartData(resEdited)//dummy data
                ScorecardCountryAnalysis.CreateTrendChart(CountryName)
            });
        }
    }
}

ScorecardCountryAnalysis.CreateTrendChart = function(CountryName){
	$("#TrendChart").html('');
    var DataReference = ko.mapping.toJS(ScorecardAnalysis.DataReference);
    var SeriesDatas = [
        {
            name: "Budget",
            field: "Budget",
            categoryField: "Month",
            color: "#8497B0",
        },
        {
            name: "Actual (Cumulative)",
            field: "ActualCummulative",
            categoryField: "Month",
            color: "#FF9933",
        },
        {
            name: "Actual (MoM)",
            field: "ActualMoM",
            categoryField: "Month",
            color: "#BF9000",
        }
    ]

    if(DataReference.Type == "spot"){
        SeriesDatas = [
            {
                name: "Target",
                field: "Target",
                categoryField: "Month",
                color: "black",
            },
            {
                name: "Budget",
                field: "Budget",
                categoryField: "Month",
                color: "#8497B0",
            },
            {
                name: "Actual (Cumulative)",
                field: "ActualCummulative",
                categoryField: "Month",
                color: "#FF9933",
            },
            {
                name: "Actual (MoM)",
                field: "ActualMoM",
                categoryField: "Month",
                color: "#BF9000",
            }
        ]
    }

    $("#TrendChart").kendoChart({
        dataSource: ScorecardCountryAnalysis.TrendChartData(),
        title: {
            text: CountryName,
        },
        legend: {
            position: "bottom",
            font:chartFont,
        },
        seriesDefaults: {
            type: "line",
            style: "smooth",
     //        markers: {
		   //  	visible: false,
		  	// },
        },
        series: SeriesDatas,
        valueAxis: {
            labels: {
                template: "#:value#",
                font:chartFont,
            },
            majorGridLines: {
                visible: false
            }
        },
        categoryAxis: {
            majorGridLines: {
                visible: false
            },
            labels: {
                font:chartFont,
            },
        },
        tooltip: {
            visible: true,
            format: "{0}",
            template: "#= category #: #= value #",
            font:chartFont,
        }
    });
}