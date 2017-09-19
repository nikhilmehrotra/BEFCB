var ScorecardInitiative = {
    Processing:ko.observable(false),
    ScorecardColor:["#A1AFC2"],
    Data:ko.observableArray([]),
    DataGraphic: ko.observable(),
	DetailData:ko.observableArray(),
	Name: ko.observable(''),
    Max:ko.observable(0),
    seeDirect: ko.observable(true),
    seeIndirect: ko.observable(true)
}
ScorecardInitiative.seeDirectIndirect = function(status){
    if(status == 1){
        //direct
        a = !ScorecardInitiative.seeDirect();
        ScorecardInitiative.seeDirect(a);
    } else if(status == 2){
        //indirect
        a = !ScorecardInitiative.seeIndirect();
        ScorecardInitiative.seeIndirect(a);
    }
}
ScorecardInitiative.Get = function(){
	ScorecardInitiative.GetData();	
}
ScorecardInitiative.GetSCDetail = function(index){
    return false;
    var Businessmetrics = ScorecardInitiative.Data()[index].Businessmetrics;
    var initiative = [];
    for(var i in Businessmetrics){
        initiative = initiative.concat(Businessmetrics[i].Initiative);
    }
    ScorecardInitiative.GetDetail(initiative,undefined,data);
}
ScorecardInitiative.GetKeyDetail = function(data){
    ScorecardInitiative.GetDetail(data.Initiative,undefined,data);
}
ScorecardInitiative.GetBarStyle = function(val){
    var style = "width:"+val+"%;";
    if(val==0){
        style += "display:none;"
    }

    return style;
}
ScorecardInitiative.GetRAGColor = function(val) {
    var result = "";
    result += "";
    switch (val) {
        case "red":
            result += "#f74e4e";
            break;
        case "amber":
            result += "#ffd24d";
            break;
        case "green":
            result += "#6ac17b";
            break;
        default:
            result += "#F2F7FC";
            break;
    }
    return result;
}

ScorecardInitiative.GetDirectDetail = function(data){
    ScorecardInitiative.GetDetail(data.Initiative,1,data);
}
ScorecardInitiative.GetIndirectDetail = function(data){
    ScorecardInitiative.GetDetail(data.Initiative,2,data);
}

ScorecardInitiative.GetData = function(){
    ScorecardInitiative.Processing(true);
    var url = "/web-cb/scorecardinitiative/getdata";
    var parm = {
        Region:c.Filter.RegionOne(),
        Country:c.Filter.CountryOne(),
    };
    ajaxPost(url,parm,function(res){
        var dataSource = res.Data;
        var maxinitiative = Enumerable.From(dataSource).SelectMany("$.Businessmetrics").Max("$.Initiative.length")+1;
        for(var i in dataSource){
            for(var b in dataSource[i].Businessmetrics){
                // Initiatieves
                var max = dataSource[i].Businessmetrics[b].Initiative.length;
                var maxPercentage = maxinitiative > 0 ? (dataSource[i].Businessmetrics[b].Initiative.length/maxinitiative*100) : 0;
                dataSource[i].Businessmetrics[b].DirectPercentage = max > 0 ? (dataSource[i].Businessmetrics[b].Direct/max*maxPercentage) : 0;
                dataSource[i].Businessmetrics[b].IndirectPercentage = max > 0 ? (dataSource[i].Businessmetrics[b].Indirect/max*maxPercentage) : 0;
                

                // RAG
                var RAGPeriod = CurrentMonth;
                var RAGTotal = dataSource[i].Businessmetrics[b].RAG.length;
                var RAGDisplay = 3;
                if(RAGTotal>0){
                    RAGPeriod = jsonDate(Enumerable.From(dataSource[i].Businessmetrics[b].RAG).OrderBy(function(x){
                        return jsonDate(x);
                    }).FirstOrDefault().Period);
                }
                if(RAGDisplay != RAGTotal){
                    var RAGRemaining = RAGDisplay-RAGTotal;
                    for(var r=0;r<RAGRemaining;r++){
                        RAGPeriod.setMonth(RAGPeriod.getMonth() - 1);
                        dataSource[i].Businessmetrics[b].RAG.unshift({
                            BmId:dataSource[i].Businessmetrics[b].Id,
                            BusinessMetric:dataSource[i].Businessmetrics[b].Description,
                            Period:RAGPeriod.toJSON(),
                            Rag:""
                        });
                    }
                }

                // Quarter Mapping
                for(var q in dataSource[i].Businessmetrics[b]){
                    
                }

            }
        }
        ScorecardInitiative.Processing(false);
        ScorecardInitiative.createDataGraphic(maxinitiative);
        ScorecardInitiative.Data(dataSource);

        setTimeout(function() {
            ScorecardInitiative.SyncHeightRAG()
        }, 200);
    });
}

ScorecardInitiative.createDataGraphic = function(Max){
	MaxPercentage = ( (Max-1)/Max ) * 100;
	SpacePercentage = ( 1/Max ) * 100;
	RangePercentage = MaxPercentage / (Max - 1);
	ArrayRangePercentage = [];

	for (i = 1; i < Max; i++) { 
		MarginLeft = (i == 1) ? RangePercentage/2 : 0;
   	ArrayRangePercentage.push({
   		RangePercentage: RangePercentage,
   		MarginLeft: MarginLeft,
   		Number: i,
   	}) 
	}

	ScorecardInitiative.DataGraphic({
		a: MaxPercentage,
		b: SpacePercentage,
		c: ArrayRangePercentage,
		// d: Max,
	});
}

ScorecardInitiative.SyncHeightRAG = function(){
    _.each($('.nopadding'),function(v,i){
      tmp = $(v).parent().height();
      $(v).children().height(tmp)
    })
}

ScorecardInitiative.GetDetail = function(data,type,bmdata){	
	if(type!==undefined){
		data = _.filter(data, function(e){return e.type == type})
	}
    _.each(data, function(v){
        v.SCName = bmdata.SCName;
        switch(v.type){
            case 0:
                v.DirectorIndirect = "N/A";
                break;
            case 1:
                v.DirectorIndirect = "Direct";
                break;
            case 2:
                v.DirectorIndirect = "Indirect";
                break;
            default:break;
        }
        v.MetricName = bmdata.Description;
        v.MetricColor = ScorecardInitiative.GetRAGColor(bmdata.Display);
    })

    ScorecardInitiative.DetailData(data);
    $("#ScorecardInitiative-Detail").modal("show");
    
}
