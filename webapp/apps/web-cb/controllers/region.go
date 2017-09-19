package controllers

import (
	m "eaciit/scb-apps/webapp/apps/web-cb/models"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	// tk "github.com/eaciit/toolkit"
	// "strings"
	// "strconv"
	bson "gopkg.in/mgo.v2/bson"
	// "time"
)

type RegionController struct {
	*BaseController
}

func (c *RegionController) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.NoLog = true
	k.Config.LayoutTemplate = "_layout_dedicated.html"
	k.Config.IncludeFiles = []string{"shared/sidebar.html"}
	k.Config.OutputType = knot.OutputTemplate
	return nil
}
func (c *RegionController) GetPrototype(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true
	return c.SetResultInfo(false, "", m.NewRegion())
}

func (c *RegionController) GetData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		IsDeleted bool
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	result := []m.RegionModel{}

	csr, err := c.Ctx.Connection.NewQuery().From("Region").Cursor(nil) //.Where(db.Eq("isdeleted", parm.IsDeleted))
	err = csr.Fetch(&result, 0, false)
	csr.Close()
	if err != nil {
		return nil
	}

	return c.SetResultInfo(false, "", result)
}

func (c *RegionController) Get(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	result := []m.RegionModel{}

	csr, err := c.Ctx.Connection.NewQuery().From("Region").Where(db.Eq("_id", "580649eae079bb325cd36179")).Cursor(nil)
	err = csr.Fetch(&result, 1, false)
	csr.Close()
	if err != nil {
		return nil
	}

	return c.SetResultInfo(false, "", result)
}

func (c *RegionController) Save(k *knot.WebContext) interface{} {
	parameter := struct {
		RegionId     string
		Major_Region string
		Region       string
		Country      string
		CountryCode  string
	}{}
	e := k.GetPayload(&parameter)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	if parameter.RegionId != "" {

		RegionIdConv := bson.ObjectIdHex(parameter.RegionId)

		data := m.NewRegion()
		data.Id = RegionIdConv
		csr, err := c.Ctx.Connection.NewQuery().From("Region").Where(db.Eq("_id", RegionIdConv)).Cursor(nil)
		err = csr.Fetch(&data, 1, false)
		csr.Close()
		if err != nil {
			return nil
		}

		data.Major_Region = parameter.Major_Region
		data.Region = parameter.Region
		data.Country = parameter.Country
		data.CountryCode = parameter.CountryCode

		e = c.Ctx.Save(data)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
	} else {
		data := m.NewRegion()
		data.Major_Region = parameter.Major_Region
		data.Region = parameter.Region
		data.Country = parameter.Country
		data.CountryCode = parameter.CountryCode

		e = c.Ctx.Save(data)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
	}

	return c.SetResultInfo(false, "", nil)
}

func (c *RegionController) Delete(k *knot.WebContext) interface{} {

	parameter := struct {
		RegionId string
	}{}
	e := k.GetPayload(&parameter)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	RegionIdConv := bson.ObjectIdHex(parameter.RegionId)

	data := m.NewRegion()
	data.Id = RegionIdConv

	csr, err := c.Ctx.Connection.NewQuery().From("Region").Where(db.Eq("_id", RegionIdConv)).Cursor(nil)
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
