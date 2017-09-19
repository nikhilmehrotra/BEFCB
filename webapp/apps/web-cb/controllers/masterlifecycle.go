package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"fmt"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
)

func (m *MController) LifeCycle(k *knot.WebContext) interface{} {
	access := m.LoadBase(k)
	k.Config.NoLog = true
	k.Config.IncludeFiles = []string{}
	k.Config.OutputType = knot.OutputTemplate
	DataAccess := Previlege{}

	for _, o := range access {
		DataAccess.Create = o["Create"].(bool)
		DataAccess.View = o["View"].(bool)
		DataAccess.Delete = o["Delete"].(bool)
		DataAccess.Process = o["Process"].(bool)
		DataAccess.Delete = o["Delete"].(bool)
		DataAccess.Edit = o["Edit"].(bool)
		DataAccess.Menuid = o["Menuid"].(string)
		DataAccess.Menuname = o["Menuname"].(string)
		DataAccess.Approve = o["Approve"].(bool)
		DataAccess.Username = o["Username"].(string)
	}

	return DataAccess
}

func (m *MController) GetLifeCycleData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	csr, e := m.Ctx.Find(new(LifeCycleModel), tk.M{})
	if e != nil {
		return m.ErrorResultInfo(e.Error(), nil)
	}
	result := []LifeCycleModel{}
	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return m.ErrorResultInfo(e.Error(), nil)
	}
	return m.SetResultInfo(false, "", result)
}

func (m *MController) DeleteLiveCircle(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		m.WriteLog(e)
	}

	result := new(LifeCycleModel)
	e = m.Ctx.GetById(result, p.Id)
	if e != nil {
		m.WriteLog(e)
	}

	e = m.Ctx.Delete(result)

	return ""
}

func (m *MController) SaveLiveCircle(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := LifeCycleModel{}

	e := k.GetPayload(&p)
	if e != nil {
		fmt.Println(e)
		return e
	}

	if p.Id == "" {
		p.Id = bson.NewObjectId()
	}

	e = m.Ctx.Save(&p)
	if e != nil {
		fmt.Println(e)
		return e
	}

	return ""

}
