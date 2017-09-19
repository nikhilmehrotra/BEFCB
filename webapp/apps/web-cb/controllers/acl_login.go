package controllers

import (
	"crypto/x509"
	"errors"
	"fmt"
	// "github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	// "github.com/eaciit/ldap"
	// "github.com/eaciit/orm/v1"
	"crypto/md5"
	"crypto/tls"
	"io"
	"io/ioutil"
	"strings"
	"time"

	"github.com/eaciit/acl/v2.0"
	"github.com/eaciit/ldap"
	"github.com/eaciit/toolkit"
)

var _aclctxErr error
var _expiredduration time.Duration

type IDTypeEnum int

const (
	IDTypeUser IDTypeEnum = iota
	IDTypeGroup
	IDTypeSession
)

func init() {
	_expiredduration = time.Minute * 30
}

func Login(username string, password string, IsUsingLDAP bool, LDAPType string, AddressLDAP string, BaseDNLDAP string, ServerNameLDAP string, LDAPCertificate []string, UserAuthAttrLDAP string, UserDNLDAP string, BindUsernameLDAP string, BindPasswordLDAP string, InsecureSkipVerify bool, BindFilterLDAP string) (sessionid string, err error) {

	tUser := new(acl.User)
	err = acl.FindUserByLoginID(tUser, username)
	toolkit.Println("# Login with Login ID : ", username)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			err = errors.New("Invalid Login ID / Password")
			return
		}
		err = errors.New(fmt.Sprintf("Found error : %v", err.Error()))
		return
	}

	if tUser.ID == "" {
		err = errors.New("Invalid Login ID / Password")
		return
	}
	LoginSuccess := false
	switch tUser.LoginType {
	case acl.LogTypeLdap:
		toolkit.Println("# ", username, "- Login using LDAP")
		if !IsUsingLDAP {
			err = errors.New("LDAP feature is disabled")
			return
		}
		tUser.LoginConf = toolkit.M{}
		tUser.LoginConf.Set("address", AddressLDAP)
		tUser.LoginConf.Set("basedn", BaseDNLDAP)
		if LDAPType != "" {
			toolkit.Println("#Type :", LDAPType)
			toolkit.Println("#ServerName : ", ServerNameLDAP)
			toolkit.Println("#Address : ", AddressLDAP)
			toolkit.Println("#BaseDN : ", BaseDNLDAP)
			if UserAuthAttrLDAP != "" {
				toolkit.Println("#UserAuthAttrLDAP : ", UserAuthAttrLDAP)
			}
			if UserDNLDAP != "" {
				toolkit.Println("#UserDNLDAP : ", UserDNLDAP)
			}
			toolkit.Println("#Username : ", username)
			tUser.LoginConf.Set("type", LDAPType)
			tlsconfig := tls.Config{}
			if ServerNameLDAP != "" {
				tlsconfig.ServerName = ServerNameLDAP
			}
			tlsconfig.InsecureSkipVerify = InsecureSkipVerify
			tlsconfig.Certificates = []tls.Certificate{}
			caCertPool := x509.NewCertPool()
			for x, c := range LDAPCertificate {
				file, err := ioutil.ReadFile(c)
				if err != nil {
					err = errors.New(fmt.Sprintf("Found error : %v", err.Error()))
				}

				s := caCertPool.AppendCertsFromPEM(file)
				toolkit.Println("Certificate ", (x + 1), " # Added : ", s)
			}
			tlsconfig.RootCAs = caCertPool
			tUser.LoginConf.Set("tlsconfig", &tlsconfig)
		}
		UserAuthAttr := "CN"
		if UserAuthAttrLDAP != "" {
			UserAuthAttr = UserAuthAttrLDAP
			username = UserAuthAttr + "=" + username
		}
		if UserDNLDAP != "" {
			if UserAuthAttrLDAP == "" {
				username = UserAuthAttr + "=" + username
			}
			username = username + "," + UserDNLDAP
		}
		if strings.Trim(BindUsernameLDAP, " ") != "" {
			toolkit.Println("Binding using : ", BindUsernameLDAP)
			LoginSuccess, err = CheckLoginLDAP(username, password, tUser.LoginConf, BindUsernameLDAP, BindPasswordLDAP, UserAuthAttrLDAP, UserDNLDAP, BindFilterLDAP)
		} else {
			toolkit.Println("Binding using : ", username)
			LoginSuccess, err = CheckLoginLDAP(username, password, tUser.LoginConf, BindUsernameLDAP, BindPasswordLDAP, UserAuthAttrLDAP, UserDNLDAP, BindFilterLDAP)
		}
		break
	case acl.LogTypeBasic:
		toolkit.Println("# ", username, "- Login using Basic Conf")
		LoginSuccess, err = CheckLoginBasic(password, tUser.Password)
		break
	default:
		toolkit.Println("# ", username, "- Login using Basic Conf")
		LoginSuccess, err = CheckLoginBasic(password, tUser.Password)
		break
	}

	if !LoginSuccess {
		return
	}

	if !tUser.Enable {
		err = errors.New("LoginID is not active")
		return
	}

	tSession := new(acl.Session)
	tSession.ID = toolkit.RandomString(32)
	tSession.UserID = tUser.ID
	tSession.LoginID = tUser.LoginID
	tSession.Created = time.Now().UTC()
	tSession.Expired = time.Now().UTC().Add(_expiredduration)

	err = acl.Save(tSession)
	if err == nil {
		sessionid = tSession.ID
	}
	toolkit.Println("# Login SUCCESS")
	return
}

func CheckLoginLDAP(username string, password string, loginconf toolkit.M, BindUsernameLDAP string, BindPasswordLDAP string, UserAuthAttrLDAP string, UserDNLDAP string, BindFilterLDAP string) (cond bool, err error) {
	toolkit.Println("# Connecting to LDAP")
	cond = false
	connectTime := time.Now()
	address := loginconf.GetString("address")
	l := GetLDAPConn(address, loginconf)
	defer l.Close()
	err = l.Connect()
	if err != nil {
		toolkit.Println("#ERROR Connecting to LDAP", err.Error())
		return
	}
	if strings.Trim(BindUsernameLDAP, " ") != "" {
		// Bind Through Config
		err = l.Bind(BindUsernameLDAP, BindPasswordLDAP)
	} else {
		err = l.Bind(username, password)
	}
	if strings.Trim(BindUsernameLDAP, " ") != "" && err == nil {
		// from Login FORM
		err = l.Bind(username, password)
		if err == nil {
			cond = true
		} else {
			toolkit.Println("#ERROR Binding to LDAP with username : ", username, " - ", err.Error())
		}
	} else {
		if err == nil {
			cond = true
		} else {
			toolkit.Println("#ERROR Binding to LDAP  with username : ", username, " - ", err.Error())
		}
	}
	toolkit.Println("# Closing LDAP Connection")
	toolkit.Println("# Connection Time : ", time.Since(connectTime).Seconds(), "s")

	return

}

func GetLDAPConn(address string, config toolkit.M) *ldap.Connection {
	l := new(ldap.Connection)

	switch config.GetString("type") {
	case "ssl":
		l = ldap.NewSSLConnection(address, config.Get("tlsconfig", nil).(*tls.Config))
	case "tls":
		l = ldap.NewSSLConnection(address, config.Get("tlsconfig", nil).(*tls.Config))
	default:
		l = ldap.NewConnection(address)
	}

	return l
}

func CheckLoginBasic(spassword, upassword string) (cond bool, err error) {
	cond = false

	tPass := md5.New()
	io.WriteString(tPass, spassword)
	ePassword := toolkit.Sprintf("%x", tPass.Sum(nil))
	if ePassword == upassword {
		cond = true
	} else {
		err = errors.New("Invalid Login ID / Password")
	}
	return
}

func FindDataLdap(addr, basedn, filter string, param toolkit.M) (arrtkm []toolkit.M, err error) {
	arrtkm = make([]toolkit.M, 0, 0)

	l := GetLDAPConn(addr, param)
	err = l.Connect()
	if err != nil {
		return
	}
	defer l.Close()

	if param.Has("username") {
		err = l.Bind(toolkit.ToString(param["username"]), toolkit.ToString(param["password"]))
		if err != nil {
			return
		}
	}

	attributes := make([]string, 0, 0)
	if param.Has("attributes") {
		attributes = param["attributes"].([]string)
	}
	// filter = "(*" + filter + "*)"
	search := ldap.NewSearchRequest(basedn,
		ldap.ScopeWholeSubtree,
		ldap.DerefAlways,
		0,
		0,
		false,
		filter,
		attributes,
		nil)

	sr, err := l.Search(search)

	for _, v := range sr.Entries {
		tkm := toolkit.M{}

		for _, str := range attributes {
			if len(v.GetAttributeValues(str)) > 1 {
				tkm.Set(str, v.GetAttributeValues(str))
			} else {
				tkm.Set(str, v.GetAttributeValue(str))
			}
		}

		if len(tkm) > 0 {
			arrtkm = append(arrtkm, tkm)
		}
	}

	return
}
