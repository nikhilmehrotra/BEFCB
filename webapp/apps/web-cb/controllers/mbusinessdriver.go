package controllers

import (
	// . "eaciit/crm/commons"
	// . "eaciit/crm/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"fmt"
	// "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	// "strings"
	// "time"
)

func (m *MController) BusinessDriver(k *knot.WebContext) interface{} {
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

func (m *MController) GetBusinessDriverList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type xsorting struct {
		Field string
		Dir   string
	}
	p := struct {
		Skip int
		Take int
		Sort []xsorting
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		m.WriteLog(e)
	}

	// csr, e := m.Ctx.Find(new(BusinessDriverModel), toolkit.M{})
	csr, e := m.Ctx.Find(new(BusinessDriverModel), toolkit.M{}.Set("order", []string{"Seq"}).Set("skip", p.Skip).Set("limit", p.Take))
	defer csr.Close()

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, "-"+p.Sort[0].Field)
			} else {
				arrsort = append(arrsort, p.Sort[0].Field)
			}
		}
		csr, e = m.Ctx.Find(new(BusinessDriverModel), toolkit.M{}.Set("order", arrsort).Set("skip", p.Skip).Set("limit", p.Take))
		defer csr.Close()

		// fmt.Println(p.Sort[0].Field, p.Sort[0].Dir)
	}
	results := make([]BusinessDriverModel, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error()
	}
	results2 := results

	query := toolkit.M{}.Set("AGGR", "$sum")
	csr, err := m.Ctx.Find(new(BusinessDriverModel), query)
	defer csr.Close()
	if err != nil {
		return err.Error()
	}

	data := struct {
		Data  []BusinessDriverModel
		Total int
	}{
		Data:  results2,
		Total: csr.Count(),
	}
	// fmt.Println(data)
	return data
}

func (m *MController) GetBusinessDriver(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		m.WriteLog(e)
	}

	result := new(BusinessDriverModel)
	e = m.Ctx.GetById(result, p.Id)
	if e != nil {
		m.WriteLog(e)
	}

	return result
}

func (m *MController) DeleteBusinessDriver(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		m.WriteLog(e)
	}

	result := new(BusinessDriverModel)
	e = m.Ctx.GetById(result, p.Id)
	if e != nil {
		m.WriteLog(e)
	}

	e = m.Ctx.Delete(result)

	return ""

	// result2 := new(CustomerModel)
	// query := toolkit.M{}.Set("where", dbox.Eq("industry._id", p.Id))
	// csr, _ := m.Ctx.Find(result2, query)
	// defer csr.Close()
	// resultcustomer := make([]CustomerModel, 0)
	// csr.Fetch(&resultcustomer, 0, false)

	// if len(resultcustomer) != 0 {
	// 	return "This data are reference to another data"
	// } else {
	// 	e = m.Ctx.Delete(result)

	// 	return ""
	// }
}

func (m *MController) SaveBusinessDriver(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := BusinessDriverModel{}

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
