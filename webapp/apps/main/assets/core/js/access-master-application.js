
master.dataApplication = ko.observableArray([])
master.dataApplicationsForDropDown = ko.computed(function () {
    return _.sortBy(master.dataApplication(), 'ID').map(function (d) {
        if (d.ID.toLowerCase() !== d.Name.toLowerCase() && d.Name !== '') {
            return { text: d.Name, value: d.ID }
        }
    
        return { text: d.ID, value: d.ID }
    })
}, master.dataApplication)

master.refreshGridApplication = function (callback) {
    master.dataApplication([])

    viewModel.ajaxPostCallback('/main/access/getapplication', {}, function (data) {
        master.dataApplication(_.sortBy(data, 'ID'))

        var config = {
            dataSource: {
                data: master.dataApplication(),
                schema: {
                    model: {
                        fields: {
                            ID: { type: 'string', editable: false },
                            Name: { type: 'string' },
                            LandingURL: { type: 'string' }
                        }
                    }
                }
            },
            sortable: true,
            filterable: true,
            editable: true,
            columns: [{
                field: 'ID',
                title: 'Application Id'
            }, {
                field: 'Name',
                title: 'Name'
            }, {
                field: 'LandingURL',
                title: 'Landing URL'
            }]
        }

        $('.grid-application').replaceWith('<div class="grid-application"></div>')
        $('.grid-application').kendoGrid(config)
        
        // ===== get user data

        master.dataGroup([])

        viewModel.ajaxPostCallback('/main/access/getgroup', {}, function (data) {

            data.forEach(function (d) {
                if (d.Applications == undefined || d.Applications == null) {
                    d.Applications = []
                }
            })
            master.dataGroup(data)

            var userAppColumns = master.dataApplication().map(function (d) {
                return {
                    field: 'groupapp-' + d.ID,
                    title: (d.Name !== '') ? d.Name : d.ID,
                    width: 100,
                    template: function (k) {
                        var checked = (k.Applications.indexOf(d.ID) > -1) ? 'checked' : ''
                        if (k.IsImportant) {
                            checked = 'checked disabled'
                        }

                        return '<input data-field="' + d.ID + '" type="checkbox" ' + checked + ' />'
                    }
                }
            })

            var config = {
                dataSource: {
                    data: data,
                    pageSize: 10,
                },
                columns: [{
                    field: 'Title',
                    title: 'Group Name',
                    width: 120,
                    locked: true
                }, {
                    headerTemplate: '<center><i class="fa fa-angle-double-left" style="font-weight: bold; font-size: 14px;"></i></center>',
                    width: 26,
                    locked: true
                }].concat(userAppColumns)
            }

            $('.grid-group-application').replaceWith('<div class="grid-group-application"></div>')
            $('.grid-group-application').kendoGrid(config)
        })
    })
}

master.saveApplication = function () {
    var payload = JSON.parse(kendo.stringify(
        $('.grid-application').data('kendoGrid').dataSource.data()))

    viewModel.ajaxPostCallback('/main/access/saveapplication', payload, function (data) {
        swal({
            title: 'Success',
            text: 'Changes saved',
            type: 'success',
            timer: 2000,
            showConfirmButton: false
        })
    
        $('#modal-access-menu').modal('hide')
        master.refreshGridApplication()
    })
}

master.saveGroupApplication = function () {
    var payload = []
    $('.grid-group-application').data('kendoGrid').dataSource.data().forEach(function (d) {
        var apps = []
        $('.grid-group-application').find('.k-grid-content [data-uid="' + d.uid + '"] input[data-field]').each(function (i, e) {
            if ($(e).is(':checked')) {
                apps.push($(e).attr('data-field'))
            }
        })

        var rowGroup = master.dataGroup().find(function (g) {
            return g._id == d._id
        })
        if (rowGroup !== undefined) {
            if (rowGroup.IsImportant) {
                return
            }
            
            rowGroup.Applications = apps
            payload.push(rowGroup)
        }
    })

    viewModel.ajaxPostCallback('/main/access/savegroups', payload, function (data) {
        swal({
            title: 'Success',
            text: 'Changes saved',
            type: 'success',
            timer: 2000,
            showConfirmButton: false
        })

        master.refreshGridApplication()
    })
}