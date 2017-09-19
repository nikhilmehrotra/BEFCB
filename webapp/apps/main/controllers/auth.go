package controllers

import (
	// "eaciit/scb-apps/webapp/apps/main/models"
	"github.com/eaciit/acl/v2.0"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// bson "gopkg.in/mgo.v2/bson"
	"log"
	"strings"
)

type AuthController struct {
	*BaseController
}

func (c *AuthController) Login(k *knot.WebContext) interface{} {
	c.SetResponseTypeHTML(k)
	k.Config.LayoutTemplate = ""

	return c.SetViewData(nil)
}

func (c *AuthController) Logout(k *knot.WebContext) interface{} {
	c.SetResponseTypeHTML(k)

	return c.SetViewData(nil)
}

func (c *AuthController) DoLogin(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := struct {
		Username string
		Password string
		Key      int
	}{}

	err := k.GetPayload(&payload)
	if err != nil {
		log.Println("# Login FAILED : ", err.Error())
		return c.SetResultError(err.Error(), nil)
	}
	IsAdminLogin := false
	if payload.Key == 1 {
		IsAdmin := strings.Contains(c.BaseDNLDAP, payload.Username)
		if IsAdmin == false {
			log.Println("# Login FAILED : ", err.Error())
			return c.SetResultError(err.Error(), nil)
		}
		IsAdminLogin = true

	}

	sessionId, _, err := acl.Login(payload.Username, payload.Password, c.IsUsingLDAP, c.LDAPType, c.AddressLDAP, c.BaseDNLDAP, c.ServerNameLDAP, c.LDAPCertificate, c.UserAuthAttrLDAP, c.UserDNLDAP, c.BindUsernameLDAP, c.BindPasswordLDAP, c.InsecureSkipVerify, c.BindFilterLDAP, c.LDAP_DATA, IsAdminLogin)
	if err != nil {
		log.Println("# Login FAILED : ", err.Error())
		return c.SetResultError(err.Error(), nil)
	}

	activeUser := new(acl.User)
	err = acl.FindUserByLoginID(activeUser, payload.Username)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}
	result := new(acl.Session)
	groupResult := new(acl.Group)

	csr, err := c.Ctx.Connection.NewQuery().From(result.TableName()).Where(db.Eq("id", sessionId)).Cursor(nil)
	err = csr.Fetch(&result, 1, false)
	csr.Close()
	if err != nil {
		return nil
	}
	groupTitle := []string{}
	for _, x := range activeUser.Groups {
		csr, err := c.Ctx.Connection.NewQuery().From(groupResult.TableName()).Where(db.Eq("id", x)).Cursor(nil)
		err = csr.Fetch(&groupResult, 1, false)
		csr.Close()
		if err != nil {
			return nil
		}
		resultTitle := groupResult.Title
		groupTitle = append(groupTitle, resultTitle)
	}

	title := strings.Join(groupTitle, ",")
	groups := strings.Join(activeUser.Groups, ",")

	if k.Session(SESSION_FIRSTNAME) != nil {
		result.FirstName = k.Session(SESSION_FIRSTNAME).(string)
	}

	if k.Session(SESSION_LASTNAME) != nil {
		result.LastName = k.Session(SESSION_LASTNAME).(string)
	}

	result.Country = activeUser.Country
	if activeUser.Country == "" {
		result.Country = "GLOBAL"
	}
	result.CountryCode = activeUser.CountryCode
	if activeUser.CountryCode == "" {
		result.CountryCode = "GLOBAL"
	}
	result.FullName = activeUser.FullName
	result.Group = groups
	result.GroupDescription = title

	e := acl.Save(result)
	if e != nil {
		return nil
	}

	k.SetSession(SESSION_KEY, sessionId)
	k.SetSession(SESSION_USERNAME, payload.Username)

	c.PrepareCurrentUserData(k)

	redirect := GetConfig().GetString("landingpage")
	apps, err := c.GetApplicationByUserName(payload.Username)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}
	if len(apps) == 1 {
		redirect = `/` + apps[0].ID + `/` + strings.Trim(apps[0].LandingURL, ` /`)
	}

	return c.SetResultOK(tk.M{}.
		Set(SESSION_KEY, sessionId).
		Set(SESSION_USERNAME, payload.Username).
		Set("redirect", redirect))
}

func (c *AuthController) DoLogout(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	sessionId := tk.ToString(k.Session(SESSION_KEY, ""))
	acl.Logout(sessionId)

	k.SetSession(SESSION_KEY, "")
	k.SetSession(SESSION_USERNAME, "")
	c.Redirect(k, "auth", "login")

	return nil
}

func (c *AuthController) GetUserInfo(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	return c.SetResultOK(LoginDataUser)
}
