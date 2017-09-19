package controllers

import (
	m "eaciit/scb-apps/webapp/apps/web-cb/models"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// "strings"
	// "strconv"
	bson "gopkg.in/mgo.v2/bson"
	"time"
)

type BEFSponsorController struct {
	*BaseController
}

func (c *BEFSponsorController) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	// c.Action(k, "Open BEF Sponsor Master Page")
	c.Action(k, "BEF Sponsor", "Open BEF Sponsor Master Page", "", "", "", "", "")
	Scorecard := c.GetAccess(k, "SCORECARD")
	Initiative := c.GetAccess(k, "INITIATIVE")

	BEFSponsor := c.GetAccess(k, "BEFSPONSOR")
	k.Config.NoLog = true

	k.Config.LayoutTemplate = "_layout-v2.html"

	PartialFiles := []string{}
	PartialFiles = append(PartialFiles, "shared/sidebar.html")
	PartialFiles = append(PartialFiles, "shared/initiative_summary.html")
	PartialFiles = append(PartialFiles, "shared/initiative_chart.html")
	PartialFiles = append(PartialFiles, "shared/metric_upload.html")
	PartialFiles = append(PartialFiles, "dashboard/initiative.html")

	PartialFiles = append(PartialFiles, "befsponsor/form.html")

	k.Config.IncludeFiles = PartialFiles

	k.Config.OutputType = knot.OutputTemplate
	return tk.M{}.Set("BEFSponsor", BEFSponsor).Set("Scorecard", Scorecard).Set("Initiative", Initiative)
}
func (c *BEFSponsorController) GetPrototype(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true
	return c.SetResultInfo(false, "", m.NewBEFSponsor())
}

func (c *BEFSponsorController) GetData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		IsDeleted bool
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	result := []m.BEFSponsorModel{}

	csr, err := c.Ctx.Connection.NewQuery().From("BEFSponsor").Where(db.Eq("isdeleted", parm.IsDeleted)).Cursor(nil)
	err = csr.Fetch(&result, 0, false)
	csr.Close()
	if err != nil {
		return nil
	}

	return c.SetResultInfo(false, "", result)
}

func (c *BEFSponsorController) Get(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	result := new(m.BEFSponsorModel)

	csr, err := c.Ctx.Connection.NewQuery().From("BEFSponsor").Where(db.Contains("_id", "59015f508d5f771490100b34")).Cursor(nil)
	err = csr.Fetch(&result, 0, false)
	csr.Close()
	if err != nil {
		return nil
	}

	return c.SetResultInfo(false, "", result)
}

func (c *BEFSponsorController) Save(k *knot.WebContext) interface{} {

	parameter := struct {
		SponsorId string
		Name      string
	}{}
	e := k.GetPayload(&parameter)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	if parameter.SponsorId != "" {
		// c.Action(k, "Update BEF Sponsor")
		c.Action(k, "BEF Sponsor", "Update BEF Sponsor", "", "", "", "", "")

		SponsorIdConv := bson.ObjectIdHex(parameter.SponsorId)

		data := m.NewBEFSponsor()
		data.Id = SponsorIdConv
		csr, err := c.Ctx.Connection.NewQuery().From("BEFSponsor").Where(db.Eq("_id", SponsorIdConv)).Cursor(nil)
		err = csr.Fetch(&data, 1, false)
		csr.Close()
		if err != nil {
			return nil
		}

		data.Name = parameter.Name
		data.Updated_By = k.Session("username").(string)
		data.Updated_Date = time.Now().UTC()
		e = c.Ctx.Save(data)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
	} else {
		// c.Action(k, "Create New BEF Sponsor")
		c.Action(k, "BEF Sponsor", "Create New BEF Sponsor", "", "", "", "", "")

		data := m.NewBEFSponsor()
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

func (c *BEFSponsorController) Delete(k *knot.WebContext) interface{} {
	// c.Action(k, "Delete BEF Sponsor")
	c.Action(k, "BEF Sponsor", "Delete BEF Sponsor", "", "", "", "", "")

	parameter := struct {
		SponsorID string
	}{}
	e := k.GetPayload(&parameter)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	SponsorIdConv := bson.ObjectIdHex(parameter.SponsorID)

	data := m.NewBEFSponsor()
	data.Id = SponsorIdConv

	csr, err := c.Ctx.Connection.NewQuery().From("BEFSponsor").Where(db.Eq("_id", SponsorIdConv)).Cursor(nil)
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
