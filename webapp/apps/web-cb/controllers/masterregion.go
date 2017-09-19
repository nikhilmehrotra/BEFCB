package controllers

import (
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type MasterRegion struct {
	*BaseController
}

func (c *MasterRegion) GetData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	ret := ResultInfo{}
	regionList := make([]tk.M, 0)

	csr, err := c.Ctx.Connection.NewQuery().From("Region").Order("Country").
		Cursor(nil)
	defer csr.Close()
	err = csr.Fetch(&regionList, 0, true)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	ret.Data = regionList
	return ret
}
