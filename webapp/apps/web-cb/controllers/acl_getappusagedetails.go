package controllers

import (
	// helper "eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	// "github.com/eaciit/acl/v2.0"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
	// "strings"
	"time"
)

func (c *AclController) GetAppUsageDetails(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	// Get PayLoad
	type SubFilters struct {
		Field    string
		Operator string
		Value    string
		Type     string
	}
	type Filters struct {
		Field    string
		Operator string
		Value    string
		Type     string
		Filters  []SubFilters
	}
	type Filter struct {
		Filters []Filters
		Logic   string
	}
	type Sort struct {
		Field string
		Dir   string
	}
	parm := struct {
		Skip     int
		Take     int
		Page     int
		PageSize int
		Filter   Filter
		Sort     []Sort

		// Priority Filter
		StartDate  string
		FinishDate string
		LoginID    []string
		FullName   []string
		Country    []string
		UserRoles  []string
		Action     []string
		Modules    []string
	}{}
	e := k.GetPayload(&parm)

	result := []tk.M{}

	query := []*db.Filter{}
	sort := []string{}

	if len(parm.LoginID) > 0 {
		loginId := []interface{}{}
		for _, x := range parm.LoginID {
			eachLoginId := x
			loginId = append(loginId, eachLoginId)
		}
		query = append(query, db.In("userid", loginId...))
	}
	if len(parm.FullName) > 0 {
		fullName := []interface{}{}
		for _, x := range parm.FullName {
			eachFullName := x
			fullName = append(fullName, eachFullName)
		}
		query = append(query, db.In("fullname", fullName...))
	}
	if len(parm.Country) > 0 {
		country := []interface{}{}
		for _, x := range parm.Country {
			eachCountry := x
			country = append(country, eachCountry)
		}
		query = append(query, db.In("country", country...))
	}
	if len(parm.UserRoles) > 0 {
		userRoles := []interface{}{}
		for _, x := range parm.UserRoles {
			eachUserRoles := x
			userRoles = append(userRoles, eachUserRoles)
		}
		query = append(query, db.In("group", userRoles...))
	}
	if len(parm.Action) > 0 {
		action := []interface{}{}
		for _, x := range parm.Action {
			eachAction := x
			action = append(action, eachAction)
		}
		query = append(query, db.In("do", action...))
	}
	// tk.Println(parm.Modules)
	if len(parm.Modules) > 0 {
		modules := []interface{}{}
		for _, x := range parm.Modules {
			eachModules := x
			modules = append(modules, eachModules)
		}
		query = append(query, db.In("module", modules...))
	}

	if parm.StartDate != "" {
		StartDate, _ := time.Parse("20060102150405", parm.StartDate)
		query = append(query, db.Gte("dateaccess", StartDate))
	}
	if parm.FinishDate != "" {
		FinishDate, _ := time.Parse("20060102150405", parm.FinishDate)
		query = append(query, db.Lte("expiredtime", FinishDate.AddDate(0, 0, 1)))
	}

	queryFilter := []*db.Filter{}

	if len(parm.Filter.Filters) > 0 {
		for _, f := range parm.Filter.Filters {

			if len(f.Filters) == 0 {
				if f.Type == "date" {
					tmpdate, _ := time.Parse("20060102150405", f.Value)
					tmpdate = tmpdate.UTC()
					value := tmpdate
					switch f.Operator {
					case "eq":
						queryFilter = append(queryFilter, db.Eq(f.Field, value))
					case "gte":
						queryFilter = append(queryFilter, db.Gte(f.Field, value))
					case "lte":
						queryFilter = append(queryFilter, db.Lte(f.Field, value))
					default:
						break
					}
				} else {
					value := f.Value
					switch f.Operator {
					case "contains":
						queryFilter = append(queryFilter, db.Contains(f.Field, value))
					case "eq":
						queryFilter = append(queryFilter, db.Eq(f.Field, value))
					case "startswith":
						queryFilter = append(queryFilter, db.Startwith(f.Field, value))
					case "endswith":
						queryFilter = append(queryFilter, db.Endwith(f.Field, value))
					default:
						break
					}

				}
				if len(queryFilter) > 0 {
					if parm.Filter.Logic == "and" {
						query = append(query, db.And(queryFilter...))
					} else {
						query = append(query, db.Or(queryFilter...))
					}
				}

			} else {
				// tk.Println("MASUK MULTIPLE")
				MultipleFilter := []*db.Filter{}
				for _, x := range parm.Filter.Filters {
					// tk.Println("Filters", x.Filters)
					for _, y := range x.Filters {
						// tk.Println("ObjFilter", y)
						if y.Type == "date" {
							tmpdate, _ := time.Parse("20060102150405", y.Value)
							tmpdate = tmpdate.UTC()
							value := tmpdate
							switch y.Operator {
							case "eq":
								MultipleFilter = append(MultipleFilter, db.Eq(y.Field, value))
							case "gte":
								MultipleFilter = append(MultipleFilter, db.Gte(y.Field, value))
							case "lte":
								MultipleFilter = append(MultipleFilter, db.Lte(y.Field, value))
							default:
								break
							}
						} else {
							value := y.Value
							switch y.Operator {
							case "contains":
								MultipleFilter = append(MultipleFilter, db.Contains(y.Field, value))
							case "eq":
								MultipleFilter = append(MultipleFilter, db.Eq(y.Field, value))
							case "startswith":
								MultipleFilter = append(MultipleFilter, db.Startwith(y.Field, value))
							case "endswith":
								MultipleFilter = append(MultipleFilter, db.Endwith(y.Field, value))
							default:
								break
							}

						}
					}

				}
				if len(MultipleFilter) > 0 {
					if parm.Filter.Logic == "and" {
						query = append(query, db.And(MultipleFilter...))
					} else {
						query = append(query, db.Or(MultipleFilter...))
					}
				}
			}
		}
	}

	// Get Sort Value
	if len(parm.Sort) > 0 {
		for _, x := range parm.Sort {
			if x.Dir == "asc" {
				sort = append(sort, x.Field)
			} else {
				sort = append(sort, "-"+x.Field)
			}
		}
	} else {
		sort = append(sort, "-dateaccess")
	}

	d := NewLogModel()

	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(query...).Skip(parm.Skip).Take(parm.Take).Order(sort...).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	returnvalue := struct {
		Data  []tk.M
		Total int
	}{
		Data:  result,
		Total: csr.Count(),
	}

	return c.SetResultInfo(false, "", returnvalue)
}

func (c *AclController) GetAppUsageDetailsXLS(k *knot.WebContext) interface{} {
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

	// Get PayLoad
	type SubFilters struct {
		Field    string
		Operator string
		Value    string
		Type     string
	}
	type Filters struct {
		Field    string
		Operator string
		Value    string
		Type     string
		Filters  []SubFilters
	}
	type Filter struct {
		Filters []Filters
		Logic   string
	}
	type Sort struct {
		Field string
		Dir   string
	}
	parm := struct {
		Filter Filter
		Sort   []Sort

		// Priority Filter
		StartDate  string
		FinishDate string
		LoginID    []string
		FullName   []string
		Country    []string
		UserRoles  []string
		Action     []string
		Modules    []string
	}{}
	e := k.GetPayload(&parm)

	result := []tk.M{}

	query := []*db.Filter{}
	sort := []string{}

	if len(parm.LoginID) > 0 {
		loginId := []interface{}{}
		for _, x := range parm.LoginID {
			eachLoginId := x
			loginId = append(loginId, eachLoginId)
		}
		query = append(query, db.In("userid", loginId...))
	}
	if len(parm.FullName) > 0 {
		fullName := []interface{}{}
		for _, x := range parm.FullName {
			eachFullName := x
			fullName = append(fullName, eachFullName)
		}
		query = append(query, db.In("fullname", fullName...))
	}
	if len(parm.Country) > 0 {
		country := []interface{}{}
		for _, x := range parm.Country {
			eachCountry := x
			country = append(country, eachCountry)
		}
		query = append(query, db.In("country", country...))
	}
	if len(parm.UserRoles) > 0 {
		userRoles := []interface{}{}
		for _, x := range parm.UserRoles {
			eachUserRoles := x
			userRoles = append(userRoles, eachUserRoles)
		}
		query = append(query, db.In("group", userRoles...))
	}
	if len(parm.Action) > 0 {
		action := []interface{}{}
		for _, x := range parm.Action {
			eachAction := x
			action = append(action, eachAction)
		}
		query = append(query, db.In("do", action...))
	}
	// tk.Println(parm.Modules)
	if len(parm.Modules) > 0 {
		modules := []interface{}{}
		for _, x := range parm.Modules {
			eachModules := x
			modules = append(modules, eachModules)
		}
		query = append(query, db.In("module", modules...))
	}

	if parm.StartDate != "" {
		StartDate, _ := time.Parse("20060102150405", parm.StartDate)
		query = append(query, db.Gte("dateaccess", StartDate))
	}
	if parm.FinishDate != "" {
		FinishDate, _ := time.Parse("20060102150405", parm.FinishDate)
		query = append(query, db.Lte("expiredtime", FinishDate.AddDate(0, 0, 1)))
	}

	queryFilter := []*db.Filter{}

	if len(parm.Filter.Filters) > 0 {
		for _, f := range parm.Filter.Filters {

			if len(f.Filters) == 0 {
				if f.Type == "date" {
					tmpdate, _ := time.Parse("20060102150405", f.Value)
					tmpdate = tmpdate.UTC()
					value := tmpdate
					switch f.Operator {
					case "eq":
						queryFilter = append(queryFilter, db.Eq(f.Field, value))
					case "gte":
						queryFilter = append(queryFilter, db.Gte(f.Field, value))
					case "lte":
						queryFilter = append(queryFilter, db.Lte(f.Field, value))
					default:
						break
					}
				} else {
					value := f.Value
					switch f.Operator {
					case "contains":
						queryFilter = append(queryFilter, db.Contains(f.Field, value))
					case "eq":
						queryFilter = append(queryFilter, db.Eq(f.Field, value))
					case "startswith":
						queryFilter = append(queryFilter, db.Startwith(f.Field, value))
					case "endswith":
						queryFilter = append(queryFilter, db.Endwith(f.Field, value))
					default:
						break
					}

				}
				if len(queryFilter) > 0 {
					if parm.Filter.Logic == "and" {
						query = append(query, db.And(queryFilter...))
					} else {
						query = append(query, db.Or(queryFilter...))
					}
				}

			} else {
				// tk.Println("MASUK MULTIPLE")
				MultipleFilter := []*db.Filter{}
				for _, x := range parm.Filter.Filters {
					// tk.Println("Filters", x.Filters)
					for _, y := range x.Filters {
						// tk.Println("ObjFilter", y)
						if y.Type == "date" {
							tmpdate, _ := time.Parse("20060102150405", y.Value)
							tmpdate = tmpdate.UTC()
							value := tmpdate
							switch y.Operator {
							case "eq":
								MultipleFilter = append(MultipleFilter, db.Eq(y.Field, value))
							case "gte":
								MultipleFilter = append(MultipleFilter, db.Gte(y.Field, value))
							case "lte":
								MultipleFilter = append(MultipleFilter, db.Lte(y.Field, value))
							default:
								break
							}
						} else {
							value := y.Value
							switch y.Operator {
							case "contains":
								MultipleFilter = append(MultipleFilter, db.Contains(y.Field, value))
							case "eq":
								MultipleFilter = append(MultipleFilter, db.Eq(y.Field, value))
							case "startswith":
								MultipleFilter = append(MultipleFilter, db.Startwith(y.Field, value))
							case "endswith":
								MultipleFilter = append(MultipleFilter, db.Endwith(y.Field, value))
							default:
								break
							}

						}
					}

				}
				if len(MultipleFilter) > 0 {
					if parm.Filter.Logic == "and" {
						query = append(query, db.And(MultipleFilter...))
					} else {
						query = append(query, db.Or(MultipleFilter...))
					}
				}
			}
		}
	}

	// Get Sort Value
	if len(parm.Sort) > 0 {
		for _, x := range parm.Sort {
			if x.Dir == "asc" {
				sort = append(sort, x.Field)
			} else {
				sort = append(sort, "-"+x.Field)
			}
		}
	} else {
		sort = append(sort, "-dateaccess")
	}

	d := NewLogModel()

	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(query...).Order(sort...).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	//xls
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
	cell.Value = "Login ID"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Full Name"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Country"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "User Role(s)"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Action"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "What Changed"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Old Value"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "New Value"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Action Time"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Resource Url"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Sources"

	for _, data := range result {
		dateaccess := data.Get("dateaccess").(time.Time).UTC().Format("Monday, Jan _2 2006, 15:04:05")

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("userid")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("fullname")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("country")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("groupdescription")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("do")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("whatchanged")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("oldvalue")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("newvalue")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = dateaccess

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("requesturi")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("sources")
		// sources := data.GetString("sources")
		// sources = "ainur_kece.xlsx"
		// url := "www.google.com"
		// tk.Println(k.Request.URL)
		// cell.SetFormula("HYPERLINK(\"" + url + "\", \"" + sources + "\")")

	}
	ExcelFilename := "Application Usage Details " + time.Now().Format("20060102150405") + ".xlsx"

	err = file.Save(c.DownloadPath + "/" + ExcelFilename)

	if err != nil {
		tk.Printf(err.Error())
	}

	return c.SetResultInfo(false, "", ExcelFilename)
}
