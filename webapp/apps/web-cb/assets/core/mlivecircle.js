model.Mode = ko.observable("");
var lc = {
    IdP:ko.observable(""),
    Id: ko.observable(""),
    Name: ko.observable(""),
    Seq: ko.observable(0),
};

lc.getDataLC = function() {
    model.Mode("Process");
    var url = "/m/getlivecircledata";
    var param = {};
    ajaxPost(url, param, function(res) {
        lc.Grid(res.data);
        model.Mode("");
    });
};

lc.Grid = function(data) {
    $("#gridMasterLC").kendoGrid({
        dataSource: {
            data: data,
            pageSize: 10,
        },
        resizable: true,
        sortable: true,
        pageable: {
            refresh: true,
            pageSizes: true,
            buttonCount: 5
        },
        columns: [{
            field: "LiveCircleId",
            title: "Id",
            width: 100,
        }, {
            field: "Name",
            title: "Name",
            width: 100,
        }, {
            field: "Seq",
            title: "Seq",
            width: 100,
        },
        {
            template:"<button onclick='lc.Edit(\"#:Id #\",\"#:LiveCircleId #\",\"#:Name #\",\"#:Seq #\")' class='btn btn-sm btn-orange'><i class='fa fa-pencil'></i></button> "+
            "<button onclick='lc.Delete(\"#:Id #\")' class='btn btn-sm btn-lightred'><i class='fa fa-trash'></i></button>",
            title:"Actions",
            width:50
            } ],
    });
};

lc.addNew = function() {
    $('#modalLC').modal('show');
    lc.resetData();
}

lc.resetData = function(){
    lc.Id("");
    lc.Name("");
    lc.Seq("");
}

lc.saveData = function(){
     if (!model.isFormValid(".modal-body")) {
          return;
        }

      var url = "/m/savelivecircle";
      var param = {
        Id:lc.IdP(),
        LiveCircleId :lc.Id(),
        Name :lc.Name(),
        Seq : parseInt(lc.Seq()),
      };
     
      $('#gridMasterLC').html('');
      model.Processing(true);
      ajaxPost(url, param, function(data) {
        if(data=='') {
          $('#modalLC').modal('hide');
          swal('Success', 'Data has been saved successfully!', 'success');
          lc.getDataLC();         
          model.Processing(false);
        } else {
          swal("Warning", data, "error");
          model.Processing(false);
        }
      }, undefined);
}

lc.Edit = function(i,l,n,s){
$('#modalLC').modal('show');
    lc.IdP(i);
    lc.Id(l);
    lc.Name(n);
    lc.Seq(s);
}

lc.Delete = function(id){
    // console.log(i)
    swal({
          title: "Are you sure?",
          text: "Your will not be able to recover this data",
          type: "warning",
          showCancelButton: true,
          confirmButtonClass: "btn-danger",
          confirmButtonText: "Yes, delete it!",
          closeOnConfirm: false
        },
        function(res){
          if(res) {
          $('#gridMasterLC').html('');
          model.Processing(true);
            var url = "/m/deletelivecircle";
            var param = { id: id };
            ajaxPost(url, param, function(data){
              if(data != ""){
                swal('Warning', data, 'error');
                $('#modalLC').modal('hide');
                lc.getDataLC();   
                model.Processing(false);
              }else{
                swal('Success', 'Data has been deleted!', 'success');
                $('#modalLC').modal('hide');
                lc.getDataLC();   
                model.Processing(false);
              }

            }, undefined);
          }
          else {
            model.Processing(false);
          }
        });
    
}

model.isFormValid = function (selector) {
    model.resetValidation(selector);
    var $validator = $(selector).data("kendoValidator");
    return ($validator.validate());
};

model.resetValidation = function (selectorID) {
    var $form = $(selectorID).data("kendoValidator");
    if ($form == undefined) {
        $(selectorID).kendoValidator();
        $form = $(selectorID).data("kendoValidator");
    }

    $form.hideMessages();
};
$(document).ready(function() {
    lc.getDataLC();
});