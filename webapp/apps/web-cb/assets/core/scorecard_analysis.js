var ScorecardAnalysis = {
    Name:ko.observable(""), //Business Metrics Name that you've click
    Processing:ko.observable(false),
    Data:ko.observableArray([]),
    MetricsCountryRankingData:ko.observableArray([]),
    IsRelativeCumulative:ko.observable(false),
    IsMetricsCountryRanking:ko.observable(false),
    // DataSource
    DataReference:ko.observable({
        LastPeriod: new Date(),
    }),
    BusinessMetrics:ko.observable(""),
    Region:ko.observable(""),
    Country:ko.observable(""),
    StatusScorecardFullyearRankingTree:ko.observable(false),
    StatusScorecardMetricCountryRankingTree:ko.observable(false),
}
ScorecardAnalysis.Get = function(data){
    ScorecardAnalysis.Region(c.Filter.RegionOne());
    ScorecardAnalysis.Country(c.Filter.CountryOne());  
    ScorecardMetricCountryRanking.CountryValue("")
    ScorecardFullyearRanking.CountryValue('')  
    ScorecardCountryAnalysis.ActiveCountry('')
    
    $("#TrendChart").html('');
    if(typeof data !== "undefined"){
        ScorecardAnalysis.IsMetricsCountryRanking(false);
        if(data.Type=="cumulative"){
            ScorecardAnalysis.IsRelativeCumulative(true);
        }else{
            ScorecardAnalysis.IsRelativeCumulative(false);
        }
        ScorecardAnalysis.DataReference(data);
        ScorecardAnalysis.Name(data.DataPoint);
        ScorecardAnalysis.GetData();
        ScorecardFullyearProjection.AfterRender(false)
        $('#ScorecardCountryAnalysisandFullYear .nondef').removeClass('active')
        $('#ScorecardCountryAnalysisandFullYear .def').addClass('active')	
        $('#data-visualisation3').html('')
        if(!ScorecardAnalysis.StatusScorecardFullyearRankingTree()){
            ScorecardFullyearRanking.InitTree()
            ScorecardAnalysis.StatusScorecardFullyearRankingTree(true)
        }
        ResetCountryRegionFilter2_FYR()
        ScorecardCountryAnalysis.Period(getUTCDate(data.LastPeriod))
    }else{
        ScorecardAnalysis.Name("Metrics Country Ranking");
        ScorecardAnalysis.IsMetricsCountryRanking(true);
        ScorecardAnalysis.GetMetricsRankingData();
        $('#ScorecardCountryAnalysisandFullYear .nondef').removeClass('active')
        $('#ScorecardCountryAnalysisandFullYear .def').removeClass('active')
        if(!ScorecardAnalysis.StatusScorecardMetricCountryRankingTree()){
            ScorecardMetricCountryRanking.InitTree()
            ScorecardAnalysis.StatusScorecardMetricCountryRankingTree(true)
        }
        ResetCountryRegionFilter2_MCR()
        ScorecardCountryAnalysis.Period(getUTCDate(new Date()))
    }
}
ScorecardAnalysis.Render = function(DataSource){
    if(typeof DataSource !== "undefined" && DataSource !== null){
        ScorecardCountryAnalysis.Get(DataSource.CountryAnalysis);
        ScorecardFullyearRanking.Get(DataSource.FullYearRanking);
    }
}
ScorecardAnalysis.RenderMetricsCountryRanking = function(DataSource){
    if(typeof DataSource !== "undefined" && DataSource !== null){
        var sort = ScorecardMetricCountryRanking.SortBy();
        // console.log(sort)
        for(var x in DataSource){
            var arr = DataSource[x].DataList;
            var dList = [];
            switch(DataSource[x].Type){
                case "cumulative":
                    if(sort==="gaptotarget"){
                        dList = Enumerable.From(arr).OrderByDescending("$.PercentGap").ToArray();
                    }else{
                        dList = Enumerable.From(arr).OrderByDescending("$.YoY").ToArray();
                    }
                    break;
                case "spot":
                        if(DataSource[x].IsHigherIsBetter){
                            dList = Enumerable.From(arr).OrderByDescending("$.Actual").ToArray();
                        }else{
                            dList = Enumerable.From(arr).OrderBy("$.Actual").ToArray();
                        }
                    break;
                default:
                    dList = Enumerable.From(arr).OrderByDescending("$.PercentGap").ToArray();
                    break;
            }
            for(var d in dList){
                dList[d].Rank = (parseInt(d)+1);
            }
            DataSource[x].DataList = dList;
        }
        ScorecardAnalysis.MetricsCountryRankingData(DataSource)
        // ScorecardMetricCountryRanking.Get();   
    }
}
ScorecardAnalysis.GetData = function(){
    var DataReference = ko.mapping.toJS(ScorecardAnalysis.DataReference);
    ScorecardAnalysis.Processing(true)
    var url = "/web-cb/scorecardanalysis/getdata";
    var parm = {
        BusinessMetrics:DataReference.Id,
        Region:ScorecardAnalysis.Region(),
        Country:ScorecardAnalysis.Country(),
        Period:kendo.toString(getUTCDate(DataReference.LastPeriod),"yyyyMMdd"),
    };
    ajaxPost(url,parm,function(res){
        ScorecardAnalysis.Processing(false)
        if (res.IsError){
            swal("Error!",res.Message,"error");
        }
        ScorecardAnalysis.Data(res.Data);
        ScorecardAnalysis.Render(res.Data);
        $("#Scorecard-Analysis").modal("show");
    });
}
ScorecardAnalysis.GetMetricsRankingData = function(){
    ScorecardAnalysis.Processing(true)
    var url = "/web-cb/scorecardanalysis/getmetricsrankingdata";
    var parm = {
        Region:ScorecardAnalysis.Region(),
        Country:ScorecardAnalysis.Country(),
        SortBy:ScorecardMetricCountryRanking.SortBy()
        // Period:kendo.toString(getUTCDate(DataReference.LastPeriod),"yyyyMMdd"),
    };
    ajaxPost(url,parm,function(res){
        ScorecardAnalysis.Processing(false)
        if (res.IsError){
            swal("Error!",res.Message,"error");
        }
        // console.log("----",res.Data)
        // ScorecardAnalysis.MetricsCountryRankingData(res.Data);
        ScorecardAnalysis.RenderMetricsCountryRanking(res.Data);
        $("#Scorecard-Analysis").modal("show");
        ScorecardMetricCountryRanking.Get()
    });
}