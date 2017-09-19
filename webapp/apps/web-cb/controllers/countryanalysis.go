package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// "strings"
	"math"
	// "strconv"
	"time"
)

type CountryAnalysisController struct {
	*BaseController
}

func (c *CountryAnalysisController) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputTemplate
	// c.Action(k, "Open Country Analysis Page")
	c.Action(k, "Country Analysis", "Open Country Analysis Page", "", "", "", "", "")

	CountryAnalysis := c.GetAccess(k, "COUNTRYANALYSIS")
	Initiative := c.GetAccess(k, "INITIATIVE")
	Scorecard := c.GetAccess(k, "SCORECARD")
	k.Config.LayoutTemplate = "_layout-v2.html"
	PartialFiles := []string{}
	PartialFiles = append(PartialFiles, "shared/sidebar.html")
	PartialFiles = append(PartialFiles, "shared/initiative_summary.html")
	PartialFiles = append(PartialFiles, "shared/initiative_chart.html")
	PartialFiles = append(PartialFiles, "shared/metric_upload.html")
	PartialFiles = append(PartialFiles, "dashboard/initiative.html")

	k.Config.IncludeFiles = PartialFiles
	return tk.M{}.Set("CountryAnalysis", CountryAnalysis).Set("Initiative", Initiative).Set("Scorecard", Scorecard)
}

func (c *CountryAnalysisController) GetData(k *knot.WebContext) interface{} {
	// c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		BMId           string
		Breakdown      string
		RelevantFilter string
		Period         string
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	csr, err := c.Ctx.Find(new(BusinessDriverL1Model), tk.M{}.Set("where", db.Eq("businessmetric.id", parm.BMId)))
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	BDL1 := make([]BusinessDriverL1Model, 0)
	err = csr.Fetch(&BDL1, 0, false)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()

	// UserCountry := ""
	// if k.Session("countrycode") != nil {
	// 	UserCountry = k.Session("countrycode").(string)
	// 	csr, err = c.Ctx.Connection.NewQuery().From("Region").Where(db.Eq("CountryCode", UserCountry)).
	// 		Cursor(nil)
	// } else {
	csr, err = c.Ctx.Connection.NewQuery().From("Region").
		Cursor(nil)
	// }
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	CountryList := []tk.M{}
	err = csr.Fetch(&CountryList, 0, false)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()

	csr, err = c.Ctx.Connection.NewQuery().From("Region").Group("Major_Region").
		Cursor(nil)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	RegionList := []tk.M{}
	err = csr.Fetch(&RegionList, 0, false)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()

	resultArr := []tk.M{}

	// fmt.Println(parm.Year.Year())
	// startDate := time.Now()
	// endDate := time.Now()
	// layout := "01-01-2000"
	// if parm.Year.Year() == 3000 {
	// 	//berarti ambil 2017 saja
	// 	startDate, _ = time.Parse(layout, "01-01-2017")
	// 	endDate, _ = time.Parse(layout, "31-12-2017")
	// } else if parm.Year.Year() == 1000 {
	// 	//ambil dari awal 2010 sampai akhir year now
	// 	startDate, _ = time.Parse(layout, "01-01-2010")
	// 	endDate, _ = time.Parse(layout, "3=1-12-"+strconv.Itoa(time.Now().Year()))
	// } else {
	// 	startDate = parm.Year
	// 	endDate = parm.Year.AddDate(0, 1, 0)
	// }

	// Period := tk.String2Date(parm.Period, "MMMMM-YYYY")
	Period, _ := time.Parse("January-2006", parm.Period)
	Period = Period.UTC()
	// fmt.Println(tk.JsonString(BDL1[0].BusinessMetric))

	if parm.Breakdown == "region" {
		for _, o := range RegionList {
			result := tk.M{}
			// RAG := ""

			RegionName := o.Get("_id").(tk.M).Get("Major_Region").(string)
			// csr, err = c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Where(db.And(db.Eq("majorregion", RegionName), db.Eq("year", parm.Year), db.Eq("bmid", parm.BMId))).Order("-period").Cursor(nil)
			// csr, err = c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Where(db.And(db.Eq("majorregion", RegionName), db.Or(db.Gte("period", startDate), db.Lt("period", endDate)), db.Eq("bmid", parm.BMId))).Order("-period").Cursor(nil)
			csr, err = c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Where(db.And(db.Eq("countrycode", RegionName), db.Eq("period", Period), db.Eq("bmid", parm.BMId))).Order("-period").Cursor(nil)

			if err != nil {
				return c.ErrorResultInfo(err.Error(), nil)
			}
			BusinessMetricsDataArr := []BusinessMetricsDataModel{}
			err = csr.Fetch(&BusinessMetricsDataArr, 1, false)
			if err != nil {
				return c.ErrorResultInfo(err.Error(), nil)
			}
			csr.Close()

			if len(BusinessMetricsDataArr) > 0 {
				oo := BusinessMetricsDataArr[0]

				higherIsBetter := false
				for _, oi := range BDL1[0].BusinessMetric {
					if oi.Id == parm.BMId {
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
							higherIsBetter = true
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

							// budgetProjection := tk.Div(findBudget, float64(oo.Period.Month())) * 12
							result.Set("BudgetProjection", findBudget)
							result.Set("NABudgetProjection", false)
							if math.IsNaN(findBudget) {
								result.Set("BudgetProjection", 0)
								result.Set("NABudgetProjection", true)
							}
						}
					}
				}

				result.Set("Name", RegionName)
				// result.Set("RAG", RAG)
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
				resultArr = append(resultArr, result)
			}
		}
	} else if parm.Breakdown == "country" {
		for _, o := range CountryList {
			result := tk.M{}
			// RAG := ""

			CountryName := o.Get("Country").(string)
			RegionName := o.Get("Major_Region").(string)
			if CountryName != RegionName {
				csr, err = c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Where(db.And(db.Eq("country", CountryName), db.Eq("period", Period), db.Eq("bmid", parm.BMId))).Order("-period").Cursor(nil)
				if err != nil {
					return c.ErrorResultInfo(err.Error(), nil)
				}
				BusinessMetricsDataArr := []BusinessMetricsDataModel{}
				err = csr.Fetch(&BusinessMetricsDataArr, 1, false)
				if err != nil {
					return c.ErrorResultInfo(err.Error(), nil)
				}
				csr.Close()

				if len(BusinessMetricsDataArr) > 0 {
					oo := BusinessMetricsDataArr[0]
					higherIsBetter := false
					for _, oi := range BDL1[0].BusinessMetric {
						if oi.Id == parm.BMId {
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
					if NAActual || NATarget {
						result.Set("NAGap", true)
						result.Set("Gap", 0)
					} else {
						result.Set("NAGap", false)
						result.Set("Gap", (1-result.GetFloat64("ActualPercentage"))*(-1))
					}
					resultArr = append(resultArr, result)
				}
			}
		}
	} else {
		result := tk.M{}
		// RAG := ""

		RegionName := "GLOBAL"
		csr, err = c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Where(db.And(db.Eq("countrycode", RegionName), db.Eq("period", Period), db.Eq("bmid", parm.BMId))).Order("-period").Cursor(nil)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
		BusinessMetricsDataArr := []BusinessMetricsDataModel{}
		err = csr.Fetch(&BusinessMetricsDataArr, 1, false)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
		csr.Close()

		if len(BusinessMetricsDataArr) > 0 {
			oo := BusinessMetricsDataArr[0]
			higherIsBetter := false
			for _, oi := range BDL1[0].BusinessMetric {
				if oi.Id == parm.BMId {
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
					// if oo.Budget != 130895111188 {
					// 	budget = oo.Budget
					// }

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
						// budgetProjection := tk.Div(findBudget, float64(oo.Period.Month())) * 12
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

			result.Set("Name", RegionName)
			// result.Set("RAG", RAG)
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
			resultArr = append(resultArr, result)
		}
	}

	// for _, o := range RegionList {

	// Country := o.Get("Country").(string)
	// MajorRegion := o.Get("Major_Region").(string)

	// for _, oi := range BDL1 {
	// 	for _, ooii := range oi.BusinessMetric {
	// 		// csr, err := c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Where(db.And(db.Eq("region", MajorRegion), db.Eq("country", Country), db.Eq("bmid", ooii.Id))).
	// 		csr, err := c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Where(db.And(db.Eq("year", parm.Year), db.Eq("bmid", ooii.Id))).
	// 			Cursor(nil)
	// 		if err != nil {
	// 			return c.ErrorResultInfo(err.Error(), nil)
	// 		}
	// 		BusinessMetricsDataArr := []tk.M{}
	// 		err = csr.Fetch(&BusinessMetricsDataArr, 0, false)
	// 		if err != nil {
	// 			return c.ErrorResultInfo(err.Error(), nil)
	// 		}
	// 		csr.Close()

	// 		// fmt.Println(BusinessMetricsDataArr, ooii.Id, MajorRegion, Country)
	// 		if len(BusinessMetricsDataArr) > 0 {
	// 			Name = ooii.Country
	// 			RAG = ooii.
	// 		}
	// 	}
	// }
	// }
	// for _, o := range RegionList {

	// }
	// fmt.Println(CountryList)
	// csr, err = c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Where(db.And(db.Eq("year", parm.Year), db.Eq("bmid", parm.BMId))).
	// 	Cursor(nil)
	// if err != nil {
	// 	return c.ErrorResultInfo(err.Error(), nil)
	// }
	// // BusinessMetricsDataArr := []tk.M{}
	// BusinessMetricsDataArr := []BusinessMetricsDataModel{}
	// err = csr.Fetch(&BusinessMetricsDataArr, 2, false)
	// if err != nil {
	// 	return c.ErrorResultInfo(err.Error(), nil)
	// }
	// csr.Close()

	// for i, o := range BusinessMetricsDataArr {
	// 	if i == 0 {
	// 		result := tk.M{}
	// 		RAG:= "Sad", 0.0, 0.0, 0.0

	// 		if parm.Breakdown == "country" {
	// 			Name = o.Country
	// 		} else if parm.Breakdown == "region" {
	// 			Name = o.Region
	// 		}

	// 		Target = o.Target
	// 		Actual = o.ActualYTD
	// 		Gap = (Target - Actual) / Target

	// 		for _, oi := range BDL1 {
	// 			for _, ooii := range oi.BusinessMetric {
	// 				if ooii.Id == parm.BMId {
	// 					RAG = ooii.Display
	// 				}
	// 			}
	// 		}

	// 		result.Set("Name", Name)
	// 		result.Set("RAG", RAG)
	// 		result.Set("Target", Target)
	// 		result.Set("Actual", Actual)
	// 		result.Set("Gap", Gap)
	// 		resultArr = append(resultArr, result)
	// 	}
	// }

	return c.SetResultInfo(false, "", resultArr)
}
