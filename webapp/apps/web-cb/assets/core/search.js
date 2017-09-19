/*
 * @Author: Ainur
 * @Date:   2016-11-17 13:40:31
 * @Last Modified by:   Ainur
 * @Last Modified time: 2017-04-07 09:53:11
 */

 function debounce(func, wait, immediate) {
  var timeout;
  return function() {
    var context = this, args = arguments;
    var later = function() {
      timeout = null;
      if (!immediate) func.apply(context, args);
    };
    var callNow = immediate && !timeout;
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
    if (callNow) func.apply(context, args);
  };
};

var search = {
  SelectedItem: ko.observable(""),
  Processing: ko.observable(false),
  Data: ko.observableArray([]),
}
search.GetData = function (keyword) {
//  console.log("SEARCH");
  var url = "/web-cb/search/getresult";
  var parm = {};
  if (keyword === undefined) {
    parm["keyword"] = search.SelectedItem().Keyword;
    Initiative.CurrentSearchKeyword('');
  } else {
    parm["keyword"] = keyword;
    Initiative.CurrentSearchKeyword(keyword);
  }

  parm["keywordEscaped"] = parm["keyword"]
    .replace(/\(/g, "\\(")
    .replace(/\)/g, "\\)")
    .replace(/\+/g, "\\+")
    .replace(/\-/g, "\\-")
    .replace(/\//g, "\\/")
    .replace(/\./g, "\\.")
    .replace(/\&/g, "\\&")
    .replace(/\,/g, "\\,")

  parm["InitiativeType"] = c.SelectedTab();
  var Initiatives = []
  if(c.DataSource() != undefined && c.DataSource().Data != undefined && c.DataSource().Data.Project != undefined){
    Initiatives = Enumerable.From(c.DataSource().Data.Project).Select("$.InitiativeID").ToArray();
  }
  parm["Initiatives"] = Initiatives
  ajaxPost(url, parm, function (res) {
    if (res.Data.ProjectName.length === 0 && res.Data.ProjectManager.length === 0 &&
            res.Data.AccountExecutive.length === 0 && res.Data.TechnologyLead.length === 0 &&
            res.Data.ProblemStatement.length === 0 && res.Data.ProjectDescription.length === 0 &&
            res.Data.Filenames.length === 0 && res.Data.TaskName.length === 0 &&
            res.Data.TaskOwner.length === 0 && res.Data.TaskStatement.length === 0 &&
            res.Data.TaskDesc.length === 0) {
      swal("Error!", "No Matching on choosen keyword.", "error");
      return false;
    } else {
      $("#search-result").modal("show");
      search.Render(res.Data);
    }
  });
}
search.Render = function (dataSource) {
//  console.log(dataSource);
  $("#sr-data").html("");
  $("#sr-data").kendoPanelBar({
    expandMode: "multiple"
  });
  var dataArr = [];
  var pnList = [];
  var pmList = [];
  var aeList = [];
  var tlList = [];
  var psList = [];
  var pdList = [];
  var fnList = [];

  var tnList = [];
  var toList = [];
  var tsList = [];
  var tdList = [];

  if (dataSource.ProjectName.length > 0) {
    $.each(dataSource.ProjectName, function (idx, obj) {
      pnList.push({"encoded": false, "text": "<a href='#' onclick='Initiative.Get(\"" + obj.Id + "\");'> " + obj.Keyword + "</a>", "value": obj.Id});
    });
    dataArr.push({text: "<b>Project Name ( " + dataSource.ProjectName.length + " )</b>", value: "VOID", encoded: false, items: pnList});
  }

  if (dataSource.ProjectManager.length > 0) {

    $.each(dataSource.ProjectManager, function (idx, obj) {
//      pmList.push({"encoded": false, "text": "<a href='#' onclick='Initiative.Get(\"" + obj.Id + "\");'> " + obj.Keyword + "</a>" + " - " + obj.ParentInfo, "value": obj.Id});
      var textInsider = "<table style='width:100%'> <tbody><tr><td style='width:50%; vertical-align:top;'><a href='#' onclick=\"Initiative.Get('" + obj.Id + "');\">" + obj.Keyword + "</a></td><td style='width:50%'>" + obj.ParentInfo + "</td></tr></tbody></table>";
      pmList.push({"encoded": false, "text": textInsider, "value": obj.Id});
    });
    dataArr.push({text: "<b>Project Manager ( " + dataSource.ProjectManager.length + " )</b>", value: "VOID", encoded: false, items: pmList});
  }
  if (dataSource.AccountExecutive.length > 0) {
    $.each(dataSource.AccountExecutive, function (idx, obj) {
      var textInsider = "<table style='width:100%'> <tbody><tr><td style='width:50%; vertical-align:top;'><a href='#' onclick=\"Initiative.Get('" + obj.Id + "');\">" + obj.Keyword + "</a></td><td style='width:50%'>" + obj.ParentInfo + "</td></tr></tbody></table>";
      aeList.push({"encoded": false, "text": textInsider, "value": obj.Id});
    });
    dataArr.push({text: "<b>Account Executive ( " + dataSource.AccountExecutive.length + " )</b>", value: "VOID", encoded: false, items: aeList});
  }
  if (dataSource.TechnologyLead.length > 0) {
    $.each(dataSource.TechnologyLead, function (idx, obj) {
      var textInsider = "<table style='width:100%'> <tbody><tr><td style='width:50%; vertical-align:top;'><a href='#' onclick=\"Initiative.Get('" + obj.Id + "');\">" + obj.Keyword + "</a></td><td style='width:50%'>" + obj.ParentInfo + "</td></tr></tbody></table>";
      tlList.push({"encoded": false, "text": textInsider, "value": obj.Id});
    });
    dataArr.push({text: "<b>Technology Lead ( " + dataSource.TechnologyLead.length + " )</b>", value: "VOID", encoded: false, items: tlList});
  }
  if (dataSource.ProblemStatement.length > 0) {
    $.each(dataSource.ProblemStatement, function (idx, obj) {
      var textInsider = "<table style='width:100%'> <tbody><tr><td style='width:50%; vertical-align:top;'><a href='#' onclick=\"Initiative.Get('" + obj.Id + "');\">" + obj.Keyword + "</a></td><td style='width:50%'>" + obj.ParentInfo + "</td></tr></tbody></table>";
      psList.push({"encoded": false, "text": textInsider, "value": obj.Id});
    });
    dataArr.push({text: "<b>Problem Statement ( " + dataSource.ProblemStatement.length + " )</b>", value: "VOID", encoded: false, items: psList});
  }
  if (dataSource.ProjectDescription.length > 0) {
    $.each(dataSource.ProjectDescription, function (idx, obj) {
      var textInsider = "<table style='width:100%'> <tbody><tr><td style='width:50%; vertical-align:top;'><a href='#' onclick=\"Initiative.Get('" + obj.Id + "');\">" + obj.Keyword + "</a></td><td style='width:50%'>" + obj.ParentInfo + "</td></tr></tbody></table>";
      pdList.push({"encoded": false, "text": textInsider, "value": obj.Id});
    });
    dataArr.push({text: "<b>Project Description ( " + dataSource.ProjectDescription.length + " )</b>", value: "VOID", encoded: false, items: pdList});
  }
  if (dataSource.Filenames.length > 0) {
    $.each(dataSource.Filenames, function (idx, obj) {
      var textInsider = "<table style='width:100%'> <tbody><tr><td style='width:50%; vertical-align:top;'><a href='#' onclick=\"Initiative.Get('" + obj.Id + "');\">" + obj.Keyword + "</a></td><td style='width:50%'>" + obj.ParentInfo + "</td></tr></tbody></table>";
      fnList.push({"encoded": false, "text": textInsider, "value": obj.Id});
    });
    dataArr.push({text: "<b>Attachments ( " + dataSource.Filenames.length + " )</b>", value: "VOID", encoded: false, items: fnList});
  }
  //Task

  if (dataSource.TaskName.length > 0) {
    $.each(dataSource.TaskName, function (idx, obj) {
      var textInsider = "<table style='width:100%'> <tbody><tr><td style='width:50%; vertical-align:top;'><a href='#' onclick=\"Task.Get('" + obj.Id + "');\">" + obj.Keyword + "</a></td><td style='width:50%'>" + obj.ParentInfo + "</td></tr></tbody></table>";
      tnList.push({"encoded": false, "text": textInsider, "value": obj.Id});
    });
    dataArr.push({text: "<b>Task Name ( " + dataSource.TaskName.length + " )</b>", value: "VOID", encoded: false, items: tnList});
  }

  if (dataSource.TaskOwner.length > 0) {
    $.each(dataSource.TaskOwner, function (idx, obj) {
      var textInsider = "<table style='width:100%'> <tbody><tr><td style='width:50%; vertical-align:top;'><a href='#' onclick=\"Task.Get('" + obj.Id + "');\">" + obj.Keyword + "</a></td><td style='width:50%'>" + obj.ParentInfo + "</td></tr></tbody></table>";
      toList.push({"encoded": false, "text": textInsider, "value": obj.Id});
    });
    dataArr.push({text: "<b>Task Owner ( " + dataSource.TaskOwner.length + " )</b>", value: "VOID", encoded: false, items: toList});
  }

  if (dataSource.TaskStatement.length > 0) {
    $.each(dataSource.TaskStatement, function (idx, obj) {
      var textInsider = "<table style='width:100%'> <tbody><tr><td style='width:50%; vertical-align:top;'><a href='#' onclick=\"Task.Get('" + obj.Id + "');\">" + obj.Keyword + "</a></td><td style='width:50%'>" + obj.ParentInfo + "</td></tr></tbody></table>";
      tsList.push({"encoded": false, "text": textInsider, "value": obj.Id});
    });
    dataArr.push({text: "<b>Task Statement ( " + dataSource.TaskStatement.length + " )</b>", value: "VOID", encoded: false, items: tsList});
  }


  if (dataSource.TaskDesc.length > 0) {
    $.each(dataSource.TaskDesc, function (idx, obj) {
      var textInsider = "<table style='width:100%'> <tbody><tr><td style='width:50%; vertical-align:top;'><a href='#' onclick=\"Task.Get('" + obj.Id + "');\">" + obj.Keyword + "</a></td><td style='width:50%'>" + obj.ParentInfo + "</td></tr></tbody></table>";
      tdList.push({"encoded": false, "text": textInsider, "value": obj.Id});
    });
    dataArr.push({text: "<b>Task Description ( " + dataSource.TaskDesc.length + " )</b>", value: "VOID", encoded: false, items: tdList});
  }

  var srdata = $("#sr-data").data("kendoPanelBar");
  srdata.append(dataArr);
};

search.onSelect = function (e) {
//  console.log(e);
//  var item = e.item;
//  index = item.parentsUntil(".k-panelbar", ".k-item").map(function () {
//    return $(this).index();
//  }).get().reverse();
//  console.log(index);
};

// search.FilterInitiative = function (keyword) {
//   keyword = keyword.toLowerCase()
//   var data = c.DataSource().Data
//   var newTableSourcesVer2 = ko.mapping.toJS(data.TableSourcesVer2Backup);
//   var newInitiative = []
//   newTableSourcesVer2.forEach(function (d) {
//     d.Initiatives = d.Initiatives.filter(function (e) {
//       var projectManager = (e.hasOwnProperty('ProjectManager') 
//           ? e.ProjectManager.toLowerCase() : "")
//       var accountableExecutive = (e.hasOwnProperty('AccountableExecutive') 
//           ? e.AccountableExecutive.toLowerCase() : "")
//       var projectName = (e.hasOwnProperty('ProjectName') 
//           ? e.ProjectName.toLowerCase() : "")

//       var isCond1Match = projectManager.indexOf(keyword) > -1
//       var isCond2Match = accountableExecutive.indexOf(keyword) > -1
//       var isCond3Match = projectName.indexOf(keyword) > -1
//       return isCond1Match || isCond2Match || isCond3Match
//     })
//     if(keyword != ''){
//         if(d.Initiatives.length != 0){
//           d.Initiatives.forEach(function (data){
//             newInitiative.push(data)    
//           })
//         }
//     }else{
//         newInitiative = c.AllInitiateSource()
//     }

//   })
//   data.Project = newInitiative
//   data.TableSourcesVer2 = ko.mapping.fromJS(newTableSourcesVer2)()
//   c.DataSource({
//     Data: data
//   })
// }

search.FilterInitiative = function (keyword) {
  // console.log("keyword : ", keyword)
  keyword = keyword.toLowerCase()
  var data = c.DataSource().Data
  
  var newTableSourcesVer3 = SortInitiative.Active() ? data.TableSourcesVer3AlignVer : data.TableSourcesVer3BackupAll;
  var dataTableSourcesVer3 = [];
  if(keyword != ""){
    _.each(newTableSourcesVer3, function(v,i){
      var newTableSourcesVer2 = ko.mapping.toJS(newTableSourcesVer3[i].TableSourcesVer2);
      var newInitiative = []
      newTableSourcesVer2.forEach(function (d) {
        d.Initiatives = d.Initiatives.filter(function (e) {
          var projectManager = (e.hasOwnProperty('ProjectManager') 
              ? e.ProjectManager.join().toLowerCase() : "")
          var accountableExecutive = (e.hasOwnProperty('AccountableExecutive') 
              ? e.AccountableExecutive.join().toLowerCase() : "")
          var projectName = (e.hasOwnProperty('ProjectName') 
              ? e.ProjectName.toLowerCase() : "")
          var techLead = (e.hasOwnProperty('TechnologyLead') 
              ? e.TechnologyLead.join().toLowerCase() : "")
          var Owner = (e.hasOwnProperty('Owner') 
              ? e.Owner.toLowerCase() : "")
          var Statement = (e.hasOwnProperty('Statement') 
              ? e.Statement.toLowerCase() : "")
          var Description = (e.hasOwnProperty('Description') 
              ? e.Description.toLowerCase() : "")

          var isCond1Match = projectManager.indexOf(keyword) > -1
          var isCond2Match = accountableExecutive.indexOf(keyword) > -1
          var isCond3Match = projectName.indexOf(keyword) > -1
          var isCond4Match = techLead.indexOf(keyword) > -1
          var isCond5Match = Owner.indexOf(keyword) > -1
          var isCond6Match = Statement.indexOf(keyword) > -1
          var isCond7Match = Description.indexOf(keyword) > -1
          return isCond1Match || isCond2Match || isCond3Match || isCond4Match || isCond5Match || isCond6Match || isCond7Match
        })
        if(keyword != ''){
            if(d.Initiatives.length != 0){
              d.Initiatives.forEach(function (data){
                newInitiative.push(data)
                // console.log("--->", data, data.ProjectName)   
              })
            }
        }else{
            newInitiative = c.AllInitiateSource()
        }
      })
      data.Project = newInitiative

      TableSourcesVer2Mapping = ko.mapping.fromJS(newTableSourcesVer2)
      var Id = v.Id
      var BDIdDefault = v.BDIdDefault
      // console.log(Id, BDIdDefault)
      dataTableSourcesVer3.push({"Id":Id,"BDIdDefault":BDIdDefault,"TableSourcesVer2": TableSourcesVer2Mapping})
    })
  } else{
    dataTableSourcesVer3 = newTableSourcesVer3;
  }
  data.TableSourcesVer3(dataTableSourcesVer3)
  // data.TableSourcesVer3(newTableSourcesVer3)

  // console.log(data.TableSourcesVer3())

  c.DataSource({
    Data: data
  })
  // setTimeout(function() {
  //     SortInitiative.sycHeight();
  // }, 1000);
}

$(document).ready(function () {
  setTimeout(function () {
    $("#search-input").on('keyup', debounce(function (e) {
      var keyword = $(this).val()
      if (e.keyCode == 13) {
        search.GetData(keyword.trim());
      }
      Initiative.CurrentSearchKeyword(keyword);
      search.FilterInitiative(keyword)
      redipsInit();
    }, 300))
  }, 200)

  $("#search-input").kendoAutoComplete({
    filter: "contains",
    minLength: 2,
    select: function (e) {
      var dataItem = this.dataItem(e.item.index());
      search.SelectedItem(dataItem);
      search.GetData();
    },
    change: function (e) {
      if (e.type == 'keydown') {
        if (e.keyCode == 13) {
          var val = e.sender.value().trim();
          search.GetData(val);
        }
      }
    },
    dataTextField: 'Keyword',
    dataValueField: 'Keyword',
    dataSource: {
      serverFiltering: true,
      transport: {
        read: function (e) {
          var parm = {};
          parm["keyword"] = e.data.filter.filters[0].value;
          parm["InitiativeType"] = c.SelectedTab();
          var Initiatives = []
          if(c.DataSource() != undefined && c.DataSource().Data != undefined && c.DataSource().Data.Project != undefined){
            Initiatives = Enumerable.From(c.DataSource().Data.Project).Select("$.InitiativeID").ToArray();
          }
          parm["Initiatives"] = Initiatives
          if (parm.keyword.trim() === "" && parm.keyword.length < 2) {
            e.success([]);
            return false;
          }
          // For Testing
//          var d = [{displayfield: "Test", rujak: "enak"}, {displayfield: "Test2", rujak: "asin"}];
//          e.success(d);
//          return false;

          var url = "/web-cb/search/autocomplete";
          ajaxPost(url, parm, function (res) {
            if (res === "") {
              e.success([]);
              return false;
            } else if (typeof res === "string" && res.indexOf("html") >= 0) {
              c.Logout();
            } else if (res.isError) {
              e.success([]);
              return false;
            }
            e.success(res.Data);
          }, function (e) {
            alert(e);
          });
        }
      }
    },
    // placeholder: "🔍 Search..",
    // placeholder:"&#xf002; Search...",
  });
});