package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"time"
)

func (m *MController) SharedAgenda(k *knot.WebContext) interface{} {
	m.LoadBase(k)
	k.Config.NoLog = true
	k.Config.LayoutTemplate = "_layout_dedicated.html"
	k.Config.IncludeFiles = []string{"shared/sidebar.html"}
	k.Config.OutputType = knot.OutputTemplate
	return nil
}

func (m *MController) SharedAgendaGetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	///get summery business driver
	crs, err := m.Ctx.Connection.NewQuery().From(NewSummaryBusinessDriverModel().TableName()).Cursor(nil)
	defer crs.Close()
	if !tk.IsNilOrEmpty(err) {
		return m.ErrorResultInfo(err.Error(), nil)
	}

	bdresult := []SummaryBusinessDriverModel{}
	err = crs.Fetch(&bdresult, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return m.ErrorResultInfo(err.Error(), nil)
	}

	bdData := map[string]string{}
	for _, v := range bdresult {
		if _, exist := bdData[v.Parentid+"|"+v.Idx]; !exist {
			bdData[v.Parentid+"|"+v.Idx] = v.Name
		}
	}

	crssa, err := m.Ctx.Connection.NewQuery().From(NewSharedAgendaModel().TableName()).Order("name", "seq").Cursor(nil) //m.Ctx.Find(NewSharedAgendaModel(), tk.M{})
	defer crssa.Close()
	if !tk.IsNilOrEmpty(err) {
		return m.ErrorResultInfo(err.Error(), nil)
	}

	saresult := []SharedAgendaModel{}
	err = crssa.Fetch(&saresult, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return m.ErrorResultInfo(err.Error(), nil)
	}

	result := []SharedAgendaModel{} //tk.Ms{}
	for _, v := range saresult {
		v.BusinessDriverName = bdData[v.SCId+"|"+v.BDId]
		result = append(result, v)
	}
	res := ResultInfo{}
	res.Data = result
	res.Total = len(result)
	return res
}

type ScoreCard struct {
	Text  string
	Value string
}
type SortScoreCard []ScoreCard

func (p SortScoreCard) Len() int           { return len(p) }
func (p SortScoreCard) Less(i, j int) bool { return p[i].Text < p[j].Text }
func (p SortScoreCard) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (c *MController) ScorecardList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	crs, err := c.Ctx.Connection.NewQuery().From(NewSummaryBusinessDriverModel().TableName()).
		Where(dbox.Ne("parentname", "Financial Framework ")).Group("parentid", "parentname").Order("parentid").Cursor(nil)
	defer crs.Close()
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	data := tk.Ms{}
	err = crs.Fetch(&data, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	result := []ScoreCard{}
	for _, v := range data {
		tom := v.Get("_id").(tk.M)
		o := ScoreCard{}
		o.Text = tom.GetString("parentid")
		o.Value = tom.GetString("parentname")
		result = append(result, o)
	}
	sort.Sort(SortScoreCard(result))

	return c.SetResultInfo(false, "", result)
}

func (c *MController) BusinessDriverList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		ParentId string
	}{}
	err := k.GetPayload(&parm)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	crs, err := c.Ctx.Find(NewSummaryBusinessDriverModel(), tk.M{}.Set("where", dbox.Eq("parentid", parm.ParentId)))
	defer crs.Close()
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	result := []SummaryBusinessDriverModel{}
	err = crs.Fetch(&result, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	return c.SetResultInfo(false, "", result)
}

func (c *MController) SharedAgendaSave(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	sa := SharedAgendaModel{}
	err := k.GetPayload(&sa)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	if sa.Id == "" {
		sa.Id = bson.NewObjectId()
		sa.CreatedDate = time.Now()
		sa.CreatedBy = k.Session("username").(string)
	}
	sa.UpdatedBy = k.Session("username").(string)
	sa.UpdatedDate = time.Now()

	//get scorecardname
	crs, err := c.Ctx.Connection.NewQuery().From(NewSummaryBusinessDriverModel().TableName()).
		Where(dbox.And(dbox.Eq("parentid", sa.SCId), dbox.Eq("Id", sa.BDId))).Group("Name", "parentname").Cursor(nil)
	defer crs.Close()
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	data := tk.Ms{}
	err = crs.Fetch(&data, 1, false)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	sa.SCName = data[0].Get("_id").(map[string]interface{})["parentname"].(string)
	sa.BusinessDriverName = data[0].Get("_id").(map[string]interface{})["Name"].(string)

	err = c.Ctx.Save(&sa)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	return c.SetResultInfo(false, "Save Success", nil)
}

func (c *MController) SharedAgendaDelete(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	param := struct {
		Id string
	}{}
	err := k.GetPayload(&param)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	m := NewSharedAgendaModel()
	err = c.Ctx.GetById(m, bson.ObjectIdHex(param.Id))
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	m.IsDeleted = true
	err = c.Ctx.Save(m)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	return c.SetResultInfo(false, "Delete Success", nil)
}

func (c *MController) GetBy(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	param := struct {
		Id string
	}{}
	err := k.GetPayload(&param)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	m := NewSharedAgendaModel()
	err = c.Ctx.GetById(m, bson.ObjectIdHex(param.Id))
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	res := ResultInfo{}
	res.Data = m
	return res
}
