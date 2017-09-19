/*
* @Author: Ainur
* @Date:   2016-10-30 11:05:10
* @Last Modified by:   Ainur
* @Last Modified time: 2017-04-07 09:53:26
*/
var Task = {
    Processing:ko.observable(true),
    Mode:ko.observable(""),
    Phase:ko.observable(1),
    TaskType:ko.observable(""),
    // Form Input
    Map:ko.observableArray([]),
    FormValue:ko.observable(),
    RegionList:ko.observableArray([]),
    CountryList:ko.observableArray([]),
    Data:{
        SCCategory:"",
        LifeCycleId:"",
        BusinessDriverId:"",
        EX: false,
        OE: false,
        Name:"",
        Owner:"",
        Statement:"",
        Description:"",
        IsGlobal: true,
        Region: [],
        Country: [],
    },
    SelectedData:ko.observable(),
    LifeCycleList:ko.observableArray([]),
    BusinessDriverList:ko.observableArray([]),
    AllBusinessDriverList:ko.observableArray([]),
    SCCategoryList:ko.observableArray([]),
}

Task.GetCountryDS = function(e){
    // var d = Task.FormValue().Region();
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
    
    Task.CountryList(arr);
    Task.FormValue().Country([])
}

Task.Close = function(){
    $("#Task").modal("hide");
    Task.Reset();
}

Task.Add = function(formValue){
    Task.RegionList(c.RegionList());
    Task.CountryList(c.CountryList());
    Task.GetDataSource();
    $("#Task").modal("show");
    Task.Phase(1);
    if(Task.Mode()===""){
        Task.Mode("new");
    }
    if(formValue.Id !== undefined){
            Task.Map([]);
            Task.Map.push(new MapDataTask({   
            SCCategory: formValue.SCCategory,
            LifeCycle: formValue.LifeCycleId,
            BusinessDriver: formValue.BusinessDriverId,
        }));
        Task.FormValue(ko.mapping.fromJS(formValue));
    }
    else{
        Task.Map([]);
        Task.Map.push(new MapDataTask());
        Task.FormValue(ko.mapping.fromJS(Task.Data));
    }
    Task.Processing(false);
}
Task.Save = function(){
    Task.Processing(false);
    var parm = ko.mapping.toJS(Task.FormValue());
    parm.Map = ko.mapping.toJS(Task.Map());
    parm.TaskType =  Task.TaskType();
    parm.Mode = Task.Mode();
    if(parm.IsGlobal){
       parm.Country = [];
       parm.Region = []; 
    }else{
        for(var i in parm.Country){
            // formData.append("country"+i, parm.Country[i]);
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

        // for(var i in parm.Region){
        //     formData.append("region"+i, parm.Region[i]);
        // }
    }
    ajaxPost("/web-cb/task/save",parm,function(res){
        if(res.IsError){
            swal("", res.Message, "info");
            return false;
        }        
        Task.Close();
    })
}

function MapDataTask(data){
    console.log(data)
    var self = this;
    if(data!==undefined){
        var selectedBD = [];
        if(data.SCCategory===null || data.SCCategory===""){
            selectedBD = Task.AllBusinessDriverList();
        }else{
            selectedBD = Enumerable.From(Task.AllBusinessDriverList()).Where("$.Parentid==='"+data.SCCategory+"'").ToArray();
        }
        self = {LifeCycle:ko.observable(data.LifeCycle),BDList:ko.observableArray(selectedBD),SCCategory:ko.observable(data.SCCategory),BusinessDriver:ko.observable(data.BusinessDriver)};
    }else{
        var selectedBD = Task.AllBusinessDriverList();
        self = {LifeCycle:ko.observable(""),BDList:ko.observableArray(selectedBD),SCCategory:ko.observable(""),BusinessDriver:ko.observable("")};    
    }
    

    self.SCCategory.subscribe(function(val){
        console.log(val);
        if(val===null || val===""){
            self.BDList(Initiative.AllBusinessDriverList());
        }else{
            var selectedBD = Enumerable.From(Task.AllBusinessDriverList()).Where("$.Parentid==='"+val+"'").ToArray();
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

Task.NewMap = function(){
    Task.Map.push(new MapDataTask());
}

Task.GetDataSource = function(){
    var activeTab = $("#dashboard .nav-pills li.active a").attr("href").replace("#","");
    Task.TaskType(activeTab);
    var sources = c.DataSource().Data;
    Task.LifeCycleList(sources.MasterLifeCycle);
    Task.BusinessDriverList(sources.SummaryBusinessDriver);
    Task.AllBusinessDriverList(sources.AllSummaryBusinessDriver);
    Task.SCCategoryList(Scorecard.Data());
}
Task.Reset = function(){
    Task.Map([]);
    Task.Map.push(new MapDataTask());
    Task.FormValue(ko.mapping.fromJS(Task.Data));
}

Task.Close = function(){
    $("#Task").modal("hide");
    Task.Reset();
    Task.Mode("");
    c.GetData();
}

Task.Get = function(Id){
    var sources = c.DataSource().Data.TaskList;
    var d = Enumerable.From(sources).Where("$.Id === '"+Id+"'").FirstOrDefault();
    Task.Mode("");
    Task.SelectedData(d);
    $("#Task").modal("show");
}

// Task.Remove = function(){
//     var d = ko.mapping.toJS(Task.SelectedData());
//     ajaxPost("/task/remove",{Id:d.Id},function(res){
//         if(res.IsError){
//             swal("", res.message, "info");
//             return false;
//         }
//         Task.Close();
//     });
// }

Task.Remove = function(){
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
    var d = ko.mapping.toJS(Task.SelectedData());
    ajaxPost("/web-cb/task/remove",{Id:d.Id},function(res){
        if(res.IsError){
            swal("", res.message, "info");
            return false;
        }
        swal("Deleted!", "Your Task file has been deleted.", "success");
        Task.Close();
    });     
    });
}
Task.Edit = function(){
    Task.Mode("edit")
    var d = ko.mapping.toJS(Task.SelectedData());
    Task.Add(d);
}
// =====