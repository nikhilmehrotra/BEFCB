package controllers

import (
	"eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	// "strings"
	"time"
)

type TaskController struct {
	*BaseController
}

func (c *TaskController) Remove(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Id string
	}{}
	e := k.GetPayload(&parm)
	tk.Println(parm.Id)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	tk.Println(parm.Id)
	e = c.Ctx.DeleteMany(new(TaskModel), dbox.Eq("_id", bson.ObjectIdHex(parm.Id)))
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	return c.SetResultInfo(false, "", nil)
}
func (c *TaskController) Save(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Mode             string
		Id               string
		LifeCycleId      string
		BusinessDriverId string
		Name             string
		Owner            string
		Statement        string
		Description      string
		TaskType         string
		SCCategory       string
		IsGlobal         bool
		Region           []string
		Country          []string
		Map              []tk.M
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), "")
	}
	for _, i := range parm.Map {
		d := NewTaskModel()
		d.Id = bson.NewObjectId()
		d.LifeCycleId = i.Get("LifeCycle").(string)
		d.BusinessDriverId = i.Get("BusinessDriver").(string)
		d.Name = parm.Name
		d.Owner = parm.Owner
		d.Statement = parm.Statement
		d.Description = parm.Description
		d.DateCreated = time.Now()
		d.TaskType = parm.TaskType
		d.IsGlobal = parm.IsGlobal
		d.Region = parm.Region
		d.Country = parm.Country
		d.SCCategory = i.Get("SCCategory").(string)

		if d.BusinessDriverId == "BD6" {
			d.TaskType = "SupportingEnablers"
		} else {
			d.TaskType = "KeyEnablers"
		}

		if parm.Mode == "edit" {
			e := c.Ctx.DeleteMany(d, dbox.Eq("_id", bson.ObjectIdHex(parm.Id)))
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
		}
		err = c.Ctx.Save(d)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), "")
		}
	}
	return helper.CreateResult(true, nil, "")

}

func (c *TaskController) MoveUpdate(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	p := TaskModel{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}
	tk.Println(p)
	//c.NewQuery().Delete().From(new(GenDataBrowserNotInTmp).TableName()).Where(dbox.Eq("ID", genIDTempTable)).Exec(nil)

	c.Ctx.Connection.NewQuery().Delete().From("Task").Where(dbox.Eq("_id", p.Id)).Exec(nil)

	p.Id = bson.NewObjectId()
	e = c.Ctx.Save(&p)
	if e != nil {
		c.WriteLog(e)
	}

	result := tk.M{}

	result.Set("Res", "OK")

	return result
}
