var dashboard = {}
viewModel.dashboard = dashboard

dashboard.dataApplication = ko.observableArray([])

dashboard.getCurrentUserApplication = function () {
    viewModel.ajaxPostCallback('/main/public/getapplication', {}, function (data) {
        data.forEach(function (d) {
            if (d.Name == '') {
                d.Name = d.Id
            }
            if (d.LandingURL[0] !== '' && d.LandingURL[0] !== '/') {
                d.LandingURL = '/'
            }
        })
        dashboard.dataApplication(data)
    })
}

$(function () {
    dashboard.getCurrentUserApplication()
})