package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type SharedAgendaController struct {
	*BaseController
}

func (c *SharedAgendaController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	businessDriverList := make([]BusinessDriverL1Model, 0)
	csr, e := c.Ctx.Connection.NewQuery().From(new(BusinessDriverL1Model).TableName()).Where(dbox.Ne("idx", "TBD3")).Order("seq").Cursor(nil)
	e = csr.Fetch(&businessDriverList, 0, false)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	csr.Close()

	result := []tk.M{}
	for _, i := range businessDriverList {
		data := tk.M{}
		e = tk.StructToM(i, &data)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}

		// Get Business Driver as per SCID
		BusinessDriverList := make([]SummaryBusinessDriverModel, 0)
		csr, e = c.Ctx.Connection.NewQuery().From(new(SummaryBusinessDriverModel).TableName()).Where(dbox.Eq("parentid", i.Idx)).Order("Seq").Cursor(nil)
		e = csr.Fetch(&BusinessDriverList, 0, false)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
		csr.Close()
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
		data.Set("BusinessDriverList", BusinessDriverList)

		// Get Shared Agenda Driver as per SCID
		SharedAgendaDrivers := make([]SharedAgendaModel, 0)
		csr, e = c.Ctx.Find(new(SharedAgendaModel), tk.M{}.Set("where", dbox.Eq("scid", i.Idx)))
		e = csr.Fetch(&SharedAgendaDrivers, 0, false)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
		csr.Close()
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}

		ResSharedAgendaDrivers := make([]SharedAgendaModel, 0)
		// Mapping Shared Agenda as per BD Index
		IsAny := false
		for _, bd := range BusinessDriverList {
			IsAny = false
			for _, sa := range SharedAgendaDrivers {
				if bd.Idx == sa.BDId {
					IsAny = true
					ResSharedAgendaDrivers = append(ResSharedAgendaDrivers, sa)
				}
			}
			if !IsAny {
				tempSA := SharedAgendaModel{}
				tempSA.Id = bson.NewObjectId()
				tempSA.SCId = i.Idx
				tempSA.SCName = i.Name
				tempSA.BDId = bd.Idx
				tempSA.Leads = []string{}
				tempSA.Seq = -1
				ResSharedAgendaDrivers = append(ResSharedAgendaDrivers, tempSA)
			}
		}
		if len(ResSharedAgendaDrivers) > 0 {
			data.Set("SharedAgendaDrivers", ResSharedAgendaDrivers)
			result = append(result, data)
		}

	}
	return c.SetResultInfo(false, "", result)
}

func (c *SharedAgendaController) Save(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		SharedAgendaList []SharedAgendaModel
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	SharedAgendaDrivers := make([]SharedAgendaModel, 0)
	csr, e := c.Ctx.Find(new(SharedAgendaModel), nil)
	e = csr.Fetch(&SharedAgendaDrivers, 0, false)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	SADrivers := map[string]SharedAgendaModel{}
	for _, x := range SharedAgendaDrivers {
		SADrivers[x.BDId] = x
	}
	e = c.Ctx.DeleteMany(new(SharedAgendaModel), nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	AnyChanges := false
	for _, i := range parm.SharedAgendaList {
		i.Id = bson.NewObjectId()
		i.CreatedDate = time.Now()
		i.CreatedBy = k.Session("username").(string)
		eData := SADrivers[i.BDId]
		if i.Name != eData.Name {
			AnyChanges = true
			c.Action(k, "Shared Agenda", "Update Shared Agenda", "Shared Agenda Drivers", eData.Name, i.Name, "", "")
		}
		oldLeads := strings.Join(eData.Leads, ",")
		newLeads := strings.Join(i.Leads, ",")
		if oldLeads != newLeads {
			AnyChanges = true
			c.Action(k, "Shared Agenda", "Update Shared Agenda", "Shared Agenda Drivers", oldLeads, newLeads, "", "")
		}
		if i.RAG != eData.RAG {
			AnyChanges = true
			c.Action(k, "Shared Agenda", "Update Shared Agenda", "Shared Agenda Drivers", eData.RAG, i.RAG, "", "")
		}
		e = c.Ctx.Save(&i)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
	}
	if AnyChanges {
		c.Action(k, "Shared Agenda", "Save Changes for Shared Agenda", "", "", "", "", "")
	}
	return c.SetResultInfo(false, "", nil)
}
