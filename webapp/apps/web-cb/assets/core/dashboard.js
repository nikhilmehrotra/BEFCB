function monthDiff(d1, d2) {
    var months;
    months = (d2.getFullYear() - d1.getFullYear()) * 12;
    months -= d1.getMonth() + 1;
    months += d2.getMonth();
    return months <= 0 ? 0 : months;
}
var c = {
	OwnedData:ko.observableArray([]),
	PullRequest:ko.observable(0),
	LCWidth: ko.observable(0), //Set LifeCycle Width
	// LCWidth:ko.observable(screen.width-250), //Set LifeCycle Width
	LCFieldValue : ko.observable(100),
	SelectedSC: ko.observable(""),
	ActiveBDFilter:ko.observableArray([]),
	SelectedTab: ko.observable("KeyEnablers"),
	AllInitiateSource: ko.observable(true),
	Processing:ko.observable(true),
	HeaderText:ko.observable("Scorecard"),
	Loading:ko.observable(false),
	AllDataList:ko.observableArray([]),
	DataSource:ko.observable(),
	Loading:ko.observable(false),
	ColorList:ko.observableArray(["#0077AC","#009D44","#001143","#4EBF45","#6E7378","#00B3EB","#002BAE","#199489","#0077AC","#009D44","#001143","#4EBF45","#6E7378","#00B3EB","#002BAE","#199489"]),
	Filter:{
		Low:ko.observable(true),
		Medium:ko.observable(true),
		High:ko.observable(true),
		Primary:ko.observable(true),
		Secondary:ko.observable(true),
		CBLead:ko.observable(true),
		BankWide:ko.observable(true),
		YtdComplete:ko.observable(true),
		Remaining:ko.observable(true),
		Task:ko.observable(true),
		Investment:ko.observable(true),
		IsExellerator:ko.observable(true),
		IsOperationalExcellence:ko.observable(true),
		Region:ko.observableArray([]),
		Country:ko.observableArray([]),
		RegionCountry: ko.observableArray([]),
		RegionCountryScorecard: ko.observable(false),
		RegionOne:ko.observable(''),
		CountryOne:ko.observable(''),
		Search:ko.observable(""),
		IsBE:ko.observable(true),
		IsBP:ko.observable(true),
		StartDate:ko.observable(""),
		EndDate:ko.observable(""),
		DisplayColor:ko.observable('')
	},
	// DataSource
	UserList:ko.observableArray([]),
	LifeCycleList:ko.observableArray([]),
	RegionalData:ko.observableArray([]),
	InitiativeRegion:ko.observableArray([]),
	InitiativeCountry:ko.observableArray([]),
	RegionList:ko.observableArray([]),
	CountryList:ko.observableArray([]),
	iconAllExpand : ko.observable(true),
	dataSourceForTreeView:ko.observable(),

	//Header-Detail
	hd:{
		cblead: ko.observable(),
		bwide: ko.observable(),
		iHigh: ko.observable(),
		iMedium: ko.observable(),
		iLow:ko.observable(),
		ytd: ko.observable(),
		remain: ko.observable(),
		selected: ko.observable(),
		see: ko.observable(true),
		hideHeader: ko.observable(true)
	}
};
c.GetUserList = function(){
	ajaxPost("/web-cb/acluser/getuserlist",{},function(res){
		c.UserList(res.Data);
	});
}
c.GetUserFullname = function(loginid){
	var d = Enumerable.From(c.UserList()).Where("$.loginid === '"+loginid+"'").FirstOrDefault();
	return d.fullname;
}
c.IsComplete = function (eachRaw, comparator) {
	var each = ko.mapping.toJS(eachRaw)

	var setAsComplete = each.hasOwnProperty('SetAsComplete') ? each.SetAsComplete : false
	if (setAsComplete) {
		return (comparator == 1)
	}

	each.StartDate = getUTCDate(each.StartDate);
	each.FinishDate = getUTCDate(each.FinishDate);
	var Milestones = each.Milestones;
	for(var m in Milestones){
		Milestones[m].StartDate = getUTCDate(Milestones[m].StartDate);
		Milestones[m].EndDate = getUTCDate(Milestones[m].EndDate);
		Milestones[m].DaysBetween = Initiative.GetDaysBetween(Milestones[m].StartDate,Milestones[m].EndDate);
	}
    Milestones = Enumerable.From(Milestones).Where("$.Name.trim() !== ''").ToArray()
    var ProgressCompletion = 0;
    if (Milestones.length == 0){
    	if(Now>each.FinishDate){
    		ProgressCompletion = 100;
    	}else if(Now>=each.StartDate){
            ProgressCompletion = (Initiative.GetDaysBetween(each.StartDate,each.FinishDate) == 0 ? 0 : Initiative.GetDaysBetween(each.StartDate,Now)/Initiative.GetDaysBetween(each.StartDate,each.FinishDate))*100;
        }
    } else {
        var TotalDays = Enumerable.From(Milestones).Sum("$.DaysBetween");
        
    	
        for(var mt in Milestones){            
            var StartDate = Milestones[mt].StartDate;
            var EndDate = Milestones[mt].EndDate;
            
            var weight = TotalDays == 0 ? 0 : Initiative.GetDaysBetween(StartDate,EndDate)/TotalDays;
            var progress = 0;
            if(Now>=StartDate){
                if(Now>=EndDate){
                    progress = weight;
                }else{
                    progress = (Initiative.GetDaysBetween(StartDate,EndDate) === 0 ? 0 : Initiative.GetDaysBetween(StartDate,Now)/Initiative.GetDaysBetween(StartDate,EndDate))*weight;
                }
            }
            ProgressCompletion+=(progress*100);
        }

    }
    if(ProgressCompletion>=100){
    	return (comparator == 2)
    }

    return (comparator == 0)
}

c.Completion = ko.computed(function(){
	if(c.DataSource() !== undefined){
		return Enumerable.From(c.DataSource().Data.SummaryBusinessDriver).Sum("$.Completion")/c.DataSource().Data.SummaryBusinessDriver.length;
	}else{
		return 0;
	}
})
c.Filter.Region.subscribe(function(newVal){
	var arr = [];
	if(newVal.length > 0){
		for(var i in newVal){
			if(newVal[i]=="GLOBAL" ){
				arr = []
				for(i in c.RegionalData()){
					arr.push({"_id": c.RegionalData()[i].Country })
				}
				break
			}else{
				var temp_arr = Enumerable.From(c.RegionalData()).Where("$.Major_Region === '"+newVal[i]+"'").GroupBy("$.Country").Select("{_id:$.Key()}").ToArray();
				arr = arr.concat(temp_arr);
			}
		}
	} else{
		for(i in c.RegionalData()){
			arr.push({"_id": c.RegionalData()[i].Country })
		}
	}

  c.CountryList(arr);
  c.CountryList.unshift({"_id":"GLOBAL"})
  c.Filter.Country([])
  c.GetData();
});


// c.Filter.RegionOne.subscribe(function(){
//     Scorecard.GetData();
// });
// c.Filter.CountryOne.subscribe(function(){
//     Scorecard.GetData();
// });
c.Filter.RegionCountry.subscribe(function(){
	c.GetData();
});
c.Filter.StartDate.subscribe(function(){
    c.GetData();
});
c.Filter.EndDate.subscribe(function(){
    c.GetData();
});
c.Filter.Country.subscribe(function(){
    c.GetData();
});
c.Filter.Low.subscribe(function(){
	c.GetData();
});
c.Filter.Medium.subscribe(function(){
	c.GetData();
});
c.Filter.High.subscribe(function(){
	c.GetData();
});
c.Filter.Primary.subscribe(function(){
	c.GetData();
});
c.Filter.Secondary.subscribe(function(){
	c.GetData();
});
c.Filter.YtdComplete.subscribe(function(){
	c.GetData();
});
c.Filter.Remaining.subscribe(function(){
	c.GetData();
});
c.Filter.Investment.subscribe(function(){
	c.GetData();
});
c.Filter.CBLead.subscribe(function(){
	c.GetData();
});
c.Filter.BankWide.subscribe(function(){
	c.GetData();
});

c.Filter.DisplayColor.subscribe(function () {
	c.GetData();
});

c.Filter.Task.subscribe(function(){
	c.GetData();
});
c.Filter.IsExellerator.subscribe(function(){
	c.GetData();
});
c.Filter.IsOperationalExcellence.subscribe(function(){
	c.GetData();
});

c.Filter.IsBE.subscribe(function(){
	c.GetData();
});
c.Filter.IsBP.subscribe(function(){
	c.GetData();
});

c.SelectedTab.subscribe(function(){
	setTimeout(function(){
		if(SortInitiative.Active() && c.SelectedTab() == "Initiative" ){
			Initiative.FixedHeader();
		}
		if(!SortInitiative.Active() || c.SelectedTab() != "Initiative" ){
			Initiative.RemoveFixedHeader();
		}
	},300)
	c.GetData();
});


// c.ActiveBDFilter.subscribe(function(){
// 	setTimeout(function(){
// 		c.GetData();
// 	},100);
// });

var topTable = 0;
c.GetHeader = function(isPNG){
	var yourheder = "Commercial Banking ";
	yourheder = yourheder + c.HeaderText();
	if(c.SelectedTab() == 'Initiative' || c.SelectedTab() == 'SharedAgenda'){
		if(c.Filter.Country().length > 0){
			var countries = "("+c.Filter.Country().join(', ')+")";
			yourheder = yourheder+" : Country "+countries;
		}else if(c.Filter.Region().length > 0){
			var regions = "("+c.Filter.Region().join(', ')+")";
			yourheder = yourheder+" : Region "+regions;
		}else{
			yourheder= yourheder+" : Global";
		}
		
	}else{
		if(c.Filter.CountryOne() != "" && c.Filter.CountryOne() != "Country"){
			var country = "("+c.Filter.CountryOne()+")";
			yourheder = yourheder+ " : Country "+country;
		}else if(c.Filter.RegionOne() != "" && c.Filter.RegionOne() != "Region"){
			var region = "("+c.Filter.RegionOne()+")";
			yourheder = yourheder+" : Regional Scorecard "+region;
		}else{
			yourheder= yourheder+" : Global Scorecard";
		}

	}
	// console.log(yourheder)
	var $tableChart = $('#tabelChart'), $redis = $('#redips-drag');
	$tableChart.height($redis.height());
	$('#ContentHeader').text(yourheder);


	if(SortInitiative.Active()){
		// Initiative.RemoveFixedHeader();
		$(window).scrollTop(0)
		$('#ContentHeader').css({
			"position": "fixed",
		    "z-index": "3",
		    "top": "80px",
		    "width": "100%",
		    "background": "white"
		})
		var topheader = $('#app-header .navbar').height();
	    var topheight = $('#top-header').height();
	    var iabdHeight = $('#initiativeAllBDF').height() ;
	    var crumbsMargin = $('#BDFilter').width();
	    var iabdTop = (topheight+topheader) ;
	    var bdfTop = iabdTop + iabdHeight;
	    $('#initiativeAllBDF').css({
	        "position": "fixed",
	        "top": (iabdTop+25)+"px",
	        "z-index": "2",
	        "background-color": "white",
	        "width": "100%",
	        "height": iabdHeight+"px",
	    })
	    $('#selectall-initiative').css({
	        "width":"15%"
	    })
	    $('#crumbs').css({
	        "position": "fixed",
	        "top": (iabdTop+25)+"px",
	        "bottom": "0px",
	        "left": "0px",
	        "right": "0px",
	        "z-index": "4",
	        "margin-left": crumbsMargin+"px",
	        "width": "82.5%",
	    })
	    
	    $('#BDFilter').css({
	        "position": "absolute",
	        "top": (bdfTop+25)+"px",
	    })

    	$('#iFooter').css("margin-top", (bdfTop+25)+"px")
	    if(topTable===0){
	    	topTable = $("#tabelChart").position().top;
		}
		if(isPNG){
			$("#tabelChart").css({"top":(topTable+10)+"px"})
		}else{
	    	$("#tabelChart").css({"top":(topTable+15)+"px"})
		}
    
	}else{
		$('#ContentHeader').css({
			"position": "",
		    "z-index": "",
		    "top": "",
		    "width": "",
		    "background": ""
		})
	}
}
c.RestoreToNormal = function(){
	var topheader = $('#app-header .navbar').height();
    var topheight = $('#top-header').height();
    var iabdHeight = $('#initiativeAllBDF').height() ;
    var crumbsMargin = $('#BDFilter').width();
    var iabdTop = (topheight+topheader) ;
    var bdfTop = iabdTop + iabdHeight;
    
    $('#initiativeAllBDF').css({
        "position": "fixed",
        "top": (iabdTop+10)+"px",
        "z-index": "2",
        "background-color": "white",
        "width": "100%",
        "height": iabdHeight+"px",
    })
    $('#selectall-initiative').css({
        "width":"15%"
    })
    $('#crumbs').css({
        "position": "fixed",
        "top": (iabdTop+10)+"px",
        "bottom": "0px",
        "left": "0px",
        "right": "0px",
        "z-index": "4",
        "margin-left": crumbsMargin+"px",
        "width": "82.5%",
    })
    
    $('#BDFilter').css({
        "position": "absolute",
        "top": (bdfTop+10)+"px",
    })

	$('#iFooter').css("margin-top", (bdfTop+10)+"px")
    $("#tabelChart").css({"top":(topTable-25)+"px"})
}
c.ExportAsPDF = function () {
	kendo.pdf.defineFont({
	    // "Open Sans": "/static/fonts/OpenSans-Regular.ttf",
	    "Helvetica Neue": "/web-cb/static/fonts/HelveticaNeue.ttf",
	});
	c.GetHeader(false);
	var $tableChart = $('#tabelChart'), $redis = $('#redips-drag');
	var tableChartHeight = $tableChart.height();
	var $container = $('.tab-content');
    $('#ContentHeader').show();
	kendo.drawing.drawDOM($container,{
		template: $("#page-template").html(), 
	}).then(function (group) {
		var title = $($('#secondmenu .active a')[0]).text() + ".pdf";
        kendo.drawing.pdf.saveAs(group, title);
        $('#ContentHeader').hide();
		if(SortInitiative.Active()){
			c.RestoreToNormal();
    	}else{
    		$tableChart.height(tableChartHeight);
    	}
    })

}

c.ExportAsPNG = function () {
	c.GetHeader(true);
	var $tableChart = $('#tabelChart'), $redis = $('#redips-drag')
	var tableChartHeight = $tableChart.height()
	var $container = $('.tab-content');
	// $tableChart.height($redis.height())
	$('.tab-content').css('background', 'white').css('padding','10px');
	
	$('#ContentHeader').show();
	// $('#ContentHeader').show()
	// $('#ContentHeader').css({
	// 	"display" : "block",
	// 	"opacity" : "1"
	// })
	// kendo.drawing.drawDOM($container).then(function (group) {
	// 	var title = $('#secondmenu .active a').text() + ".png";
 //        kendo.drawing.exportImage.saveAs(group, title);
 //        $('#ContentHeader').hide()
	// 	$tableChart.height(tableChartHeight);
 //    })
	kendo.drawing.drawDOM($container) //#qrMail
	.then(function(group) {
		$('#ContentHeader').show();
		// Render the result as a PNG image
		return kendo.drawing.exportImage(group);
	})
	.done(function(data) {
		// Save the image file
		kendo.saveAs({
			dataURI: data,
			fileName: $($('#secondmenu .active a')[0]).text()+ ".png",
			// proxyURL: "http://demos.telerik.com/kendo-ui/service/export"
		});
		$('#ContentHeader').hide();
		if(SortInitiative.Active()){
			c.RestoreToNormal();
    	}
		$('.tab-content').css('background', 'white').css('padding','0');
	});
}

c.Clear = function (selector) {
	$('#search-input').data('kendoAutoComplete').value('');	
	setTimeout(function () {
		Initiative.CurrentSearchKeyword('');
		search.FilterInitiative('')
		redipsInit();
	}, 100)
}

c.AddGLobal = function(){
	c.RemoveGlobal()
	var value={_id:"GLOBAL"}
	c.RegionList.unshift(value)
	c.CountryList.unshift(value)
	
}

c.RemoveGlobal = function(){
	var value="GLOBAL"
	c.RegionList.remove(function (item) { return item._id ==value; })
	c.CountryList.remove(function (item) { return item._id ==value; })
}

c.MappingDataSource = function(sources){
	// console.log(sources)
	var allInitiate = [];

	for(var i in sources){
		// Check for Completed Initiative
		for(var ip in sources[i].Project){
			sources[i].Project[ip].IsCompleted = false;

			if(sources[i].Project[ip].SetAsComplete){
				sources[i].Project[ip].IsCompleted = true;

				sources[i].Project[ip].StartDate = getUTCDate(sources[i].Project[ip].StartDate);
				sources[i].Project[ip].FinishDate = getUTCDate(sources[i].Project[ip].FinishDate);
				var Milestones = sources[i].Project[ip].Milestones;

				for(var m in Milestones){
					Milestones[m].Id = parseInt(m);
					Milestones[m].StartDate = getUTCDate(Milestones[m].StartDate);
					Milestones[m].EndDate = getUTCDate(Milestones[m].EndDate);
					Milestones[m].DaysBetween = Initiative.GetDaysBetween(Milestones[m].StartDate,Milestones[m].EndDate);
					
					if(Milestones[m].Seq == undefined){
						Milestones[m].Seq = parseInt(m) + 1;
					}

				}

				// sources[i].Project[ip].FinishDate =  getUTCDate(sources[i].Project[ip].CompletedDate); -- to enable tick mark logic in Chart Section [Logic:Ask Ainur]
			}else{

				sources[i].Project[ip].StartDate = getUTCDate(sources[i].Project[ip].StartDate);
				sources[i].Project[ip].FinishDate = getUTCDate(sources[i].Project[ip].FinishDate);
				var Milestones = sources[i].Project[ip].Milestones;

				for(var m in Milestones){
					Milestones[m].Id = parseInt(m);
					Milestones[m].StartDate = getUTCDate(Milestones[m].StartDate);
					Milestones[m].EndDate = getUTCDate(Milestones[m].EndDate);
					Milestones[m].DaysBetween = Initiative.GetDaysBetween(Milestones[m].StartDate,Milestones[m].EndDate);
					
					if(Milestones[m].Seq == undefined){
						Milestones[m].Seq = parseInt(m) + 1;
					}

				}

        Milestones = Enumerable.From(Milestones).Where("$.Name.trim() !== ''").ToArray()
        var ProgressCompletion = 0;
        if (Milestones.length == 0){
        	if(Now>sources[i].Project[ip].FinishDate){
                ProgressCompletion = 100;
            }else if(Now>=sources[i].Project[ip].StartDate){
                ProgressCompletion = (Initiative.GetDaysBetween(sources[i].Project[ip].StartDate,sources[i].Project[ip].FinishDate) == 0 ? 0 : Initiative.GetDaysBetween(sources[i].Project[ip].StartDate,Now)/Initiative.GetDaysBetween(sources[i].Project[ip].StartDate,sources[i].Project[ip].FinishDate))*100;
            }
        } else {
            var TotalDays = Enumerable.From(Milestones).Sum("$.DaysBetween");
            
        	
            for(var mt in Milestones){            
                var StartDate = Milestones[mt].StartDate;
                var EndDate = Milestones[mt].EndDate;
                
                var weight = TotalDays == 0 ? 0 : Initiative.GetDaysBetween(StartDate,EndDate)/TotalDays;
                var progress = 0;
                if(Now>=StartDate){
                    if(Now>=EndDate){
                        progress = weight;
                    }else{
                        progress = (Initiative.GetDaysBetween(StartDate,EndDate) === 0 ? 0 : Initiative.GetDaysBetween(StartDate,Now)/Initiative.GetDaysBetween(StartDate,EndDate))*weight;
                    }
                }
                ProgressCompletion+=(progress*100);
            }

        }

    //     if(sources[i].Project[ip]._id == "58b7a45b8365f42bc1515711"){
				// 	console.log(Milestones)
				// }

        if(ProgressCompletion>=100){
        	sources[i].Project[ip].IsCompleted = true;
        }
			}
		}

		if(!c.Filter.YtdComplete()){
			sources[i].Project = Enumerable.From(sources[i].Project).Where("!$.IsCompleted").ToArray();
		}
		if(!c.Filter.Remaining()){
			sources[i].Project = Enumerable.From(sources[i].Project).Where("$.IsCompleted").ToArray();
		}


		
		allInitiate = allInitiate.concat(sources[i].Project)

		for(var ii in allInitiate){
			if(allInitiate[ii].SetAsComplete == undefined || allInitiate[ii].SetAsComplete == null){
				allInitiate[ii].SetAsComplete = false;
			}

			if(allInitiate[ii].CompletedDate == undefined || allInitiate[ii].CompletedDate == null){
				allInitiate[ii].CompletedDate = new Date();
			}

			if(allInitiate[ii].DisplayProgress == undefined || allInitiate[ii].DisplayProgress == null){
				allInitiate[ii].DisplayProgress = "green";
			}

			if(allInitiate[ii].Sponsor == undefined || allInitiate[ii].Sponsor == null){
				allInitiate[ii].Sponsor = "";
			}
		}

		var lifecycle = c.LifeCycleList();

		sources[i].SummaryBusinessDriver = _.orderBy(sources[i].SummaryBusinessDriver, 'SeqParent');

		zz = _.groupBy(sources[i].SummaryBusinessDriver, 'Parentid')
		lenghtparent = sources[i].SummaryBusinessDriver.length
		showparentList = [0]
		totalparentList = []
		tmp = 0;
		for(zi in zz){
			totalparentList.push(zz[zi].length)
			tmp += zz[zi].length
			if(tmp < lenghtparent){
				showparentList.push(tmp)
			}
		}

		var businessDriver = sources[i].SummaryBusinessDriver;
		_.each(businessDriver, function(v,i){
			// _.each(v.BusinessMetrics, function(vv,ii){

			// 	if(vv.ActualData != null){
			// 			zz = _.find(vv.ActualData, function (o) { return o.Flag == 'C1'; });
			// 			if(zz != undefined){
			// 				businessDriver[i].BusinessMetrics[ii].tmpCurrentValue = zz.Value;
			// 			} else{
			// 				vv.ActualData = _.sortBy(vv.ActualData, 'Period');
			// 				businessDriver[i].BusinessMetrics[ii].tmpCurrentValue = vv.ActualData[vv.ActualData.length-1].Value;
			// 			}
			// 	} else{
			// 		businessDriver[i].BusinessMetrics[ii].tmpCurrentValue = 0;
			// 	}


			// 	if(vv.ActualData != null){
			// 		vv.ActualData = _.sortBy(vv.ActualData, 'Period')
			// 		businessDriver[i].BusinessMetrics[ii].tmpMinValue = (vv.ActualData[0] != undefined) ? vv.ActualData[0].Value : 0;
			// 		// console.log(vv.ActualData[vv.ActualData.length-1].Value)
			// 	} else{
			// 		businessDriver[i].BusinessMetrics[ii].tmpMinValue = 0
			// 	}
			// })
			showparent = $.inArray(i, showparentList) != -1;
			totalparent = $.inArray(i, showparentList) != -1 ? totalparentList[$.inArray(i, showparentList)] : 0;
			businessDriver[i].ShowParent = showparent;
			businessDriver[i].TotalParent = totalparent;
			businessDriver[i].ShowHide = ko.observable(false);
			businessDriver[i].ShowHideSecLvl = ko.observable(false);
		})

		sources[i].BusinessDriverData = [];
		for(var bd in businessDriver){
			businessDriver[bd]._id = businessDriver[bd].Id;
			businessDriver[bd].Id = businessDriver[bd].Idx;

			// Mapping Data Metrics Chart
			businessDriver[bd].BusinessMetrics = businessDriver[bd].BusinessMetrics === null ? [] : businessDriver[bd].BusinessMetrics;
			// var BusinessMetrics = businessDriver[bd].BusinessMetrics;
			// for(var bm in BusinessMetrics){
			// 	BusinessMetrics[bm].MaxValue = BusinessMetrics[bm].TargetValue
			// 	BusinessMetrics[bm].MinValue = BusinessMetrics[bm].ActualData.length === 0 ? 0 : BusinessMetrics[bm].ActualData[0].value;
			// }

			var Initiatives = Enumerable.From(sources[i].Project).Where("$.BusinessDriverId === '"+businessDriver[bd].Id+"'").ToArray();
			var DataCompletion = Enumerable.From(Initiatives).Sum("$.ProgressCompletion")
			businessDriver[bd].Completion = Initiatives.length > 0 ? DataCompletion/Initiatives.length : 0;
			var TargetRange = [businessDriver[bd].MinValue,businessDriver[bd].MaxValue];
			var d = {
				Id:businessDriver[bd].Id,
				Parentid:businessDriver[bd].Parentid,
				Parentname:businessDriver[bd].Parentname,
				Category:businessDriver[bd].Category,
				ShowParent:businessDriver[bd].ShowParent,
				TotalParent:businessDriver[bd].TotalParent,
				Name:businessDriver[bd].Name,
				DataPoint:businessDriver[bd].DataPoint,
				Initiative:Initiatives.length,
				Completion:businessDriver[bd].Completion,
				MinTarget:businessDriver[bd].MinLiabilities,
				MaxTarget:businessDriver[bd].MaxLiabilities,
				MetricType:businessDriver[bd].MetricType,
				IsEditCompletion:false,
				EditableCompletion:ko.observable(businessDriver[bd].Completion),
				TargetRange:TargetRange,
			};
			sources[i].BusinessDriverData.push(d);
		}

		// Creating Table Sources Data
		sources[i].TableSources = [];
		var Connector = "";
		for(var bd in businessDriver){
			var bdData = {
				Idx:bd,
				Id:businessDriver[bd].Id,
				Name:businessDriver[bd].Name,
				Parentid:businessDriver[bd].Parentid,
				Parentname:businessDriver[bd].Parentname,
				Category:businessDriver[bd].Category,
				ShowParent:businessDriver[bd].ShowParent,
				TotalParent:businessDriver[bd].TotalParent,
				LifeCycle:[],
				ShowHide:false,
				ShowHideSecLvl:false,
			}
			for(var lc in lifecycle){
				var lcData = {
					Idx:lc,
					Id:lifecycle[lc].Id,
					Name:lifecycle[lc].Name,
					Seq:lifecycle[lc].Seq,
				};
				if(Connector!==""){
					Connector += ",#"+bdData.Id+lcData.Id;
				}else{
					Connector += "#"+bdData.Id+lcData.Id;
				}
				var Initiatives = Enumerable.From(sources[i].Project).Where("$.LifeCycleId === '"+lcData.Id+"' && $.BusinessDriverId === '"+bdData.Id+"'").ToArray()
				// var Initiatives = Enumerable.From(sources[i].Project).Where("$.LifeCycleId === '"+lcData.Id+"'").ToArray()
				for(var init in Initiatives){
					Initiatives[init].IsTask = false;
				}
				var Tasks = Enumerable.From(sources[i].TaskList).Where("$.LifeCycleId === '"+lcData.Id+"' && $.BusinessDriverId === '"+bdData.Id+"'").ToArray()
				// var Tasks = Enumerable.From(sources[i].TaskList).Where("$.LifeCycleId === '"+lcData.Id+"'").ToArray()
		    for(var t in Tasks){
					Initiatives.push({
				    "_id" : "task"+ Tasks[t].Id,
				    "ProjectName" : Tasks[t].Name,
				    // "Name" : t.Name,
				    "IsTask": true,
				    "BDId" : Tasks[t].BusinessDriverId,
				    "LCId": Tasks[t].LifeCycleId,
				    "InitiativeType" : Tasks[t].TaskType,
				    "BusinessDriverImpact" : "Primary",
				    "CBLedInitiatives":false,
				    "type":"",
				    "InitiativeID":"",
				    "Attachments":[],
				    "BusinessImpact":"",
				    "ProgressCompletion":0,
				    "EX":false,
				    "OE":false,
					});

				}
				lcData.Initiatives = Initiatives;
				bdData.LifeCycle.push(lcData);
			}
			sources[i].TableSources.push(bdData);
		}

		//New Table Source
		sources[i].TableSourcesVer2 = [];
		for(var lc in lifecycle){
			var lcData = {
				Idx:lc,
				Id:lifecycle[lc].Id,
				Name:lifecycle[lc].Name,
				Seq:lifecycle[lc].Seq,
			};

			// var Initiatives = Enumerable.From(sources[i].Project).Where("$.LifeCycleId === '"+lcData.Id+"' && $.BusinessDriverId === '"+bdData.Id+"'").ToArray()
			var Initiatives = Enumerable.From(sources[i].Project).Where("$.LifeCycleId === '"+lcData.Id+"'").ToArray()
			// console.log(Initiatives);
			for(var init in Initiatives){
				Initiatives[init].IsTask = false;
			}
			// var Tasks = Enumerable.From(sources[i].TaskList).Where("$.LifeCycleId === '"+lcData.Id+"' && $.BusinessDriverId === '"+bdData.Id+"'").ToArray()
			var Tasks = Enumerable.From(sources[i].TaskList).Where("$.LifeCycleId === '"+lcData.Id+"'").ToArray()
	    for(var t in Tasks){
        // console.log("-->koplo-->",Tasks[t])
				Initiatives.push({
			    "_id" : "task"+ Tasks[t].Id,
			    "ProjectName" : Tasks[t].Name,
			    // "Name" : t.Name,
			    "IsTask": true,
			    "BDId" : Tasks[t].BusinessDriverId,
			    "LCId": Tasks[t].LifeCycleId,
			    "InitiativeType" : Tasks[t].TaskType,
			    "BusinessDriverImpact" : "Primary",
			    "CBLedInitiatives":false,
			    "type":"",
			    "InitiativeID":"",
			    "Attachments":[],
			    "BusinessImpact":"",
			    "ProgressCompletion":0,
			    "EX":false,
			    "OE":false,
          "Owner":Tasks[t].Owner,
          "Statement":Tasks[t].Statement,
          "Description":Tasks[t].Description,
				});

			}
			lcData.Initiatives = Initiatives;
			sources[i].TableSourcesVer2.push(lcData);
		}

		sources[i].TableSourcesVer3AlignVer = [];

		for(var bd in sources[i].BusinessDriverL1){
			var TableSourcesVer2 = []
			var BDIdDefaultObj = _.find(businessDriver, function (o) { return o.Parentid == sources[i].BusinessDriverL1[bd].Idx; });
			var BDIdDefault = BDIdDefaultObj == undefined ? "" : BDIdDefaultObj.Idx
			var bdData = {
				Id:sources[i].BusinessDriverL1[bd].Idx,
				BDIdDefault:BDIdDefault,
				Name:sources[i].BusinessDriverL1[bd].Name,
			}
			
			for(var lc in lifecycle){
				var lcData = {
					Idx:lc,
					Id:lifecycle[lc].Id,
					Name:lifecycle[lc].Name,
					Seq:lifecycle[lc].Seq,
				};

				var Initiatives = Enumerable.From(sources[i].Project).Where("$.LifeCycleId === '"+lcData.Id+"' && $.SCCategory === '"+bdData.Id+"'").ToArray()
				// var Initiatives = Enumerable.From(sources[i].Project).Where("$.LifeCycleId === '"+lcData.Id+"'").ToArray()
				// console.log(Initiatives);
				for(var init in Initiatives){
					Initiatives[init].IsTask = false;
				}
				var Tasks = Enumerable.From(sources[i].TaskList).Where("$.LifeCycleId === '"+lcData.Id+"' && $.SCCategory === '"+bdData.Id+"'").ToArray()
				// var Tasks = Enumerable.From(sources[i].TaskList).Where("$.LifeCycleId === '"+lcData.Id+"'").ToArray()
		    for(var t in Tasks){
					Initiatives.push({
				    "_id" : "task"+ Tasks[t].Id,
				    "ProjectName" : Tasks[t].Name,
				    // "Name" : t.Name,
				    "IsTask": true,
				    "BDId" : Tasks[t].BusinessDriverId,
				    "LCId": Tasks[t].LifeCycleId,
				    "InitiativeType" : Tasks[t].TaskType,
				    "BusinessDriverImpact" : "Primary",
				    "CBLedInitiatives":false,
				    "type":"",
				    "InitiativeID":"",
				    "Attachments":[],
				    "BusinessImpact":"",
				    "ProgressCompletion":0,
				    "EX":false,
				    "OE":false,
					});

				}
				lcData.Initiatives = Initiatives;
				TableSourcesVer2.push(lcData);
			}

			var a = {
				"Id":bdData.Id,
				"BDIdDefault":bdData.BDIdDefault,
				"TableSourcesVer2":ko.mapping.fromJS(TableSourcesVer2)
			}
			// console.log(a)
			sources[i].TableSourcesVer3AlignVer.push(a)
		}

		var ScorecardList = []
		_.each(_.groupBy(sources[i].SummaryBusinessDriver, 'Parentid'), function(v,i){
		  tmpLvl1 = {}
		  tmpLvl1.Id = i;
		  tmpLvl1.Name = v[0].Parentname;
		  tmpLvl1.Data = []
		  _.each(v, function(vv,ii){
		    tmpLvl2 = {}
		    tmpLvl2.Id = vv.Id;
		    tmpLvl2.Name = vv.Name;
		    tmpLvl1.Data.push(tmpLvl2)
		  })
		  ScorecardList.push(tmpLvl1)
		})

		sources[i].ScorecardList = ScorecardList;

		sources[i].Connector = Connector;
		sources[i].TableSources = ko.mapping.fromJS(sources[i].TableSources)();
		sources[i].TableSourcesVer2 = ko.mapping.fromJS(sources[i].TableSourcesVer2)();
		sources[i].TableSourcesVer2Backup = ko.mapping.fromJS(sources[i].TableSourcesVer2)();
		sources[i].TableSourcesVer3 = ko.observableArray([]);
		sources[i].TableSourcesVer3BackupAll = [{Id:"All","BDIdDefault":"",TableSourcesVer2: sources[i].TableSourcesVer2}];
		if (SortInitiative.Active()){
        sources[i].TableSourcesVer3(sources[i].TableSourcesVer3AlignVer)
    } else{
        sources[i].TableSourcesVer3(sources[i].TableSourcesVer3BackupAll)
    }
		
		// sources[i].TableSourcesVer3AlignVer = [{TableSourcesVer2: sources[i].TableSourcesVer2},{TableSourcesVer2: sources[i].TableSourcesVer2},{TableSourcesVer2: sources[i].TableSourcesVer2}];
		//console.log(sources[i]);
	}
	var iHigh = Enumerable.From(sources["Initiative"].Project).Where("$.BusinessImpact === 'High'").Count();
	var iMedium = Enumerable.From(sources["Initiative"].Project).Where("$.BusinessImpact === 'Medium'").Count();
	var iLow = Enumerable.From(sources["Initiative"].Project).Where("$.BusinessImpact === 'Low'").Count();
	var typeLed = Enumerable.From(sources["Initiative"].Project).Where("$.type === 'CBLED'").Count();
	var typeWide = Enumerable.From(sources["Initiative"].Project).Where("$.type === 'BANKWIDE'").Count();
	var iYTD = Enumerable.From(sources["Initiative"].Project).Where("$.IsCompleted").Count();
	var iRemain =  Enumerable.From(sources["Initiative"].Project).Where("!$.IsCompleted").Count();

	if(!c.Filter.Investment()){
			sources[i].Project = Enumerable.From(sources[i].Project).Where("$.InvestmentId === ''").ToArray();
	}

	if(c.Filter.YtdComplete()){
		c.hd.ytd(iYTD);	
	}
	if(c.Filter.Remaining()){
		c.hd.remain(iRemain);	
	}
	if(c.Filter.High()){
		c.hd.iHigh(iHigh);	
	}
	if(c.Filter.Medium()){
		c.hd.iMedium(iMedium);	
	}
	if(c.Filter.Low()){
		c.hd.iLow(iLow);	
	}
	if(c.Filter.CBLead()){
		c.hd.cblead(typeLed);
	}
	if(c.Filter.BankWide())	{
		c.hd.bwide(typeWide);	
	}
	c.AllInitiateSource(allInitiate);
	c.DataSource({
		Data: sources["Initiative"]
	});
	c.AllDataList(sources["Initiative"].Project)



	// setTimeout(function() {
	// 	if (SortInitiative.Active()){
 //      SortInitiative.sycHeight();
 //    }
 //  	}, 1000);

	// console.log(c.DataSource().Data.AllSummaryBusinessDriver);
	// var lengthHeader = $(".lcHeader").length;
	// var widthCl = Math.round(( screen.width - ($("#progressChart").width() + $(".redips-mark").width() + 30))/lengthHeader);
	// var widthEx = ($(".scrolldiv ").find("table tbody tr").find(".redips-mark").width()+19)*2;
	// $(".lcHeader" ).css({"width": widthCl+"px"});
	// $(".lcvalue" ).css({"width": widthCl+"px"});
	// $("#boxAllExpand").css({"padding-righ": widthEx +"px"});
	// $("#crumbs").css({"width":(screen.width-($("#progressChart").width()+130))+"px"});
	// $("#crumbsData").css({"width":(screen.width-($("#progressChart").width()+130))+"px"});
	// c.LCFieldValue(filedvalue);
}

c.ChangeCountry = function(e){
    var d = e.sender._old;
    var arr = [];
    if(d != "Region"){
            var temp_arr = Enumerable.From(c.RegionalData()).Where("$.Major_Region === '"+d+"'").GroupBy("$.Country").Select("{_id:$.Key()}").ToArray();
            arr = arr.concat(temp_arr);
    } else{
        for(i in c.RegionalData()){
            arr.push({"_id": c.RegionalData()[i].Country })
        }
    }
            
    c.CountryList(arr);   
    c.Filter.CountryOne('');
}

c.TooltipsterBuild = function(){
	var InitiativeType = ["KeyEnablers", "SupportingEnablers"]
	var bdAll = Enumerable.From(c.AllInitiateSource()).Select("$.BusinessDriverId").Distinct().ToArray();
	var lcAll = Enumerable.From(c.AllInitiateSource()).Select("$.LifeCycleId").Distinct().ToArray();
	// InitiativeType.forEach(function(IT){
		// var source = Enumerable.From(c.AllInitiateSource()).Where("$.InitiativeType == '" + c.SelectedTab() + "'").ToArray();
		var source = Enumerable.From(c.AllInitiateSource()).Where("$.InitiativeType == 'KeyEnablers'").ToArray();

		var source2 = ko.mapping.toJS(c.DataSource().Data.TableSourcesVer2);
		// console.log("-->",source)

		bdAll.forEach(function(bd){
			//var title = Enumerable.From(c.AllInitiateSource()).FirstOrDefault()
			var konten = "<center>Details: </center>";
			var total = Enumerable.From(source).Where("$.BusinessDriverId == '" + bd + "'").ToArray();
			var primaryCount = Enumerable.From(source).Where("$.BusinessDriverImpact == 'Primary' && $.BusinessDriverId == '" + bd + "'").ToArray();
			var SecondaryCount = Enumerable.From(source).Where("$.BusinessDriverImpact == 'Secondary' && $.BusinessDriverId == '" + bd + "'").ToArray();
			var CBLITrueCount = Enumerable.From(source).Where("$.CBLedInitiatives == true && $.BusinessDriverId == '" + bd + "'").ToArray();
			var high =  Enumerable.From(source).Where("$.BusinessDriverId == '" + bd + "' && $.BusinessImpact == 'High'" ).ToArray();
			var middle =  Enumerable.From(source).Where("$.BusinessDriverId == '" + bd + "' && $.BusinessImpact == 'Middle'" ).ToArray();
			var low =  Enumerable.From(source).Where("$.BusinessDriverId == '" + bd + "' && $.BusinessImpact == 'Low'" ).ToArray();

			//CBLedInitiatives
			// konten += "<br /># of Primary Initiatives: "+primaryCount.length;
			// konten += "<br /># of Secondary Initiatives: "+SecondaryCount.length;
			konten += "<br /># of CB Led Initiatives: "+CBLITrueCount.length;
			konten += "<br /># of Bank Wide Initiatives: -";
			konten += "<br /># of High: " + high.length;
			konten += "<br /># of Middle: "+middle.length;
			konten += "<br /># of Low: "+low.length;
			konten += "<br /># of Task: "+total.length;
			if($('a[toltipsterid='+bd+']').hasClass("tooltipstered"))
				$('a[toltipsterid='+bd+']').tooltipster("destroy")
			$('a[toltipsterid='+bd+']').tooltipster({
				contentAsHTML: true,
				content: konten,
				trigger: 'hover'
			});
		})

		lcAll.forEach(function(lc){
			var konten = "<center>Details: </center>";//InitiativeType
			// var total = Enumerable.From(source).Where("$.LifeCycleId == '" + lc + "'").ToArray();
			var datamentah = Enumerable.From(source2).Where("$.Id == '" + lc + "'").ToArray();
			var datamentah = (datamentah.length > 0) ? datamentah[0].Initiatives : [];
			// var total = datamentah.length

			// var primaryCount = Enumerable.From(source).Where("$.BusinessDriverImpact == 'Primary' && $.LifeCycleId == '" + lc + "'").ToArray();
			// var SecondaryCount = Enumerable.From(source).Where("$.BusinessDriverImpact == 'Secondary' && $.LifeCycleId == '" + lc + "'").ToArray();
			// var CBLITrueCount = Enumerable.From(source).Where("$.CBLedInitiatives == true && $.LifeCycleId == '" + lc + "'").ToArray();
			var total = Enumerable.From(datamentah).Where("$.IsTask == true").ToArray();
			var CBLITrueCount = Enumerable.From(datamentah).Where("$.CBLedInitiatives == true").ToArray();
			var high =  Enumerable.From(datamentah).Where("$.BusinessImpact == 'High'" ).ToArray();
			var middle =  Enumerable.From(datamentah).Where("$.BusinessImpact == 'Middle'" ).ToArray();
			var low =  Enumerable.From(datamentah).Where("$.BusinessImpact == 'Low'" ).ToArray();

			//CBLedInitiatives
			// konten += "<br /># of Primary Initiatives: "+primaryCount.length;
			// konten += "<br /># of Secondary Initiatives: "+SecondaryCount.length;
			konten += "<br /># of CB Led Initiatives: "+CBLITrueCount.length;
			konten += "<br /># of Bank Wide Initiatives: -";
			konten += "<br /># of High: " + high.length;
			konten += "<br /># of Middle: "+middle.length;
			konten += "<br /># of Low: "+low.length;
			konten += "<br /># of Task: "+total.length;
			if($('a[toltipsterid='+lc+']').hasClass("tooltipstered"))
				$('a[toltipsterid='+lc+']').tooltipster("destroy");
			$('a[toltipsterid='+lc+']').tooltipster({
				contentAsHTML: true,
				content: konten,
				trigger: 'hover'
			});
		})
	// })
}

c.SelectALLSC = function(){
	c.SelectedSC("");
	ResetCountryRegionFilter();
	c.ActiveBDFilter([])
	c.Filter.StartDate("");
	c.Filter.EndDate("");
	c.Filter.DisplayColor("")
	c.Clear('#search-input')
	c.GetData();
	
}
c.SetActiveBDFilter = function(obj){
	var d = ko.mapping.toJS(obj);
	// console.log('BDFilter', d)
	if(c.ActiveBDFilter().indexOf(d.Idx)>=0){
		c.ActiveBDFilter.remove(d.Idx);
	}else{
		c.ActiveBDFilter.push(d.Idx);
	}
	c.GetData();
}
c.GetBDClass = function(obj){
	var d = ko.mapping.toJS(obj);
	if(c.ActiveBDFilter().indexOf(d.Idx)>=0){
		return "active";
	}
	return "";
}
c.GetDefaultSCStyle = function(data){
	var BMLength = data.BusinessMetric().length;
	var height = ((BMLength>0?'height:'+(((BMLength*23)-BMLength) + 2)+'px;':''));
	if(SortInitiative.Active()){
		var tempData = Enumerable.From(c.DataSource().Data.TableSourcesVer3()).Where("$.Id === '"+data.Idx()+"'").FirstOrDefault();
		var InitiativePerLifeCycle = [];
		if(tempData!==undefined){
			// console.log(tempData);
			InitiativePerLifeCycle = ko.mapping.toJS(tempData.TableSourcesVer2());
		}
		var max = 0;
		for(var i in InitiativePerLifeCycle){
			var ilength = InitiativePerLifeCycle[i].Initiatives.length;
			if(ilength>max){
				max = ilength;
			}
		}
		var defaultHeight = (((BMLength*23)-BMLength) + 2);
		var iHeight = 50*max;
		if(iHeight > defaultHeight){
			return "height:"+iHeight+"px;";
		}else{
			return height;
		}
	}else{
		return height;
	}
}
c.GetSCStyleContent  = function(data,index){
	if(Scorecard.Data()[index]===undefined){
		return "";
	}
	var SCData = ko.mapping.toJS(Scorecard.Data()[index]);
	var BMLength = SCData.BusinessMetric.length;
	var height = ((BMLength>0?'height:'+(((BMLength*23)-BMLength) + 2)+'px;':''));
	var InitiativePerLifeCycle = ko.mapping.toJS(data.TableSourcesVer2);
	var max = 0;
	for(var i in InitiativePerLifeCycle){
		var ilength = InitiativePerLifeCycle[i].Initiatives.length;
		if(ilength>max){
			max = ilength;
		}
	}
	var defaultHeight = (((BMLength*23)-BMLength) + 2);
	var iHeight = 50*max;
	var res = defaultHeight;
	if(iHeight > defaultHeight){
		res = iHeight;
		height = "height:"+iHeight+"px;";
	}

	if(SortInitiative.Active()&&c.SelectedSC()!==""){
		var BMTotal =  Scorecard.BMTotal();
		var BDTotal = Scorecard.BDTotal()+2;
		var HeightTotal = BMTotal*23;
		var i = SCData.BusinessDriverList.length
		var heightValue = (i/BDTotal*HeightTotal);
		if(SCData.Idx === c.SelectedSC()){
		    if(i>6){
		    	i+=10;
		    }else if(i>=5){
		    	i+=7;
		    }else if(i>4){
		    	i+=5;
		    }else if(i>2){
		    	i+=4;
		    }else{
		    	i+=4;
		    }
			heightValue = (i/BDTotal*HeightTotal);
		}


		if(heightValue>res){
			res = heightValue;
			height = "height:"+heightValue+"px;";
		}
	}


	return height;
}
c.GetSCStyle = function(obj,isForSUBMenu){
	// console.log(obj);
	var BMTotal =  Scorecard.BMTotal();
	var BDTotal = Scorecard.BDTotal()+2;
	var HeightTotal = BMTotal*23;
	var d = ko.mapping.toJS(obj);
	var i = d.BusinessDriverList.length
	var heightValue = (i/BDTotal*HeightTotal);
	if(d.Idx === c.SelectedSC()){
	    if(i>6){
	    	i+=10;
	    }else if(i>=5){
	    	i+=7;
	    }else if(i>4){
	    	i+=5;
	    }else if(i>2){
	    	i+=4;
	    }else{
	    	i+=4;
	    }
		heightValue = (i/BDTotal*HeightTotal);
	}
	
	if(SortInitiative.Active()){
		var tempData = Enumerable.From(c.DataSource().Data.TableSourcesVer3()).Where("$.Id === '"+obj.Idx()+"'").FirstOrDefault()
		var InitiativePerLifeCycle = [];
		if(tempData!== undefined){
			InitiativePerLifeCycle = ko.mapping.toJS(tempData.TableSourcesVer2());
		}
		var max = 0;
		for(var i in InitiativePerLifeCycle){
			var ilength = InitiativePerLifeCycle[i].Initiatives.length;
			if(ilength>max){
				max = ilength;
			}
		}
		var iHeight = 50*max;
		if(iHeight>heightValue){
			heightValue = iHeight;
		}
		var BMLength = d.BusinessMetric.length;
		var defaultHeight = (((BMLength*23)-BMLength) + 2);
		if(defaultHeight>heightValue){
			heightValue = defaultHeight;
		}
	}
	var height = "height:"+heightValue+"px;";
	return height;
	// if(isForSUBMenu){
	// 	// c.SyncHeightForSortInitiativeSelectedSC.push(height)
	// 	return "height:"+heightValue+"px;";
	// }else{
	// 	// return "height:"+height+"px;padding-top:"+((height/2)-15)+"px;";
	// 	return "height:"+heightValue+"px;"
	// }
}
c.Get = function(obj){
	
	var d = ko.mapping.toJS(obj);
	// console.log('ob',d)
	c.SelectedSC(d.Idx);
	c.ActiveBDFilter([]);
	for(var i in d.BusinessDriverList){
		c.ActiveBDFilter.push(d.BusinessDriverList[i].Idx)
	}
	c.GetData();
	// setTimeout(function() {
 //      SortInitiative.sycHeight();
 //  }, 1000);
}

c.GetData = function(IsRefreshingInitiative){
	// console.log("Getting Data");
	c.Processing(true);
	var url = "/web-cb/dashboard/getpaneldata";
	var parm = ko.mapping.toJS(c.Filter);
	// if(c.SelectedTab()=="Initiative"){
	var countries = [];
	var regions = [];
    parm.RegionCountry.forEach(function(yy){
    	var isCountries = true;
    	var isGlobal = false;
	    c.RegionList().forEach(function(xx){
        	var region = xx._id
	        	if(yy == region){
	        		// regions.push(yy);
	        		isCountries = false;
	        	}
	        	if(yy == "GLOBAL"){
	        		isGlobal = true;
	        	}
    	});
    	if(isCountries){
    		countries.push(yy);
    	}else{
    		regions.push(yy);
    	}
    	if(isGlobal){
    		countries.push("GLOBAL")
    	}
    });
	parm.Country = countries;
	parm.Region = regions;
	// }else{
	// 	parm.Country = parm.Country == "Select.."? "" : parm.Country;
	// 	parm.Region = parm.Region == "Select.."? "" : parm.Region;
	// }
	parm.BDFilter = c.ActiveBDFilter();
	ajaxPost(url,parm,function(res){
		c.OwnedData(res.OwnedInitiative);
		res.MasterLifeCycle = c.LifeCycleList();
		c.MappingDataSource({Initiative:res});
		if(IsRefreshingInitiative){
			Initiative.Get(Initiative.OpenId(),true);
		}
		setTimeout(function(){

			c.Processing(false);
			redipsInit();
		}, 500);
	});


};
c.PullRequest.subscribe(function(val){
	if(val===0){
		c.GetData();
	}
})
c.Init = function(){
	c.PullRequest(2);
	ajaxPost("/web-cb/masterregion/getdata",{},function(res){
		c.CountryList(Enumerable.From(res.Data).GroupBy("$.Country").Select("{_id:$.Key()}").ToArray());
		c.RegionList(Enumerable.From(res.Data).GroupBy("$.Major_Region").Select("{_id:$.Key()}").ToArray());
		c.RegionalData(res.Data);
		c.PullRequest(c.PullRequest()-1)
	});
	ajaxPost("/web-cb/m/getlifecycledata",{},function(res){
		var arr = res.Data;
		for(var i in arr){
			arr[i].Id = arr[i].LifeCycleId;
		}
		c.LifeCycleList(arr);
		c.PullRequest(c.PullRequest()-1)
	});
}

c.shortname = function(n, len) {
  if(n.length <= len) {
      return n;
  }
  n = n.substr(0, len) + (n.length > len ? '...' : '');
  return n;
}

c.ExpandAll = function(){
	var iconAllExpand = c.iconAllExpand();
	// var iconAllExpandLain = !iconAllExpand;

	_.each(c.DataSource().Data.TableSources, function(v,i){
		c.DataSource().Data.TableSources[i].ShowHide(iconAllExpand);
    c.DataSource().Data.SummaryBusinessDriver[i].ShowHide(iconAllExpand);
    if(!iconAllExpand){
    	$('#caret'+v.Parentid()).removeClass('fa-caret-down');
      $('#caret'+v.Parentid()).addClass('fa-caret-right');
    } else{
      $('#caret'+v.Parentid()).removeClass('fa-caret-right');
      $('#caret'+v.Parentid()).addClass('fa-caret-down');
    }
	})

	if (iconAllExpand == true){
		c.iconAllExpand(false);
		$("#iconAllExpand").removeClass("fa fa-plus").addClass("fa fa-minus");
	}else{
		c.iconAllExpand(true);
		$("#iconAllExpand").removeClass("fa fa-minus").addClass("fa fa-plus");
	}
	c.SyncSCHeight()
}

c.SetSIWidth = function(){
	var div = $('<div>').addClass('col-sm-2').hide();
	$('body').append(div);
	var width = div.css('width');
	var parseWidth = parseInt(width);
	div.remove();
	var calcWidth = (parseWidth - 40) / 2 ;
	// document.getElementById("aScorecardTabMenu").style.width = calcWidth+"px";
	// document.getElementById("aInitiativeTabMenu").style.width = calcWidth+"px";
}

c.RagFlagDefined = function(data){
	var flagColor = (ko.toJS(data.DisplayProgress) == 'amber')? 'amber' : (ko.toJS(data.DisplayProgress) == 'green' || ko.toJS(data.DisplayProgress) == '')? 'green' : (ko.toJS(data.DisplayProgress) == undefined)? 'undefined' : 'red';
	return flagColor;
}

c.screenIsSmall = function(){
	return $(window).width();
}

$(document).ready(function(){


	c.GetUserList();
	c.Init();
	c.SyncHeightForSortInitiative = ko.observableArray([]);
	c.SyncHeightForSortInitiativeSelectedSC = ko.observableArray([]);
	$('[data-toggle="tooltip"]').tooltip();
	$(window).resize(function(){
		setTimeout(function() {c.LCWidth( $('.lcHeader').width() - 22 )}, 30);
		$('.overallprogress').width($('#btncollapseatas').width() + 60)
	});

	$("#gaugeTop").kendoRadialGauge({
      pointer: {
          value: 80
      },
      scale: {
          minorUnit: 5,
          startAngle: 0,
          endAngle: 180,
          max: 100,
          color:  "#fff",
          labels: {
              visible: false
              // position: "inside",
              // color : "#fff"
          },
          minorTicks: {
            size: 3
          },
          majorTicks : {
            size: 5
          } ,
          rangeSize :4,
          ranges: [
          {
            from: 0,
            to: 0,
            color: "#999999"
          },{
            from: 1,
            to: 30,
            color: "#ff0000"
          },{
            from: 30,
            to: 50,
            color: "#ff9900"
          },{
            from: 50,
            to: 60,
            color: "#ffff00"
          },{
            from: 60,
            to: 70,
            color: "#4dff4d"
          },{
            from: 70,
            to: 100,
            color: "#00b300"
          }
          ]
      }
  });

});
