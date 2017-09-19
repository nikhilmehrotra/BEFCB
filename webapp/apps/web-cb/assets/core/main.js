/* 
 @Created : 20161026 - Ainur
 @Note : Clearing code from layout
 */
function GridColumn(id, userId, gridId) {
    this.Name = id + "_" + userId + "_" + gridId;
    this.GridId = gridId;
    this.Init = function () {
        var grid = $('#' + this.GridId).data('kendoGrid');

        var gridCols = grid.columns;
        var datas = readCookie(this.Name);
        if (datas == null) {
            var visibleCols = [];
            jQuery.each(gridCols, function (index) {
                if (!this.hidden) {
                    visibleCols.push(this.field);
                }
            });
            createCookie(this.Name, visibleCols.join('|'));
        } else {
            var showColumns = datas.split('|');
            jQuery.each(gridCols, function (index) {
                if (showColumns.indexOf(this.field) < 0) {
                    grid.hideColumn(this.field);
                }
            });
        }
    };
    this.GetColumns = function () {
        return readCookie(this.Name);
    };
    this.AddColumn = function (columnName) {
        var datas = this.GetColumns();
        var cols = [];
        if (datas != null) {
            cols = datas.split('|');
            if (cols.indexOf(columnName) < 0) {
                cols.push(columnName);
            }
        }

        eraseCookie(this.Name);
        createCookie(this.Name, cols.join('|'));
    };
    this.RemoveColumn = function (columnName) {
        var datas = this.GetColumns();
        var cols = [];
        if (datas != null) {
            cols = datas.split('|');
            if (cols.indexOf(columnName) > -1) {
                jQuery.each(cols, function (index) {
                    if (cols[index] == columnName) {
                        cols.splice(index, 1);
                    }
                });
            }
        }

        eraseCookie(this.Name);
        createCookie(this.Name, cols.join('|'));
    };
}
function removeSpace(value) {
    return value.replace(" ", "");
}
function createCookie(name, value, days) {
    days = days || 365;
    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        var expires = "; expires=" + date.toGMTString();
    } else
        var expires = "";
    document.cookie = name + "=" + value + expires + "; path=/";
}

function readCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ')
            c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) == 0)
            return c.substring(nameEQ.length, c.length);
    }
    return null;
}

function eraseCookie(name) {
    createCookie(name, "", -1);
}

function getUTCDate(strdate) {
    var d = moment.utc(strdate);
    return new Date(d.year(), d.month(), d.date(), 0, 0, 0)
}
function getUTCDateFull(strdate) {
    var d = moment.utc(strdate);
    return new Date(d.year(), d.month(), d.date(), d.hours(), d.minutes(), d.seconds())
}

function toUTC(d) {
    var year = d.getFullYear();
    var month = d.getMonth();
    var date = d.getDate();
    var hours = d.getHours();
    var minutes = d.getMinutes();
    var seconds = d.getSeconds();
    return moment(Date.UTC(year, month, date, hours, minutes, seconds)).toISOString();
}
function getStringQuery(name) {
    var url = window.location.href;
    name = name.replace(/[\[\]]/g, "\\$&");
    var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
            results = regex.exec(url);
    if (!results)
        return null;
    if (!results[2])
        return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
}
function Logout() {
    localStorage.clear();
    var cookies = document.cookie.split(";");
    for (var i in cookies) {
        eraseCookie(cookies[i].split("=")[0]);
    }
    // window.location.href = "/logout/do";
    window.location.href = "/acluser/default";

}

function MenuItem(id, url, title, submenus, baseURL) {
    var obj = {
        _id: ko.observable(id),
        Title: ko.observable(title == undefined ? id : title),
        Url: ko.observable(url.replace("~/", baseURL)),
        Submenus: ko.observableArray([])
    };

    var arr = submenus;
    for (var i in arr) {
        obj.Submenus.push(
                new MenuItem(
                        arr[i]._id,
                        arr[i].Url,
                        arr[i].Title,
                        arr[i].Submenus,
                        baseURL
                        )
                );
    }
    return obj;
}
;

function BreadCrumb(id, title, url, cssClass, action) {
    var obj = {
        _id: ko.observable(id),
        Title: ko.observable(title === undefined ? id : title),
        Url: url,
        Action: action,
        CssClass: cssClass
    };

    return obj;
}


model.ParentId = ko.observableArray([]);
model.PageId = ko.observable("");
model.MainMenus = ko.observableArray([]);
model.RoleGroup = ko.observable("");
model.BreadCrumbs = ko.observableArray([]);
model.VisibleDropdown = ko.observableArray([]);
model.menu = ko.observable('');


function convert(array) {
    var map = {};
    for (var i = 0; i < array.length; i++) {
        var obj = array[i];
        obj.Submenus = [];

        map[obj.Id] = obj;

        var parent = obj.Parent || '-';
        if (!map[parent]) {
            map[parent] = {
                Submenus: []
            };
        }
        map[parent].Submenus.push(obj);
    }
    return map['-'].Submenus;
}

function MenuItem(Id, Url, Title, Submenus, baseURL) {
    var obj = {
        Id: ko.observable(Id),
        Title: ko.observable(Title === undefined ? Id : Title),
        Url: ko.observable(baseURL + Url),
        Submenus: ko.observableArray([])
    };
    var arr = Submenus;
    for (var i in arr) {
        obj.Submenus.push(
                new MenuItem(
                        arr[i]._id,
                        arr[i].url,
                        arr[i].title,
                        arr[i].submenu,
                        baseURL
                        )
                );
    }
    return obj;
}
;

function convert(array) {
    var map = {};
    for (var i = 0; i < array.length; i++) {
        var obj = array[i];
        obj.Submenus = [];

        map[obj.Id] = obj;

        var parent = obj.Parent || '-';
        if (!map[parent]) {
            map[parent] = {
                Submenus: []
            };
        }
        map[parent].Submenus.push(obj);
    }
    return map['-'].Submenus;
}
model.GetDataMenu = function () {
    // var url = "/menusetting/getmenutop";
    var url = "/web-cb/acluser/getlisttopmenu";
    var param = {
    };
    ajaxPost(url, param, function (dataMenu) {
        var arrData = [];
        var baseURL = "";
        for (var i in dataMenu) {
            arrData.push({
                "Id": dataMenu[i]._id,
                "_id": dataMenu[i]._id,
                "Url": dataMenu[i].url,
                "Title": dataMenu[i].title,
                "IndexMenu": dataMenu[i].index,
                "submenu": dataMenu[i].submenu
            });
        }
        ;
        var sortMenu = Enumerable.From(arrData).OrderBy("$.IndexMenu").ToArray();
        var dataTree = convert(sortMenu);
        var arr = dataTree;
        for (var i in arr) {
            model.MainMenus.push(
                    new MenuItem(
                            arr[i]._id,
                            arr[i].Url,
                            arr[i].Title,
                            arr[i].submenu,
                            baseURL
                            )
                    );
        }

    });
};
model.Init = function () {
    model.GetDataMenu();
};
$(document).ready(function () {
    model.Init();
});


