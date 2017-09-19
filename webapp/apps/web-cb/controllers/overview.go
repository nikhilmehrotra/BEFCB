package controllers

import (
	// . "eaciit/scb-apps/webapp/apps/web-cb/models"
	// "github.com/eaciit/cast"
	// db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// "gopkg.in/mgo.v2/bson"
	// "strings"
	// "fmt"
	// "time"
)

type OverviewController struct {
	*BaseController
}

func (c *OverviewController) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputTemplate
	Overview := c.GetAccess(k, "OVERVIEW")
	Scorecard := c.GetAccess(k, "SCORECARD")
	Initiative := c.GetAccess(k, "INITIATIVE")
	k.Config.LayoutTemplate = "_layout-v2.html"
	PartialFiles := []string{}
	PartialFiles = append(PartialFiles, "shared/sidebar.html")
	PartialFiles = append(PartialFiles, "shared/initiative_summary.html")
	PartialFiles = append(PartialFiles, "shared/initiative_chart.html")
	PartialFiles = append(PartialFiles, "shared/metric_upload.html")
	PartialFiles = append(PartialFiles, "dashboard/initiative.html")

	k.Config.IncludeFiles = PartialFiles
	return tk.M{}.Set("Overview", Overview).Set("Initiative", Initiative).Set("Scorecard", Scorecard)
}
