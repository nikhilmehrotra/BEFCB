package webext

import (
	"bufio"
	. "eaciit/scb-apps/webapp/apps/web-cb/controllers"
	"log"
	"os"
	"strings"
	"time"

	"path/filepath"
	"strconv"

	"github.com/eaciit/acl/v2.0"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/orm"
	"github.com/eaciit/toolkit"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
)

func Register() *knot.App {
	log.Println("___INIT START_____")
	conn, err := PrepareConnection()
	if err != nil {
		log.Println(err)
	}

	acl_conn, err := PrepareACLConnection()
	if err != nil {
		log.Println(err)
	}
	uploadPath := PrepareUploadPath()
	config := ReadConfig()
	ctx := orm.New(conn)
	acl_ctx := orm.New(acl_conn)

	baseCtrl := new(BaseController)
	baseCtrl.Ctx = ctx
	baseCtrl.AclCtx = acl_ctx
	baseCtrl.UploadPath = uploadPath
	baseCtrl.DownloadPath = config["downloadPath"]
	baseCtrl.TemplatePath = config["templatePath"]
	baseCtrl.UploadMetricPath = config["metricFilePath"]
	baseCtrl.ServerNameLDAP = config["ServerNameLDAP"]
	baseCtrl.IsUsingLDAP, err = strconv.ParseBool(config["isUsingLDAP"])
	baseCtrl.UserAuthAttrLDAP = config["userauthattrLDAP"]
	baseCtrl.UserDNLDAP = config["userdnLDAP"]

	baseCtrl.BindUsernameLDAP = config["BindUsernameLDAP"]
	baseCtrl.BindPasswordLDAP = config["BindPasswordLDAP"]
	baseCtrl.BindFilterLDAP = config["BindFilterLDAP"]
	baseCtrl.InsecureSkipVerify, err = strconv.ParseBool(config["InsecureSkipVerify"])

	if err != nil {
		log.Println(err)
	}
	baseCtrl.AddressLDAP = config["addressLDAP"]
	baseCtrl.BaseDNLDAP = config["basednLDAP"]
	baseCtrl.LDAPType = config["LDAPType"]
	LDAPCertificates := config["LDAPCertificate"]
	if LDAPCertificates != "" {
		tempLDAPCertificates := strings.Split(LDAPCertificates, "|")
		if len(tempLDAPCertificates) > 0 {
			baseCtrl.LDAPCertificate = []string{}
			for _, i := range tempLDAPCertificates {
				baseCtrl.LDAPCertificate = append(baseCtrl.LDAPCertificate, i)
			}
		}

	}
	//acl db
	acl.SetExpiredDuration(time.Hour * 2)
	err = acl.SetDb(acl_conn)
	if err != nil {
		log.Println(err)
	}

	err = PrepareDefaultUser()
	if err != nil {
		log.Println(err)
	}
	//==
	appName := "web-cb"
	baseCtrl.APP_NAME = appName
	app := knot.NewApp(appName)
	app.ViewsPath = filepath.Join("apps", appName, "views") + "/"
	app.Register(&AclUserController{baseCtrl})
	app.Register(&AclSysAdminController{baseCtrl})
	app.Register(&LoginController{baseCtrl})
	app.Register(&LogoutController{baseCtrl})
	app.Register(&DashboardController{baseCtrl})
	app.Register(&AdoptionModuleController{baseCtrl})

	// app.Register(&MasterLiveCircleController{baseCtrl})
	app.Register(&InitiativeOwnerController{baseCtrl})
	app.Register(&AclController{baseCtrl})
	app.Register(&MetricUploadController{baseCtrl})
	app.Register(&MenuSettingController{baseCtrl})
	app.Register(&MController{baseCtrl})
	app.Register(&InitiativeController{baseCtrl})
	app.Register(&TaskController{baseCtrl})
	app.Register(&ScorecardController{baseCtrl})
	app.Register(&ScorecardDetailController{baseCtrl})
	app.Register(&SearchController{baseCtrl})
	app.Register(&MasterCountryController{baseCtrl})
	app.Register(&MasterRegion{baseCtrl})
	app.Register(&BusinessDriverController{baseCtrl})
	app.Register(&InitiativeMasterController{baseCtrl})
	app.Register(&TestController{baseCtrl})
	// app.Register(&InitiativeDataController{baseCtrl})
	app.Register(&OverviewController{baseCtrl})
	app.Register(&SharedAgendaController{baseCtrl})
	app.Register(&CountryAnalysisController{baseCtrl})
	app.Register(&BusinessMetricsController{baseCtrl})
	app.Register(&StagingAreaController{baseCtrl})
	app.Register(&ScorecardInitiativeController{baseCtrl})
	app.Register(&BEFSponsorController{baseCtrl})
	app.Register(&RegionController{baseCtrl})
	app.Register(&ScorecardAnalysisController{baseCtrl})

	app.Static("static", filepath.Join("apps", appName, "assets"))
	app.LayoutTemplate = "_layout.html"
	knot.RegisterApp(app)
	log.Println("___INIT FINISH_____")
	return app
}

func PrepareConnection() (dbox.IConnection, error) {
	config := ReadConfig()
	ci := &dbox.ConnectionInfo{config["host"], config["database"], config["username"], config["password"], toolkit.M{}.Set("timeout", 10)}
	toolkit.Printfn("Connecting to: %s\n", toolkit.JsonString(ci))
	c, e := dbox.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}

	return c, nil
}
func PrepareACLConnection() (dbox.IConnection, error) {
	config := ReadConfig()
	ci := &dbox.ConnectionInfo{config["host"], config["acldatabase"], config["username"], config["password"], toolkit.M{}.Set("timeout", 10)}
	toolkit.Printfn("Connecting to: %s\n", toolkit.JsonString(ci))
	c, e := dbox.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}

	return c, nil
}

func PrepareUploadPath() string {
	config := ReadConfig()
	return config["uploadPath"]
}

func ReadConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(wd + "apps/web-cb/conf/app.conf")
	if err == nil {
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}

			sval := strings.Split(string(line), "=")
			if len(sval) > 2 {
				var tempval = ""
				for xi, x := range sval {
					if xi != 0 {
						if xi == 1 {
							tempval += x
						} else {
							tempval += "=" + x
						}
					}
				}
				ret[sval[0]] = tempval
			} else if len(sval) > 0 {
				ret[sval[0]] = sval[1]
			}
		}
	} else {
		log.Println(err.Error())
	}

	return ret
}

func PrepareDefaultUser() (err error) {
	username := "eaciit"
	password := "Password.1"

	user := new(acl.User)
	err = acl.FindUserByLoginID(user, username)

	if err == nil || user.LoginID == username {
		return
	}

	user.ID = toolkit.RandomString(32)
	user.LoginID = username
	user.FullName = username
	user.Password = password
	user.Enable = true

	err = acl.Save(user)
	if err != nil {
		return
	}
	err = acl.ChangePassword(user.ID, password)
	if err != nil {
		return
	}

	toolkit.Printf(`Default user "%s" with standard password has been created%s`, username, "\n")

	return
}
