package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"github.com/eaciit/knot/knot.v1"
)

type BusinessDriverController struct {
	*BaseController
}

func (c *BusinessDriverController) GetData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	result := make([]SummaryBusinessDriverModel, 0)

	csr, err := c.Ctx.Connection.NewQuery().From("SummaryBusinessDriver").
		Cursor(nil)
	defer csr.Close()
	err = csr.Fetch(&result, 0, true)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	return c.SetResultInfo(false, "", result)
}
