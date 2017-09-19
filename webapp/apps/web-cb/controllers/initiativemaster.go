package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	// "github.com/eaciit/cast"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// "gopkg.in/mgo.v2/bson"
	// "strings"
	"fmt"
	// "reflect"
	// "time"
)

type InitiativeMasterController struct {
	*BaseController
}

func (c *InitiativeMasterController) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	// c.Action(k, "Open Initiative Master Page")
	c.Action(k, "Initiative Master", "Open Initiative Master Page", "", "", "", "", "")

	InitiativeMaster := c.GetAccess(k, "INITIATIVEMASTER")
	Initiative := c.GetAccess(k, "INITIATIVE")
	Scorecard := c.GetAccess(k, "SCORECARD")
	k.Config.LayoutTemplate = "_layout-v2.html"

	PartialFiles := []string{}
	PartialFiles = append(PartialFiles, "shared/sidebar.html")
	PartialFiles = append(PartialFiles, "shared/initiative_summary.html")
	PartialFiles = append(PartialFiles, "shared/initiative_chart.html")
	PartialFiles = append(PartialFiles, "shared/metric_upload.html")
	PartialFiles = append(PartialFiles, "dashboard/initiative.html")

	k.Config.IncludeFiles = PartialFiles

	k.Config.OutputType = knot.OutputTemplate
	return tk.M{}.Set("InitiativeMaster", InitiativeMaster).Set("Initiative", Initiative).Set("Scorecard", Scorecard)
}

func (c *InitiativeMasterController) GetOwnedInitiative(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	InitiativeAccess := c.GetAccess(k, "INITIATIVE")
	Initiative := new(AccessibilityModel)
	e := tk.MtoStruct(InitiativeAccess, &Initiative)
	k.Config.OutputType = knot.OutputJson
	query := []*db.Filter{}

	csr, e := c.Ctx.Connection.NewQuery().From(new(InitiativeModel).TableName()).Cursor(nil)
	if Initiative.Global.Owned || Initiative.Region.Owned || Initiative.Country.Owned {
		fullname := ""
		if k.Session("fullname") != nil {
			fullname = k.Session("fullname").(string)
		}
		query = append(query, db.Or(db.Eq("ProjectManager", fullname), db.Eq("AccountableExecutive", fullname), db.Eq("TechnologyLead", fullname)))
		csr, e = c.Ctx.Connection.NewQuery().From(new(InitiativeModel).TableName()).Where(db.And(query...)).Cursor(nil)
	}

	result := make([]InitiativeModel, 0)
	e = csr.Fetch(&result, 0, false)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	csr.Close()
	OwnedInitiative := []string{}
	for _, x := range result {
		OwnedInitiative = append(OwnedInitiative, x.Id.Hex())
	}
	return c.SetResultInfo(false, "", OwnedInitiative)
}
func (c *InitiativeMasterController) InitiativeMasterData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	// c.Action(k, "")
	c.Action(k, "Initiative Master", "Open Initiative Master Data", "", "", "", "", "")

	k.Config.OutputType = knot.OutputJson

	csr, e := c.Ctx.Find(new(InitiativeModel), nil)
	result := make([]InitiativeModel, 0)
	e = csr.Fetch(&result, 0, false)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	csr.Close()
	return result
}

func (c *InitiativeMasterController) SummaryBusinessDriver(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	csr, e := c.Ctx.Find(new(SummaryBusinessDriverModel), tk.M{})
	defer csr.Close()

	results := make([]SummaryBusinessDriverModel, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error()
	}

	return results
}

func (c *InitiativeMasterController) LCList(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	csr, e := c.Ctx.Find(new(LifeCycleModel), tk.M{})
	defer csr.Close()

	results := make([]LifeCycleModel, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error()
	}

	return results
}

type InitiativeMasterSaveData struct {
	models []InitiativeModel
}

func (c *InitiativeMasterController) InitiativeMasterSave(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	type Logs struct {
		NewValue    string
		OldValue    string
		Whatchanged string
	}

	parameter := struct {
		InitiativeData InitiativeModel
		Logs           []Logs
	}{}

	e := k.GetPayload(&parameter)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	for _, la := range parameter.Logs {
		c.Action(k, "Initiative Master", "Save updated data on Initiative Master", la.Whatchanged, la.OldValue, la.NewValue, "", "")
	}

	data := NewInitiativeModel()
	csr, err := c.Ctx.Connection.NewQuery().From(data.TableName()).Where(db.Eq("_id", parameter.InitiativeData.Id)).Cursor(nil)
	err = csr.Fetch(&data, 1, false)
	csr.Close()
	if err != nil {
		return nil
	}

	if parameter.InitiativeData.Id == data.Id {

		data.InvestmentId = parameter.InitiativeData.InvestmentId
		data.ProjectName = parameter.InitiativeData.ProjectName
		data.StartDate = parameter.InitiativeData.StartDate
		data.FinishDate = parameter.InitiativeData.FinishDate
		data.ProjectManager = parameter.InitiativeData.ProjectManager
		data.AccountableExecutive = parameter.InitiativeData.AccountableExecutive
		data.TechnologyLead = parameter.InitiativeData.TechnologyLead
		data.ProblemStatement = parameter.InitiativeData.ProblemStatement
		data.ProjectDescription = parameter.InitiativeData.ProjectDescription
		data.ProgressCompletion = parameter.InitiativeData.ProgressCompletion
		data.PlannedCost = parameter.InitiativeData.PlannedCost
		data.IsGlobal = parameter.InitiativeData.IsGlobal
		data.Region = parameter.InitiativeData.Region
		data.Country = parameter.InitiativeData.Country
		data.BusinessImpact = parameter.InitiativeData.BusinessImpact
		data.EX = parameter.InitiativeData.EX
		data.OE = parameter.InitiativeData.OE
		data.SCCategory = parameter.InitiativeData.SCCategory
		data.BusinessDriverImpact = parameter.InitiativeData.BusinessDriverImpact
		data.LifeCycleId = parameter.InitiativeData.LifeCycleId
	}

	e = c.Ctx.Save(data)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	return c.SetResultInfo(false, "", nil)
}

func (c *InitiativeMasterController) InitiativeMasterSaveOld(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	p := InitiativeModel{}

	e := k.GetPayload(&p)
	if e != nil {
		fmt.Println(e)
		return e
	}

	if p.Id == "" {
		return "error"
	} else {
		// result := new(InitiativeModel)
		// e = c.Ctx.GetById(result, p.Id)
		// if e != nil {
		// 	return e
		// }

		// result.CommentList = p.CommentList

		if p.IsGlobal {
			p.Country = []string{}
			p.Region = []string{}
		}

		e = c.Ctx.Save(&p)
		if e != nil {
			return e
		}

		return p
	}
}

func (c *InitiativeMasterController) InitiativeMasterRemove(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := InitiativeModel{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	e = c.Ctx.DeleteMany(new(InitiativeModel), db.Eq("_id", parm.Id))
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	return parm
}
