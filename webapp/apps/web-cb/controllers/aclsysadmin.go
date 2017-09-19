package controllers

import (
	"eaciit/scb-apps/webapp/apps/web-cb/helper"
	// . "eaciit/scb-apps/webapp/apps/web-cb/models"
	"github.com/eaciit/acl/v2.0"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"gopkg.in/gomail.v2"
	// "gopkg.in/mgo.v2/bson"
	// "strconv"
	"strings"
)

type AclSysAdminController struct {
	*BaseController
}

func (a *AclSysAdminController) Default(k *knot.WebContext) interface{} {
	a.LoadBase(k)
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.LayoutTemplate = "_layout_dedicated.html"
	k.Config.IncludeFiles = []string{"shared/sidebar.html", "aclsysadmin/user.html", "aclsysadmin/session.html", "aclsysadmin/access.html", "aclsysadmin/group.html", "aclsysadmin/changePassword.html"}
	return ""
}

func (a *AclSysAdminController) SaveDataUser(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	resdata := ResultInfo{}
	resdata.IsError = true
	resdata.Message = "Save data user"
	tUser := new(acl.User)

	payload := toolkit.M{}

	err := k.GetPayload(&payload)
	if err != nil {
		resdata.Message = toolkit.Sprintf("get payload found : %v", err.Error())
		return resdata
	}

	d := payload.GetString("_id")

	toolkit.Println(payload)

	if d != "" {
		tUser.ID = toolkit.ToString(d)
	} else {
		tUser.ID = toolkit.RandomString(12) + toolkit.ToString(toolkit.RandInt(3))

	}

	enable := payload.Get("enable").(bool)
	group := payload.Get("groups").([]interface{})
	gr := []string{}
	for _, val := range group {
		gr = append(gr, val.(string))
	}

	// toolkit.Println(payload.Get("groups").)
	// group := strings.Split(payload.GetString("groups"), ",")
	// toolkit.Println(gr)
	tUser.LoginID = payload.GetString("loginid")
	tUser.FullName = payload.GetString("fullname")
	tUser.Email = payload.GetString("email")
	tUser.Groups = gr
	tUser.LoginType = acl.LogTypeBasic
	// tUser.LoginType = acl.LogTypeLdap

	tUser.Enable = enable

	if err := acl.Save(tUser); err != nil {
		resdata.Message = toolkit.Sprintf("110. save user found : %v", err.Error())
		return resdata
	}

	if payload.GetString("oldpassword") != payload.GetString("password") && payload.GetString("password") != "" {
		if err := acl.ChangePassword(tUser.ID, payload.GetString("password")); err != nil {
			return err
		}
	}

	resdata.IsError = false
	return helper.CreateResult(true, resdata, "Insert Success")
}

func (a *AclSysAdminController) DeleteDataUser(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson
	resdata := ResultInfo{}
	resdata.IsError = true
	resdata.Message = "Delete data user"

	payload := toolkit.M{}

	err := k.GetPayload(&payload)
	if err != nil {
		resdata.Message = toolkit.Sprintf("get payload found : %v", err.Error())
		return resdata
	}

	//userid
	tUser := new(acl.User)
	if payload.Has("userid") && payload.GetString("userid") != "" {
		tUser.ID = payload.GetString("userid")
		err = acl.Delete(tUser)
		if err != nil {
			resdata.Message = toolkit.Sprintf("delete user found : %v", err.Error())
		} else {
			resdata.IsError = false
		}
	} else {
		resdata.Message = toolkit.Sprintf("User ID not found")
	}

	return resdata
}

func (a *AclSysAdminController) ResetPasswordByAdmin(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := k.GetPayload(&payload); err != nil {
		return toolkit.M{}.Set("message", err.Error()).Set("success", false)
	}

	if !payload.Has("email") || !payload.Has("baseurl") {
		return toolkit.M{}.Set("message", "Data is not complete").Set("success", false)
	}

	uname, tokenid, err := acl.ResetPassword(payload.GetString("email"))
	if err != nil {
		return toolkit.M{}.Set("message", err.Error()).Set("success", false)
	}

	linkstr := toolkit.Sprintf("<a href='%v/acluser/confirmreset?1=%v&2=%v'>Click</a>", payload.GetString("baseurl"), uname, tokenid)

	mailmsg := toolkit.Sprintf("Hi, <br/><br/> We received a request to reset your password, <br/><br/>")
	mailmsg = toolkit.Sprintf("%vFollow the link below to set a new password : <br/><br/> %v <br/><br/>", mailmsg, linkstr)
	mailmsg = toolkit.Sprintf("%vIf you don't want to change your password, you can ignore this email <br/><br/> Thanks,</body></html>", mailmsg)

	m := gomail.NewMessage()

	m.SetHeader("From", "admin.support@eaciit.com")
	m.SetHeader("To", payload.GetString("email"))

	m.SetHeader("Subject", "[no-reply] Self password reset")
	m.SetBody("text/html", mailmsg)

	d := gomail.NewPlainDialer("smtp.office365.com", 587, "admin.support@eaciit.com", "B920Support")
	err = d.DialAndSend(m)

	if err != nil {
		return toolkit.M{}.Set("message", err.Error()).Set("success", false)
	}

	return toolkit.M{}.Set("message", payload.GetString("email")).Set("success", true)
}

func (a *AclSysAdminController) GetDataUser(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson
	resdata := ResultInfo{}
	resdata.IsError = true
	resdata.Message = "Get data user"

	payload := toolkit.M{}

	err := k.GetPayload(&payload)
	if err != nil {
		a.WriteLog(err)
		resdata.Message = toolkit.Sprintf("get payload found : %v", err.Error())
		return resdata
	}

	var filter *dbox.Filter
	if payload.Has("userid") {
		// filter = dbox.And(dbox.Eq("userid", UserID), dbox.Eq("purpose", TokenPurpose))
		filter = dbox.Eq("_id", payload.GetString("userid"))
	}

	//	c, err := acl.Find(tUser, filter, toolkit.M{}.Set("take", take).Set("skip", skip))
	c, err := acl.Find(new(acl.User), filter, nil)
	if err != nil {
		return a.SetResultInfo(true, "109. Cursor found error", err)
	}
	defer c.Close()

	ds := []toolkit.M{}
	err = c.Fetch(&ds, 0, false)
	if err != nil {
		return a.SetResultInfo(true, "115. "+err.Error(), nil)
	}

	resdata.IsError = false
	resdata.Total = len(ds)
	resdata.Data = ds

	return resdata
}

func (a *AclSysAdminController) GetDataAccess(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson
	resdata := ResultInfo{}
	resdata.IsError = true
	resdata.Message = "Get data access"

	// payload := toolkit.M{}

	// err := k.GetPayload(&payload)
	// if err != nil {
	// 	a.WriteLog(err)
	// 	output.Set("message", toolkit.Sprintf("get payload found : %v", err.Error()))
	// 	return output
	// }

	var filter *dbox.Filter
	// filter = dbox.And(dbox.Eq("userid", UserID), dbox.Eq("purpose", TokenPurpose))

	//	c, err := acl.Find(tUser, filter, toolkit.M{}.Set("take", take).Set("skip", skip))
	c, err := acl.Find(new(acl.Access), filter, nil)
	if err != nil {
		return a.SetResultInfo(true, "109. Cursor found error", err)
	}
	defer c.Close()

	ds := []toolkit.M{}
	err = c.Fetch(&ds, 0, false)
	if err != nil {
		return a.SetResultInfo(true, "115. "+err.Error(), nil)
	}

	resdata.IsError = false
	resdata.Total = len(ds)
	resdata.Data = ds

	return resdata
}

func (a *AclSysAdminController) GetDataGroup(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson
	resdata := ResultInfo{}
	resdata.IsError = true
	resdata.Message = "Get data group"

	// payload := toolkit.M{}

	// err := k.GetPayload(&payload)
	// if err != nil {
	// 	a.WriteLog(err)
	// 	output.Set("message", toolkit.Sprintf("get payload found : %v", err.Error()))
	// 	return output
	// }

	var filter *dbox.Filter
	// filter = dbox.And(dbox.Eq("userid", UserID), dbox.Eq("purpose", TokenPurpose))

	//	c, err := acl.Find(tUser, filter, toolkit.M{}.Set("take", take).Set("skip", skip))
	c, err := acl.Find(new(acl.Group), filter, nil)
	if err != nil {
		return a.SetResultInfo(true, "109. Cursor found error", err)
	}
	defer c.Close()

	ds := []toolkit.M{}
	err = c.Fetch(&ds, 0, false)
	if err != nil {
		return a.SetResultInfo(true, "115. "+err.Error(), nil)
	}

	resdata.IsError = false
	resdata.Total = len(ds)
	resdata.Data = ds

	return resdata
}

func (a *AclSysAdminController) GetDataSession(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson
	resdata := ResultInfo{}
	resdata.IsError = true
	resdata.Message = "Get data session login"

	c, err := a.Ctx.Connection.NewQuery().From(new(acl.Session).TableName()).Order("-created").Cursor(nil)
	if err != nil {
		return a.SetResultInfo(true, "109. Cursor found error", err)
	}
	defer c.Close()

	ds := []toolkit.M{}
	err = c.Fetch(&ds, 0, false)
	if err != nil {
		return a.SetResultInfo(true, "115. "+err.Error(), nil)
	}

	resdata.IsError = false
	resdata.Total = len(ds)
	resdata.Data = ds

	return resdata
}

func (a *AclSysAdminController) SaveMenu(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson
	resdata := ResultInfo{}
	resdata.IsError = true
	resdata.Message = "Save data user"

	payload := toolkit.M{}

	err := k.GetPayload(&payload)
	if err != nil {
		resdata.Message = toolkit.Sprintf("get payload found : %v", err.Error())
		return resdata
	}
	toolkit.Println("===========enable=======")
	toolkit.Println(k.Request.FormValue("Enable"))

	tMenu := new(acl.Access)
	tMenu.ID = payload.GetString("ID")
	tMenu.Title = payload.GetString("Title")
	tMenu.Category = 2
	tMenu.ParentId = payload.GetString("ParentId")
	tMenu.Url = payload.GetString("Url")
	tMenu.Index = payload.GetInt("Index")
	tMenu.Group1 = payload.GetString("Group1")
	// tMenu.Enable, _ = strconv.ParseBool()
	tMenu.Enable = payload.Get("Enable").(bool)

	if err := acl.Save(tMenu); err != nil {
		resdata.Message = toolkit.Sprintf("110. save user found : %v", err.Error())
		return resdata
	}

	resdata.IsError = false
	return resdata
}

func (a *AclSysAdminController) DeleteDataMenu(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson
	resdata := ResultInfo{}
	resdata.IsError = true
	resdata.Message = "Delete data Menu"

	payload := toolkit.M{}

	err := k.GetPayload(&payload)
	if err != nil {
		resdata.Message = toolkit.Sprintf("get payload found : %v", err.Error())
		return resdata
	}

	tMenu := new(acl.Access)
	if payload.Has("Id") && payload.GetString("Id") != "" {
		tMenu.ID = payload.GetString("Id")
		err = acl.Delete(tMenu)
		if err != nil {
			resdata.Message = toolkit.Sprintf("delete menu found : %v", err.Error())
		} else {
			resdata.IsError = false
		}
	} else {
		resdata.Message = toolkit.Sprintf("Menu ID not found")
	}

	return resdata
}

func (a *AclSysAdminController) SaveGroupUser(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	payload := toolkit.M{}

	err := k.GetPayload(&payload)
	if err != nil {
		return err.Error()
	}

	group := new(acl.Group)
	grants := strings.Split(payload.GetString("listRole"), ",")
	for _, v := range grants {
		access := acl.AccessGrant{}
		access.AccessID = v
		access.AccessValue = 15
		group.Grants = append(group.Grants, access)
	}

	if strings.Compare(payload.GetString("valueEnable"), "true") == 0 {
		group.Enable = true
	} else {
		group.Enable = false
	}

	group.ID = payload.GetString("Id")
	group.Title = payload.GetString("Name")
	group.GroupType = 0

	e := acl.Save(group)
	if e != nil {
		return helper.CreateResult(false, nil, "gagal")
	}

	return helper.CreateResult(true, payload, "success")
}

func (a *AclSysAdminController) GetGroupById(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	payload := toolkit.M{}

	err := k.GetPayload(&payload)
	if err != nil {
		return err.Error()
	}

	group := new(acl.Group)
	e := acl.FindByID(group, payload.GetString("id"))

	if e != nil {
		return helper.CreateResult(false, nil, "Group User Not Found")
	}

	return group
}

func (a *AclSysAdminController) DeleteGroup(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	payload := toolkit.M{}

	err := k.GetPayload(&payload)
	if err != nil {
		return err.Error()
	}

	idgroup := payload.GetString("id")
	var filter *dbox.Filter
	filter = dbox.Eq("groups", idgroup)
	c, err := acl.Find(new(acl.User), filter, nil)
	if err != nil {
		return a.SetResultInfo(true, "109. Cursor found error", err)
	}
	defer c.Close()
	tes := []toolkit.M{}
	err = c.Fetch(&tes, 0, false)

	if len(tes) > 0 {
		return helper.CreateResult(false, nil, "Delete Group Failed")
	}

	group := new(acl.Group)
	err = acl.FindByID(group, payload.GetString("id"))

	if err != nil {
		return helper.CreateResult(false, nil, "Group User Not Found")
	}

	err = acl.Delete(group)
	if err != nil {
		return helper.CreateResult(false, nil, "Delete Group Failed")
	}

	return helper.CreateResult(true, nil, "Delete Group Success")
}
