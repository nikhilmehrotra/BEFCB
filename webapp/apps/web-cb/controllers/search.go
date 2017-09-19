package controllers

import (
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"sort"
	"strings"

	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type SearchController struct {
	*BaseController
}
type KeywordStruct struct {
	Id         string
	Keyword    string
	ParentInfo string
	Score      int
}
type ListKey []KeywordStruct

func (l ListKey) Len() int { return len(l) }
func (l ListKey) Less(i, j int) bool {
	//	if l[i].Score < l[j].Score {
	//		return true
	//	}
	//	if l[i].Score == l[j].Score && l[i].Keyword < l[j].Keyword {
	//		return true
	//	} else {
	//		return false
	//	}
	//	return l[i].Score < l[j].Score
	return l[i].Keyword < l[j].Keyword
}
func (l ListKey) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

func (c *BaseController) AutoComplete(k *knot.WebContext) interface{} {
	//	c.LoadBase(k)
	InitiativeAccess := c.GetAccess(k, "INITIATIVE")
	Initiative := new(AccessibilityModel)
	e := tk.MtoStruct(InitiativeAccess, &Initiative)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	k.Config.OutputType = knot.OutputJson
	frm := struct {
		Keyword        string
		InitiativeType string
		Initiatives    []string
	}{}
	err := k.GetPayload(&frm)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	if frm.InitiativeType != "SupportingEnablers" {
		frm.InitiativeType = "KeyEnablers"
	}
	frm.Keyword = strings.ToLower(frm.Keyword)
	ret := ResultInfo{}

	InitiativeIDs := []interface{}{}
	for _, x := range frm.Initiatives {
		InitiativeIDs = append(InitiativeIDs, x)
	}
	query := []*db.Filter{}
	query = append(query, db.Eq("InitiativeType", frm.InitiativeType))
	query = append(query, db.In("InitiativeID", InitiativeIDs...))
	query = append(query, db.Or(db.Contains("ProjectName", frm.Keyword), db.Contains("ProjectManager", frm.Keyword),
		db.Contains("AccountableExecutive", frm.Keyword), db.Contains("TechnologyLead", frm.Keyword)))

	UserCountry := ""
	if k.Session("country") != nil {
		UserCountry = k.Session("country").(string)
	}

	if UserCountry != "" && (Initiative.Global.Read || Initiative.Region.Read || Initiative.Country.Read) == true {
		query = append(query, db.Or(db.Eq("Country", UserCountry), db.Eq("IsGlobal", true)))
	} else if UserCountry != "" && (Initiative.Global.Read || Initiative.Region.Read || Initiative.Country.Read) == false {
		query = append(query, db.Or(db.Eq("Country", UserCountry), db.Eq("IsGlobal", false)))
	}
	qryClause := tk.M{}.Set("select",
		"_id,InitiativeType,ProjectName,ProjectManager,AccountableExecutive,TechnologyLead").
		Set("where", db.And(query...))

	crs, err := c.Ctx.Find(NewInitiativeModel(), qryClause)
	defer crs.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	data := make([]InitiativeModel, 0)
	err = crs.Fetch(&data, 0, false)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	dataMap := make(map[string]KeywordStruct, 0)
	for _, dt := range data {
		//		tk.Println("Proj.Name ", dt.ProjectName, "; PM ", dt.ProjectManager, "; AM ", dt.AccountableExecutive, "; TL ", dt.TechnologyLead)
		score, val := c.getScore(dt.ProjectName, frm.Keyword)
		if score > 0 {
			dataMap[val] = KeywordStruct{Keyword: val, Score: score}
		}
		score, val = c.getScore(strings.Join(dt.ProjectManager, ","), frm.Keyword)
		if score > 0 {
			dataMap[val] = KeywordStruct{Keyword: val, Score: score}
		}
		score, val = c.getScore(strings.Join(dt.AccountableExecutive, ","), frm.Keyword)
		if score > 0 {
			dataMap[val] = KeywordStruct{Keyword: val, Score: score}
		}
		score, val = c.getScore(strings.Join(dt.TechnologyLead, ","), frm.Keyword)
		if score > 0 {
			dataMap[val] = KeywordStruct{Keyword: val, Score: score}
		}
	}
	qryTaskClause := tk.M{}.Set("select", "_id,Name,Owner").
		Set("where", db.Or(db.Contains("name", frm.Keyword), db.Contains("owner", frm.Keyword)))
	crsTask, err := c.Ctx.Find(NewTaskModel(), qryTaskClause)
	defer crsTask.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	dataTask := make([]TaskModel, 0)
	err = crsTask.Fetch(&dataTask, 0, false)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	for _, dt := range dataTask {
		score, val := c.getScore(dt.Name, frm.Keyword)
		if score > 0 {
			dataMap[val] = KeywordStruct{Keyword: val, Score: score}
		}
		score, val = c.getScore(dt.Owner, frm.Keyword)
		if score > 0 {
			dataMap[val] = KeywordStruct{Keyword: val, Score: score}
		}
	}
	dataList := c.rankAutoComplete(dataMap)
	//	for _, dt := range dataList {
	//		tk.Printfn("DATALIST %#v ", dt)
	//	}
	if len(dataList) > 50 {
		ret.Data = dataList[:49]
	} else {
		ret.Data = dataList
	}
	return ret
}

func (c *BaseController) GetResult(k *knot.WebContext) interface{} {
	//	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	frm := struct {
		Keyword        string
		KeywordEscaped string
		InitiativeType string
		Initiatives    []string
	}{}
	err := k.GetPayload(&frm)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	if frm.InitiativeType != "SupportingEnablers" {
		frm.InitiativeType = "KeyEnablers"
	}
	frm.Keyword = strings.ToLower(frm.Keyword)
	frm.KeywordEscaped = strings.ToLower(frm.KeywordEscaped)
	InitiativeIDs := []interface{}{}
	for _, x := range frm.Initiatives {
		InitiativeIDs = append(InitiativeIDs, x)
	}
	ret := ResultInfo{}
	qryClause := tk.M{}.Set("where", db.And(db.Eq("InitiativeType", frm.InitiativeType), db.In("InitiativeID", InitiativeIDs...),
		db.Or(
			db.Contains("ProjectName", frm.Keyword), db.Contains("ProjectManager", frm.Keyword),
			db.Contains("AccountableExecutive", frm.Keyword), db.Contains("TechnologyLead", frm.Keyword),
			db.Contains("ProblemStatement", frm.Keyword), db.Contains("ProjectDescription", frm.Keyword),
			db.Contains("Attachments.filename", frm.Keyword),
		)))
	crs, err := c.Ctx.Find(NewInitiativeModel(), qryClause)
	defer crs.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	data := make([]InitiativeModel, 0)
	err = crs.Fetch(&data, 0, false)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	tk.Println("-----", frm.KeywordEscaped, data)

	dataMapPN := make(map[string]KeywordStruct, 0)
	dataMapPM := make(map[string]KeywordStruct, 0)
	dataMapAM := make(map[string]KeywordStruct, 0)
	dataMapTL := make(map[string]KeywordStruct, 0)
	dataMapPS := make(map[string]KeywordStruct, 0)
	dataMapPD := make(map[string]KeywordStruct, 0)
	dataMapFN := make(map[string]KeywordStruct, 0)
	for _, dt := range data {
		//		tk.Println("Proj.Name ", dt.ProjectName, "; PM ", dt.ProjectManager, "; AM ", dt.AccountableExecutive, "; TL ", dt.TechnologyLead)
		//		tk.Println("Proj.Name ", dt.ProjectName, "; PM ", dt.ProjectDescription)
		score, val := c.getScore(dt.ProjectName, frm.Keyword)
		if score > 0 {
			dataMapPN[dt.Id.Hex()] = KeywordStruct{Id: dt.Id.Hex(), Keyword: val, Score: score}
		}
		score, val = c.getScore(strings.Join(dt.ProjectManager, ","), frm.Keyword)
		if score > 0 {
			dataMapPM[dt.Id.Hex()] = KeywordStruct{Id: dt.Id.Hex(), Keyword: val, ParentInfo: dt.ProjectName, Score: score}
		}
		score, val = c.getScore(strings.Join(dt.AccountableExecutive, ","), frm.Keyword)
		if score > 0 {
			dataMapAM[dt.Id.Hex()] = KeywordStruct{Id: dt.Id.Hex(), Keyword: val, ParentInfo: dt.ProjectName, Score: score}
		}
		score, val = c.getScore(strings.Join(dt.TechnologyLead, ","), frm.Keyword)
		if score > 0 {
			dataMapTL[dt.Id.Hex()] = KeywordStruct{Id: dt.Id.Hex(), Keyword: val, ParentInfo: dt.ProjectName, Score: score}
		}
		score, val = c.getScore(dt.ProblemStatement, frm.Keyword)
		if score > 0 {
			val = c.getPrettyDescriptionText(dt.ProblemStatement, frm.Keyword)
			dataMapPS[dt.Id.Hex()] = KeywordStruct{Id: dt.Id.Hex(), Keyword: val, ParentInfo: dt.ProjectName, Score: score}
		}
		score, val = c.getScore(dt.ProjectDescription, frm.Keyword)
		if score > 0 {
			val = c.getPrettyDescriptionText(dt.ProjectDescription, frm.Keyword)

			dataMapPD[dt.Id.Hex()] = KeywordStruct{Id: dt.Id.Hex(), Keyword: val, ParentInfo: dt.ProjectName, Score: score}
		}
		for _, atFN := range dt.Attachments {
			score, val = c.getScore(atFN.FileName, frm.Keyword)
			if score > 0 {
				dataMapFN[dt.Id.Hex()] = KeywordStruct{Id: dt.Id.Hex(), Keyword: val, ParentInfo: dt.ProjectName, Score: score}
			}

		}
	}
	dataListPN := c.rankAutoComplete(dataMapPN)
	dataListPM := c.rankAutoComplete(dataMapPM)
	dataListAM := c.rankAutoComplete(dataMapAM)
	dataListTL := c.rankAutoComplete(dataMapTL)
	dataListPS := c.rankAutoComplete(dataMapPS)
	dataListPD := c.rankAutoComplete(dataMapPD)
	dataListFN := c.rankAutoComplete(dataMapFN)

	qryTaskClause := tk.M{}.Set("where", db.Or(db.Contains("name", frm.Keyword), db.Contains("owner", frm.Keyword), db.Contains("statement", frm.Keyword), db.Contains("description", frm.Keyword)))
	crsTask, err := c.Ctx.Find(NewTaskModel(), qryTaskClause)
	defer crsTask.Close()
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	dataTask := make([]TaskModel, 0)
	err = crsTask.Fetch(&dataTask, 0, false)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	dataMapTN := make(map[string]KeywordStruct, 0)
	dataMapTO := make(map[string]KeywordStruct, 0)
	dataMapTS := make(map[string]KeywordStruct, 0)
	dataMapTD := make(map[string]KeywordStruct, 0)
	//	tk.Printfn("%#v", dataTask)
	for _, dt := range dataTask {
		if score, val := c.getScore(dt.Name, frm.Keyword); score > 0 {
			dataMapTN[dt.Id.Hex()] = KeywordStruct{Id: dt.Id.Hex(), Keyword: val, Score: score}
		}
		if score, val := c.getScore(dt.Owner, frm.Keyword); score > 0 {
			dataMapTO[dt.Id.Hex()] = KeywordStruct{Id: dt.Id.Hex(), Keyword: val, ParentInfo: dt.Name, Score: score}
		}
		if score, val := c.getScore(dt.Statement, frm.Keyword); score > 0 {
			val = c.getPrettyDescriptionText(dt.Statement, frm.Keyword)
			dataMapTS[dt.Id.Hex()] = KeywordStruct{Id: dt.Id.Hex(), Keyword: val, ParentInfo: dt.Name, Score: score}
		}
		if score, val := c.getScore(dt.Description, frm.Keyword); score > 0 {
			val = c.getPrettyDescriptionText(dt.Description, frm.Keyword)
			dataMapTD[dt.Id.Hex()] = KeywordStruct{Id: dt.Id.Hex(), Keyword: val, ParentInfo: dt.Name, Score: score}
		}
	}
	dataListTN := c.rankAutoComplete(dataMapTN)
	dataListTO := c.rankAutoComplete(dataMapTO)
	dataListTS := c.rankAutoComplete(dataMapTS)
	dataListTD := c.rankAutoComplete(dataMapTD)

	ret.Data = tk.M{}.
		Set("Keyword", frm.Keyword).
		Set("ProjectName", dataListPN).
		Set("ProjectManager", dataListPM).
		Set("AccountExecutive", dataListAM).
		Set("TechnologyLead", dataListTL).
		Set("ProblemStatement", dataListPS).
		Set("ProjectDescription", dataListPD).
		Set("Filenames", dataListFN).
		Set("TaskName", dataListTN).
		Set("TaskOwner", dataListTO).
		Set("TaskStatement", dataListTS).
		Set("TaskDesc", dataListTD)

	return ret
}

func (c *BaseController) getScore(val, search string) (int, string) {
	tval := strings.ToLower(val)
	score := 0
	if strings.Contains(tval, search) {
		score += 5
		if strings.HasPrefix(tval, search) {
			score += 3
		}
	}
	return score, val
}

func (c *BaseController) getPrettyDescriptionText(text, keyword string) string {
	strReal := text
	search := keyword
	str := strings.ToLower(text)
	idx := strings.Index(str, search)
	val := ""
	if idx > -1 {
		if idx == 0 {
			if (idx + 50) > len(strReal) {
				val = strReal[idx:]
			} else {
				val = strReal[idx:(idx + 50)]
				lastIndex := strings.LastIndex(val, " ")
				val = val[:lastIndex]
			}
		} else {
			lastIndex := strings.LastIndex(strReal[:idx], " ")
			idx = lastIndex + 1
			if (idx + 50) > len(strReal) {
				val = strReal[idx:]
			} else {
				val = strReal[idx:(idx + 50)]
				lastIndex = strings.LastIndex(val, " ")
				val = val[:lastIndex]
			}
		}
	}
	return val
}

func (c *BaseController) rankAutoComplete(dataMap map[string]KeywordStruct) ListKey {
	ret := make(ListKey, len(dataMap))
	cx := 0
	for key := range dataMap {
		ret[cx] = dataMap[key]
		cx++
	}
	sort.Sort(sort.Reverse(ret))
	return ret
}
