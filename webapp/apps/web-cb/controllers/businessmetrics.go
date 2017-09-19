package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/helper"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type BusinessMetricsController struct {
	*BaseController
}

func (c *BusinessMetricsController) GetData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	pipes := []tk.M{}
	unwind := tk.M{}
	match := tk.M{}

	result := make([]tk.M, 0)
	err := Deserialize(`
		{"$unwind":"$businessmetric"}
	`, &unwind)
	pipes = append(pipes, unwind)
	err = Deserialize(`
		{"$match":{"businessmetric.bdid":{"$ne":""}}}
	`, &match)
	pipes = append(pipes, match)

	csr, err := c.Ctx.Connection.NewQuery().Command("pipe", pipes).
		From("BusinessDriverL1").
		Cursor(nil)
	err = csr.Fetch(&result, 0, false)
	csr.Close()
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	return c.SetResultInfo(false, "", result)
}

func (c *BusinessMetricsController) GetAllData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	pipes := []tk.M{}
	unwind := tk.M{}
	project := tk.M{}

	result := make([]tk.M, 0)

	err := Deserialize(`
    {"$unwind": "$businessmetric"}
	`, &unwind)
	err = Deserialize(`
    {"$project": 
        {
        	"id" : "$businessmetric.id",
        	"description": "$businessmetric.description", 
            "MetricType":"$businessmetric.MetricType",
        	"DecimalFormat": "$businessmetric.DecimalFormat",
            "naactual": "$businessmetric.naactual",
            "natarget": "$businessmetric.natarget",
            "valuetype": "$businessmetric.valuetype",
            "type": "$businessmetric.type"
       	} 
    }
	`, &project)

	pipes = append(pipes, unwind)
	pipes = append(pipes, project)

	csr, err := c.Ctx.Connection.NewQuery().Command("pipe", pipes).
		From("BusinessDriverL1").
		Cursor(nil)
	err = csr.Fetch(&result, 0, false)
	csr.Close()
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	return c.SetResultInfo(false, "", result)
}

// func (c *BusinessMetricsController) GetDataPoint(k *knot.WebContext) interface{} {
// 	c.LoadBase(k)
// 	k.Config.OutputType = knot.OutputJson
// 	pipes := []tk.M{}
// 	unwind := tk.M{}
// 	project := tk.M{}

// 	result := make([]tk.M, 0)
// 	err := Deserialize(`
//     {"$unwind": "$businessmetric"}
// 	`, &unwind)
// 	err = Deserialize(`
//     {"$project":
//         {"id" : "$businessmetric.id","DataPoint": "$businessmetric.DataPoint"}
//     }
// 	`, &project)

// 	pipes = append(pipes, unwind)
// 	pipes = append(pipes, project)

// 	csr, err := c.Ctx.Connection.NewQuery().Command("pipe", pipes).
// 		From("BusinessDriverL1").
// 		Cursor(nil)
// 	err = csr.Fetch(&result, 0, false)
// 	csr.Close()
// 	if err != nil {
// 		return c.SetResultInfo(true, err.Error(), nil)
// 	}

// 	return c.SetResultInfo(false, "", result)
// }

func (c *BusinessMetricsController) GetUnassignedData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	pipes := []tk.M{}
	unwind := tk.M{}

	result := make([]tk.M, 0)
	err := Deserialize(`
		{"$unwind":"$businessmetric"}
	`, &unwind)
	pipes = append(pipes, unwind)
	csr, err := c.Ctx.Connection.NewQuery().Command("pipe", pipes).
		From("BusinessDriverL1").
		Cursor(nil)
	err = csr.Fetch(&result, 0, false)
	csr.Close()
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	return c.SetResultInfo(false, "", result)
}
