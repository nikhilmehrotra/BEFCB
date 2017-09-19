package app_main

import (
	"eaciit/scb-apps/webapp/apps/main/controllers"
	"eaciit/scb-apps/webapp/helper"
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func init() {
	type ForgetMe struct{}

	appFolderPath := helper.GetCurrentFolderPath(ForgetMe{})
	appName := helper.GetCurrentFolderName(ForgetMe{})

	// ==== start
	helper.Println("Registering", appName, "@", appFolderPath)

	// ==== get config
	config := helper.ReadConfig(ForgetMe{})

	// ==== prepare database connection
	conn, err := helper.PrepareConnection(ForgetMe{})
	if err != nil {
		helper.Println(err.Error())
		os.Exit(0)
	}

	// ==== configure acl
	acl.SetExpiredDuration(time.Second * time.Duration(config.GetFloat64("loginexpired")))
	err = acl.SetDb(conn)
	if err != nil {
		helper.Println(err.Error())
		os.Exit(0)
	}

	// ==== save connection to controller context
	ctx := orm.New(conn)
	baseCtrl := new(controllers.BaseController)
	// baseCtrl.NoLogin = true
	baseCtrl.Conn = conn
	baseCtrl.AppName = appName
	baseCtrl.Ctx = ctx
	baseCtrl.IsUsingLDAP, err = strconv.ParseBool(config.GetString("isUsingLDAP"))
	baseCtrl.ServerNameLDAP = config.GetString("ServerNameLDAP")
	baseCtrl.UserAuthAttrLDAP = config.GetString("userauthattrLDAP")
	baseCtrl.UserDNLDAP = config.GetString("userdnLDAP")
	appPath := config.GetString("AppPath")
	if strings.Trim(appPath, " ") != "" {
		appFolderPath = appPath
	}

	baseCtrl.BindUsernameLDAP = config.GetString("BindUsernameLDAP")

	baseCtrl.IsUsingEncryptedPassword, err = strconv.ParseBool(config.GetString("IsUsingEncryptedPassword"))
	if err != nil {
		log.Println(err)
	}

	baseCtrl.BindPasswordLDAP = config.GetString("BindPasswordLDAP")
	if strings.Trim(baseCtrl.BindPasswordLDAP, " ") != "" && baseCtrl.IsUsingEncryptedPassword {
		BindPasswordLDAPResult, err := helper.Decode(baseCtrl.BindPasswordLDAP)
		if err == nil {
			baseCtrl.BindPasswordLDAP = BindPasswordLDAPResult
		}
	}
	baseCtrl.BindFilterLDAP = config.GetString("BindFilterLDAP")
	baseCtrl.InsecureSkipVerify, err = strconv.ParseBool(config.GetString("InsecureSkipVerify"))
	if err != nil {
		log.Println(err)
	}

	baseCtrl.AddressLDAP = config.GetString("addressLDAP")
	baseCtrl.BaseDNLDAP = config.GetString("basednLDAP")
	baseCtrl.LDAPType = config.GetString("LDAPType")
	LDAPCertificates := config.GetString("LDAPCertificate")
	if LDAPCertificates != "" {
		tempLDAPCertificates := strings.Split(LDAPCertificates, "|")
		if len(tempLDAPCertificates) > 0 {
			baseCtrl.LDAPCertificate = []string{}
			for _, i := range tempLDAPCertificates {
				baseCtrl.LDAPCertificate = append(baseCtrl.LDAPCertificate, i)
			}
		}

	}

	baseCtrl.LDAP_DATA.LDAP_DATA_FullName = config.GetString("LDAP_DATA_FullName")
	baseCtrl.LDAP_DATA.LDAP_DATA_FirstName = config.GetString("LDAP_DATA_FirstName")
	baseCtrl.LDAP_DATA.LDAP_DATA_LastName = config.GetString("LDAP_DATA_LastName")
	baseCtrl.LDAP_DATA.LDAP_DATA_UserCountry = config.GetString("LDAP_DATA_UserCountry")

	// create default access data for the first time
	err = helper.PrepareDefaultData()
	if err != nil {
		helper.Println(err.Error())
	}

	// create the application
	app := knot.NewApp(appName)
	app.LayoutTemplate = "_layout.html"
	app.ViewsPath = filepath.Join(appFolderPath, "views") + tk.PathSeparator
	helper.Println("Configure view location", app.ViewsPath)

	// register routes
	app.Register(&(controllers.AuthController{BaseController: baseCtrl}))
	app.Register(&(controllers.AdmController{BaseController: baseCtrl}))
	app.Register(&(controllers.AccessController{BaseController: baseCtrl}))
	app.Register(&(controllers.DashboardController{BaseController: baseCtrl}))
	app.Register(&(controllers.PublicController{BaseController: baseCtrl}))
	app.Static("static", filepath.Join(appFolderPath, "assets"))

	knot.RegisterApp(app)
}
