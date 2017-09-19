package controllers

import (
	"eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"time"
	// "regexp"
	"bytes"
	"fmt"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	xl "github.com/tealeg/xlsx"
	// "gopkg.in/mgo.v2/bson"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type ScorecardController struct {
	*BaseController
}

func (c *ScorecardController) GetData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	parm := struct {
		Region  string
		Country string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	result := make([]BusinessDriverL1Model, 0)
	resultSC := make([]BusinessDriverL1Model, 0)
	csr, err := c.Ctx.Connection.NewQuery().From("BusinessDriverL1").Order("seq").Cursor(nil)
	err = csr.Fetch(&result, 0, true)
	csr.Close()
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	csr, err = c.Ctx.Connection.NewQuery().From("BusinessDriverL1").Order("seq").Cursor(nil)
	err = csr.Fetch(&resultSC, 0, true)
	csr.Close()
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	if parm.Region == "" || parm.Region == "Region" {
		parm.Region = "GLOBAL"
	}
	if parm.Country == "" || parm.Country == "Country" {
		parm.Country = parm.Region
	}

	// GetRegionList
	RegionalDataList := []string{parm.Region}
	if parm.Country == "GLOBAL" {
		pipe := []tk.M{
			tk.M{}.Set("$group", tk.M{}.Set("_id", "$Major_Region")),
			tk.M{}.Set("$sort", tk.M{}.Set("_id", -1)),
		}
		RegionList := []tk.M{}
		csr, err = c.Ctx.Connection.NewQuery().Command("pipe", pipe).From(new(RegionModel).TableName()).Cursor(nil)
		csr.Fetch(&RegionList, 0, false)
		csr.Close()
		if err != nil {
			RegionList = []tk.M{}
		}

		RegionalDataList = []string{}
		for _, rg := range RegionList {
			RegionalDataList = append(RegionalDataList, rg.GetString("_id"))
		}
	}

	// Country Data
	CountryDataList := []string{}
	if parm.Region != parm.Country {
		CountryDataList = append(CountryDataList, parm.Country)
		CountryData := new(RegionModel)
		csr, err = c.Ctx.Connection.NewQuery().From(new(RegionModel).TableName()).Where(dbox.Eq("Country", parm.Country)).Cursor(nil)
		csr.Fetch(&CountryData, 1, false)
		csr.Close()
		if err != nil {
			RegionalDataList = []string{}
		} else {
			RegionalDataList = []string{CountryData.Major_Region}
		}
		CountryDataList = []string{parm.Country}
	}
	// tk.Println(RegionalDataList)
	// tk.Println(CountryDataList)

	query := []*dbox.Filter{}
	for bd := range result {
		for bm := range result[bd].BusinessMetric {
			if result[bd].BusinessMetric[bm].DecimalFormat == "" {
				result[bd].BusinessMetric[bm].DecimalFormat = "0"
			}
			result[bd].BusinessMetric[bm].Display = ""
			dataList := []BusinessMetricsDataModel{}
			data := BusinessMetricsDataModel{}
			// tk.Println(result[bd].BusinessMetric[bm].Id)
			query = append(query[0:0], dbox.Eq("bmid", result[bd].BusinessMetric[bm].Id))
			if parm.Region == parm.Country {
				query = append(query, dbox.Eq("majorregion", "GLOBAL"))
			}

			query = append(query, dbox.Eq("country", "GLOBAL"))
			csr, err := c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(dbox.And(query...)).Order("period").Cursor(nil)
			err = csr.Fetch(&dataList, 0, false)
			csr.Close()
			if err != nil {
				result[bd].BusinessMetric[bm].CurrentPeriod = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
				result[bd].BusinessMetric[bm].CurrentValue = 0
				result[bd].BusinessMetric[bm].TargetValue = 0
				result[bd].BusinessMetric[bm].BaseLineValue = 0
				result[bd].BusinessMetric[bm].MetricFiles = []MetricFile{}
				continue
			}
			if len(dataList) > 0 {
				isAny := false
				for idx, i := range dataList {
					if i.ActualYTD != 130895111188 && !i.NAActual {
						data = dataList[idx]
						isAny = true
					}
				}
				if !isAny {
					data = dataList[0]
				}
				// tk.Println("ANY : ")
				// tk.Println(isAny)
			}
			temp := result[bd].BusinessMetric[bm].UpdatedDate
			updatedDate, _ := strconv.Atoi(temp.Format("20060102150405"))
			dataUpdatedDate, _ := strconv.Atoi(data.UpdatedDate.Format("20060102150405"))
			if data.Year > 0 && (updatedDate == 10101000000 || dataUpdatedDate > updatedDate) {
				result[bd].BusinessMetric[bm].CurrentYTDValue = data.ActualYTD
				result[bd].BusinessMetric[bm].CurrentValue = data.ActualYTD
				result[bd].BusinessMetric[bm].CurrentPeriod = data.Period
				result[bd].BusinessMetric[bm].CurrentPeriodStr = data.Period.Format("20060102")
				result[bd].BusinessMetric[bm].TargetValue = data.Target
				result[bd].BusinessMetric[bm].TargetPeriod = data.Period
				result[bd].BusinessMetric[bm].TargetPeriodStr = data.Period.Format("20060102")
				result[bd].BusinessMetric[bm].BaseLineValue = data.Baseline
				result[bd].BusinessMetric[bm].BaseLinePeriod = data.Period
				result[bd].BusinessMetric[bm].BaseLinePeriodStr = data.Period.Format("20060102")
				// tk.Println("DataPoint", result[bd].BusinessMetric[bm].DataPoint)
				// tk.Println("decformat", result[bd].BusinessMetric[bm].DecimalFormat)

				if data.Budget == 130895111188 {
					data.Budget = 0
					data.NABudget = true
				}

				if result[bd].BusinessMetric[bm].ValueType == 0 {
					result[bd].BusinessMetric[bm].CurrentYTDValueVsBudget = data.Budget - data.ActualYTD
				} else {
					result[bd].BusinessMetric[bm].CurrentYTDValueVsBudget = data.ActualYTD - data.Budget
				}

				result[bd].BusinessMetric[bm].NABaseline = data.NABaseline
				result[bd].BusinessMetric[bm].NAActual = data.NAActual
				result[bd].BusinessMetric[bm].NATarget = data.NATarget
				result[bd].BusinessMetric[bm].Display = data.RAG

				if result[bd].BusinessMetric[bm].CurrentYTDValue == 130895111188 {
					result[bd].BusinessMetric[bm].CurrentYTDValue = 0
					result[bd].BusinessMetric[bm].NAActual = true
				}
				if result[bd].BusinessMetric[bm].CurrentValue == 130895111188 {
					result[bd].BusinessMetric[bm].CurrentValue = 0
					result[bd].BusinessMetric[bm].NAActual = true
				}
				if result[bd].BusinessMetric[bm].BaseLineValue == 130895111188 {
					result[bd].BusinessMetric[bm].BaseLineValue = 0
					result[bd].BusinessMetric[bm].NABaseline = true
				}
				if result[bd].BusinessMetric[bm].TargetValue == 130895111188 {
					result[bd].BusinessMetric[bm].TargetValue = 0
					result[bd].BusinessMetric[bm].NATarget = true
				}

				if result[bd].BusinessMetric[bm].NAActual {
					data.NABudget = true
				}
				result[bd].BusinessMetric[bm].NABudget = data.NABudget
				// tk.Println("data.ActualYTD : ", data.ActualYTD)
				// tk.Println("data.Budget : ", data.Budget)
				// tk.Println("ActualVsBudget : ", result[bd].BusinessMetric[bm].CurrentYTDValueVsBudget)
				// if len(dataList) > 1 {
				// 	result[bd].BusinessMetric[bm].PreviousValue = dataList[1].ActualYTD
				// 	if result[bd].BusinessMetric[bm].PreviousValue == 130895111188 {
				// 		result[bd].BusinessMetric[bm].PreviousValue = 0
				// 	}
				// }
			} else {
				result[bd].BusinessMetric[bm].CurrentPeriod = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
				result[bd].BusinessMetric[bm].CurrentValue = 0
				result[bd].BusinessMetric[bm].TargetValue = 0
				result[bd].BusinessMetric[bm].BaseLineValue = 0
				result[bd].BusinessMetric[bm].CurrentYTDValueVsBudget = 0
				result[bd].BusinessMetric[bm].CurrentYTDValue = 0
				result[bd].BusinessMetric[bm].NABaseline = true
				result[bd].BusinessMetric[bm].NAActual = true
				result[bd].BusinessMetric[bm].NATarget = true
				result[bd].BusinessMetric[bm].NABudget = true
			}
			//Data Regional
			result[bd].BusinessMetric[bm].RegionalData = []MetricPartialData{}
			for _, rdj := range RegionalDataList {
				ner := MetricPartialData{}
				ner.Name = rdj
				query = append(query[0:0], dbox.Eq("bmid", result[bd].BusinessMetric[bm].Id))
				query = append(query, dbox.Eq("country", ner.Name))
				query = append(query, dbox.Eq("period", data.Period.UTC()))
				dataRegion := BusinessMetricsDataModel{}
				csr, err = c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(dbox.And(query...)).Cursor(nil)
				err = csr.Fetch(&dataRegion, 1, false)
				csr.Close()
				// tk.Println("###")
				// tk.Println(result[bd].BusinessMetric[bm].Id)
				// tk.Println(ner.Name)
				// tk.Println(data.Period.UTC())
				if err != nil {
					// tk.Println("ERROR")
					ner.NAActual = true
					ner.NABudget = true
				} else {
					if dataRegion.ActualYTD == 130895111188 {
						ner.CurrentYTDValue = 0
						ner.NAActual = true
					} else {
						ner.CurrentYTDValue = dataRegion.ActualYTD
						ner.NAActual = dataRegion.NAActual
					}
					if dataRegion.ActualYTD == 130895111188 || dataRegion.Budget == 130895111188 {
						ner.CurrentYTDValueVsBudget = 0
						ner.NABudget = true
					} else {
						if result[bd].BusinessMetric[bm].ValueType == 0 {
							ner.CurrentYTDValueVsBudget = dataRegion.Budget - dataRegion.ActualYTD
						} else {
							ner.CurrentYTDValueVsBudget = dataRegion.ActualYTD - dataRegion.Budget
						}
						// ner.CurrentYTDValueVsBudget = dataRegion.Budget - dataRegion.ActualYTD
						ner.NABudget = dataRegion.NABudget
					}

					ner.Rag = dataRegion.RAG
				}

				result[bd].BusinessMetric[bm].RegionalData = append(result[bd].BusinessMetric[bm].RegionalData, ner)
			}

			//Data Country
			result[bd].BusinessMetric[bm].CountryData = []MetricPartialData{}
			for _, cdl := range CountryDataList {
				her := MetricPartialData{}
				her.Name = cdl
				query = append(query[0:0], dbox.Eq("bmid", result[bd].BusinessMetric[bm].Id))
				query = append(query, dbox.Eq("country", her.Name))
				query = append(query, dbox.Eq("period", data.Period.UTC()))
				dataCountry := BusinessMetricsDataModel{}
				csr, err = c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(dbox.And(query...)).Cursor(nil)
				err = csr.Fetch(&dataCountry, 1, false)
				csr.Close()

				if err != nil {
					her.NABaseline = true
					her.NABudget = true
					her.NAActual = true
					her.NATarget = true
				} else {
					if dataCountry.Baseline == 130895111188 {
						her.Baseline = 0
						her.NABaseline = true
					} else {
						her.Baseline = dataCountry.Baseline
						her.NABaseline = dataCountry.NABaseline
					}
					if dataCountry.ActualYTD == 130895111188 {
						her.CurrentYTDValue = 0
						her.NAActual = true
					} else {
						her.CurrentYTDValue = dataCountry.ActualYTD
						her.NAActual = dataCountry.NAActual
					}
					if dataCountry.ActualYTD == 130895111188 || dataCountry.Budget == 130895111188 {
						her.CurrentYTDValueVsBudget = 0
						her.NABudget = true
					} else {
						if result[bd].BusinessMetric[bm].ValueType == 0 {
							her.CurrentYTDValueVsBudget = dataCountry.Budget - dataCountry.ActualYTD
						} else {
							her.CurrentYTDValueVsBudget = dataCountry.ActualYTD - dataCountry.Budget
						}
						// her.CurrentYTDValueVsBudget = dataCountry.Budget - dataCountry.ActualYTD
						her.NABudget = dataCountry.NABudget
					}
					if dataCountry.Target == 130895111188 {
						her.Target = 0
						her.NATarget = true
					} else {
						her.Target = dataCountry.Target
						her.NATarget = dataCountry.NATarget
					}
					her.Rag = dataCountry.RAG
				}

				result[bd].BusinessMetric[bm].CountryData = append(result[bd].BusinessMetric[bm].CountryData, her)
			}

		}
	}

	query = []*dbox.Filter{}
	for bd := range resultSC {
		for bm := range resultSC[bd].BusinessMetric {
			if resultSC[bd].BusinessMetric[bm].DecimalFormat == "" {
				resultSC[bd].BusinessMetric[bm].DecimalFormat = "0"
			}
			resultSC[bd].BusinessMetric[bm].Display = ""
			dataList := []BusinessMetricsDataModel{}
			data := BusinessMetricsDataModel{}
			// tk.Println(resultSC[bd].BusinessMetric[bm].Id)
			query = append(query[0:0], dbox.Eq("bmid", resultSC[bd].BusinessMetric[bm].Id))
			if parm.Region == parm.Country {
				query = append(query, dbox.Eq("majorregion", "GLOBAL"))
			}

			query = append(query, dbox.Eq("country", parm.Country))
			csr, err := c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(dbox.And(query...)).Order("period").Cursor(nil)
			err = csr.Fetch(&dataList, 0, false)
			csr.Close()
			if err != nil {
				resultSC[bd].BusinessMetric[bm].CurrentPeriod = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
				resultSC[bd].BusinessMetric[bm].CurrentValue = 0
				resultSC[bd].BusinessMetric[bm].TargetValue = 0
				resultSC[bd].BusinessMetric[bm].BaseLineValue = 0
				resultSC[bd].BusinessMetric[bm].MetricFiles = []MetricFile{}
				continue
			}
			if len(dataList) > 0 {
				isAny := false
				for idx, i := range dataList {
					if i.ActualYTD != 130895111188 && !i.NAActual {
						data = dataList[idx]
						isAny = true
					}
				}
				if !isAny {
					data = dataList[0]
				}
				// tk.Println("ANY : ")
				// tk.Println(isAny)
			}
			temp := resultSC[bd].BusinessMetric[bm].UpdatedDate
			updatedDate, _ := strconv.Atoi(temp.Format("20060102150405"))
			dataUpdatedDate, _ := strconv.Atoi(data.UpdatedDate.Format("20060102150405"))
			if data.Year > 0 && (updatedDate == 10101000000 || dataUpdatedDate > updatedDate) {
				resultSC[bd].BusinessMetric[bm].CurrentYTDValue = data.ActualYTD
				resultSC[bd].BusinessMetric[bm].CurrentValue = data.ActualYTD
				resultSC[bd].BusinessMetric[bm].CurrentPeriod = data.Period
				resultSC[bd].BusinessMetric[bm].CurrentPeriodStr = data.Period.Format("20060102")
				resultSC[bd].BusinessMetric[bm].TargetValue = data.Target
				resultSC[bd].BusinessMetric[bm].TargetPeriod = data.Period
				resultSC[bd].BusinessMetric[bm].TargetPeriodStr = data.Period.Format("20060102")
				resultSC[bd].BusinessMetric[bm].BaseLineValue = data.Baseline
				resultSC[bd].BusinessMetric[bm].BaseLinePeriod = data.Period
				resultSC[bd].BusinessMetric[bm].BaseLinePeriodStr = data.Period.Format("20060102")
				// tk.Println("DataPoint", resultSC[bd].BusinessMetric[bm].DataPoint)
				// tk.Println("decformat", resultSC[bd].BusinessMetric[bm].DecimalFormat)

				if data.Budget == 130895111188 {
					data.Budget = 0
					data.NABudget = true
				}

				if resultSC[bd].BusinessMetric[bm].ValueType == 0 {
					resultSC[bd].BusinessMetric[bm].CurrentYTDValueVsBudget = data.Budget - data.ActualYTD
				} else {
					resultSC[bd].BusinessMetric[bm].CurrentYTDValueVsBudget = data.ActualYTD - data.Budget
				}

				resultSC[bd].BusinessMetric[bm].NABaseline = data.NABaseline
				resultSC[bd].BusinessMetric[bm].NAActual = data.NAActual
				resultSC[bd].BusinessMetric[bm].NATarget = data.NATarget
				resultSC[bd].BusinessMetric[bm].Display = data.RAG

				if resultSC[bd].BusinessMetric[bm].CurrentYTDValue == 130895111188 {
					resultSC[bd].BusinessMetric[bm].CurrentYTDValue = 0
					resultSC[bd].BusinessMetric[bm].NAActual = true
				}
				if resultSC[bd].BusinessMetric[bm].CurrentValue == 130895111188 {
					resultSC[bd].BusinessMetric[bm].CurrentValue = 0
					resultSC[bd].BusinessMetric[bm].NAActual = true
				}
				if resultSC[bd].BusinessMetric[bm].BaseLineValue == 130895111188 {
					resultSC[bd].BusinessMetric[bm].BaseLineValue = 0
					resultSC[bd].BusinessMetric[bm].NABaseline = true
				}
				if resultSC[bd].BusinessMetric[bm].TargetValue == 130895111188 {
					resultSC[bd].BusinessMetric[bm].TargetValue = 0
					resultSC[bd].BusinessMetric[bm].NATarget = true
				}

				if resultSC[bd].BusinessMetric[bm].NAActual {
					data.NABudget = true
				}
				resultSC[bd].BusinessMetric[bm].NABudget = data.NABudget
				// tk.Println("data.ActualYTD : ", data.ActualYTD)
				// tk.Println("data.Budget : ", data.Budget)
				// tk.Println("ActualVsBudget : ", resultSC[bd].BusinessMetric[bm].CurrentYTDValueVsBudget)
				// if len(dataList) > 1 {
				// 	resultSC[bd].BusinessMetric[bm].PreviousValue = dataList[1].ActualYTD
				// 	if resultSC[bd].BusinessMetric[bm].PreviousValue == 130895111188 {
				// 		resultSC[bd].BusinessMetric[bm].PreviousValue = 0
				// 	}
				// }
			} else {
				resultSC[bd].BusinessMetric[bm].CurrentPeriod = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
				resultSC[bd].BusinessMetric[bm].CurrentValue = 0
				resultSC[bd].BusinessMetric[bm].TargetValue = 0
				resultSC[bd].BusinessMetric[bm].BaseLineValue = 0
				resultSC[bd].BusinessMetric[bm].CurrentYTDValueVsBudget = 0
				resultSC[bd].BusinessMetric[bm].CurrentYTDValue = 0
				resultSC[bd].BusinessMetric[bm].NABaseline = true
				resultSC[bd].BusinessMetric[bm].NAActual = true
				resultSC[bd].BusinessMetric[bm].NATarget = true
				resultSC[bd].BusinessMetric[bm].NABudget = true
			}
			//Data Regional
			resultSC[bd].BusinessMetric[bm].RegionalData = []MetricPartialData{}
			for _, rdj := range RegionalDataList {
				ner := MetricPartialData{}
				ner.Name = rdj
				query = append(query[0:0], dbox.Eq("bmid", resultSC[bd].BusinessMetric[bm].Id))
				query = append(query, dbox.Eq("country", ner.Name))
				query = append(query, dbox.Eq("period", data.Period.UTC()))
				dataRegion := BusinessMetricsDataModel{}
				csr, err = c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(dbox.And(query...)).Cursor(nil)
				err = csr.Fetch(&dataRegion, 1, false)
				csr.Close()
				// tk.Println("###")
				// tk.Println(resultSC[bd].BusinessMetric[bm].Id)
				// tk.Println(ner.Name)
				// tk.Println(data.Period.UTC())
				if err != nil {
					// tk.Println("ERROR")
					ner.NAActual = true
					ner.NABudget = true
				} else {
					if dataRegion.ActualYTD == 130895111188 {
						ner.CurrentYTDValue = 0
						ner.NAActual = true
					} else {
						ner.CurrentYTDValue = dataRegion.ActualYTD
						ner.NAActual = dataRegion.NAActual
					}
					if dataRegion.ActualYTD == 130895111188 || dataRegion.Budget == 130895111188 {
						ner.CurrentYTDValueVsBudget = 0
						ner.NABudget = true
					} else {
						if resultSC[bd].BusinessMetric[bm].ValueType == 0 {
							ner.CurrentYTDValueVsBudget = dataRegion.Budget - dataRegion.ActualYTD
						} else {
							ner.CurrentYTDValueVsBudget = dataRegion.ActualYTD - dataRegion.Budget
						}
						// ner.CurrentYTDValueVsBudget = dataRegion.Budget - dataRegion.ActualYTD
						ner.NABudget = dataRegion.NABudget
					}

					ner.Rag = dataRegion.RAG
				}

				resultSC[bd].BusinessMetric[bm].RegionalData = append(resultSC[bd].BusinessMetric[bm].RegionalData, ner)
			}

			//Data Country
			resultSC[bd].BusinessMetric[bm].CountryData = []MetricPartialData{}
			for _, cdl := range CountryDataList {
				her := MetricPartialData{}
				her.Name = cdl
				query = append(query[0:0], dbox.Eq("bmid", resultSC[bd].BusinessMetric[bm].Id))
				query = append(query, dbox.Eq("country", her.Name))
				query = append(query, dbox.Eq("period", data.Period.UTC()))
				dataCountry := BusinessMetricsDataModel{}
				csr, err = c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(dbox.And(query...)).Cursor(nil)
				err = csr.Fetch(&dataCountry, 1, false)
				csr.Close()

				if err != nil {
					her.NABaseline = true
					her.NABudget = true
					her.NAActual = true
					her.NATarget = true
				} else {
					if dataCountry.Baseline == 130895111188 {
						her.Baseline = 0
						her.NABaseline = true
					} else {
						her.Baseline = dataCountry.Baseline
						her.NABaseline = dataCountry.NABaseline
					}
					if dataCountry.ActualYTD == 130895111188 {
						her.CurrentYTDValue = 0
						her.NAActual = true
					} else {
						her.CurrentYTDValue = dataCountry.ActualYTD
						her.NAActual = dataCountry.NAActual
					}
					if dataCountry.ActualYTD == 130895111188 || dataCountry.Budget == 130895111188 {
						her.CurrentYTDValueVsBudget = 0
						her.NABudget = true
					} else {
						if resultSC[bd].BusinessMetric[bm].ValueType == 0 {
							her.CurrentYTDValueVsBudget = dataCountry.Budget - dataCountry.ActualYTD
						} else {
							her.CurrentYTDValueVsBudget = dataCountry.ActualYTD - dataCountry.Budget
						}
						// her.CurrentYTDValueVsBudget = dataCountry.Budget - dataCountry.ActualYTD
						her.NABudget = dataCountry.NABudget
					}
					if dataCountry.Target == 130895111188 {
						her.Target = 0
						her.NATarget = true
					} else {
						her.Target = dataCountry.Target
						her.NATarget = dataCountry.NATarget
					}
					her.Rag = dataCountry.RAG
				}

				resultSC[bd].BusinessMetric[bm].CountryData = append(resultSC[bd].BusinessMetric[bm].CountryData, her)
			}

		}
	}

	res := []tk.M{}
	for i,x:= range resultSC {
		temp := tk.M{} 
		err = tk.StructToM(x, &temp)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
		mbr := tk.M{}
		err = tk.StructToM(result[i], &mbr)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
		temp.Set("BusinessMetricMBR",mbr.Get("BusinessMetric"))
		
		res = append(res,temp)
	}
	return c.SetResultInfo(false, "", res)
}

func (c *ScorecardController) SaveOrderData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	c.Action(k, "Scorecard", "Order Key Metrics", "", "", "", "", "")
	k.Config.OutputType = knot.OutputJson

	payload := []struct {
		SectionIndex int
		Description  string
		OrderIndex   int
	}{}

	err := k.GetPayload(&payload)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	result := make([]*BusinessDriverL1Model, 0)
	csr, err := c.Ctx.Connection.NewQuery().From("BusinessDriverL1").Order("seq").Cursor(nil)
	defer csr.Close()
	err = csr.Fetch(&result, 0, true)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	for i := range result {
		isChangesHappen := false

		for j, eachMetric := range result[i].BusinessMetric {
			for _, eachIndex := range payload {

				if i == eachIndex.SectionIndex && eachMetric.Description == eachIndex.Description {
					isChangesHappen = true
					result[i].BusinessMetric[j].OrderIndex = eachIndex.OrderIndex
				}
			}
		}

		if isChangesHappen {
			err = c.Ctx.Save(result[i])
			if err != nil {
				return c.SetResultInfo(true, err.Error(), nil)
			}
		}
	}

	return c.SetResultInfo(false, "", result)
}

func (c *ScorecardController) SaveBD1(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		BusinessDriver BusinessDriverL1Model
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	d := new(BusinessDriverL1Model)
	csr, e := c.Ctx.Find(new(BusinessDriverL1Model), tk.M{}.Set("where", dbox.Eq("_id", parm.BusinessDriver.Id)))
	e = csr.Fetch(&d, 1, false)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	csr.Close()
	d.Description = parm.BusinessDriver.Description

	d.UpdatedDate = time.Now()
	e = c.Ctx.Save(d)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	return c.SetResultInfo(false, "", nil)
}

func (c *ScorecardController) SavePosition(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		BusinessDriver BusinessDriverL1Model
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	d := new(BusinessDriverL1Model)
	csr, e := c.Ctx.Find(new(BusinessDriverL1Model), tk.M{}.Set("where", dbox.Eq("_id", parm.BusinessDriver.Id)))
	e = csr.Fetch(&d, 1, false)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	csr.Close()
	bm := []BusinessMetric{}
	for idx, i := range d.BusinessMetric {
		changes := ""
		for tempidx, temp := range parm.BusinessDriver.BusinessMetric {
			if idx == tempidx && i.BDId != temp.BDId {
				changes = temp.BDId
				if temp.BDId == "" {
					changes = "empty"
				}
				break
			}
		}
		if changes != "" {
			if changes == "empty" {
				i.BDId = ""
			} else {
				i.BDId = changes
			}
		}
		bm = append(bm, i)
	}
	d.BusinessMetric = bm
	d.UpdatedDate = time.Now()
	e = c.Ctx.Save(d)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	return c.ReMappingSummaryBD(k)
}

func (c *ScorecardController) SaveBusinessMetrics(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Region         string
		Country        string
		Id             int
		BusinessMetric []BusinessMetric
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	if parm.Region == "" || parm.Region == "Region" {
		parm.Region = "GLOBAL"
	}
	if parm.Country == "" || parm.Country == "Country" {
		parm.Country = parm.Region
	}

	d := new(BusinessDriverL1Model)
	csr, e := c.Ctx.Find(new(BusinessDriverL1Model), tk.M{}.Set("where", dbox.Eq("_id", parm.Id)))
	e = csr.Fetch(&d, 1, false)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	csr.Close()
	ExistingBusinessMetrics := map[string]BusinessMetric{}
	for _, x := range d.BusinessMetric {
		ExistingBusinessMetrics[x.Id] = x
	}
	d.BusinessMetric = []BusinessMetric{}
	// ReMap Period in [d.BusinessMetric]
	// query := []*dbox.Filter{}
	for _, i := range parm.BusinessMetric {
		if i.Id == "" {
			NewId, e := c.GetNextIdSeq("businessmetric")
			if e != nil {
				return c.ErrorResultInfo(e.Error(), nil)
			}
			i.Id = tk.Sprintf("BM%v", NewId)
			c.Action(k, "Scorecard", "Add Scorecard Key Metrics", "Metric Heading", "", i.DataPoint, "", "")
			newvalue := "Higher Is Better"
			if i.ValueType == 0 {
				newvalue = "Lower Is Better"
			}
			c.Action(k, "Scorecard", "Add Scorecard Key Metrics", "Metric Denomination", "", newvalue, "", "")
			c.Action(k, "Scorecard", "Add Scorecard Key Metrics", "Metric Type", "", i.Type, "", "")
			c.Action(k, "Scorecard", "Add Scorecard Key Metrics", "DecimalFormat", "", i.DecimalFormat, "", "")
			c.Action(k, "Scorecard", "Add Scorecard Key Metrics", "Metric Type", "", i.MetricType, "", "")
			c.Action(k, "Scorecard", "Add Scorecard Key Metrics", "Description", "", i.Description, "", "")
		} else {
			eData := ExistingBusinessMetrics[i.Id]
			if i.DataPoint != eData.DataPoint {
				c.Action(k, "Scorecard", "Update Scorecard Key Metrics", "Metric Heading", eData.DataPoint, i.DataPoint, "", "")
			}
			if i.ValueType != eData.ValueType {
				oldvalue := "Higher Is Better"
				if eData.ValueType == 0 {
					oldvalue = "Lower Is Better"
				}
				newvalue := "Lower Is Better"
				if i.ValueType == 0 {
					newvalue = "Lower Is Better"
				}
				c.Action(k, "Scorecard", "Update Scorecard Key Metrics", "Metric Denomination", oldvalue, newvalue, "", "")
			}
			if i.Type != eData.Type {
				c.Action(k, "Scorecard", "Update Scorecard Key Metrics", "Metric Type", eData.Type, i.Type, "", "")
			}
			if i.DecimalFormat != eData.DecimalFormat {
				c.Action(k, "Scorecard", "Update Scorecard Key Metrics", "DecimalFormat", eData.DecimalFormat, i.DecimalFormat, "", "")
			}
			if i.MetricType != eData.MetricType {
				c.Action(k, "Scorecard", "Update Scorecard Key Metrics", "Metric Type", eData.MetricType, i.MetricType, "", "")
			}
			if i.Description != eData.Description {
				c.Action(k, "Scorecard", "Update Scorecard Key Metrics", "Description", eData.Description, i.Description, "", "")
			}
			delete(ExistingBusinessMetrics, i.Id)
		}
		if i.CurrentPeriodStr != "" {
			i.CurrentPeriod, _ = time.Parse("20060102", i.CurrentPeriodStr)
		}
		if i.TargetPeriodStr != "" {
			i.TargetPeriod, _ = time.Parse("20060102", i.TargetPeriodStr)
		}
		ActualData := []ActualValue{}
		for _, actual := range i.ActualData {
			actual.Period, _ = time.Parse("20060102", actual.PeriodStr)
			ActualData = append(ActualData, actual)
		}
		i.ActualData = ActualData

		d.BusinessMetric = append(d.BusinessMetric, i)

		// Update Data
		// MetricData := BusinessMetricsDataModel{}
		// query = append(query[0:0], dbox.Eq("bmid", i.Id))
		// if parm.Region == parm.Country {
		// 	query = append(query, dbox.Eq("majorregion", parm.Region))
		// }
		// query = append(query, dbox.Eq("country", parm.Country))
		// csr, err := c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(dbox.And(query...)).Order("-period").Cursor(nil)
		// err = csr.Fetch(&MetricData, 1, false)
		// csr.Close()
		// if err != nil {
		// 	Now := time.Now()
		// 	// tk.Println(Now)
		// 	tempPeriod := time.Date(Now.Year(), Now.Month(), 1, 0, 0, 0, 0, time.UTC)
		// 	MData := NewBusinessMetricsDataModel()
		// 	MData.Period = tempPeriod
		// 	MData.Year = tempPeriod.Year()
		// 	MData.BusinessName = "COMMERCIAL BANKING"
		// 	MData.SCId = d.Id
		// 	MData.ScorecardCategory = d.Name
		// 	MData.BMId = i.Id
		// 	MData.BusinessMetric = i.DataPoint
		// 	MData.BusinessMetricDescription = i.Description
		// 	MData.MajorRegion = parm.Region
		// 	MData.Region = parm.Country
		// 	MData.CountryCode = parm.Country
		// 	MData.Country = parm.Country
		// 	if parm.Region != parm.Country {
		// 		RegionData := tk.M{}
		// 		csr, err := c.Ctx.Connection.NewQuery().From("Region").Where(dbox.Eq("Country", parm.Country)).Cursor(nil)
		// 		err = csr.Fetch(&RegionData, 1, false)
		// 		csr.Close()
		// 		if err == nil {
		// 			MData.MajorRegion = RegionData.GetString("Major_Region")
		// 			MData.CountryCode = RegionData.GetString("CountryCode")
		// 			MData.Region = RegionData.GetString("Region")
		// 		}
		// 	}

		// 	MData.Baseline = i.BaseLineValue
		// 	MData.ActualYTD = i.CurrentValue
		// 	MData.Target = i.TargetValue
		// 	MData.UpdatedDate = time.Now().Add(360)
		// 	MData.UpdatedBy = k.Session("username", "admin").(string)

		// 	MData.NABaseline = i.NABaseline
		// 	MData.NAActual = i.NAActual
		// 	MData.NATarget = i.NATarget
		// 	MData.RAG = i.Display
		// 	if i.NABaseline {
		// 		MData.Baseline = 130895111188
		// 	}
		// 	if i.NAActual {
		// 		MData.ActualYTD = 130895111188
		// 	}
		// 	if i.NATarget {
		// 		MData.Target = 130895111188
		// 	}
		// 	e = c.Ctx.Save(MData)
		// 	if e != nil {
		// 		return c.ErrorResultInfo(e.Error(), nil)
		// 	}

		// 	// Updating to Temp Data
		// 	nd := new(BusinessMetricsDataTempModel)
		// 	temp := tk.M{}
		// 	err = tk.StructToM(MData, &temp)
		// 	if err != nil {
		// 		return c.ErrorResultInfo(err.Error(), nil)
		// 	}
		// 	err = tk.MtoStruct(temp, nd)
		// 	if err != nil {
		// 		return c.ErrorResultInfo(err.Error(), nil)
		// 	}
		// 	err = c.Ctx.DeleteMany(nd, dbox.And(dbox.Eq("bmid", nd.BMId), dbox.Eq("period", nd.Period.UTC()), dbox.Eq("countrycode", nd.CountryCode)))
		// 	if err != nil {
		// 		return c.ErrorResultInfo(err.Error(), nil)
		// 	}
		// 	err = c.Ctx.Save(nd)
		// 	if err != nil {
		// 		return c.ErrorResultInfo(err.Error(), nil)
		// 	}
		// } else {

		// 	MetricData.Baseline = i.BaseLineValue
		// 	MetricData.ActualYTD = i.CurrentValue
		// 	MetricData.Target = i.TargetValue
		// 	MetricData.UpdatedDate = time.Now()
		// 	MetricData.UpdatedBy = k.Session("username", "admin").(string)
		// 	MetricData.NABaseline = i.NABaseline
		// 	MetricData.NAActual = i.NAActual
		// 	MetricData.NATarget = i.NATarget

		// 	if i.NABaseline {
		// 		MetricData.Baseline = 130895111188
		// 	}
		// 	if i.NAActual {
		// 		MetricData.ActualYTD = 130895111188
		// 	}
		// 	if i.NATarget {
		// 		MetricData.Target = 130895111188
		// 	}

		// 	MetricData.RAG = i.Display
		// 	e = c.Ctx.Save(&MetricData)
		// 	if e != nil {
		// 		return c.ErrorResultInfo(e.Error(), nil)
		// 	}

		// 	// Updating to Temp Data
		// 	nd := new(BusinessMetricsDataTempModel)
		// 	temp := tk.M{}
		// 	err = tk.StructToM(MetricData, &temp)
		// 	if err != nil {
		// 		return c.ErrorResultInfo(err.Error(), nil)
		// 	}
		// 	err = tk.MtoStruct(temp, nd)
		// 	if err != nil {
		// 		return c.ErrorResultInfo(err.Error(), nil)
		// 	}
		// 	err = c.Ctx.DeleteMany(nd, dbox.And(dbox.Eq("bmid", nd.BMId), dbox.Eq("period", nd.Period.UTC()), dbox.Eq("countrycode", nd.CountryCode)))
		// 	if err != nil {
		// 		return c.ErrorResultInfo(err.Error(), nil)
		// 	}
		// 	err = c.Ctx.Save(nd)
		// 	if err != nil {
		// 		return c.ErrorResultInfo(err.Error(), nil)
		// 	}

		// }
		// tk.Println(MetricData)
		// tk.Println("##")
	}
	for _, eData := range ExistingBusinessMetrics {
		c.Action(k, "Scorecard", "Remove Scorecard Key Metrics", "Metric Heading", eData.DataPoint, "", "", "")
		oldvalue := "Higher Is Better"
		if eData.ValueType == 0 {
			oldvalue = "Lower Is Better"
		}
		c.Action(k, "Scorecard", "Remove Scorecard Key Metrics", "Metric Denomination", oldvalue, "", "", "")
		c.Action(k, "Scorecard", "Remove Scorecard Key Metrics", "Metric Type", eData.Type, "", "", "")
		c.Action(k, "Scorecard", "Remove Scorecard Key Metrics", "DecimalFormat", eData.DecimalFormat, "", "", "")
		c.Action(k, "Scorecard", "Remove Scorecard Key Metrics", "Metric Type", eData.MetricType, "", "", "")
		c.Action(k, "Scorecard", "Remove Scorecard Key Metrics", "Description", eData.Description, "", "", "")

	}
	d.UpdatedDate = time.Now()
	e = c.Ctx.Save(d)
	c.Action(k, "Scorecard", "Save Changes for Scorecard Key Metrics", "", "", "", "", "")
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	return c.ReMappingSummaryBD(k)
}

func (c *ScorecardController) ReMappingSummaryBD(k *knot.WebContext) interface{} {
	temp := make([]SummaryBusinessDriverModel, 0)
	csr, e := c.Ctx.Find(new(SummaryBusinessDriverModel), nil)
	e = csr.Fetch(&temp, 0, false)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	for _, x := range temp {
		x.BusinessMetrics = []BusinessMetric{}
		e = c.Ctx.Save(&x)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
	}
	csr.Close()

	data := make([]BusinessDriverL1Model, 0)
	csr, e = c.Ctx.Find(new(BusinessDriverL1Model), nil)
	e = csr.Fetch(&data, 0, false)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	csr.Close()

	for _, d := range data {
		for _, i := range d.BusinessMetric {
			bd := new(SummaryBusinessDriverModel)
			if i.BDId != "" {
				csr, e = c.Ctx.Find(new(SummaryBusinessDriverModel), tk.M{}.Set("where", dbox.And(dbox.Eq("Id", i.BDId))))
				e = csr.Fetch(&bd, 1, false)
				if e != nil {
					return c.ErrorResultInfo(e.Error(), nil)
				}
				bd.BusinessMetrics = []BusinessMetric{}
				// bd.BusinessMetrics = append(bd.BusinessMetrics, i)
				csr.Close()
				e = c.Ctx.Save(bd)
				if e != nil {
					return c.ErrorResultInfo(e.Error(), nil)
				}
			}
		}
	}

	for _, d := range data {
		for _, i := range d.BusinessMetric {
			bd := new(SummaryBusinessDriverModel)
			if i.BDId != "" {
				csr, e = c.Ctx.Find(new(SummaryBusinessDriverModel), tk.M{}.Set("where", dbox.And(dbox.Eq("Id", i.BDId))))
				e = csr.Fetch(&bd, 1, false)
				if e != nil {
					return c.ErrorResultInfo(e.Error(), nil)
				}
				// bd.BusinessMetrics = []BusinessMetric{}
				bd.BusinessMetrics = append(bd.BusinessMetrics, i)
				csr.Close()
				e = c.Ctx.Save(bd)
				if e != nil {
					return c.ErrorResultInfo(e.Error(), nil)
				}
			}
		}
	}
	return c.SetResultInfo(false, "", nil)

}

// ReMappingSummaryBD
func (c *ScorecardController) SyncBMToMasterData(k *knot.WebContext, BD *BusinessDriverL1Model, MetricId string, data *BusinessMetricsDataModel) error {
	selectedMetrics := BusinessMetric{}
	for _, i := range BD.BusinessMetric {
		if i.Id == MetricId {
			isAvailable := false
			for _, a := range i.ActualData {
				if a.Period == data.Period {
					isAvailable = true
					a.Value = data.Actual
					// Set Current Value after set actual value
					if i.CurrentPeriodStr == a.PeriodStr {
						i.CurrentValue = data.Actual
					}
				}
			}
			if !isAvailable {
				adata := ActualValue{}
				adata.Period = data.Period
				adata.PeriodStr = data.Period.Format("20060102")
				adata.Value = data.Actual
				i.ActualData = append(i.ActualData, adata)
				// type ActualValue struct {
				// 	Period    time.Time
				// 	PeriodStr string
				// 	Value     float64
				// 	Flag      string
				// }
			}
			i.TargetValue = data.Target
			i.TargetPeriod = data.Period
			i.TargetPeriodStr = data.Period.Format("20060102")
			i.BaseLineValue = data.Baseline
			i.BaseLinePeriod = data.Period
			i.BaseLinePeriodStr = data.Period.Format("20060102")
			selectedMetrics = i
		}
	}

	for counter, i := range BD.BusinessMetric {
		if i.Id == MetricId {
			i = selectedMetrics
			BD.BusinessMetric[counter].UpdatedDate = time.Now()
		}
	}

	e := c.Ctx.Save(BD)
	if e != nil {
		tk.Println("errorrr >>> ", e)
		return e
	}
	return nil
}

func (c *ScorecardController) MetricValue(val float64) string {
	if val == 130895111188 {
		return "N/A"
	} else {
		return strconv.FormatFloat(val, 'f', -1, 64)
	}
}

func (c *ScorecardController) ProcessBMData(k *knot.WebContext) interface{} {
	// c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Data          []tk.M
		UploadOptions string
	}{}

	err := k.GetPayload(&parm)
	c.Action(k, "Scorecard", "Process Metrics Data", "( "+strconv.Itoa(len(parm.Data))+" ) Metrics", "", "", "", "")

	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	result := "678deae8142518f6910077508ac4901c"
	for _, i := range parm.Data {
		DataPoint := i.GetString("DataPoint")
		// FileName := i.GetString("OriginalFileName")
		TrueFileName := i.GetString("FileName")
		c.Action(k, "Scorecard", "Process Metrics Data", DataPoint, "", "", "file", TrueFileName)
		datalist := i.Get("DataList").([]interface{})
		metricId := i.GetString("MetricsID")
		result += metricId + "s"
		/*masterBD := new(BusinessDriverL1Model)
		csr, err := c.Ctx.Find(new(BusinessDriverL1Model), tk.M{}.Set("where", dbox.Eq("businessmetric.id", metricId))) //
		err = csr.Fetch(&masterBD, 1, false)
		if err != nil {
			return c.SetResultInfo(true, err.Error(), nil)
		}
		csr.Close()*/
		bmdata := []BusinessMetricsDataTempModel{}
		m := NewBusinessMetricsNotificationModel()

		existngBM := make([]BusinessMetricsDataTempModel, 0)
		csr, err := c.Ctx.Find(new(BusinessMetricsDataTempModel), tk.M{}.Set("where", dbox.And(dbox.Eq("bmid", metricId)))) //
		err = csr.Fetch(&existngBM, 0, false)
		csr.Close()
		if err != nil {
			return c.SetResultInfo(true, err.Error(), nil)
		}
		// tk.Println(tk.JsonString(existngBM))
		for _, d := range datalist {
			tmp, err := tk.ToM(d)
			if err != nil {
				return c.SetResultInfo(true, err.Error(), nil)
			}

			if parm.UploadOptions == "owner" {
				data := NewBusinessMetricsDataTempModel()
				data.Source = TrueFileName

				err = tk.MtoStruct(tmp, data)
				if err != nil {
					return c.SetResultInfo(true, err.Error(), nil)
				}

				if data.ActualYTD == 130895111188 && tk.IsNilOrEmpty(data.RAG) {
					continue
				}

				//delete before save. if there are more than 1 BM with same year, country and region.
				isAny := false
				for _, i := range existngBM {
					if data.CountryCode == i.CountryCode && data.Period.UTC() == i.Period.UTC() {
						isAny = true
						c.Ctx.Delete(&i)

						// if data.Baseline == 130895111188 {
						// data.Baseline = i.Baseline
						// }
						// if data.Actual == 130895111188 {
						// data.Actual = i.Actual
						// }

						if data.ActualYTD == 130895111188 {
							data.ActualYTD = i.ActualYTD
						} else {
							if data.ActualYTD != i.ActualYTD {
								c.Action(k, "Scorecard", "Update Key Metrics Data ", "YTD Actual", c.MetricValue(i.ActualYTD), c.MetricValue(data.ActualYTD), "file", TrueFileName)
							}
						}
						if data.RAG != i.RAG {
							c.Action(k, "Scorecard", "Update Key Metrics Data ", "YTD Actual", strings.ToUpper(i.RAG), strings.ToUpper(data.RAG), "file", TrueFileName)
						}

						// if data.FullYearForecast == 130895111188 {
						// data.FullYearForecast = i.FullYearForecast
						// }
						// if data.Target == 130895111188 {
						// data.Target = i.Target
						// }
						data.Baseline = i.Baseline
						data.Actual = i.Actual
						data.FullYearForecast = i.FullYearForecast
						data.Target = i.Target
						data.NABudget = i.NABudget
						data.Budget = i.Budget
					}

				}

				if isAny == false {
					c.Action(k, "Scorecard", "Add Key Metrics Data ", "YTD Actual", "", c.MetricValue(data.ActualYTD), "file", TrueFileName)
					if data.RAG != "" {
						c.Action(k, "Scorecard", "Add Key Metrics Data ", "RAG", "", strings.ToUpper(data.RAG), "file", TrueFileName)
					}
					data.NABudget = true
					data.Budget = 130895111188
					data.Baseline = 130895111188
					data.FullYearForecast = 130895111188
					data.Target = 130895111188
					data.Actual = 130895111188
				}

				if data.ActualYTD != 130895111188 {
					bmdata = append(bmdata, *data)
				}

				err = c.Ctx.Save(data)
				if err != nil {
					return c.SetResultInfo(true, err.Error(), nil)
				}

				// save to notification
				m.BMId = data.BMId
				m.BMName = data.BusinessMetric
				m.Period = data.Period
				m.Type = "owner"
				m.Updated_Date = time.Now()
				m.Updated_By = data.UpdatedBy
				m.Source = data.Source
			} else {
				// if irow == 0 {
				financeColumn := i.Get("finCol").([]interface{})
				for _, fCol := range financeColumn {
					data := NewBusinessMetricsDataTempModel()
					data.Source = TrueFileName

					err = tk.MtoStruct(tmp, data)
					if err != nil {
						return c.SetResultInfo(true, err.Error(), nil)
					}

					period := tk.String2Date(fCol.(string), "MMMyyyy").UTC()
					countrycode := tmp.GetString("CountryCode")
					budgetPeriod := tmp.GetFloat64(fCol.(string))

					if budgetPeriod == 130895111188 && tmp.GetFloat64("Baseline") == 130895111188 && tmp.GetFloat64("Target") == 130895111188 {
						continue
					}

					isAny := false

					data.Period = period
					data.Budget = budgetPeriod
					data.Baseline = tmp.GetFloat64("Baseline")
					data.Target = tmp.GetFloat64("Target")
					for _, i := range existngBM {
						if countrycode == i.CountryCode && period == i.Period.UTC() {
							isAny = true
							c.Ctx.Delete(&i)
							data.Id = i.Id
							if data.Baseline == 130895111188 {
								data.Baseline = i.Baseline
							} else {
								if data.Baseline != i.Baseline {
									c.Action(k, "Scorecard", "Update Key Metrics Data ", "Baseline", c.MetricValue(i.Baseline), c.MetricValue(data.Baseline), "file", TrueFileName)
								}
							}
							if data.Target == 130895111188 {
								data.Target = i.Target
							} else {
								if data.Budget != i.Budget {
									c.Action(k, "Scorecard", "Update Key Metrics Data ", "Target", c.MetricValue(i.Target), c.MetricValue(data.Target), "file", TrueFileName)
								}
							}
							// data.Budget = budgetPeriod
							if data.Budget == 130895111188 {
								data.Budget = i.Budget
							} else {
								if data.Budget != i.Budget {
									c.Action(k, "Scorecard", "Update Key Metrics Data ", "Budget", c.MetricValue(i.Budget), c.MetricValue(data.Budget), "file", TrueFileName)
								}
							}
							data.Actual = i.Actual
							data.ActualYTD = i.ActualYTD
							data.NAActual = i.NAActual
							data.RAG = i.RAG
						}

						// data.Period = period
					}
					if isAny == false {
						data.Actual = 0
						data.ActualYTD = 130895111188
						data.NAActual = true
						data.RAG = ""
						c.Action(k, "Scorecard", "Add Key Metrics Data ", "Baseline", "", c.MetricValue(data.Baseline), "file", TrueFileName)
						c.Action(k, "Scorecard", "Add Key Metrics Data ", "Target", "", c.MetricValue(data.Target), "file", TrueFileName)
						c.Action(k, "Scorecard", "Add Key Metrics Data ", "Budget", "", c.MetricValue(data.Budget), "file", TrueFileName)

					}

					if data.Budget != 130895111188 {
						bmdata = append(bmdata, *data)
					}

					err = c.Ctx.Save(data)
					if err != nil {
						return c.SetResultInfo(true, err.Error(), nil)
					}

					// save to notification
					m.BMId = data.BMId
					m.BMName = data.BusinessMetric
					m.Period = time.Time{}
					m.Type = "finance"
					m.Updated_Date = time.Now()
					m.Updated_By = data.UpdatedBy
					m.Source = data.Source
				}
			}
		}

		// m.BMData = bmdata
		// err = c.Ctx.Save(m)
		// if err != nil {
		// 	return c.SetResultInfo(true, err.Error(), nil)
		// }
		if !tk.IsNilOrEmpty(m.BMId) {
			m.BMData = bmdata
			err = c.Ctx.Save(m)
			if err != nil {
				return c.SetResultInfo(true, err.Error(), nil)
			}
		}

	}
	return c.SetResultInfo(false, "", strings.ToLower(result))
}

func SpaceMap(str string) string {
	nospace := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
	// reg, _ := regexp.Compile("[:/<>*?|]+")
	// result := reg.ReplaceAllString(nospace,"")
	// return result
	r := strings.NewReplacer("/", "-",
		"\\'", "-")
	result := r.Replace(nospace)
	return result
}

func (c *ScorecardController) GetBMFileTemplate(k *knot.WebContext) interface{} {
	type regionModel struct {
		MajorRegion string
		CountryName string
		CountryCode string
	}
	c.LoadBase(k)
	c.Action(k, "Scorecard", "Download Template Metric Owner", "", "", "", "", "")

	k.Config.OutputType = knot.OutputJson
	frm := struct {
		BusinessName []string
		BmId         []string
		BmName       []string
		Date         string
	}{}
	err := k.GetPayload(&frm)
	if err != nil {
		c.ErrorResultInfo(err.Error(), nil)
	}
	// listCountry := c.loadMasterCountryList()
	ret := ResultInfo{}
	dataRegion := make([]tk.M, 0)
	listDataRegion := make([]regionModel, 0)
	crs, err := c.Ctx.Connection.NewQuery().From("Region").Order("Major_Region", "Country").Cursor(nil)
	defer crs.Close()
	if err != nil {
		c.ErrorResultInfo(err.Error(), nil)
	}
	err = crs.Fetch(&dataRegion, 0, false)
	if err != nil {
		c.ErrorResultInfo(err.Error(), nil)
	}
	for _, dt := range dataRegion {
		// cntry := listCountry[dt.GetString("Country")]
		tmp := regionModel{
			MajorRegion: dt.GetString("Major_Region"),
			CountryName: dt.GetString("Country"),
			// CountryCode: cntry,
			CountryCode: dt.GetString("CountryCode"),
		}
		listDataRegion = append(listDataRegion, tmp)
	}
	fileName := SpaceMap(frm.BusinessName[0]) + "_" + SpaceMap(frm.BmName[0]) + "_" + frm.Date + ".xlsx"
	xlsFile := xl.NewFile()
	sheet, err := xlsFile.AddSheet(frm.BusinessName[0])
	if err != nil {
		c.ErrorResultInfo(err.Error(), nil)
	}
	headers := []string{
		"BUSINESS NAME",
		"METRIC 1",
		"REGION",
		"COUNTRY",
		"COUNTRY CODE",
		"BASELINE",
		"MONTHLY ACTUAL",
		"YTD ACTUAL",
		"FULL ACTUAL FORECAST",
		"FULL YEAR TARGET",
	}

	font := xl.NewFont(11, "Calibri")
	style := xl.NewStyle()
	style.Font = *font

	fontHdr := xl.NewFont(11, "Calibri")
	fontHdr.Bold = true
	styleHdr := xl.NewStyle()
	styleHdr.Font = *fontHdr
	row := sheet.AddRow()
	listMetricType := c.loadMasterBusinessL1Type()
	listDecimalType := c.loadMasterBusinessL1DecimalType()
	// fmt.Println(listDecimalType)

	//	tk.Printfn("MetricList %v", listMetricType)
	for _, txt := range headers {
		cell := row.AddCell()
		cell.SetStyle(styleHdr)
		cell.SetValue(txt)
	}

	for ctx0, metric := range frm.BmName {
		// fmt.Println(listMetricType, listMetricType[frm.BmId[ctx0]], frm.BmId[ctx0], ctx0)
		row = sheet.AddRow()
		cell := row.AddCell()
		cell.SetStyle(style)
		cell.SetValue(frm.BusinessName[0])

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.SetValue(strings.TrimSpace(metric))

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.SetValue("GLOBAL")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.SetValue("GLOBAL")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.SetValue("GLOBAL")

		bmType := listMetricType[frm.BmId[ctx0]]

		decimalType := "0"

		if listDecimalType[frm.BmId[ctx0]] != "" {
			decimalType = listDecimalType[frm.BmId[ctx0]]
		}

		cell = row.AddCell()
		cell.SetStyle(style)
		c.setCellFormat(cell, bmType, decimalType)
		cell.SetValue(nil)

		cell = row.AddCell()
		cell.SetStyle(style)
		c.setCellFormat(cell, bmType, decimalType)
		cell.SetValue(nil)

		cell = row.AddCell()
		cell.SetStyle(style)
		c.setCellFormat(cell, bmType, decimalType)
		cell.SetValue(nil)

		cell = row.AddCell()
		cell.SetStyle(style)
		c.setCellFormat(cell, bmType, decimalType)
		cell.SetValue(nil)

		cell = row.AddCell()
		cell.SetStyle(style)
		c.setCellFormat(cell, bmType, decimalType)
		cell.SetValue(nil)

		cell = row.AddCell()
		cell.SetStyle(style)
		c.setCellFormat(cell, bmType, decimalType)
		cell.SetValue(nil)

		for ctx, data := range listDataRegion {
			row = sheet.AddRow()
			cell := row.AddCell()
			cell.SetStyle(style)
			cell.SetValue(frm.BusinessName[0])

			cell = row.AddCell()
			cell.SetStyle(style)
			cell.SetValue(strings.TrimSpace(metric))

			cell = row.AddCell()
			cell.SetStyle(style)
			cell.SetValue(data.MajorRegion)

			cell = row.AddCell()
			cell.SetStyle(style)
			cell.SetValue(data.CountryName)

			cell = row.AddCell()
			cell.SetStyle(style)
			cell.SetValue(data.CountryCode)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			if ctx == len(listDataRegion)-1 || (ctx > 0 && data.MajorRegion != listDataRegion[ctx+1].MajorRegion) {
				row = sheet.AddRow()
				cell := row.AddCell()
				cell.SetStyle(style)
				cell.SetValue(frm.BusinessName[0])

				cell = row.AddCell()
				cell.SetStyle(style)
				cell.SetValue(strings.TrimSpace(metric))

				cell = row.AddCell()
				cell.SetStyle(style)
				cell.SetValue(data.MajorRegion)

				cell = row.AddCell()
				cell.SetStyle(style)
				cell.SetValue(data.MajorRegion + " TOTAL")

				cell = row.AddCell()
				cell.SetStyle(style)
				cell.SetValue(data.MajorRegion)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

			}
		}
	}
	err = xlsFile.Save(c.TemplatePath + "/" + fileName)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	ret.Data = fileName
	return ret
}

func (c *ScorecardController) GetBMFileTemplateMultiple(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	c.Action(k, "Scorecard", "Download Template Metric Owner", "", "", "", "", "")
	type regionModel struct {
		MajorRegion string
		CountryName string
		CountryCode string
	}

	k.Config.OutputType = knot.OutputJson
	frm := struct {
		BusinessName   []string
		DownloadOption string
		BmName         []tk.M
		Date           string
	}{}
	err := k.GetPayload(&frm)
	if err != nil {
		c.ErrorResultInfo(err.Error(), nil)
	}
	ret := ResultInfo{}

	dataRegion := make([]tk.M, 0)
	crs, err := c.Ctx.Connection.NewQuery().From("Region").Order("Major_Region", "Country").Cursor(nil)
	defer crs.Close()
	if err != nil {
		c.ErrorResultInfo(err.Error(), nil)
	}
	err = crs.Fetch(&dataRegion, 0, false)
	if err != nil {
		c.ErrorResultInfo(err.Error(), nil)
	}

	listDataRegion := make([]regionModel, 0)
	for _, dt := range dataRegion {
		// cntry := listCountry[dt.GetString("Country")]
		tmp := regionModel{
			MajorRegion: dt.GetString("Major_Region"),
			CountryName: dt.GetString("Country"),
			CountryCode: dt.GetString("CountryCode"),
		}
		listDataRegion = append(listDataRegion, tmp)
	}

	fileName := frm.BusinessName[0] + "_Finance Team Template.xlsx" //SpaceMap(frm.BusinessName[0]) + ".xlsx"
	if frm.DownloadOption != "finance" {
		fileName = frm.BusinessName[0] + "_Metric Owner Template_" + frm.Date + ".xlsx"
	}
	/// if exist remove file
	if _, err := os.Stat(filepath.Join(c.TemplatePath, fileName)); tk.IsNilOrEmpty(err) {
		err = os.Remove(filepath.Join(c.TemplatePath, fileName))
		if !tk.IsNilOrEmpty(err) {
			return c.ErrorResultInfo(err.Error(), nil)
		}
	}

	//// process creating new xlsx
	xlsFile := xl.NewFile()
	font := xl.NewFont(11, "Calibri")
	style := xl.NewStyle()
	border := xl.NewBorder("thin", "thin", "thin", "thin")
	style.Font = *font
	style.Border = *border

	fontHdr := xl.NewFont(12, "Calibri")
	fontHdr.Bold = true
	styleHdr := xl.NewStyle()
	styleHdr.Font = *fontHdr
	styleHdr.Font.Color = "000000"
	styleHdr.Alignment.WrapText = true
	styleHdr.Fill = *xl.NewFill("solid", "a5c5f7", "a5c5f7")
	styleHdr.Border = *border

	styleFillGrey := xl.NewStyle()
	styleFillGrey.Fill = *xl.NewFill("solid", "d1cfcf", "d1cfcf")
	styleFillGrey.Border = *border

	sheet, err := xlsFile.AddSheet(frm.BusinessName[0])
	if err != nil {
		c.ErrorResultInfo(err.Error(), nil)
	}

	headers := []string{}
	if frm.DownloadOption == "owner" {
		headers = []string{
			"BUSINESS NAME",
			"DATA MONTH",
			"METRIC",
			"REGION",
			"COUNTRY",
			"COUNTRY CODE",
			"METRIC DENOMINATION",
			"DECIMAL POINT",
			"YTD ACTUAL",
			"RAG",
		}
	} else {
		currentyear, _, _ := time.Now().Date()
		headers = []string{
			"BUSINESS NAME",
			"METRIC",
			"REGION",
			"COUNTRY",
			"COUNTRY CODE",
			"BASELINE",
			"Jan " + tk.ToString(currentyear),
			"Feb " + tk.ToString(currentyear),
			"Mar " + tk.ToString(currentyear),
			"Apr " + tk.ToString(currentyear),
			"May " + tk.ToString(currentyear),
			"Jun " + tk.ToString(currentyear),
			"Jul " + tk.ToString(currentyear),
			"Aug " + tk.ToString(currentyear),
			"Sep " + tk.ToString(currentyear),
			"Oct " + tk.ToString(currentyear),
			"Nov " + tk.ToString(currentyear),
			"Dec " + tk.ToString(currentyear),
			"FULL YEAR TARGET",
		}
	}
	row := sheet.AddRow()
	for _, txt := range headers {
		cell := row.AddCell()
		cell.SetStyle(styleHdr)
		cell.SetValue(txt)
	}

	// tk.Println("BusinessName", frm.BusinessName)

	for _, metric := range frm.BmName {
		// datas := make([]tk.M, 0)
		// csr, err := c.Ctx.Connection.NewQuery().From("BusinessDriverL1").Where(dbox.Eq(field, frm.BusinessName))Cursor(nil)
		// if err != nil {
		// 	return nil
		// }
		// csr.Fetch(&datas, 1, false)
		// csr.Close()

		listMetricType := c.loadMasterBusinessL1Type()
		listDecimalType := c.loadMasterBusinessL1DecimalType()

		// tk.Println("listMetricType", listMetricType)
		// tk.Println("listDecimalType", listDecimalType)

		// fmt.Println(listMetricType, listMetricType[frm.BmId[ctx0]], frm.BmId[ctx0], ctx0)
		row = sheet.AddRow()
		cell := row.AddCell()
		cell.SetStyle(style)
		cell.SetValue(frm.BusinessName[0])

		if frm.DownloadOption == "owner" {
			cell = row.AddCell()
			cell.SetStyle(style)
			cell.SetValue(frm.Date)
		}

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.SetValue(strings.TrimSpace(metric.GetString("text")))

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.SetValue("GLOBAL")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.SetValue("GLOBAL")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.SetValue("GLOBAL")

		bmType := listMetricType[metric.GetString("metricId")]
		decimalType := "0"
		if listDecimalType[metric.GetString("metricId")] != "" {
			decimalType = listDecimalType[metric.GetString("metricId")]
		}
		// tk.Println("decimalType", decimalType)

		if frm.DownloadOption == "owner" {
			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType) //metric denom
			if bmType == "DOLLAR" {
				cell.SetValue("Dollar Value ($)")
			} else if bmType == "PERCENTAGE" {
				cell.SetValue("Percentage (%)")
			} else if bmType == "NUMERIC" {
				cell.SetValue("Numeric Value")
			} else {
				cell.SetValue("Numeric Value")
			}

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType) //decpoint
			if decimalType == "0" || decimalType == "" {
				cell.SetValue("No Decimal")
			} else if decimalType == "1" {
				cell.SetValue("0.0")
			} else if decimalType == "2" {
				cell.SetValue("0.00")
			}

			cell = row.AddCell()
			cell.SetStyle(styleFillGrey)
			c.setCellFormat(cell, bmType, decimalType) //ytd actual
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(styleFillGrey)
			// c.setCellFormat(cell, bmType, decimalType) //rag
			cell.SetValue(nil)
		} else {
			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType) //baseline
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType)
			cell.SetValue(nil)

			cell = row.AddCell()
			cell.SetStyle(style)
			c.setCellFormat(cell, bmType, decimalType) //target
			cell.SetValue(nil)
		}

		for ctx, data := range listDataRegion {
			row = sheet.AddRow()
			cell := row.AddCell()
			cell.SetStyle(style)
			cell.SetValue(frm.BusinessName[0])

			if frm.DownloadOption == "owner" {
				cell = row.AddCell()
				cell.SetStyle(style)
				cell.SetValue(frm.Date)
			}

			cell = row.AddCell()
			cell.SetStyle(style)
			cell.SetValue(strings.TrimSpace(metric.GetString("text")))

			cell = row.AddCell()
			cell.SetStyle(style)
			cell.SetValue(data.MajorRegion)

			cell = row.AddCell()
			cell.SetStyle(style)
			cell.SetValue(data.CountryName)

			cell = row.AddCell()
			cell.SetStyle(style)
			cell.SetValue(data.CountryCode)

			if frm.DownloadOption == "owner" {
				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType) //metric denom
				if bmType == "DOLLAR" {
					cell.SetValue("Dollar Value ($)")
				} else if bmType == "PERCENTAGE" {
					cell.SetValue("Percentage (%)")
				} else if bmType == "NUMERIC" {
					cell.SetValue("Numeric Value")
				} else {
					cell.SetValue("Numeric Value")
				}

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType) //decpoint
				if decimalType == "0" || decimalType == "" {
					cell.SetValue("No Decimal")
				} else if decimalType == "1" {
					cell.SetValue("0.0")
				} else if decimalType == "2" {
					cell.SetValue("0.00")
				}

				cell = row.AddCell()
				cell.SetStyle(styleFillGrey)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(styleFillGrey)
				// c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)
			} else {
				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType) //baseline
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType)
				cell.SetValue(nil)

				cell = row.AddCell()
				cell.SetStyle(style)
				c.setCellFormat(cell, bmType, decimalType) //target
				cell.SetValue(nil)
			}

			if ctx == len(listDataRegion)-1 || (ctx > 0 && data.MajorRegion != listDataRegion[ctx+1].MajorRegion) {
				row = sheet.AddRow()
				cell := row.AddCell()
				cell.SetStyle(style)
				cell.SetValue(frm.BusinessName[0])

				if frm.DownloadOption == "owner" {
					cell = row.AddCell()
					cell.SetStyle(style)
					cell.SetValue(frm.Date)
				}

				cell = row.AddCell()
				cell.SetStyle(style)
				cell.SetValue(strings.TrimSpace(metric.GetString("text")))

				cell = row.AddCell()
				cell.SetStyle(style)
				cell.SetValue(data.MajorRegion)

				cell = row.AddCell()
				cell.SetStyle(style)
				cell.SetValue(data.MajorRegion + " TOTAL")

				cell = row.AddCell()
				cell.SetStyle(style)
				cell.SetValue(data.MajorRegion)

				if frm.DownloadOption == "owner" {
					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType) //metric denom
					if bmType == "DOLLAR" {
						cell.SetValue("Dollar Value ($)")
					} else if bmType == "PERCENTAGE" {
						cell.SetValue("Percentage (%)")
					} else if bmType == "NUMERIC" {
						cell.SetValue("Numeric Value")
					} else {
						cell.SetValue("Numeric Value")
					}

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType) //decpoint
					if decimalType == "0" || decimalType == "" {
						cell.SetValue("No Decimal")
					} else if decimalType == "1" {
						cell.SetValue("0.0")
					} else if decimalType == "2" {
						cell.SetValue("0.00")
					}

					cell = row.AddCell()
					cell.SetStyle(styleFillGrey)
					c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(styleFillGrey)
					// c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)
				} else {
					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType) //baseline
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType)
					cell.SetValue(nil)

					cell = row.AddCell()
					cell.SetStyle(style)
					c.setCellFormat(cell, bmType, decimalType) //target
					cell.SetValue(nil)
				}
			}
		}
	}

	err = xlsFile.Save(c.TemplatePath + "/" + fileName)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	//// process creating new xlsx
	ret.Data = fileName
	return ret
}

func (c *ScorecardController) setCellFormat(cell *xl.Cell, btype string, btype2 string) {
	// fmt.Println(btype2)
	switch btype {
	case "PERCENTAGE":
		if btype2 == "0" {
			cell.SetFloatWithFormat(0, "#0%")
		} else if btype2 == "1" {
			cell.SetFloatWithFormat(0, "#0.0%")
		} else {
			cell.SetFloatWithFormat(0, "#0.00%")
		}

	case "NUMERIC":
		if btype2 == "0" {
			cell.SetFloatWithFormat(0, "#,##0")
		} else if btype2 == "1" {
			cell.SetFloatWithFormat(0, "#,##0.0")
		} else {
			cell.SetFloatWithFormat(0, "#,##0.00")
		}
	default:
		if btype2 == "0" {
			cell.SetFloatWithFormat(0, "#,##0")
		} else if btype2 == "1" {
			cell.SetFloatWithFormat(0, "#,##0.0")
		} else {
			cell.SetFloatWithFormat(0, "#,##0.00")
		}
		break
	}
}

func (c *ScorecardController) loadMasterBusinessL1Type() map[string]string {
	ret := map[string]string{}
	pipe := []tk.M{
		tk.M{}.Set("$unwind", "$businessmetric"),
		tk.M{}.Set("$group", tk.M{}.Set("_id", tk.M{}.Set("bmid", "$businessmetric.id").Set("bmtype", "$businessmetric.MetricType"))),
	}
	crs, err := c.Ctx.Connection.NewQuery().Command("pipe", pipe).From("BusinessDriverL1").Cursor(nil)
	if err != nil {
		return ret
	}
	ds := []tk.M{}
	err = crs.Fetch(&ds, 0, false)
	defer crs.Close()
	if err != nil {
		return ret
	}
	for _, dt := range ds {
		tmp := dt["_id"].(tk.M)
		ret[tmp.GetString("bmid")] = tmp.GetString("bmtype")
	}
	return ret
}

func (c *ScorecardController) loadMasterBusinessL1DecimalType() map[string]string {
	ret := map[string]string{}
	pipe := []tk.M{
		tk.M{}.Set("$unwind", "$businessmetric"),
		tk.M{}.Set("$group", tk.M{}.Set("_id", tk.M{}.Set("bmid", "$businessmetric.id").Set("decimalformat", "$businessmetric.DecimalFormat"))),
	}
	crs, err := c.Ctx.Connection.NewQuery().Command("pipe", pipe).From("BusinessDriverL1").Cursor(nil)
	if err != nil {
		return ret
	}
	ds := []tk.M{}
	err = crs.Fetch(&ds, 0, false)
	defer crs.Close()
	if err != nil {
		return ret
	}
	for _, dt := range ds {
		tmp := dt["_id"].(tk.M)
		ret[tmp.GetString("bmid")] = tmp.GetString("decimalformat")
	}
	return ret
}

func (c *ScorecardController) loadMasterCountryList() map[string]string {
	ret := map[string]string{}
	data := make([]tk.M, 0)
	crs, err := c.Ctx.Connection.NewQuery().From("CountryList").Cursor(nil)
	defer crs.Close()
	if err != nil {
		return ret
	}
	err = crs.Fetch(&data, 0, false)
	if err != nil {
		return ret
	}
	for _, dt := range data {
		ret[dt.GetString("BEFCountryName")] = dt.GetString("Alpha2Code")
	}
	return ret
}

// checking for owner template
func (c *ScorecardController) CheckUploadMetrics(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	err, payload := c.uploadTemplate(k)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	exelFile, e := xl.OpenFile(filepath.Join(c.UploadMetricPath, payload.GetString("file")))
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	exist := tk.M{}
	isExist := false
	bmNameExist := ""
	for idx, row := range exelFile.Sheets[0].Rows {
		if idx == 0 || len(row.Cells) < 3 {
			continue
		}
		/*dateUpload, _ := row.Cells[1].String()
		periodUpload, _ := time.Parse("200601", payload.GetString("date"))
		tk.Println(dateUpload, periodUpload.Format("Jan 2006"))
		if dateUpload != periodUpload.Format("Jan 2006") {
			return c.ErrorResultInfo("Date upload not identical with data month file", nil)
		}*/

		kolom2, _ := row.Cells[2].String()
		if bmNameExist != kolom2 {
			bmNameExist = kolom2
			period, _ := time.Parse("200601", payload.GetString("date"))

			d := make([]BusinessMetricsDataModel, 0)
			csr, e := c.Ctx.Find(new(BusinessMetricsDataModel), tk.M{}.Set("where", dbox.And(dbox.Eq("businessmetric", kolom2), dbox.Eq("period", period))))
			e = csr.Fetch(&d, 0, false)
			if e != nil {
				return c.ErrorResultInfo(e.Error(), nil)
			}
			csr.Close()
			// tk.Println("period >>>>>", tk.JsonString(d))
			if len(d) > 0 {
				exist.Set(d[0].BusinessMetric, tk.Date2String(d[0].Period, "MMM yyyy"))
				isExist = true
			}
		}
	}

	//remove file
	err = os.Remove(filepath.Join(c.UploadMetricPath, payload.GetString("file")))
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	if isExist {
		bmane := ""
		for key, v := range exist {
			bmane += key + " - " + v.(string) + "\n"
		}
		return c.SetResultInfo(false, bmane+" have already been uploaded, are you sure you want to upload again?", tk.M{"exist": exist, "payload": payload})
	}
	return c.SetResultInfo(false, "", tk.M{"payload": payload}) //c.SetResultInfo(false, "Data for this metric and this month have already been uploaded, are you sure you want to upload again?", isExist)
}

// checking for finance template
func (c *ScorecardController) CheckUploadFinanceMetrics(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	err, payload := c.uploadTemplate(k)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	exelFile, e := xl.OpenFile(filepath.Join(c.UploadMetricPath, payload.GetString("file")))
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	exist := tk.M{}
	isExist := false
	bmNameExist := ""
	// uploadperiod := []time.Time{}
	for idx, row := range exelFile.Sheets[0].Rows {
		if idx == 0 {
			continue
		}

		kolom2, _ := row.Cells[1].String()
		if bmNameExist != kolom2 {
			bmNameExist = kolom2

			d := tk.Ms{}
			y, _, _ := time.Now().Date()
			firstYear := time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC)
			lastYear := time.Date(y, 12, 1, 0, 0, 0, 0, time.UTC)
			pipe := []tk.M{
				tk.M{
					"$match": tk.M{
						"businessmetric": kolom2,
						"period":         tk.M{"$gte": firstYear, "$lte": lastYear},
					},
				},
				tk.M{
					"$group": tk.M{
						"_id": tk.M{
							"businessmetric": "$businessmetric",
							"period":         "$period",
						},
					},
				},
				tk.M{
					"$sort": tk.M{
						"_id.period": 1,
					},
				},
			}
			csr, e := c.Ctx.Connection.NewQuery().Command("pipe", pipe).From("BusinessMetricsData").Cursor(nil)
			defer csr.Close()
			e = csr.Fetch(&d, 0, false)
			if e != nil {
				return c.ErrorResultInfo(e.Error(), nil)
			}

			if len(d) > 0 {
				monthExist := []string{}
				for _, data := range d {
					r := data.Get("_id").(tk.M)
					monthid := 1
					for colid, subrow := range row.Cells {
						if colid > 5 && colid < 18 {
							periodvalue, _ := subrow.String()
							perioddata := time.Month(monthid).String() + " " + tk.ToString(y)
							if perioddata == r.Get("period").(time.Time).Format("January 2006") && !tk.IsNilOrEmpty(periodvalue) {
								monthExist = append(monthExist, r.Get("period").(time.Time).Format("Jan 2006"))

								exist.Set(kolom2, monthExist)
								isExist = true
							}
							monthid++
						}
					}
				}
			}
		}
	}

	//remove file
	err = os.Remove(filepath.Join(c.UploadMetricPath, payload.GetString("file")))
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	if isExist {
		bmane := ""
		for key, v := range exist {
			date := "["
			for i, dateexist := range v.([]string) {
				date += dateexist
				cekkosong := i + 1
				if cekkosong < len(v.([]string)) {
					date += ", "
				}
			}
			date += "]"
			bmane += key + " - " + date + "\n"
		}
		return c.SetResultInfo(false, bmane+" had been uploaded, are you sure you want to upload again?", tk.M{"exist": exist, "payload": payload})
	}

	return c.SetResultInfo(false, "", tk.M{"payload": payload})
}

func (c *ScorecardController) UploadMetricsTemplate(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	c.Action(k, "Scorecard", "Upload Metrics Data", "", "", "", "", "")
	k.Config.OutputType = knot.OutputJson
	res := ResultInfo{}

	err, payload := c.uploadTemplate(k)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	result := make([]*BusinessDriverL1Model, 0)
	csr, err := c.Ctx.Connection.
		NewQuery().
		From(new(BusinessDriverL1Model).TableName()).
		Cursor(nil)
	if csr != nil {
		defer csr.Close()
	}
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	err = csr.Fetch(&result, 0, true)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	// tk.Println(result)

	err, tempData := c.multiSheetProcess(k, payload, result, c.UploadMetricPath)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	resultReMapping := c.ReMappingSummaryBD(k).(ResultInfo)
	if resultReMapping.IsError {
		return resultReMapping
	}

	res.Data = tempData
	return res
}

func (c *ScorecardController) uploadTemplate(k *knot.WebContext) (error, tk.M) {
	var uploadPath string
	payloadAsObject := tk.M{}
	// config := helper.ReadConfig()

	// if value, ok := config["metricFilePath"]; ok {
	// 	uploadPath = value
	// } else {
	// 	wd, _ := os.Getwd()
	// 	uploadPath = filepath.Join(wd, "metricfiles")
	// }
	uploadPath = c.UploadMetricPath
	reader, err := k.Request.MultipartReader()
	if err != nil {
		return err, payloadAsObject
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		key := part.FormName()
		if key == "date" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)

			payloadAsObject.Set(key, buf.String())
		} else if key == "file" {
			randomString := tk.RandomString(20)
			extension := filepath.Ext(part.FileName())
			newFileName := fmt.Sprintf("%s%s", randomString, extension)
			payloadAsObject.Set(key, newFileName)
			payloadAsObject.Set("original-file", part.FileName())
			// tk.Println(uploadPath)
			fileLocation := filepath.Join(uploadPath, newFileName)
			dst, err := os.Create(fileLocation)
			if dst != nil {
				defer dst.Close()
			}
			if err != nil {
				return err, payloadAsObject
			}

			if _, err := io.Copy(dst, part); err != nil {
				return err, payloadAsObject
			}
		} else if key == "downloadOption" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)

			payloadAsObject.Set(key, buf.String())
		}
	}

	return nil, payloadAsObject
}

func (c *ScorecardController) ExtractMetrics(k *knot.WebContext, BMName string, metricId string, MetricType string, fileName, downloadOption string) (tk.Ms, error) {
	// tk.Println("DOWNLOAD")
	DataList := tk.Ms{} //make([]*BusinessMetricsDataModel, 0)
	//extract only selected metrics
	// tk.Println("Starting to extract ", metricId)
	k.Config.OutputType = knot.OutputJson

	regionCol := make([]tk.M, 0)
	getMajorRegion := map[string]string{}
	getRegion := map[string]string{}
	getCountry := map[string]string{}
	// country := map[string]string{}

	csr, e := c.Ctx.Connection.NewQuery().From("Region").Cursor(nil)
	defer csr.Close()
	e = csr.Fetch(&regionCol, 0, true)
	if e != nil {
		return DataList, e
	}

	for _, i := range regionCol {
		getCountry[i.GetString("CountryCode")] = i.GetString("Country")
		getRegion[i.GetString("Country")] = i.GetString("Region")
		getMajorRegion[i.GetString("Country")] = i.GetString("Major_Region")
	}

	var fileLocation string
	fileLocation = c.UploadMetricPath
	// config := helper.ReadConfig()

	// if value, ok := config["metricFilePath"]; ok {
	// 	fileLocation = value
	// } else {
	// 	wd, _ := os.Getwd()
	// 	fileLocation = filepath.Join(wd, "metricfiles")
	// }
	// tk.Println("Get Data Master ", metricId)
	masterBD := new(BusinessDriverL1Model)
	// tk.Println("Metric Id ", metricId)
	// tk.Println("Metric type ", MetricType)
	csr, e = c.Ctx.Find(new(BusinessDriverL1Model), tk.M{}.Set("where", dbox.Eq("businessmetric.id", metricId))) //
	e = csr.Fetch(&masterBD, 1, false)
	if e != nil {
		// tk.Println("Err")
		// tk.Println(e.Error())
		// return c.ErrorResultInfo(e.Error(), nil)
		return DataList, e
	}
	csr.Close()
	var fileInfo = NewMetricFile()
	var bm = BusinessMetric{}
	for _, mtc := range masterBD.BusinessMetric {
		if mtc.Id == metricId {
			bm = mtc
			for _, fl := range mtc.MetricFiles {
				if fl.FileName == fileName {
					fileInfo = &fl
				}
			}
		}
	}

	year := 1990
	month := 1
	if fileInfo.MonthYear != 0 {
		fileLocation = filepath.Join(fileLocation, fileInfo.FileName)
		my := tk.ToString(fileInfo.MonthYear)
		year = tk.ToInt(my[0:4], "")
		month = tk.ToInt(my[4:], "")
	}

	// tk.Println("File Location >>>> ", fileLocation)
	exelFile, e := xl.OpenFile(fileLocation)
	if e != nil {
		return DataList, e
	}

	// datas := []BusinessMetricsDataModel{}
	var shit = exelFile.Sheets[0]
	for idxRow, row := range shit.Rows {
		if idxRow == 0 || len(row.Cells) < 3 {
			continue
		}

		businessMetric, countryCode, country := "", "", ""
		data := tk.M{} //NewBusinessMetricsDataModel()
		scorecardCategory, _ := row.Cells[0].String()

		// 130895111188 = Magic Number to set NA
		if downloadOption == "owner" {
			businessMetric, _ = row.Cells[2].String()
			country, _ = row.Cells[4].String()
			countryCode, _ = row.Cells[5].String()
			if country == "" {
				country = getCountry[countryCode]
			}
			// metricsDenomination, _ := row.Cells[6].String()

			ytdActual, _ := row.Cells[8].Float()
			data.Set("ActualYTD", ytdActual)
			tempActualYtdData, _ := row.Cells[8].String()
			tempActualYtdData = strings.Trim(tempActualYtdData, " ")
			if tempActualYtdData == "" || tempActualYtdData == "NaN" {
				data.Set("ActualYTD", 130895111188)
			}

			rag := ""
			for idx, v := range row.Cells {
				if idx == 9 {
					rag, _ = v.String()
				}
			}
			if strings.ToLower(rag) == "r" {
				rag = "red"
			} else if strings.ToLower(rag) == "a" {
				rag = "amber"
			} else if strings.ToLower(rag) == "g" {
				rag = "green"
			}
			data.Set("RAG", strings.ToLower(rag))
		} else {
			businessMetric, _ = row.Cells[1].String()
			country, _ = row.Cells[3].String()
			countryCode, _ = row.Cells[4].String()
			if country == "" {
				country = getCountry[countryCode]
			}

			baseline, _ := row.Cells[5].Float()
			data.Set("Baseline", baseline)
			tempBaselineData, _ := row.Cells[5].String()
			tempBaselineData = strings.Trim(tempBaselineData, " ")
			if tempBaselineData == "" || tempBaselineData == "NaN" {
				data.Set("Baseline", 130895111188)
			}

			/// Jan - Dec
			i := 1
			for idx, v := range row.Cells {
				if idx > 5 && idx < 18 {
					value, _ := v.String()
					currentyear, _, _ := time.Now().Date()
					tmp := time.Date(currentyear, time.Month(i), 1, 0, 0, 0, 0, time.Time{}.Location()).Format("Jan2006")

					data.Set(tmp, tk.ToFloat64(value, 2, tk.RoundingAuto))
					tempMonthData, _ := v.String()
					tempMonthData = strings.Trim(tempMonthData, " ")
					if tempMonthData == "" {
						data.Set(tmp, 130895111188)
					}
					i++
				}
			}

			yearlyTarget, _ := row.Cells[18].Float()
			data.Set("Target", yearlyTarget)
			tempTargetData, _ := row.Cells[18].String()
			tempTargetData = strings.Trim(tempTargetData, " ")
			if tempTargetData == "" {
				data.Set("Target", 130895111188)
			}
		}

		data.Set("Period", time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC))
		data.Set("SCId", masterBD.Id)
		data.Set("BMId", bm.Id)
		data.Set("Year", year)
		BusinessName, _ := row.Cells[0].String()
		data.Set("BusinessName", BusinessName)
		data.Set("DecimalFormat", bm.DecimalFormat)
		data.Set("MetricType", bm.MetricType)
		data.Set("ScorecardCategory", scorecardCategory)
		data.Set("BusinessMetric", businessMetric)
		data.Set("BusinessMetricDescription", bm.Description)
		data.Set("CountryCode", countryCode)
		if data.GetString("CountryCode") == "GLOBAL" {
			data.Set("MajorRegion", "GLOBAL")
			data.Set("Region", "GLOBAL")
			data.Set("Country", "GLOBAL")
		} else {
			reg := getRegion[country]
			if reg != "" {
				data.Set("MajorRegion", getMajorRegion[country])
				data.Set("Region", reg)
				data.Set("Country", country)
			}
		}

		data.Set("CreatedDate", time.Now())
		data.Set("CreatedBy", k.Session("username", "admin").(string))
		data.Set("UpdatedDate", time.Now())
		data.Set("UpdatedBy", k.Session("username", "admin").(string))
		//if data selected is not the same with excel, it will ignore this row
		// tk.Println("ngeprint rows<<<<<<<<<<<<", tk.JsonString(data))
		bmDescription := strings.Replace(bm.Description, "\n", "", -1)
		if strings.Trim(bmDescription, " ") != strings.Trim(businessMetric, " ") {
			if strings.Trim(bm.DataPoint, " ") != strings.Trim(businessMetric, " ") {
				continue
			}
		}

		// Handling Region and Country for Total Data
		if data.GetString("Region") == "" && data.GetString("Country") == "" {
			data.Set("MajorRegion", data.GetString("CountryCode"))
			data.Set("Region", data.GetString("CountryCode"))
			data.Set("Country", data.GetString("CountryCode"))
		}

		if MetricType == "PERCENTAGE" {
			if downloadOption == "owner" {
				if data.GetFloat64("ActualYTD") == 130895111188 {
					data.Set("ActualYTD", data.GetFloat64("ActualYTD"))
				} else {
					data.Set("ActualYTD", data.GetFloat64("ActualYTD")*100)
				}
			} else {
				if data.GetFloat64("Baseline") == 130895111188 {
					data.Set("Baseline", data.GetFloat64("Baseline"))
				} else {
					data.Set("Baseline", data.GetFloat64("Baseline")*100)
				}

				/// Jan - Dec
				for i := 1; i < 13; i++ {
					currentyear, _, _ := time.Now().Date()
					tmp := time.Date(currentyear, time.Month(i), 1, 0, 0, 0, 0, time.Time{}.Location()).Format("Jan2006")
					if data.GetFloat64(tmp) == 130895111188 {
						data.Set(tmp, data.GetFloat64(tmp))
					} else {
						data.Set(tmp, data.GetFloat64(tmp)*100)
					}
				}

				if data.GetFloat64("Target") == 130895111188 {
					data.Set("Target", data.GetFloat64("Target"))
				} else {
					data.Set("Target", data.GetFloat64("Target")*100)
				}
			}
		}
		// if tempActualData != "" {
		DataList = append(DataList, data)
	}

	return DataList, e

}

func (c *ScorecardController) multiSheetProcess(k *knot.WebContext, payloadAsObject tk.M, result []*BusinessDriverL1Model, uploadPath string) (error, tk.Ms) {
	tempData := tk.Ms{}
	exelFile, e := xl.OpenFile(filepath.Join(uploadPath, payloadAsObject.GetString("file")))
	if e != nil {
		return e, tempData
	}

	i := 0
	dataexist := ""
	for idx, row := range exelFile.Sheets[0].Rows {
		if idx == 0 || len(row.Cells) < 3 {
			continue
		}

		if payloadAsObject.GetString("downloadOption") == "owner" {
			dateUpload, _ := row.Cells[1].String()
			periodUpload, _ := time.Parse("200601", payloadAsObject.GetString("date"))
			if dateUpload != periodUpload.Format("Jan 2006") {
				return tk.Error("Date upload not identical with data month file"), nil
			}
		}

		kolom2 := ""
		for idx, _ := range row.Cells {
			if payloadAsObject.GetString("downloadOption") == "owner" {
				if idx == 2 {
					kolom2, _ = row.Cells[idx].String()
				}
			} else {
				if idx == 1 {
					kolom2, _ = row.Cells[idx].String()
				}
			}
		}

		// kolom2, _ := row.Cells[2].String()
		for _, d := range result {
			for j, ds := range d.BusinessMetric {
				if strings.TrimSpace(ds.Description) == kolom2 {
					if dataexist != kolom2 {
						dataexist = kolom2

						o := new(MetricFile)
						o.MonthYear = tk.ToInt(payloadAsObject.GetString("date"), tk.RoundingAuto)
						o.FileName = payloadAsObject.GetString("file")
						o.OriginalFileName = payloadAsObject.GetString("original-file")
						o.UploadedAt = time.Now()
						o.UploaderName = k.Session("username", "admin").(string)

						DataPoint := ds.DataPoint
						MetricType := ds.MetricType
						d.BusinessMetric[j].MetricFiles = append(d.BusinessMetric[j].MetricFiles, *o)

						err := c.Ctx.Save(d)
						if err != nil {
							return err, tempData
						}
						temp := tk.M{}
						err = tk.StructToM(o, &temp)
						if err != nil {
							return err, tempData
						}

						DataList, err := c.ExtractMetrics(k, DataPoint, ds.Id, MetricType, o.FileName, payloadAsObject.GetString("downloadOption"))
						if err != nil {
							return err, tempData
						}

						if payloadAsObject.GetString("downloadOption") == "finance" {
							var months = []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
							tmpdate := []string{}
							for i := 0; i < 12; i++ {
								currentyear, _, _ := time.Now().Date()
								tmp := time.Date(currentyear, months[i], 1, 0, 0, 0, 0, time.Time{}.Location()).Format("Jan2006")
								tmpdate = append(tmpdate, tmp)
							}
							temp.Set("finCol", tmpdate)
						}

						temp.Set("DataList", DataList)
						temp.Set("MetricsID", ds.Id)
						temp.Set("DataPoint", DataPoint)
						temp.Set("MetricType", MetricType)
						tempData = append(tempData, temp)
					}
				}
			}
		}
		i++
	}

	return nil, tempData
}

func (c *ScorecardController) ExportExcelMBR(k *knot.WebContext) interface{} {
	// c.LoadBase(k)
	c.Action(k, "Scorecard", "Save Scorecard to Excel", "", "", "", "", "")

	k.Config.OutputType = knot.OutputJson
	parm := struct {
		SCData      string
		RegionData  []string
		CountryData []string
		Data        []tk.M
	}{}
	err := k.GetPayload(&parm)

	if err != nil {
		tk.Println(err.Error())
	}
	// tk.Println(parm)
	err = helper.Deserialize(parm.SCData, &parm.Data)
	// tk.Println("Deserializee..")
	if err != nil {
		tk.Println(err.Error())
	}

	var file *xl.File
	var sheet *xl.Sheet
	var row *xl.Row
	var cell *xl.Cell

	style := xl.NewStyle()
	font := xl.NewFont(11, "calibri")
	fontBold := xl.NewFont(11, "calibri")
	fontBold.Bold = true
	fontKeyMetrics := xl.NewFont(11, "calibri")
	// fontKeyMetrics.Bold = true
	// fontKeyMetrics.Color = "005c84"

	border := xl.NewBorder("thin", "thin", "thin", "thin")
	topBorderOnly := xl.NewBorder("", "", "thin", "")
	bottomBorderOnly := xl.NewBorder("", "", "", "thin")
	leftrightBorderOnly := xl.NewBorder("thin", "thin", "", "")
	topLeftRightBorderOnly := xl.NewBorder("thin", "thin", "thin", "")

	style.Font = *font
	style.Border = *border

	styleHeader := xl.NewStyle()
	styleHeader.Font = *fontBold
	styleHeader.Alignment.WrapText = true
	styleHeader.Alignment.Horizontal = "center"
	styleHeader.Alignment.Vertical = "center"
	styleHeader.Border = *border

	styleKeyMetrics := xl.NewStyle()
	styleKeyMetrics.Font = *fontKeyMetrics
	styleKeyMetrics.Alignment.Horizontal = "left"
	styleKeyMetrics.Alignment.Vertical = "center"
	styleKeyMetrics.Border = *leftrightBorderOnly
	styleKeyMetrics.Fill = *xl.NewFill("solid", "ffffff", "ffffff")

	styleFirstKeyMetrics := xl.NewStyle()
	styleFirstKeyMetrics.Font = *fontKeyMetrics
	styleFirstKeyMetrics.Alignment.Horizontal = "left"
	styleFirstKeyMetrics.Alignment.Vertical = "center"
	styleFirstKeyMetrics.Border = *topBorderOnly
	styleFirstKeyMetrics.Fill = *xl.NewFill("solid", "ffffff", "ffffff")

	styleLastKeyMetrics := xl.NewStyle()
	styleLastKeyMetrics.Font = *fontKeyMetrics
	styleLastKeyMetrics.Alignment.Horizontal = "left"
	styleLastKeyMetrics.Alignment.Vertical = "center"
	styleLastKeyMetrics.Border = *bottomBorderOnly
	styleLastKeyMetrics.Fill = *xl.NewFill("solid", "ffffff", "ffffff")

	lastRowBorder := xl.NewStyle()
	lastRowBorder.Font = *font
	lastRowBorder.Border = *topBorderOnly

	styleScorecard := xl.NewStyle()
	styleScorecard.Font = *font
	styleScorecard.Alignment.Horizontal = "left"
	styleScorecard.Alignment.Vertical = "center"
	styleScorecard.Border = *border

	file = xl.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		c.ErrorResultInfo(err.Error(), nil)
	}

	row = sheet.AddRow()
	row.Height = 5
	row = sheet.AddRow()
	row.Height = 5

	row = sheet.AddRow()
	row.AddCell()
	cell = row.AddCell()
	cell.Value = "Framework"
	cell.Merge(0, 1)
	cell.SetStyle(styleHeader)
	// cell.SetStyle(styleHeader)

	cell = row.AddCell()
	cell.Value = "Key Metrics"
	cell.Merge(0, 1)
	cell.SetStyle(styleHeader)
	// cell.SetStyle(styleHeader)

	cell = row.AddCell()
	cell.Value = "GLOBAL"
	cell.Merge(3, 0)
	cell.SetStyle(styleHeader)
	// cell.SetStyle(styleHeader)

	cell = row.AddCell()
	cell.Value = ""
	cell.SetStyle(styleHeader)
	cell = row.AddCell()
	cell.Value = ""
	cell.SetStyle(styleHeader)

	for _, region := range parm.RegionData {
		cell = row.AddCell()
		cell.Value = ""
		cell.SetStyle(styleHeader)
		cell = row.AddCell()
		cell.Value = region
		cell.Merge(1, 0)
		cell.SetStyle(styleHeader)
	}

	for _, country := range parm.CountryData {
		cell = row.AddCell()
		cell.Value = ""
		cell.SetStyle(styleHeader)
		cell = row.AddCell()
		cell.Value = country
		cell.Merge(1, 0)
		cell.SetStyle(styleHeader)
	}

	cell = row.AddCell()
	cell.Value = ""
	cell.SetStyle(styleHeader)
	cell = row.AddCell()
	cell.Value = "Last Updated"
	cell.Merge(0, 1)
	cell.SetStyle(styleHeader)

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = ""
	// cell.SetStyle(styleHeader)
	cell = row.AddCell()
	cell.Value = ""
	cell.SetStyle(styleHeader)
	cell = row.AddCell()
	cell.Value = ""
	cell.SetStyle(styleHeader)
	cell = row.AddCell()
	cell.Value = "2016 Baseline"
	cell.SetStyle(styleHeader)

	cell = row.AddCell()
	cell.Value = "2017 Target"
	cell.SetStyle(styleHeader)

	cell = row.AddCell()
	cell.Value = "2017 YTD"
	cell.SetStyle(styleHeader)

	cell = row.AddCell()
	cell.Value = "YTD Actual vs YTD Budget B/(W)"
	cell.SetStyle(styleHeader)

	for _, _ = range parm.RegionData {
		cell = row.AddCell()
		cell.Value = "2017 YTD"
		cell.SetStyle(styleHeader)

		cell = row.AddCell()
		cell.Value = "YTD Actual vs YTD Budget B/(W)"
		cell.SetStyle(styleHeader)
	}

	for _, _ = range parm.CountryData {
		cell = row.AddCell()
		cell.Value = "2017 YTD"
		cell.SetStyle(styleHeader)

		cell = row.AddCell()
		cell.Value = "YTD Actual vs YTD Budget B/(W)"
		cell.SetStyle(styleHeader)
	}
	cell = row.AddCell()
	cell.SetStyle(styleHeader)

	result := []tk.M{}
	result = parm.Data

	for _, ii := range result {

		row = sheet.AddRow()
		row.AddCell()
		cell = row.AddCell()
		cell.Value = ii.GetString("Name")
		BusinessMetrics := ii.Get("BusinessMetric").([]interface{})
		cell.Merge(0, len(BusinessMetrics)-1)
		cell.SetStyle(styleScorecard)
		for index, a := range BusinessMetrics {
			styleValue := xl.NewStyle()
			styleValue.Font = *font
			styleValue.Alignment.Horizontal = "center"
			styleValue.Alignment.Vertical = "center"
			styleValue.Fill = *xl.NewFill("solid", "ffffff", "ffffff")

			styleFillRed := xl.NewStyle()
			styleFillRed.Font = *font
			styleFillRed.Fill = *xl.NewFill("solid", "f74e4e", "f74e4e")
			styleFillRed.Alignment.Horizontal = "center"
			styleFillRed.Alignment.Vertical = "center"

			styleFillAmber := xl.NewStyle()
			styleFillAmber.Font = *font
			styleFillAmber.Fill = *xl.NewFill("solid", "FFD24D", "FFD24D")
			styleFillAmber.Alignment.Horizontal = "center"
			styleFillAmber.Alignment.Vertical = "center"

			styleFillGreen := xl.NewStyle()
			styleFillGreen.Font = *font
			styleFillGreen.Fill = *xl.NewFill("solid", "6AC17B", "6AC17B")
			styleFillGreen.Alignment.Horizontal = "center"
			styleFillGreen.Alignment.Vertical = "center"

			styleFillDefault := xl.NewStyle()
			styleFillDefault.Font = *font
			styleFillDefault.Fill = *xl.NewFill("solid", "ffffff", "ffffff")
			styleFillDefault.Alignment.Horizontal = "center"
			styleFillDefault.Alignment.Vertical = "center"

			if index == 0 {
				styleValue.Border = *topLeftRightBorderOnly
				styleFillRed.Border = *topLeftRightBorderOnly
				styleFillAmber.Border = *topLeftRightBorderOnly
				styleFillGreen.Border = *topLeftRightBorderOnly
				styleFillDefault.Border = *topLeftRightBorderOnly
			} else if index == (len(BusinessMetrics) - 1) {
				styleValue.Border = *bottomBorderOnly
				styleFillRed.Border = *bottomBorderOnly
				styleFillAmber.Border = *bottomBorderOnly
				styleFillGreen.Border = *bottomBorderOnly
				styleFillDefault.Border = *bottomBorderOnly
			} else {
				styleValue.Border = *leftrightBorderOnly
				styleFillRed.Border = *leftrightBorderOnly
				styleFillAmber.Border = *leftrightBorderOnly
				styleFillGreen.Border = *leftrightBorderOnly
				styleFillDefault.Border = *leftrightBorderOnly
			}

			b := a.(map[string]interface{})
			if index >= 1 {
				row = sheet.AddRow()
				row.AddCell()
				cell = row.AddCell()
				cell.SetStyle(styleValue)
			}

			cell = row.AddCell()
			cell.Value = b["Description"].(string)
			if index == 0 {
				cell.SetStyle(styleFirstKeyMetrics)
			} else if index == (len(BusinessMetrics) - 1) {
				cell.SetStyle(styleLastKeyMetrics)
			} else {
				cell.SetStyle(styleKeyMetrics)
			}
			NABaseline := b["NABaseline"].(bool)
			NATarget := b["NATarget"].(bool)
			NAActual := b["NAActual"].(bool)
			NABudget := b["NABudget"].(bool)
			DecimalFormat := b["DecimalFormat"].(string)
			MetricType := b["MetricType"].(string)
			floatFormat := "#,##0"
			devider := 1.0
			switch MetricType {
			case "PERCENTAGE":
				devider = 100
				if DecimalFormat == "0" {
					floatFormat = "#0%"
				} else if DecimalFormat == "1" {
					floatFormat = "#0.0%"
				} else {
					floatFormat = "#0.00%"
				}
				break
			case "NUMERIC":
				devider = 1
				if DecimalFormat == "0" {
					floatFormat = "#,##0"
				} else if DecimalFormat == "1" {
					floatFormat = "#,##0.0"
				} else {
					floatFormat = "#,##0.00"
				}
			default:
				devider = 1
				if DecimalFormat == "0" {
					floatFormat = "#,##0"
				} else if DecimalFormat == "1" {
					floatFormat = "#,##0.0"
				} else {
					floatFormat = "#,##0.00"
				}
				break
			}

			cell = row.AddCell()
			if NABaseline {
				cell.SetValue(nil)
			} else {
				Baseline, isOk := b["BaseLineValue"].(float64)
				if Baseline == 130895111188 || !isOk {
					cell.SetValue(nil)
				} else {
					cell.SetFloatWithFormat(Baseline/devider, floatFormat)
				}
			}
			cell.SetStyle(styleValue)

			cell = row.AddCell()
			if NATarget {
				cell.SetValue(nil)
			} else {
				Target, isOk := b["TargetValue"].(float64)
				if Target == 130895111188 || !isOk {
					cell.SetValue(nil)
				} else {
					cell.SetFloatWithFormat(Target/devider, floatFormat)
				}
			}
			cell.SetStyle(styleValue)

			cell = row.AddCell()
			if NAActual {
				cell.SetValue(nil)
			} else {
				ActualYTD, isOk := b["CurrentYTDValue"].(float64)
				if ActualYTD == 130895111188 || !isOk {
					cell.SetValue(nil)
				} else {
					cell.SetFloatWithFormat(ActualYTD/devider, floatFormat)
				}
			}

			rag := b["Display"].(string)
			switch rag {
			case "red":
				cell.SetStyle(styleFillRed)
				break
			case "amber":
				cell.SetStyle(styleFillAmber)
				break
			case "green":
				cell.SetStyle(styleFillGreen)
				break
			default:
				cell.SetStyle(styleFillDefault)
				break
			}

			cell = row.AddCell()
			if NABudget {
				cell.SetValue(nil)
			} else {
				Budget, isOk := b["CurrentYTDValueVsBudget"].(float64)
				if Budget == 130895111188 || !isOk {
					cell.SetValue(nil)
				} else {
					cell.SetFloatWithFormat(Budget/devider, floatFormat)
				}
			}
			cell.SetStyle(styleValue)

			RegionalData := b["RegionalData"].([]interface{})
			for _, rd := range RegionalData {
				rdata := rd.(map[string]interface{})
				NAActual := rdata["NAActual"].(bool)
				NABudget := rdata["NABudget"].(bool)
				cell = row.AddCell()
				if NAActual {
					cell.SetValue(nil)
				} else {
					ActualYTD, isOk := rdata["CurrentYTDValue"].(float64)
					if ActualYTD == 130895111188 || !isOk {
						cell.SetValue(nil)
					} else {
						cell.SetFloatWithFormat(ActualYTD/devider, floatFormat)
					}
				}
				rag := rdata["Rag"].(string)
				switch rag {
				case "red":
					cell.SetStyle(styleFillRed)
					break
				case "amber":
					cell.SetStyle(styleFillAmber)
					break
				case "green":
					cell.SetStyle(styleFillGreen)
					break
				default:
					cell.SetStyle(styleFillDefault)
					break
				}

				cell = row.AddCell()
				if NABudget {
					cell.SetValue(nil)
				} else {
					Budget, isOk := rdata["CurrentYTDValueVsBudget"].(float64)
					if Budget == 130895111188 || !isOk {
						cell.SetValue(nil)
					} else {
						cell.SetFloatWithFormat(Budget/devider, floatFormat)
					}
				}
				cell.SetStyle(styleValue)
			}

			CountryData := b["CountryData"].([]interface{})
			for _, cd := range CountryData {
				cdata := cd.(map[string]interface{})
				NAActual := cdata["NAActual"].(bool)
				NABudget := cdata["NABudget"].(bool)
				cell = row.AddCell()
				if NAActual {
					cell.SetValue(nil)
				} else {
					ActualYTD, isOk := cdata["CurrentYTDValue"].(float64)
					if ActualYTD == 130895111188 || !isOk {
						cell.SetValue(nil)
					} else {
						cell.SetFloatWithFormat(ActualYTD/devider, DecimalFormat)
					}
				}
				rag := cdata["Rag"].(string)
				switch rag {
				case "red":
					cell.SetStyle(styleFillRed)
					break
				case "amber":
					cell.SetStyle(styleFillAmber)
					break
				case "green":
					cell.SetStyle(styleFillGreen)
					break
				default:
					cell.SetStyle(styleFillDefault)
					break
				}

				cell = row.AddCell()
				if NABudget {
					cell.SetValue("")
				} else {
					Budget, isOk := cdata["CurrentYTDValueVsBudget"].(float64)
					if Budget == 130895111188 || !isOk {
						cell.SetValue("N/A")
					} else {
						cell.SetFloatWithFormat(Budget/devider, DecimalFormat)
					}
				}
				cell.SetStyle(styleValue)
			}

			cell = row.AddCell()
			cell.SetValue(b["LastUpdate"].(string))
			cell.SetStyle(styleValue)

		}

	}

	// Add Last Row to Fix Border
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = ""
	cell = row.AddCell()
	cell.SetStyle(lastRowBorder)
	cell = row.AddCell()
	cell.SetStyle(lastRowBorder)
	cell = row.AddCell()
	cell.SetStyle(lastRowBorder)

	cell = row.AddCell()
	cell.SetStyle(lastRowBorder)

	cell = row.AddCell()
	cell.SetStyle(lastRowBorder)

	cell = row.AddCell()
	cell.SetStyle(lastRowBorder)

	for _, _ = range parm.RegionData {
		cell = row.AddCell()
		cell.SetStyle(lastRowBorder)

		cell = row.AddCell()
		cell.SetStyle(lastRowBorder)
	}

	for _, _ = range parm.CountryData {
		cell = row.AddCell()
		cell.SetStyle(lastRowBorder)

		cell = row.AddCell()
		cell.SetStyle(lastRowBorder)
	}
	cell = row.AddCell()
	cell.SetStyle(lastRowBorder)

	sheet.SetColWidth(0, 0, 2)
	sheet.SetColWidth(1, 1, 20)
	sheet.SetColWidth(2, 2, 30)

	ExcelFilename := "CB_Operational Scorecard.xlsx"
	err = file.Save(c.TemplatePath + "/" + ExcelFilename)
	// tk.Println("Excel Exported ")
	if err != nil {
		tk.Println(err)
		c.ErrorResultInfo(err.Error(), nil)
	}

	return c.SetResultInfo(false, "", ExcelFilename)
}
