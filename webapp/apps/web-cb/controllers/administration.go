package controllers

import (
	// db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	// tk "github.com/eaciit/toolkit"
	// "encoding/json"
	// . "eaciit/scbocir/models"
	// "sort"
	// "strings"
	// "reflect"
)

type AdministrationController struct {
	*BaseController
}

func (c *AdministrationController) Default(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	k.Config.NoLog = true
	k.Config.IncludeFiles = []string{"administration/user.html", "administration/session.html", "administration/access.html", "administration/group.html"}
	k.Config.OutputType = knot.OutputTemplate
	DataAccess := Previlege{}

	for _, o := range access {
		DataAccess.Create = o["Create"].(bool)
		DataAccess.View = o["View"].(bool)
		DataAccess.Delete = o["Delete"].(bool)
		DataAccess.Process = o["Process"].(bool)
		DataAccess.Delete = o["Delete"].(bool)
		DataAccess.Edit = o["Edit"].(bool)
		DataAccess.Menuid = o["Menuid"].(string)
		DataAccess.Menuname = o["Menuname"].(string)
		DataAccess.Approve = o["Approve"].(bool)
		DataAccess.Username = o["Username"].(string)
	}

	return DataAccess
}
