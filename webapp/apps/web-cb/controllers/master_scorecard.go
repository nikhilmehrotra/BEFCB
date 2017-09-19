package controllers

import (
	// . "eaciit/crm/commons"
	// . "eaciit/crm/helper"
	m "eaciit/scb-apps/webapp/apps/web-cb/models"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// "gopkg.in/mgo.v2/bson"
	// "strings"
	// "time"
)

func (m *MController) Scorecard(k *knot.WebContext) interface{} {
	access := m.LoadBase(k)
	k.Config.NoLog = true
	k.Config.LayoutTemplate = "_layout_dedicated.html"
	k.Config.IncludeFiles = []string{"shared/sidebar.html"}
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

func (c *MController) GetPrototypeScorecard(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true
	return c.SetResultInfo(false, "", m.NewBusinessDriverL1Model())
}

func (c *MController) GetDataScorecard(k *knot.WebContext) interface{} {
	tk.Println("GetDataScorecard")
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	result := []m.BusinessDriverL1Model{}

	csr, err := c.Ctx.Connection.NewQuery().From("BusinessDriverL1").Cursor(nil)
	err = csr.Fetch(&result, 0, false)
	csr.Close()
	if err != nil {
		tk.Println(err)
	}
	tk.Println("Query OK GetDataScorecard")
	return c.SetResultInfo(false, "", result)
}

func (c *MController) GetScorecard(k *knot.WebContext) interface{} {
	tk.Println("getscorecard")
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	result := []m.BusinessDriverL1Model{}

	csr, err := c.Ctx.Connection.NewQuery().From("BusinessDriverL1").Where(db.Eq("idx", "TBD1")).Cursor(nil)
	err = csr.Fetch(&result, 1, false)
	csr.Close()
	if err != nil {
		tk.Println(err)
	}
	tk.Println("Query OK getscorecard")
	return c.SetResultInfo(false, "", result)
}
func (c *MController) SaveScorecard(k *knot.WebContext) interface{} {
	parameter := struct {
		ScoreCardId string
		Name        string
		Description string
	}{}
	e := k.GetPayload(&parameter)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	sid := parameter.ScoreCardId

	if sid != "" {

		data := m.NewBusinessDriverL1Model()
		csr, err := c.Ctx.Connection.NewQuery().From("BusinessDriverL1").Where(db.Eq("_id", sid)).Cursor(nil)
		err = csr.Fetch(&data, 1, false)
		csr.Close()
		if err != nil {
			return nil
		}

		data.Name = parameter.Name
		data.Description = parameter.Description

		e = c.Ctx.Save(data)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
	} else {
		data := m.NewBusinessDriverL1Model()
		data.Name = parameter.Name
		data.Description = parameter.Description

		e = c.Ctx.Save(data)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
	}

	return c.SetResultInfo(false, "", nil)
}

func (c *MController) DeleteScorecard(k *knot.WebContext) interface{} {

	parameter := struct {
		ScoreCardId string
	}{}
	e := k.GetPayload(&parameter)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	// ScoreCardIdConv := bson.ObjectIdHex(parameter.ScoreCardId)
	sid := parameter.ScoreCardId

	data := m.NewBusinessDriverL1Model()

	csr, err := c.Ctx.Connection.NewQuery().From("BusinessDriverL1").Where(db.Eq("_id", sid)).Cursor(nil)
	err = csr.Fetch(&data, 1, false)
	csr.Close()
	if err != nil {
		return nil
	}

	// data.IsDeleted = true

	e = c.Ctx.Save(data)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	return c.SetResultInfo(false, "", data)
}

//Get_Scorecard
//Save_Scorecard
//Delete_Scorecard
