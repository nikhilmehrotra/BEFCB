var Now = new Date();
// var SeqLast = ko.observable(0);
var Initiative = {
    ScorecardData:ko.observableArray([]),
    UserList:ko.observableArray([]),
    LogMessage: ko.observableArray([]),
    OpenId:ko.observable(""),
    Processing:ko.observable(true),
    Mode:ko.observable(""),
    ModeEditedProgress:ko.observable(false),
    Phase:ko.observable(1),
    InitiativeType:ko.observable(""),
    CurrentSearchKeyword: ko.observable(''),

    LCWidth: ko.observable(100), //Set LifeCycle Width
    ColorList:ko.observableArray(["#0077AC","#009D44","#001143","#4EBF45","#6E7378","#00B3EB","#002BAE","#199489","#0077AC","#009D44","#001143","#4EBF45","#6E7378","#00B3EB","#002BAE","#199489"]),
    // Form Input
    Attachments:ko.observableArray([]),
    Map:ko.observableArray([]),
    FormValue:ko.observable(),
    Data:{
        Id:0,
        ProjectName: "",
        ProjectDriver:"",
        BusinessDriverId: "",
        BusinessDriverImpact: "",
        CBLedInitiatives: false,
        EX: false,
        OE: false,
        StartDate: new Date(),
        FinishDate: new Date(),
        ProblemStatement: "",
        ProjectDescription: "",
        ProjectManager: "",
        BusinessImpact: "",
        InvestmentId: "",
        AccountableExecutive: "",
        TechnologyLead: "",
        ProgressCompletion: 0,
        PlannedCost: 0,
        ProjectClassification: "",
        Attachments: [],
        InitiativeType: "",
        InitiativeID: "",
        LifeCycleId: "",
        SubLifeCycleId: "",
        IsGlobal:false,
        Region: [],
        Country:[],
        Type:"",
        ImprovedEfficiency:false,
        ClientExperience: false,    
        OperationalImprovement:false,      
        CSRIncrease:false, 
        TurnAroundTime:false, 
        ImprovedEfficiencyCurrent: 0,     
        ImprovedEfficiencyTarget: 0,  
        ClientExperienceCurrent: 0,       
        ClientExperienceTarget: 0,  
        OperationalImprovementCurrent: 0, 
        OperationalImprovementTarget: 0, 
        CSRIncreaseCurrent: 0,             
        CSRIncreaseTarget: 0, 
        TurnAroundTimeCurrent: 0,         
        TurnAroundTimeTarget: 0, 
        Milestones:[],
        CommentList:[],
        SetAsComplete: false,
        CompletedDate: new Date(),
        DisplayProgress: "green",
        Sponsor: "",
        IsInitiativeTracked: false,
        MetricBenchmark: "",
        AdoptionScoreDenomination: "",
        UsefulResources: "",
    },
    TmpComment: ko.observable(),
    SelectedData:ko.observable(),
    // Instance
    Milestone:{
        Id: 0,
        Name:"",
        StartDate:null,
        EndDate:null,
        Country:[],
        Completed:false,
        CompletedDate:new Date(),
        Seq: 0,
    },
    // Sources
    InitiativeOwnerList:ko.observableArray([]),
    BusinessMetricList:ko.observableArray([]),
    DirectIndirectList:ko.observableArray([
        { text: "N/A", value: 0 },  
        { text: "Direct", value: 1 },  
        { text: "Indirect", value: 2 },  
    ]),
    TypeList:ko.observableArray([
        { text: "Bank Wide", value: "BANKWIDE" },  
        { text: "CB LED Initiative", value: "CBLED" },  
    ]),
    AllBusinessDriverList:ko.observableArray([]),
    RegionList:ko.observableArray([]),
    CountryList:ko.observableArray([ ]),
    LifeCycleList:ko.observableArray([]),
    SubLifeCycleList:ko.observableArray([]),
    SCCategoryList:ko.observableArray([]), 
    BusinessDriverList:ko.observableArray([]), 
    BusinessImpactList:ko.observableArray([
        { text: "Low", value: "Low" },  
        { text: "Medium", value: "Medium" },
        { text: "High", value: "High" }
    ]),
    BusinessDriverImpact:ko.observableArray([        
        { text: "Primary", value: "Primary" },
        { text: "Secondary", value: "Secondary" }, 
    ]),
    ProjectClassificationList:ko.observableArray([        
        { text: "Small", value: "Small" },
        { text: "Medium", value: "Medium" },
        { text: "Large", value: "Large" },  
    ]),
    DisplayColorList:ko.observableArray([
        {name:"Red",value:"red"},
        {name:"Amber",value:"amber"},
        {name:"Green",value:"green"},
    ]),
    MetricsTypeList :ko.observableArray([
        {name: "Dollar Value ($)", value: "DOLLAR"},
        {name: "Numeric Value", value: "NUMERIC"},
        {name: "Percentage (%)", value: "PERCENTAGE"},
    ]),
    SponsorList:ko.observableArray([])
}
Initiative.CheckOwnedAccess = function(){
    var id = Initiative.OpenId();
    var sources = [];
    if(typeof c !== undefined){
        sources = c.OwnedData();
    }else{
        sources = scchart.OwnedData();
    }
    if(sources.indexOf(id) >= 0 ){
        return true;
    }else{
        return false;
    }
}

Initiative.Phase.subscribe(function(newval){
    if(newval==2){
        var formValue = ko.toJS(Initiative.FormValue());
        var color = (formValue.DisplayProgress == 'amber') ? '#ffd24d' : (formValue.DisplayProgress == 'green') ? '#6ac17b' : '#f74e4e';
        $('.pgbaredit').children("span").remove();
        $('#dashboard .k-progressbar .k-state-selected .k-progress-status-wrap .k-progress-status').remove();
        $('#dashboard .k-progressbar .k-state-selected').css('background-color', color).css('border-color', color);
        if(parseInt(kendo.toString(Initiative.FormValue().ProgressCompletion(), 'n0')) == 100 ){
            var classColorAfter = "progressbar"
            classColorAfter = (formValue.DisplayProgress == 'amber') ? classColorAfter+'yellowafter' : (formValue.DisplayProgress == 'green') ? classColorAfter+'greenafter' : classColorAfter+'redafter';
            $('.pgbaredit').addClass(classColorAfter);
        } else if (Initiative.FormValue().ProgressCompletion() != 0) {
            // $('.pgbaredit').children("span").remove();
            $('#dashboard .k-progressbar .k-state-selected .k-progress-status-wrap').css("width","");
        } else if (Initiative.FormValue().ProgressCompletion() >= 0 && Initiative.FormValue().ProgressCompletion() <= 5) {
            $('.pgbaredit .k-progress-status-wrap').css({'color':'black', 'text-align': 'left'});
        }
    }
})

Initiative.GetSponsorList = function(){
    ajaxPost("/web-cb/befsponsor/getdata",{},function(res){
        Initiative.SponsorList(res.Data);
    });
}

Initiative.GetCountryDS = function(e){
    // var d = Initiative.FormValue().Region();
    var d = e.sender._old;

    var arr = [];
    if(d.length > 0){
        for(var i in d){
            var temp_arr = Enumerable.From(c.RegionalData()).Where("$.Major_Region === '"+d[i]+"'").GroupBy("$.Country").Select("{_id:$.Key()}").ToArray();
            arr = arr.concat(temp_arr);
        }
    } else{
        for(i in c.RegionalData()){
            // console.log(c.RegionalData()[i])
            arr.push({"_id": c.RegionalData()[i].Country })
        }
    }
    
    Initiative.CountryList(arr);
    Initiative.FormValue().Country([])
}

Initiative.CalcTableChart = function(){
    var shitstein = $('#app-header .navbar').height()+$('#top-header').height()+ $('#crumbs').height()
    var top = shitstein*1.1;
    var left = $('#BDFilter').width();
    return "top:"+top+"px;left:"+left+"px;";
}

Initiative.FixedHeader = function(){
    Initiative.RemoveFixedHeader()
    var topheader = $('#app-header .navbar').height();
    var topheight = $('#top-header').height();
    var iabdHeight = $('#initiativeAllBDF').height() *1.5;
    var crumbsMargin = $('#BDFilter').width()*0.95;
    var iabdTop = (topheight+topheader) *1.1;
    var bdfTop = iabdTop + iabdHeight;
    var sidebarwidth = 0;
    var initiativewidth = $('#initiativeAllBDF').width();
    
    if(Sidebar.IsExpand()){
        sidebarwidth = "180px"
    }else{
        sidebarwidth = "40px"
    }
    $('#app-header').css({
        "z-index":4, 
        "position": "fixed",
        "right":"0",
        "left":sidebarwidth,
        "width":"auto",
        "-webkit-transition": "all 250ms ease-in-out",
        "-moz-transition": "all 250ms ease-in-out",
        "-ms-transition": "all 250ms ease-in-out",
        "-o-transition": "all 250ms ease-in-out",
        "transition": "all 250ms ease-in-out",
    })
    $('#top-header').css({
        "position": "fixed",
        "z-index": "3",
        "background-color": "#fff",
        "top": topheader+"px",
        "right":"0",
        "left":sidebarwidth,
        "width":"auto",
        "-webkit-transition": "all 250ms ease-in-out",
        "-moz-transition": "all 250ms ease-in-out",
        "-ms-transition": "all 250ms ease-in-out",
        "-o-transition": "all 250ms ease-in-out",
        "transition": "all 250ms ease-in-out",
    })
    // $('#initiativeAllBDF').css({
    //     "position": "fixed",
    //     "top": iabdTop+"px",
    //     "z-index": "2",
    //     "background-color": "white",
    //     "width": "100%",
    //     "height": iabdHeight+"px",
    // })
    $('#initiativeAllBDF').css({
        "background-color": "white",
        "padding-left":"5px",
    })
    $('#abalabal').css({
        "position":"fixed",
        "z-index":"999",
        "top": iabdTop+"px",
        "left":sidebarwidth,
        "right":"0",
        "-webkit-transition": "all 250ms ease-in-out",
        "-moz-transition": "all 250ms ease-in-out",
        "-ms-transition": "all 250ms ease-in-out",
        "-o-transition": "all 250ms ease-in-out",
        "transition": "all 250ms ease-in-out",
    })

    // $('#selectall-initiative').css({
    //     "width":"15%"
    // })
    // $('#crumbs').css({
    //     "position": "fixed",
    //     "top": iabdTop+"px",
    //     "bottom": "0px",
    //     "left": "0px",
    //     "right": "0px",
    //     "z-index": "2",
    //     "margin-left": crumbsMargin+"px",
    //     "width": "82.5%",
    // })
    $('#BDFilter').css({
        "position": "absolute",
        "top": bdfTop-7+"px",
            "width":initiativewidth+30+"px",
    })
    $('#tabelChart').css({
        "top": bdfTop-7+"px",
            "left":initiativewidth+30+"px",

    })
    setTimeout(function() {
        var initiativewidth = $('#initiativeAllBDF').width();
        $('#BDFilter').css({
            // "position": "absolute",
            // "top": bdfTop-7+"px",
            "width":initiativewidth+30+"px",
            "-webkit-transition": "all 250ms ease-in-out",
            "-moz-transition": "all 250ms ease-in-out",
            "-ms-transition": "all 250ms ease-in-out",
            "-o-transition": "all 250ms ease-in-out",
            "transition": "all 250ms ease-in-out",
        })
        $('#tabelChart').css({
            // "top": bdfTop-7+"px",
            "left":initiativewidth+30+"px",
            "-webkit-transition": "all 250ms ease-in-out",
            "-moz-transition": "all 250ms ease-in-out",
            "-ms-transition": "all 250ms ease-in-out",
            "-o-transition": "all 250ms ease-in-out",
            "transition": "all 250ms ease-in-out",
        })
    }, 500);

    $('#iFooter').css("margin-top", bdfTop+"px")
}

Initiative.RemoveFixedHeader = function(){
    $('#app-header').css({
        "z-index":3, 
        "position": "",
        "width":"100%",
        "left":"0",
        "right":"0",
    })
    $('#top-header').css({
        "position": "",
        "z-index": "3",
        "background-color": "",
        "top": "",
        "width":"100%",
        "left":"0",
        "right":"0",
    })
    // $('#initiativeAllBDF').css({
    //     "position": "",
    //     "top": "",
    //     "z-index": "3",
    //     "background-color": "",
    //     "width": "",
    //     "height": "",
    // })
    $('#initiativeAllBDF').css({
        "background-color": "",
        "padding-left":"15px",
    })
    $('#abalabal').css({
        "position":"",
        "z-index":"999",
        "top":"",
        "left":"0",
    })
    // $('#selectall-initiative').css({
    //     "width":"100%"
    // })
    // $('#crumbs').css({
    //     "position": "",
    //     "top": "",
    //     "bottom": "",
    //     "left": "",
    //     "right": "",
    //     "z-index": "3",
    //     "margin-left": "",
    //     "width": "98.5%",
    // })
    $('#tabelChart').css({
        // "height":" auto !important",
        // "position": "absolute",
        "top": "",
        "left": "",
    })
    $('#BDFilter').css({
        "position": "",
        "top": "",
        "width":"",
    })
    $('#iFooter').css("margin-top", "")
}

// Initiative.Data.Region.subscribe(function(newVal){
    // console.log(newVal)
    // var arr = [];
    // if(newVal.length > 0){
    //     for(var i in newVal){
    //         var temp_arr = Enumerable.From(c.RegionalData()).Where("$.Major_Region === '"+newVal[i]+"'").GroupBy("$.Country").Select("{_id:$.Key()}").ToArray();
    //         arr = arr.concat(temp_arr);
    //     }
    // } else{
    //     for(i in c.RegionalData()){
    //         // console.log(c.RegionalData()[i])
    //         arr.push({"_id": c.RegionalData()[i].Country })
    //     }
    // }
    
  // c.CountryList(arr);
  // c.Filter.Country([])
  // c.GetData();
// });
function KeyMetric(data){
    var self = this;
    if(data!==undefined){
        self.BMId = ko.observable(data.BMId);
        self.DirectIndirect = ko.observable(data.DirectIndirect);
    }else{
        self.BMId = ko.observable("");
        self.DirectIndirect = ko.observable(0);
    }
    return self;
}
Initiative.AddKeyMetric = function(d){
    d.KeyMetrics.push(new KeyMetric(undefined))
}
Initiative.RemoveKeyMetric = function(bdIndex,keyIndex){
    Initiative.Map()[bdIndex].KeyMetrics.remove(Initiative.Map()[bdIndex].KeyMetrics()[keyIndex]);
}
function MapData(data){
    var self = this;
    if(data!==undefined){
        var selectedBD = [];
        if(data.SCCategory===null || data.SCCategory===""){
            selectedBD = Initiative.AllBusinessDriverList();
        }else{
            selectedBD = Enumerable.From(Initiative.AllBusinessDriverList()).Where("$.Parentid==='"+data.SCCategory+"'").ToArray();
        }
        var selectedLC = Enumerable.From(Initiative.LifeCycleList()).Where("$.Id==='"+data.LifeCycle+"'").FirstOrDefault();
        self = {SubLifeCycleList:ko.observableArray(selectedLC.SubLC),BDList:ko.observableArray(selectedBD),SCCategory:ko.observable(data.SCCategory),LifeCycle:ko.observable(data.LifeCycle),SubLifeCycle:ko.observable(data.SubLifeCycle),BusinessDriver:ko.observable(data.BusinessDriver),BusinessImpact:ko.observable(data.BusinessImpact),ImpactonBusinessDriver:ko.observable(data.ImpactonBusinessDriver),Type:ko.observable(data.Type)};
        self.KeyMetrics = ko.observableArray([]);
        for(var k in data.KeyMetrics){
            self.KeyMetrics.push(new KeyMetric(data.KeyMetrics[k]))
        }
    }else{
        var selectedBD = Initiative.AllBusinessDriverList();
        self = {SubLifeCycleList:ko.observableArray([]),BDList:ko.observableArray(selectedBD),SCCategory:ko.observable(""),LifeCycle:ko.observable(""),SubLifeCycle:ko.observable(""),BusinessDriver:ko.observable(""),BusinessImpact:ko.observable("Primary"),ImpactonBusinessDriver:ko.observable(""),Type:ko.observable("")};    
        self.KeyMetrics = ko.observableArray([]);
    }
    self.LifeCycle.subscribe(function(val){
        var selectedLC = Enumerable.From(Initiative.LifeCycleList()).Where("$.Id==='"+val+"'").FirstOrDefault();
        if(selectedLC===undefined){
            self.SubLifeCycleList([]);
        }else{
            self.SubLifeCycleList(selectedLC.SubLC);
        }
    })
    self.SCCategory.subscribe(function(val){
        if(val===null || val===""){
            self.BDList(Initiative.AllBusinessDriverList());
        }else{
            var selectedBD = Enumerable.From(Initiative.AllBusinessDriverList()).Where("$.Parentid==='"+val+"'").ToArray();
            if(selectedBD===undefined){
                self.BDList([]);
            }else{
                self.BDList(selectedBD);
            }
        }
        
        self.BusinessDriver("");
    });
    
    return self;
}

Initiative.GetSCCategory = function(id){
    var d = Enumerable.From(ko.mapping.toJS(Initiative.ScorecardData())).Where("$.Idx === '"+id+"'").FirstOrDefault();
    if(d===undefined){
        return "";
    }else{
        return d.Name+" - ";
    }
}

// SelectIFile - Upload function through attachments view list in Overview initiative
Initiative.SelectIFile = function(){
    var file = document.getElementById('attfile').files;
    if (file.length === 0){
        Initiative.FormValue().Attachments([]);
    }
    var formData = new FormData();
    Initiative.Attachments([]);
    for(var i = 0; i < file.length; i++){
        Initiative.Attachments.push(ko.mapping.fromJS({
            filename:file[i].name,
            description:""
        }))
    }
}

// Select File - Upload function through initiative form
Initiative.SelectFile = function(){
    var file = document.getElementById('attfile').files;
    if (file.length === 0){
        Initiative.FormValue().Attachments([]);
    }
    var formData = new FormData();
    for(var i = 0; i < file.length; i++){
        Initiative.FormValue().Attachments.push(ko.mapping.fromJS({
            filename:file[i].name,
            description:""
        }))
    }
}

Initiative.NewMap = function(){
    Initiative.Map.push(new MapData());
}
Initiative.SaveComplete = function(){
    Initiative.Reset();
    Initiative.Mode("");
    Initiative.ModeEditedProgress(false)
    c.GetData(true); //Is Refreshing Initiative
    $("#ScorecardInitiative-Detail").modal("hide");
    ScorecardInitiative.GetData();
}
Initiative.Close = function(){
    $("#initiative").modal("hide");
    Initiative.Reset();
    Initiative.Mode("");
    Initiative.ModeEditedProgress(false)
    c.GetData();
    $("#ScorecardInitiative-Detail").modal("hide");
    ScorecardInitiative.GetData();
}
Initiative.GetDataSource = function(){
    var activeTab = "InitiativeTab";
    Initiative.InitiativeType(activeTab);
    var sources = undefined
    if(typeof c !== "undefined"){
        sources = c.DataSource().Data;
    }else{
        sources = scchart.DataSource().Data;
    }
     
    Initiative.LifeCycleList(sources.MasterLifeCycle);
    var arr = sources.BusinessDriverList;
    for(var i in arr){
        var d = Enumerable.From(sources.AllSummaryBusinessDriver).Where("$.Idx === '"+arr[i].BusinessDriverId+"'").FirstOrDefault();
        if(d!== undefined){
            arr[i].SCCategory = d.Parentid;
        }else{
            arr[i].SCCategory = "";
        }
    }
    Initiative.BusinessDriverList(arr);
    Initiative.AllBusinessDriverList(sources.AllSummaryBusinessDriver);
    

}
Initiative.GetBusinessDriver = function(id){
    if (id != ""){
        var row = Initiative.AllBusinessDriverList().find(function (d) {
            return d.Idx == id
        })
        
        if (row !== undefined) {
            return row.Name
        }
    }

    return ""
}
Initiative.RemoveMilestone = function(obj){
    console.log(obj.Id())
    $(".row[index='"+obj.Id()+"']").remove();
    Initiative.FormValue().Milestones.remove(obj);
    Initiative.GetProgressCompletion();
}
Initiative.AddMilestone = function(){
    tmp = Initiative.Milestone;
    tmp.Id = Initiative.FormValue().Milestones().length;
    tmp.Seq = Initiative.FormValue().Milestones().length + 1;

    var d = ko.mapping.fromJS(tmp);

    d.Name.subscribe(function(){
        Initiative.GetProgressCompletion();
    });
    d.StartDate.subscribe(function(){
        Initiative.GetProgressCompletion();
    });
    d.EndDate.subscribe(function(){
        Initiative.GetProgressCompletion();
    });
    
    Initiative.FormValue().Milestones.push(d);
        Initiative.SortablesListMilestone();
}
Initiative.ShowDetail = function(formValue){
    $("#initiative").modal("show");
    Initiative.Phase(2);
    Initiative.Mode("show");
    if(formValue.KeyMetrics===undefined){
        formValue.KeyMetrics = [];
    }
    Initiative.Map([{
        LifeCycle: formValue.LifeCycleId,
        SCCategory:formValue.SCCategory,
        SubLifeCycle: formValue.SubLifeCycleId,
        BusinessDriver: formValue.BusinessDriverId,
        BusinessImpact: formValue.BusinessDriverImpact,
        ImpactonBusinessDriver: formValue.BusinessImpact,
        KeyMetrics:formValue.KeyMetrics
    }]);
    Initiative.FormValue(ko.mapping.fromJS(formValue));
}
Initiative.Reset = function(){
    Initiative.Map([]);
    Initiative.Map.push(new MapData());
    Initiative.FormValue(ko.mapping.fromJS(Initiative.Data));
}
Initiative.Add = function(formValue,isHelping){
    Initiative.GetSponsorList()
    Initiative.SCCategoryList(Scorecard.Data());
    Initiative.RegionList(c.RegionList());
    Initiative.CountryList(c.CountryList());
    
    // console log(formValue.BusinessDriverId);
    Initiative.GetDataSource();
    $("#initiative").modal("show");
    Initiative.Phase(1);
    if(formValue !== undefined && formValue.InitiativeID !== undefined){
        Initiative.Map([]);
        if(formValue.KeyMetrics===undefined){
            formValue.KeyMetrics = [];
        }
        Initiative.Map.push(new MapData({   
            LifeCycle: formValue.LifeCycleId,
            SCCategory:formValue.SCCategory == undefined ? "" : formValue.SCCategory,
            SubLifeCycle: formValue.SubLifeCycleId == undefined ? "" : formValue.SubLifeCycleId,
            BusinessDriver: formValue.BusinessDriverId,
            BusinessImpact: formValue.BusinessDriverImpact,
            ImpactonBusinessDriver: formValue.BusinessImpact,
            Type: formValue.type,
            KeyMetrics:formValue.KeyMetrics
        }));
        // Get SubLC DataSource
        var selectedLC = Enumerable.From(Initiative.LifeCycleList()).Where("$.Id==='"+formValue.LifeCycleId+"'").FirstOrDefault();
        if(selectedLC===undefined){
            Initiative.SubLifeCycleList([]);
        }else{
            Initiative.SubLifeCycleList(selectedLC.SubLC);
        }
        formValue.FinishDate = new Date(formValue.FinishDate);
        formValue.StartDate = new Date(formValue.StartDate);
        formValue.ImprovedEfficiency = false;
        if(formValue.ImprovedEfficiencyCurrent!==0||formValue.ImprovedEfficiencyTarget!==0){
            formValue.ImprovedEfficiency = true;
        }
        formValue.ClientExperience = false;   
        if(formValue.ClientExperienceCurrent!==0||formValue.ClientExperienceTarget!==0){
            formValue.ClientExperience = true;
        }
        formValue.OperationalImprovement = false;
        if(formValue.OperationalImprovementCurrent!==0||formValue.OperationalImprovementTarget!==0){
            formValue.OperationalImprovement = true;
        }
        formValue.CSRIncrease = false;
        if(formValue.CSRIncreaseCurrent!==0||formValue.CSRIncreaseTarget!==0){
            formValue.CSRIncrease = true;
        }
        formValue.TurnAroundTime = false; 
        if(formValue.TurnAroundTimeCurrent!==0||formValue.TurnAroundTimeTarget!==0){
            formValue.TurnAroundTime = true;
        }
        for(var m in formValue.Milestones){
            formValue.Milestones[m].StartDate = new Date(formValue.Milestones[m].StartDate)
            formValue.Milestones[m].EndDate = new Date(formValue.Milestones[m].EndDate)
            formValue.Milestones[m].CompletedDate = new Date(formValue.Milestones[m].CompletedDate)
            if(formValue.Milestones[m].Completed === undefined){
                formValue.Milestones[m].Completed = false;
            }
        }
        Initiative.FormValue(ko.mapping.fromJS(formValue));
        for(var m in Initiative.FormValue().Milestones()){
            var d = Initiative.FormValue().Milestones()[m];
            d.Name.subscribe(function(){
                Initiative.GetProgressCompletion();
            });
            d.StartDate.subscribe(function(){
                Initiative.GetProgressCompletion();
            });
            d.EndDate.subscribe(function(){
                Initiative.GetProgressCompletion();
            });
        }
        Initiative.GetProgressCompletion();
        // console.log(Initiative.FormValue());
    }
    else{
        Initiative.Mode("new");
        Initiative.Map([]);
        Initiative.Map.push(new MapData());
        Initiative.FormValue(ko.mapping.fromJS(Initiative.Data));
        Initiative.FormValue().IsGlobal(true);
        Initiative.FormValue().Milestones.push(ko.mapping.fromJS(Initiative.Milestone));

    }
    if(!isHelping && formValue !== undefined){
    
        var color = (formValue.DisplayProgress == 'amber') ? '#ffd24d' : (formValue.DisplayProgress == 'green') ? '#6ac17b' : '#f74e4e';
        $('#dashboard .k-progressbar .k-state-selected').css('background-color', color).css('border-color', color);

        /// custom progress bar style
        var classColor = (formValue.DisplayProgress == 'amber') ? 'progressbaryellow' : (formValue.DisplayProgress == 'green') ? 'progressbargreen' : 'progressbarred';
        $('.pgbaredit').addClass(classColor).css("margin", "0 5px");

        Initiative.FormValue().DisplayProgress.subscribe(function(value){
            var color = (value == 'amber') ? '#ffd24d' : (value == 'green') ? '#6ac17b' : '#f74e4e';
            $('#dashboard .k-progressbar .k-state-selected').css('background-color', color).css('border-color', color);
        });

        var tmpProgressCompletion = Initiative.FormValue().ProgressCompletion()

        $(".completedisable").prop('disabled',Initiative.FormValue().SetAsComplete());
        $('[name="DisplayProgress"]').data('kendoDropDownList').enable(!Initiative.FormValue().SetAsComplete())
        $('[name="ProjectClassification"]').data('kendoDropDownList').enable(!Initiative.FormValue().SetAsComplete())
        $('[name="pclassval"]').each(function(i,e){
            $(e).data('kendoMultiSelect').enable(!Initiative.FormValue().SetAsComplete())  
        })
        $('#completedisableregion').data('kendoMultiSelect').enable(!Initiative.FormValue().SetAsComplete())
        $('#completedisablecountry').data('kendoMultiSelect').enable(!Initiative.FormValue().SetAsComplete())

        $('#ProjectManager').data('kendoMultiSelect').enable(!Initiative.FormValue().SetAsComplete())
        $('#AccountableExecutive').data('kendoMultiSelect').enable(!Initiative.FormValue().SetAsComplete())
        $('#TechnologyLead').data('kendoMultiSelect').enable(!Initiative.FormValue().SetAsComplete())
        $('#Sponsor').data('kendoDropDownList').enable(!Initiative.FormValue().SetAsComplete())
        $('#PlannedCost').data('kendoNumericTextBox').enable(!Initiative.FormValue().SetAsComplete())

        $('[name="AdoptionScoreDenomination"]').data('kendoDropDownList').enable(!Initiative.FormValue().SetAsComplete())
        $('#MetricBenchmark').prop('disabled',Initiative.FormValue().SetAsComplete());
        $('#UsefulResources').prop('disabled',Initiative.FormValue().SetAsComplete());
        $('#IsInitiativeTracked').bootstrapSwitch('disabled',Initiative.FormValue().SetAsComplete());

        Initiative.FormValue().SetAsComplete.subscribe(function(value){
            $(".completedisable").prop('disabled',value)
            $('[name="DisplayProgress"]').data('kendoDropDownList').enable(!value)
            $('[name="ProjectClassification"]').data('kendoDropDownList').enable(!value)
            $('[name="pclassval"]').each(function(i,e){
                $(e).data('kendoMultiSelect').enable(!value)  
            })
            $('#completedisableregion').data('kendoMultiSelect').enable(!value)
            $('#completedisablecountry').data('kendoMultiSelect').enable(!value)

            $('#ProjectManager').data('kendoMultiSelect').enable(!value)
            $('#AccountableExecutive').data('kendoMultiSelect').enable(!value)
            $('#TechnologyLead').data('kendoMultiSelect').enable(!value)
            $('#Sponsor').data('kendoDropDownList').enable(!value)
            $('#PlannedCost').data('kendoNumericTextBox').enable(!value)

            $('[name="AdoptionScoreDenomination"]').data('kendoDropDownList').enable(!value)
            $('#MetricBenchmark').prop('disabled',value);
            $('#UsefulResources').prop('disabled',value);
            $('#IsInitiativeTracked').bootstrapSwitch('disabled',value);

            if(value){
                Initiative.FormValue().ProgressCompletion(100)
                var classColorAfter = "progressbar"
                classColorAfter = (formValue.DisplayProgress == 'amber') ? classColorAfter+'yellowafter' : (formValue.DisplayProgress == 'green') ? classColorAfter+'greenafter' : classColorAfter+'redafter';
                $('.pgbaredit').addClass(classColorAfter);
            } else if (value >= 0 && value <= 5) {
                $('.pgbaredit .k-progress-status-wrap').css({'color':'black', 'text-align': 'left'});
            } else {
                Initiative.FormValue().ProgressCompletion(tmpProgressCompletion)
                $('.pgbaredit').children("span").remove();
                $('.pgbaredit .k-state-selected .k-progress-status-wrap').css("width","");
            }   
            // console.log("sip", value, tmpProgressCompletion)
            // var color = (value == 'amber') ? '#ffd24d' : (value == 'green') ? '#6ac17b' : '#f74e4e';
            // $('#dashboard .k-progressbar .k-state-selected').css('background-color', color).css('border-color', color);
        });
    }else{
        Initiative.Processing(false);
    }

    Initiative.Processing(false);
    setTimeout(function() {
        $("#IsInitiativeTracked").bootstrapSwitch();
        $('#IsInitiativeTracked').on('switchChange.bootstrapSwitch', function(event, state) {
          Initiative.FormValue().IsInitiativeTracked(state)
        });
    }, 300);

}
// Initiative.Remove = function(){
//     var d = ko.mapping.toJS(Initiative.SelectedData());
//     ajaxPost("/initiative/remove",{Id:d.InitiativeID},function(res){
//         if(res.IsError){
//             swal("", res.message, "info");
//             return false;
//         }
//         Initiative.Close();
//     });
// }
Initiative.ClearAttachment = function(){
    $("#attfile").val("")
    Initiative.Attachments([]);
}
Initiative.UploadAttachment = function(){
    Initiative.Processing(true);
    var d = ko.mapping.toJS(Initiative.SelectedData());
    d.Id = d._id;
    Initiative.Add(d,true);
    // var att = Initiative.Attachments();
    // $(att).each(function(e,d){
    //     Initiative.FormValue().Attachments.push(d)
    // })
    Initiative.Save(true);
    Initiative.ClearAttachment();
    Initiative.SaveComplete();
}
Initiative.Remove = function(){
    swal({
      title: "Are you sure?",
      text: "You will not be able to recover this!",
      type: "warning",
      showCancelButton: true,
      confirmButtonColor: "#DD6B55",
      confirmButtonText: "Yes, remove it!",
      closeOnConfirm: false
    },
    function(){
    var d = ko.mapping.toJS(Initiative.SelectedData());
    ajaxPost("/web-cb/initiative/remove",{Id:d._id},function(res){
        if(res.IsError){
            swal("", res.message, "info");
            return false;
        }
        swal("Deleted!", "Your Task file has been deleted.", "success");
        Initiative.Close();
    });     
    });
}

Initiative.Edit = function(){
    Initiative.Mode("edit")
    var d = ko.mapping.toJS(Initiative.SelectedData());
    // d.StartDate = getUTCDate(d.StartDate);
    // d.FinishDate = getUTCDate(d.FinishDate);
    d.Id = d._id;

    Initiative.Add(d);
    setTimeout(function() {
        Initiative.SortablesListMilestone();
    }, 300);
}
Initiative.Save = function(isUploadingAttachment){
    var validator = $("#initiativeform").data("kendoValidator");
    if(validator==undefined){
       validator= $("#initiativeform").kendoValidator().data("kendoValidator");
    }
    
    if (validator.validate()) {
        Initiative.GetProgressCompletion();
        Initiative.Processing(true);
        var formData = new FormData();
        var Attachments = Initiative.FormValue().Attachments();
        var Milestones = Initiative.FormValue().Milestones();
        var AttachedFile = [];
        var file = document.getElementById('attfile') !== null ? document.getElementById('attfile').files : [];
        if(isUploadingAttachment){
            var InitiativeAttachment = Initiative.Attachments();
            for(var i = 0; i < file.length; i++){
                formData.append("FileUpload"+i, file[i]);
                AttachedFile.push(file[i].name);
                var fileDescription = Enumerable.From(InitiativeAttachment).Where(function(x){return x.filename() === file[i].name; }).FirstOrDefault().description();
                formData.append("FileDescription"+i, fileDescription);
            }
        }
        var existingfiletotal = 0;
        for(var i in Attachments){
            if(AttachedFile.indexOf(Attachments[i].filename())<0){
                existingfiletotal++;
                formData.append("ExistingAttachment"+i+"FileName", Attachments[i].filename());
                formData.append("ExistingAttachment"+i+"Description", Attachments[i].description());
                formData.append("ExistingAttachment"+i+"UpdatedBy", Attachments[i].updated_by());
                formData.append("ExistingAttachment"+i+"UpdatedDate", kendo.toString(getUTCDate(moment(Attachments[i].updated_date())),"yyyyMMdd"));
            }
        }
        // for(var i in Milestones){
        //     console.log(Milestones[i])
        //     if(Milestones[i].Name().trim()!==""){
        //         formData.append("Milestone"+i, Milestones[i].Name());
        //         formData.append("MilestoneStartDate"+i, kendo.toString(Milestones[i].StartDate(),'yyyyMMdd'));
        //         formData.append("MilestoneEndDate"+i, kendo.toString(Milestones[i].EndDate(),'yyyyMMdd'));
        //         formData.append("MilestoneCountry"+i, Milestones[i].Country().join("|"));
        //         formData.append("MilestoneCompleted"+i, Milestones[i].Completed());
        //         formData.append("MilestoneCompletedDate"+i, kendo.toString(Milestones[i].CompletedDate(),'yyyyMMdd'));
        //         formData.append("MilestoneSeq"+i, parseInt(i)+1);
        //     }
        // }
        var Milestones = Initiative.FormValue().Milestones();

        $('.Milestone-List .row').each(function (i, e) {
            id = $(e).attr("index");
            // console.log(ko.mapping.toJS(Milestones[id]), id)
            z = _.find(Milestones, function(v,i){return v.Id() == id})
            
            if(z.Name().trim()==""){
                
            } else{
                formData.append("Milestone"+i, z.Name().trim());
                formData.append("MilestoneStartDate"+i, kendo.toString(z.StartDate(),'yyyyMMdd'));
                formData.append("MilestoneEndDate"+i, kendo.toString(z.EndDate(),'yyyyMMdd'));
                formData.append("MilestoneCountry"+i, z.Country().join("|"));
                formData.append("MilestoneCompleted"+i, z.Completed());
                formData.append("MilestoneCompletedDate"+i, kendo.toString(z.CompletedDate(),'yyyyMMdd'));
                formData.append("MilestoneSeq"+i, parseInt(i)+1);
            }

            // console.log($(e).attr("index"))

        });

        if(isUploadingAttachment){
            formData.append("filetotal", file.length);
        }else{
            formData.append("filetotal", 0);
        }
        formData.append("existingfiletotal", existingfiletotal);
        formData.append("milestonetotal", Milestones.length);
        var parm = ko.mapping.toJS(Initiative.FormValue());

        if(parm.IsGlobal){
           parm.Country = [];
           parm.Region = []; 
        }else{

            for(var i in parm.Country){
                formData.append("country"+i, parm.Country[i]);
                //search region for this country
                var tmpCountry = parm.Country[i]
                tmpCountry = tmpCountry.trim()
                // console.log('tmpCountry', tmpCountry)
                var tmpRegionObj = _.find(c.RegionalData(), function(v,i){return v.Country == tmpCountry})
                // console.log('tmpRegionObj',tmpRegionObj)
                if(tmpRegionObj != undefined){
                    var Major_Region = tmpRegionObj.Major_Region
                    // console.log('Major_Region',Major_Region)

                    //search region on parm.Region
                    var tmpRegionObjForm = _.find(parm.Region, function(v,i){return v == Major_Region})
                    // console.log('tmpRegionObjForm',tmpRegionObjForm)
                    if(tmpRegionObjForm == undefined){
                        //if region not in parm.Region append to region
                        parm.Region.push(Major_Region);
                    }
                }

            }

            if(parm.Country.length == 0){
                parm.Region = []
            }

            for(var i in parm.Region){
                formData.append("region"+i, parm.Region[i]);
            }
        }
        formData.append("regiontotal", parm.Region.length);
        formData.append("countrytotal", parm.Country.length);

        parm.StartDate = kendo.toString(parm.StartDate,"yyyyMMdd");
        parm.FinishDate = kendo.toString(parm.FinishDate,"yyyyMMdd");
        parm.InitiativeType = Initiative.InitiativeType();
        if(isUploadingAttachment){
            parm.mode = "edit";
        }else{
            parm.mode = Initiative.Mode();
        }
        // console.log(parm.SetAsComplete)
        if(parm.SetAsComplete){
            parm.CompletedDate = kendo.toString(new Date(),"yyyyMMdd");
            parm.ProgressCompletion = 100;
        }

        // console.log(parm)

        if(!parm.ImprovedEfficiency){
            parm.ImprovedEfficiencyCurrent = 0;
            parm.ImprovedEfficiencyTarget = 0;
        }
        if(!parm.ClientExperience){
            parm.ClientExperienceCurrent = 0;
            parm.ClientExperienceTarget = 0;
        }
        if(!parm.OperationalImprovement){
            parm.OperationalImprovementCurrent = 0;
            parm.OperationalImprovementTarget = 0;
        }
        if(!parm.CSRIncrease){
            parm.CSRIncreaseCurrent = 0;
            parm.CSRIncreaseTarget = 0;
        }
        if(!parm.TurnAroundTime){
            parm.TurnAroundTimeCurrent = 0;
            parm.TurnAroundTimeTarget = 0;
        }

        if(parm.Sponsor == "Select.."){
            parm.Sponsor = "";
        }

        for(var i in parm){
            formData.append(i,parm[i]);
        }
        var iMap = ko.mapping.toJS(Initiative.Map()); //Initiative Map
        for(var i in iMap){
            var iMapData = iMap[i].LifeCycle+"|"+iMap[i].SubLifeCycle+"|"+iMap[i].BusinessDriver+"|"+iMap[i].BusinessImpact+"|"+iMap[i].ImpactonBusinessDriver+"|"+iMap[i].Type+"|"+iMap[i].SCCategory;
            formData.append("Map"+i, iMapData);
            var KeyMetrics = "";        
            var MKeyMetrics = iMap[i].KeyMetrics;
            for(var k in MKeyMetrics){
                if(MKeyMetrics[k].BMId!==""){
                    if(k==0){
                        KeyMetrics += MKeyMetrics[k].BMId+","+MKeyMetrics[k].DirectIndirect
                    }else{
                        KeyMetrics += "|"+MKeyMetrics[k].BMId+","+MKeyMetrics[k].DirectIndirect
                    }
                }
            }
            formData.append("KeyMetrics"+i, KeyMetrics);
        }
        formData.append("maptotal", iMap.length);

        $.ajax({
            url: "/web-cb/initiative/save",
            data: formData,
            contentType: false,
            dataType: "json",
            mimeType: 'multipart/form-data',
            processData: false,
            type: 'POST',
            success: function (res) {
                if(!res.success){
                    swal("", res.message, "info");
                    return false;
                }
                if(isUploadingAttachment){
                    swal("", "Upload File Success", "success");
                }else{
                    Initiative.OpenId(res.data.Id);
                    Initiative.SaveComplete();
                }
                // c.SyncSCHeight();
            }
        });
    }
}

Initiative.AddComment = function(){
    Initiative.Processing(true);
    formData = Initiative.SelectedData();
    if(formData.Id === undefined){
        formData.Id = formData._id;
    }
    tmpNewComment = {Username: localStorage.getItem("Username"),DateInput: new Date(), Comment: Initiative.TmpComment(), Editable: false}
    formData.CommentList.push(tmpNewComment)
    // tmpInitiativeData = formData.InitiativeData
    formData.InitiativeData = ""
    ajaxPost("/web-cb/dashboard/commentsave",formData,function(res){

        Initiative.SaveLogComment("Insert", tmpNewComment, formData.Id)

        var sources = [];
        Initiative.SaveComplete();;
        Initiative.TmpComment("")
        // $("#commenttab li").removeClass("active");$("#commentbox .tab-pane").removeClass("active");
        // $($("#commenttab li")[2]).addClass("active");$($("#commentbox .tab-pane")[2]).addClass("active");

        // $(".modal-body.data-list .nav.nav-pills li").removeClass("active");$("#commentbox .tab-pane").removeClass("active");
        // $($(".modal-body.data-list .nav.nav-pills li")[2]).addClass("active");$($("#commentbox .tab-pane")[2]).addClass("active");
        setTimeout(Initiative.scrollComment(), 300);
        
    })
}

// inject span.mark-search on text
Initiative.MarkTextByKeyword = function (target) {
    var regex = new RegExp(Initiative.CurrentSearchKeyword(), 'ig');
    var replacement = '<span class="mark-search">' + Initiative.CurrentSearchKeyword() + '</span>'

    return target.replace(regex, replacement)
}

/*
 * wrapper for Initiative.SelectedData()
 * this wont change the original value,
 * just wrap it on ko.computed
 */
Initiative.SelectedDataMarkTextByKeyword = ko.computed(function () {
    if (Initiative.SelectedData() === undefined) {
        return undefined;
    }

    var raw = ko.mapping.toJS(Initiative.SelectedData);
    var searchKeyword = Initiative.CurrentSearchKeyword();

    if (searchKeyword !== '') {
        var fieldsToBeMarked = [
            'InvestmentId',
            'AccountableExecutive', 
            'ProjectClassification',
            'ProjectDescription',
            'ProjectDriver',
            'ProjectManager',
            'ProjectName',
            'TechnologyLead',
            'ProblemStatement',
        ]

        /*
         * inject span.mark-search on every field's value
         * that contains the search keyword
         */
        fieldsToBeMarked.forEach(function (field) {
            // console.log(field)
            if (!raw.hasOwnProperty(field)) {
                return;
            }

            var isFound = (raw[field].toLowerCase().indexOf(searchKeyword.toLowerCase()) > -1)
            if (isFound) {
                raw[field] = Initiative.MarkTextByKeyword(raw[field])
            }
        })
    }

    return raw;
}, Initiative.SelectedData)


Task.SelectedDataMarkTextByKeyword = ko.computed(function () {
    if (Task.SelectedData() === undefined) {
        return undefined;
    }

    var raw = ko.mapping.toJS(Task.SelectedData);
    var searchKeyword = Initiative.CurrentSearchKeyword();

    if (searchKeyword !== '') {
        var fieldsToBeMarked = [
            'Name',
            'Owner',
            'Statement',
            'Description'
        ]

        /*
         * inject span.mark-search on every field's value
         * that contains the search keyword
         */
        fieldsToBeMarked.forEach(function (field) {
            // console.log(field)
            if (!raw.hasOwnProperty(field)) {
                return;
            }

            var isFound = (raw[field].toLowerCase().indexOf(searchKeyword.toLowerCase()) > -1)
            if (isFound) {
                raw[field] = Initiative.MarkTextByKeyword(raw[field])
            }
        })
    }

    return raw;
}, Task.SelectedData)

Initiative.GetUserList = function(){
    ajaxPost("/web-cb/acluser/getuserlist",{},function(res){
        Initiative.UserList(res.Data);
    });
}
Initiative.GetUserFullname = function(loginid){
    var d = Enumerable.From(Initiative.UserList()).Where("$.loginid === '"+loginid+"'").FirstOrDefault();
    if(typeof d == "undefined"){
        return ""
    }
    return d.fullname;
}
Initiative.Get = function(Id,IsRefreshing){
    Initiative.Processing(true);
    Initiative.OpenId(Id);
    Initiative.Attachments([]);
    Initiative.GetDataSource();
    // var sources = c.DataSource().Data.Project;
    var sources = [];
    if(typeof c !== "undefined"){
        sources = JSON.parse(JSON.stringify(c.AllDataList()));
    }else{
        sources = JSON.parse(JSON.stringify(scchart.DataSource().Data.Project));
    }
    var d = Enumerable.From(sources).Where("$._id === '"+Id+"'").FirstOrDefault();
    if (typeof d === "undefined"){
        return false;
    }
    d.updated_by_fullname = Initiative.GetUserFullname(d.updated_by);
    var SelectedLifeCycle = Enumerable.From(sources).Where("$.InitiativeID === '"+d.InitiativeID+"'").GroupBy("$.LifeCycleId").Select("$.Key()").ToArray();

    d.SelectedLifeCycle = SelectedLifeCycle;
    d.InitiativeData =  Enumerable.From(sources).Where("$.InitiativeID === '"+d.InitiativeID+"'").ToArray();
    for(var m in d.Milestones){
        if(d.Milestones[m].Completed === undefined){
            d.Milestones[m].Completed = false;
        }
    }
    if(d.CommentList == undefined){
        d.CommentList = []
    }

    for(i in d.CommentList){
        d.CommentList[i].Editable = false;
        // console.log(d.CommentList[i])
    }

    if(d.ProgressCompletion===undefined){
        d.ProgressCompletion = parseInt(Math.random()*100);
    }

    if(d.InitiativeBenefits===undefined){
        d.InitiativeBenefits = parseInt(Math.random()*100);
    }
    if(d.IsGlobal===undefined){
        d.IsGlobal = true;
        d.Country = [];
        d.Region = [];
    }
    d.IsViewAttachment = ko.observable(false);
    Initiative.Mode("");

    var currentTabIndex = 0;
    $("#initiative-data ul > li").each(function(i,e){
        if($(e).attr("class") == "active"){
            currentTabIndex = i;
        }
        $(e).removeAttr("class");

    })

    d.IsInitiativeTracked = (d.IsInitiativeTracked == undefined) ? false : d.IsInitiativeTracked;
    d.MetricBenchmark = (d.MetricBenchmark == undefined) ? "" : d.MetricBenchmark;
    d.AdoptionScoreDenomination = (d.AdoptionScoreDenomination == undefined) ? "" : d.AdoptionScoreDenomination;
    d.UsefulResources = (d.UsefulResources == undefined) ? "" : d.UsefulResources;

    Initiative.SelectedData(d);
    $($($("#initiative-data ul > li")[currentTabIndex]).children()).click()
    Initiative.Processing(false);

    // Set Color Progress as per RAG
    var color = (d.DisplayProgress == 'amber') ? '#ffd24d' : (d.DisplayProgress == 'green') ? '#6ac17b' : '#f74e4e';
    $('#dashboard .k-progressbar .k-state-selected').css('background-color', color).css('border-color', color);

    /// custom progress bar style
    var classColor = (d.DisplayProgress == 'amber') ? 'progressbaryellow' : (d.DisplayProgress == 'green') ? 'progressbargreen' : 'progressbarred';
    $('.pgbar').addClass(classColor);
    if (d.ProgressCompletion == 100) {
        var classColorAfter = "progressbar"
        classColorAfter = (d.DisplayProgress == 'amber') ? classColorAfter+'yellowafter' : (d.DisplayProgress == 'green') ? classColorAfter+'greenafter' : classColorAfter+'redafter';
        $('.pgbar').addClass(classColorAfter);
    } else if (d.ProgressCompletion != 0) {
        $('.pgbar').children("span").remove();
        $('.pgbar .k-state-selected .k-progress-status-wrap').css("width","");
    } else if (d.ProgressCompletion >= 0 && d.ProgressCompletion <= 5) {
        $('.pgbar .k-progress-status-wrap').css({'color':'black', 'text-align': 'left'});
    }

    $("#initiative").modal("show");

    // mark label if contains search keyword
    var searchKeyword = Initiative.CurrentSearchKeyword();
    if (searchKeyword !== '') {
        searchKeyword = "work"
        $('#initiative-data #1b label').each(function (i, e) {
            if ($(e).children().size() > 0) {
                return;
            }

            var target = $(e).parent().html();
            var isFound = (target.toLowerCase().indexOf(searchKeyword.toLowerCase()) > -1);
            // console.log($(e).html(), isFound,target.toLowerCase(), searchKeyword.toLowerCase())
            if (isFound) {
                var replacedText = Initiative.MarkTextByKeyword(target);
                $(e).parent().html(replacedText);
            }
        })
    }

    
}
Initiative.SortablesListMilestone = function(){
    $('.Milestone-List').each(function (i, e) {

        var config = {
            containment: "parent"
        };
        $(e).sortable(config);

    });
}
Initiative.GetFile = function(obj){
    var url = '/web-cb/static/ifile/'+Initiative.SelectedData()._id+"/"+encodeURIComponent(obj.filename)
    var http = new XMLHttpRequest();
    http.open('HEAD', url, false);
    http.send();
    if (http.status!=404){
        window.open(url,'_blank');
    }else{
        window.open('/web-cb/static/ifile/'+Initiative.SelectedData().InitiativeID+"/"+encodeURIComponent(obj.filename),'_blank');
    }
}

Initiative.ChangeStartDate = function(e){
    var Data = Initiative.FormValue();
    Data.StartDate(e.sender._old)
    Initiative.GetProgressCompletion()
}

Initiative.ChangeFinishDate = function(e){
    var Data = Initiative.FormValue();
    // console.log(Data.StartDate(), Data.FinishDate())
    Data.FinishDate(e.sender._old)
    Initiative.GetProgressCompletion()
}

Initiative.GetProgressCompletion = function(Data){
    var Data = Initiative.FormValue();
    console.log(Data)
    if( Data !== undefined && (Data.StartDate() == null || Data.FinishDate() == null) ){

    }
    else{
        var Milestones = Data.Milestones();
        Milestones = Enumerable.From(Milestones).Where("$.Name().trim() !== ''").ToArray()
        var ProgressCompletion = 0;
        // console.log( Data.StartDate(),"--", Data.FinishDate() );
        if (Milestones.length == 0){
            if(Now>=Data.FinishDate()){
                ProgressCompletion = 100;
            }else if(Now>=Data.StartDate()){
                ProgressCompletion = (Initiative.GetDaysBetween(Data.StartDate(),Data.FinishDate()) == 0 ? 0 : Initiative.GetDaysBetween(Data.StartDate(),Now)/Initiative.GetDaysBetween(Data.StartDate(),Data.FinishDate()))*100;
            }
        } else {
            tmpMilestones = [];
            for(var i in Milestones){
                var StartDate = Milestones[i].StartDate();
                var EndDate = Milestones[i].EndDate();
                // console.log(StartDate, EndDate)
                if(StartDate !== null && EndDate !== null){
                    tmpMilestones.push(Milestones[i])
                }
            } 

            var TotalDays = Enumerable.From(tmpMilestones).Sum("Initiative.GetDaysBetween($.StartDate(),$.EndDate())");
            
            for(var i in tmpMilestones){            
                var StartDate = Milestones[i].StartDate();
                var EndDate = Milestones[i].EndDate();

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
        Initiative.FormValue().ProgressCompletion(ProgressCompletion);
        /*if (Initiative.FormValue().ProgressCompletion() < 100){
            var formValue = ko.toJS(Initiative.FormValue());
            var classColorAfter = "progressbar"
            classColorAfter = (formValue.DisplayProgress == 'amber') ? classColorAfter+'yellowafter' : (formValue.DisplayProgress == 'green') ? classColorAfter+'greenafter' : classColorAfter+'redafter';
            $('.pgbaredit').removeClass(classColorAfter);
        }else{
            var formValue = ko.toJS(Initiative.FormValue());
            var classColorAfter = "progressbar"
            classColorAfter = (formValue.DisplayProgress == 'amber') ? classColorAfter+'yellowafter' : (formValue.DisplayProgress == 'green') ? classColorAfter+'greenafter' : classColorAfter+'redafter';
            $('.pgbaredit').addClass(classColorAfter);
        }*/
        return true;
    }
}

Initiative.GetDaysBetween = function(date1, date2 ){
     //Get 1 day in milliseconds
  var one_day=1000*60*60*24;

  // Convert both dates to milliseconds
  var date1_ms = date1.getTime();
  var date2_ms = date2.getTime();

  // Calculate the difference in milliseconds
  var difference_ms = date2_ms - date1_ms;
    
  // Convert back to days and return
  return Math.round((difference_ms/one_day),2); 
}
Initiative.GetTime = function(date1){
    return date1.getTime();

}

Initiative.DisplayColorProgress = function(e){

}

Initiative.GetCompletion = function(data){
    var Completion = 0;
    // console.log('completionData',data);
    if(data.SetAsComplete){
        Completion = 100;
    }else{
        data.StartDate = getUTCDate(data.StartDate);
        data.FinishDate = getUTCDate(data.FinishDate);
        var Milestones = data.Milestones;
        for(var m in Milestones){
            Milestones[m].StartDate = getUTCDate(Milestones[m].StartDate);
            Milestones[m].EndDate = getUTCDate(Milestones[m].EndDate);
            Milestones[m].DaysBetween = Initiative.GetDaysBetween(Milestones[m].StartDate,Milestones[m].EndDate);
        }
        Milestones = Enumerable.From(Milestones).Where("$.Name.trim() !== ''").ToArray()
        var ProgressCompletion = 0;
        if (Milestones.length == 0){
            if(Now>data.FinishDate){
                ProgressCompletion = 100;
            }else if(Now>=data.StartDate){
                ProgressCompletion = (Initiative.GetDaysBetween(data.StartDate,data.FinishDate) == 0 ? 0 : Initiative.GetDaysBetween(data.StartDate,Now)/Initiative.GetDaysBetween(data.StartDate,data.FinishDate))*100;
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
            Completion = 100;
        }else{
            Completion = ProgressCompletion;
        }
    }
    
    return Completion;


}

Initiative.scrollComment = function(){
    $(".box-body").animate({
    scrollTop: $(".direct-chat-messages").height()
    }, 1000);
}
Initiative.GetDataScorecard = function(){
    var url = "/web-cb/scorecard/getdata";
    var parm = {
        Region:"",
        Country:"",
    };
    ajaxPost(url, parm, function(res) {
        if (res.IsError) {
            swal("", res.Message, "info");
        }
        Initiative.ScorecardData(res.Data);
    });
}
Initiative.AddNewLine = function (arr){
    var str =''
    for(var i =0; i< arr.length;i++){
        str+=arr[i]+',';
    }
    var result = '';

    if(str.length < 29){
        return str
    }
    while (str.length > 0) {
        result += str.substring(0, 29) + '\n';
        str = str.substring(29);
    }
    
    return result;

}
Initiative.Linkify = function (inputText) {
    var replacedText, replacePattern1, replacePattern2, replacePattern3;

    //URLs starting with http://, https://, or ftp://
    replacePattern1 = /(\b(https?|ftp):\/\/[-A-Z0-9+&@#\/%?=~_|!:,.;]*[-A-Z0-9+&@#\/%=~_|])/gim;
    replacedText = inputText.replace(replacePattern1, '<a href="$1" target="_blank">$1</a>');

    //URLs starting with "www." (without // before it, or it'd re-link the ones done above).
    replacePattern2 = /(^|[^\/])(www\.[\S]+(\b|$))/gim;
    replacedText = replacedText.replace(replacePattern2, '$1<a href="http://$2" target="_blank">$2</a>');

    //Change email addresses to mailto:: links.
    replacePattern3 = /(([a-zA-Z0-9\-\_\.])+@[a-zA-Z\_]+?(\.[a-zA-Z]{2,6})+)/gim;
    replacedText = replacedText.replace(replacePattern3, '<a href="mailto:$1">$1</a>');

    return replacedText;
}

Initiative.Init = function(){
    Initiative.GetDataScorecard();
    Initiative.GetUserList();
    ajaxPost("/web-cb/businessmetrics/getunassigneddata",{},function(res){
        var businessmetrics = [];
        for(var i in res.Data){
            businessmetrics = businessmetrics.concat(res.Data[i].businessmetric)
        }
        businessmetrics = Enumerable.From(businessmetrics).OrderBy("$.description").ToArray();

        // console.log(businessmetrics)

        hasil = sortStringNumber(100,businessmetrics, 'description');

        Initiative.BusinessMetricList(hasil);
    });
}
$(document).ready(function(){
    Initiative.Init();
})

Initiative.RemoveAttachment = function(a,b,c){
    // console.log(a,b,c)
    var parm = {};
    parm.Initiative_id = a._id;
    parm.InitiativeID = a.InitiativeID
    parm.Attachment = c;
    parm.AttachmentIndex = b;
    console.log(parm)

    swal({
      title: "Are you sure?",
      text: "You will not be able to recover this!",
      type: "warning",
      showCancelButton: true,
      confirmButtonColor: "#DD6B55",
      confirmButtonText: "Yes, delete it!",
      closeOnConfirm: false
    },
    function(){
        Initiative.Processing(true);
        ajaxPost("/web-cb/initiative/removeattacment",parm,function(res){
            Initiative.SaveComplete();
            swal("Deleted!", "Your Attachment has been removed.", "success");
        })
    });
}

Initiative.EditComment = function(a,b,c){
    d = Initiative.SelectedData()
    // console.log(a,b,c)
    d.CommentList[b].Editable = true;

    Initiative.SelectedData(d)

        $("#commenttab li").removeClass("active");$("#commentbox .tab-pane").removeClass("active");
        $($("#commenttab li")[2]).addClass("active");$($("#commentbox .tab-pane")[2]).addClass("active");

        $(".modal-body.data-list .nav.nav-pills li").removeClass("active");$("#commentbox .tab-pane").removeClass("active");
        $($(".modal-body.data-list .nav.nav-pills li")[2]).addClass("active");$($("#commentbox .tab-pane")[2]).addClass("active");

}

Initiative.CancelUpdateComment = function(a,b,c){
    d = Initiative.SelectedData()
    // console.log(a,b,c)
    d.CommentList[b].Editable = false;

    Initiative.SelectedData(d)

        $("#commenttab li").removeClass("active");$("#commentbox .tab-pane").removeClass("active");
        $($("#commenttab li")[2]).addClass("active");$($("#commentbox .tab-pane")[2]).addClass("active");

        $(".modal-body.data-list .nav.nav-pills li").removeClass("active");$("#commentbox .tab-pane").removeClass("active");
        $($(".modal-body.data-list .nav.nav-pills li")[2]).addClass("active");$($("#commentbox .tab-pane")[2]).addClass("active");

}

Initiative.SaveComment = function(a,b,c){
    Initiative.Processing(true);
    var d = ko.mapping.toJS(Initiative.SelectedData());
    d.Id = d._id;
    // console.log(a,b,c)
    tmpNewComment = [];
    for(i in d.CommentList){
        if(i != b){
            tmpNewComment.push(d.CommentList[i])
        }
    }
    c.updated_by = localStorage.getItem("Username");
    c.updated_date = new Date();
    tmpNewComment.push(c)
    d.CommentList = tmpNewComment;

    // Initiative.SelectedData(d)

    // Initiative.Processing(true);
    // formData = Initiative.SelectedData();
    // if(formData.Id === undefined){
    //     formData.Id = formData._id;
    // }
    // formData.CommentList.push({Username: localStorage.getItem("Username"),DateInput: new Date(), Comment: Initiative.TmpComment(), Editable: false})
    // tmpInitiativeData = formData.InitiativeData
    d.InitiativeData = ""
    ajaxPost("/web-cb/dashboard/commentsave",d,function(res){

        Initiative.SaveLogComment("Update", c, d.Id)

        var sources = [];
        Initiative.SaveComplete();;
        Initiative.TmpComment("")
        // $("#commenttab li").removeClass("active");$("#commentbox .tab-pane").removeClass("active");
        // $($("#commenttab li")[2]).addClass("active");$($("#commentbox .tab-pane")[2]).addClass("active");

        // $(".modal-body.data-list .nav.nav-pills li").removeClass("active");$("#commentbox .tab-pane").removeClass("active");
        // $($(".modal-body.data-list .nav.nav-pills li")[2]).addClass("active");$($("#commentbox .tab-pane")[2]).addClass("active");
        setTimeout(Initiative.scrollComment(), 300);
        swal("Update!", "Your Comment has been updated.", "success");
          
    })
}

Initiative.DeleteComment = function(a,b,c){
    swal({
      title: "Are you sure?",
      text: "You will not be able to recover this!",
      type: "warning",
      showCancelButton: true,
      confirmButtonColor: "#DD6B55",
      confirmButtonText: "Yes, delete it!",
      closeOnConfirm: false
    },
    function(){
        Initiative.Processing(true);
        var d = ko.mapping.toJS(Initiative.SelectedData());
        d.Id = d._id;
        // console.log(a,b,c)
        tmpNewComment = [];
        for(i in d.CommentList){
            if(i != b){
                tmpNewComment.push(d.CommentList[i])
            }
        }
        // c.updated_by = localStorage.getItem("Username");
        // c.updated_date = new Date();
        // tmpNewComment.push(c)
        d.CommentList = tmpNewComment;

        // Initiative.SelectedData(d)

        // Initiative.Processing(true);
        // formData = Initiative.SelectedData();
        // if(formData.Id === undefined){
        //     formData.Id = formData._id;
        // }
        // formData.CommentList.push({Username: localStorage.getItem("Username"),DateInput: new Date(), Comment: Initiative.TmpComment(), Editable: false})
        // tmpInitiativeData = formData.InitiativeData
        d.InitiativeData = ""
        ajaxPost("/web-cb/dashboard/commentsave",d,function(res){

            Initiative.SaveLogComment("Delete", c, d.Id)

            var sources = [];
            Initiative.SaveComplete();;
            Initiative.TmpComment("")
            // $("#commenttab li").removeClass("active");$("#commentbox .tab-pane").removeClass("active");
            // $($("#commenttab li")[2]).addClass("active");$($("#commentbox .tab-pane")[2]).addClass("active");

            // $(".modal-body.data-list .nav.nav-pills li").removeClass("active");$("#commentbox .tab-pane").removeClass("active");
            // $($(".modal-body.data-list .nav.nav-pills li")[2]).addClass("active");$($("#commentbox .tab-pane")[2]).addClass("active");
            setTimeout(Initiative.scrollComment(), 300);
            swal("Deleted!", "Your Comment has been remove.", "success");
        })
    });

}

Initiative.SaveLogComment = function(action, comment, initiativeId){
    var d = {};
    d.Id = "";
    d.IdInitiative = initiativeId;
    d.Comment = comment;
    d.Action = action;
    ajaxPost("/web-cb/dashboard/logcommentsave",d,function(res){

    })
}

Initiative.GetLogComment = function(comment){
    // var d = {};
    // console.log(comment)
    Initiative.LogMessage([])
    ajaxPost("/web-cb/dashboard/logcommentget",comment,function(res){
        console.log(res)
        Initiative.LogMessage(res.Data)
        $("#LogComment").modal("show")
    })
}
Initiative.GetUserData = function(){
    ajaxPost("/web-cb/initiativeowner/getdata",{},function(res){
        if(res.IsError){
            swal("", res.Message, "info");
            return false;
        }
        Initiative.InitiativeOwnerList(res.Data);
    })
}
$(document).ready(function(){
    Initiative.GetUserData();
})