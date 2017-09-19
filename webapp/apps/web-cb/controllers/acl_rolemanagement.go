package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"strings"
	"time"

	"github.com/eaciit/acl/v2.0"
	"github.com/eaciit/dbox"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
	"gopkg.in/mgo.v2/bson"
	// "fmt"
)

func (c *AclController) GetRoleList(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
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

	for _, x := range result {

		resultUser := []tk.M{}
		u := new(acl.User)
		csr, e := c.AclCtx.Connection.NewQuery().From(u.TableName()).Where(db.Eq("groups", x.GetString("id"))).Cursor(nil)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
		e = csr.Fetch(&resultUser, 0, false)
		csr.Close()
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
		if len(resultUser) > 0 {
			x.Set("IsAvailableUser", true)
		} else {
			x.Set("IsAvailableUser", false)
		}
	}
	return c.SetResultInfo(false, "", result)
}

func (c *AclController) GetRoleManagementReferences(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	result := tk.M{}

	result.Set("GroupData", new(acl.Group))
	return c.SetResultInfo(false, "", result)
}

func (c *AclController) GetRole(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	parm := struct {
		Id string
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	result := tk.M{}
	GroupData := new(acl.Group)
	csr, e := c.AclCtx.Connection.NewQuery().From(GroupData.TableName()).Where(db.Eq("id", parm.Id)).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&GroupData, 1, false)
	csr.Close()

	AccessibilityData := []AccessibilityModel{}
	csr, e = c.AclCtx.Connection.NewQuery().From(new(AccessibilityModel).TableName()).Where(db.Eq("roleid", GroupData.ID)).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&AccessibilityData, 0, false)
	csr.Close()

	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	result.Set("GroupData", GroupData)
	result.Set("AccessibilityData", AccessibilityData)
	return c.SetResultInfo(false, "", result)

}

func (c *AclController) SaveRole(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Mode          string
		Data          acl.Group
		Accessibility []AccessibilityModel
	}{}

	err := k.GetPayload(&parm)
	c.Action(k, "Role Management", "Save Role", "", "", "", "", "")

	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	parm.Data.ID = "CB_" + strings.ToUpper(parm.Data.ID)
	data := new(acl.Group)
	if parm.Mode == "create" {
		data.ID = parm.Data.ID
	} else {
		// tk.Println(parm.Data.ID)
		csr, err := c.AclCtx.Connection.NewQuery().From(data.TableName()).Where(db.Eq("id", parm.Data.ID)).Cursor(nil)
		err = csr.Fetch(&data, 1, false)
		csr.Close()
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
	}

	data.ID = parm.Data.ID
	data.Enable = parm.Data.Enable
	data.IsPartOfInitiativeOwner = parm.Data.IsPartOfInitiativeOwner
	data.Title = parm.Data.Title
	data.Grants = parm.Data.Grants
	data.Owner = parm.Data.Owner
	data.GroupType = parm.Data.GroupType
	data.GroupConf = parm.Data.GroupConf
	data.MemberConf = parm.Data.MemberConf

	// Reset Grant Data
	data.Grants = []acl.AccessGrant{}

	AccessData := new(AccessibilityModel)
	err = c.AclCtx.DeleteMany(AccessData, dbox.Eq("roleid", data.ID))
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	// Checking Accessibility
	AccessibilitiesData := map[string]bool{}
	IDData := tk.M{}
	for _, i := range parm.Accessibility {
		i.ID = bson.NewObjectId()
		IDData.Set(i.AccessID, i.ID)
		i.RoleID = data.ID
		if i.Global.Create == true {
			i.AllowStatus = true
		}
		if i.Global.Read == true {
			i.AllowStatus = true
		}
		if i.Global.Update == true {
			i.AllowStatus = true
		}
		if i.Global.Delete == true {
			i.AllowStatus = true
		}
		if i.Global.Owned == true {
			i.AllowStatus = true
		}
		if i.Global.Curtain == true {
			i.AllowStatus = true
		}
		if i.Global.Upload == true {
			i.AllowStatus = true
		}
		if i.Region.Create == true {
			i.AllowStatus = true
		}
		if i.Region.Read == true {
			i.AllowStatus = true
		}
		if i.Region.Update == true {
			i.AllowStatus = true
		}
		if i.Region.Delete == true {
			i.AllowStatus = true
		}
		if i.Region.Owned == true {
			i.AllowStatus = true
		}
		if i.Region.Curtain == true {
			i.AllowStatus = true
		}
		if i.Region.Upload == true {
			i.AllowStatus = true
		}
		if i.Country.Create == true {
			i.AllowStatus = true
		}
		if i.Country.Read == true {
			i.AllowStatus = true
		}
		if i.Country.Update == true {
			i.AllowStatus = true
		}
		if i.Country.Delete == true {
			i.AllowStatus = true
		}
		if i.Country.Owned == true {
			i.AllowStatus = true
		}
		if i.Country.Curtain == true {
			i.AllowStatus = true
		}
		if i.Country.Upload == true {
			i.AllowStatus = true
		}

		if i.Url == "#" && i.AllowStatus && !AccessibilitiesData[i.ParentID] {
			AccessibilitiesData[i.ParentID] = true
		}
		err = c.AclCtx.Save(&i)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}

		access := []acl.AccessTypeEnum{}
		if i.AllowStatus == true || i.AccessID == "LOGOUT" {
			access = append(access, 15)
		}
		if len(access) > 0 {
			data.Grant(i.AccessID, access...)
			if i.ParentID != "" {
				data.Grant(i.ParentID, access...)
			}
		}
	}

	for x, _ := range AccessibilitiesData {
		for _, i := range parm.Accessibility {
			if i.AccessID == x {
				i.ID = IDData[i.AccessID].(bson.ObjectId)
				i.AllowStatus = true
				i.RoleID = data.ID
				err = c.AclCtx.Save(&i)
				if err != nil {
					return c.ErrorResultInfo(err.Error(), nil)
				}
				break
			}
		}
	}

	data.Applications = []string{"web-cb"}
	err = acl.Save(data)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	return c.SetResultInfo(false, "", nil)
}

func (c *AclController) RemoveRole(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	c.Action(k, "Role Management", "Remove Role", "", "", "", "", "")

	k.Config.OutputType = knot.OutputJson
	// data := tk.M{}
	parm := struct {
		Id string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	data := new(acl.Group)

	csr, err := c.AclCtx.Connection.NewQuery().From(data.TableName()).Where(db.Eq("_id", parm.Id)).Cursor(nil)
	err = csr.Fetch(&data, 1, false)
	csr.Close()
	if err != nil {
		return nil
	}

	err = acl.Delete(data)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	return c.SetResultInfo(false, "", data)
}

func (c *AclController) ExportXLSRoleManagemet(k *knot.WebContext) interface{} {
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

	fontHdr := xlsx.NewFont(11, "Calibri")
	fontHdr.Bold = true
	styleHdr := xlsx.NewStyle()
	styleHdr.Font = *fontHdr

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
	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = "Role"

	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = "Role Description"

	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = "Permissions"

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
		cell.Value = strings.Replace(data.GetString("id"), "CB_", "", -1)

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
				granttext += data2.AccessID
				granttext += " : "
				granttext += granttext2
				granttext += "\n"
			}

		}
		// datas := data.Get("grants").([]interface{})

		// for _,data2 := range datas{
		// 	// fmt.Println(data2)
		// 	dat,e := tk.ToM(data2)
		// 	if e != nil {
		// 		return c.ErrorResultInfo(e.Error(), nil)
		// 	}
		// 	fmt.Println(dat,dat.GetString("accessid") )
		// 	granttext += dat.GetString("accessid")
		// 	granttext += " : "
		// 	granttext += "\n"
		// }

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = granttext
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
