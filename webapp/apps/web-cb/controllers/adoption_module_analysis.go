package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"time"
)

func GetRAGValue(val string) float64 {
	result := 0.0
	switch val {
	case "red":
		result = 1
		break
	case "amber":
		result = 2
		break
	case "green":
		result = 3
		break
	default:
		result = 0
		break
	}
	return result
}
func (c *AdoptionModuleController) GetAnalysisData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		InitiativeId string
		Year         int
		Country      string
		MetricID     string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	result := tk.M{}
	MetricData := []tk.M{}
	InitiativeData := InitiativeModel{}
	DetailData := []MilestoneValueModel{}
	GoLiveDate := map[string]time.Time{}
	AnalyticData := map[string]float64{}
	NAData := map[string]bool{}
	IsResultAvailable := false

	// Get Metrics Data
	pipes, match, group, project, sort := []tk.M{}, []tk.M{}, tk.M{}, tk.M{}, tk.M{}
	match = append(match, tk.M{"initiativeid": parm.InitiativeId})
	match = append(match, tk.M{"year": parm.Year})
	if parm.Country != "" {
		match = append(match, tk.M{"country": parm.Country})
	}
	match = append(match, tk.M{"metricid": tk.M{"$ne": ""}})
	err = Deserialize(`
		{"$group":{"_id":{
	        "MetricId":"$metricid",
	        "MetricName":"$metricname",
	        "Denomination":"$denomination",
	        "Description":"$description"
	    }}}
	`, &group)
	err = Deserialize(`
		{"$project":{
	        "MetricId":"$_id.MetricId",
	        "MetricName":"$_id.MetricName",
	        "Denomination":"$_id.Denomination",
	        "Description":"$_id.Description"
	    }}
	`, &project)
	err = Deserialize(`
		{"$sort":{"MetricName":1}}
	`, &sort)
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": match}})
	pipes = append(pipes, group)
	pipes = append(pipes, project)
	pipes = append(pipes, sort)
	csr, err := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From(new(MilestoneValueModel).TableName()).Cursor(nil)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	err = csr.Fetch(&MetricData, 0, false)
	csr.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	result.Set("MetricData", MetricData)
	TotalValue := false
	if len(MetricData) > 0 {
		IsResultAvailable = true
		// Get Initiative Country Data: InitiativeId
		csr, err = c.Ctx.Connection.NewQuery().From(new(InitiativeModel).TableName()).Where(db.Eq("InitiativeID", parm.InitiativeId)).Take(1).Cursor(nil)
		err = csr.Fetch(&InitiativeData, 1, false)
		csr.Close()
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
		InitiativeCountry := []string{}
		for _, x := range InitiativeData.Milestones {
			LiveDate := x.EndDate
			for _, country := range x.Country {
				if country != "" {
					if !stringInSlice(country, InitiativeCountry) {
						InitiativeCountry = append(InitiativeCountry, country)
						GoLiveDate[country] = LiveDate
					}
				}
			}
		}

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
		if parm.MetricID == "" {
			parm.MetricID = MetricData[0].GetString("MetricId")
		}
		query = append(query, db.Eq("metricid", parm.MetricID))
		if parm.Country != "" {
			query = append(query, db.Eq("country", parm.Country))
		}
		csr, err = c.Ctx.Connection.NewQuery().From(new(MilestoneValueModel).TableName()).Where(db.And(query...)).Cursor(nil)
		err = csr.Fetch(&DetailData, 0, false)
		csr.Close()
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}

		// Get Total RAG Data
		query = []*db.Filter{}
		query = append(query, db.Eq("initiativeid", parm.InitiativeId))
		query = append(query, db.Eq("year", parm.Year))
		if parm.MetricID == "" {
			parm.MetricID = MetricData[0].GetString("MetricId")
		}
		query = append(query, db.Eq("metricid", parm.MetricID))
		query = append(query, db.Eq("country", "TOTAL"))
		TotalData := []MilestoneValueModel{}
		csr, err = c.Ctx.Connection.NewQuery().From(new(MilestoneValueModel).TableName()).Where(db.And(query...)).Cursor(nil)
		err = csr.Fetch(&TotalData, 0, false)
		csr.Close()
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
		for _, x := range TotalData {
			AnalyticData["RAGJan"] = GetRAGValue(x.RAGJan)
			AnalyticData["RAGFeb"] = GetRAGValue(x.RAGFeb)
			AnalyticData["RAGMar"] = GetRAGValue(x.RAGMar)
			AnalyticData["RAGApr"] = GetRAGValue(x.RAGApr)
			AnalyticData["RAGMay"] = GetRAGValue(x.RAGMay)
			AnalyticData["RAGJun"] = GetRAGValue(x.RAGJun)
			AnalyticData["RAGJul"] = GetRAGValue(x.RAGJul)
			AnalyticData["RAGAug"] = GetRAGValue(x.RAGAug)
			AnalyticData["RAGSep"] = GetRAGValue(x.RAGSep)
			AnalyticData["RAGOct"] = GetRAGValue(x.RAGOct)
			AnalyticData["RAGNov"] = GetRAGValue(x.RAGNov)
			AnalyticData["RAGDec"] = GetRAGValue(x.RAGDec)
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

		// Get Analytic Data
		// Data Source => Cumulative, MoM
		AnalyticData["SumJan"] = 0.0
		AnalyticData["SumFeb"] = 0.0
		AnalyticData["SumMar"] = 0.0
		AnalyticData["SumApr"] = 0.0
		AnalyticData["SumMay"] = 0.0
		AnalyticData["SumJun"] = 0.0
		AnalyticData["SumJul"] = 0.0
		AnalyticData["SumAug"] = 0.0
		AnalyticData["SumSep"] = 0.0
		AnalyticData["SumOct"] = 0.0
		AnalyticData["SumNov"] = 0.0
		AnalyticData["SumDec"] = 0.0
		NASumJan := true
		NASumFeb := true
		NASumMar := true
		NASumApr := true
		NASumMay := true
		NASumJun := true
		NASumJul := true
		NASumAug := true
		NASumSep := true
		NASumOct := true
		NASumNov := true
		NASumDec := true

		AnalyticData["SumTotalJan"] = 0.0
		AnalyticData["SumTotalFeb"] = 0.0
		AnalyticData["SumTotalMar"] = 0.0
		AnalyticData["SumTotalApr"] = 0.0
		AnalyticData["SumTotalMay"] = 0.0
		AnalyticData["SumTotalJun"] = 0.0
		AnalyticData["SumTotalJul"] = 0.0
		AnalyticData["SumTotalAug"] = 0.0
		AnalyticData["SumTotalSep"] = 0.0
		AnalyticData["SumTotalOct"] = 0.0
		AnalyticData["SumTotalNov"] = 0.0
		AnalyticData["SumTotalDec"] = 0.0
		NASumTotalJan := true
		NASumTotalFeb := true
		NASumTotalMar := true
		NASumTotalApr := true
		NASumTotalMay := true
		NASumTotalJun := true
		NASumTotalJul := true
		NASumTotalAug := true
		NASumTotalSep := true
		NASumTotalOct := true
		NASumTotalNov := true
		NASumTotalDec := true

		for _, x := range DetailData {

			if x.TotalValue {
				TotalValue = true
			}
			// Actual Data
			if !x.NAJan {
				AnalyticData["Jan"] = x.Jan
				if NASumJan {
					NASumJan = false
				}
			}
			if !x.NAFeb {
				AnalyticData["Feb"] = x.Feb
				if NASumFeb {
					NASumFeb = false
				}
			}
			if !x.NAMar {
				AnalyticData["Mar"] = x.Mar
				if NASumMar {
					NASumMar = false
				}
			}
			if !x.NAApr {
				AnalyticData["Apr"] = x.Apr
				if NASumApr {
					NASumApr = false
				}
			}
			if !x.NAMay {
				AnalyticData["May"] = x.May
				if NASumMay {
					NASumMay = false
				}
			}
			if !x.NAJun {
				AnalyticData["Jun"] = x.Jun
				if NASumJun {
					NASumJun = false
				}
			}
			if !x.NAJul {
				AnalyticData["Jul"] = x.Jul
				if NASumJul {
					NASumJul = false
				}
			}
			if !x.NAAug {
				AnalyticData["Aug"] = x.Aug
				if NASumAug {
					NASumAug = false
				}
			}
			if !x.NASep {
				AnalyticData["Sep"] = x.Sep
				if NASumSep {
					NASumSep = false
				}
			}
			if !x.NAOct {
				AnalyticData["Oct"] = x.Oct
				if NASumOct {
					NASumOct = false
				}
			}
			if !x.NANov {
				AnalyticData["Nov"] = x.Nov
				if NASumNov {
					NASumNov = false
				}
			}
			if !x.NADec {
				AnalyticData["Dec"] = x.Dec
				if NASumDec {
					NASumDec = false
				}
			}
			AnalyticData["SumJan"] += x.Jan
			AnalyticData["SumFeb"] += x.Feb
			AnalyticData["SumMar"] += x.Mar
			AnalyticData["SumApr"] += x.Apr
			AnalyticData["SumMay"] += x.May
			AnalyticData["SumJun"] += x.Jun
			AnalyticData["SumJul"] += x.Jul
			AnalyticData["SumAug"] += x.Aug
			AnalyticData["SumSep"] += x.Sep
			AnalyticData["SumOct"] += x.Oct
			AnalyticData["SumNov"] += x.Nov
			AnalyticData["SumDec"] += x.Dec

			// Total Data
			if !x.NATotalJan {
				AnalyticData["TotalJan"] = x.TotalJan
				if NASumTotalJan {
					NASumTotalJan = false
				}
			}
			if !x.NATotalFeb {
				AnalyticData["TotalFeb"] = x.TotalFeb
				if NASumTotalFeb {
					NASumTotalFeb = false
				}
			}
			if !x.NATotalMar {
				AnalyticData["TotalMar"] = x.TotalMar
				if NASumTotalMar {
					NASumTotalMar = false
				}
			}
			if !x.NATotalApr {
				AnalyticData["TotalApr"] = x.TotalApr
				if NASumTotalApr {
					NASumTotalApr = false
				}
			}
			if !x.NATotalMay {
				AnalyticData["TotalMay"] = x.TotalMay
				if NASumTotalMay {
					NASumTotalMay = false
				}
			}
			if !x.NATotalJun {
				AnalyticData["TotalJun"] = x.TotalJun
				if NASumTotalJun {
					NASumTotalJun = false
				}
			}
			if !x.NATotalJul {
				AnalyticData["TotalJul"] = x.TotalJul
				if NASumTotalJul {
					NASumTotalJul = false
				}
			}
			if !x.NATotalAug {
				AnalyticData["TotalAug"] = x.TotalAug
				if NASumTotalAug {
					NASumTotalAug = false
				}
			}
			if !x.NATotalSep {
				AnalyticData["TotalSep"] = x.TotalSep
				if NASumTotalSep {
					NASumTotalSep = false
				}
			}
			if !x.NATotalOct {
				AnalyticData["TotalOct"] = x.TotalOct
				if NASumTotalOct {
					NASumTotalOct = false
				}
			}
			if !x.NATotalNov {
				AnalyticData["TotalNov"] = x.TotalNov
				if NASumTotalNov {
					NASumTotalNov = false
				}
			}
			if !x.NATotalDec {
				AnalyticData["TotalDec"] = x.TotalDec
				if NASumTotalDec {
					NASumTotalDec = false
				}
			}
			AnalyticData["SumTotalJan"] += x.TotalJan
			AnalyticData["SumTotalFeb"] += x.TotalFeb
			AnalyticData["SumTotalMar"] += x.TotalMar
			AnalyticData["SumTotalApr"] += x.TotalApr
			AnalyticData["SumTotalMay"] += x.TotalMay
			AnalyticData["SumTotalJun"] += x.TotalJun
			AnalyticData["SumTotalJul"] += x.TotalJul
			AnalyticData["SumTotalAug"] += x.TotalAug
			AnalyticData["SumTotalSep"] += x.TotalSep
			AnalyticData["SumTotalOct"] += x.TotalOct
			AnalyticData["SumTotalNov"] += x.TotalNov
			AnalyticData["SumTotalDec"] += x.TotalDec

		}

		// Actual
		AnalyticData["CumulativeJan"] = AnalyticData["SumJan"]
		AnalyticData["CumulativeFeb"] = AnalyticData["SumJan"] + AnalyticData["SumFeb"]
		AnalyticData["CumulativeMar"] = AnalyticData["SumJan"] + AnalyticData["SumFeb"] + AnalyticData["SumMar"]
		AnalyticData["CumulativeApr"] = AnalyticData["SumJan"] + AnalyticData["SumFeb"] + AnalyticData["SumMar"] + AnalyticData["SumApr"]
		AnalyticData["CumulativeMay"] = AnalyticData["SumJan"] + AnalyticData["SumFeb"] + AnalyticData["SumMar"] + AnalyticData["SumApr"] + AnalyticData["SumMay"]
		AnalyticData["CumulativeJun"] = AnalyticData["SumJan"] + AnalyticData["SumFeb"] + AnalyticData["SumMar"] + AnalyticData["SumApr"] + AnalyticData["SumMay"] + AnalyticData["SumJun"]
		AnalyticData["CumulativeJul"] = AnalyticData["SumJan"] + AnalyticData["SumFeb"] + AnalyticData["SumMar"] + AnalyticData["SumApr"] + AnalyticData["SumMay"] + AnalyticData["SumJun"] + AnalyticData["SumJul"]
		AnalyticData["CumulativeAug"] = AnalyticData["SumJan"] + AnalyticData["SumFeb"] + AnalyticData["SumMar"] + AnalyticData["SumApr"] + AnalyticData["SumMay"] + AnalyticData["SumJun"] + AnalyticData["SumJul"] + AnalyticData["SumAug"]
		AnalyticData["CumulativeSep"] = AnalyticData["SumJan"] + AnalyticData["SumFeb"] + AnalyticData["SumMar"] + AnalyticData["SumApr"] + AnalyticData["SumMay"] + AnalyticData["SumJun"] + AnalyticData["SumJul"] + AnalyticData["SumAug"] + AnalyticData["SumSep"]
		AnalyticData["CumulativeOct"] = AnalyticData["SumJan"] + AnalyticData["SumFeb"] + AnalyticData["SumMar"] + AnalyticData["SumApr"] + AnalyticData["SumMay"] + AnalyticData["SumJun"] + AnalyticData["SumJul"] + AnalyticData["SumAug"] + AnalyticData["SumSep"] + AnalyticData["SumOct"]
		AnalyticData["CumulativeNov"] = AnalyticData["SumJan"] + AnalyticData["SumFeb"] + AnalyticData["SumMar"] + AnalyticData["SumApr"] + AnalyticData["SumMay"] + AnalyticData["SumJun"] + AnalyticData["SumJul"] + AnalyticData["SumAug"] + AnalyticData["SumSep"] + AnalyticData["SumOct"] + AnalyticData["SumNov"]
		AnalyticData["CumulativeDec"] = AnalyticData["SumJan"] + AnalyticData["SumFeb"] + AnalyticData["SumMar"] + AnalyticData["SumApr"] + AnalyticData["SumMay"] + AnalyticData["SumJun"] + AnalyticData["SumJul"] + AnalyticData["SumAug"] + AnalyticData["SumSep"] + AnalyticData["SumOct"] + AnalyticData["SumNov"] + AnalyticData["SumDec"]

		NAData["NASumJan"] = NASumJan
		NAData["NASumFeb"] = NASumFeb
		NAData["NASumMar"] = NASumMar
		NAData["NASumApr"] = NASumApr
		NAData["NASumMay"] = NASumMay
		NAData["NASumJun"] = NASumJun
		NAData["NASumJul"] = NASumJul
		NAData["NASumAug"] = NASumAug
		NAData["NASumSep"] = NASumSep
		NAData["NASumOct"] = NASumOct
		NAData["NASumNov"] = NASumNov
		NAData["NASumDec"] = NASumDec

		// Total
		AnalyticData["CumulativeTotalJan"] = AnalyticData["SumTotalJan"]
		AnalyticData["CumulativeTotalFeb"] = AnalyticData["SumTotalJan"] + AnalyticData["SumTotalFeb"]
		AnalyticData["CumulativeTotalMar"] = AnalyticData["SumTotalJan"] + AnalyticData["SumTotalFeb"] + AnalyticData["SumTotalMar"]
		AnalyticData["CumulativeTotalApr"] = AnalyticData["SumTotalJan"] + AnalyticData["SumTotalFeb"] + AnalyticData["SumTotalMar"] + AnalyticData["SumTotalApr"]
		AnalyticData["CumulativeTotalMay"] = AnalyticData["SumTotalJan"] + AnalyticData["SumTotalFeb"] + AnalyticData["SumTotalMar"] + AnalyticData["SumTotalApr"] + AnalyticData["SumTotalMay"]
		AnalyticData["CumulativeTotalJun"] = AnalyticData["SumTotalJan"] + AnalyticData["SumTotalFeb"] + AnalyticData["SumTotalMar"] + AnalyticData["SumTotalApr"] + AnalyticData["SumTotalMay"] + AnalyticData["SumTotalJun"]
		AnalyticData["CumulativeTotalJul"] = AnalyticData["SumTotalJan"] + AnalyticData["SumTotalFeb"] + AnalyticData["SumTotalMar"] + AnalyticData["SumTotalApr"] + AnalyticData["SumTotalMay"] + AnalyticData["SumTotalJun"] + AnalyticData["SumTotalJul"]
		AnalyticData["CumulativeTotalAug"] = AnalyticData["SumTotalJan"] + AnalyticData["SumTotalFeb"] + AnalyticData["SumTotalMar"] + AnalyticData["SumTotalApr"] + AnalyticData["SumTotalMay"] + AnalyticData["SumTotalJun"] + AnalyticData["SumTotalJul"] + AnalyticData["SumTotalAug"]
		AnalyticData["CumulativeTotalSep"] = AnalyticData["SumTotalJan"] + AnalyticData["SumTotalFeb"] + AnalyticData["SumTotalMar"] + AnalyticData["SumTotalApr"] + AnalyticData["SumTotalMay"] + AnalyticData["SumTotalJun"] + AnalyticData["SumTotalJul"] + AnalyticData["SumTotalAug"] + AnalyticData["SumTotalSep"]
		AnalyticData["CumulativeTotalOct"] = AnalyticData["SumTotalJan"] + AnalyticData["SumTotalFeb"] + AnalyticData["SumTotalMar"] + AnalyticData["SumTotalApr"] + AnalyticData["SumTotalMay"] + AnalyticData["SumTotalJun"] + AnalyticData["SumTotalJul"] + AnalyticData["SumTotalAug"] + AnalyticData["SumTotalSep"] + AnalyticData["SumTotalOct"]
		AnalyticData["CumulativeTotalNov"] = AnalyticData["SumTotalJan"] + AnalyticData["SumTotalFeb"] + AnalyticData["SumTotalMar"] + AnalyticData["SumTotalApr"] + AnalyticData["SumTotalMay"] + AnalyticData["SumTotalJun"] + AnalyticData["SumTotalJul"] + AnalyticData["SumTotalAug"] + AnalyticData["SumTotalSep"] + AnalyticData["SumTotalOct"] + AnalyticData["SumTotalNov"]
		AnalyticData["CumulativeTotalDec"] = AnalyticData["SumTotalJan"] + AnalyticData["SumTotalFeb"] + AnalyticData["SumTotalMar"] + AnalyticData["SumTotalApr"] + AnalyticData["SumTotalMay"] + AnalyticData["SumTotalJun"] + AnalyticData["SumTotalJul"] + AnalyticData["SumTotalAug"] + AnalyticData["SumTotalSep"] + AnalyticData["SumTotalOct"] + AnalyticData["SumTotalNov"] + AnalyticData["SumTotalDec"]

		NAData["NASumTotalJan"] = NASumTotalJan
		NAData["NASumTotalFeb"] = NASumTotalFeb
		NAData["NASumTotalMar"] = NASumTotalMar
		NAData["NASumTotalApr"] = NASumTotalApr
		NAData["NASumTotalMay"] = NASumTotalMay
		NAData["NASumTotalJun"] = NASumTotalJun
		NAData["NASumTotalJul"] = NASumTotalJul
		NAData["NASumTotalAug"] = NASumTotalAug
		NAData["NASumTotalSep"] = NASumTotalSep
		NAData["NASumTotalOct"] = NASumTotalOct
		NAData["NASumTotalNov"] = NASumTotalNov
		NAData["NASumTotalDec"] = NASumTotalDec

	}
	result.Set("ActiveMetric", parm.MetricID)
	result.Set("IsResultAvailable", IsResultAvailable)
	result.Set("DetailData", DetailData)
	result.Set("AnalyticData", AnalyticData)
	result.Set("NASumData", NAData)
	result.Set("GoLiveDate", GoLiveDate)
	result.Set("TotalValue", TotalValue)
	return c.SetResultInfo(false, "", result)
}
