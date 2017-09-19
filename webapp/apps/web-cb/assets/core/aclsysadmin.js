var adm = {}
adm.searchGrid = ko.observable(false);
adm.addUser = ko.observable(false);
adm.saveUser = ko.observable(false);
adm.updateUser = ko.observable(false);
adm.addGroup = ko.observable(false);
adm.listUser = ko.observableArray([]);
adm.updateGroup = ko.observable(false);
adm.listGroup = ko.observableArray([]);
adm.listSession = ko.observableArray([]);
adm.nameFilter = ko.observable("");
adm.dataChangePassword = ko.observable({oldPassword: '', newPassword: '', reNewPassword: ''});
adm.onSearch = ko.observable("");
adm.addNewUser = ko.observable({
  loginid: "", 
  fullname: "", 
  email: "", 
  password: "", 
  rePassword:"", 
  groups: [], 
  enable: "",
});
adm.btnSaveGroup = ko.observable(false);
adm.FormValue=ko.observable();
adm.UserFormValue=ko.observable();
adm.listMenu  = ko.observableArray();
adm.ListMenuTree  = ko.observableArray([]);
adm.treelistView  = ko.observableArray();
adm.UserName = ko.observable("");
adm.Fullname = ko.observable("");
adm.Email = ko.observable("");
adm.Password = ko.observable("");
adm.RePassword = ko.observable("");
adm.valueGroup = ko.observableArray([]);
adm.valueEnable = ko.observable('true');
adm.Id = ko.observable("");
adm.LastPassword = ko.observable("");
adm.Lastphoto = ko.observable("");
adm.Data = {
  ID : ko.observable(""),
  Title : ko.observable(""),
  ParentId : ko.observable(""),
  Url : ko.observable(""),
  Index : ko.observable(""),
  Group1 :ko.observable(""),
  Group2 :ko.observable(""),
  Group3 :ko.observable(""),
  Enable : ko.observable(""),
  Category :ko.observable(""),
}
adm.btnAddMenu = ko.observable(true);
adm.btnSaveMenu = ko.observable(false);
adm.btnEditMenu = ko.observable(false);
adm.btnDeleteMenu = ko.observable(false);
adm.isNewUser = ko.observable(false);

adm.sourceDataMenu = ko.observable();
var grp = {
  Id          : ko.observable(""),
  Name        : ko.observable(""),
  valueEnable : ko.observable(""),
  groupType   : ko.observable(0),
  listRole    : ko.observable("")
};

grp.reset = function(){
  grp.Id("");
  grp.Name("");
  grp.valueEnable("");
  grp.groupType(0);
  grp.listRole();
  grp.checked(false);
};

grp.checked = ko.observable(false);
$('#enablegroup').prop('checked', true);

$(function () {
  var container = $("#formAddGroup");
  kendo.init(container);
  container.kendoValidator({
    rules: {
      checkalphabet: function (input) {
        if (input.is("[data-checkalphabet-msg]") && input.val() != "") {
          if(!/^[_a-zA-Z]*$/g.test(input.val())) {
            return false;
          }
        }

        return true;
      }
    }
  });

});

grp.doSaveGroup = function(){
  var validator = $("#formAddGroup").data("kendoValidator")

  if (validator.validate()) {
    adm.updateGroup(false);
    var roles = $('input[name=role]:checked').map(function(){
      return $(this).val();
    }).get();
                      
    grp.listRole(roles.join(","));

    if(grp.checked()){
      grp.valueEnable("true")
    }else{
      grp.valueEnable("false")
    }

    var Id = grp.Id();
    grp.Id(Id.toUpperCase());

    ajaxPost('/web-cb/aclsysadmin/savegroupuser', grp, function(data){
      swal("Success","Insert/Update Data", "success");
      adm.createGridGroupAdministrator();
      grp.reset();
      adm.addGroup(false);
      adm.btnSaveGroup(false);
      // $("#addusernew").show();
    });
  }
}

grp.showEdit = function(value){
  adm.addGroup(true);
  adm.btnSaveGroup(false);
  adm.updateGroup(true);
  var param = {
    id: value
  };

  ajaxPost("/web-cb/aclsysadmin/getgroupbyid", param, function(data){
    grp.Id(data._id);   
    grp.Name(data.Title);
    grp.groupType(data.GroupType);
    grp.checked(data.Enable);

    checkboxes = document.getElementsByName('role');
    for(var i=0; i < data.Grants.length; i++) {
      document.getElementById("check-all-new" + data.Grants[i].AccessID).checked = true;
    }
  });
}

grp.doDelete = function(value){
 var param = {
    id:value
  };
    var url = "/web-cb/aclsysadmin/deletegroup";
    swal({
            title: "Are you sure?",
            text: "Are you sure delete this group!",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: '#DD6B55',
            confirmButtonText: 'Yes, I am sure!',
            cancelButtonText: "No, cancel it!",
            closeOnConfirm: false,
            closeOnCancel: false
        },
        function(isConfirm) {
            if (isConfirm) {
                ajaxPost(url, param, function(data){
                    if (data.success == true){
                        swal("Success!",data.Message,"success");
                         adm.createGridGroupAdministrator();
                    }else{
                        swal("Error!",data.Message,"error");
                    }
                });
            } else {
                swal("Cancelled", "Cancelled Delete Group", "error");
            }
        });
}

grp.getRole = function() {
  var param = {};
  var url = "";
  grp.listRole([]);
  ajaxPost(url, param, function(res){
    var dataRole = Enumerable.From(res).OrderBy("$.name").ToArray();
    for (var r in dataRole){
      grp.listRole.push({
        "text": dataRole[r].name,
        "value": dataRole[r].name,
      })
    }
  })
}

grp.getTopMenu = function(){
  var param = {};
  var url = "/web-cb/aclsysadmin/getdataaccess";
  ajaxPost(url, param, function(res){
    if(res.IsError != true){
      var dataMenu = res.Data;
      var newRecords = [];
      for (var d in dataMenu){
        newRecords.push({
          "Id": dataMenu[d].id,
          "IndexMenu": dataMenu[d].index,
          "PageId": dataMenu[d]._id,
          "Parent": dataMenu[d].parentid,
          "Title": dataMenu[d].title,
          "Checkall": false,
          "Access": false,
          "Create": false,
          "Edit": false,
          "Delete": false,
          "View": false,
          "Approve": false,
          "Process": false,
        });
      }
      grp.GetDataMenu(newRecords);
    }else{
      return swal("Error!", res.Message, "error");
    }
  });
}

grp.GetDataMenu = function(e) {
  var myData = new kendo.data.TreeListDataSource({
    data: e,
    schema: {
      model: {
        id: "Id",
        parentId: "Parent",
        fields: {
          Id: {field: "Id", type: "string"},
          Title: {field: "Title", type: "string"},
          parentId: {field: "Parent", type: "string"}
        },
        expanded: false,
      }
    }
  });
  if ($("#MasterGridMenu").data("kendoTreeList") !== undefined) {
    $("#MasterGridMenu").data("kendoTreeList").setDataSource(myData);
    return;
  }
  $("#MasterGridMenu").kendoTreeList({
    dataSource: myData,

    columns: [
      { 
        field: "Title",
        title:"Available Menu", 
        width: 200 
      },
      { 
        field:"Checkall",
        title:"Show", 
        width: 50,
        attributes:{"class": "align-center"},
        template: "<input id='check-all-new#:id#' name='role' value='#:id#' class='rolecheck-value-check-all' type='checkbox' onclick='grp.Checkall(#:id#)' #: Checkall==true ? 'checked' : '' #/>"
      },
      {
        field:"Access",
        title:"Access",
        width:50,
        attributes: {"class": "align-center"},
        // template:"#if(parentId != '' || Haschild == false){#<input id='check-Access-#:Id #' class='rolecheck-value-Access' onclick='grp.unCheck(#:Id#)' type='checkbox' #: Access==true ? 'checked' : '' #/>#}#"              
        template:"<input id='check-Access-#:id #' name='access' class='rolecheck-value-Access' onclick='grp.checkAccess(\"#: Id #\")' type='checkbox' #: Access==true ? 'checked' : '' #/>"
      },
      {
        field:"Create",
        title:"Create",
        width:50,
        attributes: {"class": "align-center"},
        template:"<input id='check-Create-#:id #' name='create' class='rolecheck-value-Create' onclick='grp.unCheck(#:id#)' type='checkbox' #: Create==true ? 'checked' : '' #/>"              
      },
      {
        field:"Edit",
        title:"Edit",
        width:50,
        attributes: {"class": "align-center"},
        template:"<input id='check-Edit-#:id #' name='edit' class='rolecheck-value-Edit' onclick='grp.unCheck(#:id#)' type='checkbox' #: Edit==true ? 'checked' : '' #/>"  
      },
      {
        field:"Delete",
        title:"Delete",
        width:50,
        attributes: {"class": "align-center"},
        template:"<input id='check-Delete-#:id #' name='delete' class='rolecheck-value-Delete' onclick='grp.unCheck(#:id#)' type='checkbox' #: Delete==true ? 'checked' : '' #/>"
      },
      {
        field:"View",
        title:"View",
        width:50,
        attributes: {"class": "align-center"},
        template:"<input id='check-View-#:id #' name='view' class='rolecheck-value-View' onclick='grp.unCheck(#:id#)' type='checkbox' #: View==true ? 'checked' : '' #/>"
      }
    ],
  }); 
}

grp.unCheck = function(Menuid){
  if(!$("#check-Access-"+Menuid).prop('checked') 
      || !$("#check-Create-"+Menuid).prop('checked') 
      || !$("#check-Edit-"+Menuid).prop('checked') 
      || !$("#check-Delete-"+Menuid).prop('checked') 
      || !$("#check-View-"+Menuid).prop('checked') 
      || !$("#check-Approve-"+Menuid).prop('checked') 
      || !$("#check-Process-"+Menuid).prop('checked')){
          $('#check-all'+Menuid).prop('checked', false);
          $('#check-all-new'+Menuid).prop('checked', false);
  }else if($("#check-Access-"+Menuid).prop('checked') 
      == true && $("#check-Create-"+Menuid).prop('checked')
      == true && $("#check-Edit-"+Menuid).prop('checked') 
      == true && $("#check-Delete-"+Menuid).prop('checked') 
      == true && $("#check-View-"+Menuid).prop('checked') 
      == true && $("#check-Approve-"+Menuid).prop('checked') 
      == true && $("#check-Process-"+Menuid).prop('checked') 
      == true){
          $('#check-all'+Menuid).prop('checked', true);
          $('#check-all-new'+Menuid).prop('checked', true); 
  }
}

grp.Checkall = function(Menuid){
  // $('#check-all'+Menuid).change(function(){
  //     $("#check-Access-"+Menuid).prop('checked', $(this).prop('checked'));
  //     $("#check-Create-"+Menuid).prop('checked', $(this).prop('checked'));
  //     $("#check-Edit-"+Menuid).prop('checked', $(this).prop('checked'));
  //     $("#check-Delete-"+Menuid).prop('checked', $(this).prop('checked'));
  //     $("#check-View-"+Menuid).prop('checked', $(this).prop('checked'));
  //     $("#check-Approve-"+Menuid).prop('checked', $(this).prop('checked'));
  //     $("#check-Process-"+Menuid).prop('checked', $(this).prop('checked'));
  // });
  $('#check-all-new'+Menuid).change(function(){
      $("#check-Access-"+Menuid).prop('checked', $(this).prop('checked'));
      $("#check-Create-"+Menuid).prop('checked', $(this).prop('checked'));
      $("#check-Edit-"+Menuid).prop('checked', $(this).prop('checked'));
      $("#check-Delete-"+Menuid).prop('checked', $(this).prop('checked'));
      $("#check-View-"+Menuid).prop('checked', $(this).prop('checked'));
      $("#check-Approve-"+Menuid).prop('checked', $(this).prop('checked'));
      $("#check-Process-"+Menuid).prop('checked', $(this).prop('checked'));
  });
}

adm.reset = function(){
  adm.UserName("");
  adm.Fullname("");
  adm.Email("");
  adm.Password("");
  adm.RePassword("");
  adm.valueGroup([]);
  adm.Id("");
  $('#enableuser').prop('checked', false);
  $("#img-ex").attr("src","/static/img/noimage.jpg");
  $("#imgval").val("");
}

adm.addFormGroup = function() {
  adm.addGroup(true);
  adm.btnSaveGroup(true);
}

adm.createGridUserAdministration = function(){
 // var param = { Userid: $("#searchUser").data("kendoMultiSelect").value()};
  var url = "/web-cb/aclsysadmin/getdatauser";  
  var param = {};
  ajaxPost(url, param, function (datas){
    for(var i in datas.Data){
      datas.Data[i].user_group = datas.Data[i].groups.join(", ");
    }
    adm.listUser(datas.Data);
    $("#gridUserAdministrator").html("");
    $("#gridUserAdministrator").kendoGrid({
      dataSource: {
        data: datas.Data,
        pageSize : 10 
        },
        columns: [
        {
          field:"loginid",
          title:"Login ID",
          width:150,
        },
        {
          field:"fullname",
          title:"Full Name",
          width:150,
        },
        {
          field:"email",
          title:"E-mail",
          width:150,
        },
        {
          field:"user_group",
          title:"Groups",
        },
        {
          field:"enable",
          title:"Status",
          width:100,
          template:function(d){
            if(d.enable){
              return "<span class=\'label label-success status-info\' style=\'font-size: 10px;font-weight: normal;text-transform: uppercase;\'>enabled</span>";
            }else{
              return "<span class=\'label label-default status-info\' style=\'font-size: 10px;font-weight: normal;text-transform: uppercase;\'>disabled</span>";
            }
          },
          attributes:{class:'text-center'}
        },
        {
          title:"",
          width: 80,
          attributes:{class:'text-center'},
          template:"<button onclick='adm.editData(\"#:id #\")' class='btn btn-xs btn-warning'><i class='fa fa-pencil'></i></button> <button onclick='adm.doDelete(\"#:id #\")' class='btn btn-xs btn-danger'><i class='fa fa-trash'></i></button>"
        }],
        filterable:true,
      sortable: true,
      pageable: {
          refresh: true,
          pageSizes: true,
          buttonCount: 5
      },
      // toolbar: ["excel"],
      excel: {
        fileName: "User List Name.xlsx",
        allPages: true
      },
      // height: 380,
    });
  }) 


}

adm.createGridGroupAdministrator = function(){
  var url = "/web-cb/aclsysadmin/getdatagroup"; 
  ajaxPost(url, {}, function (datas){
    adm.listGroup(datas.Data)
  });

  setTimeout(function() {
    $("#gridGroupAdministrator").html("");
    $("#gridGroupAdministrator").kendoGrid({
      dataSource: {
        data: adm.listGroup(),
        pageSize : 10 
        },
        columns: [
        {
          field:"id",
          title:'ID'
        },
        {
          field:"title",
          title:'Title'
        },
        {
          field:"enable",
          title:'Enable'
        },
        {
          field:"grouptype",
          title:'Group type'
        },
        {
          title:"Action",
          width: 120,
          template:"<button onclick='grp.showEdit(\"#:id #\")' class='btn btn-xs btn-warning'><i class='fa fa-pencil'></i></button> <button onclick='grp.doDelete(\"#:id #\")' class='btn btn-xs btn-danger'><i class='fa fa-trash'></i></button>"
        }],
      sortable: true,
      pageable: {
          refresh: true,
          pageSizes: true,
          buttonCount: 5
      },
      // toolbar: ["excel"],
      excel: {
        fileName: "Group List Name.xlsx",
        allPages: true
      },
      // height: 380, 
    });
  }, 1000);
  
}

adm.createGridAccessAdministrator = function(){
  var url = '/web-cb/aclsysadmin/getdataaccess'; 

  ajaxPost(url, {}, function (datas){
    setTimeout(function() {
    $("#gridAccessAdministrator").html("");
    $("#gridAccessAdministrator").kendoGrid({
      dataSource:{
        data:datas.Data,
        pageSize : 10 ,
      }, 
      columns: [
        {
          field:"id",
          title:'ID'
        },
        {
          field:"title",
          title:'Title'
        },
        {
          field:"group1",
          title:'Group 1'
        },
        {
          field:"group2",
          title:'Group 2'
        },
        {
          field:"group3",
          title:'Group 3'
        },
          ],
      sortable: true,
      pageable: {
          refresh: true,
          pageSizes: true,
          buttonCount: 5
      },
      height: 380,
    });

  }, 1000); 
  });


}

adm.createGridSessionAdministrator = function(){
  var url = "/web-cb/aclsysadmin/getdatasession"; 

  ajaxPost(url, {}, function (datas){
    var dataSource = datas.Data;
    for(var i in dataSource){
      dataSource[i].created = getUTCDateFull(dataSource[i].created);
      dataSource[i].expired = getUTCDateFull(dataSource[i].expired);
    }
    $("#gridSessionAdministrator").html("");
    $("#gridSessionAdministrator").kendoGrid({
      dataSource: {
       data: dataSource,
       pageSize : 10  
      },
      columns: [
        {
          field:"id",
          title:'ID'
        },
        {
          field:"loginid",
          title:'Login ID'
        },
        {
          field:"created",
          title:'Created',
          template:"#:kendo.toString(created,'dddd, MMMM dd, yyyy h:mm:ss tt')#"
        },
        {
          field:"expired",
          title:'Expired',
          template:"#:kendo.toString(expired,'dddd, MMMM dd, yyyy h:mm:ss tt')#"
        }],
      sortable: true,
      pageable: {
          refresh: true,
          pageSizes: true,
          buttonCount: 5
      },

    });
  });
}

adm.currentChangePassword = function(){

  if ( adm.dataChangePassword().newPassword !==  adm.dataChangePassword().reNewPassword ){
    swal("Error!", "New Password not match with re New Password", "error")
    return;
  }

  var url = "/web-cb/acluser/savenewpassword"
  var sessionid = localStorage.getItem("sessionid")  
  var param = {
    newpassword: adm.dataChangePassword().newPassword,
    sessionid: sessionid,
    oldpassword:adm.dataChangePassword().oldPassword  
  }

  ajaxPost(url, param, function (datas){
    if (datas == null){
     swal("Success","Change Passsword Success", "success");
     $("#newPassword").val('');
     $("#reNewPassword").val('');
     $("#oldPassword").val('');
    }else {
      swal("Error!", datas, "error")
      $("#newPassword").val('');
      $("#reNewPassword").val('');
      $("#oldPassword").val('');
    }

  });
}

adm.AddNewUserLogin = function(){
  var validator = $("#formAddUser").data("kendoValidator");
  if(validator==undefined){
       validator= $("#formAddUser").kendoValidator().data("kendoValidator");
  }
  if (validator.validate()) {
      var url = '/web-cb/aclsysadmin/savedatauser';  
      var password = adm.Password();
      var rePassword = adm.RePassword();
      if (password != rePassword){
        swal("Error!", 'Password did not match', "error");
        return;
      }
      adm.updateUser(false); 
      var parm = {
        _id:adm.Id(),
        loginid:adm.UserName(),
        fullname:adm.Fullname(),
        email:adm.Email(),    
        password:adm.Password(),
        oldpassword:adm.LastPassword(),
        groups:adm.valueGroup(),
        enable:$("#enableuser").is(':checked'),  
      }
      
      ajaxPost(url, parm, function(data){
        swal("Success","Change Data Success", "success");
        adm.addUser(false);
        adm.saveUser(false);
        adm.updateUser(false);
        adm.reset();
        adm.createGridUserAdministration();
      });
    }

 
}
adm.editData = function(s){
  adm.addUser(true);
  adm.updateUser(true);  
  adm.isNewUser(false);
  $("#addusernew").hide();
  var Selected = Enumerable.From(adm.listUser()).Where("$.id === '"+s+"'").ToArray();
  for (var i in Selected){
    var x = Selected[i];
    adm.Id(x.id);
    adm.UserName(x.loginid);
    adm.Fullname(x.fullname);
    adm.Email(x.email);
    adm.Password("");
    adm.RePassword("");
    adm.valueGroup(x.groups);
    adm.LastPassword(x.password);
    $('#enableuser').prop('checked', x.enable);
  }
}

adm.addusernew = function(){
  adm.addUser(true);
  adm.saveUser(true);
  adm.isNewUser(true);
  $("#addusernew").hide();
  adm.reset();
}
adm.doDelete = function(s){
    var param = {
       userid: s
    }
    var url = "/web-cb/aclsysadmin/deletedatauser";
    swal({
            title: "Are you sure?",
            text: "Are you sure delete this user!",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: '#DD6B55',
            confirmButtonText: 'Yes, I am sure!',
            cancelButtonText: "No, cancel it!",
            closeOnConfirm: false,
            closeOnCancel: false
        },
        function(isConfirm) {
            if (isConfirm) {
                ajaxPost(url, param, function(data){
                    if (data.IsError == false){
                        swal("Success!",data.Message,"success");
                        adm.createGridUserAdministration();
                    }else{
                        swal("Error!",data.Message,"error");
                    }
                });
            } else {
                swal("Cancelled", "Cancelled Delete User", "error");
            }
    });
}

adm.loadMenu = function(){
    var url = '/web-cb/aclsysadmin/getdataaccess'; 
    var param = {
    };
    adm.listMenu([{title: "[TOP LEVEL]", Id: ""}]);
    adm.treelistView([]);
    ajaxPost(url, param, function (data) {
        adm.sourceDataMenu(data.Data);
        var dataMenu = data.Data;
        var sortdataMenu =  Enumerable.From(dataMenu).OrderBy("$.parentid").ThenBy("$.index").ToArray();
        var dataTree =  adm.convert(sortdataMenu);
        
        var spacer = "--";
        var listSubmenu = [];
        for (var i in dataTree){
            if (dataTree[i].Submenus.length  != 0 ){
                adm.listMenu.push({
                        "title" : spacer + " " + dataTree[i].title,
                        "Id" : dataTree[i]._id
                    });
                adm.subMenuMaster(dataTree[i].Submenus, spacer);
                //=================== 
                adm.treelistView.push({
                    "_id" : dataTree[i]._id,
                    "title" : dataTree[i].title,
                    "url" : "#",
                    "icon" : dataTree[i].icon,
                    "parentid" : dataTree[i].parentid,
                    "index" : dataTree[i].index,
                    "enable" : dataTree[i].enable,
                });
                adm.subtreelist(dataTree[i].Submenus);

            }else{
                adm.listMenu.push({
                        "title" : spacer + " " + dataTree[i].title,
                        "Id" : dataTree[i]._id
                    });

                adm.treelistView.push({
                    "_id" : dataTree[i]._id,
                    "title" : dataTree[i].title,
                    "url" : "#",
                    "icon" : dataTree[i].icon,
                    "parentid" : dataTree[i].parentid,
                    "index" : dataTree[i].index,
                    "enable" : dataTree[i].enable,
                });

            }
        }

        var sortdataTree =  adm.convert(adm.treelistView());
        adm.ListMenuTree(sortdataTree);

        var inline = new kendo.data.HierarchicalDataSource({
                        data: adm.ListMenuTree(),
                        schema: {
                            model: {
                                children: "Submenus"
                            }
                        }
                    }); 

        var treeview = $("#menu-list").kendoTreeView({
                animation: false,
                template: kendo.template($("#menulist-template").html()),
                dataTextField: "Title",
                dataSource:inline,
                select: adm.selectDirFolder,
                loadOnDemand: false
            }).data("kendoTreeView");
            treeview.expand(".k-item");
    });
}

adm.selectDirFolder = function(e){
    adm.disableFormMenu();
    adm.btnAddMenu(true);
    adm.btnSaveMenu(false);
    adm.btnEditMenu(true);
    adm.btnDeleteMenu(true);
    // $('#btnSaveAccess').html("<i class='fa fa-save'></i> Update")
    var data = $('#menu-list').data('kendoTreeView').dataItem(e.node);
    var SelectId = "";
    if (data._id == undefined){
        SelectId = data.Id
    }else{
        SelectId = data._id
    }
    var Selected = Enumerable.From(adm.sourceDataMenu()).Where("$.id === '"+SelectId+"'").ToArray();
    for(var i in Selected){
      adm.Data.ID(Selected[i].id);
      adm.Data.Title(Selected[i].title);
      adm.Data.ParentId(Selected[i].parentid);
      adm.Data.Url(Selected[i].url); 
      adm.Data.Index(Selected[i].index);
      adm.Data.Group1(Selected[i].group1);
      $('#checkEnable').prop('checked', Selected[i].enable);
    }
 }

adm.editMenu = function(){
adm.btnAddMenu(false);
adm.enableFormMenu();
adm.btnSaveMenu(true);
$('#btnSaveAccess').html("<i class='fa fa-save'></i> Update")
}

adm.deleteMenu = function(){
    var param = {
        Id : adm.Data.ID()
    }
    var url = "/web-cb/aclsysadmin/deletedatamenu";
    swal({
            title: "Are you sure?",
            text: "Are you sure remove this menu!",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: '#DD6B55',
            confirmButtonText: 'Yes, I am sure!',
            cancelButtonText: "No, cancel it!",
            closeOnConfirm: false,
            closeOnCancel: false
        },
        function(isConfirm) {
            if (isConfirm) {
                ajaxPost(url, param, function(data){
                    if (data.IsError == false){
                        swal("Success!",data.Message,"success");
                        adm.loadMenu();
                    }else{
                        swal("Error!",data.Message,"error");
                    }
                });
            } else {
                swal("Cancelled", "Cancelled Delete Menu", "error");
            }
        });

}

adm.convert = function (array){
    var map = {};
    for(var i = 0; i < array.length; i++){
        var obj = array[i];
        obj.Submenus= [];

        map[obj._id] = obj;

        var parent = obj.parentid || '-';
        if(!map[parent]){
            map[parent] = {
                Submenus: []
            };
        }
        map[parent].Submenus.push(obj);
    }
    return map['-'].Submenus;
}

adm.subMenuMaster = function(SubData, spacer){
    spacer += "--";
    for (var i in SubData){
            if (SubData[i].Submenus.length != 0 ){
                adm.listMenu.push({
                        "title" : spacer + " " + SubData[i].title,
                        "Id" : SubData[i]._id
                    });
                adm.subMenuMaster(SubData[i].Submenus, spacer);
            }else{
                adm.listMenu.push({
                        "title" : spacer + " " + SubData[i].title,
                        "Id" : SubData[i]._id
                    });
            }
        }
}

adm.subtreelist = function(SubData){
    for (var i in SubData){
            if (SubData[i].Submenus.length != 0 ){
                adm.treelistView.push({
                        "_id" : SubData[i]._id,
                        "title" : SubData[i].title,
                        "url" : "#",
                        "icon" : SubData[i].icon,
                        "parentid" : SubData[i].parentid,
                        "index" : SubData[i].index,
                        "enable" : SubData[i].enable,
                    });
                adm.subtreelist(SubData[i].Submenus);
            }else{
                adm.treelistView.push({
                        "_id" : SubData[i]._id,
                        "title" : SubData[i].title,
                        "url" : "#",
                        "icon" : SubData[i].icon,
                        "parentid" : SubData[i].parentid,
                        "index" : SubData[i].index,
                        "enable" : SubData[i].enable,
                    });
            }
        }
}

adm.disableFormMenu = function(){
  $("#admparent").data("kendoDropDownList").enable(false);
  $("#admgroup").data("kendoDropDownList").enable(false);
  $("#admid").attr('disabled','disabled');
  $("#admtitle").attr('disabled','disabled');
  $("#admurl").attr('disabled','disabled');
  $("#admindex").attr('disabled','disabled');
  $("#checkEnable").attr('disabled','disabled');
}
adm.enableFormMenu = function(){
  $("#admparent").data("kendoDropDownList").enable(true);
  $("#admgroup").data("kendoDropDownList").enable(true);
  $("#admid").removeAttr('disabled');
  $("#admtitle").removeAttr('disabled');
  $("#admurl").removeAttr('disabled');
  $("#admindex").removeAttr('disabled');
  $("#checkEnable").removeAttr('disabled');
}

adm.showMenu = function(){
  adm.resetMenu();
  adm.disableFormMenu();
  adm.btnAddMenu(true);
  adm.btnSaveMenu(false);
  adm.btnEditMenu(false);
  adm.btnDeleteMenu(false);
  $('#btnSaveAccess').html("<i class='fa fa-save'></i> Save")
}
adm.addMenu = function(){
  adm.Data.ID("");
  adm.Data.Title("");
  adm.Data.ParentId("");
  adm.Data.Url(""); 
  adm.Data.Index("");
  adm.Data.Group1("");
  adm.enableFormMenu();
  adm.btnAddMenu(false);
  adm.btnSaveMenu(true);
  adm.btnEditMenu(false);
  adm.btnDeleteMenu(false);
  $('#btnSaveAccess').html("<i class='fa fa-save'></i> Save")
}

adm.resetMenu = function(){
  adm.FormValue(ko.mapping.fromJS(adm.Data));
}
adm.cancelMenu = function(){
  adm.btnAddMenu(true);
  adm.btnSaveMenu(false);
  adm.disableFormMenu();
}
adm.saveMenu = function(){
  adm.btnAddMenu(true);
  adm.btnSaveMenu(false);
  var url = '/web-cb/aclsysadmin/savemenu'; 
  var parm = ko.mapping.toJS(adm.FormValue());
  parm.Enable = $('#checkEnable').is(":checked");
  parm.Index = parseInt(parm.Index);
  ajaxPost(url,parm,function(res){
      if(res.IsError){
          swal("", res.Message, "info");
          return false;
      } 
      adm.resetMenu();
      adm.loadMenu();
      adm.disableFormMenu();
  })
}



adm.getGroup = function(){
var url = "/web-cb/aclsysadmin/getdataaccess" 
  var param = {}

  ajaxPost(url, param, function (datas){
   
    if (datas.Data != null){
     // swal("Success","Change Data Success", "success");
    }

  });

}

function allEmpoyee(){
  ajaxPost("/web-cb/aclsysadmin/getdatauser", {}, function(data) {
     adm.listUser(data.Data);         
  });
}
function readURL(event){
 var getImagePath = URL.createObjectURL(event.target.files[0]);
 $('#img-ex').attr("src",getImagePath);
}

$(document).on('change', '.btn-file :file', function() {
  var input = $(this),
      numFiles = input.get(0).files ? input.get(0).files.length : 1,
      label = input.val().replace(/\\/g, '/').replace(/.*\//, '');  
      input.trigger('fileselect', [numFiles, label]);
});
// ==============================================
$(function(){
  grp.getTopMenu();
  allEmpoyee();
  adm.createGridUserAdministration();
  adm.createGridGroupAdministrator();
  adm.createGridSessionAdministrator();
  adm.loadMenu();
  // adm.createGridAccessAdministrator();
  $("#grouptype").kendoNumericTextBox({
    min: 0
  });


  // for QOS only
  $('.btn-file :file').on('fileselect', function(event, numFiles, label) {        
        var input = $(this).parents('.input-group').find(':text'),
            log = numFiles > 1 ? numFiles + ' files selected' : label;        
        if( input.length ) {
            input.val(log);
        } else {
            if( log ) alert(log);
        }
        
    });
  // =================
})
