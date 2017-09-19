package controllers

import (
	// "eaciit/scb-apps/webapp/apps/main/models"
	// "github.com/eaciit/acl/v2.0"
	// db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	// tk "github.com/eaciit/toolkit"
	// bson "gopkg.in/mgo.v2/bson"
	// "log"
	// "strings"
)

type AdmController struct {
	*BaseController
}

func (c *AdmController) Login(k *knot.WebContext) interface{} {
	c.SetResponseTypeHTML(k)
	k.Config.LayoutTemplate = ""

	return c.SetViewData(nil)
}
