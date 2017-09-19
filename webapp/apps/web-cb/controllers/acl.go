package controllers

import (
	h "eaciit/scb-apps/webapp/apps/web-cb/helper"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AclController struct {
	*BaseController
}

func (c *AclController) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	c.Action(k, "Configuration", "Open Configuration Page", "", "", "", "", "")
	Configuration := c.GetAccess(k, "APPCONFIG")
	ConfigUser := c.GetAccess(k, "CONFIGUSER")
	ConfigRole := c.GetAccess(k, "CONFIGROLE")
	ApplicationAuditTrail := c.GetAccess(k, "APPLICATIONAUDITTRAIL")
	UserAuditTrail := c.GetAccess(k, "USERAUDITTRAIL")
	Scorecard := c.GetAccess(k, "SCORECARD")
	Initiative := c.GetAccess(k, "INITIATIVE")
	ApplicationUsageDetails := c.GetAccess(k, "APPLICATIONUSAGEDETAILS")
	k.Config.LayoutTemplate = "_layout-v2.html"
	PartialFiles := []string{}
	PartialFiles = append(PartialFiles, "shared/sidebar.html")
	PartialFiles = append(PartialFiles, "shared/initiative_summary.html")
	PartialFiles = append(PartialFiles, "shared/initiative_chart.html")
	PartialFiles = append(PartialFiles, "shared/metric_upload.html")
	PartialFiles = append(PartialFiles, "dashboard/initiative.html")

	PartialFiles = append(PartialFiles, "acl/user_management.html")
	PartialFiles = append(PartialFiles, "acl/role_management.html")
	PartialFiles = append(PartialFiles, "acl/menu_management.html")
	PartialFiles = append(PartialFiles, "acl/user_loginactivity.html")
	PartialFiles = append(PartialFiles, "acl/user_audittrail.html")
	PartialFiles = append(PartialFiles, "acl/app_usagedetails.html")

	k.Config.IncludeFiles = PartialFiles
	k.Config.OutputType = knot.OutputTemplate
	return tk.M{}.Set("Configuration", Configuration).Set("ConfigUser", ConfigUser).Set("ConfigRole", ConfigRole).Set("ApplicationAuditTrail", ApplicationAuditTrail).Set("UserAuditTrail", UserAuditTrail).Set("ApplicationUsageDetails", ApplicationUsageDetails).Set("Scorecard", Scorecard).Set("Initiative", Initiative)
}

func (c *AclController) EncodeDecode(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	c.Action(k, "Configuration", "Open Encode - Decode Page", "", "", "", "", "")
	Configuration := c.GetAccess(k, "APPCONFIG")
	ConfigUser := c.GetAccess(k, "CONFIGUSER")
	ConfigRole := c.GetAccess(k, "CONFIGROLE")
	ApplicationAuditTrail := c.GetAccess(k, "APPLICATIONAUDITTRAIL")
	UserAuditTrail := c.GetAccess(k, "USERAUDITTRAIL")
	Initiative := c.GetAccess(k, "INITIATIVE")
	k.Config.LayoutTemplate = "_layout-v2.html"
	PartialFiles := []string{}
	PartialFiles = append(PartialFiles, "shared/sidebar.html")
	PartialFiles = append(PartialFiles, "shared/initiative_summary.html")
	PartialFiles = append(PartialFiles, "shared/initiative_chart.html")
	PartialFiles = append(PartialFiles, "shared/metric_upload.html")
	PartialFiles = append(PartialFiles, "dashboard/initiative.html")

	PartialFiles = append(PartialFiles, "acl/encrypt.html")
	PartialFiles = append(PartialFiles, "acl/decrypt.html")

	k.Config.IncludeFiles = PartialFiles
	k.Config.OutputType = knot.OutputTemplate
	return tk.M{}.Set("Configuration", Configuration).Set("ConfigUser", ConfigUser).Set("ConfigRole", ConfigRole).Set("ApplicationAuditTrail", ApplicationAuditTrail).Set("UserAuditTrail", UserAuditTrail).Set("Initiative", Initiative)
}

func (c *AclController) Encrypt(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	c.Action(k, "Configuration", "Open Encypt Page", "", "", "", "", "")
	Configuration := c.GetAccess(k, "APPCONFIG")
	ConfigUser := c.GetAccess(k, "CONFIGUSER")
	ConfigRole := c.GetAccess(k, "CONFIGROLE")
	ApplicationAuditTrail := c.GetAccess(k, "APPLICATIONAUDITTRAIL")
	UserAuditTrail := c.GetAccess(k, "USERAUDITTRAIL")
	Initiative := c.GetAccess(k, "INITIATIVE")
	k.Config.LayoutTemplate = "_layout-v2.html"
	PartialFiles := []string{}
	PartialFiles = append(PartialFiles, "shared/sidebar.html")
	PartialFiles = append(PartialFiles, "shared/initiative_summary.html")
	PartialFiles = append(PartialFiles, "shared/initiative_chart.html")
	PartialFiles = append(PartialFiles, "shared/metric_upload.html")
	PartialFiles = append(PartialFiles, "dashboard/initiative.html")

	k.Config.IncludeFiles = PartialFiles
	k.Config.OutputType = knot.OutputTemplate
	return tk.M{}.Set("Configuration", Configuration).Set("ConfigUser", ConfigUser).Set("ConfigRole", ConfigRole).Set("ApplicationAuditTrail", ApplicationAuditTrail).Set("UserAuditTrail", UserAuditTrail).Set("Initiative", Initiative)
}
func (c *AclController) DoEncrypt(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Text string
	}{}
	err := k.GetPayload(&parm)
	result, err := h.Encode(parm.Text)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	return c.SetResultInfo(false, "", result)
}
func (c *AclController) Decrypt(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	c.Action(k, "Configuration", "Open Encypt Page", "", "", "", "", "")
	Configuration := c.GetAccess(k, "APPCONFIG")
	ConfigUser := c.GetAccess(k, "CONFIGUSER")
	ConfigRole := c.GetAccess(k, "CONFIGROLE")
	ApplicationAuditTrail := c.GetAccess(k, "APPLICATIONAUDITTRAIL")
	UserAuditTrail := c.GetAccess(k, "USERAUDITTRAIL")
	Initiative := c.GetAccess(k, "INITIATIVE")
	k.Config.LayoutTemplate = "_layout-v2.html"
	PartialFiles := []string{}
	PartialFiles = append(PartialFiles, "shared/sidebar.html")
	PartialFiles = append(PartialFiles, "shared/initiative_summary.html")
	PartialFiles = append(PartialFiles, "shared/initiative_chart.html")
	PartialFiles = append(PartialFiles, "shared/metric_upload.html")
	PartialFiles = append(PartialFiles, "dashboard/initiative.html")

	k.Config.IncludeFiles = PartialFiles
	k.Config.OutputType = knot.OutputTemplate
	return tk.M{}.Set("Configuration", Configuration).Set("ConfigUser", ConfigUser).Set("ConfigRole", ConfigRole).Set("ApplicationAuditTrail", ApplicationAuditTrail).Set("UserAuditTrail", UserAuditTrail).Set("Initiative", Initiative)
}

func (c *AclController) DoDecrypt(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Text string
	}{}
	err := k.GetPayload(&parm)
	result, err := h.Decode(parm.Text)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	return c.SetResultInfo(false, "", result)
}
