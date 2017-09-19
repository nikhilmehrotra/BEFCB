package controllers

import (
	// helper "eaciit/scb-apps/webapp/apps/web-cb/helper"
	m "eaciit/scb-apps/webapp/apps/web-cb/models"
	"strings"
	"time"

	"github.com/eaciit/acl/v2.0"
	db "github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
	bson "gopkg.in/mgo.v2/bson"
)

func (c *AclController) GetUserList(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	result := []tk.M{}
	// c.AclAclCtx.Find(, parms)
	d := new(acl.User)
	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.Startwith("groups", "CB_")).Cursor(nil)
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

func (c *AclController) GetUserManagementReferences(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	result := tk.M{}

	// Get Group List
	group_list := []tk.M{}
	d := new(acl.Group)
	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.Startwith("_id", "CB_")).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&group_list, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	// Get Region list
	region_list := make([]tk.M, 0)
	csr, e = c.Ctx.Connection.NewQuery().From("Region").Order("Country").
		Cursor(nil)
	e = csr.Fetch(&region_list, 0, true)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	result.Set("GroupList", group_list)
	result.Set("RegionList", region_list)
	result.Set("UserData", new(acl.User))

	return c.SetResultInfo(false, "", result)
}

func (c *AclController) GetUser(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	parm := struct {
		Id string
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	result := tk.M{}

	csr, e := c.AclCtx.Connection.NewQuery().From("acl_users").Where(db.Eq("id", parm.Id)).Cursor(nil)
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

func (c *AclController) SaveUser(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Mode string
		Data acl.User
	}{}

	err := k.GetPayload(&parm)

	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	data := new(acl.User)
	auditdata := m.NewUserAuditTrailModel()
	userlogin := k.Session("username").(string)
	Olddata := new(acl.User)

	if parm.Mode == "create" {
		data.ID = tk.RandomString(32)
		auditdata.UserID = userlogin
		auditdata.SessionID = userlogin
		auditdata.ActionTime = time.Now().UTC()
		auditdata.UserIDChanged = parm.Data.LoginID
		auditdata.TypeOfChange = "Add User"
	} else {
		csr, err := c.AclCtx.Connection.NewQuery().From(data.TableName()).Where(db.Eq("_id", parm.Data.ID)).Cursor(nil)
		err = csr.Fetch(&data, 1, false)
		csr.Close()
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
		auditdata.UserID = userlogin
		auditdata.SessionID = userlogin
		auditdata.ActionTime = time.Now().UTC()
		auditdata.UserIDChanged = parm.Data.LoginID
		auditdata.TypeOfChange = "Update User"

		csr, err = c.AclCtx.Connection.NewQuery().From(data.TableName()).Where(db.Eq("_id", parm.Data.ID)).Cursor(nil)
		err = csr.Fetch(&Olddata, 1, false)
		csr.Close()
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}

	}

	// tk.Println(parm.Data.LoginType)
	if parm.Data.LoginType == 1 {
		data.LoginType = acl.LogTypeLdap
	} else {
		data.LoginType = acl.LogTypeBasic
	}

	data.LoginID = parm.Data.LoginID
	data.FullName = parm.Data.FullName
	data.FirstName = parm.Data.FirstName
	data.LastName = parm.Data.LastName
	data.Email = parm.Data.Email
	data.Groups = parm.Data.Groups
	data.Enable = parm.Data.Enable
	data.CountryCode = parm.Data.CountryCode
	data.Country = parm.Data.Country
	if data.LoginID != "" && parm.Data.LoginID != Olddata.LoginID {
		auditdata.Id = bson.NewObjectId()
		auditdata.FieldChanged = "UserName"
		auditdata.NewValue = parm.Data.LoginID
		auditdata.OldValue = Olddata.LoginID
		err = c.AclCtx.Save(auditdata)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}

	}

	if data.FullName != "" && parm.Data.FullName != Olddata.FullName {
		auditdata.Id = bson.NewObjectId()
		auditdata.FieldChanged = "FullName"
		auditdata.NewValue = parm.Data.FullName
		auditdata.OldValue = Olddata.FullName
		err = c.AclCtx.Save(auditdata)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}

	}

	if data.Email != "" && parm.Data.Email != Olddata.Email {
		auditdata.Id = bson.NewObjectId()
		auditdata.FieldChanged = "Email"
		auditdata.NewValue = parm.Data.Email
		auditdata.OldValue = Olddata.Email
		err = c.AclCtx.Save(auditdata)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}

	}

	OldGroup := Olddata.Groups
	NewGroup := parm.Data.Groups
	anychanges := false
	olddatastr := strings.Join(OldGroup, "|")

	// tk.Println("OLDGROP", OldGroup)
	// tk.Println("olddatastr", olddatastr)
	// tk.Println("NewGroup", NewGroup)

	for _, i := range NewGroup {
		if strings.Index(olddatastr, i) <= -1 {
			anychanges = true
		}

	}

	if anychanges == true {
		auditdata.Id = bson.NewObjectId()
		auditdata.FieldChanged = "User Role"
		auditdata.NewValue = strings.Join(parm.Data.Groups, ",")
		auditdata.OldValue = olddatastr
		err = c.AclCtx.Save(auditdata)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}

	}

	err = acl.Save(data)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	if parm.Mode == "create" {
		if data.LoginID != "" && parm.Data.LoginID != Olddata.LoginID {
			c.Action(k, "User Management", "Create New User", "LoginID", Olddata.LoginID, parm.Data.LoginID, "", "")
		}
		if data.FullName != "" && parm.Data.FullName != Olddata.FullName {
			c.Action(k, "User Management", "Create New User", "FullName", Olddata.FullName, parm.Data.FullName, "", "")
		}
		if data.Email != "" && parm.Data.Email != Olddata.Email {
			c.Action(k, "User Management", "Create New User", "Email", Olddata.Email, parm.Data.Email, "", "")
		}
	} else {
		if data.LoginID != "" && parm.Data.LoginID != Olddata.LoginID {
			c.Action(k, "User Management", "Update User", "LoginID", Olddata.LoginID, parm.Data.LoginID, "", "")
		}
		if data.FullName != "" && parm.Data.FullName != Olddata.FullName {
			c.Action(k, "User Management", "Update User", "FullName", Olddata.FullName, parm.Data.FullName, "", "")
		}
		if data.Email != "" && parm.Data.Email != Olddata.Email {
			c.Action(k, "User Management", "Update User", "Email", Olddata.Email, parm.Data.Email, "", "")
		}
		if parm.Data.Country != Olddata.Country {
			OldDatas := ""
			NewDatas := ""

			if Olddata.Country == "" {
				OldDatas = "GLOBAL"
			}

			if parm.Data.Country == "" {
				NewDatas = "GLOBAL"
			}

			c.Action(k, "User Management", "Update User", "Country", OldDatas, NewDatas, "", "")
		}
		if parm.Data.Enable != Olddata.Enable {
			newvalue := "DISABLED"
			oldvalue := "DISABLED"
			if parm.Data.Enable {
				newvalue = "ENABLED"
			}
			if Olddata.Enable {
				oldvalue = "ENABLED"
			}
			c.Action(k, "User Management", "Update User", "Status", newvalue, oldvalue, "", "")
		}
		if anychanges {
			groups := map[string]string{}
			oldrole := []string{}
			newrole := []string{}

			// Get Group List
			group_list := []acl.Group{}
			d := new(acl.Group)
			csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.Startwith("_id", "CB_")).Cursor(nil)
			if e != nil {
				return c.ErrorResultInfo(e.Error(), nil)
			}
			e = csr.Fetch(&group_list, 0, false)
			csr.Close()
			if e != nil {
				return c.ErrorResultInfo(e.Error(), nil)
			}
			for _, x := range group_list {
				groups[x.ID] = x.Title
			}
			for _, x := range Olddata.Groups {
				oldrole = append(oldrole, groups[x])
			}
			for _, x := range parm.Data.Groups {
				newrole = append(newrole, groups[x])
			}
			oldvalue := strings.Join(oldrole, ",")
			newvalue := strings.Join(newrole, ",")
			c.Action(k, "User Management", "Update User", "Role", oldvalue, newvalue, "", "")
		}
	}

	if strings.Trim(parm.Data.Password, "") != "" {
		err = acl.ChangePassword(data.ID, parm.Data.Password)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
	}

	return c.SetResultInfo(false, "", nil)
}

func (c *AclController) RemoveUser(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	// c.Action(k, "Remove User Management")
	c.Action(k, "User Management", "Remove User Management", "", "", "", "", "")

	k.Config.OutputType = knot.OutputJson
	// data := tk.M{}
	parm := struct {
		Id string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	data := new(acl.User)

	csr, err := c.AclCtx.Connection.NewQuery().From(data.TableName()).Where(db.Eq("_id", parm.Id)).Cursor(nil)
	err = csr.Fetch(&data, 1, false)
	csr.Close()
	if err != nil {
		return nil
	}

	auditdata := m.NewUserAuditTrailModel()
	userlogin := k.Session("username").(string)

	auditdata.UserID = userlogin
	auditdata.SessionID = userlogin
	auditdata.ActionTime = time.Now().UTC()
	auditdata.UserIDChanged = data.LoginID
	auditdata.TypeOfChange = "Delete User"

	auditdata.FieldChanged = ""
	auditdata.NewValue = ""
	auditdata.OldValue = ""

	err = c.AclCtx.Save(auditdata)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	err = acl.Delete(data)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	return c.SetResultInfo(false, "", data)
}

func (c *AclController) ExportXLSUserManagement(k *knot.WebContext) interface{} {
	// c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		tk.Printf(err.Error())
	}

	result := []tk.M{}

	d := new(acl.User)
	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.Startwith("groups", "CB_")).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	now := time.Now().UTC()
	times := now.Format("01/02/2006")
	rundate := "Run Date : "
	rundates := rundate + times

	//loginid
	users := k.Session("username").(string)
	user := "User : "
	username := user + users

	font := xlsx.NewFont(11, "Calibri")
	style := xlsx.NewStyle()
	style.Font = *font

	fontHdr := xlsx.NewFont(11, "Calibri")
	fontHdr.Bold = true
	styleHdr := xlsx.NewStyle()
	styleHdr.Font = *fontHdr

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Appilcation : BEF CB"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = ""

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = rundates

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = ""

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = username

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = ""

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "UserId"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "User Name"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "BankId"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Profile"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Country"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Last Loggged In"

	for _, data := range result {

		loginid := data.GetString("loginid")

		results := []tk.M{}
		csr, e := c.AclCtx.Connection.NewQuery().From("acl_log").Where(db.Eq("userid", loginid)).Order("-logintime").Cursor(nil)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
		e = csr.Fetch(&results, 1, false)
		csr.Close()
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}

		// tk.Println("RESULT", results)

		group := data.Get("groups").([]interface{})
		grup := []string{}

		for _, i := range group {
			g := i.(string)
			grup = append(grup, g)
			// tk.Println("grup", grup)
		}
		groups := strings.Join(grup, ",")
		gr := strings.Replace(groups, "CB_", "", -1)

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("id")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("loginid")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("id")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = gr

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = ""

		for _, i := range results {

			logintime, _ := time.Parse(time.RFC3339, i.GetString("logintime"))
			dates := logintime.Format("01/02/2006")

			cell = row.AddCell()
			cell.SetStyle(style)
			cell.Value = dates
		}

	}
	ExcelFilename := "User Metrics Report " + time.Now().Format("20060102150405") + ".xlsx"

	err = file.Save(c.DownloadPath + "/" + ExcelFilename)

	if err != nil {
		tk.Printf(err.Error())
	}

	return c.SetResultInfo(false, "", ExcelFilename)
}
