var ov = {
	Loading: ko.observable(true),
	Initiatives:ko.observable(),
	Metrics:ko.observable(),
	Data:ko.observable(),
	ClientLifeCycle:ko.observable(),
  MaxClientLifeCycle: ko.observable(),
	RegionMetric:ko.observable(),
	MD:{
		Declined:ko.observable(),
		Unchanged:ko.observable(),
		Improved:ko.observable(),
	},
	MS:{
		Red:ko.observable(),
		Amber:ko.observable(),
		Green:ko.observable(),
	},
	IC:{
		BankWide:ko.observable(),
		CBLead:ko.observable(),
		Active:ko.observable(),
		Completed:ko.observable(),
	}
}

ov.GetData = function(){
	ajaxPost("/web-cb/dashboard/getdata",{},function(res){
    console.log(res)
		if(res.IsError){
			swal("", res.Message, "info");
            return false;
		}
		console.log(res.Data)
		// ov.Data(res.Data);
		ov.Mapping(res.Data);
		setTimeout(function(){
			ov.renderchart1();
			ov.renderchart2();
			ov.renderchart3();
			ov.renderchart4();
			ov.renderchart5();
		}, 500);
		// ov.Render(res.Data);
	})
}

ov.GetBarStyle = function(data){
  var maxheight = 180;
  var satuan = maxheight / (ov.MaxClientLifeCycle() + 1)
  height = 'height: 0px;'
  if (data != 0){
    tinggi = data * satuan;
    height = 'height: ' + tinggi + 'px;';
    if(tinggi > 20){
      height += 'padding-top: ' + tinggi * 0.4 + 'px;';
    }
  }
  // console.log(height)
  return height;
}

ov.Mapping =function(dataSource){
	ov.Data(dataSource);
	ov.Initiatives(dataSource.Initiatives)
	ov.Metrics(dataSource.Metrics)
	ov.ClientLifeCycle(dataSource.ClientLifeCycle)

  tmpMaxCLC = 0;
  for(i in dataSource.ClientLifeCycle){
    tmpTotal = dataSource.ClientLifeCycle[i].BankWide + dataSource.ClientLifeCycle[i].CBLead;
    if (tmpMaxCLC < tmpTotal){
      tmpMaxCLC = tmpTotal;
    }
  }
  ov.MaxClientLifeCycle(tmpMaxCLC);

	// ### Metric Direction
	ov.MD.Declined(dataSource.MetricDirection.Declined)
	ov.MD.Improved(dataSource.MetricDirection.Improved)
	ov.MD.Unchanged(dataSource.MetricDirection.Unchanged)
	// ### Metric Staus
	ov.MS.Amber(dataSource.MetricStatus.Amber)
	ov.MS.Green(dataSource.MetricStatus.Green)
	ov.MS.Red(dataSource.MetricStatus.Red)
	// ### Initiative Comparison
	ov.IC.Active(dataSource.InitiativeComparison.Active)
	ov.IC.Completed(dataSource.InitiativeComparison.Completed)
	ov.IC.BankWide(dataSource.InitiativeComparison.BankWide)
	ov.IC.CBLead(dataSource.InitiativeComparison.CBLed)
	//### Count Majour Reogoion
	ov.RegionMetric(dataSource.NumberOfMetricRegion)
	// ov.CountRegion(dataSource.NumberOfMetric)
}

ov.CountRegion = function(source){
	var regionMetric = []
	var regionList = Enumerable.From(source).Select("$.MajorRegion").Distinct().ToArray();
	console.log('region',regionList)
	c.RegionList().forEach(function(xx){
		var region = xx._id
		if(region === "GLOBAL"){
			return
		}
		var filteredRegion = Enumerable.From(source).Where("$.MajorRegion == '"+region+"'").ToArray();
		var green = Enumerable.From(filteredRegion).Sum("$.GREEN");
		var amber = Enumerable.From(filteredRegion).Sum("$.AMBER");
		var red = Enumerable.From(filteredRegion).Sum("$.RED");
		regionMetric.push({"MajorRegion":region, "GREEN":green, "AMBER":amber, "RED":red})
	})
	ov.RegionMetric(regionMetric)
}

ov.Render = function(dataSource){
	// Data Format - it will came from ov.GetData
	// This sample is for top 2 chart,  and donut 
	dataSource = {
		Initiatives:106,
		Metrics:30,
		Summary:[
			{Name:"Financial Framework",Initiatives:7,Metrics:7},
			{Name:"Client,Growth & Collaboration ",Initiatives:18,Metrics:6}
			// dst.. base on Scorecard Category that we have
		],
		InitiativeComparison:{
			BankWide:63,
			CBLed:43,
			Active:63,
			Completed:43,
		},
		MetricDirection:{
			Declined:4,
			Unchanged:15,
			Improved:11,
		},
		MetricStatus:{
			Red:7,
			Amber:13,
			Green:10,
		},
		ClientLifeCycle:[
			{LifeCycle:"Prospecting & On Board",CBLead:17,BankWide:10},
			{LifeCycle:"Credit",CBLead:4,BankWide:6},
			// and so on
		],
		NumberOfMetric:[
			{
				Country:"BOTSWANA",
				MajorRegion:"AME",
				AMBER:13,
				GREEN:15,
				RED:2,
			},
			{
				Country:"SINGAPORE",
				MajorRegion:"ASA",
				AMBER:13,
				GREEN:15,
				RED:2,
			},
			{
				Country:"TAIWAN",
				MajorRegion:"GCNA",
				AMBER:13,
				GREEN:15,
				RED:2,
			},
			// and so on....
		]
	}
}

ov.renderchart1 = function(){
	var dataSource = ov.Data().Summary
	var barChart1Data = Enumerable.From(dataSource).GroupBy("$.Country").Select("{_id:$.Key()}").ToArray()
	$("#barchart-1").kendoChart({
		dataSource:dataSource,
		chartArea: {
			height: 180,
			width: 360,
			margin: {
				top:25,
			  right: 25
			}
		},
		seriesDefaults: {
			type: "bar",
			overlay: {
				gradient: "none"
			},
			labels: {
				font: "10px Helvetica Neue, Helvetica, Arial, sans-serif",
				// color: "#fff",
				visible: true,
				position:"outsideEnd",
				background: "transparent",
				margin: {
				        // top: -20,
				        left:-1, 
		      	}
			},
			gap: 3,
			color: colorcustom,
			border: {
				color: bordercustom,
			},
		},
	   series:[{
	      // data: barChart1Data,
	      field: "Initiatives"
	   }],
	   categoryAxis :{
	   	// categories: ["Financial Framework", "Clients, Growth, & Collaboration", "Digitization & Analytics", "Risk & Control", "Efficiency, Productivity, & Service Quality", "People, Culture, & Conduct"],
	   	field: "Name",
			labels: {
				font: "10px Helvetica Neue, Helvetica, Arial, sans-serif",
				// color: "#fff",
				visible: true,
				mirror: true,
				// position:"insie",
				background: "transparent",
				padding: {
			        left: -5
			    },
				margin: {
				        top: -25,
		      	},
            border:{
              color:"transparent"
            },
			},
			majorGridLines: {
				visible: false
			},
			line:{
				visible:false
			}
		},
	   valueAxis:{
			max: _.max(dataSource.Initiatives),
			majorGridLines: {
				visible: false,
			},
			line: {
				visible: false
			},
			labels:{
				visible:false,
			},
		},
	});

	setTimeout(function () {
		$('#barchart-1 g:first g:first > g:eq(9) text').attr('fill', '#0e4361')
	}, 200)

	function colorcustom(e){
		if (e.index == 0){
			return "#0e4361"
		}else if (e.index == 1){
			return "#126390"
		}else if (e.index == 2){
			return "#1a82bf"
		}else if (e.index == 3){
			return "#479acc"
		}else if (e.index == 4){
			return "#50D0FF"
		}else if (e.index == 5){
			return "#9bd9ff"
		};
	}
	function bordercustom(e){
		if (e.index == 0){
			return "#0e4361"
		}else if (e.index == 1){
			return "#126390"
		}else if (e.index == 2){
			return "#1a82bf"
		}else if (e.index == 3){
			return "#479acc"
		}else if (e.index == 4){
			return "#50D0FF"
		}else if (e.index == 5){
			return "#9bd9ff"
		};
	}
}

ov.renderchart2 = function(){
	var dataSource = ov.Data().Summary
	$("#barchart-2").kendoChart({
		dataSource:dataSource,
		chartArea: {
			height: 180,
			width: 360,
			margin: {
				top:25
			}
		},
		seriesDefaults: {
			type: "bar",
			overlay: {
				gradient: "none"
			},
			labels: {
				font: "10px Helvetica Neue, Helvetica, Arial, sans-serif",
				// color: "#fff",
				visible: true,
				position:"outsideEnd",
				background: "transparent",
				margin: {
				        // top: -20,
				        left:-1, 
		      	}
			},
			gap: 3,
			color: colorcustom,
			border: {
				color: bordercustom,
			},
		},
	   series:[{
	      field:"Metrics"
	   }],
	   categoryAxis :{
			// categories: ["Financial Framework", "Clients, Growth, & Collaboration", "Digitization & Analytics", "Risk & Control", "Efficiency, Productivity, & Service Quality", "People, Culture, & Conduct"],
			field:"Name",
			labels: {
				font: "10px Helvetica Neue, Helvetica, Arial, sans-serif",
				// color: "#fff",
				visible: true,
				mirror: true,
				background: "transparent",
				padding: {
	        		left: -5
	      		},
				margin: {
		        	top: -25,
		      	},
            border:{
              color:"transparent"
            },
			},
			majorGridLines: {
				visible: false
			},
			line:{
				visible:false
			}
		},
	   valueAxis:{
			max: _.max(dataSource.Metrics),
			majorGridLines: {
				visible: false,
			},
			line: {
				visible: false
			},
			labels:{
				visible:false,
			},
		},
	});

	setTimeout(function () {
		$('#barchart-2 g:first g:first > g:eq(9) text').attr('fill', '#0e4361')
	}, 200)

	function colorcustom(e){
		if (e.index == 0){
			return "#0e4361"
		}else if (e.index == 1){
			return "#126390"
		}else if (e.index == 2){
			return "#1a82bf"
		}else if (e.index == 3){
			return "#479acc"
		}else if (e.index == 4){
			return "#50D0FF"
		}else if (e.index == 5){
			return "#9bd9ff"
		};
	}
	function bordercustom(e){
		if (e.index == 0){
			return "#0e4361"
		}else if (e.index == 1){
			return "#126390"
		}else if (e.index == 2){
			return "#1a82bf"
		}else if (e.index == 3){
			return "#479acc"
		}else if (e.index == 4){
			return "#50D0FF"
		}else if (e.index == 5){
			return "#9bd9ff"
		};
	}
}

ov.renderchart3 = function(){
	$("#donutchart-1").kendoChart({
    legend: {

    	visible:false,
      position: "right",
 	      margin:{
					visible:true,
					left:-10,
					top:-60
				},

      },
      chartArea: {
			height : 160,
			width : 160,
			background: "transparent",
			margin:{
				visible:true,
				// left:-10,
				top:-20
			},
		},
    seriesDefaults: {
      labels: {
        template: "#= kendo.format('{0:P0}', percentage)#",
        position: "center",
        visible: true,
        background: "transparent",
        // padding: 7,
        color: "#fff",
		    font: "9px Helvetica Neue, Helvetica Neue, Helvetica, Arial, sans-serif",
        border:{
          color:"transparent"
        }
      }
    },
    series: [{
      type: "donut",
      // holeSize: 10,
      size:30,
      data: [{
        category: "Bank Wide",
        value: ov.IC.BankWide()
      },{
        category: "CB Led",
        value: ov.IC.CBLead()
      }],
      overlay:{
       	gradient: "none"
      }
    }],
    seriesColors:["#0e4361","#50D0FF"],
    valueAxis:{
			visible:false,
			labels: {
				font: "10px Helvetica Neue, Helvetica, Arial, sans-serif",
				visible: false,
			},
			majorGridLines: {
				visible: false
			},
		},
  })
  $("#donutchart-2").kendoChart({
    legend: {
    	visible:false,
      position: "right",
	      margin:{
				visible:true,
				left:-10,
				top:-60
			},
    },
    chartArea: {
			height : 160,
			width : 160,
			background: "transparent",
			margin:{
				visible:true,
				// left:-10,
				top:-20
			},
		},
    seriesDefaults: {
      labels: {
        template: "#= kendo.format('{0:P0}', percentage)#",
        position: "center",
        visible: true,
        background: "transparent",
        // padding: 7,
        color: "#fff",
				font: "9px Helvetica Neue, Helvetica Neue, Helvetica, Arial, sans-serif",
        border:{
          color:"transparent"
        }
      }
    },
    series: [{
      type: "donut",
      // holeSize: 30,
      size:30,
      data: [{
        category: "Active",
        value: ov.IC.Active()
      },{
        category: "Completed",
        value: ov.IC.Completed()
      }],
      overlay:{
        gradient: "none"
      }
    }],
    seriesColors:["#0e4361","#50D0FF"],
    valueAxis:{
			visible:false,
			labels: {
				font: "10px Helvetica Neue, Helvetica, Arial, sans-serif",
				visible: false,
			},
			majorGridLines: {
				visible: false
			},
		},
  })
}

ov.renderchart4 = function(){

	var dataSource = ov.ClientLifeCycle()
	
	var scwidth = $("#stackchart-1-width").width();
	$("#stackchart-1").kendoChart({
		dataSource:dataSource,
		title: {
		  text: "Client Life Cycle",
			font: "14px Helvetica Neue, Helvetica, Arial, sans-serif",
			color: "#0e4361",
			margin:{
		   		left:-90,
		   	},
		},
		chartArea: {
	   	height: 210,
	   	width: scwidth,
		   	padding:0,
		   	margin:{
		   		right:-90,
		   	},
		},
    legend: {
    	visible:false,
	    // position: "custom",
	    // orientation: "vertical",
	    // offsetX: 500,
	    // width: 200,
	  },
    seriesDefaults: {
      stack: true,
      overlay: {
				gradient: "none"
			},
			gap: 0.8,
			labels: {
				font: "11px Helvetica Neue, Helvetica, Arial, sans-serif",
				color: "#fff",
				visible: true,
				position:"center",
				background: "transparent"
			},
    },
    series: [
	    { name: "CB Led", field: "CBLead"},
	    { name: "Bank Wide", field:"BankWide"},
	  ],
	  seriesColors: ["#50D0FF", "#0e4361"],
    valueAxis: {
      // max: 19,
      title: {
	      text: "Number of Initiatives",
				font: "12px Helvetica Neue, Helvetica, Arial, sans-serif",
	      color: "#0e4361",
	      margin:{
	      	right:-40
	      }
	    },
      labels: {
	      visible: false
	    },
      line: {
        visible: false
      },
      majorGridLines: {
        visible: false
      }
    },
    categoryAxis: {
      // categories: ["Prospecting & On Boarding", "Credit", "Documentation", "Account Opening", "RM Client and Account Management", "Operations Account Servicing", "Monitoring & Off Boarding"],
      // labels: {
      //   visual: function(e) {

      //     var html = $('<div style="width:70px;background:#0e4061;text-align:center;font: 10px Helvetica Neue, Helvetica, Arial, sans-serif"><b>' + e.text + '</b></div>').appendTo(document.body);
      //     var visual = new kendo.drawing.Group();
      //     var rect = e.rect;
      //     kendo.drawing.drawDOM(html).done(function(group) {
      //       html.remove();
      //       var layout = new kendo.drawing.Layout(rect, {
      //         justifyContent: "center"
      //       });
      //       layout.append(group);
      //       layout.reflow();
      //       visual.append(layout);
      //     });
      //     return visual;
      //   }
      // },
      line: {
        visible: false
      },
      majorGridLines: {
        visible: false
      }
    },
	});

	// setTimeout(function () {
	 //    var left = $('#stackchart-1 [clip-path]').offset().left - 65;
		// var width = $('#stackchart-1 [clip-path]:eq(0)')[0].getBoundingClientRect().width + 30;
		// console.log('left',left)
		// console.log('width',width)
		// $('.stackchartlegend-1').closest('table').css('margin-left', left + 'px')
		// $('.stackchartlegend-1').closest('table').attr('width', width)
	// }, 200)
}

ov.top5DataList = ko.observableArray([])

ov.renderchart5 = function(){
	$("#stackchart-2").kendoChart({
		dataSource: ov.RegionMetric(),
		title: {
		  text: "Region",
			font: "14px Helvetica Neue, Helvetica, Arial, sans-serif",
			color: "#0e4361"
		},
		chartArea: {
	   	height: 210,
	   	width: 225,
	   	margin:{
	   		// left:20
	   	}
		},
		legend: {
    	visible:false,
	    position: "custom",
	    orientation: "vertical",
	    width: 200
	  },
    seriesDefaults: {
      stack: true,
      overlay: {
				gradient: "none"
			},
			gap: 0.5,
			labels: {
				font: "10px Helvetica Neue, Helvetica, Arial, sans-serif",
				color: "#fff",
				visible: true,
				position:"center",
				background: "transparent"
			},
    },
    series: [
	    { name: "GREEN",field:"GREEN"},
	    { name: "AMBER" ,field:"AMBER"},
		{ name: "RED",field:"RED"},
    ],
	  seriesColors: ["#36bc9b" , "#f7bb3b", "#db4453"],
    valueAxis: {
      // max: 30,
      title: {
	      text: "Number of Metrics",
				font: "12px Helvetica Neue, Helvetica, Arial, sans-serif",
	      color: "#0e4361"
	    },
      labels: {
	      visible: false
	    },
      line: {
        visible: false
      },
      majorGridLines: {
        visible: false
      }
    },
    categoryAxis: {
      // categories: ["AME", "ASA", "GCNA"],
      labels: {
        visual: function(e) {

          var html = $('<div style="width:50px;background:#0e4061;text-align:center;font: 10px Helvetica Neue, Helvetica, Arial, sans-serif"><b>' + e.text + '</b></div>').appendTo(document.body);
          var visual = new kendo.drawing.Group();
          var rect = e.rect;
          kendo.drawing.drawDOM(html).done(function(group) {
            html.remove();
            var layout = new kendo.drawing.Layout(rect, {
              justifyContent: "center"
            });
            layout.append(group);
            layout.reflow();
            visual.append(layout);
          });
          return visual;
        }
      },
      line: {
        visible: false
      },
      majorGridLines: {
        visible: false
      }
    },
	});



	// var data = "BOTSWANA;10;9;11|GHANA;4;22;4|KENYA;17;1;12|NIGERIA;12;10;8|TANZANIA UNITED REPUBLIC OF;5;16;9|UGANDA;9;14;7|ZAMBIA;10;13;7|ZIMBABWE;9;13;8|BAHRAIN;10;17;3|JORDAN;7;14;9|OMAN;6;12;12|PAKISTAN;12;4;14|QATAR;5;16;9|UNITED ARAB EMIRATES;16;1;13|INDONESIA;13;14;3|MALAYSIA;13;6;11|SINGAPORE;15;1;14|THAILAND;14;12;4|VIETNAM;6;21;3|BANGLADESH;16;8;6|INDIA;2;13;15|NEPAL;8;17;5|SRI LANKA;5;13;12|CHINA;4;16;10|HONG KONG;15;14;1|TAIWAN;18;2;10|KOREA, REPUBLIC OF;13;8;9"
	// data = data.split("|");
	// var dataList = [];
	// for(var i in data){
	// 	var arr = data[i].split(";");
	// 	var d = {Country:arr[0],RED:parseInt(arr[1]),AMBER:parseInt(arr[2]),GREEN:parseInt(arr[3])};
	// 	dataList.push(d)
	// }
	var data = ov.Data().NumberOfMetric
	var dataList = Enumerable.From(data).OrderByDescending("$.GREEN").Take(5).ToArray();
	var MaxValue = Enumerable.From(dataList).Max("$.RED+$.AMBER+$.GREEN")+10;
	ov.top5DataList(dataList);

	$("#stackchart-3").kendoChart({
		dataSource:{
			data:dataList
		},
		title: {
		  text: "Top 5 Countries (by # of Green)",
			font: "14px Helvetica Neue, Helvetica, Arial, sans-serif",
			color: "#0e4361"
		},
		chartArea: {
	   	height: 210,
	   	width: 300,
		},
    legend: {
      visible: false
    },
    seriesDefaults: {
      stack: true,
      overlay: {
				gradient: "none"
			},
			gap: 0.3,
			labels: {
				font: "10px Helvetica Neue, Helvetica, Arial, sans-serif",
				color: "#fff",
				visible: true,
				position:"center",
				background: "transparent"
			},
    },
    series: [
	    { name: "GREEN",field:"GREEN"},
	    { name: "AMBER" ,field:"AMBER"},
		{ name: "RED",field:"RED"},
    ],
	  seriesColors: ["#36bc9b" , "#f7bb3b", "#db4453"],

    valueAxis: {
      max: MaxValue,
      title: {
	      text: "Number of Metrics",
				font: "12px Helvetica Neue, Helvetica, Arial, sans-serif",
	      color: "#0e4361"
	    },
      labels: {
	      visible: false
	    },
      line: {
        visible: false
      },
      majorGridLines: {
        visible: false
      }
    },
    categoryAxis: {
    	field:"Country",
      // categories: ["Thailand", "Singapore", "Kenya", "Botswana", "India"],
      labels: {
      	visible:false,
        visual: function(e) {
          var html = $('<div style="width:50px;background:#0e4061;text-align:center;font: 10px Helvetica Neue, Helvetica, Arial, sans-serif"><b>' + e.text + '</b></div>').appendTo(document.body);
          var visual = new kendo.drawing.Group();
          var rect = e.rect;
          kendo.drawing.drawDOM(html).done(function(group) {
            html.remove();
            var layout = new kendo.drawing.Layout(rect, {
              justifyContent: "center"
            });
            layout.append(group);
            layout.reflow();
            visual.append(layout);
          });
          return visual;
        }
      },
      line: {
        visible: false
      },
      majorGridLines: {
        visible: false
      }
    },
	});
}
ov.DataList = ko.observableArray([]);
ov.renderchart6 = function(){
	// var data = "BOTSWANA;10;9;11|GHANA;4;22;4|KENYA;17;1;12|NIGERIA;12;10;8|TANZANIA UNITED REPUBLIC OF;5;16;9|UGANDA;9;14;7|ZAMBIA;10;13;7|ZIMBABWE;9;13;8|BAHRAIN;10;17;3|JORDAN;7;14;9|OMAN;6;12;12|PAKISTAN;12;4;14|QATAR;5;16;9|UNITED ARAB EMIRATES;16;1;13|INDONESIA;13;14;3|MALAYSIA;13;6;11|SINGAPORE;15;1;14|THAILAND;14;12;4|VIETNAM;6;21;3|BANGLADESH;16;8;6|INDIA;2;13;15|NEPAL;8;17;5|SRI LANKA;5;13;12|CHINA;4;16;10|HONG KONG;15;14;1|TAIWAN;18;2;10|KOREA, REPUBLIC OF;13;8;9"
	var data = ov.Data().NumberOfMetric;
	var dataList = [];
	for(var i in data){
		dataList[i]=data[i]
		if(data[i].Country =='TANZANIA UNITED REPUBLIC OF'){
			dataList[i].Country = 'TANZANIA'
		}else if(data[i].Country =='UNITED ARAB EMIRATES'){
			dataList[i].Country = 'UAE'
		}else if(dataList[i].Country == 'KOREA, REPUBLIC OF'){
			dataList[i].Country = 'KOREA'
		}
		
	}
	// var dataList = ov.Data().NumberOfMetric;
	var MaxValue = Enumerable.From(dataList).Max("$.RED+$.AMBER+$.GREEN")+10;
	// ov.DataList(dataList);
	// var seriesList = [
	// 	{ name: "RED", data: [] },
	//     { name: "AMBER", data: [] },
	//     { name: "GREEN", data: [] },
	// ]
	// for(var i in dataList){
	// 	seriesList[0].data.push(dataList[i].RED)
	// 	seriesList[0].data.push(dataList[i].RED)
	// 	seriesList[0].data.push(dataList[i].RED)
	// 	// seriesList[s].data.push(dataList[i].RED)
	// }
	// console.log(seriesList);

	$("#stackchart-4").kendoChart({
        dataSource: {
            data:dataList
        },
		chartArea: {
	   	height: 346,
	   	width: 1150,
	   	margin:{
				visible:true,
				top:-60
			},
		},
    legend: {
    	visible:false,
	    position: "custom",
	    orientation: "vertical",
	    width: 200
	  },
    seriesDefaults: {
      stack: true,
      overlay: {
				gradient: "none"
			},
			gap: 0.1,
			labels: {
				font: "11px Helvetica Neue, Helvetica, Arial, sans-serif",
				color: "#fff",
				visible: true,
				position:"center",
				background: "transparent"
			},
    },
    series: [

	    { name: "GREEN",field:"GREEN"},
	    { name: "AMBER" ,field:"AMBER"},
		{ name: "RED",field:"RED"},
	    
    ],
	  seriesColors: ["#36bc9b" , "#f7bb3b", "#db4453"],
    valueAxis: {
      max: MaxValue,
      title: {
	      text: "Number of Metrics",
				font: "16px Helvetica Neue, Helvetica, Arial, sans-serif",
	      color: "#0e4361"
	    },
      labels: {
	      visible: false
	    },
      line: {
        visible: false
      },
      majorGridLines: {
        visible: false
      }
    },
    categoryAxis: {
    	field: "Country",
      // categories: ["Thailand", "Singapore", "Kenya", "Botswana", "India", "Thailand", "Singapore", "Kenya", "Botswana", "India", "Thailand", "Singapore", "Kenya", "Botswana", "India", "Thailand", "Singapore", "Kenya", "Botswana", "India", "Thailand", "Singapore", "Kenya", "Botswana", "India"],
      labels: {
      	visible:true,
      	font:"8px Helvetica Neue",
      	rotation:20,
        visual: function(e) {
          // var html = $('<div style="width:50px;background:#0e4061;text-align:center;font: 8px Helvetica Neue"><b>' + e.text + '</b></div>').appendTo(document.body);
          // var visual = new kendo.drawing.Group();
          // var rect = e.rect;
          // kendo.drawing.drawDOM(html).done(function(group) {
          //   html.remove();
          //   var layout = new kendo.drawing.Layout(rect, {
          //     justifyContent: "center"
          //   });
          //   layout.append(group);
          //   layout.reflow();
          //   visual.append(layout);
          // });
          // return visual;
        }
      },
      line: {
        visible: true
      },
      majorGridLines: {
        visible: false
      }
    },
	});
	function shortLabels(value) {
   if (value.length > 3) {
      value = value.substring(0, 5);
      return value;
   }
}

}
ov.ExportAsPDF2 = function () {
  	$('.btn-export-pdf').hide()
    $('.btn-export-png').hide()
  	$('.popupcustom').hide()
  	$('.ovheader').show()

    kendo.pdf.defineFont({
        "Helvetica Neue": "/web-cb/static/fonts/HelveticaNeue.ttf"
    });

    kendo.drawing.drawDOM($('#app-content > .col-md-12')).then(function (group) {
        var title = "Dashboard.pdf"
        kendo.drawing.pdf.saveAs(group, title);
        $('.btn-export-png').show()
    		$('.btn-export-pdf').show()
        $('.popupcustom').show()
    		$('.ovheader').hide()
    })
}

ov.ExportAsPDF = function () {
  $('.btn-export-pdf').hide();
  $('.btn-export-png').hide();
  $('.popupcustom').hide();
  $('.ovheader').show()
  $('#app-content > .col-md-12').css('background', 'white').css('padding','15px');

  kendo.pdf.defineFont({
      "Helvetica Neue": "/web-cb/static/fonts/HelveticaNeue.ttf"
  });

  // $("g text").attr("fill","#000").parent().attr("opacity",1).find("path").hide()

    // kendo.pdf.defineFont({
    //     "Open Sans": "/static/fonts/OpenSans-Regular.ttf",
    //     "Helvetica Neue": "/static/fonts/HelveticaNeue.ttf"
    // });

    // kendo.drawing.drawDOM($('#app-content > .col-md-12')).then(function (group) {
    //     var title = "Dashboard.png"
    //     kendo.drawing.pdf.saveAs(group, title);
    // });

    kendo.drawing.drawDOM($('#app-content > .col-md-12')) //#qrMail
    .then(function(group) {
      // Render the result as a PNG image
      return kendo.drawing.exportPDF(group);
    })
    .done(function(data) {
      // Save the image file
      kendo.saveAs({
        dataURI: data,
        fileName: "Dashboard.pdf",
        // proxyURL: "http://demos.telerik.com/kendo-ui/service/export"
      });
      $('.ovheader').hide();
      $('.btn-export-pdf').show();
      $('.btn-export-png').show();
      $('.popupcustom').show()
      $('#app-content > .col-md-12').css('background', 'white').css('padding','0');
    });
}

ov.ExportAsPNG = function () {
  $('.btn-export-pdf').hide();
  $('.btn-export-png').hide();
  $('.ovheader').show()
  $('#app-content > .col-md-12').css('background', 'white').css('padding','15px');
  // $("g text").attr("fill","#000").parent().attr("opacity",1).find("path").hide()

    // kendo.pdf.defineFont({
    //     "Open Sans": "/static/fonts/OpenSans-Regular.ttf",
    //     "Helvetica Neue": "/static/fonts/HelveticaNeue.ttf"
    // });

    // kendo.drawing.drawDOM($('#app-content > .col-md-12')).then(function (group) {
    //     var title = "Dashboard.png"
    //     kendo.drawing.pdf.saveAs(group, title);
    // });

    kendo.drawing.drawDOM($('#app-content > .col-md-12')) //#qrMail
    .then(function(group) {
      // Render the result as a PNG image
      return kendo.drawing.exportImage(group);
    })
    .done(function(data) {
      // Save the image file
      kendo.saveAs({
        dataURI: data,
        fileName: "Dashboard.png",
        // proxyURL: "http://demos.telerik.com/kendo-ui/service/export"
      });
      $('.ovheader').hide();
      $('.btn-export-pdf').show();
      $('.btn-export-png').show();
      $('#app-content > .col-md-12').css('background', 'white').css('padding','0');
    });
    
    // console.log('hello')
}

ov.Loading.subscribe(function(az){
	// #### calculate stackchart1 legend
	if(!az){
		setTimeout(function(){
		    // var left = $('#stackchart-1 [clip-path]').offset().left - 45;
			// var width = $('#stackchart-1 [clip-path]:eq(0)')[0].getBoundingClientRect().width + 30;
			// console.log('left',left)
			// console.log('width',width)
			// $('.stackchartlegend-1').closest('table').css('margin-left', left+'px')
			// $('.stackchartlegend-1').closest('table').attr('width', width)
		},500);
	}

});

$(function(){
	ov.GetData();
	ov.Loading(true);

	setTimeout( function(){	
		ov.Loading(false);
	},500);

});