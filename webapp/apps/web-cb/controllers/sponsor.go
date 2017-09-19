package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type SponsorController struct {
	*BaseController
}

func (c *SponsorController) GetData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	result := make([]SponsorModel, 0)

	csr, err := c.Ctx.Connection.NewQuery().From(new(SponsorModel).TableName()).
		Cursor(nil)
	defer csr.Close()
	err = csr.Fetch(&result, 0, true)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	return c.SetResultInfo(false, "", result)
}

func (c *SponsorController) InitData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	initiateData := []string{
		"Humphrey, Mike",
		"Marrs, Anna",
		"Thomas, Alistair",
		"Murray, Stuart",
		"Phebey, Tom",
		"Arora, Jiten",
		"Siva, Ve",
		"Walker, Mark",
		"Pandey, Prachit",
		"Rathnam, Venkatesh"}
	result := tk.M{}
	for _, i := range initiateData {
		d := NewBEFSponsor()
		d.Name = i
		e := c.Ctx.Save(d)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), "")
		}
	}
	return c.SetResultInfo(false, "", result)
}
