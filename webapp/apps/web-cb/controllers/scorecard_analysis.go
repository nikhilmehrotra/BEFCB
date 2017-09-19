package controllers

import (
	// "eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"math"
	"time"
	// "regexp"
	// "bytes"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// xl "github.com/tealeg/xlsx"
	// "gopkg.in/mgo.v2/bson"
	// "io"
	// "os"
	// "path/filepath"
	// "strconv"
	// "strings"
	// "unicode"
	"sort"
)

type ScorecardAnalysisController struct {
	*BaseController
}

type SortGAP []tk.M

func (a SortGAP) Len() int           { return len(a) }
func (a SortGAP) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortGAP) Less(i, j int) bool { return a[i].GetFloat64("Gap") > a[j].GetFloat64("Gap") }

type SortActual []tk.M

func (a SortActual) Len() int           { return len(a) }
func (a SortActual) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortActual) Less(i, j int) bool { return a[i].GetFloat64("Actual") < a[j].GetFloat64("Actual") }

type SortActualDescending []tk.M

func (a SortActualDescending) Len() int      { return len(a) }
func (a SortActualDescending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortActualDescending) Less(i, j int) bool {
	return a[i].GetFloat64("Actual") > a[j].GetFloat64("Actual")
}

type ParmScorecardAnalysis struct {
	BusinessMetrics string
	Region          string
	Country         string
	Period          string
}

func (c *ScorecardAnalysisController) GetData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		BusinessMetrics string
		Region          string
		Country         string
		Period          string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	// tk.Println("PARM.BMID", parm.BusinessMetrics)

	result := tk.M{}
	result.Set("CountryAnalysis", c.GetCountryAnalysisData(k, parm.BusinessMetrics, parm.Period))
	// result.Set("CountryAnalysis2", c.GetCountryAnalysisData2(k, parm.BusinessMetrics, parm.Period))
	result.Set("FullYearRanking", c.GetFullYearRankingData(k, parm))
	return c.SetResultInfo(false, "", result)
}
func (c *ScorecardAnalysisController) GetCountryAnalysisData(k *knot.WebContext, BusinessMetrics string, Periods string) interface{} {

	k.Config.OutputType = knot.OutputJson

	csr, err := c.Ctx.Find(new(BusinessDriverL1Model), tk.M{}.Set("where", db.Eq("businessmetric.id", BusinessMetrics)))
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	BDL1 := make([]BusinessDriverL1Model, 0)
	err = csr.Fetch(&BDL1, 0, false)
	if err != nil {
		return err.Error()
	}
	csr.Close()

	csr, err = c.Ctx.Connection.NewQuery().From("Region").
		Cursor(nil)
	if err != nil {
		return err.Error()
	}
	CountryList := []tk.M{}
	err = csr.Fetch(&CountryList, 0, false)
	if err != nil {
		return err.Error()
	}
	csr.Close()

	csr, err = c.Ctx.Connection.NewQuery().From("Region").Group("Major_Region").
		Cursor(nil)
	if err != nil {
		return err.Error()
	}
	RegionList := []tk.M{}
	err = csr.Fetch(&RegionList, 0, false)
	if err != nil {
		return err.Error()
	}
	csr.Close()

	resultArr := []tk.M{}
	resultLastYear := []tk.M{}

	Period, _ := time.Parse("20060102", Periods)
	Period = Period.UTC()

	LastYear := Period.Year() - 1
	LastYearPeriod := time.Date(LastYear, Period.Month(), 1, 0, 0, 0, 0, time.UTC)

	for _, o := range CountryList {
		result := tk.M{}
		// RAG := ""

		CountryName := o.Get("Country").(string)
		RegionName := o.Get("Major_Region").(string)
		if CountryName != RegionName {
			csr, err = c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Where(db.And(db.Eq("country", CountryName), db.Eq("period", Period), db.Eq("bmid", BusinessMetrics))).Order("-period").Cursor(nil)
			if err != nil {
				return err.Error()
			}
			BusinessMetricsDataArr := []BusinessMetricsDataModel{}
			err = csr.Fetch(&BusinessMetricsDataArr, 1, false)
			if err != nil {
				return err.Error()
			}
			csr.Close()

			if len(BusinessMetricsDataArr) > 0 {
				oo := BusinessMetricsDataArr[0]
				higherIsBetter := false
				for _, oi := range BDL1[0].BusinessMetric {
					if oi.Id == BusinessMetrics {
						target, actualytd, remainingBudget, remainingBudgetOpposite := 0.0, 0.0, 0.0, 0.0
						// if oo.Baseline != 130895111188 {
						// 	baseline = oo.Baseline
						// }
						if oo.Target != 130895111188 {
							target = oo.Target
						}
						if oo.ActualYTD != 130895111188 {
							actualytd = oo.ActualYTD
						}

						if oo.RemainingBudget != 130895111188 {
							remainingBudget = oo.RemainingBudget
						}
						if oo.RemainingBudgetOpposite != 130895111188 {
							remainingBudgetOpposite = oo.RemainingBudgetOpposite
						}

						if oi.ValueType == 0 {
							findActual := tk.Div(target, actualytd)
							if target > 0 && actualytd < 0 {
								findActual = tk.Div((actualytd - target), actualytd)
							}
							calcBudget := actualytd - remainingBudgetOpposite
							findBudget := tk.Div(calcBudget, target)

							result.Set("ActualPercent", findActual)
							result.Set("NAActualPercent", false)
							if math.IsNaN(findActual) {
								result.Set("ActualPercent", 0)
								result.Set("NAActualPercent", true)
							}

							projection := tk.Div(findActual, float64(oo.Period.Month())) * 12
							result.Set("Projection", projection)
							result.Set("NAProjection", false)
							if math.IsNaN(findActual) {
								result.Set("Projection", 0)
								result.Set("NAProjection", true)
							}
							// budgetProjection := tk.Div(findBudget, float64(oo.Period.Month())) * 12
							result.Set("BudgetProjection", findBudget)
							result.Set("NABudgetProjection", false)
							if math.IsNaN(findBudget) {
								result.Set("BudgetProjection", 0)
								result.Set("NABudgetProjection", true)
							}

						} else {
							findActual := tk.Div(actualytd, target)
							if target < 0 && actualytd > 0 {
								findActual = tk.Div((target - actualytd), target)
							}
							calcBudget := actualytd + remainingBudget
							findBudget := tk.Div(calcBudget, target)

							result.Set("ActualPercent", findActual)
							result.Set("NAActualPercent", false)
							if math.IsNaN(findActual) {
								result.Set("ActualPercent", 0)
								result.Set("NAActualPercent", true)
							}

							projection := tk.Div(findActual, float64(oo.Period.Month())) * 12 /*(findActual / float64(oo.Period.Month())) * 12*/
							result.Set("Projection", projection)
							result.Set("NAProjection", false)
							if math.IsNaN(findActual) {
								result.Set("Projection", 0)
								result.Set("NAProjection", true)
							}

							// budgetProjection := tk.Div(findBudget, float64(oo.Period.Month())) * 12 /*(findActual / float64(oo.Period.Month())) * 12*/
							result.Set("BudgetProjection", findBudget)
							result.Set("NABudgetProjection", false)
							if math.IsNaN(findBudget) {
								result.Set("BudgetProjection", 0)
								result.Set("NABudgetProjection", true)
							}
							higherIsBetter = true
						}
					}
				}

				result.Set("Name", CountryName)
				result.Set("RAG", oo.RAG)
				NATarget, NAActual := false, false
				if oo.Target == 130895111188 {
					NATarget = true
					result.Set("NATarget", true)
					result.Set("Target", 0)
				} else {
					result.Set("NATarget", false)
					result.Set("Target", oo.Target)
				}
				if oo.ActualYTD == 130895111188 {
					NAActual = true
					result.Set("NAActual", true)
					result.Set("Actual", 0)
					result.Set("ActualPercentage", 0)
				} else {
					result.Set("NAActual", false)
					result.Set("Actual", oo.ActualYTD)
					if higherIsBetter {
						if oo.Target < 0 && oo.ActualYTD > 0 {
							result.Set("ActualPercentage", tk.Div((oo.Target-oo.ActualYTD), oo.Target))
						} else {
							result.Set("ActualPercentage", tk.Div((oo.ActualYTD), (oo.Target)))
						}
					} else {
						if oo.Target > 0 && oo.ActualYTD < 0 {
							result.Set("ActualPercentage", tk.Div((oo.ActualYTD-oo.Target), oo.ActualYTD))
						} else {
							result.Set("ActualPercentage", tk.Div(oo.Target, oo.ActualYTD))
						}
					}
				}
				OverTarget := 0.0
				if NAActual || NATarget {
					result.Set("NAGap", true)
					result.Set("Gap", 0)
					result.Set("OverTarget", OverTarget)
				} else {
					result.Set("NAGap", false)
					result.Set("Gap", (1-result.GetFloat64("ActualPercentage"))*(-1))
					OverTarget = result.GetFloat64("ActualPercentage") - 1
					if OverTarget < 0 {
						OverTarget = 0
					}
					result.Set("OverTarget", OverTarget)
				}
				resultArr = append(resultArr, result)
			}
		}
	}

	for _, o := range CountryList {
		result := tk.M{}
		// RAG := ""

		CountryName := o.Get("Country").(string)
		RegionName := o.Get("Major_Region").(string)
		if CountryName != RegionName {
			// tk.Println("LastPeriod", LastPeriod)
			csr, err = c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Where(db.And(db.Eq("country", CountryName), db.Eq("period", LastYearPeriod), db.Eq("bmid", BusinessMetrics))).Order("-period").Cursor(nil)
			if err != nil {
				return err.Error()
			}
			BusinessMetricsDataArr := []BusinessMetricsDataModel{}
			err = csr.Fetch(&BusinessMetricsDataArr, 1, false)
			if err != nil {
				return err.Error()
			}
			csr.Close()

			if len(BusinessMetricsDataArr) > 0 {
				oo := BusinessMetricsDataArr[0]
				higherIsBetter := false
				for _, oi := range BDL1[0].BusinessMetric {
					if oi.Id == BusinessMetrics {
						target, actualytd, remainingBudget, remainingBudgetOpposite := 0.0, 0.0, 0.0, 0.0
						// if oo.Baseline != 130895111188 {
						// 	baseline = oo.Baseline
						// }
						if oo.Target != 130895111188 {
							target = oo.Target
						}
						if oo.ActualYTD != 130895111188 {
							actualytd = oo.ActualYTD
						}

						if oo.RemainingBudget != 130895111188 {
							remainingBudget = oo.RemainingBudget
						}
						if oo.RemainingBudgetOpposite != 130895111188 {
							remainingBudgetOpposite = oo.RemainingBudgetOpposite
						}

						if oi.ValueType == 0 {
							findActual := tk.Div(target, actualytd)
							if target > 0 && actualytd < 0 {
								findActual = tk.Div((actualytd - target), actualytd)
							}
							calcBudget := actualytd - remainingBudgetOpposite
							findBudget := tk.Div(calcBudget, target)

							result.Set("ActualPercent", findActual)
							result.Set("NAActualPercent", false)
							if math.IsNaN(findActual) {
								result.Set("ActualPercent", 0)
								result.Set("NAActualPercent", true)
							}

							projection := tk.Div(findActual, float64(oo.Period.Month())) * 12
							result.Set("Projection", projection)
							result.Set("NAProjection", false)
							if math.IsNaN(findActual) {
								result.Set("Projection", 0)
								result.Set("NAProjection", true)
							}
							// budgetProjection := tk.Div(findBudget, float64(oo.Period.Month())) * 12
							result.Set("BudgetProjection", findBudget)
							result.Set("NABudgetProjection", false)
							if math.IsNaN(findBudget) {
								result.Set("BudgetProjection", 0)
								result.Set("NABudgetProjection", true)
							}

						} else {
							findActual := tk.Div(actualytd, target)
							if target < 0 && actualytd > 0 {
								findActual = tk.Div((target - actualytd), target)
							}
							calcBudget := actualytd + remainingBudget
							findBudget := tk.Div(calcBudget, target)

							result.Set("ActualPercent", findActual)
							result.Set("NAActualPercent", false)
							if math.IsNaN(findActual) {
								result.Set("ActualPercent", 0)
								result.Set("NAActualPercent", true)
							}

							projection := tk.Div(findActual, float64(oo.Period.Month())) * 12 /*(findActual / float64(oo.Period.Month())) * 12*/
							result.Set("Projection", projection)
							result.Set("NAProjection", false)
							if math.IsNaN(findActual) {
								result.Set("Projection", 0)
								result.Set("NAProjection", true)
							}

							// budgetProjection := tk.Div(findBudget, float64(oo.Period.Month())) * 12 /*(findActual / float64(oo.Period.Month())) * 12*/
							result.Set("BudgetProjection", findBudget)
							result.Set("NABudgetProjection", false)
							if math.IsNaN(findBudget) {
								result.Set("BudgetProjection", 0)
								result.Set("NABudgetProjection", true)
							}
							higherIsBetter = true
						}
					}
				}

				result.Set("Name", CountryName)
				result.Set("RAG", oo.RAG)
				NATarget, NAActual := false, false
				if oo.Target == 130895111188 {
					NATarget = true
					result.Set("NATarget", true)
					result.Set("Target", 0)
				} else {
					result.Set("NATarget", false)
					result.Set("Target", oo.Target)
				}
				if oo.ActualYTD == 130895111188 {
					NAActual = true
					result.Set("NAActual", true)
					result.Set("Actual", 0)
					result.Set("ActualPercentage", 0)
				} else {
					result.Set("NAActual", false)
					result.Set("Actual", oo.ActualYTD)
					if higherIsBetter {
						if oo.Target < 0 && oo.ActualYTD > 0 {
							result.Set("ActualPercentage", tk.Div((oo.Target-oo.ActualYTD), oo.Target))
						} else {
							result.Set("ActualPercentage", tk.Div((oo.ActualYTD), (oo.Target)))
						}
					} else {
						if oo.Target > 0 && oo.ActualYTD < 0 {
							result.Set("ActualPercentage", tk.Div((oo.ActualYTD-oo.Target), oo.ActualYTD))
						} else {
							result.Set("ActualPercentage", tk.Div(oo.Target, oo.ActualYTD))
						}
					}
				}
				OverTarget := 0.0
				if NAActual || NATarget {
					result.Set("NAGap", true)
					result.Set("Gap", 0)
					result.Set("OverTarget", OverTarget)
				} else {
					result.Set("NAGap", false)
					result.Set("Gap", (1-result.GetFloat64("ActualPercentage"))*(-1))
					OverTarget = result.GetFloat64("ActualPercentage") - 1
					if OverTarget < 0 {
						OverTarget = 0
					}
					result.Set("OverTarget", OverTarget)
				}
				resultLastYear = append(resultLastYear, result)
			}
		}
	}

	for _, x := range resultArr {
		action := true
		for _, y := range resultLastYear {
			// tk.Println("ACTUALPRG", y.GetFloat64("ActualPercent"))
			if x.GetString("Name") == y.GetString("Name") {
				x.Set("LastYearActual", y.GetFloat64("Actual"))
				x.Set("LastYearNAActual", y.Get("NAActual"))
				x.Set("LastYearActualPercent", y.GetFloat64("ActualPercent"))
				x.Set("LastYearActualPercentage", y.GetFloat64("ActualPercentage"))
				x.Set("LastYearBudgetProjection", y.GetFloat64("BudgetProjection"))
				x.Set("LastYearGap", y.GetFloat64("Gap"))
				x.Set("LastYearNAActualPercent", y.Get("NAActualPercent"))
				x.Set("LastYearNABudgetProjection", y.Get("NABudgetProjection"))
				x.Set("LastYearNAGap", y.Get("NAGap"))
				x.Set("LastYearNAProjection", y.Get("NAProjection"))
				x.Set("LastYearNATarget", y.Get("NATarget"))
				x.Set("LastYearName", y.GetString("Name"))
				x.Set("LastYearOverTarget", y.GetFloat64("OverTarget"))
				x.Set("LastYearProjection", y.GetFloat64("Projection"))
				x.Set("LastYearTarget", y.GetFloat64("Target"))
				x.Set("LastYearRAG", y.GetString("RAG"))
				action = false
				break
			}
		}

		if len(resultLastYear) == 0 || action {
			x.Set("LastYearActual", 0)
			x.Set("LastYearNAActual", true)
			x.Set("LastYearActualPercent", 0)
			x.Set("LastYearActualPercentage", 0)
			x.Set("LastYearBudgetProjection", 0)
			x.Set("LastYearGap", 0)
			x.Set("LastYearNAActualPercent", true)
			x.Set("LastYearNABudgetProjection", true)
			x.Set("LastYearNAGap", true)
			x.Set("LastYearNAProjection", true)
			x.Set("LastYearNATarget", true)
			x.Set("LastYearName", x.GetString("Name"))
			x.Set("LastYearOverTarget", 0)
			x.Set("LastYearProjection", 0)
			x.Set("LastYearTarget", 0)
			x.Set("LastYearRAG", "")
		}

	}
	return resultArr
}
func (c *ScorecardAnalysisController) GetCountryAnalysisTrendlineData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		BusinessMetrics string
		Region          string
		Country         string
		Period          string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	csr, err := c.Ctx.Find(new(BusinessDriverL1Model), tk.M{}.Set("where", db.Eq("businessmetric.id", parm.BusinessMetrics)))
	if err != nil {
		return err.Error()
	}
	BDL1 := make([]BusinessDriverL1Model, 0)
	err = csr.Fetch(&BDL1, 0, false)
	if err != nil {
		return err.Error()
	}
	csr.Close()

	csr, err = c.Ctx.Connection.NewQuery().From("Region").Where(db.Eq("Country", parm.Country)).Cursor(nil)
	if err != nil {
		return err.Error()
	}
	CountryList := []tk.M{}
	err = csr.Fetch(&CountryList, 0, false)
	if err != nil {
		return err.Error()
	}
	csr.Close()

	Period, _ := time.Parse("20060102", parm.Period)
	eachmonth := [12][]tk.M{}

	for _, o := range CountryList {

		CountryName := o.Get("Country").(string)
		RegionName := o.Get("Major_Region").(string)
		if CountryName != RegionName {

			csr, err = c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Where(db.And(db.Eq("country", CountryName), db.Eq("year", Period.Year()), db.Eq("bmid", parm.BusinessMetrics))).Order("-period").Cursor(nil)
			if err != nil {
				return err.Error()
			}
			BusinessMetricsDataArr := []BusinessMetricsDataModel{}
			err = csr.Fetch(&BusinessMetricsDataArr, 0, false)
			if err != nil {
				return err.Error()
			}
			csr.Close()

			if len(BusinessMetricsDataArr) > 0 {
				// tk.Printf("bussiness =====%d\n", len(BusinessMetricsDataArr))
				for _, oo := range BusinessMetricsDataArr {
					result := tk.M{}
					higherIsBetter := false
					MetricType := ""
					for _, oi := range BDL1[0].BusinessMetric {
						if oi.Id == parm.BusinessMetrics {
							MetricType = oi.MetricType
							target, actualytd, budget := 0.0, 0.0, 0.0
							// target, actualytd, budget, remainingBudget, remainingBudgetOpposite := 0.0, 0.0, 0.0, 0.0, 0.0
							if oo.Target != 130895111188 {
								target = oo.Target
							}
							if oo.ActualYTD != 130895111188 {
								actualytd = oo.ActualYTD
							}
							if oo.Budget != 130895111188 {
								budget = oo.Budget
							}
							// if oo.RemainingBudget != 130895111188 {
							// 	remainingBudget = oo.RemainingBudget
							// }
							// if oo.RemainingBudgetOpposite != 130895111188 {
							// 	remainingBudgetOpposite = oo.RemainingBudgetOpposite
							// }
							result.Set("Region", RegionName)
							result.Set("Actual", actualytd)
							if oi.ValueType == 0 {

								findActual := tk.Div(target, actualytd)
								if target > 0 && actualytd < 0 {
									findActual = tk.Div((actualytd - target), actualytd)
								}
								// calcBudget := actualytd - remainingBudgetOpposite
								// findBudget := tk.Div(calcBudget, target)

								result.Set("ActualMoM", findActual)

								if math.IsNaN(findActual) {
									result.Set("ActualMoM", 0)

								}
								if oo.Budget != 130895111188 {
									result.Set("Budget", budget)
								}

								// if math.IsNaN(findBudget) {
								// 	result.Set("Budget", 0)

								// }

							} else {
								findActual := tk.Div(actualytd, target)
								if target < 0 && actualytd > 0 {
									findActual = tk.Div((target - actualytd), target)
								}
								// calcBudget := actualytd + remainingBudget
								// findBudget := tk.Div(calcBudget, target)

								result.Set("ActualMoM", findActual)

								if math.IsNaN(findActual) {
									result.Set("ActualMoM", 0)

								}
								if oo.Budget != 130895111188 {
									result.Set("Budget", budget)
								}
								// if math.IsNaN(findBudget) {
								// 	result.Set("Budget", 0)

								// }
								higherIsBetter = true
							}
						}
					}

					result.Set("Name", CountryName)
					result.Set("MetricType", MetricType)
					result.Set("Month", oo.Period.Month().String()[:3])

					if oo.Target == 130895111188 {

					} else {

						result.Set("Target", oo.Target)
					}
					if oo.ActualYTD == 130895111188 {

					} else {

						result.Set("ActualCummulative", oo.ActualYTD)
						// if higherIsBetter {
						// 	result.Set("ActualMoM", tk.Div((oo.ActualYTD), (oo.Target)))
						// } else {
						// 	result.Set("ActualMoM", tk.Div(oo.Target, oo.ActualYTD))
						// }
						if higherIsBetter {
							if oo.Target < 0 && oo.ActualYTD > 0 {
								result.Set("ActualMoM", tk.Div((oo.Target-oo.ActualYTD), oo.Target))
							} else {
								result.Set("ActualMoM", tk.Div((oo.ActualYTD), (oo.Target)))
							}
						} else {
							if oo.Target > 0 && oo.ActualYTD < 0 {
								result.Set("ActualMoM", tk.Div((oo.ActualYTD-oo.Target), oo.ActualYTD))
							} else {
								result.Set("ActualMoM", tk.Div(oo.Target, oo.ActualYTD))
							}
						}
					}

					eachmonth[int(oo.Period.Month())-1] = append(eachmonth[int(oo.Period.Month())-1], result)

				}
			}
		}
	}
	ActualYTD := 0.0
	for _, s := range eachmonth {
		for _, x := range s {
			YtdValue := x.GetFloat64("ActualCummulative")
			ActualMoM := 0.0
			MetricType := x.GetString("MetricType")
			if MetricType == "spot" {
				ActualMoM = YtdValue
			} else {
				ActualMoM = YtdValue - ActualYTD
			}
			if x.Get("ActualCummulative") != nil {
				x.Set("ActualMoM", ActualMoM)
			} else {
				x.Unset("ActualMoM")
			}
			ActualYTD = YtdValue
		}
	}
	for _, s := range eachmonth {
		sort.Sort(SortGAP(s))
	}

	return c.SetResultInfo(false, "", eachmonth)

}

func (c *ScorecardAnalysisController) GetFullYearRankingData(k *knot.WebContext, parm ParmScorecardAnalysis) interface{} {
	// tk.Println(parm)
	csr, err := c.Ctx.Find(new(BusinessDriverL1Model), tk.M{}.Set("where", db.Eq("businessmetric.id", parm.BusinessMetrics)))
	if err != nil {
		return err.Error()
	}
	BDL1 := make([]BusinessDriverL1Model, 0)
	err = csr.Fetch(&BDL1, 0, false)
	if err != nil {
		return err.Error()
	}
	csr.Close()

	csr, err = c.Ctx.Connection.NewQuery().From("Region").
		Cursor(nil)
	if err != nil {
		return err.Error()
	}
	CountryList := []tk.M{}
	err = csr.Fetch(&CountryList, 0, false)
	if err != nil {
		return err.Error()
	}
	csr.Close()

	Period, _ := time.Parse("20060102", parm.Period)
	// Period = Period.UTC()
	// resultArr := []tk.M{}
	metrictype := [12]string{}
	valuetype := [12]bool{}
	eachmonth := [12][]tk.M{}

	for _, o := range CountryList {

		CountryName := o.Get("Country").(string)
		RegionName := o.Get("Major_Region").(string)
		if CountryName != RegionName {

			csr, err = c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Where(db.And(db.Eq("country", CountryName), db.Eq("year", Period.Year()), db.Eq("bmid", parm.BusinessMetrics))).Order("-period").Cursor(nil)
			if err != nil {
				return err.Error()
			}
			BusinessMetricsDataArr := []BusinessMetricsDataModel{}
			err = csr.Fetch(&BusinessMetricsDataArr, 0, false)
			if err != nil {
				return err.Error()
			}
			csr.Close()

			if len(BusinessMetricsDataArr) > 0 {
				// oo := BusinessMetricsDataArr[0]
				// tk.Printf("bussiness =====%d\n", len(BusinessMetricsDataArr))
				for _, oo := range BusinessMetricsDataArr {
					result := tk.M{}
					higherIsBetter := false
					Type := ""
					for _, oi := range BDL1[0].BusinessMetric {
						if oi.Id == parm.BusinessMetrics {
							Type = oi.Type
							target, actualytd, remainingBudget, remainingBudgetOpposite := 0.0, 0.0, 0.0, 0.0
							// if oo.Baseline != 130895111188 {
							// 	baseline = oo.Baseline
							// }
							if oo.Target != 130895111188 {
								target = oo.Target
							}
							if oo.ActualYTD != 130895111188 {
								actualytd = oo.ActualYTD
							}

							if oo.RemainingBudget != 130895111188 {
								remainingBudget = oo.RemainingBudget
							}
							if oo.RemainingBudgetOpposite != 130895111188 {
								remainingBudgetOpposite = oo.RemainingBudgetOpposite
							}

							if oi.ValueType == 0 {

								findActual := tk.Div(target, actualytd)
								if target > 0 && actualytd < 0 {
									findActual = tk.Div((actualytd - target), actualytd)
								}
								calcBudget := actualytd - remainingBudgetOpposite
								findBudget := tk.Div(calcBudget, target)

								result.Set("ActualPercent", findActual)
								result.Set("NAActualPercent", false)
								if math.IsNaN(findActual) {
									result.Set("ActualPercent", 0)
									result.Set("NAActualPercent", true)
								}

								projection := tk.Div(findActual, float64(oo.Period.Month())) * 12
								result.Set("Projection", projection)
								result.Set("NAProjection", false)
								if math.IsNaN(findActual) {
									result.Set("Projection", 0)
									result.Set("NAProjection", true)
								}
								// budgetProjection := tk.Div(findBudget, float64(oo.Period.Month())) * 12
								result.Set("BudgetProjection", findBudget)
								result.Set("NABudgetProjection", false)
								if math.IsNaN(findBudget) {
									result.Set("BudgetProjection", 0)
									result.Set("NABudgetProjection", true)
								}

							} else {
								findActual := tk.Div(actualytd, target)
								if target < 0 && actualytd > 0 {
									findActual = tk.Div((target - actualytd), target)
								}
								calcBudget := actualytd + remainingBudget
								findBudget := tk.Div(calcBudget, target)

								result.Set("ActualPercent", findActual)
								result.Set("NAActualPercent", false)
								if math.IsNaN(findActual) {
									result.Set("ActualPercent", 0)
									result.Set("NAActualPercent", true)
								}

								projection := tk.Div(findActual, float64(oo.Period.Month())) * 12 /*(findActual / float64(oo.Period.Month())) * 12*/
								result.Set("Projection", projection)
								result.Set("NAProjection", false)
								if math.IsNaN(findActual) {
									result.Set("Projection", 0)
									result.Set("NAProjection", true)
								}

								// budgetProjection := tk.Div(findBudget, float64(oo.Period.Month())) * 12 /*(findActual / float64(oo.Period.Month())) * 12*/
								result.Set("BudgetProjection", findBudget)
								result.Set("NABudgetProjection", false)
								if math.IsNaN(findBudget) {
									result.Set("BudgetProjection", 0)
									result.Set("NABudgetProjection", true)
								}
								higherIsBetter = true
							}
						}
					}
					result.Set("RegionName", RegionName)
					result.Set("Name", CountryName)
					result.Set("CountryCode", oo.CountryCode)
					result.Set("RAG", oo.RAG)
					NATarget, NAActual := false, false
					if oo.Target == 130895111188 {
						NATarget = true
						result.Set("NATarget", true)
						result.Set("Target", 0)
					} else {
						result.Set("NATarget", false)
						result.Set("Target", oo.Target)
					}
					if oo.ActualYTD == 130895111188 {
						NAActual = true
						result.Set("NAActual", true)
						result.Set("Actual", 0)
						result.Set("ActualPercentage", 0)
					} else {
						result.Set("NAActual", false)
						result.Set("Actual", oo.ActualYTD)
						if higherIsBetter {
							if oo.Target < 0 && oo.ActualYTD > 0 {
								result.Set("ActualPercentage", tk.Div((oo.Target-oo.ActualYTD), oo.Target))
							} else {
								result.Set("ActualPercentage", tk.Div((oo.ActualYTD), (oo.Target)))
							}
						} else {
							if oo.Target > 0 && oo.ActualYTD < 0 {
								result.Set("ActualPercentage", tk.Div((oo.ActualYTD-oo.Target), oo.ActualYTD))
							} else {
								result.Set("ActualPercentage", tk.Div(oo.Target, oo.ActualYTD))
							}
						}
					}
					if NAActual || NATarget {
						result.Set("NAGap", true)
						result.Set("Gap", 0)
					} else {
						result.Set("NAGap", false)
						result.Set("Gap", (1-result.GetFloat64("ActualPercentage"))*(-1))
					}

					// tk.Printf("BULAN %v %d\n", oo.Period.Month(), int(oo.Period.Month()))
					// tk.Println("INI RESULT LHO", result.GetString("Gap"))
					metrictype[int(oo.Period.Month())-1] = Type
					valuetype[int(oo.Period.Month())-1] = higherIsBetter
					eachmonth[int(oo.Period.Month())-1] = append(eachmonth[int(oo.Period.Month())-1], result)
					// resultArr = append(resultArr, result)
				}
			}
		}
	}

	for x, s := range eachmonth {
		mtype := metrictype[x]
		vtype := valuetype[x]
		if mtype == "spot" {
			if vtype {
				sort.Sort(SortActualDescending(s))
			} else {
				sort.Sort(SortActual(s))

			}
		} else {
			sort.Sort(SortGAP(s))
		}
	}

	// return c.SetResultInfo(false, "", eachmonth)
	return eachmonth
}

func (c *ScorecardAnalysisController) GetMetricsRankingData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	BDL1 := make([]BusinessDriverL1Model, 0)

	csr, err := c.Ctx.Connection.NewQuery().From("BusinessDriverL1").Order("seq").Cursor(nil)
	defer csr.Close()
	err = csr.Fetch(&BDL1, 0, true)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	query := []*db.Filter{}
	result := []tk.M{}
	for bd := range BDL1 {
		for bm := range BDL1[bd].BusinessMetric {
			if BDL1[bd].BusinessMetric[bm].DecimalFormat == "" {
				BDL1[bd].BusinessMetric[bm].DecimalFormat = "0"
			}
			BDL1[bd].BusinessMetric[bm].Display = ""
			bmdataList := []BusinessMetricsDataModel{}
			data := BusinessMetricsDataModel{}
			query = append(query[0:0], db.Eq("bmid", BDL1[bd].BusinessMetric[bm].Id))
			query = append(query, db.Eq("majorregion", "GLOBAL"))
			query = append(query, db.Eq("country", "GLOBAL"))
			csr, err := c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(db.And(query...)).Order("period").Cursor(nil)
			err = csr.Fetch(&bmdataList, 0, false)
			csr.Close()
			if err != nil {
				continue
			}
			if len(bmdataList) > 0 {
				isAny := false
				for idx, i := range bmdataList {
					if i.ActualYTD != 130895111188 && !i.NAActual {
						data = bmdataList[idx]
						isAny = true
					}
				}
				if !isAny {
					data = bmdataList[0]
				}
			}

			BMType := BDL1[bd].BusinessMetric[bm].Type
			MetricsPeriod := data.Period

			DataList := []tk.M{}

			bmdataList = []BusinessMetricsDataModel{}
			// Get Data base on latest Period [Metrics Period]
			query = append(query[0:0], db.Eq("period", MetricsPeriod))
			query = append(query, db.Eq("bmid", data.BMId))
			query = append(query, db.Nin("country", "GLOBAL"))
			query = append(query, db.Ne("country", "AME"))
			query = append(query, db.Ne("country", "ASA"))
			query = append(query, db.Ne("country", "GCNA"))
			csr, err = c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(db.And(query...)).Cursor(nil)
			err = csr.Fetch(&bmdataList, 0, false)
			if err != nil {
				continue
			}
			csr.Close()
			higherIsBetter := false
			if BDL1[bd].BusinessMetric[bm].ValueType != 0 {
				higherIsBetter = true
			}
			for _, x := range bmdataList {
				tmp := tk.M{}
				tmp.Set("Region", x.MajorRegion)
				tmp.Set("CountryCode", x.CountryCode)
				tmp.Set("CountryName", x.Country)
				Actual, PercentGap := 0.0, 0.0

				NATarget, NAActual := false, false
				div := 0.0
				if x.Target == 130895111188 {
					NATarget = true
				} else {
					NATarget = false
				}
				if x.ActualYTD == 130895111188 {
					NAActual = true
					Actual = 0
				} else {
					Actual = x.ActualYTD
					if higherIsBetter {
						div = tk.Div((x.ActualYTD), (x.Target))
					} else {
						div = tk.Div(x.Target, x.ActualYTD)
					}
				}
				if NAActual || NATarget {
					PercentGap = 0
				} else {
					PercentGap = (1 - div) * (-1)
				}

				tmp.Set("Actual", Actual)
				tmp.Set("PercentGap", PercentGap)
				DataList = append(DataList, tmp)
			}

			//LastYearDataList
			LastYearDataList := []tk.M{}
			LastYear := MetricsPeriod.Year() - 1
			LastYearPeriod := time.Date(LastYear, MetricsPeriod.Month(), 1, 0, 0, 0, 0, time.UTC)

			bmdataList = []BusinessMetricsDataModel{}
			// Get Data base on latest Period [Metrics Period]
			query = append(query[0:0], db.Eq("period", LastYearPeriod))
			query = append(query, db.Eq("bmid", data.BMId))
			query = append(query, db.Nin("country", "GLOBAL"))
			query = append(query, db.Ne("country", "AME"))
			query = append(query, db.Ne("country", "ASA"))
			query = append(query, db.Ne("country", "GCNA"))
			csr, err = c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(db.And(query...)).Cursor(nil)
			err = csr.Fetch(&bmdataList, 0, false)
			if err != nil {
				continue
			}
			csr.Close()
			higherIsBetter = false
			if BDL1[bd].BusinessMetric[bm].ValueType != 0 {
				higherIsBetter = true
			}
			for _, x := range bmdataList {
				tmp := tk.M{}
				tmp.Set("Region", x.MajorRegion)
				tmp.Set("CountryCode", x.CountryCode)
				tmp.Set("CountryName", x.Country)
				Actual, PercentGap := 0.0, 0.0

				NATarget, NAActual := false, false
				div := 0.0
				if x.Target == 130895111188 {
					NATarget = true
				} else {
					NATarget = false
				}
				if x.ActualYTD == 130895111188 {
					NAActual = true
					Actual = 0
				} else {
					Actual = x.ActualYTD
					if higherIsBetter {
						div = tk.Div((x.ActualYTD), (x.Target))
					} else {
						div = tk.Div(x.Target, x.ActualYTD)
					}
				}
				if NAActual || NATarget {
					PercentGap = 0
				} else {
					PercentGap = (1 - div) * (-1)
				}

				tmp.Set("Actual", Actual)
				tmp.Set("PercentGap", PercentGap)
				LastYearDataList = append(LastYearDataList, tmp)
			}
			//YoY % = (2017 YTD Actual - 2016 YTD Actual)/2016 YTD Actual
			for _, x := range DataList {
				for _, j := range LastYearDataList {
					if x.GetString("CountryName") == j.GetString("CountryName") {
						x.Set("LastYearActual", j.Get("Actual"))
						x.Set("LastYearPercentGap", j.Get("PercentGap"))
						actual2016 := j.Get("Actual").(float64)
						actual2017 := x.Get("Actual").(float64)
						LastYearYoY := tk.Div(actual2017-actual2016, actual2016)
						x.Set("YoY", LastYearYoY)
					}
				}
			}

			each := tk.M{}
			each.Set("Id", BDL1[bd].BusinessMetric[bm].Id)
			each.Set("Name", BDL1[bd].BusinessMetric[bm].Description)
			each.Set("Period", MetricsPeriod.Format("20060102"))
			each.Set("Type", BMType)
			each.Set("IsHigherIsBetter", higherIsBetter)
			each.Set("DataList", DataList)
			result = append(result, each)
		}
	}

	return c.SetResultInfo(false, "", result)
}
