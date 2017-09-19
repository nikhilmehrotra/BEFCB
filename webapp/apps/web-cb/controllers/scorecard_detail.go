package controllers

import (
	m "eaciit/scb-apps/webapp/apps/web-cb/models"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// . "github.com/testcase/webapp/helper"
	"time"
)

type ScorecardDetailController struct {
	*BaseController
}

func (c *ScorecardDetailController) GetData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	parm := struct {
		Region  string
		Country string
		SCId    string
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
	// tk.Println("REGION : ", parm.Region)
	// tk.Println("Country : ", parm.Country)

	result := []tk.M{}
	// Get Result..

	csr, err := c.Ctx.Connection.NewQuery().From("ScorecardDetailCategory").Where(dbox.Eq("scid", parm.SCId)).Order("seq").Cursor(nil)
	err = csr.Fetch(&result, 0, false)
	csr.Close()
	if err != nil {
		return nil
	}

	for _, i := range result {

		sdetail := []tk.M{}
		query := []*dbox.Filter{}
		query = append(query, dbox.Eq("scdetailcategoryid", i.GetString("sdcid")))

		if parm.Region == parm.Country {
			query = append(query, dbox.Eq("majorregion", parm.Region))
		}
		query = append(query, dbox.Eq("country", parm.Country))
		query = append(query, dbox.Eq("scid", parm.SCId))

		csr, err := c.Ctx.Connection.NewQuery().From("ScorecardDetail").Where(dbox.And(query...)).Cursor(nil)
		err = csr.Fetch(&sdetail, 0, false)
		csr.Close()
		if err != nil {
			return nil
		}

		for _, y := range sdetail {
			ytd := 0
			naytd := true
			ragactual := ""

			if y.Get("nadec").(bool) != naytd {
				ytd = y.GetInt("dec")
				ragactual = y.GetString("ragdec")

			} else if y.Get("nanov").(bool) != naytd {
				ytd = y.GetInt("nov")
				ragactual = y.GetString("ragnov")

			} else if y.Get("naoct").(bool) != naytd {
				ytd = y.GetInt("oct")
				ragactual = y.GetString("ragoct")

			} else if y.Get("nasep").(bool) != naytd {
				ytd = y.GetInt("sep")
				ragactual = y.GetString("ragsep")

			} else if y.Get("naaug").(bool) != naytd {
				ytd = y.GetInt("aug")
				ragactual = y.GetString("ragaug")

			} else if y.Get("najul").(bool) != naytd {
				ragactual = y.GetString("ragjul")
				ytd = y.GetInt("jul")

			} else if y.Get("najun").(bool) != naytd {
				ytd = y.GetInt("jun")
				ragactual = y.GetString("ragjun")

			} else if y.Get("namay").(bool) != naytd {
				ytd = y.GetInt("may")
				ragactual = y.GetString("ragmay")

			} else if y.Get("naapr").(bool) != naytd {
				ytd = y.GetInt("apr")
				ragactual = y.GetString("ragapr")

			} else if y.Get("namar").(bool) != naytd {
				ytd = y.GetInt("mar")
				ragactual = y.GetString("ragmar")

			} else if y.Get("nafeb").(bool) != naytd {
				ytd = y.GetInt("feb")
				ragactual = y.GetString("ragfeb")

			} else if y.Get("najan").(bool) != naytd {
				ytd = y.GetInt("jan")
				ragactual = y.GetString("ragjan")

			}

			y.Set("ytd", ytd)
			y.Set("rag", ragactual)

		}

		i.Set("ListMetric", sdetail)

	}

	return c.SetResultInfo(false, "", result)

}

func (c *ScorecardDetailController) GetDataSync(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	parm := struct {
		Region   string
		Country  string
		SCId     string
		SCDId    string
		BMId     string
		DetailID string
		Year     int
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

	result := tk.M{}
	BMData := []m.BusinessMetricsDataModel{}
	// Get Result..
	query := []*dbox.Filter{}
	query = append(query, dbox.Eq("bmid", parm.BMId))
	// query = append(query, dbox.Eq("scid", parm.SCId))
	if parm.Region == parm.Country {
		query = append(query, dbox.Eq("majorregion", parm.Region))
	}
	query = append(query, dbox.Eq("country", parm.Country))
	query = append(query, dbox.Eq("year", parm.Year))

	csr, err := c.Ctx.Connection.NewQuery().From("BusinessMetricsData").Order("period").Where(dbox.And(query...)).Cursor(nil)
	err = csr.Fetch(&BMData, 0, false)
	csr.Close()
	if err != nil {
		return nil
	}

	// tk.Println(len(BMData))

	for _, bm := range BMData {
		period_data := bm.Period
		period := period_data.Month().String()
		// tk.Println(period)

		switch period {
		case "January":
			if bm.ActualYTD != 130895111188 && bm.NAActual != true {
				result.Set("jan", bm.ActualYTD)
				result.Set("najan", false)
			}
			if bm.Baseline != 130895111188.0 && bm.NABaseline != true {
				result.Set("baseline", bm.Baseline)
				result.Set("nabaseline", bm.NABaseline)
			}
			if bm.Target != 130895111188.0 && bm.NATarget != true {
				result.Set("target", bm.Target)
				result.Set("natarget", bm.NATarget)
			}
			if bm.RAG != "" {
				result.Set("ragjan", bm.RAG)
				result.Set("naragjan", false)
			}
		case "February":
			if bm.ActualYTD != 130895111188 && bm.NAActual != true {
				result.Set("feb", bm.ActualYTD)
				result.Set("nafeb", false)
			}
			if bm.Baseline != 130895111188.0 && bm.NABaseline != true {
				result.Set("baseline", bm.Baseline)
				result.Set("nabaseline", bm.NABaseline)
			}
			if bm.Target != 130895111188.0 && bm.NATarget != true {
				result.Set("target", bm.Target)
				result.Set("natarget", bm.NATarget)
			}
			if bm.RAG != "" {
				result.Set("ragfeb", bm.RAG)
				result.Set("naragfeb", false)
			}
		case "March":
			if bm.ActualYTD != 130895111188 && bm.NAActual != true {
				result.Set("mar", bm.ActualYTD)
				result.Set("namar", false)
			}
			if bm.Baseline != 130895111188.0 && bm.NABaseline != true {
				result.Set("baseline", bm.Baseline)
				result.Set("nabaseline", bm.NABaseline)
			}
			if bm.Target != 130895111188.0 && bm.NATarget != true {
				result.Set("target", bm.Target)
				result.Set("natarget", bm.NATarget)
			}
			if bm.RAG != "" {
				result.Set("ragmar", bm.RAG)
				result.Set("naragmar", false)
			}
		case "April":
			if bm.ActualYTD != 130895111188 && bm.NAActual != true {
				result.Set("apr", bm.ActualYTD)
				result.Set("naapr", false)
			}
			if bm.Baseline != 130895111188.0 && bm.NABaseline != true {
				result.Set("baseline", bm.Baseline)
				result.Set("nabaseline", bm.NABaseline)
			}
			if bm.Target != 130895111188.0 && bm.NATarget != true {
				result.Set("target", bm.Target)
				result.Set("natarget", bm.NATarget)
			}
			if bm.RAG != "" {
				result.Set("ragapr", bm.RAG)
				result.Set("naragapr", false)
			}
		case "May":
			if bm.ActualYTD != 130895111188 && bm.NAActual != true {
				result.Set("may", bm.ActualYTD)
				result.Set("namay", false)
			}
			if bm.Baseline != 130895111188.0 && bm.NABaseline != true {
				result.Set("baseline", bm.Baseline)
				result.Set("nabaseline", bm.NABaseline)
			}
			if bm.Target != 130895111188.0 && bm.NATarget != true {
				result.Set("target", bm.Target)
				result.Set("natarget", bm.NATarget)
			}
			if bm.RAG != "" {
				result.Set("ragmay", bm.RAG)
				result.Set("naragmay", false)
			}
		case "June":
			if bm.ActualYTD != 130895111188 && bm.NAActual != true {
				result.Set("jun", bm.ActualYTD)
				result.Set("najun", false)
			}
			if bm.Baseline != 130895111188.0 && bm.NABaseline != true {
				result.Set("baseline", bm.Baseline)
				result.Set("nabaseline", bm.NABaseline)
			}
			if bm.Target != 130895111188.0 && bm.NATarget != true {
				result.Set("target", bm.Target)
				result.Set("natarget", bm.NATarget)
			}
			if bm.RAG != "" {
				result.Set("ragjun", bm.RAG)
				result.Set("naragjun", false)
			}
		case "July":
			if bm.ActualYTD != 130895111188 && bm.NAActual != true {
				result.Set("jul", bm.ActualYTD)
				result.Set("najul", false)
			}
			if bm.Baseline != 130895111188.0 && bm.NABaseline != true {
				result.Set("baseline", bm.Baseline)
				result.Set("nabaseline", bm.NABaseline)
			}
			if bm.Target != 130895111188.0 && bm.NATarget != true {
				result.Set("target", bm.Target)
				result.Set("natarget", bm.NATarget)
			}
			if bm.RAG != "" {
				result.Set("ragjul", bm.RAG)
				result.Set("naragjul", false)
			}
		case "August":
			if bm.ActualYTD != 130895111188 && bm.NAActual != true {
				result.Set("aug", bm.ActualYTD)
				result.Set("naaug", false)
			}
			if bm.Baseline != 130895111188.0 && bm.NABaseline != true {
				result.Set("baseline", bm.Baseline)
				result.Set("nabaseline", bm.NABaseline)
			}
			if bm.Target != 130895111188.0 && bm.NATarget != true {
				result.Set("target", bm.Target)
				result.Set("natarget", bm.NATarget)
			}
			if bm.RAG != "" {
				result.Set("ragaug", bm.RAG)
				result.Set("naragaug", false)
			}
		case "September":
			if bm.ActualYTD != 130895111188 && bm.NAActual != true {
				result.Set("sep", bm.ActualYTD)
				result.Set("nasep", false)
			}
			if bm.Baseline != 130895111188.0 && bm.NABaseline != true {
				result.Set("baseline", bm.Baseline)
				result.Set("nabaseline", bm.NABaseline)
			}
			if bm.Target != 130895111188.0 && bm.NATarget != true {
				result.Set("target", bm.Target)
				result.Set("natarget", bm.NATarget)
			}
			if bm.RAG != "" {
				result.Set("ragsep", bm.RAG)
				result.Set("naragsep", false)
			}
		case "October":
			if bm.ActualYTD != 130895111188 && bm.NAActual != true {
				result.Set("oct", bm.ActualYTD)
				result.Set("naoct", false)
			}
			if bm.Baseline != 130895111188.0 && bm.NABaseline != true {
				result.Set("baseline", bm.Baseline)
				result.Set("nabaseline", bm.NABaseline)
			}
			if bm.Target != 130895111188.0 && bm.NATarget != true {
				result.Set("target", bm.Target)
				result.Set("natarget", bm.NATarget)
			}
			if bm.RAG != "" {
				result.Set("ragoct", bm.RAG)
				result.Set("naragoct", false)
			}
		case "November":
			if bm.ActualYTD != 130895111188 && bm.NAActual != true {
				result.Set("nov", bm.ActualYTD)
				result.Set("nanov", false)
			}
			if bm.Baseline != 130895111188.0 && bm.NABaseline != true {
				result.Set("baseline", bm.Baseline)
				result.Set("nabaseline", bm.NABaseline)
			}
			if bm.Target != 130895111188.0 && bm.NATarget != true {
				result.Set("target", bm.Target)
				result.Set("natarget", bm.NATarget)
			}
			if bm.RAG != "" {
				result.Set("ragnov", bm.RAG)
				result.Set("naragnov", false)
			}

		case "December":
			if bm.ActualYTD != 130895111188 && bm.NAActual != true {
				result.Set("dec", bm.ActualYTD)
				result.Set("nadec", false)
			}
			if bm.Baseline != 130895111188.0 && bm.NABaseline != true {
				result.Set("baseline", bm.Baseline)
				result.Set("nabaseline", bm.NABaseline)
			}
			if bm.Target != 130895111188.0 && bm.NATarget != true {
				result.Set("target", bm.Target)
				result.Set("natarget", bm.NATarget)
			}
			if bm.RAG != "" {
				result.Set("ragdec", bm.RAG)
				result.Set("naragdec", false)
			}
			break
			// default:
			// 	result.Set("jan", 0)
			// 	result.Set("baseline", 0)
			// 	result.Set("target", 0)
			// 	break

		}

	}
	result.Set("DetailID", parm.DetailID)
	result.Set("BMId", parm.BMId)
	return c.SetResultInfo(false, "", result)

}

func (c *ScorecardDetailController) Save(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Region   string
		Country  string
		SCId     string
		DataList []m.ScorecardDetailModel
	}{}

	err := k.GetPayload(&parm)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	for _, i := range parm.DataList {
		tempData := m.ScorecardDetailModel{}
		existingName := ""
		csr, err := c.Ctx.Connection.NewQuery().From(new(m.ScorecardDetailModel).TableName()).Where(dbox.Eq("_id", i.Id)).Cursor(nil)
		err = csr.Fetch(&tempData, 1, false)
		csr.Close()
		if err != nil {
			existingName = ""
		} else {
			existingName = tempData.Name
		}

		i.Updated_Date = time.Now()
		i.Updated_By = k.Session("username").(string)
		err = c.Ctx.Save(&i)
		if err != nil {
			return c.SetResultInfo(true, err.Error(), nil)
		}

		existingData := []m.ScorecardDetailModel{}
		csr, err = c.Ctx.Connection.NewQuery().From(new(m.ScorecardDetailModel).TableName()).Where(dbox.And(dbox.Eq("scid", i.SCId), dbox.Eq("scdetailcategoryid", i.SCDetailCategoryId), dbox.Eq("name", existingName))).Cursor(nil)
		err = csr.Fetch(&existingData, 0, false)
		csr.Close()
		if err == nil {
			// tk.Println(len(existingData))
			for _, x := range existingData {
				x.Name = i.Name
				x.Description = i.Description
				x.Type = i.Type
				x.Denomination = i.Denomination
				x.ValueType = i.ValueType
				x.DecimalFormat = i.DecimalFormat
				x.MetricReference = i.MetricReference
				err = c.Ctx.Save(&x)
				if err != nil {
					return c.SetResultInfo(true, err.Error(), nil)
				}
			}
		}

	}
	return c.SetResultInfo(false, "", nil)
}
