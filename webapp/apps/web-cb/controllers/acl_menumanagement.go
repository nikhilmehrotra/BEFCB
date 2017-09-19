package controllers

import (
	// helper "eaciit/scb-apps/webapp/apps/web-cb/helper"
	// m "eaciit/scb-apps/webapp/apps/web-cb/models"
	"github.com/eaciit/acl/v2.0"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

func (c *AclController) GetMenuList(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	result := []tk.M{}
	// c.AclAclCtx.Find(, parms)
	d := new(acl.Access)
	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Order("index").Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	return c.SetResultInfo(false, "", result)
}

func (c *AclController) GetMenuManagementReferences(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	result := tk.M{}
	result.Set("MenuData", new(acl.Access))
	return c.SetResultInfo(false, "", result)
}

func (c *AclController) SaveMenu(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Mode string
		Data acl.Access
	}{}

	err := k.GetPayload(&parm)
	// c.Action(k, "Save Menu Management")
	c.Action(k, "Menu Management", "Save Menu Management", "", "", "", "", "")
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	data := new(acl.Access)
	if parm.Mode == "create" {
		data.ID = parm.Data.ID
	} else {
		// tk.Println(parm.Data.ID)
		csr, err := c.AclCtx.Connection.NewQuery().From(data.TableName()).Where(db.Eq("id", parm.Data.ID)).Cursor(nil)
		err = csr.Fetch(&data, 1, false)
		csr.Close()
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
	}

	data.ParentId = parm.Data.ParentId
	data.ID = parm.Data.ID
	data.Title = parm.Data.Title
	data.Url = parm.Data.Url
	data.Index = parm.Data.Index
	data.Enable = parm.Data.Enable
	data.Category = parm.Data.Category
	data.SpecialAccess1 = parm.Data.SpecialAccess1
	data.SpecialAccess2 = parm.Data.SpecialAccess2
	data.SpecialAccess3 = parm.Data.SpecialAccess3
	data.SpecialAccess4 = parm.Data.SpecialAccess4

	err = acl.Save(data)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	return c.SetResultInfo(false, "", nil)
}

func (c *AclController) GetMenu(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	parm := struct {
		Id string
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	result := new(acl.Access)
	d := new(acl.Access)
	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.Eq("id", parm.Id)).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&result, 1, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	return c.SetResultInfo(false, "", result)

}

func (c *AclController) RemoveMenu(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	// c.Action(k, "Delete Menu Management")
	c.Action(k, "Menu Management", "Delete Menu Management", "", "", "", "", "")

	k.Config.OutputType = knot.OutputJson
	// data := tk.M{}
	parm := struct {
		Id string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	data := new(acl.Access)

	csr, err := c.AclCtx.Connection.NewQuery().From(data.TableName()).Where(db.Eq("_id", parm.Id)).Cursor(nil)
	err = csr.Fetch(&data, 1, false)
	csr.Close()
	if err != nil {
		return nil
	}

	err = acl.Delete(data)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	return c.SetResultInfo(false, "", data)
}
