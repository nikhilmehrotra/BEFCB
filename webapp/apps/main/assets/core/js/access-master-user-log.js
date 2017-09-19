master.refreshGridLog = function () {
    viewModel.ajaxPostCallback('/main/access/getuserlog', {}, function (data) {
        var config = {
            dataSource: {
                data: data,
                schema: {
                    model: {
                        fields: {
                            Created: { type: 'date' },
                            Expired: { type: 'date' },
                        }
                    }
                },
                pageSize: 10,
            },
            pageable: true,
            sortable: true,
            filterable: true,
            columns: [{
                field: 'LoginID',
                title: 'Username'
            }, {
                field: 'Created',
                title: 'Login Date',
                template: function (d) {
                    return moment(d.Created).format('YYYY-MMM-DD HH:mm:ss')
                }
            }, {
                field: 'Expired',
                title: 'Activity End (Logout / Session Expired)',
                template: function (d) {
                    return moment(d.Expired).format('YYYY-MMM-DD HH:mm:ss')
                }
            }],
            dataBound: function () {
                viewModel.prepareTooltipsterGrid(this)
            }
        }

        $('.grid-log').replaceWith('<div class="grid-log"></div>')
        $('.grid-log').kendoGrid(config)
    })
}