var master = {}
viewModel.master = master

master.activeMenu = ko.observable('User') // User, Group, Menu, Log
master.toggleActiveMenu = function (activeMenu, obj) {
    master.activeMenu(activeMenu)

    $(obj).siblings().removeClass('active')
    $(obj).addClass('active')

    switch (activeMenu) {
        case 'User': 
            master.refreshGridUser()
        break
        case 'Group': 
            master.refreshGridGroup()
        break
        case 'AccessMenu': 
            master.refreshGridMenu()
        break
        case 'Application': 
            master.refreshGridApplication()
        break
        case 'Log': 
            master.refreshGridLog()
        break
    }
}

$(function () {
    master.refreshGridGroup()
    master.refreshGridUser()
})