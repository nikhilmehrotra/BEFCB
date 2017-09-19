package controllers

import (
	"github.com/eaciit/knot/knot.v1"
	"net/http"
	"strings"
)

type DashboardController struct {
	*BaseController
}

func (c *DashboardController) Index(k *knot.WebContext) interface{} {
	c.SetResponseTypeHTML(k)
	if !c.ValidateAccessOfRequestedURL(k) {
		apps, err := c.GetApplicationByUserName(c.GetCurrentUsername(k))
		if err != nil {
			return c.SetResultError(err.Error(), nil)
		}
		if len(apps) == 1 {
			redirect := `/` + apps[0].ID + `/` + strings.Trim(apps[0].LandingURL, ` /`)
			http.Redirect(k.Writer, k.Request, redirect, http.StatusTemporaryRedirect)
		}
		return nil
	}
	return c.SetViewData(nil)
}
