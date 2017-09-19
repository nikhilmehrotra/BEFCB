package controllers

import (
	// m "eaciit/scb-apps/webapp/apps/web-cb/models"
	// db "github.com/eaciit/dbox"
	. "eaciit/scb-apps/webapp/apps/web-cb/helper"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// "strings"
	// "strconv"
	// bson "gopkg.in/mgo.v2/bson"
	// "time"
)

type MetricUploadController struct {
	*BaseController
}

func (c *MetricUploadController) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	c.Action(k, "Metric Upload", "Open Metric Upload Page", "", "", "", "", "")
	MetricUpload := c.GetAccess(k, "METRICUPLOAD")
	Initiative := c.GetAccess(k, "INITIATIVE")
	Scorecard := c.GetAccess(k, "SCORECARD")
	k.Config.OutputType = knot.OutputTemplate
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
	k.Config.IncludeFiles = PartialFiles
	return tk.M{}.Set("MetricUpload", MetricUpload).Set("Initiative", Initiative).Set("Scorecard", Scorecard)
}
func (c *MetricUploadController) GetData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	username := k.Session("username").(string)

	// tk.Println("username", username)

	result := []tk.M{}
	pipes := []tk.M{}
	unwind2 := tk.M{}
	unwind1 := tk.M{}
	project := tk.M{}
	sort := tk.M{}

	err := Deserialize(`
	{"$unwind": "$businessmetric"}
	`, &unwind1)
	err = Deserialize(`
	{"$unwind":"$businessmetric.MetricFiles"}
	`, &unwind2)
	err = Deserialize(`
	{"$project": 
        {       
            "Name" : "$name",
            "KeyMetrics" : "$businessmetric.description",
            "Period": "$businessmetric.MetricFiles.monthyear", 
            "FileName": "$businessmetric.MetricFiles.filename",
            "UploadedDate": "$businessmetric.MetricFiles.uploadedat",
            "UploadedBy": "$businessmetric.MetricFiles.uploadername"
        } 
    }
	`, &project)

	err = Deserialize(`
	{"$sort" : {"UploadedDate": -1}}
	`, &sort)

	pipes = append(pipes, unwind1)
	pipes = append(pipes, unwind2)
	pipes = append(pipes, tk.M{"$match": tk.M{"businessmetric.MetricFiles.uploadername": username}})
	pipes = append(pipes, project)
	pipes = append(pipes, sort)

	// result := []tk.M{}

	csr, err := c.Ctx.Connection.NewQuery().Command("pipe", pipes).
		From("BusinessDriverL1"). //Where(db.Eq("businessmetric.metricfiles.uploadername", username)).
		Cursor(nil)
	err = csr.Fetch(&result, 0, true)
	csr.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)

	}

	// result.Set("Metrics", result)

	return c.SetResultInfo(false, "", result)
}
