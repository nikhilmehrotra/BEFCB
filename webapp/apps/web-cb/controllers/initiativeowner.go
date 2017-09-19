package controllers

import (
	m "eaciit/scb-apps/webapp/apps/web-cb/models"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// "strings"
	"github.com/eaciit/acl/v2.0"
	// "strconv"
	bson "gopkg.in/mgo.v2/bson"
	"time"
)

type InitiativeOwnerController struct {
	*BaseController
}

func (c *InitiativeOwnerController) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	// c.Action(k, "Open Initiative Owner Master Page")
	c.Action(k, "Initiative Owner", "Open Initiative Owner Master Page", "", "", "", "", "")

	InitiativeOwner := c.GetAccess(k, "INITIATIVEOWNER")
	Scorecard := c.GetAccess(k, "SCORECARD")
	Initiative := c.GetAccess(k, "INITIATIVE")

	k.Config.NoLog = true

	k.Config.LayoutTemplate = "_layout-v2.html"

	PartialFiles := []string{}
	PartialFiles = append(PartialFiles, "shared/sidebar.html")
	PartialFiles = append(PartialFiles, "shared/initiative_summary.html")
	PartialFiles = append(PartialFiles, "shared/initiative_chart.html")
	PartialFiles = append(PartialFiles, "shared/metric_upload.html")
	PartialFiles = append(PartialFiles, "dashboard/initiative.html")

	k.Config.IncludeFiles = PartialFiles

	k.Config.OutputType = knot.OutputTemplate
	return tk.M{}.Set("InitiativeOwner", InitiativeOwner).Set("Scorecard", Scorecard).Set("Initiative", Initiative)
}
func (c *InitiativeOwnerController) GetPrototype(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true
	return c.SetResultInfo(false, "", new(m.InitiativeOwnerModel))
}

func (c *InitiativeOwnerController) GetData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		IsDeleted bool
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	aclGroup := new(acl.Group)
	res := []acl.Group{}
	userResult := []acl.User{}
	csr, err := c.AclCtx.Connection.NewQuery().From(aclGroup.TableName()).Where(db.Eq("ispartofinitiativeowner", true)).Cursor(nil)
	err = csr.Fetch(&res, 0, false)
	csr.Close()
	if err != nil {
		return nil
	}
	grup := []interface{}{}
	for _, x := range res {
		grup = append(grup, x.ID)
	}

	aclUser := new(acl.User)
	csr, err = c.AclCtx.Connection.NewQuery().From(aclUser.TableName()).Where(db.In("groups", grup...)).Cursor(nil)
	err = csr.Fetch(&userResult, 0, false)
	csr.Close()
	if err != nil {
		return nil
	}

	// tk.Println(userResult)
	// tk.Println(grup)

	// result := []m.InitiativeOwnerModel{}

	// csr, err = c.Ctx.Connection.NewQuery().From(new(m.InitiativeOwnerModel).TableName()).Where(db.Eq("isdeleted", parm.IsDeleted)).Cursor(nil)
	// err = csr.Fetch(&result, 0, false)
	// csr.Close()
	// if err != nil {
	// 	return nil
	// }

	return c.SetResultInfo(false, "", userResult)
}

func (c *InitiativeOwnerController) Get(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	result := new(m.InitiativeOwnerModel)

	csr, err := c.Ctx.Connection.NewQuery().From(new(m.InitiativeOwnerModel).TableName()).Where(db.Contains("_id", "59015f508d5f771490100b34")).Cursor(nil)
	err = csr.Fetch(&result, 0, false)
	csr.Close()
	if err != nil {
		return nil
	}

	return c.SetResultInfo(false, "", result)
}

func (c *InitiativeOwnerController) Save(k *knot.WebContext) interface{} {

	parameter := struct {
		InitiativeOwnerID string
		Name              string
	}{}
	e := k.GetPayload(&parameter)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	user := new(acl.User)
	csr, err := c.Ctx.Connection.NewQuery().From(user.TableName()).Where(db.Eq("fullname", parameter.Name)).Cursor(nil)
	err = csr.Fetch(&user, 1, false)
	csr.Close()
	if err != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	if parameter.InitiativeOwnerID != "" {
		// c.Action(k, "Update Initiative Owner")
		c.Action(k, "Initiative Owner", "Update Initiative Owner", "", "", "", "", "")

		InitiativeOwnerIDConv := bson.ObjectIdHex(parameter.InitiativeOwnerID)

		data := m.NewInitiativeOwner()
		data.Id = InitiativeOwnerIDConv
		csr, err := c.Ctx.Connection.NewQuery().From(data.TableName()).Where(db.Eq("_id", InitiativeOwnerIDConv)).Cursor(nil)
		err = csr.Fetch(&data, 1, false)
		csr.Close()
		if err != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}

		data.LoginID = user.LoginID
		data.Name = parameter.Name
		data.Updated_By = k.Session("username").(string)
		data.Updated_Date = time.Now().UTC()
		e = c.Ctx.Save(data)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
	} else {
		c.Action(k, "Initiative Owner", "Create New Initiative Owner", "", "", "", "", "")

		data := m.NewInitiativeOwner()
		data.LoginID = user.LoginID
		data.Name = parameter.Name
		data.Created_By = k.Session("username").(string)
		data.Created_Date = time.Now().UTC()
		e = c.Ctx.Save(data)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
	}

	return c.SetResultInfo(false, "", nil)
}

func (c *InitiativeOwnerController) Delete(k *knot.WebContext) interface{} {
	// c.Action(k, "Delete Initiative Owner")
	c.Action(k, "Initiative Owner", "Delete Initiative Owner", "", "", "", "", "")

	parameter := struct {
		InitiativeOwnerID string
	}{}
	e := k.GetPayload(&parameter)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	InitiativeOwnerIDConv := bson.ObjectIdHex(parameter.InitiativeOwnerID)

	data := m.NewInitiativeOwner()
	data.Id = InitiativeOwnerIDConv

	csr, err := c.Ctx.Connection.NewQuery().From(new(m.InitiativeOwnerModel).TableName()).Where(db.Eq("_id", InitiativeOwnerIDConv)).Cursor(nil)
	err = csr.Fetch(&data, 1, false)
	csr.Close()
	if err != nil {
		return nil
	}

	data.IsDeleted = true

	e = c.Ctx.Save(data)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	return c.SetResultInfo(false, "", data)
}
