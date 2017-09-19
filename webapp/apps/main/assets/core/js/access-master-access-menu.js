

master.dataAccessAll = ko.observableArray([])
master.dataAccessAllByGroupDataForDropDown = ko.computed(function () {
    return []
    var rows = master.dataAccessAll().filter(function (d) {
        return d.ApplicationID == master.selectedGroup.ApplicationID()
    })
    
    return _.sortBy(rows, 'Title')
}, master)

master.dataAccessMenuTree = ko.observableArray([])
master.dataAccessMenuFlat = ko.observableArray([])
master.dataAccessMenuForDropdown = ko.computed(function () {
    return master.dataAccessMenuFlat().map(function (d) {
        if (d.ParentId !== "") {
            var parent = master.dataAccessMenuFlat().find(function (e) { return e._id === d.ParentId })
            if (parent !== undefined) {
                var text = (parent.Title + ' - ' + d.Title)
                return { text: text, value: d._id }
            }
        }

        return { text: d.Title, value: d._id }
    })
}, master.dataGroup)
master.newAccessMenu = function () {
    return {
        _id: "",
        Title: "",
        Icon: "",
        ParentId: "",
        Url: "#",
        Index: 1,
        ApplicationID: "",

        Category: 1,
        Group1: "",
        Group2: "",
        Group3: "",
        Enable: true,
        SpecialAccess1: "",
        SpecialAccess2: "",
        SpecialAccess3: "",
        SpecialAccess4: ""
    }
}
master.selectedAccessMenuApplicationID = ko.observable('')
master.accessMenuIsInsertMode = ko.observable(false)
master.selectedAccessMenu = ko.mapping.fromJS(master.newAccessMenu())

master.changeAccessMenuApplication = function () {
    setTimeout(function () {
        master.refreshGridMenu()
    }, 110)
}

master.refreshGridMenu = function () {
    if (master.selectedAccessMenuApplicationID() == '') {
        return
    }

    $('.access-menu-message').hide()

    master.dataAccessMenuTree([])
    master.dataAccessMenuFlat([])

    var payload = {}
    payload.applicationID = master.selectedAccessMenuApplicationID()

    viewModel.ajaxPostCallback('/main/access/getaccessmenubyapplication', payload, function (data) {
        if (data.length == 0) {
            $('.access-menu-message').show()
            return
        }

        data.forEach(function (d) {
            d.Submenu = []
        })

        var menu = []

        data.filter(function (d) {
            return d.Category == 1
        }).forEach(function (d) {
            menu.push(d)
        })

        data.filter(function (d) {
            return d.Category == 2
        }).forEach(function (d) {
            var parent = menu.find(function (e) {
                return e._id == d.ParentId
            })
            if (parent !== undefined) {
                parent.Submenu.push(d)
            }
        })

        master.dataAccessMenuTree(menu)
        master.dataAccessMenuFlat(data)

        setTimeout(function () {
            master.editAccessMenu(master.dataAccessMenuTree()[0])
        }, 100)
    })
}

master.editAccessMenu = function (data) {
    return function () {
        master.accessMenuIsInsertMode(false)
        ko.mapping.fromJS(data, master.selectedAccessMenu)
        $('#modal-access-menu').modal('show')
        
        setTimeout(function () { viewModel.isFormValid('#modal-access-menu form') }, 310)
    }
}

master.createAccessMenu = function () {
    master.accessMenuIsInsertMode(true)
    ko.mapping.fromJS(master.newAccessMenu(), master.selectedAccessMenu)
    master.selectedAccessMenu.Index(master.dataAccessMenuTree().length + 1)
    master.selectedAccessMenu.ApplicationID(master.selectedAccessMenuApplicationID())
    $('#modal-access-menu').modal('show')

    setTimeout(function () { viewModel.isFormValid('#modal-access-menu form') }, 310)
}

master.saveAccessMenu = function () {
    if (!viewModel.isFormValid('#modal-access-menu form')) {
        swal("Error!", "Some inputs are not valid", "error")
        return
    }
    
    var payload = ko.mapping.toJS(master.selectedAccessMenu)
    if (payload._id === payload.ParentId) {
        swal("Error!", "Parent ID cannot be same with current ID", "error")
        return
    }

    viewModel.ajaxPostCallback('/main/access/saveaccessmenu', payload, function (data) {
        swal({
            title: 'Success',
            text: 'Changes saved',
            type: 'success',
            timer: 2000,
            showConfirmButton: false
        })
    
        $('#modal-access-menu').modal('hide')
        master.refreshGridMenu()
    })
}

master.deleteAccessMenu = function () {
    swal({
        title: "Are you sure?",
        text: "You will not be able to recover deleted data!",
        type: "warning",
        showCancelButton: true,
        confirmButtonColor: "#DD6B55",
        confirmButtonText: "Yes, delete it!",
        closeOnConfirm: false
    }, function(){
        var payload = ko.mapping.toJS(master.selectedAccessMenu)

        viewModel.ajaxPostCallback('/main/access/deleteaccessmenu', payload, function (data) {
            swal({
                title: 'Success',
                text: 'Menu successfully deleted',
                type: 'success',
                timer: 2000,
                showConfirmButton: false
            })
        
            $('#modal-access-menu').modal('hide')
            master.refreshGridMenu()
        })
    })
}