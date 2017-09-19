package controllers

import (
	// . "eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// tk "github.com/eaciit/toolkit"
	// "fmt"
	"github.com/tealeg/xlsx"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"
)

type StagingAreaController struct {
	*BaseController
}

func (c *StagingAreaController) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	c.Action(k, "Metric Validation", "Open Country Analysis Page", "", "", "", "", "")
	StagingArea := c.GetAccess(k, "STAGINGAREA")
	Initiative := c.GetAccess(k, "INITIATIVE")
	k.Config.LayoutTemplate = "_layout-v2.html"
	PartialFiles := []string{}
	PartialFiles = append(PartialFiles, "shared/sidebar.html")
	PartialFiles = append(PartialFiles, "shared/initiative_summary.html")
	PartialFiles = append(PartialFiles, "shared/initiative_chart.html")
	PartialFiles = append(PartialFiles, "shared/metric_upload.html")
	PartialFiles = append(PartialFiles, "dashboard/initiative.html")

	k.Config.IncludeFiles = PartialFiles
	UserCountry := ""
	if k.Session("country") != nil {
		UserCountry = k.Session("country").(string)
	}
	return tk.M{}.Set("UserCountry", UserCountry).Set("StagingArea", StagingArea).Set("Initiative", Initiative)

}
func (c *StagingAreaController) getBMData(data tk.M, year int, countryCode string, bmID string, EndingPeriod time.Time) {
	query := []*dbox.Filter{}
	query = append(query, dbox.Eq("bmid", bmID))
	query = append(query, dbox.Eq("year", year))
	query = append(query, dbox.Eq("countrycode", countryCode))

	BMData := []BusinessMetricsDataTempModel{}
	csr, err := c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataTempModel).TableName()).Where(dbox.And(query...)).Cursor(nil)
	err = csr.Fetch(&BMData, 0, false)
	csr.Close()
	if err != nil {
		BMData = []BusinessMetricsDataTempModel{}
	}

	EndingMonth, _ := strconv.Atoi(EndingPeriod.Format("1"))
	baseline, ytdactual, fyforecast, target := 130895111188.0, 130895111188.0, 130895111188.0, 130895111188.0
	for i := int(1); i <= EndingMonth; i++ {
		month := strconv.Itoa(i)
		tempPeriod := strconv.Itoa(year)
		if i < 10 {
			tempPeriod += "0" + month + "01"
		} else {
			tempPeriod += month + "01"
		}
		period, _ := time.Parse("20060102", tempPeriod)
		actual := 0.0
		rag := ""
		budget := 0.0
		isAny := false
		for _, bm := range BMData {
			if bm.Period.UTC() != period {
				continue
			}
			baseline = bm.Baseline
			if bm.ActualYTD != 130895111188 {
				ytdactual = bm.ActualYTD
			}
			fyforecast = bm.FullYearForecast
			target = bm.Target
			actual = bm.ActualYTD
			rag = bm.RAG
			budget = bm.Budget
			isAny = true
			break
		}
		if !isAny {
			actual = 130895111188
			budget = 130895111188
		}
		data.Set("Budget"+period.Format("Jan2006"), budget)
		data.Set("RAG"+period.Format("Jan2006"), rag)
		data.Set(period.Format("Jan2006"), actual)
	}
	data.Set("BaseLine", baseline)
	data.Set("YTDActual", ytdactual)
	data.Set("FullYearForecast", fyforecast)
	data.Set("Target", target)
}
func (c *StagingAreaController) GetData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Year int
		To   string
		From string
		BMID string
		// Country string
		// Region  string
		Country []interface{}
		Region  []interface{}
		Search  string
	}{}

	err := k.GetPayload(&parm)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	startingMonth, endingMonth := 1, 12
	if parm.From == "" && parm.To == "" {
		parm.Year = time.Now().Year()
	} else if parm.To != "" {
		period, _ := time.Parse("200601", parm.To)
		parm.Year = period.Year()
		endingMonth = int(period.Month())
		if parm.From != "" {
			period, _ := time.Parse("200601", parm.From)
			startingMonth = int(period.Month())
		}
	} else if parm.From != "" {
		period, _ := time.Parse("200601", parm.From)
		parm.Year = period.Year()
		startingMonth = int(period.Month())
	}

	// Get Selected Period
	year := strconv.Itoa(parm.Year)
	selectedPeriod := []tk.M{}

	for i := int(startingMonth); i <= endingMonth; i++ {
		month := strconv.Itoa(i)
		tempPeriod := year
		if i < 10 {
			tempPeriod += "0" + month + "01"
		} else {
			tempPeriod += month + "01"
		}

		period, _ := time.Parse("20060102", tempPeriod)
		data := tk.M{}
		data.Set("Title", period.Format("Jan-2006"))
		data.Set("Period", period.Format("Jan2006"))
		selectedPeriod = append(selectedPeriod, data)
	}

	// Get Selected Data
	result := []tk.M{}

	BusinessDriverL1 := make([]BusinessDriverL1Model, 0)
	csr, err := c.Ctx.Find(new(BusinessDriverL1Model), nil)
	err = csr.Fetch(&BusinessDriverL1, 0, false)
	csr.Close()

	//get notification
	pipe := []tk.M{
		tk.M{
			"$group": tk.M{
				"_id": tk.M{
					"bmid":   "$bmid",
					"bmname": "$bmname",
					"bmtype": "$type",
				},
				"lastUpdatePeriod": tk.M{
					"$last": "$updated_date",
				},
			},
		},
	}
	crs, err := c.Ctx.Connection.NewQuery().Command("pipe", pipe).From(NewBusinessMetricsNotificationModel().TableName()).Cursor(nil)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	dataNotif := tk.Ms{}
	err = crs.Fetch(&dataNotif, 0, false)
	crs.Close()
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	dateUpdated := []interface{}{}
	for _, v := range dataNotif {
		dateUpdated = append(dateUpdated, v.Get("lastUpdatePeriod").(time.Time).UTC())
	}

	err, notif := c.getNotif(dateUpdated, k)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	for _, bd := range BusinessDriverL1 {
		for _, bm := range bd.BusinessMetric {
			if parm.BMID != "" && parm.BMID != bm.Id {
				continue
			}
			if parm.Search != "" && parm.BMID == "" && !strings.Contains(strings.ToLower(bm.DataPoint), parm.Search) {
				continue
			}
			bmData := tk.M{}
			bmData.Set("BMId", bm.Id)
			bmData.Set("DataPoint", bm.DataPoint)
			bmData.Set("Description", bm.Description)
			bmData.Set("NotificationFinance", notif[bm.Id+"|finance"])
			bmData.Set("NotificationOwner", notif[bm.Id+"|owner"])
			bmData.Set("DecimalFormat", bm.DecimalFormat)
			bmData.Set("MetricType", bm.MetricType)
			result = append(result, bmData)
		}
	}
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	return c.SetResultInfo(false, "", tk.M{}.Set("SelectedData", result).Set("SelectedPeriod", selectedPeriod))
}

func (c *StagingAreaController) getNotif(dateUpdated []interface{}, k *knot.WebContext) (error, map[string][]BusinessMetricsNotificationModel) {
	crs, err := c.Ctx.Connection.NewQuery().From(NewBusinessMetricsNotificationModel().TableName()).
		Where(dbox.And(dbox.In("type", []interface{}{"owner", "finance"}...), dbox.In("updated_date", dateUpdated...), dbox.Ne("hasopen_by", k.Session("username").(string)))).Cursor(nil)
	if !tk.IsNilOrEmpty(err) {
		return err, nil
	}
	tmp := []BusinessMetricsNotificationModel{}
	err = crs.Fetch(&tmp, 0, false)
	crs.Close()
	if !tk.IsNilOrEmpty(err) {
		return err, nil
	}

	notif := map[string][]BusinessMetricsNotificationModel{}
	tmp1 := []BusinessMetricsNotificationModel{}
	for _, v := range tmp {
		if _, exist := notif[v.BMId+"|"+v.Type]; !exist {
			/*if len(tmp1) > 0 {
				tmp1 = []BusinessMetricsNotificationModel{}
			}*/
			notif[v.BMId+"|"+v.Type] = append(tmp1, v)
		} else {
			notif[v.BMId+"|"+v.Type] = append(notif[v.BMId], v)
		}
	}
	return nil, notif
}

func (c *StagingAreaController) GetDataList(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Year    int
		To      string
		From    string
		BMID    string
		Country []interface{}
		Region  []interface{}
		Search  string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	StagingAreaAccess := c.GetAccess(k, "STAGINGAREA")
	StagingArea := new(AccessibilityModel)
	err = tk.MtoStruct(StagingAreaAccess, &StagingArea)

	RegionList := make([]RegionModel, 0)
	csr, err := c.Ctx.Connection.NewQuery().From(new(RegionModel).TableName()).Order("Major_Region", "Country").Cursor(nil)
	err = csr.Fetch(&RegionList, 0, false)
	csr.Close()

	BusinessDriverL1 := make([]BusinessDriverL1Model, 0)
	csr, err = c.Ctx.Find(new(BusinessDriverL1Model), tk.M{}.Set("where", dbox.Eq("businessmetric.id", parm.BMID)))
	err = csr.Fetch(&BusinessDriverL1, 0, false)
	csr.Close()

	// Get Selected Period
	var (
		EndingPeriod time.Time
	)
	startingMonth, endingMonth := 1, 12
	if parm.From == "" && parm.To == "" {
		parm.Year = time.Now().Year()
	} else if parm.To != "" {
		period, _ := time.Parse("200601", parm.To)
		parm.Year = period.Year()
		endingMonth = int(period.Month())
		if parm.From != "" {
			period, _ := time.Parse("200601", parm.From)
			startingMonth = int(period.Month())
		}
	} else if parm.From != "" {
		period, _ := time.Parse("200601", parm.From)
		parm.Year = period.Year()
		startingMonth = int(period.Month())
	}

	// Get Selected Period
	year := strconv.Itoa(parm.Year)
	selectedPeriod := []tk.M{}

	for i := int(startingMonth); i <= endingMonth; i++ {
		month := strconv.Itoa(i)
		tempPeriod := year
		if i < 10 {
			tempPeriod += "0" + month + "01"
		} else {
			tempPeriod += month + "01"
		}

		period, _ := time.Parse("20060102", tempPeriod)
		data := tk.M{}
		data.Set("Title", period.Format("Jan-2006"))
		data.Set("Period", period.Format("Jan2006"))
		selectedPeriod = append(selectedPeriod, data)
		EndingPeriod = period
	}

	result := []tk.M{}
	for _, bd := range BusinessDriverL1[0].BusinessMetric {
		if parm.BMID != "" && parm.BMID != bd.Id {
			continue
		}
		if parm.Search != "" && parm.BMID == "" && !strings.Contains(strings.ToLower(bd.DataPoint), parm.Search) {
			continue
		}

		dataList := []tk.M{}
		if !checkGlobalValue(parm.Region) && len(parm.Region) == 0 && len(parm.Country) == 0 {
			// if no param
			// or all param active
			for ri, r := range RegionList {
				if StagingArea.Global.Read {
					if ri == 0 {
						data := tk.M{}
						data.Set("CountryName", "GLOBAL")
						data.Set("CountryCode", "GLOBAL")
						c.getBMData(data, parm.Year, "GLOBAL", bd.Id, EndingPeriod)
						dataList = append(dataList, data)
					}
				}
				if StagingArea.Country.Read {
					data := tk.M{}
					data.Set("CountryName", r.Country)
					data.Set("CountryCode", r.CountryCode)
					c.getBMData(data, parm.Year, r.CountryCode, bd.Id, EndingPeriod)
					dataList = append(dataList, data)
				}
				if StagingArea.Region.Read {
					if ri == len(RegionList)-1 || (ri > 0 && r.Major_Region != RegionList[ri+1].Major_Region) {
						data := tk.M{}
						data.Set("CountryName", "TOTAL "+r.Major_Region)
						data.Set("CountryCode", r.Major_Region)
						c.getBMData(data, parm.Year, r.Major_Region, bd.Id, EndingPeriod)
						dataList = append(dataList, data)
					}
				}
			}

		} else {
			if StagingArea.Global.Read {
				for _, r := range parm.Region {
					rtmp, _ := r.(string)
					if rtmp == "GLOBAL" {
						data := tk.M{}
						data.Set("CountryName", "GLOBAL")
						data.Set("CountryCode", "GLOBAL")
						c.getBMData(data, parm.Year, "GLOBAL", bd.Id, EndingPeriod)
						dataList = append(dataList, data)
					}
				}
			}
			if StagingArea.Country.Read {
				for _, r := range parm.Country {
					rtmp, _ := r.(string)
					data := tk.M{}
					for _, rr := range RegionList {
						if rr.CountryCode == rtmp {
							data.Set("CountryName", rr.Country)
						}
					}
					data.Set("CountryCode", rtmp)
					c.getBMData(data, parm.Year, rtmp, bd.Id, EndingPeriod)
					dataList = append(dataList, data)
				}
			}
			if StagingArea.Region.Read {
				for _, r := range parm.Region {
					rtmp, _ := r.(string)

					if rtmp != "GLOBAL" {
						data := tk.M{}
						data.Set("CountryName", "TOTAL "+rtmp)
						data.Set("CountryCode", rtmp)
						c.getBMData(data, parm.Year, rtmp, bd.Id, EndingPeriod)
						dataList = append(dataList, data)
					}
				}
			}
		}

		// bmData.Set("DataList", dataList)
		result = append(result, dataList...)
	}
	// return c.SetResultInfo(false, "", tk.M{}.Set("SelectedData", result).Set("SelectedPeriod", selectedPeriod))
	return ResultInfo{Data: result}
}

func (c *StagingAreaController) Save(k *knot.WebContext) interface{} {
	c.LoadBase(k)

	k.Config.OutputType = knot.OutputJson
	type Logs struct {
		NewValue    string
		OldValue    string
		Whatchanged string
	}
	parm := struct {
		Year      int
		To        string
		From      string
		BMID      string
		Country   []interface{}
		Region    []interface{}
		Search    string
		BMList    []tk.M
		LogAction []Logs
	}{}
	err := k.GetPayload(&parm)
	result := tk.M{}
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	for _, la := range parm.LogAction {
		c.Action(k, "Metric Validation", "Save updated data on Metric Validation", la.Whatchanged, la.OldValue, la.NewValue, "", "")
	}

	startingMonth, endingMonth := 1, 12
	if parm.From == "" && parm.To == "" {
		parm.Year = time.Now().Year()
	} else if parm.To != "" {
		period, _ := time.Parse("200601", parm.To)
		parm.Year = period.Year()
		endingMonth = int(period.Month()) + 1
		if parm.From != "" {
			period, _ := time.Parse("200601", parm.From)
			startingMonth = int(period.Month())
		}
	} else if parm.From != "" {
		period, _ := time.Parse("200601", parm.From)
		parm.Year = period.Year()
		startingMonth = int(period.Month())
	}

	////get country
	regionCol := make([]tk.M, 0)
	getMajorRegion := map[string]string{}
	getRegion := map[string]string{}
	getCountry := map[string]string{}
	csr, e := c.Ctx.Connection.NewQuery().From("Region").Cursor(nil)
	e = csr.Fetch(&regionCol, 0, true)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	for _, i := range regionCol {
		getCountry[i.GetString("CountryCode")] = i.GetString("Country")
		getRegion[i.GetString("Country")] = i.GetString("Region")
		getMajorRegion[i.GetString("Country")] = i.GetString("Major_Region")
	}
	////get country

	year := strconv.Itoa(parm.Year)
	startingPeriod := time.Date(parm.Year, time.Month(startingMonth), 1, 0, 0, 0, 0, time.UTC)
	endingPeriod := time.Date(parm.Year, time.Month(endingMonth), 1, 0, 0, 0, 0, time.UTC)

	for _, bm := range parm.BMList {
		bmID := bm.GetString("BMId")
		query := []*dbox.Filter{}
		query = append(query, dbox.Gte("period", startingPeriod.UTC()))
		// if !tk.IsNilOrEmpty(parm.From) {
		// query = append(query, dbox.Gte("period", startingPeriod.UTC()))
		// }
		if tk.IsNilOrEmpty(parm.To) {
			BMData := make([]BusinessMetricsDataTempModel, 0)
			csr, err := c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataTempModel).TableName()).Where(dbox.Eq("bmid", bmID)).Order("-period").Take(1).Cursor(nil)
			err = csr.Fetch(&BMData, 0, false)
			csr.Close()
			if err != nil {
				BMData = make([]BusinessMetricsDataTempModel, 0)
				endingPeriod = BMData[0].Period
			}
			query = append(query, dbox.Lte("period", endingPeriod.UTC()))
		} else {
			query = append(query, dbox.Lte("period", endingPeriod.UTC()))
		}

		DataList := bm.Get("DataList").([]interface{})
		// query := []*dbox.Filter{}
		// query = append(query, dbox.Lte("period", endingPeriod.UTC()))

		query = append(query, dbox.Gte("period", startingPeriod.UTC()))
		query = append(query, dbox.Lte("period", endingPeriod.UTC()))

		query = append(query, dbox.Eq("bmid", bmID))
		BMData := make([]BusinessMetricsDataTempModel, 0)
		csr, err := c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataTempModel).TableName()).Where(dbox.And(query...)).Cursor(nil)
		err = csr.Fetch(&BMData, 0, false)
		csr.Close()
		if err != nil {
			BMData = make([]BusinessMetricsDataTempModel, 0)
		}

		periodDirty := []time.Time{}
		for _, d := range DataList {
			data := d.(map[string]interface{})
			CountryCode := data["CountryCode"].(string)
			BaseLine, FullYearForecast, Target := 0.0, 0.0, 0.0
			// BaseLine, FullYearForecast, Target := 0.0, 0.0, 0.0

			// Get BaseLine Value
			vBaselinetype := tk.Sprintf("%T", data["BaseLine"])
			if vBaselinetype == "string" {
				tempValue := data["BaseLine"].(string)
				tempValue = strings.Trim(tempValue, " ")
				if tempValue == "" {
					BaseLine = 130895111188
				} else {
					BaseLine, err = strconv.ParseFloat(tempValue, 64)
				}
				if err != nil {
					return c.ErrorResultInfo(err.Error(), nil)
				}
			} else {
				BaseLine = data["BaseLine"].(float64)
			}

			// Get FullYearForecast Value
			vForecasttype := tk.Sprintf("%T", data["FullYearForecast"])
			if vForecasttype == "string" {
				tempValue := data["FullYearForecast"].(string)
				tempValue = strings.Trim(tempValue, " ")
				if tempValue == "" {
					FullYearForecast = 130895111188
				} else {
					FullYearForecast, err = strconv.ParseFloat(tempValue, 64)
				}
				if err != nil {
					return c.ErrorResultInfo(err.Error(), nil)
				}
			} else {
				FullYearForecast = data["FullYearForecast"].(float64)
			}

			// Get Target Value
			vTargettype := tk.Sprintf("%T", data["Target"])
			if vTargettype == "string" {
				tempValue := data["Target"].(string)
				tempValue = strings.Trim(tempValue, " ")
				if tempValue == "" {
					Target = 130895111188
				} else {
					Target, err = strconv.ParseFloat(tempValue, 64)
				}
				if err != nil {
					return c.ErrorResultInfo(err.Error(), nil)
				}
			} else {
				Target = data["Target"].(float64)
			}

			for _, ed := range BMData {
				if CountryCode != ed.CountryCode {
					continue
				}

				for i := int(startingMonth); i <= endingMonth; i++ {
					month := strconv.Itoa(i)
					tempPeriod := year
					if i < 10 {
						tempPeriod += "0" + month + "01"
					} else {
						tempPeriod += month + "01"
					}
					period, _ := time.Parse("20060102", tempPeriod)

					vtype := tk.Sprintf("%T", data[period.Format("Jan2006")])
					NewRAG := data["RAG"+period.Format("Jan2006")].(string)
					PrevRAG := data["PrevRAG"+period.Format("Jan2006")].(string)
					btype := tk.Sprintf("%T", data["Budget"+period.Format("Jan2006")])
					if ed.Period.UTC() != period {
						continue
					}

					if vtype == "string" || btype == "string" || NewRAG != PrevRAG || vBaselinetype == "string" || vForecasttype == "string" || vTargettype == "string" {
						if ed.CountryCode == "GLOBAL" {
							// tk.Println("Masuk")
						}
						if vtype == "string" {
							tempValue := data[period.Format("Jan2006")].(string)
							tempValue = strings.Trim(tempValue, " ")
							if tempValue == "" {
								ed.ActualYTD = 130895111188
								ed.NAActual = true
							} else {
								ed.ActualYTD, err = strconv.ParseFloat(tempValue, 64)
								ed.NAActual = false
							}
						}
						if btype == "string" {
							tempValue := data["Budget"+period.Format("Jan2006")].(string)
							tempValue = strings.Trim(tempValue, " ")
							if tempValue == "" {
								ed.Budget = 130895111188
								ed.NABudget = true
							} else {
								ed.Budget, err = strconv.ParseFloat(tempValue, 64)
								ed.NABudget = false
							}
						}

						if NewRAG != PrevRAG {
							tempValue := NewRAG
							tempValue = strings.Trim(tempValue, " ")
							ed.RAG = tempValue
						}

					}
					// tk.Println(ed.Country)
					if strings.Contains(ed.Country, "TOTAL") {
						ed.Country = strings.Replace(ed.Country, "TOTAL ", "", -1)
					}
					ed.Baseline = BaseLine
					ed.FullYearForecast = FullYearForecast
					ed.Target = Target
					ed.UpdatedDate = time.Now()

					err = c.Ctx.Save(&ed)
					if err != nil {
						return c.ErrorResultInfo(err.Error(), nil)
					}
				}

			}

			//cek data baru
			for i := int(startingMonth); i <= endingMonth; i++ {
				month := strconv.Itoa(i)
				tempPeriod := year
				if i < 10 {
					tempPeriod += "0" + month + "01"
				} else {
					tempPeriod += month + "01"
				}

				period, _ := time.Parse("20060102", tempPeriod)
				vtype := tk.Sprintf("%T", data[period.Format("Jan2006")])
				btype := tk.Sprintf("%T", data["Budget"+period.Format("Jan2006")])
				baselinetype := tk.Sprintf("%T", data["BaseLine"])
				forcasttype := tk.Sprintf("%T", data["FullYearForecast"])
				targettype := tk.Sprintf("%T", data["Target"])

				if (data[period.Format("Jan2006")] != "" && vtype == "string") || (data["Budget"+period.Format("Jan2006")] != "" && btype == "string") || (data["BaseLine"] != "" && baselinetype == "string") || (data["FullYearForecast"] != "" && forcasttype == "string") || (data["Target"] != "" && targettype == "string") || data["RAG"+period.Format("Jan2006")] != "" {
					// tk.Printfn("baseline:%s, forecast:%s, target:%s, countryname:%s, tanggal:%s", data["BaseLine"], data["FullYearForecast"], data["Target"], data["CountryName"], period.Format("Jan2006"))
					flag := false
					for _, v := range BMData {
						if v.Period.UTC() == period && v.CountryCode == data["CountryCode"].(string) {
							flag = true
						}
					}
					tk.Println(flag, period)
					if !flag {
						periodDirty = append(periodDirty, period)
					}
				}
			}
		}

		if len(periodDirty) > 0 {
			periods := c.removeDuplicates(periodDirty)
			err = c.checkNewData(periods, DataList, bm, getRegion, getMajorRegion, k)
			if !tk.IsNilOrEmpty(err) {
				return c.ErrorResultInfo(err.Error(), nil)
			}
		}

	}

	return c.SetResultInfo(false, "", result)
}

func (c *StagingAreaController) checkNewData(periods []time.Time, datalist []interface{}, bm tk.M, getRegion, getMajorRegion map[string]string, k *knot.WebContext) error {
	for _, v := range periods {
		for _, d := range datalist {
			tmpDate := v.Format("Jan2006")
			data := d.(map[string]interface{})
			BaseLine, FullYearForecast, Target := 0.0, 0.0, 0.0
			var err error
			// Get BaseLine Value
			vBaselinetype := tk.Sprintf("%T", data["BaseLine"])
			if vBaselinetype == "string" {
				tempValue := data["BaseLine"].(string)
				tempValue = strings.Trim(tempValue, " ")
				if tempValue == "" {
					BaseLine = 130895111188
				} else {
					BaseLine, err = strconv.ParseFloat(tempValue, 64)
				}
				if err != nil {
					return err
				}
			} else {
				BaseLine = data["BaseLine"].(float64)
			}

			// Get FullYearForecast Value
			vForecasttype := tk.Sprintf("%T", data["FullYearForecast"])
			if vForecasttype == "string" {
				tempValue := data["FullYearForecast"].(string)
				tempValue = strings.Trim(tempValue, " ")
				if tempValue == "" {
					FullYearForecast = 130895111188
				} else {
					FullYearForecast, err = strconv.ParseFloat(tempValue, 64)
				}
				if err != nil {
					return err
				}
			} else {
				FullYearForecast = data["FullYearForecast"].(float64)
			}

			// Get Target Value
			vTargettype := tk.Sprintf("%T", data["Target"])
			if vTargettype == "string" {
				tempValue := data["Target"].(string)
				tempValue = strings.Trim(tempValue, " ")
				if tempValue == "" {
					Target = 130895111188
				} else {
					Target, err = strconv.ParseFloat(tempValue, 64)
				}
				if err != nil {
					return err
				}
			} else {
				Target = data["Target"].(float64)
			}

			////get masterBD
			masterBD := new(BusinessDriverL1Model)
			csr, err := c.Ctx.Find(new(BusinessDriverL1Model), tk.M{}.Set("where", dbox.Eq("businessmetric.id", bm.GetString("BMId")))) //
			err = csr.Fetch(&masterBD, 1, false)
			csr.Close()
			if err != nil {
				return err
			}
			if data["Budget"+v.Format("Jan2006")] == data["PrevBudget"+v.Format("Jan2006")] && data[v.Format("Jan2006")] == data["Prev"+v.Format("Jan2006")] && data["RAG"+v.Format("Jan2006")] == data["PrevRAG"+v.Format("Jan2006")] && data["BaseLine"] == "" && data["FullYearForecast"] == "" && data["Target"] == "" {
				continue
			}
			////get masterBD
			o := BusinessMetricsDataTempModel{}
			o.Id = bson.NewObjectId()
			err = c.saveNewData(o, v, tmpDate, masterBD, bm, data, getRegion, getMajorRegion, BaseLine, FullYearForecast, Target, k)
			if !tk.IsNilOrEmpty(err) {
				return err
			}
		}
	}
	return nil
}

func (c *StagingAreaController) saveNewData(o BusinessMetricsDataTempModel, v time.Time, tmpDate string, masterBD *BusinessDriverL1Model, bm tk.M, data map[string]interface{},
	getRegion, getMajorRegion map[string]string, BaseLine, FullYearForecast, Target float64, k *knot.WebContext) error {
	o.Period = v
	o.Year = tk.ToInt(v.Format("2006"), tk.RoundingAuto)
	o.BusinessName = "COMMERCIAL BANKING"
	o.SCId = masterBD.Id
	o.ScorecardCategory = "COMMERCIAL BANKING"
	o.BMId = bm.GetString("BMId")
	o.BusinessMetric = bm.GetString("DataPoint")
	o.BusinessMetricDescription = bm.GetString("Description")

	o.Country = data["CountryCode"].(string)
	o.CountryCode = data["CountryCode"].(string)
	o.MajorRegion = data["CountryCode"].(string)
	o.Region = data["CountryCode"].(string)
	if data["CountryCode"].(string) == "GLOBAL" {
		o.MajorRegion = "GLOBAL"
		o.Region = "GLOBAL"
		o.Country = "GLOBAL"

	} else {
		reg := getRegion[data["CountryName"].(string)]
		if reg != "" {
			o.MajorRegion = getMajorRegion[data["CountryName"].(string)]
			o.Region = reg
			o.Country = data["CountryName"].(string)
		}
	}

	o.Baseline = BaseLine
	o.Actual = 130895111188
	o.Budget = 130895111188
	o.ActualYTD = 130895111188
	o.FullYearForecast = FullYearForecast
	o.Target = Target
	o.CreatedDate = time.Now()
	o.CreatedBy = k.Session("username").(string)
	o.UpdatedDate = time.Now()
	o.UpdatedBy = k.Session("username").(string)
	o.IsCurrent = 0
	o.SourceUID = ""
	o.Source = ""
	o.NABaseline = false
	o.NAActual = false
	o.NABudget = false
	o.NATarget = false
	for dname, _ := range data {
		if strings.Contains(dname, tmpDate) {
			o.RAG = data["RAG"+v.Format("Jan2006")].(string)
			actualYTD := tk.ToFloat64(data[v.Format("Jan2006")], 2, tk.RoundingAuto)
			if actualYTD > 0 {
				o.ActualYTD = actualYTD
			}

			budget := tk.ToFloat64(data["Budget"+v.Format("Jan2006")], 2, tk.RoundingAuto)
			if budget > 0 {
				o.Budget = budget
			}
		}
	}

	err := c.Ctx.Save(&o)
	if !tk.IsNilOrEmpty(err) {
		return err
	}
	return nil
}

func (c *StagingAreaController) removeDuplicates(elements []time.Time) []time.Time {
	// Use map to record duplicates as we find them.
	encountered := map[time.Time]bool{}
	result := []time.Time{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func (c *StagingAreaController) Publish(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	c.Action(k, "Metric Validation", "Publish to Scorecard", "", "", "", "", "")
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Year                   int
		To                     string
		From                   string
		BMID                   string
		Country                []interface{}
		Region                 []interface{}
		RegionCountryScorecard bool
		Search                 string
		BMList                 []tk.M
	}{}
	err := k.GetPayload(&parm)
	result := tk.M{}
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	startingMonth, endingMonth := 1, 12
	if parm.From == "" && parm.To == "" {
		parm.Year = time.Now().Year()
	} else if parm.To != "" {
		period, _ := time.Parse("200601", parm.To)
		parm.Year = period.Year()
		endingMonth = int(period.Month()) + 1
		if parm.From != "" {
			period, _ := time.Parse("200601", parm.From)
			startingMonth = int(period.Month())
		}
	} else if parm.From != "" {
		period, _ := time.Parse("200601", parm.From)
		parm.Year = period.Year()
		startingMonth = int(period.Month())
	}

	// year := strconv.Itoa(parm.Year)
	startingPeriod := time.Date(parm.Year, time.Month(startingMonth), 1, 0, 0, 0, 0, time.UTC)
	endingPeriod := time.Date(parm.Year, time.Month(endingMonth), 1, 0, 0, 0, 0, time.UTC)

	/*conn, err := PrepareConnection()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}*/

	q := c.Ctx.Connection.NewQuery().SetConfig("multiexec", true).From(NewBusinessMetricsDataModel().TableName()).Save()
	d := c.Ctx.Connection.NewQuery().SetConfig("pooling", true).From(NewBusinessMetricsDataModel().TableName()).Delete()
	for _, bm := range parm.BMList {
		bmID := bm.GetString("BMId")

		query := []*dbox.Filter{}
		query = append(query, dbox.Gte("period", startingPeriod.UTC()))
		query = append(query, dbox.Lte("period", endingPeriod.UTC()))
		query = append(query, dbox.Eq("bmid", bmID))
		BMData := make([]BusinessMetricsDataTempModel, 0)
		csr, err := c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataTempModel).TableName()).Where(dbox.And(query...)).Order("countrycode", "-period").Cursor(nil)
		err = csr.Fetch(&BMData, 0, false)
		csr.Close()

		if err != nil {
			BMData = make([]BusinessMetricsDataTempModel, 0)
		}

		existingBMDataList := []BusinessMetricsDataModel{}
		csr, err = c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(dbox.And(dbox.Eq("bmid", bmID))).Cursor(nil)
		if err != nil {
			existingBMDataList = []BusinessMetricsDataModel{}
		} else {
			csr.Fetch(&existingBMDataList, 0, false)
		}
		csr.Close()
		lastBudget := 0.0
		isFoundLastBudget := false
		tempCountryCode := ""
		for _, v := range BMData {
			/*nd := new(BusinessMetricsDataModel)
			temp := tk.M{}
			err = tk.StructToM(v, &temp)
			if err != nil {
				return c.ErrorResultInfo(err.Error(), nil)
			}
			err = tk.MtoStruct(temp, nd)
			if err != nil {
				return c.ErrorResultInfo(err.Error(), nil)
			}*/
			nd := BusinessMetricsDataModel(v)
			if tempCountryCode != nd.CountryCode {
				isFoundLastBudget = false
			}
			if !isFoundLastBudget && !nd.NABudget && nd.Budget != 130895111188 {
				lastBudget = nd.Budget
				isFoundLastBudget = true
			}

			for _, x := range existingBMDataList {
				if nd.Period.UTC() == x.Period.UTC() && nd.CountryCode == x.CountryCode {
					nd.Id = x.Id
					err = d.Exec(tk.M{"data": nd}) /*c.Ctx.Delete(&nd)*/
					if err != nil {
						return c.ErrorResultInfo(err.Error(), nil)
					}
					break
				}
			}
			budgetValue := 0.0
			if !nd.NABudget && nd.Budget != 130895111188 {
				budgetValue = nd.Budget
			}
			nd.RemainingBudget = lastBudget - budgetValue
			nd.RemainingBudgetOpposite = budgetValue - lastBudget
			// tk.Println(nd.CountryCode, "|", nd.Budget, "|", prevBudget, "|", lastBudget)
			err = q.Exec(tk.M{}.Set("data", nd))
			if err != nil {
				return c.ErrorResultInfo(err.Error(), nil)
			}

			tempCountryCode = nd.CountryCode
			/*err = c.Ctx.Save(nd)
			if err != nil {
				return c.ErrorResultInfo(err.Error(), nil)
			}*/
		}
	}
	q.Close()
	d.Close()
	// conn.Close()

	return c.SetResultInfo(false, "", result)
}

func (c *StagingAreaController) ExportXLS(k *knot.WebContext) interface{} {
	// c.LoadBase(k)
	// k.Config.OutputType = knot.OutputJson
	// parm := struct {
	// 	Year    int
	// 	To      string
	// 	From    string
	// 	BMID    string
	// 	Country string
	// 	Region  string
	// 	Search  string
	// 	BMList  []tk.M
	// }{}
	// err := k.GetPayload(&parm)
	// result := ""
	// if err != nil {
	// 	// return c.SetResultInfo(true, err.Error(), nil)
	// }
	// startingMonth, endingMonth := 1, 12
	// if parm.From == "" && parm.To == "" {
	// 	parm.Year = time.Now().Year()
	// } else if parm.To != "" {
	// 	period, _ := time.Parse("200601", parm.To)
	// 	parm.Year = period.Year()
	// 	endingMonth = int(period.Month()) + 1
	// 	if parm.From != "" {
	// 		period, _ := time.Parse("200601", parm.From)
	// 		startingMonth = int(period.Month())
	// 	}
	// } else if parm.From != "" {
	// 	period, _ := time.Parse("200601", parm.From)
	// 	parm.Year = period.Year()
	// 	startingMonth = int(period.Month())
	// }

	// year := strconv.Itoa(parm.Year)
	// startingPeriod := time.Date(parm.Year, time.Month(startingMonth), 1, 0, 0, 0, 0, time.UTC)
	// endingPeriod := time.Date(parm.Year, time.Month(endingMonth), 1, 0, 0, 0, 0, time.UTC)

	// for _, bm := range parm.BMList {
	// 	bmID := bm.GetString("BMId")
	// 	DataList := bm.Get("DataList").([]interface{})
	// 	query := []*dbox.Filter{}
	// 	query = append(query, dbox.Gte("period", startingPeriod.UTC()))
	// 	query = append(query, dbox.Lte("period", endingPeriod.UTC()))
	// 	query = append(query, dbox.Eq("bmid", bmID))
	// 	BMData := make([]BusinessMetricsDataTempModel, 0)
	// 	csr, err := c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataTempModel).TableName()).Where(dbox.And(query...)).Cursor(nil)
	// 	err = csr.Fetch(&BMData, 0, false)
	// 	csr.Close()
	// 	if err != nil {
	// 		BMData = make([]BusinessMetricsDataTempModel, 0)
	// 	}
	// 	for _, d := range DataList {
	// 		data := d.(map[string]interface{})
	// 		CountryCode := data["CountryCode"].(string)
	// 		for _, ed := range BMData {
	// 			if CountryCode != ed.CountryCode {
	// 				continue
	// 			}
	// 			for i := int(startingMonth); i <= endingMonth; i++ {
	// 				month := strconv.Itoa(i)
	// 				tempPeriod := year
	// 				if i < 10 {
	// 					tempPeriod += "0" + month + "01"
	// 				} else {
	// 					tempPeriod += month + "01"
	// 				}
	// 				period, _ := time.Parse("20060102", tempPeriod)

	// 				if ed.Period.UTC() != period {
	// 					continue
	// 				}
	// 				vtype := tk.Sprintf("%T", data[period.Format("Jan2006")])
	// 				if vtype == "string" {
	// 					ed.Actual, err = strconv.ParseFloat(data[period.Format("Jan2006")].(string), 64)
	// 					err = c.Ctx.Save(&ed)
	// 					if err != nil {
	// 						return c.ErrorResultInfo(err.Error(), nil)
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// var (
	// 	EndingPeriod time.Time
	// )
	// tk.Println("Export data....")
	c.LoadBase(k)
	c.Action(k, "Metric Validation", "Save Metric Validation to Excel", "", "", "", "", "")
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Year        int
		To          string
		From        string
		BMID        string
		Country     []string
		Region      []string
		Search      string
		BMList      []interface{}
		OptionValue string
	}{}

	err := k.GetPayload(&parm)
	if err != nil {
		tk.Println(err)
		return c.SetResultInfo(true, err.Error(), nil)
	}
	startingMonth, endingMonth := 1, 12
	if parm.From == "" && parm.To == "" {
		parm.Year = time.Now().Year()
	} else if parm.To != "" {
		period, _ := time.Parse("200601", parm.To)
		parm.Year = period.Year()
		endingMonth = int(period.Month()) + 1
		if parm.From != "" {
			period, _ := time.Parse("200601", parm.From)
			startingMonth = int(period.Month())
		}

	} else if parm.From != "" {
		period, _ := time.Parse("200601", parm.From)
		parm.Year = period.Year()
		startingMonth = int(period.Month())

	}

	// Get Selected Period
	year := strconv.Itoa(parm.Year)
	selectedPeriod := []tk.M{}
	for i := int(startingMonth); i <= endingMonth; i++ {
		month := strconv.Itoa(i)
		tempPeriod := year
		if i < 10 {
			tempPeriod += "0" + month + "01"
		} else {
			tempPeriod += month + "01"
		}

		period, _ := time.Parse("20060102", tempPeriod)
		data := tk.M{}
		data.Set("Title", period.Format("Jan-2006"))
		data.Set("Period", period.Format("Jan2006"))
		selectedPeriod = append(selectedPeriod, data)
		// EndingPeriod = period
	}

	// // Get Selected Data
	result := []interface{}{}
	result = parm.BMList

	// BusinessDriverL1 := make([]BusinessDriverL1Model, 0)
	// csr, err := c.Ctx.Find(new(BusinessDriverL1Model), nil)
	// err = csr.Fetch(&BusinessDriverL1, 0, false)
	// csr.Close()

	// RegionList := make([]RegionModel, 0)
	// query := []*dbox.Filter{}
	// if parm.Region != "" {
	// 	query = append(query, dbox.Eq("Major_Region", parm.Region))
	// }
	// if parm.Country != "" {
	// 	query = append(query, dbox.Eq("CountryCode", parm.Country))
	// }
	// if len(query) > 0 {
	// 	csr, err = c.Ctx.Connection.NewQuery().From(new(RegionModel).TableName()).Where(dbox.And(query...)).Order("Major_Region", "Country").Cursor(nil)
	// } else {
	// 	csr, err = c.Ctx.Connection.NewQuery().From(new(RegionModel).TableName()).Order("Major_Region", "Country").Cursor(nil)
	// }
	// err = csr.Fetch(&RegionList, 0, false)
	// csr.Close()
	// for _, bd := range BusinessDriverL1 {
	// 	for _, bm := range bd.BusinessMetric {
	// 		if parm.BMID != "" && parm.BMID != bm.Id {
	// 			continue
	// 		}
	// 		if parm.Search != "" && parm.BMID == "" && !strings.Contains(strings.ToLower(bm.DataPoint), parm.Search) {
	// 			continue
	// 		}
	// 		bmData := tk.M{}
	// 		bmData.Set("BMId", bm.Id)
	// 		bmData.Set("DataPoint", bm.DataPoint)
	// 		bmData.Set("Description", bm.Description)
	// 		dataList := []tk.M{}
	// 		for ri, r := range RegionList {
	// 			// Add Global Value
	// 			if ri == 0 && parm.Region == "" && parm.Country == "" {
	// 				data := tk.M{}
	// 				data.Set("CountryName", "GLOBAL")
	// 				data.Set("CountryCode", "GLOBAL")
	// 				c.getBMData(data, parm.Year, "GLOBAL", bm.Id, EndingPeriod)
	// 				dataList = append(dataList, data)
	// 			}
	// 			data := tk.M{}
	// 			data.Set("CountryName", r.Country)
	// 			data.Set("CountryCode", r.CountryCode)
	// 			c.getBMData(data, parm.Year, r.CountryCode, bm.Id, EndingPeriod)
	// 			dataList = append(dataList, data)
	// 			if parm.Country != "" {
	// 				continue
	// 			}
	// 			if ri == len(RegionList)-1 || (ri > 0 && r.Major_Region != RegionList[ri+1].Major_Region) {
	// 				data := tk.M{}
	// 				data.Set("CountryName", "TOTAL "+r.Major_Region)
	// 				data.Set("CountryCode", r.Major_Region)
	// 				c.getBMData(data, parm.Year, r.Major_Region, bm.Id, EndingPeriod)
	// 				dataList = append(dataList, data)
	// 			}
	// 		}

	// 		bmData.Set("DataList", dataList)
	// 		result = append(result, bmData)
	// 	}
	// }
	// if err != nil {
	// 	return c.SetResultInfo(true, err.Error(), nil)
	// }

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	style := xlsx.NewStyle()
	font := xlsx.NewFont(12, "calibri")
	border := xlsx.NewBorder("thin", "thin", "thin", "thin")
	style.Font = *font
	style.Border = *border

	styleHeader := xlsx.NewStyle()
	styleHeader.Font = *font
	styleHeader.Font.Bold = true
	styleHeader.Font.Color = "FFFFFFFF"
	styleHeader.Alignment.WrapText = true
	styleHeader.Fill = *xlsx.NewFill("solid", "00516380", "00516380")
	styleHeader.Border = *border

	styleBD := xlsx.NewStyle()
	styleBD.Font = *font
	styleBD.Font.Bold = true
	styleBD.Font.Color = "FFFFFFFF"
	styleBD.Fill = *xlsx.NewFill("solid", "00516380", "00516380")
	styleBD.Border = *border

	styleCountryRegion := xlsx.NewStyle()
	styleCountryRegion.Font = *font
	styleCountryRegion.Font.Bold = true
	styleCountryRegion.Border = *border

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		c.ErrorResultInfo(err.Error(), nil)
	}
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "Metrics"
	cell.SetStyle(styleHeader)

	cell = row.AddCell()
	cell.Value = "Baseline 2017"
	cell.SetStyle(styleHeader)

	for _, bd := range selectedPeriod {
		cell = row.AddCell()
		cell.Value = bd.GetString("Title")
		cell.SetStyle(styleHeader)
		// tk.Println("selectedPeriod", selectedPeriod)
	}
	cell = row.AddCell()
	cell.Value = "YTD Actual 2017"
	cell.SetStyle(styleHeader)
	cell = row.AddCell()
	cell.Value = "Full Year Forecast 2017"
	cell.SetStyle(styleHeader)
	cell = row.AddCell()
	cell.Value = "Target 2018"
	cell.SetStyle(styleHeader)

	for _, bdr := range result {
		bd := bdr.(map[string]interface{})
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = bd["DataPoint"].(string)
		cell.Merge(16, 0)
		cell.SetStyle(styleBD)
		for _, bd22 := range bd["DataList"].([]interface{}) {
			bd2 := bd22.(map[string]interface{})
			row = sheet.AddRow()
			cell = row.AddCell()
			cell.Value = bd2["CountryName"].(string)
			if cell.Value == "GLOBAL" || cell.Value == "TOTAL AME" || cell.Value == "TOTAL ASA" || cell.Value == "TOTAL GCNA" {
				cell.SetStyle(styleCountryRegion)
			} else {
				cell.SetStyle(style)
			}

			cell = row.AddCell()

			baselineValue, isOk := bd2["BaseLine"].(float64)
			if baselineValue == 130895111188 || !isOk {
				cell.SetValue("N/A")

			} else {
				cell.SetValue(baselineValue)
			}
			if bd2["CountryName"].(string) == "GLOBAL" || bd2["CountryName"].(string) == "TOTAL AME" || bd2["CountryName"].(string) == "TOTAL ASA" || bd2["CountryName"].(string) == "TOTAL GCNA" {
				cell.SetStyle(styleCountryRegion)
			} else {
				cell.SetStyle(style)
			}

			for _, bd3 := range selectedPeriod {
				val := bd3["Period"].(string)
				cell = row.AddCell()

				periodValue := 0.0
				periodValues, isString := bd2[val].(string)
				// tk.Println(parm.OptionValue)
				if parm.OptionValue == "RAG" {
					styleFillRed := xlsx.NewStyle()
					border := xlsx.NewBorder("thin", "thin", "thin", "thin")
					styleFillRed.Fill = *xlsx.NewFill("solid", "f74e4e", "f74e4e")
					styleFillRed.Border = *border

					styleFillAmber := xlsx.NewStyle()
					styleFillAmber.Fill = *xlsx.NewFill("solid", "FFD24D", "FFD24D")
					styleFillAmber.Border = *border

					styleFillGreen := xlsx.NewStyle()
					styleFillGreen.Fill = *xlsx.NewFill("solid", "6AC17B", "6AC17B")
					styleFillGreen.Border = *border

					styleFillDefault := xlsx.NewStyle()
					styleFillDefault.Fill = *xlsx.NewFill("solid", "ffffff", "ffffff")
					styleFillDefault.Border = *border

					if bd2["RAG"+val].(string) == "red" {
						cell.SetStyle(styleFillRed)
					} else if bd2["RAG"+val].(string) == "amber" {
						cell.SetStyle(styleFillAmber)
					} else if bd2["RAG"+val].(string) == "green" {
						cell.SetStyle(styleFillGreen)
					} else if bd2["RAG"+val].(string) == "" {
						cell.SetStyle(styleFillDefault)
					}

					// tk.Println(color)
				} else if parm.OptionValue == "Budget" {
					if isString {
						periodValue, err = strconv.ParseFloat(periodValues, 64)
						if err != nil {
							// tk.Println(err)
						}
					} else {
						periodValue = bd2["Budget"+val].(float64)
					}

					if periodValue == 130895111188 || isString {
						cell.SetValue("N/A") //N/A
					} else {
						cell.SetValue(periodValue)
					}
					if bd2["CountryName"].(string) == "GLOBAL" || bd2["CountryName"].(string) == "TOTAL AME" || bd2["CountryName"].(string) == "TOTAL ASA" || bd2["CountryName"].(string) == "TOTAL GCNA" {
						cell.SetStyle(styleCountryRegion)
					} else {
						cell.SetStyle(style)
					}
				} else {
					if isString {
						periodValue, err = strconv.ParseFloat(periodValues, 64)
						// tk.Println("periodValue ", periodValue)
						if err != nil {
							// tk.Println(err)
						}
					} else {
						periodValue = bd2[val].(float64)
					}

					if periodValue == 130895111188 || isString {
						cell.SetValue("N/A")

					} else {
						cell.SetValue(periodValue)
					}
					if bd2["CountryName"].(string) == "GLOBAL" || bd2["CountryName"].(string) == "TOTAL AME" || bd2["CountryName"].(string) == "TOTAL ASA" || bd2["CountryName"].(string) == "TOTAL GCNA" {
						cell.SetStyle(styleCountryRegion)
					} else {
						cell.SetStyle(style)
					}
				}

			}

			cell = row.AddCell()
			ytdValue, isOk := bd2["YTDActual"].(float64)
			if ytdValue == 130895111188 || !isOk {
				cell.SetValue("N/A")

			} else {
				cell.SetValue(ytdValue)
			}

			if bd2["CountryName"].(string) == "GLOBAL" || bd2["CountryName"].(string) == "TOTAL AME" || bd2["CountryName"].(string) == "TOTAL ASA" || bd2["CountryName"].(string) == "TOTAL GCNA" {
				cell.SetStyle(styleCountryRegion)
			} else {
				cell.SetStyle(style)
			}
			cell = row.AddCell()
			forecastValue, isOk := bd2["FullYearForecast"].(float64)
			if forecastValue == 130895111188 || !isOk {
				cell.SetValue("N/A")

			} else {
				cell.SetValue(forecastValue)
			}

			if bd2["CountryName"].(string) == "GLOBAL" || bd2["CountryName"].(string) == "TOTAL AME" || bd2["CountryName"].(string) == "TOTAL ASA" || bd2["CountryName"].(string) == "TOTAL GCNA" {
				cell.SetStyle(styleCountryRegion)
			} else {
				cell.SetStyle(style)
			}
			cell = row.AddCell()
			targetValue, isOk := bd2["Target"].(float64)
			if targetValue == 130895111188 || !isOk {
				cell.SetValue("N/A")

			} else {
				cell.SetValue(targetValue)
				// tk.Println("targetValue: ", targetValue)
			}
			if bd2["CountryName"].(string) == "GLOBAL" || bd2["CountryName"].(string) == "TOTAL AME" || bd2["CountryName"].(string) == "TOTAL ASA" || bd2["CountryName"].(string) == "TOTAL GCNA" {
				cell.SetStyle(styleCountryRegion)
			} else {
				cell.SetStyle(style)
			}
		}
	}

	t := time.Now().UTC()
	times := t.Format("20060102_150405")
	err = file.Save("bef/assets/download/stagingarea_" + times + ".xlsx")
	// tk.Println("Excel Exported ")
	if err != nil {
		tk.Println(err)
		c.ErrorResultInfo(err.Error(), nil)
	}

	// fullresult := tk.M{}.Set("SelectedData", result).Set("SelectedPeriod", selectedPeriod)

	return c.SetResultInfo(false, "ok", "stagingarea_"+times+".xlsx")
}

func (c *StagingAreaController) UpdateNotification(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		BMId string
		Type string
	}{}
	err := k.GetPayload(&parm)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	pipe := []tk.M{
		tk.M{
			"$match": tk.M{
				"bmid": parm.BMId,
				"type": parm.Type,
			},
		},
		tk.M{
			"$group": tk.M{
				"_id": tk.M{
					"bmid":   "$bmid",
					"bmname": "$bmname",
					"bmtype": "$type",
				},
				"lastUpdatePeriod": tk.M{
					"$last": "$updated_date",
				},
			},
		},
	}
	crs, err := c.Ctx.Connection.NewQuery().Command("pipe", pipe).From(NewBusinessMetricsNotificationModel().TableName()).Cursor(nil)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	dataNotif := tk.Ms{}
	err = crs.Fetch(&dataNotif, 0, false)
	crs.Close()
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	crs, err = c.Ctx.Connection.NewQuery().From(NewBusinessMetricsNotificationModel().TableName()).
		Where(dbox.And(dbox.Eq("bmid", parm.BMId), dbox.Eq("type", parm.Type))).Cursor(nil)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	data := []*BusinessMetricsNotificationModel{}
	err = crs.Fetch(&data, 0, false)
	crs.Close()
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	data1 := BusinessMetricsNotificationModel{}
	for _, v := range dataNotif {
		for _, m := range data {
			// tk.Println(v.Get("lastUpdatePeriod").(time.Time).UTC(), m.Updated_Date.UTC(), v.Get("lastUpdatePeriod").(time.Time).UTC() == m.Updated_Date.UTC())
			if v.Get("lastUpdatePeriod").(time.Time).UTC() == m.Updated_Date.UTC() {
				data1 = *m
			}
		}
	}

	err = c.Ctx.DeleteMany(NewBusinessMetricsNotificationModel(), dbox.And(dbox.Eq("bmid", parm.BMId), dbox.Eq("type", parm.Type)))
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	data1.HasOpen_By = append(data1.HasOpen_By, k.Session("username").(string))
	err = c.Ctx.Save(&data1)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	// tk.Println(tk.JsonString(data1))

	return c.SetResultInfo(false, "", tk.M{})
}
