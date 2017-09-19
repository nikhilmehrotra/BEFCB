package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type MenuSettingController struct {
	*BaseController
}

func (c *MenuSettingController) Default(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	k.Config.NoLog = true
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

func (c *MenuSettingController) GetMenuTop(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	ret := ResultInfo{}
	payLoad := struct {
		Id string
	}{}
	err := k.GetPayload(&payLoad)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	var menuAccess []interface{}
	accesMenu := k.Session("roles").([]SysRolesModel)
	if accesMenu != nil {
		for _, o := range accesMenu[0].Menu {
			if o.Access == true {
				menuAccess = append(menuAccess, o.Menuid)
			}
		}
	}

	var dbFilter []*db.Filter
	dbFilter = append(dbFilter, db.Eq("Enable", true))
	dbFilter = append(dbFilter, db.In("_id", menuAccess...))

	queryTotal := tk.M{}
	query := tk.M{}
	data := make([]TopMenuModel, 0)
	total := make([]TopMenuModel, 0)
	retModel := tk.M{}

	if len(dbFilter) > 0 {
		query.Set("where", db.And(dbFilter...))
		queryTotal.Set("where", db.And(dbFilter...))
	}

	crsData, errData := c.Ctx.Find(NewTopMenuModel(), query)
	defer crsData.Close()
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}
	errData = crsData.Fetch(&data, 0, false)

	//	log.Printf("Data => %#v\n", len(data))
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	} else {
		retModel.Set("Records", data)
	}
	crsTotal, errTotal := c.Ctx.Find(NewTopMenuModel(), queryTotal)
	defer crsTotal.Close()
	if errTotal != nil {
		return c.SetResultInfo(true, errTotal.Error(), nil)
	}
	errTotal = crsTotal.Fetch(&total, 0, false)

	//	log.Printf("Total => %#v\n", len(total))
	if errTotal != nil {
		return c.SetResultInfo(true, errTotal.Error(), nil)
	} else {
		retModel.Set("Count", len(total))
	}
	ret.Data = retModel

	return ret
}

func (c *MenuSettingController) GetAccessMenu(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	url := k.Request.URL.String()
	accesMenu := k.Session("roles").([]SysRolesModel)

	access := []tk.M{}
	for _, o := range accesMenu[0].Menu {

		if o.Url == url {
			obj := tk.M{}
			obj.Set("view", o.View)
			obj.Set("create", o.Create)
			obj.Set("approve", o.Approve)
			obj.Set("delete", o.Delete)
			obj.Set("process", o.Process)
			obj.Set("edit", o.Edit)
			access = append(access, obj)
		}

	}

	return access
}

func (c *MenuSettingController) GetSelectMenu(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	ret := ResultInfo{}
	payLoad := struct {
		Id string
	}{}
	err := k.GetPayload(&payLoad)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	var dbFilter []*db.Filter
	if payLoad.Id != "" {
		dbFilter = append(dbFilter, db.Eq("_id", payLoad.Id))
	}
	//	log.Printf("QUERY=> %#v\n", query)
	queryTotal := tk.M{}
	query := tk.M{}
	data := make([]TopMenuModel, 0)
	total := make([]TopMenuModel, 0)
	retModel := tk.M{}

	if len(dbFilter) > 0 {
		query.Set("where", db.And(dbFilter...))
		queryTotal.Set("where", db.And(dbFilter...))
	}

	crsData, errData := c.Ctx.Find(NewTopMenuModel(), query)
	defer crsData.Close()
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}
	errData = crsData.Fetch(&data, 0, false)

	//	log.Printf("Data => %#v\n", len(data))
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	} else {
		retModel.Set("Records", data)
	}
	crsTotal, errTotal := c.Ctx.Find(NewTopMenuModel(), queryTotal)
	defer crsTotal.Close()
	if errTotal != nil {
		return c.SetResultInfo(true, errTotal.Error(), nil)
	}
	errTotal = crsTotal.Fetch(&total, 0, false)

	//	log.Printf("Total => %#v\n", len(total))
	if errTotal != nil {
		return c.SetResultInfo(true, errTotal.Error(), nil)
	} else {
		retModel.Set("Count", len(total))
	}
	ret.Data = retModel

	return ret
}

func (c *MenuSettingController) SaveMenuTop(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	payLoad := struct {
		Id        string
		PageId    string
		Parent    string
		Title     string
		Url       string
		IndexMenu int
		Enable    bool
	}{}
	err := k.GetPayload(&payLoad)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	mt := NewTopMenuModel()
	mt.Id = payLoad.Id
	mt.PageId = payLoad.PageId
	mt.Parent = payLoad.Parent
	mt.Title = payLoad.Title
	mt.Url = payLoad.Url
	mt.IndexMenu = payLoad.IndexMenu
	mt.Enable = payLoad.Enable
	c.Ctx.Save(mt)
	return c.SetResultInfo(false, "Menu has been successfully created.", nil)
}

func (c *MenuSettingController) DeleteMenuTop(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	payLoad := struct {
		Id string
	}{}

	err := k.GetPayload(&payLoad)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	u := NewTopMenuModel()
	err = c.Ctx.GetById(u, payLoad.Id)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	err = c.Ctx.Delete(u)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	return c.SetResultInfo(false, "Menu has been successfully created.", nil)
}

func (c *MenuSettingController) UpdateMenuTop(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	payLoad := struct {
		Id        string
		PageId    string
		Parent    string
		Title     string
		Url       string
		IndexMenu int
		Enable    bool
	}{}
	err := k.GetPayload(&payLoad)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	mt := NewTopMenuModel()
	mt.Id = payLoad.Id
	mt.PageId = payLoad.PageId
	mt.Parent = payLoad.Parent
	mt.Title = payLoad.Title
	mt.Url = payLoad.Url
	mt.IndexMenu = payLoad.IndexMenu
	mt.Enable = payLoad.Enable
	c.Ctx.Save(mt)
	return c.SetResultInfo(false, "Menu has been successfully update.", nil)
}
