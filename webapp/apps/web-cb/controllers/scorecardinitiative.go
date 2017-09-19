package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/models"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"sort"
	"time"
	// "gopkg.in/mgo.v2/bson"
)

type ScorecardInitiativeController struct {
	*BaseController
}

func (c *ScorecardInitiativeController) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	// k.Config.LayoutTemplate = ""
	// k.Config.IncludeFiles = []string{"aclsysadmin/pagecontent.html", "aclsysadmin/user.html", "aclsysadmin/session.html", "aclsysadmin/access.html", "aclsysadmin/group.html", "aclsysadmin/changePassword.html", "dashboard/search.html", "dashboard/chart.html", "dashboard/initiativeTab.html", "dashboard/scorecard.html", "dashboard/scorecard_bm.html",  "dashboard/detailbusinessdriver.html", "dashboard/task.html", "dashboard/content.html", "dashboard/initiative.html" ,"dashboard/modalclone.html"}
	return ""
}

type ScInitiatives struct {
	Id              int
	Idx             string
	Name            string
	Description     string
	Updatedby       string
	Updateddate     time.Time
	Seq             int
	Businessmetrics []BusinessMetricData
}
type ScInitiative []ScInitiatives

func (p ScInitiative) Len() int           { return len(p) }
func (p ScInitiative) Less(i, j int) bool { return p[i].Seq < p[j].Seq }
func (p ScInitiative) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type BusinessMetricData struct {
	SCName      string
	Id          string
	DataPoint   string
	Description string
	MetricType  string
	Display     string
	Type        string
	Orderindex  int
	Initiative  tk.Ms
	Na          int
	Direct      int
	Indirect    int
	RAG         []DataRag
	Quarters    []Quarter
	LastPeriod  time.Time
}
type BusinessMetricDataList []BusinessMetricData

func (p BusinessMetricDataList) Len() int           { return len(p) }
func (p BusinessMetricDataList) Less(i, j int) bool { return p[i].Orderindex < p[j].Orderindex }
func (p BusinessMetricDataList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (c *ScorecardInitiativeController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Region  string
		Country string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	if parm.Region == "" || parm.Region == "Region" {
		parm.Region = "GLOBAL"
	}
	if parm.Country == "" || parm.Country == "Country" {
		parm.Country = parm.Region
	}

	result := ResultInfo{}

	crs, err := c.Ctx.Find(NewBusinessDriverL1Model(), tk.M{})
	defer crs.Close()
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	o := []BusinessDriverL1Model{}
	err = crs.Fetch(&o, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	data := []ScInitiatives{}
	query := []*dbox.Filter{}
	for _, v := range o {
		bms := []BusinessMetricData{}
		for _, vs := range v.BusinessMetric {
			err, r := c.findByInitiativeId(vs)
			if !tk.IsNilOrEmpty(err) {
				return c.ErrorResultInfo(err.Error(), nil)
			}

			///// get RAG
			err, rag := c.findRAG(vs, parm.Region, parm.Country)
			if !tk.IsNilOrEmpty(err) {
				return c.ErrorResultInfo(err.Error(), nil)
			}

			///// get Quarters
			err, quarter := c.findQuarters(vs, parm.Region, parm.Country)
			if !tk.IsNilOrEmpty(err) {
				return c.ErrorResultInfo(err.Error(), nil)
			}

			vs.Display = ""
			Now := time.Now()
			StartPeriod := time.Date(Now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
			LastPeriod := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
			BMData := []BusinessMetricsDataModel{}
			query = append(query, dbox.Gte("bmid", StartPeriod))
			query = append(query[0:0], dbox.Eq("bmid", vs.Id))
			if parm.Region == parm.Country {
				query = append(query, dbox.Eq("majorregion", parm.Region))
			}
			query = append(query, dbox.Eq("country", parm.Country))
			csr, err := c.Ctx.Connection.NewQuery().From(new(BusinessMetricsDataModel).TableName()).Where(dbox.And(query...)).Order("period").Cursor(nil)
			if err != nil {
				vs.Display = ""
			} else {
				err = csr.Fetch(&BMData, 0, false)
				if err != nil {
					vs.Display = ""
				} else {
					for _, x := range BMData {
						if x.ActualYTD != 130895111188 && !x.NAActual {
							vs.Display = x.RAG
							LastPeriod = x.Period.UTC()
						}
					}
				}
			}
			csr.Close()
			bm := BusinessMetricData{
				SCName:      v.Name,
				Id:          vs.Id,
				DataPoint:   vs.DataPoint,
				Description: vs.Description,
				MetricType:  vs.MetricType,
				Display:     vs.Display,
				Orderindex:  vs.OrderIndex,
				Initiative:  r,
				RAG:         rag,
				Quarters:    quarter,
				LastPeriod:  LastPeriod,
				Type:        vs.Type,
			}
			if len(r) > 0 {
				for _, vr := range r {
					switch vr.GetInt("type") {
					case 0:
						bm.Na += 1
					case 1:
						bm.Direct += 1
					case 2:
						bm.Indirect += 1
					}
				}
			} else {
				bm.Na = 0
				bm.Direct = 0
				bm.Indirect = 0
			}
			bms = append(bms, bm)
		}

		sort.Sort(BusinessMetricDataList(bms))
		d := ScInitiatives{
			Id:              v.Id,
			Idx:             v.Idx,
			Name:            v.Name,
			Description:     v.Description,
			Updatedby:       v.UpdatedBy,
			Updateddate:     v.UpdatedDate,
			Seq:             v.Seq,
			Businessmetrics: bms,
		}
		data = append(data, d)
	}
	sort.Sort(ScInitiative(data))
	result.Data = data

	return result
}

func (c *ScorecardInitiativeController) findByInitiativeId(bs BusinessMetric) (error, tk.Ms) {
	pipe := []tk.M{
		tk.M{
			"$match": tk.M{
				"KeyMetrics": tk.M{"$elemMatch": tk.M{"BMId": bs.Id}},
			},
		},
	}
	crs, err := c.Ctx.Connection.NewQuery().Command("pipe", pipe).From("Initiative").Cursor(nil)
	defer crs.Close()
	if !tk.IsNilOrEmpty(err) {
		return err, nil
	}

	resultInitiative := []InitiativeModel{}
	err = crs.Fetch(&resultInitiative, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return err, nil
	}

	result := tk.Ms{}
	for _, v := range resultInitiative {
		aepm := ""
		for x, ae := range v.AccountableExecutive {
			if x == 0 {
				aepm += ae
			} else {
				aepm += ", " + ae
			}
		}
		if aepm == "" {
			for x, pm := range v.ProjectManager {
				if x == 0 {
					aepm += pm
				} else {
					aepm += ", " + pm
				}
			}
		} else {
			aepm += " / "
			for x, pm := range v.ProjectManager {
				if x == 0 {
					aepm += pm
				} else {
					aepm += ", " + pm
				}
			}
		}
		for _, vs := range v.KeyMetrics {
			if vs.BMId == bs.Id {
				result = append(result, tk.M{}.Set("id", v.Id.Hex()).Set("name", v.ProjectName).Set("type", vs.DirectIndirect).Set("aepm", aepm))
				// if _, exist := getKeyMetrics[v.Id.Hex()]; !exist {
				// 	getKeyMetrics[v.Id.Hex()] = tk.M{}.Set("id", v.Id).Set("type", vs.DirectIndirect)
				// } else {
				// 	getKeyMetrics[v.Id.Hex()].Set("id", v.Id).Set("type", vs.DirectIndirect)
				// }
			}
		}
	}

	// for _, v := range getKeyMetrics {
	// 	result = append(result, tk.M{}.Set("id", v.Get("id").(bson.ObjectId).Hex()).Set("type", v.GetInt("type")))
	// }

	return nil, result
}

type Quarter struct {
	QuarterName string
	Quarters    []DataRag
}
type Quarters []Quarter

func (p Quarters) Len() int           { return len(p) }
func (p Quarters) Less(i, j int) bool { return p[i].QuarterName < p[j].QuarterName }
func (p Quarters) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type DataRag struct {
	BmId           string
	BusinessMetric string
	Period         time.Time
	Rag            string
}
type DataRags []DataRag

func (p DataRags) Len() int           { return len(p) }
func (p DataRags) Less(i, j int) bool { return p[i].Period.Before(p[j].Period) }
func (p DataRags) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (c *ScorecardInitiativeController) findRAG(bs BusinessMetric, region, country string) (error, []DataRag) {
	ds, result := []DataRag{}, []DataRag{}
	var filter = []*dbox.Filter{}
	filter = append(filter, dbox.Eq("bmid", bs.Id))
	if region == country {
		filter = append(filter, dbox.Eq("majorregion", region))
	}
	filter = append(filter, dbox.Eq("country", country))
	crs, err := c.Ctx.Connection.NewQuery().From(NewBusinessMetricsDataModel().TableName()).Where(dbox.And(filter...)).
		Group("period", "rag", "bmid", "businessmetric").Order("-_id.period").Cursor(nil)
	defer crs.Close()
	if !tk.IsNilOrEmpty(err) {
		return err, result
	}

	o := []tk.M{}
	err = crs.Fetch(&o, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return err, result
	}

	for _, v := range o {
		tom := v.Get("_id").(tk.M)
		d := DataRag{}
		d.BmId = tom.GetString("bmid")
		d.BusinessMetric = tom.GetString("businessmetric")
		d.Period = tom.Get("period").(time.Time)
		d.Rag = tom.GetString("rag")
		ds = append(ds, d)
		// tk.Println(tk.JsonString(tom))
	}
	sort.Sort(sort.Reverse(DataRags(ds)))

	i := 0
	for _, v := range ds {
		if !tk.IsNilOrEmpty(v.Rag) {
			result = append(result, v)
			if i == 2 {
				break
			}
			i++
		}
	}

	return nil, result
}

func (c *ScorecardInitiativeController) findQuarters(bs BusinessMetric, region, country string) (error, []Quarter) {
	quarters := []Quarter{}
	majorRegion := tk.M{}
	if region == country {
		majorRegion.Set("majorregion", region)
	}
	pipe := []tk.M{
		tk.M{
			"$match": tk.M{
				"$and": []tk.M{
					tk.M{
						"bmid": bs.Id,
					},
					majorRegion,
					tk.M{
						"country": country,
					},
				},
			},
		},
		tk.M{
			"$group": tk.M{
				"_id": tk.M{
					"period":         "$period",
					"rag":            "$rag",
					"bmid":           "$bmid",
					"businessmetric": "$businessmetric",
				},
			},
		},
		tk.M{
			"$sort": tk.M{
				"_id.period": 1,
			},
		},
	}

	crs, err := c.Ctx.Connection.NewQuery().Command("pipe", pipe).From(NewBusinessMetricsDataModel().TableName()).Cursor(nil)
	if !tk.IsNilOrEmpty(err) {
		return err, nil
	}
	hasil := tk.Ms{}
	err = crs.Fetch(&hasil, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return err, nil
	}
	crs.Close()

	// defined querter1
	currentyear, _, _ := time.Now().Date()
	quarter1 := []time.Time{
		time.Date(currentyear, time.Month(1), 1, 0, 0, 0, 0, time.Time{}.Location()),
		time.Date(currentyear, time.Month(2), 1, 0, 0, 0, 0, time.Time{}.Location()),
		time.Date(currentyear, time.Month(3), 1, 0, 0, 0, 0, time.Time{}.Location()),
	}
	quarter2 := []time.Time{
		time.Date(currentyear, time.Month(4), 1, 0, 0, 0, 0, time.Time{}.Location()),
		time.Date(currentyear, time.Month(5), 1, 0, 0, 0, 0, time.Time{}.Location()),
		time.Date(currentyear, time.Month(6), 1, 0, 0, 0, 0, time.Time{}.Location()),
	}
	quarter3 := []time.Time{
		time.Date(currentyear, time.Month(7), 1, 0, 0, 0, 0, time.Time{}.Location()),
		time.Date(currentyear, time.Month(8), 1, 0, 0, 0, 0, time.Time{}.Location()),
		time.Date(currentyear, time.Month(9), 1, 0, 0, 0, 0, time.Time{}.Location()),
	}
	quarter4 := []time.Time{
		time.Date(currentyear, time.Month(10), 1, 0, 0, 0, 0, time.Time{}.Location()),
		time.Date(currentyear, time.Month(11), 1, 0, 0, 0, 0, time.Time{}.Location()),
		time.Date(currentyear, time.Month(12), 1, 0, 0, 0, 0, time.Time{}.Location()),
	}

	mapquarters := map[string][]DataRag{}
	for _, v := range hasil {
		d := DataRag{}
		tom := v.Get("_id").(tk.M)
		// m := int(tom.Get("period").(time.Time).Month())
		for i, t := range quarter1 {
			if t.UTC() == tom.Get("period").(time.Time).UTC() {
				d.BmId = tom.GetString("bmid")
				d.BusinessMetric = tom.GetString("businessmetric")
				d.Period = tom.Get("period").(time.Time)
				d.Rag = tom.GetString("rag")
				mapquarters["1"] = append(mapquarters["1"], d)

				quarter1 = append(quarter1[:i], quarter1[i+1:]...)
			}
		}

		for i, t := range quarter2 {
			if t.UTC() == tom.Get("period").(time.Time).UTC() {
				d.BmId = tom.GetString("bmid")
				d.BusinessMetric = tom.GetString("businessmetric")
				d.Period = tom.Get("period").(time.Time)
				d.Rag = tom.GetString("rag")
				mapquarters["2"] = append(mapquarters["2"], d)

				quarter2 = append(quarter2[:i], quarter2[i+1:]...)
			}
		}

		for i, t := range quarter3 {
			if t.UTC() == tom.Get("period").(time.Time).UTC() {
				d.BmId = tom.GetString("bmid")
				d.BusinessMetric = tom.GetString("businessmetric")
				d.Period = tom.Get("period").(time.Time)
				d.Rag = tom.GetString("rag")
				mapquarters["3"] = append(mapquarters["3"], d)

				quarter3 = append(quarter3[:i], quarter3[i+1:]...)
			}
		}

		for i, t := range quarter4 {
			if t.UTC() == tom.Get("period").(time.Time).UTC() {
				d.BmId = tom.GetString("bmid")
				d.BusinessMetric = tom.GetString("businessmetric")
				d.Period = tom.Get("period").(time.Time)
				d.Rag = tom.GetString("rag")
				mapquarters["4"] = append(mapquarters["4"], d)

				quarter4 = append(quarter4[:i], quarter4[i+1:]...)
			}
		}

		// tk.Println(tk.JsonString(mapquarters))
		/*m := int(tom.Get("period").(time.Time).Month())
		if m >= 1 && m <= 3 {
			d.BmId = tom.GetString("bmid")
			d.BusinessMetric = tom.GetString("businessmetric")
			d.Period = tom.Get("period").(time.Time)
			d.Rag = tom.GetString("rag")
			mapquarters["1"] = append(mapquarters["1"], d)

			if exist, i := tk.MemberIndex(monthQuarter, "1"); exist {
				monthQuarter = append(monthQuarter[:i], monthQuarter[i+1:]...)
			}
		} else if m >= 4 && m <= 6 {
			d.BmId = tom.GetString("bmid")
			d.BusinessMetric = tom.GetString("businessmetric")
			d.Period = tom.Get("period").(time.Time)
			d.Rag = tom.GetString("rag")
			mapquarters["2"] = append(mapquarters["2"], d)

			if exist, i := tk.MemberIndex(monthQuarter, "2"); exist {
				monthQuarter = append(monthQuarter[:i], monthQuarter[i+1:]...)
			}
		} else if m >= 7 && m <= 9 {
			d.BmId = tom.GetString("bmid")
			d.BusinessMetric = tom.GetString("businessmetric")
			d.Period = tom.Get("period").(time.Time)
			d.Rag = tom.GetString("rag")
			mapquarters["3"] = append(mapquarters["3"], d)

			if exist, i := tk.MemberIndex(monthQuarter, "3"); exist {
				monthQuarter = append(monthQuarter[:i], monthQuarter[i+1:]...)
			}
		} else if m >= 10 && m <= 12 {
			d.BmId = tom.GetString("bmid")
			d.BusinessMetric = tom.GetString("businessmetric")
			d.Period = tom.Get("period").(time.Time)
			d.Rag = tom.GetString("rag")
			mapquarters["4"] = append(mapquarters["4"], d)

			if exist, i := tk.MemberIndex(monthQuarter, "4"); exist {
				monthQuarter = append(monthQuarter[:i], monthQuarter[i+1:]...)
			}
		}*/
	}

	if len(quarter1) > 0 {
		for _, t := range quarter1 {
			d := DataRag{}
			d.BmId = bs.Id
			d.BusinessMetric = bs.Description
			d.Period = t
			d.Rag = ""
			mapquarters["1"] = append(mapquarters["1"], d)
		}
		sort.Sort(DataRags(mapquarters["1"]))
	}
	if len(quarter2) > 0 {
		for _, t := range quarter2 {
			d := DataRag{}
			d.BmId = bs.Id
			d.BusinessMetric = bs.Description
			d.Period = t
			d.Rag = ""
			mapquarters["2"] = append(mapquarters["2"], d)
		}
		sort.Sort(DataRags(mapquarters["2"]))
	}
	if len(quarter3) > 0 {
		for _, t := range quarter3 {
			d := DataRag{}
			d.BmId = bs.Id
			d.BusinessMetric = bs.Description
			d.Period = t
			d.Rag = ""
			mapquarters["3"] = append(mapquarters["3"], d)
		}
		sort.Sort(DataRags(mapquarters["3"]))
	}
	if len(quarter4) > 0 {
		for _, t := range quarter4 {
			d := DataRag{}
			d.BmId = bs.Id
			d.BusinessMetric = bs.Description
			d.Period = t
			d.Rag = ""
			mapquarters["4"] = append(mapquarters["4"], d)
		}
		sort.Sort(DataRags(mapquarters["4"]))
	}

	for key, v := range mapquarters {
		q := Quarter{}
		if key == "1" {
			q.QuarterName = key
			q.Quarters = v
		} else if key == "2" {
			q.QuarterName = key
			q.Quarters = v
		} else if key == "3" {
			q.QuarterName = key
			q.Quarters = v
		} else if key == "4" {
			q.QuarterName = key
			q.Quarters = v
		}

		quarters = append(quarters, q)
	}

	sort.Sort(Quarters(quarters))

	return nil, quarters
}
