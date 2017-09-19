package controllers

import (
	"eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	//	"strings"

	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	//	"gopkg.in/mgo.v2/bson"
)

type LoginController struct {
	*BaseController
}

func (c *LoginController) Default(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.LayoutTemplate = ""
	return ""
}

func (c *LoginController) Do(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson
	formData := struct {
		UserName   string
		Password   string
		RememberMe bool
	}{}
	message := ""
	isValid := false
	err := k.GetPayload(&formData)
	if err != nil {
		c.WriteLog(err)
		message = "Backend Error " + err.Error()
	}
	q := tk.M{}.Set("where", db.Eq("username", formData.UserName))
	cur, err := c.Ctx.Find(new(SysUserModel), q)
	if err != nil {
		return tk.M{}.Set("Valid", false).Set("Message", err.Error())
	}
	res := make([]SysUserModel, 0)
	//	defer c.Ctx.Close()
	defer cur.Close()
	err = cur.Fetch(&res, 0, false)
	if err != nil {
		return tk.M{}.Set("Valid", false).Set("Message", err.Error())
	}

	if len(res) > 0 {
		resUser := res[0]
		if helper.GetMD5Hash(formData.Password) == resUser.Password {
			if resUser.Enable == true {

				resroles := make([]SysRolesModel, 0)
				crsR, errR := c.Ctx.Find(new(SysRolesModel), tk.M{}.Set("where", db.Eq("name", resUser.Roles)))
				if errR != nil {
					return c.SetResultInfo(true, errR.Error(), nil)
				}
				errR = crsR.Fetch(&resroles, 0, false)
				if errR != nil {
					return c.SetResultInfo(true, errR.Error(), nil)
				}
				defer crsR.Close()

				k.SetSession("userid", string(resUser.Id))
				k.SetSession("username", resUser.Username)
				k.SetSession("usermodel", resUser)
				k.SetSession("roles", resroles)
				isValid = true

				loginlog := resUser.GetLoginLog()
				_ = c.Ctx.Save(loginlog)

			} else {
				message = "Your account is disabled, please contact administrator to enable it."
			}
		} else {
			message = "Invalid Username or password!"
		}
	} else {
		return "Invalid Username or password!"
	}

	return tk.M{}.Set("Valid", isValid).Set("Message", message)
}

// func (l *LoginController) ProcessLogin(r *knot.WebContext) interface{} {
// 	r.Config.OutputType = knot.OutputJson

// 	payload := toolkit.M{}
// 	if err := r.GetPayload(&payload); err != nil {
// 		return helper.CreateResult(false, "", err.Error())
// 	}
// ================
// func (c *LoginController) Do(k *knot.WebContext) interface{} {
// 	k.Config.NoLog = true
// 	k.Config.OutputType = knot.OutputJson
// 	formData := struct {
// 		UserName   string
// 		Password   string
// 		RememberMe bool
// 	}{}
// 	message := ""
// 	isValid := false
// 	err := k.GetPayload(&formData)

// func (c *LoginController) SaveNewPassword(k *knot.WebContext) interface{} {
// 	k.Config.NoLog = true
// 	k.Config.OutputType = knot.OutputJson

// 	payload := toolkit.M{}
// 	if err := k.GetPayload(&payload); err != nil {
// 		return tk.M{}.Set("Valid", false).Set("Message", err.Error())
// 	}

// 	if !payload.Has("newpassword") || !payload.Has("userid") {
// 		return errors.New("Data is not complete")
// 	}

// 	switch {
// 	case payload.Has("tokenid"):
// 		acl.ChangePasswordToken(toolkit.ToString(payload["userid"]), toolkit.ToString(payload["newpassword"]), toolkit.ToString(payload["tokenid"]))
// 	default:
// 		// check sessionid first
// 		savedsessionid := "" //change with get session
// 		//=======================
// 		userid, err := acl.FindUserBySessionID(savedsessionid)
// 		if err == nil && userid == toolkit.ToString(payload["userid"]) {
// 			err = acl.ChangePassword(toolkit.ToString(payload["userid"]), toolkit.ToString(payload["newpassword"]))
// 		} else if err == nil {
// 			err = errors.New("Userid is not match")
// 		}
// 	}

// 	return nil
// }

// func ResetPassword(payload toolkit.M) error {
// 	if !payload.Has("email") || !payload.Has("baseurl") {
// 		errors.New("Data is not complete")
// 	}

// 	uname, tokenid, err := acl.ResetPassword(toolkit.ToString(payload["email"]))

// 	if err != nil {
// 		err.Error()
// 	}

// 	linkstr := toolkit.Sprintf("<a href='%v/web/confirmreset?1=%v&2=%v'>Click</a>", toolkit.ToString(payload["baseurl"]), uname, tokenid)

// 	mailmsg := toolkit.Sprintf("Hi, <br/><br/> We received a request to reset your password, <br/><br/>")
// 	mailmsg = toolkit.Sprintf("%vFollow the link below to set a new password : <br/><br/> %v <br/><br/>", mailmsg, linkstr)
// 	mailmsg = toolkit.Sprintf("%vIf you don't want to change your password, you can ignore this email <br/><br/> Thanks,</body></html>", mailmsg)

// 	m := gomail.NewMessage()

// 	m.SetHeader("From", "admin.support@eaciit.com")
// 	m.SetHeader("To", toolkit.ToString(payload["email"]))

// 	m.SetHeader("Subject", "[no-reply] Self password reset")
// 	m.SetBody("text/html", mailmsg)

// 	d := gomail.NewPlainDialer("smtp.office365.com", 587, "admin.support@eaciit.com", "******")
// 	err = d.DialAndSend(m)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
