package controllers

import (
	"eaciit/scb-apps/webapp/apps/web-cb/helper"
	// . "eaciit/scb-apps/webapp/apps/web-cb/models"
	// "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// "gopkg.in/mgo.v2/bson"
	// "os"
	// "path/filepath"
	// "sort"
	// "strconv"
	// "strings"
	"time"
)

type GenericCOntroller struct {
	*BaseController
}

func (c *GenericCOntroller) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.NoLog = true
	k.Config.LayoutTemplate = "_layout_dedicated.html"
	k.Config.IncludeFiles = []string{"shared/sidebar.html", "generic/form.html"}
	k.Config.OutputType = knot.OutputTemplate
	return nil
}

func (c *GenericCOntroller) GetLastPeriod(k *knot.WebContext) time.Time {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	pipes := []tk.M{}
	group := tk.M{}

	result := tk.M{}

	err := helper.Deserialize(`
		{
		    "$group" : {
		        "_id": null , 
		        "LastPeriod":{"$max":
		                     { "$cond": [ {"$and" : [ { "$ne": [ "$actualytd", 0.0] },
		                                               { "$ne": [ "$actualytd",130895111188.0] },
		                                               {"$ne":["$naactual",true]}
		                                     ] },
		                                     "$period",
		                                     null ] }
		                  }
		             }
		}
	`, &group)

	pipes = append(pipes, group)

	csr, err := c.Ctx.Connection.NewQuery().Command("pipe", pipes).
		From("BusinessMetricsData").
		Cursor(nil)
	err = csr.Fetch(&result, 1, true)
	csr.Close()
	if err != nil {
		tk.Println(err.Error())
	}
	period := time.Now().UTC()
	period = result.Get("LastPeriod").(time.Time)

	return period

}

func (c *GenericCOntroller) GetLastPeriodByBusinessMetrics(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	pipes := []tk.M{}
	group := tk.M{}

	// result := []tk.M{}

	err := helper.Deserialize(`
		{
		    "$group" : {
		        "_id": "$bmid" , 
		        "LastPeriod":{"$max":
		                     { "$cond": [ {"$and" : [ { "$ne": [ "$actualytd", 0.0] },
		                                               { "$ne": [ "$actualytd",130895111188.0] },
		                                               {"$ne":["$naactual",true]}
		                                     ] },
		                                     "$period",
		                                     null ] }
		                  }
		             }
		}
	`, &group)

	pipes = append(pipes, group)
	aggregateResult := []tk.M{}

	csr, err := c.Ctx.Connection.NewQuery().Command("pipe", pipes).
		From("BusinessMetricsData").
		Cursor(nil)
	err = csr.Fetch(&aggregateResult, 0, true)
	csr.Close()
	if err != nil {
		tk.Println(err.Error())
	}

	// result.Set()

	return c.SetResultInfo(false, "", aggregateResult)

}
