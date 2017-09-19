package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"strings"
	"time"
)

func (c *AdoptionModuleController) GetDataDetail(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		InitiativeId string
		Year         int
		MetricId     string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	// Get Initiative Country Data: InitiativeId
	InitiativeData := InitiativeModel{}
	csr, err := c.Ctx.Connection.NewQuery().From(new(InitiativeModel).TableName()).Where(db.Eq("InitiativeID", parm.InitiativeId)).Take(1).Cursor(nil)
	err = csr.Fetch(&InitiativeData, 1, false)
	csr.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	InitiativeCountry := []string{}
	for _, x := range InitiativeData.Milestones {
		for _, country := range x.Country {
			if country != "" {
				if !stringInSlice(country, InitiativeCountry) {

					InitiativeCountry = append(InitiativeCountry, country)

				}
			}
		}

	}
	countryTotal := "TOTAL"
	InitiativeCountry = append(InitiativeCountry, countryTotal)

	// Ge Region Data
	GetCountryCode := map[string]string{}
	RegionData := []RegionModel{}
	csr, err = c.Ctx.Connection.NewQuery().From(new(RegionModel).TableName()).Cursor(nil)
	err = csr.Fetch(&RegionData, 0, false)
	csr.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	for _, x := range RegionData {
		GetCountryCode[x.Country] = x.CountryCode
	}
	// Get Metrics Data
	MetricData := []tk.M{}

	pipes, match, group, project, sort := []tk.M{}, []tk.M{}, tk.M{}, tk.M{}, tk.M{}
	match = append(match, tk.M{"initiativeid": parm.InitiativeId})
	match = append(match, tk.M{"year": parm.Year})
	match = append(match, tk.M{"metricid": tk.M{"$ne": ""}})
	err = Deserialize(`
		{"$group":{"_id":{
	        "MetricId":"$metricid",
	        "MetricName":"$metricname",
	        "Description":"$description",
	        "Denomination":"$denomination",
	        "ActualValue":"$actualvalue",
	        "TotalValue":"$totalvalue",
	        "RAGValue":"$ragvalue"
	    }}}
	`, &group)
	err = Deserialize(`
		{"$project":{
	        "MetricId":"$_id.MetricId",
	        "MetricName":"$_id.MetricName",
	        "Description":"$_id.Description",
	        "Denomination":"$_id.Denomination",
	        "ActualValue":"$_id.ActualValue",
	        "TotalValue":"$_id.TotalValue",
	        "RAGValue":"$_id.RAGValue"
	    }}
	`, &project)
	err = Deserialize(`
		{"$sort":{"MetricName":1}}
	`, &sort)
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": match}})
	pipes = append(pipes, group)
	pipes = append(pipes, project)
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

	for _, x := range MetricData {
		MetricId := x.GetString("MetricId")
		// Get Detail Data
		query := []*db.Filter{}
		query = append(query, db.Eq("initiativeid", parm.InitiativeId))
		query = append(query, db.Eq("year", parm.Year))
		query = append(query, db.Eq("metricid", MetricId))

		DetailData := []MilestoneValueModel{}
		csr, err = c.Ctx.Connection.NewQuery().From(new(MilestoneValueModel).TableName()).Where(db.And(query...)).Cursor(nil)
		err = csr.Fetch(&DetailData, 0, false)
		csr.Close()
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}

		// Mapping Existing Initiative Country Data with Existing Data
		CurrentData := []MilestoneValueModel{}
		for _, x := range InitiativeCountry {
			IsAvailable := false
			data := MilestoneValueModel{}
			for _, d := range DetailData {
				if d.Country == x {
					IsAvailable = true
					data = d
					break
				}
			}
			if !IsAvailable {
				data.InitiativeId = parm.InitiativeId
				data.CountryCode = GetCountryCode[x]
				data.Country = x
				data.Year = parm.Year
				data.NAJan = true
				data.NAFeb = true
				data.NAMar = true
				data.NAApr = true
				data.NAMay = true
				data.NAJun = true
				data.NAJul = true
				data.NAAug = true
				data.NASep = true
				data.NAOct = true
				data.NANov = true
				data.NADec = true
			}
			CurrentData = append(CurrentData, data)

		}

		DetailData = CurrentData
		x.Set("DetailData", DetailData)

	}

	return c.SetResultInfo(false, "", MetricData)
}

func (c *AdoptionModuleController) GetDataDetailOld(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		InitiativeId string
		Year         int
		MetricId     string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	// Get Initiative Country Data: InitiativeId
	InitiativeData := InitiativeModel{}
	csr, err := c.Ctx.Connection.NewQuery().From(new(InitiativeModel).TableName()).Where(db.Eq("InitiativeID", parm.InitiativeId)).Take(1).Cursor(nil)
	err = csr.Fetch(&InitiativeData, 1, false)
	csr.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	InitiativeCountry := []string{}
	for _, x := range InitiativeData.Milestones {
		for _, country := range x.Country {
			if country != "" {
				if !stringInSlice(country, InitiativeCountry) {
					InitiativeCountry = append(InitiativeCountry, country)

				}
			}
		}

	}
	CountryTotal := "TOTAL"
	InitiativeCountry = append(InitiativeCountry, CountryTotal)

	// Ge Region Data
	GetCountryCode := map[string]string{}
	RegionData := []RegionModel{}
	csr, err = c.Ctx.Connection.NewQuery().From(new(RegionModel).TableName()).Cursor(nil)
	err = csr.Fetch(&RegionData, 0, false)
	csr.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	for _, x := range RegionData {
		GetCountryCode[x.Country] = x.CountryCode
	}

	// Get Detail Data
	query := []*db.Filter{}
	query = append(query, db.Eq("initiativeid", parm.InitiativeId))
	query = append(query, db.Eq("year", parm.Year))
	query = append(query, db.Eq("metricid", parm.MetricId))

	DetailData := []MilestoneValueModel{}
	csr, err = c.Ctx.Connection.NewQuery().From(new(MilestoneValueModel).TableName()).Where(db.And(query...)).Cursor(nil)
	err = csr.Fetch(&DetailData, 0, false)
	csr.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	// Get Metrics Data
	result := tk.M{}
	MetricData := []tk.M{}
	IsResultAvailable := false

	pipes, match, group, project, sort := []tk.M{}, []tk.M{}, tk.M{}, tk.M{}, tk.M{}
	match = append(match, tk.M{"initiativeid": parm.InitiativeId})
	match = append(match, tk.M{"year": parm.Year})
	// for _, x := range DetailData {
	// 	match = append(match, tk.M{"country": x.Country})
	// }
	//match = append(match, tk.M{"country": parm.Country})
	// match = append(match, country)
	match = append(match, tk.M{"metricid": tk.M{"$ne": ""}})
	err = Deserialize(`
		{"$group":{"_id":{
	        "MetricId":"$metricid",
	        "MetricName":"$metricname",
	        "Description":"$description",
	        "Denomination":"$denomination",
	        "ActualValue":"$actualvalue",
	        "TotalValue":"$totalvalue",
	        "RAGValue":"$ragvalue"
	    }}}
	`, &group)
	err = Deserialize(`
		{"$project":{
	        "MetricId":"$_id.MetricId",
	        "MetricName":"$_id.MetricName",
	        "Description":"$_id.Description",
	        "Denomination":"$_id.Denomination",
	        "ActualValue":"$_id.ActualValue",
	        "TotalValue":"$_id.TotalValue",
	        "RAGValue":"$_id.RAGValue"
	    }}
	`, &project)
	err = Deserialize(`
		{"$sort":{"MetricName":1}}
	`, &sort)
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": match}})
	pipes = append(pipes, group)
	pipes = append(pipes, project)
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
	result.Set("MetricData", MetricData)

	if len(MetricData) > 0 {
		IsResultAvailable = true
	}

	// Mapping Existing Initiative Country Data with Existing Data
	CurrentData := []MilestoneValueModel{}
	for _, x := range InitiativeCountry {
		IsAvailable := false
		data := MilestoneValueModel{}
		for _, d := range DetailData {
			if d.Country == x {
				IsAvailable = true
				data = d
				break
			}
		}
		if !IsAvailable {
			data.InitiativeId = parm.InitiativeId
			data.CountryCode = GetCountryCode[x]
			data.Country = x
			data.Year = parm.Year
			data.NAJan = true
			data.NAFeb = true
			data.NAMar = true
			data.NAApr = true
			data.NAMay = true
			data.NAJun = true
			data.NAJul = true
			data.NAAug = true
			data.NASep = true
			data.NAOct = true
			data.NANov = true
			data.NADec = true
			data.NATotalJan = true
			data.NATotalFeb = true
			data.NATotalMar = true
			data.NATotalApr = true
			data.NATotalMay = true
			data.NATotalJun = true
			data.NATotalJul = true
			data.NATotalAug = true
			data.NATotalSep = true
			data.NATotalOct = true
			data.NATotalNov = true
			data.NATotalDec = true
		}
		CurrentData = append(CurrentData, data)

	}

	DetailData = CurrentData

	return c.SetResultInfo(false, "", tk.M{}.Set("DetailData", DetailData).Set("IsResultAvailable", IsResultAvailable).Set("MetricData", MetricData))
}

func (c *AdoptionModuleController) SaveDetail(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	parm := struct {
		MetricData   []tk.M
		InitiativeId string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	result := []*MilestoneValueModel{}
	mId := tk.M{}
	e := c.Ctx.DeleteMany(new(MilestoneValueModel), db.And(db.Eq("initiativeid", parm.InitiativeId)))
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	for _, x := range parm.MetricData {
		// tk.Println("DETAIL DATA", x.Get("DetailData"))
		for _, d := range x.Get("DetailData").([]interface{}) {
			value := NewMilestoneValue()
			data := d.(map[string]interface{})

			id := strings.ToLower(x.GetString("MetricId"))
			metricId := ""
			if strings.Contains(id, "metric") {
				if mId.GetString(id) == "" {
					metricId = tk.RandomString(32)
					mId.Set(id, metricId)

				} else {
					metricId = mId.GetString(id)
				}

			} else {
				metricId = x.GetString("MetricId")
			}

			value.MetricId = metricId
			value.MetricName = x.GetString("MetricName")
			value.Denomination = x.GetString("Denomination")
			value.Description = x.GetString("Description")

			value.InitiativeId = data["InitiativeId"].(string)
			value.Country = data["Country"].(string)
			value.CountryCode = data["CountryCode"].(string)
			value.Year = int(data["Year"].(float64))
			value.ActualValue = x["ActualValue"].(bool)
			value.TotalValue = x["TotalValue"].(bool)
			value.RAGValue = x["RAGValue"].(bool)
			value.Jan = data["Jan"].(float64)
			value.Feb = data["Feb"].(float64)
			value.Mar = data["Mar"].(float64)
			value.Apr = data["Apr"].(float64)
			value.May = data["May"].(float64)
			value.Jun = data["Jun"].(float64)
			value.Jul = data["Jul"].(float64)
			value.Aug = data["Aug"].(float64)
			value.Sep = data["Sep"].(float64)
			value.Oct = data["Oct"].(float64)
			value.Nov = data["Nov"].(float64)
			value.Dec = data["Dec"].(float64)
			value.NAJan = data["NAJan"].(bool)
			value.NAFeb = data["NAFeb"].(bool)
			value.NAMar = data["NAMar"].(bool)
			value.NAApr = data["NAApr"].(bool)
			value.NAMay = data["NAMay"].(bool)
			value.NAJun = data["NAJun"].(bool)
			value.NAJul = data["NAJul"].(bool)
			value.NAAug = data["NAAug"].(bool)
			value.NASep = data["NASep"].(bool)
			value.NAOct = data["NAOct"].(bool)
			value.NANov = data["NANov"].(bool)
			value.NADec = data["NADec"].(bool)

			if value.TotalValue == false {
				value.TotalJan = 0
				value.TotalFeb = 0
				value.TotalMar = 0
				value.TotalApr = 0
				value.TotalMay = 0
				value.TotalJun = 0
				value.TotalJul = 0
				value.TotalAug = 0
				value.TotalSep = 0
				value.TotalOct = 0
				value.TotalNov = 0
				value.TotalDec = 0
				value.NATotalJan = true
				value.NATotalFeb = true
				value.NATotalMar = true
				value.NATotalApr = true
				value.NATotalMay = true
				value.NATotalJun = true
				value.NATotalJul = true
				value.NATotalAug = true
				value.NATotalSep = true
				value.NATotalOct = true
				value.NATotalNov = true
				value.NATotalDec = true
			} else {
				value.TotalJan = data["TotalJan"].(float64)
				value.TotalFeb = data["TotalFeb"].(float64)
				value.TotalMar = data["TotalMar"].(float64)
				value.TotalApr = data["TotalApr"].(float64)
				value.TotalMay = data["TotalMay"].(float64)
				value.TotalJun = data["TotalJun"].(float64)
				value.TotalJul = data["TotalJul"].(float64)
				value.TotalAug = data["TotalAug"].(float64)
				value.TotalSep = data["TotalSep"].(float64)
				value.TotalOct = data["TotalOct"].(float64)
				value.TotalNov = data["TotalNov"].(float64)
				value.TotalDec = data["TotalDec"].(float64)
				value.NATotalJan = data["NATotalJan"].(bool)
				value.NATotalFeb = data["NATotalFeb"].(bool)
				value.NATotalMar = data["NATotalMar"].(bool)
				value.NATotalApr = data["NATotalApr"].(bool)
				value.NATotalMay = data["NATotalMay"].(bool)
				value.NATotalJun = data["NATotalJun"].(bool)
				value.NATotalJul = data["NATotalJul"].(bool)
				value.NATotalAug = data["NATotalAug"].(bool)
				value.NATotalSep = data["NATotalSep"].(bool)
				value.NATotalOct = data["NATotalOct"].(bool)
				value.NATotalNov = data["NATotalNov"].(bool)
				value.NATotalDec = data["NATotalDec"].(bool)
			}

			if value.RAGValue == false {
				value.RAGJan = ""
				value.RAGFeb = ""
				value.RAGMar = ""
				value.RAGApr = ""
				value.RAGMay = ""
				value.RAGJun = ""
				value.RAGJul = ""
				value.RAGAug = ""
				value.RAGSep = ""
				value.RAGOct = ""
				value.RAGNov = ""
				value.RAGDec = ""
			} else {
				value.RAGJan = data["RAGJan"].(string)
				value.RAGFeb = data["RAGFeb"].(string)
				value.RAGMar = data["RAGMar"].(string)
				value.RAGApr = data["RAGApr"].(string)
				value.RAGMay = data["RAGMay"].(string)
				value.RAGJun = data["RAGJun"].(string)
				value.RAGJul = data["RAGJul"].(string)
				value.RAGAug = data["RAGAug"].(string)
				value.RAGSep = data["RAGSep"].(string)
				value.RAGOct = data["RAGOct"].(string)
				value.RAGNov = data["RAGNov"].(string)
				value.RAGDec = data["RAGDec"].(string)

				if value.RAGDec != "" {
					value.LatestRAG = value.RAGDec
					value.PeriodRAG = "Dec"
				} else if value.RAGNov != "" {
					value.LatestRAG = value.RAGNov
					value.PeriodRAG = "Nov"
				} else if value.RAGOct != "" {
					value.LatestRAG = value.RAGOct
					value.PeriodRAG = "Oct"
				} else if value.RAGSep != "" {
					value.LatestRAG = value.RAGSep
					value.PeriodRAG = "Sep"
				} else if value.RAGAug != "" {
					value.LatestRAG = value.RAGAug
					value.PeriodRAG = "Aug"
				} else if value.RAGJul != "" {
					value.LatestRAG = value.RAGJul
					value.PeriodRAG = "Jul"
				} else if value.RAGJun != "" {
					value.LatestRAG = value.RAGJun
					value.PeriodRAG = "Jun"
				} else if value.RAGMay != "" {
					value.LatestRAG = value.RAGMay
					value.PeriodRAG = "May"
				} else if value.RAGApr != "" {
					value.LatestRAG = value.RAGApr
					value.PeriodRAG = "Apr"
				} else if value.RAGMar != "" {
					value.LatestRAG = value.RAGMar
					value.PeriodRAG = "Mar"
				} else if value.RAGFeb != "" {
					value.LatestRAG = value.RAGFeb
					value.PeriodRAG = "Feb"
				} else if value.RAGJan != "" {
					value.LatestRAG = value.RAGJan
					value.PeriodRAG = "Jan"
				}

			}

			value.Created_By = c.GetCurrentUser(k)
			value.Created_Date = time.Now().UTC()
			value.Updated_By = data["Updated_By"].(string)
			value.Updated_Date = time.Now().UTC()

			e = c.Ctx.Save(value)
			if e != nil {
				return c.ErrorResultInfo(e.Error(), nil)
			}
			// tk.Println("MetricName", value.MetricName)
			// tk.Println("COUNTRY", value.Country)
			// tk.Println("Jan", value.Jan)
			result = append(result, value)
		}

	}

	return c.SetResultInfo(false, "", result)
}

func (c *AdoptionModuleController) SaveDetailOld(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	parm := struct {
		MetricData   map[string]tk.M
		InitiativeId string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	result := []*MilestoneValueModel{}
	mId := tk.M{}
	e := c.Ctx.DeleteMany(new(MilestoneValueModel), db.And(db.Eq("initiativeid", parm.InitiativeId)))
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	for _, x := range parm.MetricData {
		for _, d := range x.Get("Data").([]interface{}) {
			value := NewMilestoneValue()
			data := d.(map[string]interface{})

			id := strings.ToLower(x.GetString("MetricId"))
			metricId := ""
			if strings.Contains(id, "metric") {
				if mId.GetString(id) == "" {
					metricId = tk.RandomString(32)
					mId.Set(id, metricId)

				} else {
					metricId = mId.GetString(id)
				}

			} else {
				metricId = x.GetString("MetricId")
			}

			value.MetricId = metricId
			value.MetricName = x.GetString("MetricName")
			value.Denomination = x.GetString("Denomination")
			value.Description = x.GetString("Description")

			value.InitiativeId = data["InitiativeId"].(string)
			value.Country = data["Country"].(string)
			value.CountryCode = data["CountryCode"].(string)
			value.Year = int(data["Year"].(float64))
			value.ActualValue = x["ActualValue"].(bool)
			value.TotalValue = x["TotalValue"].(bool)
			value.RAGValue = x["RAGValue"].(bool)
			value.Jan = data["Jan"].(float64)
			value.Feb = data["Feb"].(float64)
			value.Mar = data["Mar"].(float64)
			value.Apr = data["Apr"].(float64)
			value.May = data["May"].(float64)
			value.Jun = data["Jun"].(float64)
			value.Jul = data["Jul"].(float64)
			value.Aug = data["Aug"].(float64)
			value.Sep = data["Sep"].(float64)
			value.Oct = data["Oct"].(float64)
			value.Nov = data["Nov"].(float64)
			value.Dec = data["Dec"].(float64)
			value.NAJan = data["NAJan"].(bool)
			value.NAFeb = data["NAFeb"].(bool)
			value.NAMar = data["NAMar"].(bool)
			value.NAApr = data["NAApr"].(bool)
			value.NAMay = data["NAMay"].(bool)
			value.NAJun = data["NAJun"].(bool)
			value.NAJul = data["NAJul"].(bool)
			value.NAAug = data["NAAug"].(bool)
			value.NASep = data["NASep"].(bool)
			value.NAOct = data["NAOct"].(bool)
			value.NANov = data["NANov"].(bool)
			value.NADec = data["NADec"].(bool)
			value.TotalJan = data["TotalJan"].(float64)
			value.TotalFeb = data["TotalFeb"].(float64)
			value.TotalMar = data["TotalMar"].(float64)
			value.TotalApr = data["TotalApr"].(float64)
			value.TotalMay = data["TotalMay"].(float64)
			value.TotalJun = data["TotalJun"].(float64)
			value.TotalJul = data["TotalJul"].(float64)
			value.TotalAug = data["TotalAug"].(float64)
			value.TotalSep = data["TotalSep"].(float64)
			value.TotalOct = data["TotalOct"].(float64)
			value.TotalNov = data["TotalNov"].(float64)
			value.TotalDec = data["TotalDec"].(float64)
			value.NATotalJan = data["NATotalJan"].(bool)
			value.NATotalFeb = data["NATotalFeb"].(bool)
			value.NATotalMar = data["NATotalMar"].(bool)
			value.NATotalApr = data["NATotalApr"].(bool)
			value.NATotalMay = data["NATotalMay"].(bool)
			value.NATotalJun = data["NATotalJun"].(bool)
			value.NATotalJul = data["NATotalJul"].(bool)
			value.NATotalAug = data["NATotalAug"].(bool)
			value.NATotalSep = data["NATotalSep"].(bool)
			value.NATotalOct = data["NATotalOct"].(bool)
			value.NATotalNov = data["NATotalNov"].(bool)
			value.NATotalDec = data["NATotalDec"].(bool)
			value.RAGJan = data["RAGJan"].(string)
			value.RAGFeb = data["RAGFeb"].(string)
			value.RAGMar = data["RAGMar"].(string)
			value.RAGApr = data["RAGApr"].(string)
			value.RAGMay = data["RAGMay"].(string)
			value.RAGJun = data["RAGJun"].(string)
			value.RAGJul = data["RAGJul"].(string)
			value.RAGAug = data["RAGAug"].(string)
			value.RAGSep = data["RAGSep"].(string)
			value.RAGOct = data["RAGOct"].(string)
			value.RAGNov = data["RAGNov"].(string)
			value.RAGDec = data["RAGDec"].(string)

			value.Created_By = c.GetCurrentUser(k)
			value.Created_Date = time.Now().UTC()
			value.Updated_By = data["Updated_By"].(string)
			value.Updated_Date = time.Now().UTC()

			e = c.Ctx.Save(value)
			if e != nil {
				return c.ErrorResultInfo(e.Error(), nil)
			}
			// tk.Println("MetricName", value.MetricName)
			// tk.Println("COUNTRY", value.Country)
			// tk.Println("Jan", value.Jan)
			result = append(result, value)
		}

	}

	return c.SetResultInfo(false, "", result)
}
