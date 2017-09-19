package controllers

import (
	"eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"fmt"
	"github.com/disintegration/imaging"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2/bson"
	"image"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DashboardController struct {
	*BaseController
}

func (c *DashboardController) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	AccessibilityData := new(AccessibilityModel)
	group := k.Session("group", []string{}).([]string)
	roles := []interface{}{}
	for _, x := range group {
		roles = append(roles, x)
	}
	csr, _ := c.AclCtx.Connection.NewQuery().From(AccessibilityData.TableName()).Where(db.And(db.In("roleid", roles...), db.Eq("allowstatus", true))).Order("index").Cursor(nil)
	_ = csr.Fetch(&AccessibilityData, 1, false)
	csr.Close()
	access_url := "/web-cb" + AccessibilityData.Url
	if access_url != k.Request.URL.String() {
		http.Redirect(k.Writer, k.Request, access_url, http.StatusTemporaryRedirect)
	}

	Access := c.GetAccess(k, "DASHBOARD")
	Scorecard := c.GetAccess(k, "SCORECARD")
	Initiative := c.GetAccess(k, "INITIATIVE")
	SharedAgenda := c.GetAccess(k, "SHAREDAGENDA")
	METRICCOUNTRYANALYSIS := c.GetAccess(k, "METRICCOUNTRYANALYSIS")
	// c.Action(k, "Open Scorecard & Initiatives Page")
	c.Action(k, "Dashboard", "Open Scorecard & Initiatives Page", "", "", "", "", "")

	// k.Config.NoLog = true
	// k.Config.LayoutTemplate = "_layout-v2.html"
	PartialFiles := []string{}
	PartialFiles = append(PartialFiles, "shared/sidebar.html")
	PartialFiles = append(PartialFiles, "shared/initiative_summary.html")
	PartialFiles = append(PartialFiles, "shared/initiative_chart.html")
	PartialFiles = append(PartialFiles, "shared/metric_upload.html")
	PartialFiles = append(PartialFiles, "dashboard/initiative.html")

	PartialFiles = append(PartialFiles, "dashboard/scorecard_initiative.html")
	PartialFiles = append(PartialFiles, "dashboard/scorecard_analysis.html")
	PartialFiles = append(PartialFiles, "dashboard/scorecard_countryanalysis.html")
	PartialFiles = append(PartialFiles, "dashboard/scorecard_fullyearranking.html")
	PartialFiles = append(PartialFiles, "dashboard/scorecard_fullyearprojection.html")
	PartialFiles = append(PartialFiles, "dashboard/scorecard_metriccountryranking.html")
	PartialFiles = append(PartialFiles, "dashboard/scorecard_detail.html")
	PartialFiles = append(PartialFiles, "dashboard/scorecard_detail_form.html")
	PartialFiles = append(PartialFiles, "dashboard/sharedagenda.html")
	PartialFiles = append(PartialFiles, "dashboard/search.html")
	PartialFiles = append(PartialFiles, "dashboard/chart.html")
	PartialFiles = append(PartialFiles, "dashboard/initiativeTab.html")
	PartialFiles = append(PartialFiles, "dashboard/scorecard.html")
	PartialFiles = append(PartialFiles, "dashboard/scorecard_bm.html")
	PartialFiles = append(PartialFiles, "dashboard/task.html")
	PartialFiles = append(PartialFiles, "dashboard/modalclone.html")

	k.Config.IncludeFiles = PartialFiles
	k.Config.OutputType = knot.OutputTemplate
	IsFinADMIN := false
	IsProdADMIN := false

	var groups []string

	if k.Session("group") != nil {
		groups = k.Session("group").([]string)
	}

	if k.Session("username") != nil && stringInSlice("CB_FINADMIN", groups) || stringInSlice("CB_ADMIN", groups) {
		IsFinADMIN = true
		// tk.Println(k.Session("username").(string))
	}

	if stringInSlice("CB_PRODADMIN", groups) {
		IsProdADMIN = true
		IsFinADMIN = true
	}

	// if stringInSlice("CB_METRICOWNER", groups) {
	// 	c.Redirect(k, "bef/metricupload", "default")
	// 	// c.Redirect(k, controller, action)
	// }

	// tk.Println("IsMetricOwner", IsMetricOwner)
	// tk.Println(IsProdADMIN)
	// tk.Println(k.Session("username").(string))
	// tk.Println(IsFinADMIN)
	UserCountry := ""
	if k.Session("country") != nil {
		UserCountry = k.Session("country").(string)
	}
	return tk.M{}.Set("UserCountry", UserCountry).Set("IsFINADMIN", !IsFinADMIN).Set("IsDedicated", false).Set("IsProdADMIN", IsProdADMIN).Set("SharedAgenda", SharedAgenda).Set("Access", Access).Set("Initiative", Initiative).Set("Scorecard", Scorecard).Set("METRICCOUNTRYANALYSIS", METRICCOUNTRYANALYSIS)
}

type NumberMetric struct {
	Country     string
	MajorRegion string
	AMBER       int
	GREEN       int
	RED         int
}

type ClientLifeCycle struct {
	BankWide    int
	CBLead      int
	LifeCycle   string
	LifeCycleId string
}

type ClientLifeCycles []ClientLifeCycle

func (p ClientLifeCycles) Len() int           { return len(p) }
func (p ClientLifeCycles) Less(i, j int) bool { return p[i].LifeCycleId < p[j].LifeCycleId }
func (p ClientLifeCycles) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (c *DashboardController) GetData(k *knot.WebContext) interface{} {
	// c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	csr, err := c.Ctx.Find(new(BusinessDriverL1Model), tk.M{})
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	BDL1 := make([]BusinessDriverL1Model, 0)
	err = csr.Fetch(&BDL1, 0, false)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()

	MetricTotal := 0
	InitiativesTotal := 0
	InitiativeComparison := tk.M{}
	BankWideIC, CBLedIC, ActiveIC, CompletedIC := 0, 0, 0, 0
	SummaryArray := []tk.M{}
	MetricDirection := tk.M{}
	Declined, Unchanged, Improved := 0, 0, 0
	// Red, Amber, Green := 0, 0, 0

	/*get last actual YTD period*/
	pipePeriod := []tk.M{
		tk.M{
			"$match": tk.M{
				"$and": []tk.M{
					tk.M{
						"actualytd": tk.M{"$ne": 130895111188},
					},
					tk.M{
						"actualytd": tk.M{"$ne": 0},
					},
					tk.M{
						"rag": tk.M{"$ne": ""},
					},
				},
			},
		},
		tk.M{
			"$group": tk.M{
				"_id": "$period",
			},
		},
		tk.M{
			"$sort": tk.M{
				"_id": -1,
			},
		},
	}
	crs, err := c.Ctx.Connection.NewQuery().From(NewBusinessMetricsDataModel().TableName()).Command("pipe", pipePeriod).Cursor(nil)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	dataLastPeriod := tk.M{}
	err = crs.Fetch(&dataLastPeriod, 1, false)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	crs.Close()
	lastActualPeriod := dataLastPeriod.Get("_id").(time.Time)
	/*get last actual YTD period*/

	for _, o := range BDL1 {
		MetricTotal += len(o.BusinessMetric)

		query := tk.M{}.Set("where", db.Eq("SCCategory", o.Idx)).Set("AGGR", "$sum")
		csr2, err := c.Ctx.Find(new(InitiativeModel), query)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
		InitiativesTotal += csr2.Count()
		csr2.Close()

		Summary := tk.M{}
		Summary.Set("Name", o.Name)
		Summary.Set("Metrics", len(o.BusinessMetric))
		Summary.Set("Initiatives", csr2.Count())
		SummaryArray = append(SummaryArray, Summary)

		for _, oo := range o.BusinessMetric {
			/*// for Declined, Unchanged, Improved ...
			dataList := []BusinessMetricsDataModel{}
			data := BusinessMetricsDataModel{}

			csr2, err := c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(db.And(db.Eq("bmid", oo.Id), db.Eq("countrycode", "GLOBAL"))).Order("period").Cursor(nil)
			err = csr2.Fetch(&dataList, 0, false)
			csr2.Close()
			if err != nil {
				BDL1[i].BusinessMetric[ii].CurrentPeriod = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
				BDL1[i].BusinessMetric[ii].CurrentValue = 0
				BDL1[i].BusinessMetric[ii].TargetValue = 0
				BDL1[i].BusinessMetric[ii].BaseLineValue = 0
				BDL1[i].BusinessMetric[ii].MetricFiles = []MetricFile{}
				continue
			}
			if len(dataList) > 0 {
				for idx, i := range dataList {
					if i.ActualYTD != 130895111188 && !i.NAActual {
						data = dataList[idx]
					}
				}
			}
			temp := BDL1[i].BusinessMetric[ii].UpdatedDate
			updatedDate, _ := strconv.Atoi(temp.Format("20060102150405"))
			dataUpdatedDate, _ := strconv.Atoi(data.UpdatedDate.Format("20060102150405"))
			BDL1[i].BusinessMetric[ii].Display = ""
			if data.Year > 0 && (updatedDate == 10101000000 || dataUpdatedDate > updatedDate) {
				// for RAG ...
				if data.RAG == "red" {
					Red += 1
				} else if data.RAG == "amber" {
					Amber += 1
				} else if data.RAG == "green" {
					Green += 1
				}
				// .....
				BDL1[i].BusinessMetric[ii].CurrentValue = data.ActualYTD
				if BDL1[i].BusinessMetric[ii].CurrentValue == 130895111188 {
					BDL1[i].BusinessMetric[ii].CurrentValue = 0
				}
			} else {
				BDL1[i].BusinessMetric[ii].CurrentValue = 0
			}
			// .....
			*/

			if oo.MetricDirection == 0 {
				Unchanged += 1
			} else if oo.MetricDirection == 1 {
				Improved += 1
			} else if oo.MetricDirection == 2 {
				Declined += 1
			}
		}
	}

	result := tk.M{}
	result.Set("Initiatives", InitiativesTotal)
	result.Set("Metrics", MetricTotal)
	result.Set("Summary", SummaryArray)
	// result.Set("BDL1", BDL1)

	//Metric Direction
	MetricDirection.Set("Declined", Declined)
	MetricDirection.Set("Unchanged", Unchanged)
	MetricDirection.Set("Improved", Improved)
	result.Set("MetricDirection", MetricDirection)

	//Metric Status
	err, MetricStatus := c.getMetricStatus(lastActualPeriod)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	result.Set("MetricStatus", MetricStatus)

	//ClientLifeCycle
	err, ClientLifeCycleArr, BankWideIC, CBLedIC := c.getClientLifeCycle()
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	result.Set("ClientLifeCycle", ClientLifeCycleArr)
	// ....

	// InitiativeComparison...
	InitiativeComparison.Set("BankWide", BankWideIC)
	InitiativeComparison.Set("CBLed", CBLedIC)

	csr, err = c.Ctx.Find(new(InitiativeModel), tk.M{})
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	InitiativeList := make([]InitiativeModel, 0)
	err = csr.Fetch(&InitiativeList, 0, false)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()
	TotalInitiatives := 0
	for i, o := range InitiativeList {
		TotalInitiatives += 1
		if o.SetAsComplete {
			CompletedIC += 1
		} else {

			MilestonesArr := []tk.M{}
			TotalDay := 0.0
			for _, oo := range InitiativeList[i].Milestones {
				Milestones := tk.M{}
				oo.Name = strings.Trim(oo.Name, " ")
				Milestones.Set("Name", oo.Name)
				Milestones.Set("StartDate", oo.StartDate)
				Milestones.Set("EndDate", oo.EndDate)
				Milestones.Set("Country", oo.Country)
				DaysBetween := oo.EndDate.Sub(oo.StartDate).Hours() / 24
				TotalDay += DaysBetween

				Milestones.Set("DaysBetween", DaysBetween)
				MilestonesArr = append(MilestonesArr, Milestones)
			}

			//....
			ProgressCompletion := 0.0
			Now := time.Now()
			if len(MilestonesArr) == 0 {
				if Now.After(InitiativeList[i].FinishDate) {
					ProgressCompletion = 100.0
				} else {
					DaysBetweenStartnFinish := o.FinishDate.Sub(o.StartDate).Hours() / 24.0

					DaysBetweenStartnNow := Now.Sub(o.StartDate).Hours() / 24.0

					if DaysBetweenStartnFinish == 0 {
						ProgressCompletion = 0
					} else {
						ProgressCompletion = (DaysBetweenStartnNow / DaysBetweenStartnFinish) * 100.0
					}
				}
				// fmt.Println("if", ProgressCompletion)
			} else {
				// fmt.Println("total milestone", len(InitiativeList[i].Milestones))
				for _, oo := range InitiativeList[i].Milestones {
					if len(oo.Name) != 0 {
						weight := 0.0
						if TotalDay == 0.0 {
							weight = 0.0
						} else {
							weight = (oo.EndDate.Sub(oo.StartDate).Hours() / 24.0) / TotalDay
						}
						// fmt.Println("-", weight, TotalDay, oo.StartDate, oo.EndDate, oo.EndDate.Sub(oo.StartDate).Hours()/24.0)
						// fmt.Println("cuk", oo.Name, TotalDay, weight)

						progress := 0.0
						DaysBetweenStartnEnd := oo.EndDate.Sub(oo.StartDate).Hours() / 24.0
						DaysBetweenStartnNow := Now.Sub(oo.StartDate).Hours() / 24.0
						DaysBetweenEndnNow := Now.Sub(oo.EndDate).Hours() / 24.0
						if Now.After(oo.StartDate) || (DaysBetweenStartnNow < 1 && DaysBetweenStartnNow > -1) {
							if Now.After(oo.EndDate) || (DaysBetweenEndnNow < 1 && DaysBetweenEndnNow > -1) {
								progress = weight
							} else {

								if DaysBetweenStartnEnd == 0 {
									progress = 0.0
								} else {
									progress = (DaysBetweenStartnNow / DaysBetweenStartnEnd) * weight
								}
							}
						}
						// fmt.Println("progress", progress)
						ProgressCompletion += progress * 100
					}
				}
			}
			// fmt.Println(ProgressCompletion)
			if ProgressCompletion >= 100 {
				CompletedIC += 1
			}

		}
	}

	ActiveIC = TotalInitiatives - CompletedIC

	InitiativeComparison.Set("Active", ActiveIC)
	InitiativeComparison.Set("Completed", CompletedIC)
	result.Set("InitiativeComparison", InitiativeComparison)
	// ....

	// NumberOfMetric ..
	/*test speed up by rangga*/
	err, NumberOfMetricArr := c.getNumberOfMetric(BDL1, lastActualPeriod)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	result.Set("NumberOfMetric", NumberOfMetricArr)
	/*test speedu up nek bermasalah balekno awalneh ae :D*/

	// ...
	/*test speed up by rangga*/
	err, NumberOfMetricRegion := c.getNumberOfMetricByRegion(BDL1, lastActualPeriod)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	result.Set("NumberOfMetricRegion", NumberOfMetricRegion)
	/*test speedu up nek bermasalah balekno awalneh ae :D*/

	return c.SetResultInfo(false, "", result)
}

func (c *DashboardController) getMetricStatus(lastActualPeriod time.Time) (error, tk.M) {
	pipe := []tk.M{
		tk.M{
			"$match": tk.M{
				"$and": []tk.M{
					tk.M{
						"actualytd": tk.M{"$ne": 130895111188},
					},
					tk.M{
						"actualytd": tk.M{"$ne": 0},
					},
					tk.M{
						"rag": tk.M{"$ne": ""},
					},
					tk.M{
						"period": lastActualPeriod,
					},
				},
			},
		},
		tk.M{
			"$group": tk.M{
				"_id": "$rag",
				"totalRag": tk.M{
					"$sum": 1,
				},
			},
		},
	}
	crs, err := c.Ctx.Connection.NewQuery().Command("pipe", pipe).From(NewBusinessMetricsDataModel().TableName()).Cursor(nil)
	if !tk.IsNilOrEmpty(err) {
		return err, nil
	}
	data := tk.Ms{}
	err = crs.Fetch(&data, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return err, nil
	}
	crs.Close()
	MetricStatus := tk.M{}
	for _, rag := range data {
		if rag.GetString("_id") == "red" {
			MetricStatus.Set("Red", rag.GetInt("totalRag"))
		} else {
			MetricStatus.Set("Red", 0)
		}
		if rag.GetString("_id") == "amber" {
			MetricStatus.Set("Amber", rag.GetInt("totalRag"))
		} else {
			MetricStatus.Set("Amber", 0)
		}
		if rag.GetString("_id") == "green" {
			MetricStatus.Set("Green", rag.GetInt("totalRag"))
		} else {
			MetricStatus.Set("Green", 0)
		}
	}

	return nil, MetricStatus
}

func (c *DashboardController) getClientLifeCycle() (error, []ClientLifeCycle, int, int) {
	ClientLifeCycleArr := []ClientLifeCycle{}
	BankWideIC, CBLedIC := 0, 0

	csr, err := c.Ctx.Find(new(LifeCycleModel), tk.M{})
	if err != nil {
		return err, nil, BankWideIC, CBLedIC
	}
	LifeCycle := make([]LifeCycleModel, 0)
	err = csr.Fetch(&LifeCycle, 0, false)
	if err != nil {
		return err, nil, BankWideIC, CBLedIC
	}
	csr.Close()

	pipe := []tk.M{
		tk.M{
			"$group": tk.M{
				"_id": tk.M{
					"type":        "$type",
					"LifeCycleId": "$LifeCycleId",
				},
				"totalType": tk.M{"$sum": 1},
			},
		},
		tk.M{
			"$sort": tk.M{
				"_id.LifeCycleId": 1,
			},
		},
	}
	csr2, err := c.Ctx.Connection.NewQuery().Command("pipe", pipe).From(NewInitiativeModel().TableName()).Cursor(nil)
	data := tk.Ms{}
	err = csr2.Fetch(&data, 0, false)
	if err != nil {
		return err, nil, BankWideIC, CBLedIC
	}
	csr2.Close()

	lifeCycleDistinct := map[string]tk.M{}
	for _, o := range LifeCycle {
		for _, v := range data {
			tom := v.Get("_id").(tk.M)
			if tom.GetString("LifeCycleId") == o.LifeCycleId {
				if _, exist := lifeCycleDistinct[o.LifeCycleId]; !exist {
					lifeCycleDistinct[o.LifeCycleId] = tk.M{}.Set(tom.GetString("type"), v.GetInt("totalType")).Set("LifeCycle", o.Name)
				} else {
					lifeCycleDistinct[o.LifeCycleId] = lifeCycleDistinct[o.LifeCycleId].Set(tom.GetString("type"), v.GetInt("totalType"))
				}
			}
		}
	}

	for key, v := range lifeCycleDistinct {
		lc := ClientLifeCycle{}
		lc.BankWide = v.GetInt("BANKWIDE")
		lc.CBLead = v.GetInt("CBLED")
		lc.LifeCycle = v.GetString("LifeCycle")
		lc.LifeCycleId = key
		ClientLifeCycleArr = append(ClientLifeCycleArr, lc)
		CBLedIC += v.GetInt("CBLED")
		BankWideIC += v.GetInt("BANKWIDE")
	}
	sort.Sort(ClientLifeCycles(ClientLifeCycleArr))

	return nil, ClientLifeCycleArr, BankWideIC, CBLedIC
}

func (c *DashboardController) getNumberOfMetric(BDL1 []BusinessDriverL1Model, lastActualPeriod time.Time) (error, []NumberMetric) {
	csr, err := c.Ctx.Connection.NewQuery().From("Region").
		Cursor(nil)
	if err != nil {
		return err, nil
	}
	RegionList := []RegionModel{}
	err = csr.Fetch(&RegionList, 0, false)
	if err != nil {
		return err, nil
	}
	csr.Close()

	NumberOfMetricArr := []NumberMetric{}
	regions, countryFeatRegion := []string{}, []string{}
	for _, region := range RegionList {
		regions = append(regions, region.Country)
		countryFeatRegion = append(countryFeatRegion, region.Country+"|"+region.Major_Region)
	}

	pipe := []tk.M{
		tk.M{
			"$match": tk.M{
				"country": tk.M{"$in": regions},
				"rag":     tk.M{"$ne": ""},
				"period":  lastActualPeriod,
			},
		},
		tk.M{
			"$group": tk.M{
				"_id": tk.M{
					"rag":     "$rag",
					"country": "$country",
					"bmid":    "$bmid",
					"major":   "$majorregion",
				},
				"ragCount": tk.M{
					"$sum": 1,
				},
			},
		},
		tk.M{
			"$sort": tk.M{
				"_id.country": 1,
			},
		},
	}
	csr, err = c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Command("pipe", pipe).Cursor(nil)
	if err != nil {
		return err, nil
	}
	bmdataRegion := []tk.M{}
	err = csr.Fetch(&bmdataRegion, 0, false)
	csr.Close()

	distinctByRegion := map[string]tk.Ms{}
	if len(bmdataRegion) == 0 {
		copyMajorRegionList := make([]string, len(countryFeatRegion))
		copy(copyMajorRegionList, countryFeatRegion)
		for _, v := range copyMajorRegionList {
			distinctByRegion[v] = append(distinctByRegion[v], tk.M{}.Set("", 0))
		}
	}
	for _, bd := range BDL1 {
		for _, bm := range bd.BusinessMetric {
			copyMajorRegionList := make([]string, len(regions))
			copy(copyMajorRegionList, regions)
			for _, v := range bmdataRegion {
				tom := v.Get("_id").(tk.M)
				if tom.GetString("bmid") == bm.Id {
					if _, exist := distinctByRegion[tom.GetString("country")+"|"+tom.GetString("major")]; !exist {
						distinctByRegion[tom.GetString("country")+"|"+tom.GetString("major")] = append(distinctByRegion[tom.GetString("country")+"|"+tom.GetString("major")], tk.M{}.Set(tom.GetString("rag"), v.GetInt("ragCount")))
					} else {
						distinctByRegion[tom.GetString("country")+"|"+tom.GetString("major")] = append(distinctByRegion[tom.GetString("country")+"|"+tom.GetString("major")], tk.M{}.Set(tom.GetString("rag"), v.GetInt("ragCount")))
					}

					_, idx := tk.MemberIndex(copyMajorRegionList, tom.GetString("country"))
					copyMajorRegionList = append(copyMajorRegionList[:idx], copyMajorRegionList[idx+1:]...)
					if len(copyMajorRegionList) > 0 {
						for _, m := range copyMajorRegionList {
							distinctByRegion[m+"|"+tom.GetString("major")] = append(distinctByRegion[m+"|"+tom.GetString("major")], tk.M{}.Set("", 0))
						}
					}
				}
			}
		}
	}

	for key, val := range distinctByRegion {
		splt := strings.Split(key, "|")
		country := splt[0]
		majorregion := splt[1]
		AMBER, GREEN, RED := 0, 0, 0
		for _, rag := range val {
			for ragstring, ragval := range rag {
				if ragstring == "red" {
					RED += ragval.(int)
				} else if ragstring == "amber" {
					AMBER += ragval.(int)
				} else if ragstring == "green" {
					GREEN += ragval.(int)
				}
			}
		}
		NumberOfMetric := NumberMetric{}
		NumberOfMetric.Country = country
		NumberOfMetric.MajorRegion = majorregion
		if AMBER > 0 {
			NumberOfMetric.AMBER = AMBER
		}
		if GREEN > 0 {
			NumberOfMetric.GREEN = GREEN
		}
		if RED > 0 {
			NumberOfMetric.RED = RED
		}
		NumberOfMetricArr = append(NumberOfMetricArr, NumberOfMetric)
	}
	return nil, NumberOfMetricArr
}

func (c *DashboardController) getNumberOfMetricByRegion(BDL1 []BusinessDriverL1Model, lastActualPeriod time.Time) (error, []NumberMetric) {
	MajorRegionList := []string{"AME", "ASA", "GCNA"}
	pipe := []tk.M{
		tk.M{
			"$match": tk.M{
				"country": tk.M{"$in": MajorRegionList},
				"rag":     tk.M{"$ne": ""},
				"period":  lastActualPeriod,
			},
		},
		tk.M{
			"$group": tk.M{
				"_id": tk.M{
					"rag":     "$rag",
					"country": "$country",
					"bmid":    "$bmid",
				},
				"ragCount": tk.M{
					"$sum": 1,
				},
			},
		},
		tk.M{
			"$sort": tk.M{
				"_id.country": 1,
			},
		},
	}
	csr, err := c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Command("pipe", pipe).Cursor(nil)
	if err != nil {
		return err, nil
	}
	bmdata := []tk.M{}
	err = csr.Fetch(&bmdata, 0, false)
	csr.Close()
	distinctByCountry := map[string]tk.Ms{}
	if len(bmdata) == 0 {
		copyMajorRegionList := make([]string, len(MajorRegionList))
		copy(copyMajorRegionList, MajorRegionList)
		for _, v := range copyMajorRegionList {
			distinctByCountry[v] = append(distinctByCountry[v], tk.M{}.Set("", 0))
		}
	}
	for _, bd := range BDL1 {
		for _, bm := range bd.BusinessMetric {
			copyMajorRegionList := make([]string, len(MajorRegionList))
			copy(copyMajorRegionList, MajorRegionList)
			for _, v := range bmdata {
				tom := v.Get("_id").(tk.M)
				if tom.GetString("bmid") == bm.Id {
					if _, exist := distinctByCountry[tom.GetString("country")]; !exist {
						a := tk.M{}.Set(tom.GetString("rag"), v.GetInt("ragCount"))
						distinctByCountry[tom.GetString("country")] = append(distinctByCountry[tom.GetString("country")], a)
					} else {
						distinctByCountry[tom.GetString("country")] = append(distinctByCountry[tom.GetString("country")], tk.M{}.Set(tom.GetString("rag"), v.GetInt("ragCount")))
					}
					_, idx := tk.MemberIndex(copyMajorRegionList, tom.GetString("country"))
					copyMajorRegionList = append(copyMajorRegionList[:idx], copyMajorRegionList[idx+1:]...)
					if len(copyMajorRegionList) > 0 {
						for _, m := range copyMajorRegionList {
							distinctByCountry[m] = append(distinctByCountry[m], tk.M{}.Set("", 0))
						}
					}
				}
			}
		}
	}
	NumberOfMetricRegion := []NumberMetric{}
	for key, val := range distinctByCountry {
		AMBER, GREEN, RED := 0, 0, 0
		for _, rag := range val {
			for ragstring, ragval := range rag {
				if ragstring == "red" {
					RED += ragval.(int)
				} else if ragstring == "amber" {
					AMBER += ragval.(int)
				} else if ragstring == "green" {
					GREEN += ragval.(int)
				}
			}
		}
		NumberOfMetric := NumberMetric{}
		NumberOfMetric.Country = key
		NumberOfMetric.MajorRegion = key
		if AMBER > 0 {
			NumberOfMetric.AMBER = AMBER
		}
		if GREEN > 0 {
			NumberOfMetric.GREEN = GREEN
		}
		if RED > 0 {
			NumberOfMetric.RED = RED
		}
		NumberOfMetricRegion = append(NumberOfMetricRegion, NumberOfMetric)
	}
	return nil, NumberOfMetricRegion
}

func (c *DashboardController) SaveAsPdf(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	// ===== parse payload

	payload := struct {
		ImageData   string
		ImageHeight float64
		ImageWidth  float64
	}{}
	err := k.GetPayload(&payload)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	// ===== prepare path of destination image

	imageTemporaryPath := ""
	config := helper.ReadConfig()
	if value, ok := config["downloadPath"]; ok {
		imageTemporaryPath = filepath.Join(value, fmt.Sprintf("%s.png", tk.RandomString(18)))
	}
	fmt.Println("temporary image created @", imageTemporaryPath)
	defer os.Remove(imageTemporaryPath)

	// ===== convert payload base64 into image file

	err = tk.Base64ToImage(payload.ImageData, imageTemporaryPath)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	// ===== resize image

	imageTemporaryResizedPath, err := (func() (string, error) {
		imageTemporaryResizedPath := strings.Replace(imageTemporaryPath, ".png", "-resized.png", -1)
		f, err := os.Open(imageTemporaryPath)
		if err != nil {
			return "", err
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			return "", err
		}

		m := imaging.Resize(img, 1190, 0, imaging.Lanczos)

		err = imaging.Save(m, imageTemporaryResizedPath)
		if err != nil {
			return "", err
		}

		return imageTemporaryResizedPath, nil
	})()
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	defer os.Remove(imageTemporaryResizedPath)

	// ===== prepare path of downloadable pdf

	downloadablePdfName := fmt.Sprintf("%s.pdf", tk.RandomString(18))
	downloadablePdfPath := ""
	if value, ok := config["downloadPath"]; ok {
		downloadablePdfPath = filepath.Join(value, downloadablePdfName)
	}

	// ===== embed image into pdf

	pdf := gofpdf.New("P", "mm", "a2", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 11)
	pdf.Image(imageTemporaryResizedPath, 50, 4, 0, 0, false, "", 0, "")
	err = pdf.OutputFileAndClose(downloadablePdfPath)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	// c.Action(k, "Save Scorecard to PDF")
	c.Action(k, "Dashboard", "Save Scorecard to PDF", "", "", "", "", "")

	downloadLink := fmt.Sprintf("/static/download/%s", downloadablePdfName)
	return tk.M{"downloadLink": downloadLink}
}

func checkGlobalValue(val []interface{}) bool {
	for _, data := range val {
		if data == "GLOBAL" {
			return true
		}
	}
	return false
}

func (c *DashboardController) GetPanelData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	InitiativeAccess := c.GetAccess(k, "INITIATIVE")
	Initiative := new(AccessibilityModel)
	e := tk.MtoStruct(InitiativeAccess, &Initiative)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Low                     bool
		Medium                  bool
		High                    bool
		Primary                 bool
		Secondary               bool
		CBLead                  bool
		BankWide                bool
		YtdComplete             bool
		Remaining               bool
		Task                    bool
		IsExellerator           bool
		IsOperationalExcellence bool
		Region                  []interface{}
		Country                 []interface{}
		IsBE                    bool
		IsBP                    bool
		BDFilter                []interface{}
		StartDate               string
		EndDate                 string
		DisplayColor            string
		SCRegion                string
		SCCountry               string
	}{}
	e = k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	ds := []tk.M{}
	csr, e := c.Ctx.Connection.NewQuery().From("MasterLifeCycle").Order("Seq").Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&ds, 0, false)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	csr.Close()

	querybd := []*db.Filter{}
	querybd = append(querybd, db.Eq("Type", "Key Enablers"))

	if parm.IsBE != parm.IsBP {
		if parm.IsBE {
			querybd = append(querybd, db.Eq("category", "Business Enablement"))
		}
		if parm.IsBP {
			querybd = append(querybd, db.Eq("category", "Business Protection"))
		}
	}

	dtIntv := make([]SummaryBusinessDriverModel, 0)
	csr, e = c.Ctx.Find(new(SummaryBusinessDriverModel), tk.M{}.Set("where", db.And(querybd...)).Set("order", []string{"Seq"}))
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&dtIntv, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	if parm.SCRegion == "" || parm.SCRegion == "Region" {
		parm.SCRegion = "GLOBAL"
	}
	if parm.SCCountry == "" || parm.SCCountry == "Country" {
		parm.SCCountry = parm.SCRegion
	}
	query := []*db.Filter{}

	UserCountry := ""
	if k.Session("country") != nil {
		UserCountry = k.Session("country").(string)
	}

	for bd := range dtIntv {
		for bm := range dtIntv[bd].BusinessMetrics {
			BMDataList, data := []BusinessMetricsDataModel{}, BusinessMetricsDataModel{}

			query = append(query[0:0], db.Eq("bmid", dtIntv[bd].BusinessMetrics[bm].Id))

			if UserCountry == "" {

				if parm.SCRegion == parm.SCCountry {
					query = append(query, db.Eq("majorregion", parm.SCRegion))
				}
				query = append(query, db.Eq("country", parm.SCCountry))
			} else {
				query = append(query, db.Eq("country", UserCountry))
			}
			csr, err := c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(db.And(query...)).Order("period").Cursor(nil)
			err = csr.Fetch(&BMDataList, 0, false)
			csr.Close()
			if err != nil {
				dtIntv[bd].BusinessMetrics[bm].CurrentValue = 0
				dtIntv[bd].BusinessMetrics[bm].TargetValue = 0
				dtIntv[bd].BusinessMetrics[bm].BaseLineValue = 0
				dtIntv[bd].BusinessMetrics[bm].MetricFiles = []MetricFile{}
				continue
			}
			temp := dtIntv[bd].BusinessMetrics[bm].UpdatedDate
			updatedDate, _ := strconv.Atoi(temp.Format("20060102150405"))
			if len(BMDataList) > 0 {
				// data = BMDataList[0]
				// tk.Println(dtIntv[bd].BusinessMetrics[bm].Id)
				for idx, i := range BMDataList {
					// tk.Println("# ", i.Period, "|", i.ActualYTD)
					if i.ActualYTD != 130895111188 && !i.NAActual {
						data = BMDataList[idx]
					}
				}
				// tk.Println(data.Period, "|", updatedDate)
			}
			dataUpdatedDate, _ := strconv.Atoi(data.UpdatedDate.Format("20060102150405"))
			if len(BMDataList) > 0 && (updatedDate == 10101000000 || dataUpdatedDate > updatedDate) {
				dtIntv[bd].BusinessMetrics[bm].CurrentValue = data.ActualYTD
				dtIntv[bd].BusinessMetrics[bm].CurrentPeriod = data.Period
				dtIntv[bd].BusinessMetrics[bm].CurrentPeriodStr = data.Period.Format("20060102")
				dtIntv[bd].BusinessMetrics[bm].TargetValue = data.Target
				dtIntv[bd].BusinessMetrics[bm].TargetPeriod = data.Period
				dtIntv[bd].BusinessMetrics[bm].TargetPeriodStr = data.Period.Format("20060102")
				dtIntv[bd].BusinessMetrics[bm].BaseLineValue = data.Baseline
				dtIntv[bd].BusinessMetrics[bm].BaseLinePeriod = data.Period
				dtIntv[bd].BusinessMetrics[bm].BaseLinePeriodStr = data.Period.Format("20060102")
				dtIntv[bd].BusinessMetrics[bm].Display = data.RAG
			} else {
				dtIntv[bd].BusinessMetrics[bm].CurrentValue = 0
				dtIntv[bd].BusinessMetrics[bm].TargetValue = 0
				dtIntv[bd].BusinessMetrics[bm].BaseLineValue = 0
			}

			dtIntv[bd].BusinessMetrics[bm].ActualData = []ActualValue{}
			for _, BMData := range BMDataList {
				actual := ActualValue{}
				actual.Period = BMData.Period
				actual.PeriodStr = BMData.Period.Format("20060102")
				actual.Value = BMData.ActualYTD
				if dtIntv[bd].BusinessMetrics[bm].CurrentPeriod.After(BMData.Period) || dtIntv[bd].BusinessMetrics[bm].CurrentPeriod == BMData.Period {
					dtIntv[bd].BusinessMetrics[bm].ActualData = append(dtIntv[bd].BusinessMetrics[bm].ActualData, actual)
				}

			}
		}
	}

	query = []*db.Filter{}

	query = append(query, db.Eq("InitiativeType", "KeyEnablers"))
	if !parm.Primary || !parm.Secondary {
		if parm.Primary {
			query = append(query, db.Eq("BusinessDriverImpact", "Primary"))
		}
		if parm.Secondary {
			query = append(query, db.Eq("BusinessDriverImpact", "Secondary"))
		}
		if !parm.Primary && !parm.Secondary {
			query = append(query, db.Ne("BusinessDriverImpact", "Primary"))
			query = append(query, db.Ne("BusinessDriverImpact", "Secondary"))
		}
	}
	if len(parm.Region) > 0 {
		if checkGlobalValue(parm.Region) {
			query = append(query, db.Or(db.Eq("IsGlobal", true), db.In("Region", parm.Region...)))
		} else {
			query = append(query, db.In("Region", parm.Region...))
		}
	}
	if len(parm.Country) > 0 {
		if checkGlobalValue(parm.Country) {
			query = append(query, db.Or(db.Eq("IsGlobal", true), db.In("Country", parm.Country...)))
		} else {
			query = append(query, db.In("Country", parm.Country...))
		}
	}
	// if checkGlobalValue(parm.Country) || checkGlobalValue(parm.Region){
	// 	query = append(query, db.Eq("IsGlobal", true))
	// }
	if len(parm.BDFilter) > 0 {
		query = append(query, db.In("BusinessDriverId", parm.BDFilter...))
	}
	if !parm.Low {
		query = append(query, db.Ne("BusinessImpact", "Low"))
	}
	if !parm.Medium {
		query = append(query, db.Ne("BusinessImpact", "Medium"))
	}
	if !parm.High {
		query = append(query, db.Ne("BusinessImpact", "High"))
	}
	if !parm.CBLead {
		query = append(query, db.Ne("type", "CBLED"))
	}
	if !parm.BankWide {
		query = append(query, db.Ne("type", "BANKWIDE"))
	}
	if !parm.IsExellerator {
		query = append(query, db.Ne("EX", false))
	}
	if !parm.IsOperationalExcellence {
		query = append(query, db.Ne("OE", false))
	}
	if parm.StartDate != "" && parm.EndDate != "" {
		startdateTmp, _ := time.Parse(time.RFC3339, parm.StartDate)
		enddateTmp, _ := time.Parse(time.RFC3339, parm.EndDate)
		query = append(query, db.Or(db.And(db.Gte("FinishDate", startdateTmp), db.Lte("FinishDate", enddateTmp)), db.And(db.Gte("StartDate", startdateTmp), db.Lte("StartDate", enddateTmp))))
	} else if parm.StartDate != "" {
		startdateTmp, _ := time.Parse(time.RFC3339, parm.StartDate)
		query = append(query, db.Gte("FinishDate", startdateTmp))
	} else if parm.EndDate != "" {
		enddateTmp, _ := time.Parse(time.RFC3339, parm.EndDate)
		query = append(query, db.Lte("FinishDate", enddateTmp))
	}

	if parm.DisplayColor != "" {
		query = append(query, db.Eq("DisplayProgress", parm.DisplayColor))
	}
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
	csr, e = c.Ctx.Connection.NewQuery().Where(db.And(query...)).From("Initiative").Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	InitiativeList := []tk.M{}
	e = csr.Fetch(&InitiativeList, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	if Initiative.Global.Owned || Initiative.Region.Owned || Initiative.Country.Owned {
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

	for _, v := range InitiativeList {
		for key, _ := range v.Get("Attachments").([]interface{}) {
			if strings.Contains(v.Get("Attachments").([]interface{})[key].(tk.M).GetString("filename"), "+") {
				v.Get("Attachments").([]interface{})[key].(tk.M)["filename"] = strings.Split(v.Get("Attachments").([]interface{})[key].(tk.M).GetString("filename"), "+")[1]
			}
		}
	}

	query = []*db.Filter{}
	query = append(query, db.Eq("InitiativeType", "SupportingEnablers"))
	if !parm.Primary || !parm.Secondary {
		if parm.Primary {
			query = append(query, db.Eq("BusinessDriverImpact", "Primary"))
		}
		if parm.Secondary {
			query = append(query, db.Eq("BusinessDriverImpact", "Secondary"))
		}
	}

	// Task
	TaskList := []TaskModel{}
	queryTask := []*db.Filter{}
	queryTask = append(queryTask, db.Eq("tasktype", "KeyEnablers"))
	if len(parm.BDFilter) > 0 {
		queryTask = append(queryTask, db.In("businessdriverid", parm.BDFilter...))
	}
	if len(parm.Region) > 0 {
		queryTask = append(queryTask, db.Or(db.Eq("IsGlobal", true), db.In("Region", parm.Region...)))
	}
	if len(parm.Country) > 0 {
		queryTask = append(queryTask, db.Or(db.Eq("IsGlobal", true), db.In("Country", parm.Country...)))
	}
	if parm.DisplayColor != "" {
		queryTask = append(queryTask, db.Eq("DisplayProgress", parm.DisplayColor))
	}
	if parm.Task && parm.IsExellerator && parm.IsOperationalExcellence {

		csr, e = c.Ctx.Connection.NewQuery().Where(db.And(queryTask...)).From("Task").Cursor(nil)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
		e = csr.Fetch(&TaskList, 0, false)
		csr.Close()
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}

	}

	result := tk.M{}
	allBD := make([]SummaryBusinessDriverModel, 0)
	csr, e = c.Ctx.Find(new(SummaryBusinessDriverModel), nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&allBD, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	var arrsort []string
	arrsort = append(arrsort, "seq")
	csr, err := c.Ctx.Find(new(BusinessDriverL1Model), tk.M{}.Set("order", arrsort))
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	BDL1 := make([]BusinessDriverL1Model, 0)
	err = csr.Fetch(&BDL1, 0, false)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()

	result.Set("SummaryBusinessDriver", dtIntv)
	result.Set("BusinessDriverL1", BDL1)
	result.Set("MasterLifeCycle", ds)
	result.Set("Project", InitiativeList)
	result.Set("TaskList", TaskList)
	result.Set("AllSummaryBusinessDriver", allBD)

	OwnedInitiativesData := []string{}
	for _, i := range OwnedInitiatives {
		o := i.Get("_id").(bson.ObjectId).Hex()
		OwnedInitiativesData = append(OwnedInitiativesData, o)
	}
	result.Set("OwnedInitiative", OwnedInitiativesData)

	return result
}

func (c *DashboardController) MoveUpdate(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	p := InitiativeModel{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}
	//c.NewQuery().Delete().From(new(GenDataBrowserNotInTmp).TableName()).Where(db.Eq("ID", genIDTempTable)).Exec(nil)
	c.Ctx.Connection.NewQuery().Delete().From("Initiative").Where(db.Eq("InitiativeID", p.InitiativeID)).Exec(nil)

	p.Id = bson.NewObjectId()
	e = c.Ctx.Save(&p)
	if e != nil {
		c.WriteLog(e)
	}

	result := tk.M{}

	result.Set("Res", "OK")

	return result
}

func (c *DashboardController) SummaryBusinessDriverSave(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	p := SummaryBusinessDriverModel{}

	e := k.GetPayload(&p)
	if e != nil {
		tk.Println(e)
		return e
	}

	if p.Id == "" {
		p.Id = bson.NewObjectId()
	}

	e = c.Ctx.Save(&p)
	if e != nil {
		tk.Println(e)
		return e
	}

	return "success"
}

func (c *DashboardController) CommentSave(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	p := InitiativeModel{}

	e := k.GetPayload(&p)
	if e != nil {
		tk.Println(e)
		return e
	}

	if p.Id == "" {
		return "error"
	} else {
		result := new(InitiativeModel)
		e = c.Ctx.GetById(result, p.Id)
		if e != nil {
			return e
		}

		result.CommentList = p.CommentList

		e = c.Ctx.Save(result)
		if e != nil {
			return e
		}

		return result
	}
}

func (c *DashboardController) LogCommentSave(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	p := LogCommentModel{}

	e := k.GetPayload(&p)
	if e != nil {
		tk.Println(e)
		return e
	}

	p.Id = bson.NewObjectId()
	e = c.Ctx.Save(&p)
	if e != nil {
		c.WriteLog(e)
	}

	result := tk.M{}

	result.Set("Res", "OK")

	return result
}

func (c *DashboardController) LogCommentGet(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	p := Comment{}

	e := k.GetPayload(&p)
	if e != nil {
		tk.Println(e)
		return e
	}

	csr, err := c.Ctx.Find(new(LogCommentModel), tk.M{}.Set("where", db.And(db.Eq("comment.Username", p.Username), db.Eq("comment.DateInput", p.DateInput))))
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	LogComment := make([]LogCommentModel, 0)
	err = csr.Fetch(&LogComment, 0, false)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()

	result := tk.M{}
	result.Set("Res", "OK").Set("Data", LogComment)

	return result
}

func (c *DashboardController) CountHeader(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Low      bool
		Medium   bool
		High     bool
		CBLead   bool
		BankWide bool
		Task     bool
		FromDate string
		ToDate   string
		Region   []interface{}
		Country  []interface{}
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	result := tk.M{}
	querybd := []*db.Filter{}
	querybdofCB := []*db.Filter{}
	querybdofBankWide := []*db.Filter{}
	querybdofRemaining := []*db.Filter{}
	querybdofYTDComplate := []*db.Filter{}
	querybdofHigh := []*db.Filter{}
	querybdofMedium := []*db.Filter{}

	if len(parm.Region) != 0 {
		querybd = append(querybd, db.In("Region", parm.Region...))
	}

	if len(parm.Country) != 0 {
		querybd = append(querybd, db.In("Country", parm.Country...))
	}

	if parm.FromDate != "" && parm.ToDate != "" {
		FromDate, _ := time.Parse(time.RFC3339, parm.FromDate)
		ToDate, _ := time.Parse(time.RFC3339, parm.ToDate)
		querybdofCB = append(querybdofCB, db.And(db.Gte("FinishDate", FromDate), db.Lte("FinishDate", ToDate)))
		querybdofBankWide = append(querybdofBankWide, db.And(db.Gte("FinishDate", FromDate), db.Lte("FinishDate", ToDate)))
		querybdofHigh = append(querybdofHigh, db.And(db.Gte("FinishDate", FromDate), db.Lte("FinishDate", ToDate)))
		querybdofMedium = append(querybdofMedium, db.And(db.Gte("FinishDate", FromDate), db.Lte("FinishDate", ToDate)))
		querybdofRemaining = append(querybdofRemaining, db.And(db.Gte("FinishDate", FromDate), db.Lte("FinishDate", ToDate)))
		querybdofYTDComplate = append(querybdofYTDComplate, db.And(db.Gte("FinishDate", FromDate), db.Lte("FinishDate", ToDate)))

	} else {
		querybdofRemaining = append(querybdofRemaining, db.Lte("FinishDate", time.Date(2016, 12, 31, 0, 0, 0, 0, time.UTC)))
		querybdofRemaining = append(querybdofRemaining, db.Gte("FinishDate", time.Now()))
		querybdofYTDComplate = append(querybdofYTDComplate, db.Gte("FinishDate", time.Date(2016, 01, 01, 0, 0, 0, 0, time.UTC)))
		querybdofYTDComplate = append(querybdofYTDComplate, db.Lte("FinishDate", time.Now()))
	}

	querybdofCB = append(querybdofCB, db.Eq("type", "CBLED"))
	querybdofBankWide = append(querybdofBankWide, db.Eq("type", "BANKWIDE"))
	querybdofHigh = append(querybdofHigh, db.Eq("BusinessImpact", "High"))
	querybdofMedium = append(querybdofMedium, db.Eq("BusinessImpact", "Medium"))

	queryCbled := tk.M{}.Set("where", db.And(querybd...)).Set("where", db.And(querybdofCB...))
	queryBankWide := tk.M{}.Set("where", db.And(querybd...)).Set("where", db.And(querybdofBankWide...))
	queryHigh := tk.M{}.Set("where", db.And(querybd...)).Set("where", db.And(querybdofHigh...))
	queryMedium := tk.M{}.Set("where", db.And(querybd...)).Set("where", db.And(querybdofMedium...))
	queryRemaining := tk.M{}.Set("where", db.And(querybd...)).Set("where", db.And(querybdofRemaining...))
	queryYTDComplate := tk.M{}.Set("where", db.And(querybd...)).Set("where", db.And(querybdofYTDComplate...))

	crscbled, ex := c.Ctx.Find(new(InitiativeModel), queryCbled)
	defer crscbled.Close()
	if crscbled == nil {
		return c.SetResultInfo(true, "Cursor Not initialized..", nil)
	}
	dtCbled := []tk.M{}
	ex = crscbled.Fetch(&dtCbled, 0, false)
	if ex != nil {
		return c.SetResultInfo(true, ex.Error(), nil)
	}

	crscBankWide, ex := c.Ctx.Find(new(InitiativeModel), queryBankWide)
	defer crscBankWide.Close()
	if crscBankWide == nil {
		return c.SetResultInfo(true, "Cursor Not initialized..", nil)
	}
	dtBankWide := []tk.M{}

	ex = crscBankWide.Fetch(&dtBankWide, 0, false)
	if ex != nil {
		return c.SetResultInfo(true, ex.Error(), nil)
	}

	crscHigh, ex := c.Ctx.Find(new(InitiativeModel), queryHigh)
	defer crscHigh.Close()
	if crscHigh == nil {
		return c.SetResultInfo(true, "Cursor Not initialized..", nil)
	}
	dtHigh := []tk.M{}

	ex = crscHigh.Fetch(&dtHigh, 0, false)
	if ex != nil {
		return c.SetResultInfo(true, ex.Error(), nil)
	}

	crscMedium, ex := c.Ctx.Find(new(InitiativeModel), queryMedium)
	defer crscMedium.Close()
	if crscMedium == nil {
		return c.SetResultInfo(true, "Cursor Not initialized..", nil)
	}
	dtMedium := []tk.M{}

	ex = crscMedium.Fetch(&dtMedium, 0, false)
	if ex != nil {
		return c.SetResultInfo(true, ex.Error(), nil)
	}

	crsRemaining, ex := c.Ctx.Find(new(InitiativeModel), queryRemaining)
	defer crsRemaining.Close()
	if crsRemaining == nil {
		return c.SetResultInfo(true, "Cursor Not initialized..", nil)
	}
	dtRemaining := []tk.M{}

	ex = crsRemaining.Fetch(&dtRemaining, 0, false)
	if ex != nil {
		return c.SetResultInfo(true, ex.Error(), nil)
	}

	crsYTDComplate, ex := c.Ctx.Find(new(InitiativeModel), queryYTDComplate)
	defer crsYTDComplate.Close()
	if crsYTDComplate == nil {
		return c.SetResultInfo(true, "Cursor Not initialized..", nil)
	}
	dtYTDComplate := []tk.M{}

	ex = crsYTDComplate.Fetch(&dtYTDComplate, 0, false)
	if ex != nil {
		return c.SetResultInfo(true, ex.Error(), nil)
	}
	var countCBled = 0
	var countBankWide = 0
	var countHigh = 0
	var countMedium = 0

	if parm.Medium != false {
		countMedium = len(dtMedium)
	}
	if parm.High != false {
		countHigh = len(dtHigh)
	}
	if parm.CBLead != false {
		countCBled = len(dtCbled)
	}
	if parm.BankWide != false {
		countBankWide = len(dtBankWide)
	}

	result.Set("NumberofCBled ", countCBled)
	result.Set("NumberofBankwide ", countBankWide)
	result.Set("NumberofHigh ", countHigh)
	result.Set("NumberofMedium  ", countMedium)
	result.Set("NumberofRemaining  ", len(dtRemaining))
	result.Set("NumberofYTDComplate  ", len(dtYTDComplate))
	return c.SetResultInfo(false, "Total Header", result)
}

func (c *DashboardController) ScorecardInitiativeDetail(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	p := []string{}

	e := k.GetPayload(&p)
	if e != nil {
		tk.Println(e)
		return e
	}

	parm := []interface{}{}
	for _, v := range p {
		parm = append(parm, bson.ObjectIdHex(v))
	}

	crs, err := c.Ctx.Find(new(InitiativeModel), tk.M{}.Set("where", db.In("_id", parm...)))
	defer crs.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	data := make([]InitiativeModel, 0)
	err = crs.Fetch(&data, 0, false)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	return c.SetResultInfo(false, "Data", data)

	// if p.Id == "" {
	// 	return "error"
	// } else {
	// 	result := new(InitiativeModel)
	// 	e = c.Ctx.GetById(result, p.Id)
	// 	if e != nil {
	// 		return e
	// 	}

	// 	result.CommentList = p.CommentList

	// 	e = c.Ctx.Save(result)
	// 	if e != nil {
	// 		return e
	// 	}

	// 	return result
	// }
}
