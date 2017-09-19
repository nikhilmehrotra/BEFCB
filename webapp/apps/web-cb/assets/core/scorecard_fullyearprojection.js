var ScorecardFullyearProjection = {
	Processing:ko.observable(false),
	Data:ko.observableArray([]),
	YearList:ko.observableArray(),
    MonthYearList: ko.observableArray(),
    BusinessMetricsList:ko.observableArray(),
    AfterRender:ko.observable(false)
}

ScorecardFullyearProjection.Init = function(){
	month = ["January", "February", "March", "April","May","June","July","August","September","October","November","December"]
    // Gen YearList
    var CurrentYear = new Date().getFullYear();
    for(var i = CurrentYear;i>=2010;i--){
    for(ii in month){
      ScorecardFullyearProjection.MonthYearList.push({Text: month[ii]+" "+i, Value: month[ii]+"-"+i})
    }
      ScorecardFullyearProjection.YearList.push(i);
    }
    ajaxPost("/web-cb/businessmetrics/getalldata",{},function(res){
		var arr = sortStringNumber(100, res.Data, "description")
		//Enumerable.From(res.Data).OrderBy("$.businessmetric.description").ToArray();
		    newarr1 = []
		    _.each(arr,function(v,i){
		        if(v.type == undefined){
		            v.type = ""
		        }
		        if(v.type == "cumulative"){
		            newarr1.push(v)
		        }
		    })

		ScorecardFullyearProjection.BusinessMetricsList(arr);
    })

	ScorecardFullyearProjection.Processing(false)
	var DataReference = ko.mapping.toJS(ScorecardAnalysis.DataReference);
	parm = {
		BMId : DataReference.Id,
		Breakdown : "country",
		Period : kendo.toString(getUTCDate(DataReference.LastPeriod),"MMMM-yyyy"),
		RelevantFilter : "",
	}
	ajaxPost("/web-cb/countryanalysis/getdata",parm,function(res){
		ScorecardFullyearProjection.Render(res.Data);
		ScorecardFullyearProjection.Processing(true)
	})	
}

ScorecardFullyearProjection.Get = function(dataSource){
    if(!ScorecardFullyearProjection.AfterRender()){
        ScorecardFullyearProjection.Init()
        ScorecardFullyearProjection.AfterRender(true)
    }
}

ScorecardFullyearProjection.Render = function(dataSource){

	var isNegative = false;
    _.each(dataSource, function(v,i){

      if(v.Projection >= 1){
        dataSource[i].ColorActualPercent = '#36bc9b'
      } else{
        dataSource[i].ColorActualPercent = '#FF0041'
      }

      if(v.BudgetProjection >= 1){
        dataSource[i].ColorBudgetProjection = '#36bc9b'
      } else{
        dataSource[i].ColorBudgetProjection = '#FF0041'
      }

      if(v.Projection < 0 || v.ActualPercent < 0){
        isNegative = true;
      }

      //for budget
      dataSource[i].RematchingBudgetProjection = 0;
      if(v.ActualPercent >= 0 && v.BudgetProjection >= 0){
        dataSource[i].RematchingBudgetProjection = v.BudgetProjection - v.ActualPercent; //projection
        if(dataSource[i].RematchingBudgetProjection < 0){
            dataSource[i].RematchingBudgetProjection = 0; //projection
        }
      } else if(v.ActualPercent < 0 && v.BudgetProjection < 0){
        dataSource[i].RematchingBudgetProjection = v.BudgetProjection - v.ActualPercent; //projection
        if(dataSource[i].RematchingBudgetProjection >= 0){
            dataSource[i].RematchingBudgetProjection = 0; //projection
        }
      } else if(v.ActualPercent >= 0 && v.BudgetProjection < 0){
        dataSource[i].RematchingBudgetProjection = v.BudgetProjection;
      } else if(v.ActualPercent < 0 && v.BudgetProjection >= 0){
        dataSource[i].RematchingBudgetProjection = v.BudgetProjection;
      }

      dataSource[i].RematchingProjectionforBudgetUse = 0;
      if(v.BudgetProjection >= 0 && v.Projection >= 0){
        dataSource[i].RematchingProjectionforBudgetUse = v.Projection - v.BudgetProjection; //projection
        if(dataSource[i].RematchingProjectionforBudgetUse < 0){
            dataSource[i].RematchingProjectionforBudgetUse = 0; //projection
        }
      } else if(v.BudgetProjection < 0 && v.Projection < 0){
        dataSource[i].RematchingProjectionforBudgetUse = v.Projection - v.BudgetProjection; //projection
        if(dataSource[i].RematchingProjectionforBudgetUse >= 0){
            dataSource[i].RematchingProjectionforBudgetUse = 0; //projection
        }
      } else if(v.BudgetProjection >= 0 && v.Projection < 0){
        dataSource[i].RematchingProjectionforBudgetUse = v.Projection;
      } else if(v.BudgetProjection < 0 && v.Projection >= 0){
        dataSource[i].RematchingProjectionforBudgetUse = v.Projection;
      }
      //

      dataSource[i].ProjectionMinusActualPercent = 0;  
      if(v.ActualPercent >= 0 && v.Projection >= 0){
        dataSource[i].ProjectionMinusActualPercent = v.Projection - v.ActualPercent; //projection
        if(dataSource[i].ProjectionMinusActualPercent < 0){
            dataSource[i].ProjectionMinusActualPercent = 0; //projection
        }
      } else if(v.ActualPercent < 0 && v.Projection < 0){
        dataSource[i].ProjectionMinusActualPercent = v.Projection - v.ActualPercent; //projection
        if(dataSource[i].ProjectionMinusActualPercent >= 0){
            dataSource[i].ProjectionMinusActualPercent = 0; //projection
        }
      } else if(v.ActualPercent >= 0 && v.Projection < 0){
        dataSource[i].ProjectionMinusActualPercent = v.Projection;
      } else if(v.ActualPercent < 0 && v.Projection >= 0){
        dataSource[i].ProjectionMinusActualPercent = v.Projection;
      }

      dataSource[i].RematchingGap = 0;
      if(v.ActualPercent >= 0 && v.Gap >= 0){
        // dataSource[i].RematchingGap = v.Gap - v.ActualPercent; //gap
        // if(dataSource[i].RematchingGap < 0){
        //     dataSource[i].RematchingGap = 0; //gap
        // }
      } else if(v.ActualPercent < 0 && v.Gap < 0){
        dataSource[i].RematchingGap = v.Gap - v.ActualPercent; //gap
        if(dataSource[i].RematchingGap >= 0){
            dataSource[i].RematchingGap = 0; //gap
        }
      } else if(v.ActualPercent >= 0 && v.Gap < 0){
        dataSource[i].RematchingGap = v.Gap;
      } else if(v.ActualPercent < 0 && v.Gap >= 0){
        // dataSource[i].RematchingGap = v.Gap;
      }
      
    })
    // console.log(dataSource)

	var selectedBM = Enumerable.From(ScorecardFullyearProjection.BusinessMetricsList()).Where("$.id === '"+ScorecardAnalysis.DataReference().Id+"'").FirstOrDefault();
    if(selectedBM===undefined){
      selectedBM = {};
      selectedBM.DecimalFormat = "0";
      selectedBM.naactual = false;
      selectedBM.natarget = false;
      selectedBM.valuetype = false;
    }else{
      selectedBM = selectedBM;
    }
    dataSource = Enumerable.From(dataSource).OrderBy("$.Gap").ToArray()

    dataSource = _.sortBy(dataSource, function(x){return x.Projection * -1})
    dataSourceForChart = Enumerable.From(dataSource).Where("$.NAActual == false && $.NATarget == false").ToArray()
    var MaxLength = 1;
    var MinLength = -1;
    if(dataSourceForChart.length > 0){
        // console.log(dataSourceForChart[0].Projection)
        MaxLength = dataSourceForChart[0].Projection + 0.1;
        MinLength = dataSourceForChart[dataSourceForChart.length-1].Projection - 0.1;
    }
    if(!isNegative){
        MinLength = 0;
    }
    if(MaxLength >= 2){
        MaxLength = 2;
    }
    if(MinLength <= -2){
        MinLength = -2
    }

    var DVWidth = $("#data-visualisation3-width").width();
    var DataVisualisation = $("#data-visualisation3");
    DataVisualisation.html("");
    if(dataSource.length == 0 && (selectedBM.natarget||selectedBM.naactual)){
        $("#data-visualisation3").html('<div class="emptychart-wrapper"><div>Not Applicable</div></div>');
        return false;
    }
    var allCountry = 27;
    if (dataSourceForChart.length<allCountry && dataSourceForChart.length!==0) {
        for(var i=dataSourceForChart.length;i<=27;i++){
            dataSourceForChart.push({Name:""});
        }
    }
    // console.log(dataSourceForChart)
    DataVisualisation.kendoChart({
        dataSource: {
            data: dataSourceForChart
        },
        legend: {
            visible: false
        },
        chartArea:{
            width:DVWidth,
            height:700
        },
        seriesDefaults: {
            gap:0.1,
          	// height:23,
        },
        series: [{
            field: "ActualPercent",
            // colorField:"ColorActualPercent",
            color:"#8c97a0",
            stack: true,
            // type: "column",
            type: "bar",
            labels: {
                visible: true,
                font:"bold 11px sans-serif",
                color:"#333",
                // template:function(e){
                //   tmp = 0;
                //   html = "";
                  
                //   if(e.dataItem.Gap >= 0){
                //     tmp = e.dataItem.Gap * 100;
                //     html = e.category+" "+kendo.toString(tmp,'N0')+" %";

                //   } else{
                //     tmp = e.dataItem.Gap * -100;
                //     html = kendo.toString(tmp,'N0')+" % "+e.category;
                //   }

                //   return html;
                // },
                template:" #: category #",
                // rotation:breakdown=="country"?270:0,
                rotation:0,
                // template:"#:kendo.toString(value*100,'N0')#%",
                background: "transparent",
                position:"insideBase",
                padding:{
                	top:10,
                	bottom:10
                }
                // margin:{
                //  bottom:0
                // }
            },
            tooltip:{
              visible:true,
              template:"#:kendo.toString(value,'P2')#"
            },
            overlay: {
              gradient: "none"
            },
            border: {
              width: 0
            },
        },{
            field: "RematchingBudgetProjection",
            colorField:"ColorBudgetProjection",
            stack: true,
            // type: "column",
            type: "bar",
            labels: {
                visible: false,
                font:"bold 11px sans-serif",
                color:"#333",
                // template:function(e){
                //   tmp = 0;
                //   html = "";
                  
                //   if(e.dataItem.Gap >= 0){
                //     tmp = e.dataItem.Gap * 100;
                //     html = e.category+" "+kendo.toString(tmp,'N0')+" %";

                //   } else{
                //     tmp = e.dataItem.Gap * -100;
                //     html = kendo.toString(tmp,'N0')+" % "+e.category;
                //   }

                //   return html;
                // },
                template:" #: category #",
                // rotation:breakdown=="country"?270:0,
                rotation:0,
                // template:"#:kendo.toString(value*100,'N0')#%",
                background: "transparent",
                position:"insideBase",
                padding:{
                	top:10,
                	bottom:10
                }
                // margin:{
                //  bottom:0
                // }
            },
            tooltip:{
              visible:true,
              template:"#:kendo.toString(dataItem.BudgetProjection,'P2')#"
            },
            overlay: {
              gradient: "none"
            },
            border: {
              width: 0
            },
        },{
            field: "RematchingProjectionforBudgetUse",
            color:"#D4D4D4",
            stack: true,
            labels: {
                visible: false,
                font:"8px Helvetica Neue",
                color:"#333",
                // template:function(e){
                //   tmp = 0;
                //   html = "";
                  
                //   if(e.dataItem.Gap >= 0){
                //     tmp = e.dataItem.Gap * 100;
                //     html = e.category+" "+kendo.toString(tmp,'N0')+" %";

                //   } else{
                //     tmp = e.dataItem.Gap * -100;
                //     html = kendo.toString(tmp,'N0')+" % "+e.category;
                //   }

                //   return html;
                // },
                template:function(e){
                    // console.log(e.dataItem.Gap)
                    return kendo.toString(e.dataItem.Gap,'P2');
                },
                padding:{
                    top: -20,
                    bottom: -20,
                },
                // rotation:breakdown=="country"?270:0,
                rotation:0,
                // template:"#:kendo.toString(value*100,'N0')#%",
                background: "transparent",
                position:"insideBase",
                // margin:{
                //  bottom:0
                // }
            },
            border: {
              width: 0
            },
            overlay: {
              gradient: "none"
            },
            tooltip:{
              visible:true,
              template:"#:kendo.toString(dataItem.Projection,'P2')#"
            }
        }],
        valueAxis: {
            max: MaxLength,
            min: MinLength,
            majorGridLines: {
                visible: false
            },
            majorUnit: 0.1,
            visible: true,
            labels:{
                template:"#:kendo.toString(value,'P0')#"
            },
            plotBands: [
                {from: 0.999,
                 to: 1.001,
                 color: "#818FAF"}
            ]
        },
        categoryAxis: {
            field: "Name",
            majorGridLines: {
                visible: false
            },
            labels:{
                visible:false,
                // rotation:breakdown=="country"?270:0,
                rotation:0,
                position:"insideEnd",
                // margin:{
                //  bottom:-250
                // },
                
            },
            line: {
                visible: false
            }
        }
    });

    // console.log(dataSource)
    if(dataSource.length === 0 ){
        var Message = 'No data to display.';
        $("#data-visualisation3").html('<div class="emptychart-wrapper"><div>'+Message+'</div></div>');
        // $('#textleftchart').hide()
    } else{
        // $('#textleftchart').show()
    }

}