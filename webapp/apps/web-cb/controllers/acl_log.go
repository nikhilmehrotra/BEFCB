package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"github.com/eaciit/knot/knot.v1"
	// tk "github.com/eaciit/toolkit"
	"time"
)

func (c *AclController) SaveLog(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		NewValue    string
		OldValue    string
		Whatchanged string
		Do          string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	data := NewLogModel()

	if k.Session("username") != nil && k.Session("username") != "" {
		// tk.Println("DO", parm.Do)
		data.Do = parm.Do
		data.UserID = k.Session("username").(string)
		data.DateAccess = time.Now().UTC()
		data.LoginTime = k.Session("logintime").(time.Time)
		data.SessionID = k.Session("sessionid").(string)
		data.ExpiredTime = k.Session("expiredtime").(time.Time)
		data.WhatChanged = parm.Whatchanged
		data.OldValue = parm.OldValue
		data.NewValue = parm.NewValue

		e := c.AclCtx.Save(data)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}

	}

	return c.SetResultInfo(false, "", data)
}

func (c *AclController) SaveLogArray(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputTemplate
	type Logs struct {
		NewValue    string
		OldValue    string
		Whatchanged string
		Do          string
	}
	parm := struct {
		LogAction []Logs
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	data := NewLogModel()

	if k.Session("username") != nil && k.Session("username") != "" {
		for _, la := range parm.LogAction {
			data.Do = la.Do
			data.UserID = k.Session("username").(string)
			data.DateAccess = time.Now().UTC()
			data.LoginTime = k.Session("logintime").(time.Time)
			data.SessionID = k.Session("sessionid").(string)
			data.ExpiredTime = k.Session("expiredtime").(time.Time)
			data.WhatChanged = la.Whatchanged
			data.OldValue = la.OldValue
			data.NewValue = la.NewValue

			e := c.AclCtx.Save(data)
			if e != nil {
				return c.ErrorResultInfo(e.Error(), nil)
			}
		}
	}

	return c.SetResultInfo(false, "", data)
}
