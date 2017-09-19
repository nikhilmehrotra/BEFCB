package controllers

import (
	"github.com/eaciit/knot/knot.v1"
)

type AccessController struct {
	*BaseController
}

func (c *AccessController) Master(k *knot.WebContext) interface{} {
	c.SetResponseTypeHTML(k)
	if !c.ValidateAccessOfRequestedURL(k, "admin") {
		return nil
	}

	return c.SetViewData(nil)
}
