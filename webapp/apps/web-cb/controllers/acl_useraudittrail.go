package controllers

import (
	// helper "eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"

	"github.com/eaciit/acl/v2.0"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
	// bson "gopkg.in/mgo.v2/bson"
	// "strings"
	"time"
)

// Function List
// acl/getusermetricsreport
// acl/getrolemetricsreport
// acl/getaudittrailreport

func (c *AclController) GetUserMetricsReport(k *knot.WebContext) interface{} {
	// c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		tk.Printf(err.Error())
	}

	result := []tk.M{}

	d := new(acl.User)
	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.Startwith("groups", "CB_")).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	now := time.Now().UTC()
	times := now.Format("01/02/2006")
	rundate := "Run Date : "
	rundates := rundate + times

	//loginid
	users := k.Session("username").(string)
	user := "User : "
	username := user + users

	font := xlsx.NewFont(11, "Calibri")
	style := xlsx.NewStyle()
	style.Font = *font

	fontHdr := xlsx.NewFont(11, "Calibri")
	fontHdr.Bold = true
	styleHdr := xlsx.NewStyle()
	styleHdr.Font = *fontHdr

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Appilcation : BEF CB"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = ""

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = rundates

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = ""

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = username

	row = sheet.AddRow()
	row = sheet.AddRow()

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Login ID"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "BankId"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Roles"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Country"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Status"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Last Login"

	for _, data := range result {

		loginid := data.GetString("loginid")

		results := []tk.M{}
		csr, e := c.AclCtx.Connection.NewQuery().From("acl_log").Where(db.Eq("userid", loginid)).Order("-logintime").Cursor(nil)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
		e = csr.Fetch(&results, 1, false)
		csr.Close()
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}

		// tk.Println("RESULT", results)

		group := data.Get("groups").([]interface{})
		// grup := []string{}

		// for _, i := range group {
		// 	g := i.(string)
		// 	grup = append(grup, g)
		// 	// tk.Println("grup", grup)
		// }
		// groups := strings.Join(grup, ",")
		// gr := strings.Replace(groups, "CB_", "", -1)
		gr := ""
		if len(group) > 0 {
			selectedgroup := new(acl.Group)
			for _, x := range group {
				csr, e := c.AclCtx.Connection.NewQuery().From(new(acl.Group).TableName()).Where(db.Eq("_id", x)).Cursor(nil)
				if e != nil {
					return c.ErrorResultInfo(e.Error(), nil)
				}
				e = csr.Fetch(&selectedgroup, 1, false)
				csr.Close()
				gr = selectedgroup.Title
			}
		}

		row = sheet.AddRow()
		// cell = row.AddCell()
		// cell.SetStyle(style)
		// cell.Value = data.GetString("id")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("loginid")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("loginid")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = gr

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("country")

		cell = row.AddCell()
		cell.SetStyle(style)
		enable := data.Get("enable").(bool)
		if enable {
			cell.Value = "ENABLED"
		} else {
			cell.Value = "DISABLED"
		}

		for _, i := range results {

			logintime, _ := time.Parse(time.RFC3339, i.GetString("logintime"))
			dates := logintime.Format("01/02/2006")

			cell = row.AddCell()
			cell.SetStyle(style)
			cell.Value = dates
		}

	}
	ExcelFilename := "User Metrics Report " + time.Now().Format("20060102150405") + ".xlsx"
	err = file.Save(c.DownloadPath + "/" + ExcelFilename)

	if err != nil {
		tk.Printf(err.Error())
	}

	return c.SetResultInfo(false, "", ExcelFilename)
}

func (c *AclController) GetRoleMetricsReport(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		tk.Printf(err.Error())
	}

	// type Filters struct {
	// 	Field    string
	// 	Operator string
	// 	Value    string
	// 	Type     string
	// }
	// type Filter struct {
	// 	Filters []Filters
	// 	Logic   string
	// }
	// type Sort struct {
	// 	Field string
	// 	Dir   string
	// }
	// p := struct {
	// 	Filter Filter
	// 	Sort   []Sort
	// }{}
	// e := k.GetPayload(&p)

	result := []tk.M{}
	// c.Acl AclCtx.Find(, parms)
	d := new(acl.Group)
	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.Startwith("_id", "CB_")).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	font := xlsx.NewFont(11, "Calibri")
	style := xlsx.NewStyle()
	style.Font = *font
	style.Alignment.WrapText = true
	style.Alignment.Vertical = "top"

	fontHdr := xlsx.NewFont(11, "Calibri")
	fontHdr.Bold = true
	styleHdr := xlsx.NewStyle()
	styleHdr.Font = *fontHdr

	fontTitle := xlsx.NewFont(11, "Calibri")
	fontTitle.Bold = true
	styleTitle := xlsx.NewStyle()
	styleTitle.Font = *fontTitle
	styleTitle.Alignment.Horizontal = "center"
	styleTitle.Alignment.Vertical = "middle"

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Application:BEF CB"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = " "

	tnow := time.Now()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Run Date: " + tnow.Format("02/01/2006")

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = " "

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "User: " + k.Session("username").(string)

	// === new row

	row = sheet.AddRow()

	// === new row
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleTitle)
	cell.Value = "Role"

	cell = row.AddCell()
	cell.SetStyle(styleTitle)
	cell.Value = "Role Description"

	cell = row.AddCell()
	cell.SetStyle(styleTitle)
	cell.Value = "Permissions"

	cell = row.AddCell()
	cell.SetStyle(styleTitle)
	cell.Value = "Status"

	// === new row

	for _, data := range result {
		id := data.GetString("id")

		AccessibilityData := []AccessibilityModel{}
		csr, e = c.AclCtx.Connection.NewQuery().From(new(AccessibilityModel).TableName()).Where(db.Eq("roleid", id)).Cursor(nil)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
		e = csr.Fetch(&AccessibilityData, 0, false)
		csr.Close()

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("title")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = ""

		granttext := ""
		granttext2 := ""
		for _, data2 := range AccessibilityData {
			Create, Read, Update, Delete, Owned, Curtain, Upload := "", "", "", "", "", "", ""

			if data2.Global.Create {
				Create += "G"
			}
			if data2.Global.Read {
				Read += "G"
			}
			if data2.Global.Update {
				Update += "G"
			}
			if data2.Global.Delete {
				Delete += "G"
			}
			if data2.Global.Owned {
				Owned += "G"
			}
			if data2.Global.Curtain {
				Curtain += "G"
			}
			if data2.Global.Upload {
				Upload += "G"
			}

			if data2.Region.Create {
				if Create != "" {
					Create += ","
				}
				Create += "R"
			}
			if data2.Region.Read {
				if Read != "" {
					Read += ","
				}
				Read += "R"
			}
			if data2.Region.Update {
				if Update != "" {
					Update += ","
				}
				Update += "R"
			}
			if data2.Region.Delete {
				if Delete != "" {
					Delete += ","
				}
				Delete += "R"
			}
			if data2.Region.Owned {
				if Owned != "" {
					Owned += ","
				}
				Owned += "R"
			}
			if data2.Region.Curtain {
				if Curtain != "" {
					Curtain += ","
				}
				Curtain += "R"
			}
			if data2.Region.Upload {
				if Upload != "" {
					Upload += ","
				}
				Upload += "R"
			}

			if data2.Country.Create {
				if Create != "" {
					Create += ","
				}
				Create += "C"
			}
			if data2.Country.Read {
				if Read != "" {
					Read += ","
				}
				Read += "C"
			}
			if data2.Country.Update {
				if Update != "" {
					Update += ","
				}
				Update += "C"
			}
			if data2.Country.Delete {
				if Delete != "" {
					Delete += ","
				}
				Delete += "C"
			}
			if data2.Country.Owned {
				if Owned != "" {
					Owned += ","
				}
				Owned += "C"
			}
			if data2.Country.Curtain {
				if Curtain != "" {
					Curtain += ","
				}
				Curtain += "C"
			}
			if data2.Country.Upload {
				if Upload != "" {
					Upload += ","
				}
				Upload += "C"
			}

			if Create != "" {
				if granttext2 != "" {
					granttext2 += ", "
				}
				granttext2 += "Create(" + Create + ")"
			}

			if Read != "" {
				if granttext2 != "" {
					granttext2 += ", "
				}
				granttext2 += "Read(" + Read + ")"
			}

			if Update != "" {
				if granttext2 != "" {
					granttext2 += ", "
				}
				granttext2 += "Update(" + Update + ")"
			}

			if Delete != "" {
				if granttext2 != "" {
					granttext2 += ", "
				}
				granttext2 += "Delete(" + Delete + ")"
			}

			if Owned != "" {
				if granttext2 != "" {
					granttext2 += ", "
				}
				granttext2 += "Owned(" + Owned + ")"
			}

			if Curtain != "" {
				if granttext2 != "" {
					granttext2 += ", "
				}
				granttext2 += "Curtain(" + Curtain + ")"
			}

			if Upload != "" {
				if granttext2 != "" {
					granttext2 += ", "
				}
				granttext2 += "Upload(" + Upload + ")"
			}

			if granttext2 != "" {
				granttext += data2.Title
				granttext += " : "
				granttext += granttext2
				granttext += "\n"
			}

			granttext2 = ""

		}

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = granttext

		cell = row.AddCell()
		cell.SetStyle(style)
		var status = data.Get("enable").(bool)
		var Statuss = ""
		// fmt.Println(status)
		if status {
			Statuss = "ENABLED"
		} else {
			Statuss = "DISABLED"
		}
		cell.Value = Statuss
	}

	// cell = row.AddCell()
	// cell.SetStyle(styleHdr)
	// cell.Value = "Expired Date"

	// for _, data := range result {
	// 	created := data.Get("created").(time.Time).UTC().Format("Monday, Jan _2 2006, 15:04:05")
	// 	expired := data.Get("expired").(time.Time).UTC().Format("Monday, Jan _2 2006, 15:04:05")
	// 	row = sheet.AddRow()
	// 	cell = row.AddCell()
	// 	cell.SetStyle(style)
	// 	cell.Value = data.GetString("id")

	// 	cell = row.AddCell()
	// 	cell.SetStyle(style)
	// 	cell.Value = data.GetString("loginid")

	// 	cell = row.AddCell()
	// 	cell.SetStyle(style)
	// 	cell.Value = created

	// 	cell = row.AddCell()
	// 	cell.SetStyle(style)
	// 	cell.Value = expired

	// }
	ExcelFilename := "Role Metrics Report " + time.Now().Format("20060102150405") + ".xlsx"

	err = file.Save(c.DownloadPath + "/" + ExcelFilename)

	if err != nil {
		tk.Printf(err.Error())
	}

	return c.SetResultInfo(false, "", ExcelFilename)
}

func (c *AclController) GetAuditTrailReport(k *knot.WebContext) interface{} {

	// c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		tk.Printf(err.Error())
	}

	parm := struct {
		EndPeriod   string
		StartPeriod string
	}{}

	err = k.GetPayload(&parm)
	if err != nil {
		return err.Error()
	}

	startperiod, _ := time.Parse("20060102", parm.StartPeriod)
	endperiod, _ := time.Parse("20060102", parm.EndPeriod)
	endperiod = endperiod.AddDate(0, 0, 1)
	result := []tk.M{}

	d := new(UserAuditTrailModel)
	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.And(db.Gte("actiontime", startperiod), db.Lte("actiontime", endperiod))).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	now := time.Now().UTC()
	times := now.Format("01/02/2006")
	rundate := "Run Date : "
	rundates := rundate + times

	//loginid
	users := k.Session("username").(string)
	user := "User : "
	username := user + users

	font := xlsx.NewFont(11, "Calibri")
	style := xlsx.NewStyle()
	style.Font = *font

	fontHdr := xlsx.NewFont(11, "Calibri")
	fontHdr.Bold = true
	styleHdr := xlsx.NewStyle()
	styleHdr.Font = *fontHdr

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Appilcation : BEF CB"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = ""

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = rundates

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = ""

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = username

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = ""

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Date"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Time"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "BankId"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "ID Changed"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Type Of Change"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Field Changed"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Old Value"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "New Value"

	for _, data := range result {
		logintime := data.Get("actiontime").(time.Time)
		dates := logintime.Format("01/02/2006")
		time := logintime.Format("15:04:05")

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = dates

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = time

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("userid")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("useridchanged")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("typeofchange")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("fieldchanged")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("oldvalue")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("newvalue")

		// for _, i := range results {

		// 	logintime, _ := time.Parse(time.RFC3339, i.GetString("logintime"))
		// 	dates := logintime.Format("01/02/2006")

		// 	cell = row.AddCell()
		// 	cell.SetStyle(style)
		// 	cell.Value = dates
		// }

	}
	ExcelFilename := "Audit Trail Report " + time.Now().Format("20060102150405") + ".xlsx"

	err = file.Save(c.DownloadPath + "/" + ExcelFilename)

	if err != nil {
		tk.Printf(err.Error())
	}

	return c.SetResultInfo(false, "", ExcelFilename)
}
