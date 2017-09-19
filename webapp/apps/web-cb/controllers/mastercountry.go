package controllers

import (
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type MasterCountryController struct {
	*BaseController
}

func (c *MasterCountryController) GetData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	ret := ResultInfo{}
	countries := make([]tk.M, 0)

	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id": "$Country",
	}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, err := c.Ctx.Connection.NewQuery().Command("pipe", pipes).
		From("Region").
		Cursor(nil)
	defer csr.Close()
	err = csr.Fetch(&countries, 0, true)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	ret.Data = countries
	return ret
}
