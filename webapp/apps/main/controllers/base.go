package controllers

import (
	"eaciit/scb-apps/webapp/apps/main/models"
	"eaciit/scb-apps/webapp/helper"
	"errors"
	"github.com/eaciit/acl/v2.0"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"net/http"
	"strings"
)

type IBaseController interface{}

type BaseController struct {
	base                     IBaseController
	Ctx                      *orm.DataContext
	Conn                     dbox.IConnection
	AppName                  string
	NoLogin                  bool
	IsUsingEncryptedPassword bool
	IsUsingLDAP              bool
	ServerNameLDAP           string
	AddressLDAP              string
	BaseDNLDAP               string
	LDAPType                 string
	LDAPCertificate          []string
	UserDNLDAP               string
	UserAuthAttrLDAP         string
	BindUsernameLDAP         string
	BindPasswordLDAP         string
	InsecureSkipVerify       bool
	BindFilterLDAP           string

	// LDAP Attribute.GET
	LDAP_DATA acl.LDAP_DATA_LIST
}

const (
	SESSION_KEY         string = "sessionid"
	SESSION_USERNAME    string = "username"
	SESSION_FULLNAME    string = "fullname"
	SESSION_FIRSTNAME   string = "firstname"
	SESSION_LASTNAME    string = "lastname"
	SESSION_EXPIREDTIME string = "expiredtime"
	SESSION_LOGINTIME   string = "logintime"
	SESSION_GROUP       string = "group"

	//
	SESSION_COUNTRYCODE string = "countrycode"
	SESSION_COUNTRY     string = "country"
)

var (
	LoginDataUser       acl.User
	LoginDataGroups     []acl.Group
	LoginDataAccessMenu []acl.Access
	LoginDataSession    acl.Session
)

func (b *BaseController) Authenticate(k *knot.WebContext, callback, failback func()) {
	sessionid := tk.ToString(k.Session(SESSION_KEY, ""))
	if acl.IsSessionIDActive(sessionid) {
		if callback != nil {
			callback()
		}
	} else {
		k.SetSession(SESSION_KEY, "")
		if failback != nil {
			failback()
		}
	}
}

func (b *BaseController) IsLoggedIn(k *knot.WebContext) bool {
	return k.Session(SESSION_KEY, "") != ""
}

func (b *BaseController) GetCurrentUsername(k *knot.WebContext) string {
	if !b.IsLoggedIn(k) {
		return ""
	}

	return k.Session(SESSION_USERNAME).(string)
}

func (b *BaseController) PrepareCurrentUserData(k *knot.WebContext) {

	// ==== user

	username := b.GetCurrentUsername(k)
	user := new(acl.User)
	err := acl.FindUserByLoginID(user, username)
	if err != nil {
		return
	}
	k.SetSession(SESSION_FULLNAME, user.FullName)
	k.SetSession(SESSION_FIRSTNAME, "")
	k.SetSession(SESSION_LASTNAME, "")
	k.SetSession(SESSION_COUNTRYCODE, user.CountryCode)
	k.SetSession(SESSION_COUNTRY, user.Country)

	// ==== groups

	groups := make([]acl.Group, 0)
	groupAccessMenu := tk.M{}
	for _, each := range user.Groups {
		group := new(acl.Group)
		err = acl.FindByID(group, each)
		if err != nil {
			return
		}
		groups = append(groups, *group)

		for _, each := range group.Grants {
			if !groupAccessMenu.Has(each.AccessID) {
				groupAccessMenu.Set(each.AccessID, 0)
			}
		}
	}
	k.SetSession(SESSION_GROUP, user.Groups)

	// ==== access menu

	csr, err := acl.Find(new(acl.Access), nil, nil)
	defer csr.Close()
	if err != nil {
		return
	}

	accessMenuAll := make([]acl.Access, 0)
	err = csr.Fetch(&accessMenuAll, 0, false)
	if err != nil {
		return
	}

	// ==== get access for current user

	allowed := make([]acl.Access, 0)
	for _, each := range accessMenuAll {
		if groupAccessMenu.Has(each.ID) {
			allowed = append(allowed, each)
		}
	}

	// Get Login and Expired Time
	usersession := acl.Session{}
	sessionID := k.Session(SESSION_KEY, "").(string)
	csr, err = b.Ctx.Connection.NewQuery().From(new(acl.Session).TableName()).Where(dbox.Eq("_id", sessionID)).Cursor(nil)
	err = csr.Fetch(&usersession, 1, false)
	csr.Close()
	if err == nil {
		k.SetSession(SESSION_LOGINTIME, usersession.Created)
		k.SetSession(SESSION_EXPIREDTIME, usersession.Expired)
	}

	// Get Login and

	// ==== save it to memory

	LoginDataUser = *user
	LoginDataGroups = groups
	LoginDataAccessMenu = allowed
	LoginDataSession = usersession
}

func (b *BaseController) SetResponseTypeHTML(k *knot.WebContext) {
	k.Config.OutputType = knot.OutputTemplate
}

func (b *BaseController) SetResponseTypeAJAX(k *knot.WebContext) {
	k.Config.OutputType = knot.OutputJson
}

func (b *BaseController) ValidateAccessOfRequestedURL(k *knot.WebContext, grants ...string) bool {
	config := GetConfig()
	unauthorizedErrorMessage := GetUnauthorizedMessageAsQueryString(k)
	landingPagePath := strings.Trim(config.GetString("landingpage"), ` /`)

	isOK := false
	b.Authenticate(k, func() {
		if len(grants) > 0 {
			isAllowed := false
			for _, grant := range grants {
				if tk.HasMember(LoginDataUser.Groups, grant) && !isAllowed {
					isAllowed = true
				}
			}

			if !isAllowed {
				b.Redirect(k, strings.Split(landingPagePath, "/")[1], strings.Split(landingPagePath, "/")[2]+unauthorizedErrorMessage)
			}
		}
	}, func() {
		// if strings.Contains(strings.ToLower(k.Request.URL.String()),"logout") >= 0 {
		// b.Redirect(k, "auth", "login"+unauthorizedErrorMessage)
		b.Redirect(k, "auth", "login")
	})

	return isOK
}

func (b *BaseController) Redirect(k *knot.WebContext, controller string, action string) {
	urlString := "/" + b.AppName + "/" + controller + "/" + action
	http.Redirect(k.Writer, k.Request, urlString, http.StatusTemporaryRedirect)
}

func (b *BaseController) SetResultOK(data interface{}) *tk.Result {
	r := tk.NewResult()
	r.Data = data

	return r
}

func (b *BaseController) SetResultError(msg string, data interface{}) *tk.Result {
	tk.Println(msg)

	r := tk.NewResult()
	r.SetError(errors.New(msg))
	r.Data = data

	return r
}

func (b *BaseController) SetViewData(viewData tk.M) tk.M {
	if viewData == nil {
		viewData = tk.M{}
	}

	viewData.Set("AppName", b.AppName)
	return viewData
}

func (b *BaseController) GetViewBaseData(k *knot.WebContext) tk.M {
	data := tk.M{}

	if b.IsLoggedIn(k) {
		data.Set("UserData", LoginDataUser)
	}

	return data
}

func (b *BaseController) GetApplicationByUserName(username string) ([]models.ApplicationModel, error) {
	defaultValue := make([]models.ApplicationModel, 0)
	// ===== get users

	user := new(acl.User)
	err := acl.FindUserByLoginID(user, username)
	if err != nil {
		return defaultValue, err
	}

	// ===== get groups

	csrGroup, err := acl.Find(new(models.GroupModel), nil, nil)
	defer csrGroup.Close()
	if err != nil {
		return defaultValue, err
	}

	dataGroup := make([]models.GroupModel, 0)
	err = csrGroup.Fetch(&dataGroup, 0, false)
	if err != nil {
		return defaultValue, err
	}

	// ===== get applications

	csr, err := b.Ctx.Find(new(models.ApplicationModel), nil)
	defer csr.Close()
	if err != nil {
		return defaultValue, err
	}

	dataApplication := make([]models.ApplicationModel, 0)
	err = csr.Fetch(&dataApplication, 0, false)
	if err != nil {
		return defaultValue, err
	}

	// ====== get current user app

	if tk.HasMember(LoginDataUser.Groups, "admin") {
		return dataApplication, nil
	}

	result := make([]models.ApplicationModel, 0)

	for _, group := range dataGroup {
		if tk.HasMember(user.Groups, group.ID) {
			for _, app := range dataApplication {
				if tk.HasMember(group.Applications, app.ID) {
					result = append(result, app)
				}
			}
		}
	}

	return result, nil
}

func (b *BaseController) GetAccessMenuByApplicationID(appID string) ([]models.AccessMenuModel, error) {
	defaultValue := make([]models.AccessMenuModel, 0)

	csr, err := b.Ctx.Connection.
		NewQuery().
		From(new(models.AccessMenuModel).TableName()).
		Where(dbox.Eq("applicationid", appID)).
		Cursor(nil)
	defer csr.Close()
	if err != nil {
		return defaultValue, err
	}

	result := make([]models.AccessMenuModel, 0)
	err = csr.Fetch(&result, 0, false)
	if err != nil {
		return defaultValue, err
	}

	return result, nil
}

func GetConfig() tk.M {
	type ForgetMe struct{}
	config := helper.ReadConfig(ForgetMe{})
	return config
}

func GetUnauthorizedMessageAsQueryString(k *knot.WebContext) string {
	return "?NotAllowed=You don't have permission to access requested page " + strings.Split(k.Request.URL.String(), "?")[0]
}
