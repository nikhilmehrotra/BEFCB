package controllers

import (
	"eaciit/scb-apps/webapp/apps/main/models"
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type PublicController struct {
	*BaseController
}

func (c *PublicController) GetApplication(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	user := new(models.UserModel)
	err := acl.FindByID(user, LoginDataUser.ID)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	result, err := c.GetApplicationByUserName(LoginDataUser.LoginID)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(result)
}

func (c *PublicController) GetAccessMenu(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := struct{ ApplicationID string }{""}
	err := k.GetPayload(&payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	allowedAccessMenu, err := c.GetAccessMenuByApplicationID(payload.ApplicationID)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	allMenu, err := acl.GetListMenuBySessionId(k.Session(SESSION_KEY, ""))
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	result := make([]tk.M, 0)
	for _, menu := range allMenu {
		for _, allowed := range allowedAccessMenu {
			if menu.GetString("_id") == allowed.ID {
				result = append(result, menu)
			}
		}
	}

	return c.SetResultOK(result)
}
