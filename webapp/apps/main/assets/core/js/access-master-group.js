

master.dataGroup = ko.observableArray([])
master.dataGroupForDropDown = ko.computed(function () {
    return master.dataGroup().map(function (d) {
        if (d._id.toLowerCase() !== d.Title.toLowerCase()) {
            var text = d._id + " - " + d.Title
            return { text: text, value: d._id }
        }
    
        return { text: d.Title, value: d._id }
    })
}, master.dataGroup)
master.newAccessGrant = function () {
    return {
        AccessID: "",
        AccessValue: 1
    }
}
master.newGroup = function () {
    return {
        Enable: true,
        Grants: [],
        GroupConf: {},
        GroupType: 0,
        MemberConf: {},
        Owner: "",
        Title: "",
        _id: "",
        // ApplicationID: "",
        Applications: [],
        IsImportant: false
    }
}
master.selectedGroup = ko.mapping.fromJS(master.newGroup())
master.groupIsInsertMode = ko.observable(false)

master.refreshGridGroup = function () {
    master.dataGroup([])

    viewModel.ajaxPostCallback('/main/access/getgroup', {}, function (data) {

        // hacks for Access Grants
        data.forEach(function (d) {
            if (d.Grants == null || d.Grants == undefined) {
                d.Grants = []
            }

            if (d.Grants.length > 0) {
                d.Grants = d.Grants.map(function (e) {
                    return e.AccessID
                })
            }
        })

        master.dataGroup(data)

        var config = {
            dataSource: {
                data: data,
                pageSize: 10,
            },
            pageable: true,
            sortable: true,
            filterable: true,
            columns: [{
                field: '_id',
                title: 'Group ID'
            }, {
                field: 'Title',
                title: 'Name'
            }, {
                field: 'Applications',
                template: function (d) {
                    if (d.Applications.length == 1) {
                        return d.Applications[0]
                    }
                    
                    return d.Applications.map(function (g) {
                        var row = master.dataApplicationsForDropDown().find(function (l) {
                            return l.value == g.ApplicationID
                        })
                        if (row !== undefined) {
                            return row.text
                        }

                        return g
                    }).map(function (g, j) {
                        return ' ' + (j + 1) + '. ' + g
                    }).join('<br />')
                }
            // }, {
            //     title: 'Grants Access Menu',
            //     template: function (d) {
            //         return d.Grants.map(function (e) { 
            //             var menu = master.dataAccessMenuFlat().find(function (f) {
            //                 return f._id === e
            //             })
            //             if (menu !== undefined) {
            //                 return ' - ' + menu.Title
            //             }

            //             return ' - ' + e
            //         }).join('<br />')
            //     }
            }, {
                title: '&nbsp;',
                width: 80,
                attributes: { class: 'align-center' },
                template: function (d) {
                    var disabled = d.IsImportant ? 'disabled' : ''

                    return "<button class='btn btn-xs btn-primary' data-tooltipster='Edit' onclick='master.editGroup(\"" + d._id + "\")' " + disabled + "><i class='fa fa-edit'></i></button>"
                         + "&nbsp;"
                         + "<button class='btn btn-xs btn-danger' onclick='master.deleteGroup(\"" + d._id + "\")' data-tooltipster='Remove' " + disabled + "><i class='fa fa-trash'></i></button>"
                }
            }],
            dataBound: function () {
                viewModel.prepareTooltipsterGrid(this)
            }
        }

        $('.grid-group').replaceWith('<div class="grid-group"></div>')
        $('.grid-group').kendoGrid(config)
    })
}

master.editGroup = function (_id) {
    master.dataAccessAll([])

    viewModel.ajaxPostCallback('/main/access/getaccessmenu', {}, function (data) {
        if (data.length == 0) {
            return
        }

        master.dataAccessAll(data)

        var data = master.dataGroup().find(function (d) { return d._id === _id })
        
        master.groupIsInsertMode(false)
        ko.mapping.fromJS(data, master.selectedGroup)
        $('#modal-group').modal('show')

        setTimeout(function () { viewModel.isFormValid('#modal-group form') }, 310)

        setTimeout(function () {
            master.selectedGroup.Grants()
            
            setTimeout(function () {
                master.selectedGroup.Grants(data.Grants)
            }, 100)
        }, 100)
    })
}

master.createGroup = function () {
    master.dataAccessAll([])

    viewModel.ajaxPostCallback('/main/access/getaccessmenu', {}, function (data) {
        if (data.length == 0) {
            return
        }

        master.dataAccessAll(data)

        master.groupIsInsertMode(true)
        ko.mapping.fromJS(master.newGroup(), master.selectedGroup)
        $('#modal-group').modal('show')

        setTimeout(function () { viewModel.isFormValid('#modal-user form') }, 310)
    })
}

master.saveGroup = function () {
    if (!viewModel.isFormValid('#modal-group form')) {
        swal("Error!", "Some inputs are not valid", "error")
        return
    }
    
    var payload = ko.mapping.toJS(master.selectedGroup)
    if (payload.Grants.length > 0) {
        payload.Grants = payload.Grants.map(function (e) {
            var row = master.newAccessGrant()
            row.AccessID = e
            return row
        })
    }
        
    viewModel.ajaxPostCallback('/main/access/savegroup', payload, function (data) {
        swal({
            title: 'Success',
            text: 'Changes saved',
            type: 'success',
            timer: 2000,
            showConfirmButton: false
        })
    
        $('#modal-group').modal('hide')
        master.refreshGridGroup()
    })
}

master.deleteGroup = function (_id) {
    swal({
        title: "Are you sure?",
        text: "You will not be able to recover deleted data!",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#DD6B55",
        confirmButtonText: "Yes, delete it!",
        closeOnConfirm: false
    }, function(){
        var payload = master.newGroup()
        payload._id = _id

        viewModel.ajaxPostCallback('/main/access/deletegroup', payload, function (data) {
            swal({
                title: 'Success',
                text: 'Menu successfully deleted',
                type: 'success',
                timer: 2000,
                showConfirmButton: false
            })
        
            $('#modal-group').modal('hide')
            master.refreshGridGroup()
        })
    });
}

master.modalGroupChangeApplication = function () {
    setTimeout(function () {
        master.selectedGroup.Grants([])
    }, 100)
}