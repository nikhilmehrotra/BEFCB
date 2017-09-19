/*
* @Author: Ainur
* @Date:   2016-10-30 12:30:30
* @Last Modified by:   Ainur
* @Last Modified time: 2016-11-23 09:26:34
*/

var chartColorsonDetails = ["#6D6E71", "#2890C0","#6AC17B","#6D6E71","#939598","#2890C0","#6BA8D0","#6AC17B","#9FD18B"];

var DetailBusinessDriver = {
	Processing:ko.observable(true),
    Id:ko.observable(""),
	Name:ko.observable(""),
	Data:ko.observable(), 
    Idx: ko.observable("no"),
    SelectedBI:ko.observableArray([]),
}
DetailBusinessDriver.GetInitiative = function(obj){
    if(typeof obj._id === "function"){
        Initiative.Get(obj._id());
    }else{
        Initiative.Get(obj._id);
    }
}

DetailBusinessDriver.ShowHideParent = function(id){
    var sources = ko.mapping.toJS(c.DataSource().Data);
    sources.TableSources = Enumerable.From(sources.TableSources).Where("$.Parentid === '"+id+"' ").ToArray()
    // console.log("0----> ",sources.TableSources)
    var scID = id;
    _.each(sources.TableSources, function(v,i){
        id = -1;
        _.each(c.DataSource().Data.TableSources, function(vv,ii){
            if(vv.Id() == v.Id){
                id = ii;
            }
        })

        if(v.ShowHide){
            c.DataSource().Data.TableSources[id].ShowHide(false);
            c.DataSource().Data.SummaryBusinessDriver[id].ShowHide(false);
            $('#caret'+scID).removeClass('fa-caret-down');
            $('#caret'+scID).addClass('fa-caret-right');
        } else{
            c.DataSource().Data.TableSources[id].ShowHide(true);
            c.DataSource().Data.SummaryBusinessDriver[id].ShowHide(true);
            $('#caret'+scID).removeClass('fa-caret-right');
            $('#caret'+scID).addClass('fa-caret-down');
        }
    })
    c.SyncSCHeight()
}

DetailBusinessDriver.ShowHideSecLvl = function(id){
    var sources = ko.mapping.toJS(c.DataSource().Data);
    sources.TableSources = Enumerable.From(sources.TableSources).Where("$.Id === '"+id+"' ").ToArray()
    console.log("0----> ",sources.TableSources)
    var scID = id;
    _.each(sources.TableSources, function(v,i){
        id = -1;
        _.each(c.DataSource().Data.TableSources, function(vv,ii){
            if(vv.Id() == v.Id){
                id = ii;
            }
        })

        if(v.ShowHideSecLvl){
            c.DataSource().Data.TableSources[id].ShowHideSecLvl(false);
            c.DataSource().Data.SummaryBusinessDriver[id].ShowHideSecLvl(false);
        } else{
            c.DataSource().Data.TableSources[id].ShowHideSecLvl(true);
            c.DataSource().Data.SummaryBusinessDriver[id].ShowHideSecLvl(true);
        }
    })
}

DetailBusinessDriver.GetDataSource = function(idx, mode){

  var activeTab = $("#dashboard .nav-pills li.active a").attr("href").replace("#","");
  var sources = ko.mapping.toJS(c.DataSource().Data);
  sources['DetailBusinessDriverList'] = [];
  sources['TabMenuValue'] = TabMenuValue();
  if(mode == 'seclvl' ){
      sources.TableSources = Enumerable.From(sources.TableSources).Where("$.Idx === '"+idx+"' ").ToArray()
      sources.TableSources[0].ShowParent = true;
      sources.TableSources[0].TotalParent = 1;

      DetailBusinessDriver.Id(sources.TableSources[0].Id);
      DetailBusinessDriver.Name(sources.TableSources[0].Name);
  } else if(mode == 'firstlvl'){
      sources.TableSources = Enumerable.From(sources.TableSources).Where("$.Parentid === '"+idx+"' ").ToArray()

      DetailBusinessDriver.Id(sources.TableSources[0].Parentid);
      DetailBusinessDriver.Name(sources.TableSources[0].Parentname);
  }

  var tatalCompletionList = [] 
  _.each(sources.TableSources, function(vv,ii){
    sources.TableSources[ii].OverallProgress = [];
    var tatalCompletion = 0;

    for(var lc in sources.TableSources[ii].LifeCycle){
      for(var i in sources.TableSources[ii].LifeCycle[lc].Initiatives){
        // var Completion = Math.random()*100;
              // console.log(sources.TableSources[ii].LifeCycle[lc].Initiatives[i]);
        // sources.TableSources[ii].LifeCycle[lc].Initiatives[i].Completion = Completion > 100 ? 100 : Completion;
              if(!sources.TableSources[ii].LifeCycle[lc].Initiatives[i].IsTask){
                  tatalCompletion += sources.TableSources[ii].LifeCycle[lc].Initiatives[i].ProgressCompletion;
                  sources.TableSources[ii].OverallProgress.push(sources.TableSources[ii].LifeCycle[lc].Initiatives[i]);
              }
      }
    }

    tatalCompletionList.push(tatalCompletion)
  })

  if(mode == 'seclvl' ){
    sources.BusinessDriverData = Enumerable.From(sources.BusinessDriverData).Where("$.Id === '"+sources.TableSources[0].Id+"' ").ToArray()
    sources.SummaryBusinessDriver = Enumerable.From(sources.SummaryBusinessDriver).Where("$.Id === '"+sources.TableSources[0].Id+"' ").ToArray()
  } else if(mode == 'firstlvl'){
    sources.BusinessDriverData = Enumerable.From(sources.BusinessDriverData).Where("$.Parentid === '"+sources.TableSources[0].Parentid+"' ").ToArray()
    sources.SummaryBusinessDriver = Enumerable.From(sources.SummaryBusinessDriver).Where("$.Parentid === '"+sources.TableSources[0].Parentid+"' ").ToArray()
  }

  _.each(sources.TableSources, function(vv,ii){
    if(sources.TableSources[ii].OverallProgress.length==0){
      sources.SummaryBusinessDriver[ii].Completion = 0;
    }else{
      sources.SummaryBusinessDriver[ii].Completion = tatalCompletionList[ii]/sources.TableSources[ii].OverallProgress.length;
    }

    var tmpData = {}
    tmpData['Name'] = sources.TableSources[ii].Name;
    tmpData['Mode'] = mode;

    xxx = false;
    if (mode =='seclvl'){
      xxx = true;
    }

    tmpData['btnshowhide'] = xxx;
    tmpData['TableSources'] = [];
    tmpData['TableSources'].push(sources.TableSources[ii])
    tmpData['SummaryBusinessDriver'] = [];
    tmpData['SummaryBusinessDriver'].push(sources.SummaryBusinessDriver[ii])

    sources.DetailBusinessDriverList.push(tmpData)

  })

  DetailBusinessDriver.Data(ko.mapping.fromJS(sources));
}
DetailBusinessDriver.Get = function(idx){
    DetailBusinessDriver.SelectedBI([]);
    DetailBusinessDriver.Idx(idx)
    DetailBusinessDriver.Processing(true);
  	DetailBusinessDriver.GetDataSource(idx, 'seclvl')
  	$("#detail-businessdriver").modal("show");
  	DetailBusinessDriver.Render();
    c.SyncSCHeight();
}

DetailBusinessDriver.GetParent = function(idParent){
    DetailBusinessDriver.SelectedBI([]);
    // DetailBusinessDriver.Idx(idx)
    DetailBusinessDriver.Processing(true);
    DetailBusinessDriver.GetDataSource(idParent, 'firstlvl')
    $("#detail-businessdriver").modal("show");
    DetailBusinessDriver.Render();
    c.SyncSCHeight();
}

DetailBusinessDriver.Render = function(){
  var psBapake = $("#progression-summary");
  psBapake.html("");

  _.each(DetailBusinessDriver.Data().TableSources(), function(vvv,iii){

  
    var ds = ko.mapping.toJS(DetailBusinessDriver.Data().TableSources()[iii].OverallProgress());

    if(ds.length>0){
var BusinessMetrics = ko.mapping.toJS(DetailBusinessDriver.Data().SummaryBusinessDriver()[iii].BusinessMetrics());
    // BusinessMetrics = [BusinessMetrics[iii]]
      var dataSource = [];
      // Get Chart Data
    var BMMinDate = new Date(Now);
    for(var bm in BusinessMetrics){
        BMMinDate = new Date(BusinessMetrics[bm].BaseLinePeriod)
        var mPeriod = Enumerable.From(BusinessMetrics[bm].ActualData).Min("jsonDate($.Period)");
        // console.log("<<<-----",BusinessMetrics[bm], mPeriod, BMMinDate)
        if(mPeriod<BMMinDate){
            BMMinDate = mPeriod;
        }
    }
      
    var BusinessImpactList = Initiative.BusinessImpactList()
    var StartDate = Enumerable.From(ds).Min("jsonDate($.FinishDate)");
    if(BMMinDate<StartDate){
        StartDate = BMMinDate;
    }
    
    var EndDate = Enumerable.From(ds).Max("jsonDate($.FinishDate)");
    var currentDMValue = {};
    for(var i = new Date(StartDate.setDate(1));i<=new Date(EndDate.setDate(1));i.setMonth(i.getMonth()+1)){
        var d = {
              Period:kendo.toString(new Date(i),"MMM yy"),
          };
          var total = 0;
          
      for(var bi in BusinessImpactList){
        var biData = Enumerable.From(ds).Where("$.BusinessImpact === '"+BusinessImpactList[bi].value+"' && kendo.toString(jsonDate($.FinishDate),'MMM yy') === '"+d.Period+"'").ToArray()
              d[BusinessImpactList[bi].value] = biData.length;
          }
      
      for(var bm in BusinessMetrics){
        // console.log(BusinessMetrics);
        var bmData = Enumerable.From(BusinessMetrics[bm].ActualData).Where("kendo.toString(jsonDate($.Period),'MMM yy') === '"+d.Period+"'").FirstOrDefault();
        var date = new Date(new Date(i).setDate(1));
        var target = jsonDate(BusinessMetrics[bm].TargetPeriod);
        var CurrentDate = Enumerable.From(BusinessMetrics[bm].ActualData).Where("$.Value!==0").Max("jsonDate($.Period)");
        // console.log(CurrentDate);
        if(date<=CurrentDate){
          
            if(kendo.toString(i,'MMM yy') == kendo.toString(new Date(BusinessMetrics[bm].BaseLinePeriod),'MMM yy')){
              d["bm"+bm] = BusinessMetrics[bm].BaseLineValue;
            } else{
              d["bm"+bm] = bmData !== undefined ? bmData.Value : undefined;
            }
            
            // console.log(d);
            // console.log("---");
            if(d.Period === kendo.toString(CurrentDate,'MMM yy')){
                // console.log(d["bm"+bm] );
                currentDMValue[bm] = bmData !== undefined ? bmData.Value : 0;
                currentDMValue[bm+"Period"] = target == "" ? 0 : monthDiff(date,target)+1;
                // console.log("----->>>>",currentDMValue[bm+"Period"])
                currentDMValue[bm+"Idx"] = 1;
                d["ebm"+bm] = currentDMValue[bm];
            }
        }else if(date<=target){
            // console.log(BusinessMetrics[bm].TargetValue);
            // console.log(currentDMValue[bm]);
            var tvalue = currentDMValue[bm+"Period"] == 0 ? 0 : (BusinessMetrics[bm].TargetValue-currentDMValue[bm])/currentDMValue[bm+"Period"];
            // console.log(tvalue);
            // console.log(BusinessMetrics[bm].DataPoint+"|"+tvalue+"|"+currentDMValue[bm+"Idx"]+"|"+currentDMValue[bm+"Period"]+"#"+BusinessMetrics[bm].TargetValue);
            d["ebm"+bm] = currentDMValue[bm]+(tvalue * (parseInt(currentDMValue[bm+"Idx"])));
            if(d.Period === kendo.toString(target,"MMM yy")){
                d["eLabel"+bm] = d["ebm"+bm];
                // console.log(bm)
             }


            // console.log(d["ebm"+bm]);
            currentDMValue[bm+"Idx"] += 1;
        }
        // console.log("__________________________________--");
      }
          
      // d.average = BusinessImpactList.length === 0 ? 0 : total/BusinessImpactList.length;
          dataSource.push(d);
      }

    // dataSource.ps.push({Period:kendo.toString(new Date(2016,0,1),'MMM yy'),Low:2,Medium:0,High:1,liabilities:50});
    // dataSource.ps.push({Period:kendo.toString(new Date(2016,1,1),'MMM yy'),Low:0,Medium:1,High:0,liabilities:30});
    // dataSource.ps.push({Period:kendo.toString(new Date(2016,2,1),'MMM yy'),Low:0,Medium:3,High:0,liabilities:10});
    // dataSource.ps.push({Period:kendo.toString(new Date(2016,3,1),'MMM yy'),Low:0,Medium:0,High:0,liabilities:60});
    // dataSource.ps.push({Period:kendo.toString(new Date(2016,4,1),'MMM yy'),Low:1,Medium:1,High:1,liabilities:40});
    // dataSource.ps.push({Period:kendo.toString(new Date(2016,5,1),'MMM yy'),Low:0,Medium:0,High:0,liabilities:50});
    // dataSource.ps.push({Period:kendo.toString(new Date(2016,6,1),'MMM yy'),Low:0,Medium:0,High:0,liabilities:50});
    // dataSource.ps.push({Period:kendo.toString(new Date(2016,7,1),'MMM yy'),Low:0,Medium:0,High:2,liabilities:60});
    // dataSource.ps.push({Period:kendo.toString(new Date(2016,8,1),'MMM yy'),Low:1,Medium:0,High:0,liabilities:30});


    // var liabilities = 0;
    // for(var i in dataSource.ps){
    //     dataSource.ps[i].liabilities = liabilities + dataSource.ps[i].liabilities;
    //     liabilities += dataSource.ps[i].liabilities;
                
    // }

    // Get Series
    var PSSeriesList = [];
    for(var bi in BusinessImpactList){
        PSSeriesList.push({field: BusinessImpactList[bi].value,name:BusinessImpactList[bi].text,axis:"value"});
    }

    var valueAxisList = [
        {
            name:"value",
            majorGridLines: {
                visible: true
            },
            line:{
                visible:false
            },
            labels:{
                font:chartFont,
                template:"#:value<=1000?kendo.toString(value,'N0'):kendo.toString(value/1000,'N0')+'k'#"
            },
            visible: true,
            title: {
                text: "Number of Initiative",
                // color: "#333",
                font: "14px Helvetica Neue"
            }
        },
    ];

    axisCrossingValuesList = [0]
    for(var bm in BusinessMetrics){
        PSSeriesList.push({field:"bm"+bm,name:"Actual "+BusinessMetrics[bm].DataPoint,axis:BusinessMetrics[bm].DataPoint,type:"line",stack:false,labels:{visible:false},
            markers: {
              size:5
            },
        });
        now = PSSeriesList.length;
        // console.log("-----< ", jsonDate(BusinessMetrics[bm].TargetPeriod))
        PSSeriesList.push({field:"ebm"+bm,name:"Target "+BusinessMetrics[bm].DataPoint,axis:BusinessMetrics[bm].DataPoint,type:"line",stack:false,labels:{visible:false},
            markers: {
              visible: false,
            },
            legend:{
                visible:false,
            },
            tooltip:{
                visible:false
            },
            labels:{
                visible:true,
                color:"#333",
                template:"#:dataItem.eLabel===undefined?'':kendo.toString(dataItem.eLabel,'N0')#",
                position:"top"
            }
        });

        // console.log(bm, PSSeriesList[now])
        if(bm == 0){
            PSSeriesList[now].labels.template = "#:dataItem.eLabel0===undefined?'':kendo.toString(dataItem.eLabel0,'N0')#"
        } else if(bm == 1){
            PSSeriesList[now].labels.template = "#:dataItem.eLabel1===undefined?'':kendo.toString(dataItem.eLabel1,'N0')#"
        } else if(bm == 2){
            PSSeriesList[now].labels.template = "#:dataItem.eLabel2===undefined?'':kendo.toString(dataItem.eLabel2,'N0')#"
        }
        valueAxisList.push(
            {
                name:BusinessMetrics[bm].DataPoint,
                majorGridLines: {
                    visible: true
                },
                line:{
                    visible:false
                },
                labels:{
                    font:chartFont,
                    template:"#:kendo.toString(value,'N0')#"
                },
                visible: true,
                title: {
                    text: BusinessMetrics[bm].DataPoint,
                     // color: "#333",
                    font: "14px Helvetica Neue"
                }, 
                color: chartColorsonDetails[bm]
            }
        );
        axisCrossingValuesList.push(dataSource.length+bm);
    }

    // console.log(valueAxisList);

      DetailBusinessDriver.Processing(false);

    $pstmp = jQuery("<div />");
    $pstmp.attr({id:"progression-summary"+iii});
    $pstmp.appendTo("#progression-summary");

    ps = $("#progression-summary"+iii)
    // console.log(dataSource)
      
    ps.kendoChart({
        dataSource: {
            data:dataSource
        },
        title: {
            visible:true,            
            text: "Progress vs Business Impact"
           
        },
        chartArea:{
            height:300
        },
        legend: {
            visible: true,
            position:"bottom",
            labels:{
                font:chartFont
            }
        },
        seriesDefaults: {
            type: "column",
            stack:true,
            overlay: {
              gradient: "none"
            },
            labels:{
                // visible: chartLabelDisplay,
                visible:false,
                font:chartFont,
                // color:chartLabelColor,
                color:"#FFF",
                background:"transparent",
                position:"center",
                template: "#= series.field=='conversion'?kendo.toString(value, 'P2'):kendo.toString(value, 'N0') #",
                padding:0,
                margin:0,
            }
        },
        series: PSSeriesList,
        valueAxis : valueAxisList,
        seriesClick: function(e){
            if(e.series.axis==="value"){
                // console.log(e);
                var InitiativeList = ko.mapping.toJS(DetailBusinessDriver.Data().Project());
                var arr = Enumerable.From(InitiativeList).Where("$.BusinessDriverId === '"+DetailBusinessDriver.Id()+"' && $.BusinessImpact === '"+e.series.field+"' && kendo.toString(jsonDate($.FinishDate),'MMM yy') === '"+e.category+"'").ToArray();

                DetailBusinessDriver.SelectedBI(arr);
            }
        },
        // valueAxis: [
        //     {
        //         name:"value",
        //         majorGridLines: {
        //             visible: true
        //         },
        //         line:{
        //             visible:false
        //         },
        //         labels:{
        //             font:chartFont,
        //             template:"#:value<=1000?kendo.toString(value,'N0'):kendo.toString(value/1000,'N0')+'k'#"
        //         },
        //         visible: true,
        //         title: {
        //             text: "Number of Initiative",
        //             // color: "#333",
        //             font: "14px Helvetica Neue"
        //         }
        //     },
        //     {
        //         name:"percentage",
        //         majorGridLines: {
        //             visible: true
        //         },
        //         line:{
        //             visible:false
        //         },
        //         labels:{
        //             font:chartFont,
        //             template:"#:kendo.toString(value,'N0')#"
        //         },
        //         visible: true,
        //         title: {
        //             text: "Data Metrics",
        //              // color: "#333",
        //             font: "14px Helvetica Neue"
        //         }
        //     }
        // ],
        categoryAxis: {
            field: "Period",
            majorGridLines: {
                visible: false
            },
            labels:{
                font:chartFont,
                rotation:40
            },
            axisCrossingValues: axisCrossingValuesList
        },
        tooltip: {
            visible: true,
            font:chartFont,
            template: "#= series.name #: #= series.field=='conversion'?kendo.toString(value, 'P2'):kendo.toString(value, 'N0') #"
        },
        // seriesColors:chartColors,
        seriesColors:["#6D6E71", "#2890C0","#6AC17B","#6D6E71","#939598","#2890C0","#6BA8D0","#6AC17B","#9FD18B"],
    });

    // setTimeout(function(){
    //     // $("#progression-summary"+iii).data("kendoChart").refresh();
    // }, 100);

    }

  })
}