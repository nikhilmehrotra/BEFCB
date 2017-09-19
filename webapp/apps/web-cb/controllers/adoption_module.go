package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type AdoptionModuleController struct {
	*BaseController
}

func (c *AdoptionModuleController) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	c.Action(k, "Adoption Module", "Open Adoption Modul Page", "", "", "", "", "")
	AdoptionModule := c.GetAccess(k, "ADOPTIONMODULE")
	Scorecard := c.GetAccess(k, "SCORECARD")
	Initiative := c.GetAccess(k, "INITIATIVE")
	k.Config.NoLog = true
	k.Config.LayoutTemplate = "_layout-v2.html"
	PartialFiles := []string{}
	PartialFiles = append(PartialFiles, "shared/sidebar.html")
	PartialFiles = append(PartialFiles, "shared/initiative_summary.html")
	PartialFiles = append(PartialFiles, "shared/initiative_chart.html")
	PartialFiles = append(PartialFiles, "shared/metric_upload.html")
	PartialFiles = append(PartialFiles, "dashboard/initiative.html")

	PartialFiles = append(PartialFiles, "adoptionmodule/detail.html")
	PartialFiles = append(PartialFiles, "adoptionmodule/analysis.html")

	k.Config.IncludeFiles = PartialFiles
	k.Config.OutputType = knot.OutputTemplate

	UserCountry := ""
	if k.Session("country") != nil {
		UserCountry = k.Session("country").(string)
	}
	return tk.M{}.Set("AdoptionModule", AdoptionModule).Set("UserCountry", UserCountry).Set("Scorecard", Scorecard).Set("Initiative", Initiative)
}

func (c *AdoptionModuleController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := []tk.M{}
	parm := struct {
		Geography string
		Country   string
	}{}
	AdoptionModuleAccess := c.GetAccess(k, "ADOPTIONMODULE")
	AdoptionModule := new(AccessibilityModel)
	e := tk.MtoStruct(AdoptionModuleAccess, &AdoptionModule)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	InitiativeAccess := c.GetAccess(k, "INITIATIVE")
	Initiative := new(AccessibilityModel)
	e = tk.MtoStruct(InitiativeAccess, &Initiative)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	err := k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	if !AdoptionModule.Global.Read && !AdoptionModule.Region.Read && !AdoptionModule.Country.Read && !(AdoptionModule.Global.Owned || AdoptionModule.Region.Owned || AdoptionModule.Country.Owned) {
		return c.SetResultInfo(false, "", res)
	}

	UserCountry := ""
	if k.Session("country") != nil {
		UserCountry = k.Session("country").(string)
	}

	query := []*db.Filter{}
	query = append(query, db.Eq("InitiativeType", "KeyEnablers"))
	if !Initiative.Global.Read || !Initiative.Global.Read || !Initiative.Global.Read {
		geographyQuery := []*db.Filter{}
		if Initiative.Global.Read {
			geographyQuery = append(geographyQuery, db.Eq("IsGlobal", true))
		}
		if Initiative.Region.Read {
			geographyQuery = append(geographyQuery, db.Ne("Region", []string{}))
		}
		if Initiative.Country.Read {
			geographyQuery = append(geographyQuery, db.Ne("Country", []string{}))
		}
		if len(geographyQuery) > 0 {
			query = append(query, db.Or(geographyQuery...))
		}
	}
	if !(Initiative.Global.Read || Initiative.Region.Read || Initiative.Country.Read) && (Initiative.Global.Owned || Initiative.Region.Owned || Initiative.Country.Owned) {
		fullname := ""
		if k.Session("fullname") != nil {
			fullname = k.Session("fullname").(string)
		}
		query = append(query, db.Or(db.Eq("ProjectManager", fullname), db.Eq("AccountableExecutive", fullname), db.Eq("TechnologyLead", fullname)))
	}

	if UserCountry != "" && (Initiative.Global.Read || Initiative.Region.Read || Initiative.Country.Read) == true {
		query = append(query, db.Or(db.Eq("Country", UserCountry), db.Eq("IsGlobal", true)))
	} else if UserCountry != "" && (Initiative.Global.Read || Initiative.Region.Read || Initiative.Country.Read) == false {
		query = append(query, db.Or(db.Eq("Country", UserCountry), db.Eq("IsGlobal", false)))
	}

	csr, e := c.Ctx.Connection.NewQuery().Where(db.And(query...)).From("Initiative").Select("InitiativeID").Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	InitiativeList := []tk.M{}
	e = csr.Fetch(&InitiativeList, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	InitiativesData := []string{}
	for _, i := range InitiativeList {
		o := i.Get("InitiativeID").(string)
		InitiativesData = append(InitiativesData, o)
	}

	if AdoptionModule.Global.Owned || AdoptionModule.Region.Owned || AdoptionModule.Country.Owned {
		fullname := ""
		if k.Session("fullname") != nil {
			fullname = k.Session("fullname").(string)
		}
		query = append(query, db.Or(db.Eq("ProjectManager", fullname), db.Eq("AccountableExecutive", fullname), db.Eq("TechnologyLead", fullname)))
	}
	csr, e = c.Ctx.Connection.NewQuery().Where(db.And(query...)).From("Initiative").Select("InitiativeID").Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	OwnedInitiatives := []tk.M{}
	e = csr.Fetch(&OwnedInitiatives, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	OwnedInitiativesData := []string{}
	OwnedInitiativesDataID := []string{}
	for _, i := range OwnedInitiatives {
		oid := i.Get("_id").(bson.ObjectId).Hex()
		o := i.Get("InitiativeID").(string)
		OwnedInitiativesDataID = append(OwnedInitiativesDataID, oid)
		OwnedInitiativesData = append(OwnedInitiativesData, o)
	}

	result := make([]tk.M, 0)
	pipes, match, group, sort := []tk.M{}, []tk.M{}, tk.M{}, tk.M{}
	if parm.Geography != "ALL" {
		switch parm.Geography {
		case "GLOBAL":
			match = append(match, tk.M{"$or": []tk.M{
				tk.M{
					"IsGlobal": tk.M{"$eq": true},
				},
				tk.M{"$and": []tk.M{
					tk.M{
						"IsGlobal": tk.M{"$eq": false},
					},
					tk.M{
						"Region": tk.M{"$eq": []string{}},
					},
					tk.M{
						"Country": tk.M{"$eq": []string{}},
					},
				}},
			}})
			break
		// case "REGION":
		// 	match = append(match, tk.M{"Region": tk.M{"$ne": []string{}}})
		// 	break
		// case "COUNTRY":
		// 	match = append(match, tk.M{"Country": tk.M{"$ne": []string{}}})
		// 	break
		case "COUNTRY":
			match = append(match, tk.M{"$or": []tk.M{
				tk.M{
					"Region": tk.M{"$ne": []string{}},
				},
				tk.M{
					"Country": tk.M{"$ne": []string{}},
				},
			}})
			break
			break
		default:
			break
		}
	}
	err = Deserialize(`
		{"$group":{
	        "_id":{
	            "InitiativeID":"$InitiativeID",
	            "Initiatives":"$ProjectName",
	            "Description":"$ProblemStatement",
	            "Benefit":"$ProjectDescription",
	            "GoLive":"$FinishDate",
	            "RAGBrief":"$MetricBenchmark",
	            "UsefulResources":"$UsefulResources",
	            "IsGlobal":"$IsGlobal",
	            "Region":"$Region",
	            "Country":"$Country"
	        }
	    }}
	`, &group)

	err = Deserialize(`
		{"$sort":{"_id.Initiatives":1}}
	`, &sort)

	pipes = append(pipes, tk.M{"$match": tk.M{"IsInitiativeTracked": true}})
	if AdoptionModule.Global.Read || AdoptionModule.Region.Read || AdoptionModule.Country.Read {
		pipes = append(pipes, tk.M{"$match": tk.M{"InitiativeID": tk.M{"$in": InitiativesData}}})
	} else {
		pipes = append(pipes, tk.M{"$match": tk.M{"InitiativeID": tk.M{"$in": OwnedInitiativesData}}})
	}

	if parm.Country != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"Country": parm.Country}})
	}

	if parm.Geography != "ALL" {
		pipes = append(pipes, tk.M{"$match": tk.M{"$and": match}})
		if parm.Country != "" {
			unwind, unwind_country := tk.M{}, tk.M{}
			err = Deserialize(`
				{"$unwind":"$Milestones"}
			`, &unwind)
			pipes = append(pipes, unwind)
			err = Deserialize(`
				{"$unwind":"$Milestones.Country"}
			`, &unwind_country)
			pipes = append(pipes, unwind_country)
			pipes = append(pipes, tk.M{"$match": tk.M{"Milestones.Country": parm.Country}})

			err = Deserialize(`
				{"$group":{
			         "_id":{
			            "InitiativeID":"$InitiativeID",
			            "Initiatives":"$ProjectName",
			            "Description":"$ProblemStatement",
			            "Benefit":"$ProjectDescription",
			            "GoLive":"$Milestones.EndDate",
			            "RAGBrief":"$MetricBenchmark",
			            "UsefulResources":"$UsefulResources",
			            "IsGlobal":"$IsGlobal",
			            "Region":"$Region",
			            "Country":"$Milestones.Country"
			        }
			    }}
			`, &group)
		}
	}
	pipes = append(pipes, group)
	pipes = append(pipes, sort)

	csr, err = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From(new(InitiativeModel).TableName()).Cursor(nil)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	err = csr.Fetch(&result, 0, false)
	csr.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	// Get Metrics Data
	MetricData := []tk.M{}
	pipes, match, group, sort = []tk.M{}, []tk.M{}, tk.M{}, tk.M{}
	if parm.Country != "" {
		match = append(match, tk.M{"country": parm.Country})
	} else {
		match = append(match, tk.M{"country": "TOTAL"})
	}

	match = append(match, tk.M{"metricid": tk.M{"$ne": ""}})
	err = Deserialize(`
		{"$group":{"_id":{
	        "InitiativeId" : "$initiativeid",
	        "MetricId":"$metricid",
	        "MetricName":"$metricname",
	        "LatestRAG":"$latestrag"
	    }}}
	`, &group)
	err = Deserialize(`
    	{"$sort":{
    		"_id.InitiativeId":1,
    		"_id.MetricId":1
    	}}
	`, &sort)
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": match}})
	pipes = append(pipes, group)
	pipes = append(pipes, sort)
	csr, err = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From(new(MilestoneValueModel).TableName()).Cursor(nil)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	err = csr.Fetch(&MetricData, 0, false)
	csr.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	MData := map[string]bool{}
	MDataList := map[string][]tk.M{}
	for _, x := range MetricData {
		InitaitiveId := x.Get("_id").(tk.M).GetString("InitiativeId")
		MetricId := x.Get("_id").(tk.M).GetString("MetricId")
		MetricName := x.Get("_id").(tk.M).GetString("MetricName")
		tmp := MDataList[InitaitiveId]
		current_list := []tk.M{}
		if tmp != nil {
			current_list = tmp
		}
		RAG := x.Get("_id").(tk.M).GetString("LatestRAG")
		data := tk.M{"MetricId": MetricId, "MetricName": MetricName, "RAG": RAG}
		current_list = append(current_list, data)
		MDataList[InitaitiveId] = current_list
		MData[InitaitiveId] = true
	}

	for _, i := range result {
		d := i.Get("_id").(tk.M)
		d.Set("Adoption", "")
		IsHavingMetricData, _ := MData[d.GetString("InitiativeID")]
		d.Set("IsHavingMetricData", IsHavingMetricData)

		Metric, _ := MDataList[d.GetString("InitiativeID")]
		if len(Metric) == 0 {
			Metric = append(Metric, tk.M{"MetricId": "", "MetricName": "", "RAG": ""})
		}
		d.Set("MetricData", Metric)
		res = append(res, d)
	}
	return c.SetResultInfo(false, "", tk.M{"Sources": res, "OwnedInitiative": OwnedInitiativesData})
}

func (c *AdoptionModuleController) SaveExcel(k *knot.WebContext) interface{} {
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
	parm := struct {
		Geography string
		Country   string
	}{}
	err = k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	//data
	result := make([]tk.M, 0)
	pipes, match, group, sort := []tk.M{}, []tk.M{}, tk.M{}, tk.M{}
	if parm.Geography != "ALL" {
		switch parm.Geography {
		case "GLOBAL":
			match = append(match, tk.M{"$or": []tk.M{
				tk.M{
					"IsGlobal": tk.M{"$eq": true},
				},
				tk.M{"$and": []tk.M{
					tk.M{
						"IsGlobal": tk.M{"$eq": false},
					},
					tk.M{
						"Region": tk.M{"$eq": []string{}},
					},
					tk.M{
						"Country": tk.M{"$eq": []string{}},
					},
				}},
			}})
			break
		case "COUNTRY":
			match = append(match, tk.M{"$or": []tk.M{
				tk.M{
					"Region": tk.M{"$ne": []string{}},
				},
				tk.M{
					"Country": tk.M{"$ne": []string{}},
				},
			}})
			break
			break
		default:
			break
		}
	}
	err = Deserialize(`
		{"$group":{
	        "_id":{
	            "InitiativeID":"$InitiativeID",
	            "Initiatives":"$ProjectName",
	            "Description":"$ProblemStatement",
	            "Benefit":"$ProjectDescription",
	            "GoLive":"$FinishDate",
	            "RAGBrief":"$MetricBenchmark",
	            "UsefulResources":"$UsefulResources",
	            "IsGlobal":"$IsGlobal",
	            "Region":"$Region",
	            "Country":"$Country"
	        }
	    }}
	`, &group)

	err = Deserialize(`
		{"$sort":{"_id.Initiatives":1}}
	`, &sort)

	pipes = append(pipes, tk.M{"$match": tk.M{"IsInitiativeTracked": true}})
	if parm.Country != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"Country": parm.Country}})
	}

	if parm.Geography != "ALL" {
		pipes = append(pipes, tk.M{"$match": tk.M{"$and": match}})
		if parm.Country != "" {
			unwind, unwind_country := tk.M{}, tk.M{}
			err = Deserialize(`
				{"$unwind":"$Milestones"}
			`, &unwind)
			pipes = append(pipes, unwind)
			err = Deserialize(`
				{"$unwind":"$Milestones.Country"}
			`, &unwind_country)
			pipes = append(pipes, unwind_country)
			pipes = append(pipes, tk.M{"$match": tk.M{"Milestones.Country": parm.Country}})

			err = Deserialize(`
				{"$group":{
			         "_id":{
			            "InitiativeID":"$InitiativeID",
			            "Initiatives":"$ProjectName",
			            "Description":"$ProblemStatement",
			            "Benefit":"$ProjectDescription",
			            "GoLive":"$Milestones.EndDate",
			            "RAGBrief":"$MetricBenchmark",
			            "UsefulResources":"$UsefulResources",
			            "IsGlobal":"$IsGlobal",
			            "Region":"$Region",
			            "Country":"$Milestones.Country"
			        }
			    }}
			`, &group)
		}
	}
	pipes = append(pipes, group)
	pipes = append(pipes, sort)

	csr, err := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From(new(InitiativeModel).TableName()).Cursor(nil)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	err = csr.Fetch(&result, 0, false)
	csr.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	res := []tk.M{}
	for _, i := range result {
		d := i.Get("_id").(tk.M)
		d.Set("Adoption", "")
		res = append(res, d)
	}

	//xls
	font := xlsx.NewFont(11, "Calibri")
	style := xlsx.NewStyle()
	style.Font = *font

	fontHdr := xlsx.NewFont(11, "Calibri")
	fontHdr.Bold = true
	styleHdr := xlsx.NewStyle()
	styleHdr.Font = *fontHdr
	// tk.Println("parm.Country", parm.Geography)
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Adoption Module - " + parm.Geography

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = ""

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Initiatives"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Problem Statement"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Project Description"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Geography"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Go-Live Month"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "RAG Brief"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "UsefullResources"

	for _, data := range res {
		global := data.Get("IsGlobal").(bool)
		geo := ""
		if global == true {
			geo = "Global"
		} else {
			geo = ""
		}

		t := data.Get("GoLive").(time.Time)
		tFormat := t.Format("Jan 2006")

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("Initiatives")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("Description")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("Benefit")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = geo

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = tFormat

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("RAGBrief")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("UsefulResources")

	}
	ExcelFilename := "Adoption Module " + time.Now().Format("20060102150405") + ".xlsx"

	err = file.Save(c.DownloadPath + "/" + ExcelFilename)

	if err != nil {
		tk.Printf(err.Error())
	}

	return c.SetResultInfo(false, "", ExcelFilename)
}
