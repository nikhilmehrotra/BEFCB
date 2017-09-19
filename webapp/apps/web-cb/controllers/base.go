package controllers

import (
	_ "eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	// "errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/eaciit/acl/v2.0"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

type IBaseController interface {
}
type BaseController struct {
	APP_NAME           string
	base               IBaseController
	Ctx                *orm.DataContext
	AclCtx             *orm.DataContext
	UploadPath         string
	DownloadPath       string
	TemplatePath       string
	UploadMetricPath   string
	IsUsingLDAP        bool
	ServerNameLDAP     string
	AddressLDAP        string
	BaseDNLDAP         string
	LDAPType           string
	LDAPCertificate    []string
	UserDNLDAP         string
	UserAuthAttrLDAP   string
	BindUsernameLDAP   string
	BindPasswordLDAP   string
	InsecureSkipVerify bool
	BindFilterLDAP     string
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

type PageInfo struct {
	PageTitle    string
	SelectedMenu string
	Breadcrumbs  map[string]string
}

type ResultInfo struct {
	IsError bool
	Message string
	Total   int
	Data    interface{}
}

type Previlege struct {
	View     bool
	Create   bool
	Edit     bool
	Delete   bool
	Approve  bool
	Process  bool
	Menuid   string
	Menuname string
	Username string
}

// func (b *BaseController) LoadBaseUrl(k *knot.WebContext) {

// 	k.Config.NoLog = true
// 	b.IsAuthenticate(k)

// 	return
// }

func (b *BaseController) GetAccess(k *knot.WebContext, AccessID string) tk.M {
	k.Config.NoLog = true

	var groups []string

	if k.Session("group") != nil {
		groups = k.Session("group").([]string)
	}

	grup := []interface{}{}

	for _, u := range groups {
		group := u
		grup = append(grup, group)
	}

	d := new(AccessibilityModel)
	result := []AccessibilityModel{}
	Url := strings.Trim(k.Request.RequestURI, " ")
	query := []*db.Filter{}
	query = append(query, db.In("roleid", grup...))
	prefix := ""
	if k.Session("prefix") != nil {
		prefix = k.Session("prefix").(string)
	}
	Url = strings.Replace(Url, "/"+prefix, "", -1)
	if Url != "" && Url != "#" && AccessID == "" {
		query = append(query, db.Eq("url", Url))
	} else {
		query = append(query, db.Eq("accessid", AccessID))
	}
	csr, err := b.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.And(query...)).Cursor(nil)
	err = csr.Fetch(&result, 0, false)
	csr.Close()
	if err != nil {
		return nil
	}
	Status := false

	for _, i := range result {
		allowstatus := i.AllowStatus
		// Global
		if i.Global.Create {
			d.Global.Create = true
		}
		if i.Global.Read {
			d.Global.Read = true
		}
		if i.Global.Update {
			d.Global.Update = true
		}
		if i.Global.Delete {
			d.Global.Delete = true
		}
		if i.Global.Owned {
			d.Global.Owned = true
		}
		if i.Global.Curtain {
			d.Global.Curtain = true
		}
		if i.Global.Upload {
			d.Global.Upload = true
		}
		// Region
		if i.Region.Create {
			d.Region.Create = true
		}
		if i.Region.Read {
			d.Region.Read = true
		}
		if i.Region.Update {
			d.Region.Update = true
		}
		if i.Region.Delete {
			d.Region.Delete = true
		}
		if i.Region.Owned {
			d.Region.Owned = true
		}
		if i.Region.Curtain {
			d.Region.Curtain = true
		}
		if i.Region.Upload {
			d.Region.Upload = true
		}
		// Country
		if i.Country.Create {
			d.Country.Create = true
		}
		if i.Country.Read {
			d.Country.Read = true
		}
		if i.Country.Update {
			d.Country.Update = true
		}
		if i.Country.Delete {
			d.Country.Delete = true
		}
		if i.Country.Owned {
			d.Country.Owned = true
		}
		if i.Country.Curtain {
			d.Country.Curtain = true
		}
		if i.Country.Upload {
			d.Country.Upload = true
		}

		if allowstatus == true {
			Status = true
		}
		if allowstatus == false {
			Status = false
		}
	}

	if Status == true {
		d.AllowStatus = true
	} else {
		// k.SetSession("sessionid", "")
		// b.Redirect(k, "acluser", "default")
	}
	// tk.Println("AllowStatus : ", d.AllowStatus)
	returnvalue := tk.M{}
	err = tk.StructToM(d, &returnvalue)
	// tk.Println("Status : ", returnvalue)
	return returnvalue

}

func (b *BaseController) GetAccessStatus(k *knot.WebContext, AccessID string) *AccessibilityModel {
	k.Config.NoLog = true

	var groups []string

	if k.Session("group") != nil {
		groups = k.Session("group").([]string)
	}

	grup := []interface{}{}

	for _, u := range groups {
		group := u
		grup = append(grup, group)
	}

	d := new(AccessibilityModel)
	result := []AccessibilityModel{}
	Url := strings.Trim(k.Request.RequestURI, " ")
	query := []*db.Filter{}
	query = append(query, db.In("roleid", grup...))
	prefix := ""
	if k.Session("prefix") != nil {
		prefix = k.Session("prefix").(string)
	}
	Url = strings.Replace(Url, "/"+prefix, "", -1)
	if Url != "" && Url != "#" && AccessID == "" {
		query = append(query, db.Eq("url", Url))
	} else {
		query = append(query, db.Eq("accessid", AccessID))
	}
	csr, err := b.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.And(query...)).Cursor(nil)
	err = csr.Fetch(&result, 0, false)
	csr.Close()
	if err != nil {
		return nil
	}
	Status := false

	for _, i := range result {
		allowstatus := i.AllowStatus
		// Global
		if i.Global.Create {
			d.Global.Create = true
		}
		if i.Global.Read {
			d.Global.Read = true
		}
		if i.Global.Update {
			d.Global.Update = true
		}
		if i.Global.Delete {
			d.Global.Delete = true
		}
		if i.Global.Owned {
			d.Global.Owned = true
		}
		if i.Global.Curtain {
			d.Global.Curtain = true
		}
		if i.Global.Upload {
			d.Global.Upload = true
		}
		// Region
		if i.Region.Create {
			d.Region.Create = true
		}
		if i.Region.Read {
			d.Region.Read = true
		}
		if i.Region.Update {
			d.Region.Update = true
		}
		if i.Region.Delete {
			d.Region.Delete = true
		}
		if i.Region.Owned {
			d.Region.Owned = true
		}
		if i.Region.Curtain {
			d.Region.Curtain = true
		}
		if i.Region.Upload {
			d.Region.Upload = true
		}
		// Country
		if i.Country.Create {
			d.Country.Create = true
		}
		if i.Country.Read {
			d.Country.Read = true
		}
		if i.Country.Update {
			d.Country.Update = true
		}
		if i.Country.Delete {
			d.Country.Delete = true
		}
		if i.Country.Owned {
			d.Country.Owned = true
		}
		if i.Country.Curtain {
			d.Country.Curtain = true
		}
		if i.Country.Upload {
			d.Country.Upload = true
		}

		if allowstatus == true {
			Status = true
		}
		if allowstatus == false {
			Status = false
		}
	}

	if Status == true {
		d.AllowStatus = true
	} else {
		// k.SetSession("sessionid", "")
		// b.Redirect(k, "acluser", "default")
	}
	return d

}

func (b *BaseController) LoadBaseAjaxServ(k *knot.WebContext) {
	k.Config.NoLog = true

	sessionid := tk.ToString(k.Session("sessionid", ""))
	if !acl.IsSessionIDActive(sessionid) {
		// return errors.New("Session Expired")
	}

	return
}

func (b *BaseController) LoadBase(k *knot.WebContext) []tk.M {
	k.Config.NoLog = true
	b.IsAuthenticate(k)
	// access := b.AccessMenu(k)
	return []tk.M{}
}

func (b *BaseController) IsAuthenticate(k *knot.WebContext) {
	sessionid := tk.ToString(k.Session("sessionid", ""))
	if !acl.IsSessionIDActive(sessionid) {
		k.SetSession("sessionid", "")
		b.Redirect(k, "acluser", "default")
	}
	return
}

func (b *BaseController) AccessMenu(k *knot.WebContext) []tk.M {
	url := k.Request.URL.String()
	if strings.Index(url, "?") > -1 {
		url = url[:strings.Index(url, "?")]
		//		tk.Println("URL_PARSED,", url)
	}
	sessionRoles := k.Session("roles")
	access := []tk.M{}
	if sessionRoles != nil {
		accesMenu := sessionRoles.([]SysRolesModel)
		if len(accesMenu) > 0 {
			for _, o := range accesMenu[0].Menu {
				if o.Url == url {
					obj := tk.M{}
					obj.Set("View", o.View)
					obj.Set("Create", o.Create)
					obj.Set("Approve", o.Approve)
					obj.Set("Delete", o.Delete)
					obj.Set("Process", o.Process)
					obj.Set("Edit", o.Edit)
					obj.Set("Menuid", o.Menuid)
					obj.Set("Menuname", o.Menuname)
					obj.Set("Username", k.Session("username").(string))
					access = append(access, obj)
					return access
				}

			}
		}
	}
	return access
}

func (b *BaseController) IsLoggedIn(k *knot.WebContext) bool {
	if k.Session("userid") == nil {
		return false
	}
	return true
}
func (b *BaseController) GetCurrentUser(k *knot.WebContext) string {
	if k.Session("userid") == nil {
		return ""
	}
	return k.Session("username").(string)
}
func (b *BaseController) Redirect(k *knot.WebContext, controller string, action string) {
	http.Redirect(k.Writer, k.Request, "/"+controller+"/"+action, http.StatusTemporaryRedirect)
}

func (b *BaseController) WriteLog(msg interface{}) {
	log.Printf("%#v\n\r", msg)
	return
}
func (b *BaseController) SetResultInfo(isError bool, msg string, data interface{}) ResultInfo {
	r := ResultInfo{}
	r.IsError = isError
	r.Message = msg
	r.Data = data
	return r
}

func (b *BaseController) ErrorResultInfo(msg string, data interface{}) ResultInfo {
	r := ResultInfo{}
	r.IsError = true
	r.Message = msg
	r.Data = data
	return r
}
func (b *BaseController) Round(f float64) float64 {
	return math.Floor(f + .5)
}
func (b *BaseController) RoundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return b.Round(f*shift) / shift
}
func (b *BaseController) FirstMonday(year int, mn int) int {
	month := time.Month(mn)
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	d0 := (8-int(t.Weekday()))%7 + 1
	s := strconv.Itoa(year) + fmt.Sprintf("%02d", mn) + fmt.Sprintf("%02d", d0)
	ret, _ := strconv.Atoi(s)
	return ret
}

func (b *BaseController) FirstWorkDay(ym string) int {
	t, err := time.Parse("2006-01-02", ym+"-01")
	if err != nil {
		fmt.Println(err.Error())
	}
	for t.Weekday() == 0 || t.Weekday() == 6 {
		if t.Weekday() == 0 {
			t = t.AddDate(0, 0, 1)
		} else if t.Weekday() == 6 {
			t = t.AddDate(0, 0, 2)
		}
	}
	ret, _ := strconv.Atoi(t.Format("20060102"))
	return ret
}

func (b *BaseController) GetNextIdSeq(collName string) (int, error) {
	ret := 0
	mdl := NewSequenceModel()
	crs, err := b.Ctx.Find(NewSequenceModel(), tk.M{}.Set("where", db.Eq("collname", collName)))
	if err != nil {
		return -9999, err
	}
	defer crs.Close()
	err = crs.Fetch(mdl, 1, false)
	if err != nil {
		return -9999, err
	}
	ret = mdl.Lastnumber + 1
	mdl.Lastnumber = ret
	b.Ctx.Save(mdl)
	return ret, nil
}

func (b *BaseController) Action(k *knot.WebContext, module string, do string, whatchanged string, oldvalue string, newvalue string, sourcetype string, sources string) {

	data := NewLogModel()
	if k.Session("username") != nil && k.Session("username") != "" {
		data.Do = do
		data.UserID = k.Session("username").(string)
		data.DateAccess = time.Now().UTC()
		data.LoginTime = k.Session("logintime").(time.Time)
		data.SessionID = k.Session("sessionid").(string)
		data.ExpiredTime = k.Session("expiredtime").(time.Time)
		data.Module = module

		data.FullName = k.Session(SESSION_FULLNAME).(string)
		if k.Session(SESSION_FIRSTNAME) != nil {
			data.FirstName = k.Session(SESSION_FIRSTNAME).(string)
		}
		if k.Session(SESSION_LASTNAME) != nil {
			data.LastName = k.Session(SESSION_LASTNAME).(string)
		}
		if k.Session(SESSION_COUNTRY) != nil {
			data.Country = k.Session(SESSION_COUNTRY).(string)
		}
		if data.Country == "" {
			data.Country = "GLOBAL"
		}
		if k.Session(SESSION_GROUP) != nil {
			group_list := k.Session(SESSION_GROUP).([]string)
			group := strings.Join(group_list, ", ")
			data.Group = group

			group_list_wh := []interface{}{}
			for _, x := range group_list {
				group_list_wh = append(group_list_wh, x)
			}
			group_description := []acl.Group{}
			csr, _ := b.AclCtx.Connection.NewQuery().From(new(acl.Group).TableName()).Where(db.In("_id", group_list_wh...)).Cursor(nil)
			csr.Fetch(&group_description, 0, false)
			csr.Close()
			// tk.Println(group_description)
			group_description_list := []string{}
			for _, x := range group_description {
				group_description_list = append(group_description_list, x.Title)
			}
			data.GroupDescription = strings.Join(group_description_list, ", ")
		}

		data.WhatChanged = whatchanged
		data.OldValue = oldvalue
		data.NewValue = newvalue
		data.SourceType = sourcetype
		data.Sources = sources
		data.RequestURI = k.Request.RequestURI
		e := b.AclCtx.Save(data)
		if e != nil {
			return
		}
	}

	return
}
