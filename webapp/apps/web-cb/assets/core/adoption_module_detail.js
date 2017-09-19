var AdoptionModuleDetail = {
	Processing:ko.observable(false),
	InitiativeData:ko.observable(""),//## GLOBAL/REGION/COUNTRY
	InitiativeId:ko.observable(""),
	InputType:ko.observable('Actual'), //## VALUE / RAG true=>value
	InitiativeList:ko.observableArray([]),
	Data:ko.observableArray([]),
	DataBack:ko.observableArray([]),
	Mode:ko.observable("edit"),
	// YearValue:ko.observable(new Date()),
	OptionalList: ko.observableArray([
	    {Id: 'ActualYTD', Text: 'Actual YTD'},
	    {Id: 'Budget', Text: 'Budget'},
	    {Id: 'RAG', Text: 'RAG'},
	]),
	Year:ko.observable(Now.getFullYear()),
	// Data Source
	ActiveMetric:ko.observable(""),
	YearList:ko.observableArray([]),
	MetricData:ko.observableArray([]),
	MetricDataBack:ko.observableArray([]),
	MetricDataPool:ko.observable({}),
	MetricForm:{
		MetricId:ko.observable(""),
		MetricName:ko.observable(""),
		Denomination:ko.observable(""),
		Description:ko.observable(""),
		ActualValue:ko.observable(true),
		TotalValue:ko.observable(false),
		RAGValue:ko.observable(false),
	},
	MetricsTypeList : ko.observableArray([
        {name: "Dollar Value ($)", value: "DOLLAR"},
        {name: "Numeric Value", value: "NUMERIC"},
        {name: "Percentage (%)", value: "PERCENTAGE"},
    ]),
    RAGValue: ko.observable(false),
    TotalValue:ko.observable(false),
    ActualValue: ko.observable(false),
    // AvailableData: ko.observableArray([])
}
// AdoptionModuleDetail.ActiveMetric.subscribe(function(d){
// 	if(AdoptionModuleDetail.Mode() != 'edit'){
// 		AdoptionModuleDetail.GetData();
// 	}
// })
AdoptionModuleDetail.Get = function(d){
	// console.log("--",d)
	InitiativeData = d.InitiativeData
	InitiativeList  = d.InitiativeList
	// MetricData = d.MetricData;

	AdoptionModule.Mode("DETAIL");
	AdoptionModuleDetail.Mode("edit");
	// AdoptionModuleDetail.YearValue(kendo.toString(new Date(), 'yyyy'));
	// AdoptionModuleDetail.MetricData(ko.mapping.fromJS(MetricData)());
	// AdoptionModuleDetail.MetricDataBack(ko.mapping.fromJS(MetricData)());
	AdoptionModuleDetail.Year(Now.getFullYear());
	AdoptionModuleDetail.InputType('Actual');
	AdoptionModuleDetail.InitiativeId(InitiativeData.InitiativeID);
	AdoptionModuleDetail.InitiativeData(InitiativeData);
	AdoptionModuleDetail.InitiativeList(InitiativeList);
}

AdoptionModuleDetail.Save = function(){
	// console.log("save");
	AdoptionModuleDetail.Mode("view");
	
	var MetricData = ko.mapping.toJS(AdoptionModuleDetail.MetricData);
	// var MetricDataPoolTmp = ko.mapping.toJS(AdoptionModuleDetail.MetricDataPool());
	// var MetricDataPool = ko.mapping.toJS(AdoptionModuleDetail.MetricDataPool());
	var ActiveMetric = AdoptionModuleDetail.ActiveMetric();
	for(var x in MetricData){
		if(ActiveMetric == MetricData[x].MetricId){
			var data_pool = MetricData[x];

			data_pool.DetailData = ko.mapping.toJS(AdoptionModuleDetail.Data());
			data_pool.DetailDataBack = ko.mapping.toJS(AdoptionModuleDetail.DataBack());
			MetricData[x] = data_pool;
			// console.log(data_pool.DetailData, data_pool,ActiveMetric, MetricData[x].MetricId)
			// var pool  = ko.mapping.toJS(AdoptionModuleDetail.MetricDataPool);
			// MetricDataPool[MetricData[x].MetricId] = data_pool;
		}
	}
	var Report = [];
	_.each(MetricData, function(v0,i0){
		var parm = v0.DetailData;
		var oldparm = v0.DetailDataBack;
		_.each(parm, function(v,i){
			// parm.Year = parseInt(parm.Year)
			console.log(MetricData[i0], parm)
			mont = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
			
			_.each(mont, function(vv,ii){
				if(parm[i][vv+""] === ""){
					parm[i]["NA"+vv]=true;
					parm[i][vv+""]=0;
				}else if(parm[i][vv+""] !== 0){
					parm[i][vv+""] = parseFloat(parm[i][vv+""]);
					parm[i]["NA"+vv]=false;
				}

				if(oldparm[i][vv+""] === ""){
					oldparm[i]["NA"+vv]=true;
					oldparm[i][vv+""]=0;
				}else if(oldparm[i][vv+""] !== 0){
					oldparm[i][vv+""] = parseFloat(oldparm[i][vv+""]);
					oldparm[i]["NA"+vv]=false;
				}

				if(parm[i]["Total"+vv] === "" || !MetricData[i0].TotalValue){
					parm[i]["NATotal"+vv]=true;
					parm[i]["Total"+vv]=0;
				}else if(parm[i]["Total"+vv] !== 0){
					parm[i]["Total"+vv] = parseFloat(parm[i]["Total"+vv]);
					parm[i]["NATotal"+vv]=false;
				}

				if(oldparm[i]["Total"+vv] === ""){
					oldparm[i]["NATotal"+vv]=true;
					oldparm[i]["Total"+vv]=0;
				}else if(oldparm[i]["Total"+vv] !== 0){
					oldparm[i]["Total"+vv] = parseFloat(oldparm[i]["Total"+vv]);
					oldparm[i]["NATotal"+vv]=false;
				}

				if(!MetricData[i0].RAGValue){
					parm[i]["RAG"+vv]="";
				}
			})	

			_.each(v, function(vv,ii){
		        if(parm[i][ii] !== oldparm[i][ii]){

					var NewChanges = {
		                Whatchanged:parm[i].InitiativeId+" - "+parm[i].Country+" - "+ii,
		                OldValue:oldparm[i][ii],
		                NewValue:parm[i][ii],
		            }
		            Report.push(NewChanges)
		        }
		    })
		})
		MetricData[i0].DetailData = parm;
	})

	trueParam = {
		// MilestoneData:parm,
		InitiativeId:AdoptionModuleDetail.InitiativeId(),
		Report: Report,
		MetricData:MetricData,
		ActiveMetric:ActiveMetric,
	}
	// console.log(trueParam)
	
	ajaxPost("/web-cb/adoptionmodule/savedetail",trueParam,function(res){
		if(res.IsError){
			swal("Error",res.Message,"error")
			return false;
		}

		// console.log(res)
		AdoptionModuleDetail.GetData()
		// var parm = {
		// 	InitiativeId: AdoptionModuleDetail.InitiativeId(),
		// 	year:parseInt(AdoptionModuleDetail.Year()),
		// 	// MetricId:'',
		// }
		// // AdoptionModuleDetail.MetricDataPool([])
		// if(parm.InitiativeId != ""){
		// 	ajaxPost("/web-cb/adoptionmodule/getdatadetail",parm,function(res){
		// 		var tmp = []
		// 		var ActiveMetric = AdoptionModuleDetail.ActiveMetric();
		// 		var MetricData = ko.mapping.toJS(AdoptionModuleDetail.MetricData);
		// 		var Index = 0;
		// 		for(var x in MetricData){
		// 			if(MetricData[x].MetricId == ActiveMetric){
		// 				Index = x;
		// 				break;
		// 			}
		// 		}
		// 		var data = Enumerable.From(MetricData).Where("$.MetricId === '"+ActiveMetric+"'").FirstOrDefault();
		// 		var NewMetricData = res.Data.MetricData;
		// 		var ex_data = Enumerable.From(NewMetricData).Where("$.MetricId === '"+ActiveMetric+"'").FirstOrDefault();
		// 		AdoptionModuleDetail.MetricData(ko.mapping.fromJS(res.Data.MetricData)())
		// 		if(typeof ex_data == "undefined") {
		// 			for(var x in NewMetricData){
		// 				if(x == Index &&NewMetricData[x].MetricName == data.MetricName){
		// 					ex_data = NewMetricData[x]
		// 					break;
		// 				}
		// 			}
		// 			if(typeof ex_data == "undefined"){
		// 				ex_data = Enumerable.From(NewMetricData).Where("$.MetricName === '"+data.MetricName+"'").FirstOrDefault();
		// 			}
		// 		}
		// 		AdoptionModuleDetail.ActiveMetric(ex_data.MetricId);
		// 		if(ActiveMetric === ex_data.MetricId){
		// 			AdoptionModuleDetail.GetData();
		// 		}
		// 	})
		// }

	})
}
AdoptionModuleDetail.ChangeMetric = function(d){
	var data = ko.mapping.toJS(d);
	
	var ActiveMetric = AdoptionModuleDetail.ActiveMetric();
	var MetricData = ko.mapping.toJS(AdoptionModuleDetail.MetricData);
	var Index = 0;
	for(var x in MetricData){
		if(MetricData[x].MetricId == ActiveMetric){
			Index = x;
			break;
		}
	}

	if(ActiveMetric!==""){
		var DataTable = ko.mapping.toJS(AdoptionModuleDetail.Data());
		AdoptionModuleDetail.MetricData()[Index].DetailData(ko.mapping.fromJS(DataTable)())
		var data_pool = Enumerable.From(MetricData).Where("$.MetricId == '"+ActiveMetric+"'").FirstOrDefault();
		if(data.MetricId != ""){
			AdoptionModuleDetail.Data(ko.mapping.fromJS(data.DetailData)())
		}
	}
	AdoptionModuleDetail.ActiveMetric(data.MetricId);
	AdoptionModuleDetail.RAGValue(data.RAGValue)
	AdoptionModuleDetail.TotalValue(data.TotalValue)
	AdoptionModuleDetail.ActualValue(data.ActualValue)
	$(".ScorecardTabMenu").removeClass('active')
	$($(".ScorecardTabMenu")[0]).addClass('active')
	AdoptionModuleDetail.InputType('Actual')
}
AdoptionModuleDetail.EditMetric = function(d){
	var da = ko.mapping.toJS(d);
	var MetricForm = AdoptionModuleDetail.MetricForm;
	MetricForm.MetricId(da.MetricId);
	MetricForm.MetricName(da.MetricName);
	MetricForm.Denomination(da.Denomination);
	MetricForm.Description(da.Description);
	MetricForm.ActualValue(da.ActualValue);
	MetricForm.TotalValue(da.TotalValue);
	MetricForm.RAGValue(da.RAGValue);
	$("#RAGValue").bootstrapSwitch('state', da.RAGValue);
	$("#TotalValue").bootstrapSwitch('state', da.TotalValue);
	$("#adptionmetric-form").modal("show");
}
AdoptionModuleDetail.RemoveMetric = function(){
	swal({
		title: "",
		text: 'Are You Sure?',
		type: "warning",
		showCancelButton: true,
		confirmButtonClass: "btn-danger",
		confirmButtonText: "Yes, Remove Now",
		closeOnConfirm: true
	},
	function(){	
		var ActiveMetric = AdoptionModuleDetail.ActiveMetric();
		var MetricData = ko.mapping.toJS(AdoptionModuleDetail.MetricData);
		var data_pool = Enumerable.From(MetricData).Where("$.MetricId == '"+ActiveMetric+"'").FirstOrDefault();
		var newmetric = [];

		_.each(MetricData, function(v,i){
			if(v.MetricId != data_pool.MetricId){
				var newya = ko.mapping.fromJS(v)
				newmetric.push(newya)
			}
		})
		AdoptionModuleDetail.MetricData(newmetric)
		// console.log("fak", newmetric)
		lastindex = newmetric.length - 1;
		if(lastindex >= 0){
			d = ko.mapping.toJS(newmetric[lastindex])
			AdoptionModuleDetail.Data(ko.mapping.fromJS(newmetric[lastindex].DetailData)())
			AdoptionModuleDetail.ActiveMetric(d.MetricId)
			AdoptionModuleDetail.ActualValue(d.ActualValue);
			AdoptionModuleDetail.TotalValue(d.TotalValue);
			AdoptionModuleDetail.RAGValue(d.RAGValue);
		} else{
			AdoptionModuleDetail.Data([])
			AdoptionModuleDetail.ActiveMetric("")
			AdoptionModuleDetail.ActualValue(false);
			AdoptionModuleDetail.TotalValue(false);
			AdoptionModuleDetail.RAGValue(false);
		}
		
		
		//ajax remove metric
	});
}
AdoptionModuleDetail.Cancel = function(){
	var BAK = ko.mapping.toJS(AdoptionModuleDetail.DataBack())
	var BAKF = ko.mapping.fromJS(BAK)
	AdoptionModuleDetail.Data(BAKF())
	var BK = ko.mapping.toJS(AdoptionModuleDetail.MetricDataBack())
	var BKF = ko.mapping.fromJS(BK) 
	AdoptionModuleDetail.MetricData(BKF())
	AdoptionModuleDetail.Mode("view");

	// lastindex = BK.length - 1;
	if(BK.length > 0){
		d = ko.mapping.toJS(BK[0])
		AdoptionModuleDetail.Data(ko.mapping.fromJS(BK[0].DetailData)())
		AdoptionModuleDetail.ActiveMetric(d.MetricId)
		AdoptionModuleDetail.ActualValue(d.ActualValue);
		AdoptionModuleDetail.TotalValue(d.TotalValue);
		AdoptionModuleDetail.RAGValue(d.RAGValue);
		AdoptionModuleDetail.InputType("Actual")
	} else{
		AdoptionModuleDetail.Data([])
		AdoptionModuleDetail.ActiveMetric("")
		AdoptionModuleDetail.ActualValue(false);
		AdoptionModuleDetail.TotalValue(false);
		AdoptionModuleDetail.RAGValue(false);
		AdoptionModuleDetail.InputType("")
	}
}
AdoptionModuleDetail.Edit = function(){
	AdoptionModuleDetail.Mode("edit");
	var BK = ko.mapping.toJS(AdoptionModuleDetail.MetricDataBack())
	if(BK.length >= 0){
		d = ko.mapping.toJS(BK[0])
		AdoptionModuleDetail.Data(ko.mapping.fromJS(BK[0].DetailData)())
		AdoptionModuleDetail.ActiveMetric(d.MetricId)
		AdoptionModuleDetail.ActualValue(d.ActualValue);
		AdoptionModuleDetail.TotalValue(d.TotalValue);
		AdoptionModuleDetail.RAGValue(d.RAGValue);
		AdoptionModuleDetail.InputType("Actual")
	} else{
		AdoptionModuleDetail.Data([])
		AdoptionModuleDetail.ActiveMetric("")
		AdoptionModuleDetail.ActualValue(false);
		AdoptionModuleDetail.TotalValue(false);
		AdoptionModuleDetail.RAGValue(false);
		AdoptionModuleDetail.InputType("")
	}
}
AdoptionModuleDetail.Return = function(){
	AdoptionModule.Mode('ANALYSIS');
	AdoptionModuleAnalysis.GetData();
}

AdoptionModuleDetail.AddMetric = function(){
	var MetricForm = AdoptionModuleDetail.MetricForm;
	MetricForm.MetricId("");
	MetricForm.MetricName("");
	MetricForm.Denomination("");
	MetricForm.Description("");
	MetricForm.ActualValue(true);
	MetricForm.TotalValue(false);
	MetricForm.RAGValue(false);
	$("#RAGValue").bootstrapSwitch('state', false);
	$("#TotalValue").bootstrapSwitch('state', false);
	$("#adptionmetric-form").modal("show");
}
AdoptionModuleDetail.SaveMetric = function(){
	var validator = $("#adptionmetric-form").data("kendoValidator");
    if(validator==undefined){
       validator= $("#adptionmetric-form").kendoValidator().data("kendoValidator");
    }
    if (validator.validate()) {
		var MetricForm = ko.mapping.toJS(AdoptionModuleDetail.MetricForm);
		
		if(MetricForm.MetricId==""){
			MetricForm.MetricId = "Metric"+AdoptionModuleDetail.MetricData().length;
			AdoptionModuleDetail.ActiveMetric(MetricForm.MetricId)

			var parm = {
				InitiativeId: AdoptionModuleDetail.InitiativeId(),
				year:parseInt(AdoptionModuleDetail.Year()),
				// MetricId:AdoptionModuleDetail.ActiveMetric(),
			}
			if(parm.InitiativeId != ""){
				ajaxPost("/web-cb/adoptionmodule/getdatadetailold",parm,function(res){
					if(res.IsError){
						swal("Error",res.Message,"error")
						return false;
					}

					_.each(res.Data.DetailData, function(v,i){
						mont = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
					    _.each(mont, function(vv,ii){
							if(res.Data.DetailData[i]["NA"+vv+""]){
								res.Data.DetailData[i][vv+""] = "";
							}
							if(res.Data.DetailData[i]["NATotal"+vv]){
								res.Data.DetailData[i]["Total"+vv] = "";
							}
					    })
					})

					MetricForm.DetailData =res.Data.DetailData
					MetricForm.DetailDataBack = res.Data.DetailData
					// console.log(MetricForm)
					AdoptionModuleDetail.MetricData.push(ko.mapping.fromJS(MetricForm));
					AdoptionModuleDetail.Data(ko.mapping.fromJS(res.Data.DetailData))
					AdoptionModuleDetail.ChangeMetric(MetricForm);
				})
			}

		}else{
			var MetricData = ko.mapping.toJS(AdoptionModuleDetail.MetricData);
			var idx = 0;
			for(var i in MetricData){
				if(MetricData[i].MetricId == MetricForm.MetricId){
					idx = i;
					break;
				}
			}

			da = ko.mapping.toJS(AdoptionModuleDetail.MetricForm)
			AdoptionModuleDetail.RAGValue(da.RAGValue)
			AdoptionModuleDetail.TotalValue(da.TotalValue)
			$(".ScorecardTabMenu").removeClass('active')
			$($(".ScorecardTabMenu")[0]).addClass('active')
			AdoptionModuleDetail.InputType('Actual')

			AdoptionModuleDetail.MetricData()[idx].RAGValue(da.RAGValue)
			AdoptionModuleDetail.MetricData()[idx].TotalValue(da.TotalValue)
			AdoptionModuleDetail.MetricData()[idx].Denomination(da.Denomination)
			AdoptionModuleDetail.MetricData()[idx].MetricName(da.MetricName)
			AdoptionModuleDetail.MetricData()[idx].Description(da.Description)
			// console.log(AdoptionModuleDetail.MetricData()[idx]);
		}
		$("#adptionmetric-form").modal("hide");
	}
}

AdoptionModuleDetail.InitiativeId.subscribe(function(v){
	AdoptionModuleDetail.GetData()
})

AdoptionModuleDetail.Year.subscribe(function(v){
	AdoptionModuleDetail.GetData()
})

AdoptionModuleDetail.GetData = function(){
	var parm = {
		InitiativeId: AdoptionModuleDetail.InitiativeId(),
		year:parseInt(AdoptionModuleDetail.Year()),
		// MetricId:AdoptionModuleDetail.ActiveMetric(),
	}
	if(parm.InitiativeId != ""){
		ajaxPost("/web-cb/adoptionmodule/getdatadetail",parm,function(res){
			if(res.IsError){
				swal("Error",res.Message,"error")
				return false;
			}
			_.each(res.Data, function(v0,i0){
				_.each(v0.DetailData, function(v,i){
					mont = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
				    _.each(mont, function(vv,ii){
						if(res.Data[i0].DetailData[i]["NA"+vv+""]){
							res.Data[i0].DetailData[i][vv+""] = "";
						}
						if(res.Data[i0].DetailData[i]["NATotal"+vv]){
							res.Data[i0].DetailData[i]["Total"+vv] = "";
						}
				    })
				})
				res.Data[i0].DetailDataBack = res.Data[i0].DetailData;
			})

			// console.log(res.Data.MetricData)
			AdoptionModuleDetail.MetricData(ko.mapping.fromJS(res.Data)());
			AdoptionModuleDetail.MetricDataBack(ko.mapping.fromJS(res.Data)());
			if(res.Data.length > 0){
				AdoptionModuleDetail.ActiveMetric(res.Data[0].MetricId);
				AdoptionModuleDetail.ActualValue(res.Data[0].ActualValue)
				AdoptionModuleDetail.TotalValue(res.Data[0].TotalValue);
				AdoptionModuleDetail.RAGValue(res.Data[0].RAGValue);
				AdoptionModuleDetail.Data(ko.mapping.fromJS(res.Data[0].DetailData)())
				AdoptionModuleDetail.DataBack(ko.mapping.fromJS(res.Data[0].DetailData)())
			} else{
				var parm = {
					InitiativeId: AdoptionModuleDetail.InitiativeId(),
					year:parseInt(AdoptionModuleDetail.Year()),
					// MetricId:AdoptionModuleDetail.ActiveMetric(),
				}
				if(parm.InitiativeId != ""){
					ajaxPost("/web-cb/adoptionmodule/getdatadetailold",parm,function(res){
						if(res.IsError){
							swal("Error",res.Message,"error")
							return false;
						}

						_.each(res.Data.DetailData, function(v,i){
							mont = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
						    _.each(mont, function(vv,ii){
								if(res.Data.DetailData[i]["NA"+vv+""]){
									res.Data.DetailData[i][vv+""] = "";
								}
								if(res.Data.DetailData[i]["NATotal"+vv]){
									res.Data.DetailData[i]["Total"+vv] = "";
								}
						    })
						})

						AdoptionModuleDetail.Data(ko.mapping.fromJS(res.Data.DetailData)())
						AdoptionModuleDetail.DataBack(ko.mapping.fromJS(res.Data.DetailData)())
					})
				}
			}


			var BK = ko.mapping.toJS(AdoptionModuleDetail.MetricDataBack())
			// console.log(BK, BK.length)
			if(BK.length > 0){
				d = ko.mapping.toJS(BK[0])
				AdoptionModuleDetail.Data(ko.mapping.fromJS(BK[0].DetailData)())
				AdoptionModuleDetail.ActiveMetric(d.MetricId)
				AdoptionModuleDetail.ActualValue(d.ActualValue);
				AdoptionModuleDetail.TotalValue(d.TotalValue);
				AdoptionModuleDetail.RAGValue(d.RAGValue);
				AdoptionModuleDetail.InputType("Actual")
			} else{
				AdoptionModuleDetail.Data([])
				AdoptionModuleDetail.ActiveMetric("")
				AdoptionModuleDetail.ActualValue(false);
				AdoptionModuleDetail.TotalValue(false);
				AdoptionModuleDetail.RAGValue(false);
				AdoptionModuleDetail.InputType("")
			}

		})

	}
}

$(document).ready(function(){
	// // AdoptionModuleDetail.Init();
	// $("#InputType").bootstrapSwitch();
 //    $('#InputType').on('switchChange.bootstrapSwitch', function(event, state) {
 //    	AdoptionModuleDetail.InputType(state)
 //    });

    var Last10Years = AdoptionModuleDetail.Year()-10;
	var NextYears = AdoptionModuleDetail.Year()+1;
	for (var i=Last10Years;i<=NextYears;i++){
		AdoptionModuleDetail.YearList.push(i);
	}
})