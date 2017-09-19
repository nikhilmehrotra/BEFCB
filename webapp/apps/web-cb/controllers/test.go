package controllers

import (
	// "eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	// "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	// tk "github.com/eaciit/toolkit"
	"strconv"
	"time"
	// "bytes"
	// "fmt"
	// "io"
	// "os"
	// "path/filepath"
	// "strings"
)

type TestController struct {
	*BaseController
}

func (c *TestController) GetData(k *knot.WebContext) interface{} {
	// c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	result := make([]BusinessDriverL1Model, 0)

	csr, err := c.Ctx.Connection.NewQuery().From("BusinessDriverL1").Order("seq").
		Cursor(nil)
	defer csr.Close()
	err = csr.Fetch(&result, 0, true)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	metricsDatas := []*BusinessMetricsDataModel{}
	for _, res := range result {
		for _, bd := range res.BusinessMetric {
			YTDValue := 0.0
			for _, mm := range bd.ActualData {

				metricData := NewBusinessMetricsDataModel()
				metricData.Period, _ = time.Parse("20060102", mm.PeriodStr)
				// metricData.Period = time.Date(2016, 11, 1, 0, 0, 0, 0, time.Local)
				metricData.Year, _ = strconv.Atoi(metricData.Period.Format("2006"))
				metricData.SCId = res.Id
				metricData.BusinessName = "COMMERCIAL BANKING"
				metricData.ScorecardCategory = res.Name
				metricData.BMId = bd.Id                               //    string
				metricData.BusinessMetric = bd.DataPoint              //    string
				metricData.BusinessMetricDescription = bd.Description //    string
				metricData.MajorRegion = "GLOBAL"
				metricData.Region = "GLOBAL"
				metricData.Country = "GLOBAL"
				metricData.CountryCode = "GLOBAL"
				metricData.Baseline = bd.BaseLineValue

				metricData.Actual = mm.Value
				// YTDValue += mm.Value
				metricData.ActualYTD = YTDValue

				metricData.FullYearForecast = 0
				metricData.Target = bd.TargetValue
				metricData.CreatedDate = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
				metricData.CreatedBy = "eaciit"
				metricData.UpdatedDate = time.Now()
				metricData.UpdatedBy = "eaciit"

				if mm.Flag == "C1" {
					// metricData.IsCurrent = 1
					metricsDatas = append(metricsDatas, metricData)
				}

			}

		}
	}

	for _, gg := range metricsDatas {
		e := c.Ctx.Save(gg)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
	}

	return c.SetResultInfo(false, "", metricsDatas)
}
