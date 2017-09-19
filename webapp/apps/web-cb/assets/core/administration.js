var adm = {}

adm.searchGrid = ko.observable(false);
adm.addUser = ko.observable(false);


adm.createGridUserAdministration = function(){
  $("#gridUserAdministrator").html("");
  $("#gridUserAdministrator").kendoGrid({
      dataSource: [{}], 
      columns: [
        {
          field:"Count",
          title:'Login ID'
        },
        {
          field:"Count",
          title:'Fullname'
        },
        {
          field:"Count",
          title:'Email'
        },
        {
          field:"Count",
          title:'Password'
        },
        {
          field:"Count",
          title:'FTE'
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
}


adm.createGridGroupAdministrator = function(){
  $("#gridGroupAdministrator").html("");
  $("#gridGroupAdministrator").kendoGrid({
      dataSource: [{}], 
      columns: [
        {
          field:"Count",
          title:'ID'
        },
        {
          field:"Count",
          title:'Title'
        },
        {
          field:"Count",
          title:'Enable'
        },
        {
          field:"Count",
          title:'Owner'
        },
        {
          field:"Count",
          title:'Action'
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
}

adm.createGridAccessAdministrator = function(){
  $("#gridAccessAdministrator").html("");
  $("#gridAccessAdministrator").kendoGrid({
      dataSource: [{}], 
      columns: [
        {
          field:"Count",
          title:'ID'
        },
        {
          field:"Count",
          title:'Title'
        },
        {
          field:"Count",
          title:'Group 1'
        },
        {
          field:"Count",
          title:'Group 2'
        },
        {
          field:"Count",
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
}

adm.createGridSessionAdministrator = function(){
  $("#gridSessionAdministrator").html("");
  $("#gridSessionAdministrator").kendoGrid({
      dataSource: [{}], 
      columns: [
        {
          field:"Count",
          title:'Status'
        },
        {
          field:"Count",
          title:'Username'
        },
        {
          field:"Count",
          title:'Created'
        },
        {
          field:"Count",
          title:'Expired'
        },
        {
          field:"Count",
          title:'Active In'
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
}


$(function(){
  adm.createGridUserAdministration();
  adm.createGridGroupAdministrator();
  adm.createGridAccessAdministrator();
  adm.createGridSessionAdministrator();
})