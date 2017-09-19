package controllers

import (
	"github.com/eaciit/knot/knot.v1"
)

type MController struct {
	*BaseController
}

func (a *MController) Default(k *knot.WebContext) interface{} {
	// a.LoadBase(k)
	// a.LoadPartial(k, []string{})
	k.Config.OutputType = knot.OutputTemplate
	return ""
}
