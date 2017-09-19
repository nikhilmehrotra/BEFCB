package controllers

import (
	// "eaciit/scbocir/helper"
	m "eaciit/scb-apps/webapp/apps/web-cb/models"
	"errors"
	"strings"

	"github.com/eaciit/acl/v2.0"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"gopkg.in/gomail.v2"
)

type AclUserController struct {
	*BaseController
}

func (a *AclUserController) Default(k *knot.WebContext) interface{} {
	usersess := k.Session("sessionid", "")
	sessionid := ""
	if usersess != nil {
		sessionid = toolkit.ToString(usersess)
	}
	if usersess != nil && acl.IsSessionIDActive(sessionid) {
		a.Redirect(k, "dashboard", "default")
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.LayoutTemplate = ""
	return ""
}

func (a *AclUserController) GetUserList(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson

	result := make([]toolkit.M, 0)
	csr, e := a.AclCtx.Connection.NewQuery().From("acl_users").Select("loginid", "fullname", "groups").Cursor(nil)
	if e != nil {
		return a.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return a.ErrorResultInfo(e.Error(), nil)
	}

	return a.SetResultInfo(false, "", result)
}

func (a *AclUserController) ConfirmReset(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.LayoutTemplate = ""
	return ""
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (a *AclUserController) Login(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	output := toolkit.M{}.Set("valid", false).Set("message", "not yet login")

	err := k.GetPayload(&payload)
	if err != nil {
		toolkit.Println("# Login FAILED : ", err.Error())
		output.Set("message", err.Error())
		return output
	}

	sessionid, err := Login(payload.GetString("username"), payload.GetString("password"), a.IsUsingLDAP, a.LDAPType, a.AddressLDAP, a.BaseDNLDAP, a.ServerNameLDAP, a.LDAPCertificate, a.UserAuthAttrLDAP, a.UserDNLDAP, a.BindUsernameLDAP, a.BindPasswordLDAP, a.InsecureSkipVerify, a.BindFilterLDAP)
	if err != nil {
		toolkit.Println("# Login FAILED : ", err.Error())
		output.Set("message", err.Error())
		return output
	}

	output.Set("valid", true).Set("message", "login success").Set("sessionid", sessionid)
	k.SetSession("sessionid", sessionid)
	k.SetSession("username", payload.GetString("username"))

	results := acl.Session{}
	csr, errs := a.AclCtx.Connection.NewQuery().From("acl_sessions").Where(db.Eq("_id", sessionid)).Cursor(nil)
	errs = csr.Fetch(&results, 1, false)
	csr.Close()
	if errs != nil {
		toolkit.Println("# Login FAILED : ", err.Error())
		return nil
	}

	k.SetSession("logintime", results.Created)
	k.SetSession("expiredtime", results.Expired)

	result := toolkit.M{}
	csr, e := a.AclCtx.Connection.NewQuery().From("acl_users").Select("loginid", "fullname", "groups", "country", "countrycode").Where(db.Eq("loginid", payload.GetString("username"))).Cursor(nil)
	if e != nil {
		toolkit.Println("# Login FAILED : ", err.Error())
		return a.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&result, 1, false)
	csr.Close()
	if e != nil {
		toolkit.Println("# Login FAILED : ", err.Error())
		return a.ErrorResultInfo(e.Error(), nil)
	}

	group := result.Get("groups").([]interface{})
	grup := []string{}

	// username := result.Get("loginid").(string)
	fullname := ""
	if result.Get("fullname") != nil {
		fullname = result.Get("fullname").(string)
	}
	countrycode, country := "", ""
	if result.Get("country") != nil {
		country = result.Get("country").(string)
		countrycode = result.Get("countrycode").(string)
	}
	k.SetSession("fullname", fullname)
	k.SetSession("country", country)
	k.SetSession("countrycode", countrycode)
	isGlobal := false
	for _, i := range group {
		g := i.(string)
		if strings.Contains(g, "GLOBAL_") {
			isGlobal = true
		}
		grup = append(grup, g)
	}

	// toolkit.Println("GRUP", grup)

	//get last menu
	landpage := []m.AccessibilityModel{}

	csr, e = a.AclCtx.Connection.NewQuery().From("acl_accessibility").Where(db.And(db.In("roleid", group...), db.Eq("allowstatus", true))).Cursor(nil)
	if e != nil {
		return a.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&landpage, 1, false)
	csr.Close()
	if e != nil {
		return a.ErrorResultInfo(e.Error(), nil)
	}
	// toolkit.Println("LASTLINK", lastlink)

	LandingPage := ""

	for _, i := range landpage {
		LandingPage = i.Url

	}

	// toolkit.Println("LINK", LandingPage)

	k.SetSession("group", grup)
	if isGlobal {
		k.SetSession("sessionid", "")
		err := errors.New("Invalid Login ID / Password")
		output.Set("message", err.Error())
		return output
		// if stringInSlice("BUSINESSMANAGER", grup) || stringInSlice("BEFSPONSORS", grup) || stringInSlice("INITIATIVEOWNERS", grup) || stringInSlice("ACCOUNTABLEEXECUTIVES", grup) || stringInSlice("PROJECTMANAGERS", grup) || stringInSlice("TECHNOLOGYLEADS", grup) || stringInSlice("VIEWONLY", grup) || username == "global" {
		// k.SetSession("prefix", "global")
		// output.Set("prefix", "global")
	} else {
		// k.SetSession("sessionid", "")
		// err = errors.New("Username and password is incorrect")
		k.SetSession("prefix", "bef")
		output.Set("prefix", "bef")
	}
	return output.Set("LandingPage", LandingPage)
}

func (a *AclUserController) Logout(k *knot.WebContext) interface{} {
	username := ""
	if k.Session("username") != nil {
		username = k.Session("username").(string)
		toolkit.Println("# ", username, " - Logging Out")
	}
	sessionid := toolkit.ToString(k.Session("sessionid", ""))
	if sessionid != "" {
		err := acl.Logout(sessionid)
		if err != nil {
			a.WriteLog(err)
			return toolkit.M{}.Set("message", toolkit.Sprintf("logout process found : %v", err.Error()))
		}
	}

	k.SetSession("sessionid", "")
	a.Redirect(k, "login", "default")
	return ""
}

func (a *AclUserController) SaveNewPassword(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := k.GetPayload(&payload); err != nil {
		return toolkit.M{}.Set("message", err.Error())
	}

	sessionid := toolkit.ToString(k.Session("sessionid", ""))
	if !payload.Has("newpassword") || (!payload.Has("tokenid") && sessionid == "") {
		return errors.New("Data is not complete")
	}

	var err error
	switch {
	case payload.Has("tokenid"):
		err = acl.ChangePasswordToken(payload.GetString("userid"), payload.GetString("newpassword"), payload.GetString("tokenid"))
	default:
		var userid string
		userid, err = acl.FindUserBySessionID(sessionid)
		if err != nil {
			return err.Error()
		}

		if sessionid == payload.GetString("sessionid") {
			err = acl.ChangePasswordFromOld(userid, payload.GetString("newpassword"), payload.GetString("oldpassword"))
		} else if err == nil {
			err = errors.New("session is not match")
		}
	}

	if err != nil {
		return err.Error()
	}

	return nil
}

func (a *AclUserController) ResetPassword(k *knot.WebContext) interface{} {
	// k.Config.NoLog = true
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

func (a *AclUserController) GetListTopMenu(k *knot.WebContext) interface{} {
	a.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	sessionid := toolkit.ToString(k.Session("sessionid", ""))
	arrmenu, err := acl.GetListMenuBySessionId(sessionid)

	if err != nil {
		return err.Error()
	}

	return arrmenu
}
