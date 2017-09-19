
master.dataUser = ko.observableArray([])
master.newUser = function () {
    return {
        _id: "",
        LoginID: "",
        FullName: "",
        Email: "",
        Password: "",
        Groups: [],
        IsImportant: false,

        Enable: false,
        Grants: [],
        LoginType: 0,
        LoginConf: {}
    }
}
master.selectedUser = ko.mapping.fromJS(master.newUser())
master.userIsInsertMode = ko.observable(false)

master.refreshGridUser = function () {
    master.dataUser([])

    viewModel.ajaxPostCallback('/main/access/getapplication', {}, function (data) {
        master.dataApplication(_.sortBy(data, 'Id'))

        viewModel.ajaxPostCallback('/main/access/getuser', {}, function (data) {
            master.dataUser(data)

            data.forEach(function (d) {
                if (d.Applications == undefined || d.Applications == null) {
                    d.Applications = []
                }
            })

            var config = {
                dataSource: {
                    data: data,
                    pageSize: 5
                },
                pageable: true,
                sortable: true,
                filterable: true,
                columns: [{
                    field: 'LoginID',
                    title: 'Username'
                }, {
                    field: 'Email',
                    title: 'Email'
                }, {
                    field: 'FullName',
                    title: 'Name'
                }, {
                    title: 'Group',
                    template: function (d) {
                        if (d.Groups.length == 1) {
                            return d.Groups[0]
                        }

                        return d.Groups.map(function (k, i) {
                            return ' ' + (i + 1) + '. ' + k
                        }).join('<br />')
                    }
                }, {
                    title: '&nbsp;',
                    width: 80,
                    attributes: { class: 'align-center' },
                    template: function (d) {
                        var disabled = (d.IsImportant) ? 'disabled' : '';

                        return "<button class='btn btn-xs btn-primary' data-tooltipster='Edit' onclick='master.editUser(\"" + d._id + "\")' " + disabled + "><i class='fa fa-edit'></i></button>"
                            + "&nbsp;"
                            + "<button class='btn btn-xs btn-danger' data-tooltipster='Remove' " + disabled + "><i class='fa fa-trash' onclick='master.deleteUser(\"" + d._id + "\")'></i></button>"
                    }
                }],
                dataBound: function () {
                    viewModel.prepareTooltipsterGrid(this)
                }
            }

            $('.grid-user').replaceWith('<div class="grid-user"></div>')
            $('.grid-user').kendoGrid(config)
        }, {
            loader: false
        })
    })
}

master.editUser = function (_id) {
    var data = master.dataUser().find(function (d) { return d._id === _id })
    
    master.userIsInsertMode(false)
    ko.mapping.fromJS(data, master.selectedUser)
    $('#modal-user').modal('show')

    setTimeout(function () { viewModel.isFormValid('#modal-user form') }, 310)
}

master.createUser = function () {
    master.userIsInsertMode(true)
    ko.mapping.fromJS(master.newUser(), master.selectedUser)
    $('#modal-user').modal('show')

    setTimeout(function () { viewModel.isFormValid('#modal-user form') }, 310)
}

master.saveUser = function () {
    if (!viewModel.isFormValid('#modal-user form')) {
        swal("Error!", "Some inputs are not valid", "error")
        return
    }
    
    var payload = ko.mapping.toJS(master.selectedUser)

    viewModel.ajaxPostCallback('/main/access/saveuser', payload, function (data) {
        swal({
            title: 'Success',
            text: 'Changes saved',
            type: 'success',
            timer: 2000,
            showConfirmButton: false
        })
    
        $('#modal-user').modal('hide')
        master.refreshGridUser()
    })
}

master.deleteUser = function (_id) {
    swal({
        title: "Are you sure?",
        text: "You will not be able to recover deleted data!",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#DD6B55",
        confirmButtonText: "Yes, delete it!",
        closeOnConfirm: false
    }, function(){
        var payload = master.newUser()
        payload._id = _id

        viewModel.ajaxPostCallback('/main/access/deleteuser', payload, function (data) {
            swal({
                title: 'Success',
                text: 'Menu successfully deleted',
                type: 'success',
                timer: 2000,
                showConfirmButton: false
            })
        
            $('#modal-user').modal('hide')
            master.refreshGridUser()
        })
    });
}