package controllers

import (
	// helper "eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"strings"
	"time"
	// "fmt"
	"github.com/eaciit/acl/v2.0"
	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
)

func (c *AclController) GetUserLoginActivityList(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	resultuser := []tk.M{}

	// Get PayLoad
	type SubFilters struct {
		Field    string
		Operator string
		Value    string
		Type     string
	}
	type Filters struct {
		Field    string
		Operator string
		Value    string
		Type     string
		Filters  []SubFilters
	}
	type Filter struct {
		Filters []Filters
		Logic   string
	}
	type Sort struct {
		Field string
		Dir   string
	}
	parm := struct {
		Skip     int
		Take     int
		Page     int
		PageSize int
		Filter   Filter
		Sort     []Sort

		// Priority Filter
		StartDate  string
		FinishDate string
		LoginID    []string
		FullName   []string
		Country    []string
		UserRoles  []string
		Action     []string
		Modules    []string
	}{}
	e := k.GetPayload(&parm)

	// Get LoginID Data
	d := new(acl.User)
	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.Startwith("groups", "CB_")).Select("loginid").Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&resultuser, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	logins := []interface{}{}
	for _, i := range resultuser {
		login := i.Get("loginid")
		logins = append(logins, login)
	}

	result := []tk.M{}
	query := []*db.Filter{}
	sort := []string{}

	// Apply Top Filter
	if parm.StartDate != "" {
		StartDate, _ := time.Parse("20060102150405", parm.StartDate)
		query = append(query, db.Gte("created", StartDate))
	}
	if parm.FinishDate != "" {
		FinishDate, _ := time.Parse("20060102150405", parm.FinishDate)
		query = append(query, db.Lte("expired", FinishDate.AddDate(0, 0, 1)))
	}

	if len(parm.LoginID) > 0 {
		loginId := []interface{}{}
		for _, x := range parm.LoginID {
			eachLoginId := x
			loginId = append(loginId, eachLoginId)
		}
		query = append(query, db.In("loginid", loginId...))
	}
	if len(parm.FullName) > 0 {
		fullName := []interface{}{}
		for _, x := range parm.FullName {
			eachFullName := x
			fullName = append(fullName, eachFullName)
		}
		query = append(query, db.In("fullname", fullName...))
	}
	if len(parm.Country) > 0 {
		country := []interface{}{}
		for _, x := range parm.Country {
			eachCountry := x
			country = append(country, eachCountry)
		}
		query = append(query, db.In("country", country...))
	}
	if len(parm.UserRoles) > 0 {
		userRoles := []interface{}{}
		for _, x := range parm.UserRoles {
			eachUserRoles := x
			userRoles = append(userRoles, eachUserRoles)
		}
		query = append(query, db.In("group", userRoles...))
	}
	if len(parm.Action) > 0 {
		actions := []interface{}{}
		for _, x := range parm.Action {
			// tk.Println(x)
			actions = append(actions, x)
		}
		LogResults := make([]LogModel, 0)
		csr, e = c.AclCtx.Connection.NewQuery().From(new(LogModel).TableName()).Where(db.In("do", actions...)).Select("sessionid").Cursor(nil)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
		e = csr.Fetch(&LogResults, 0, false)
		csr.Close()
		sessionsid := []interface{}{}
		ex_sessionsid := ""
		for _, x := range LogResults {
			if !strings.Contains(ex_sessionsid, x.SessionID) {
				sessionsid = append(sessionsid, x.SessionID)
				ex_sessionsid += x.SessionID + "|"
			}
		}
		// tk.Println(ex_sessionsid)
		query = append(query, db.In("_id", sessionsid...))
	}
	if len(parm.Modules) > 0 {
		modules := []interface{}{}
		for _, x := range parm.Modules {
			// tk.Println(x)
			modules = append(modules, x)
		}
		LogResults := make([]LogModel, 0)
		csr, e = c.AclCtx.Connection.NewQuery().From(new(LogModel).TableName()).Where(db.In("module", modules...)).Select("sessionid").Cursor(nil)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
		e = csr.Fetch(&LogResults, 0, false)
		csr.Close()
		sessionsid := []interface{}{}
		ex_sessionsid := ""
		for _, x := range LogResults {
			if !strings.Contains(ex_sessionsid, x.SessionID) {
				sessionsid = append(sessionsid, x.SessionID)
				ex_sessionsid += x.SessionID + "|"
			}
		}
		// tk.Println(ex_sessionsid)
		query = append(query, db.In("_id", sessionsid...))
	}

	query = append(query, db.In("loginid", logins...))
	queryFilter := []*db.Filter{}

	for _, f := range parm.Filter.Filters {
		if len(parm.Filter.Filters) > 0 {
			if len(f.Filters) == 0 {
				if f.Type == "date" {
					tmpdate, _ := time.Parse("20060102150405", f.Value)
					tmpdate = tmpdate.UTC()
					value := tmpdate
					switch f.Operator {
					case "eq":
						queryFilter = append(queryFilter, db.Eq(f.Field, value))
					case "gte":
						queryFilter = append(queryFilter, db.Gte(f.Field, value))
					case "lte":
						queryFilter = append(queryFilter, db.Lte(f.Field, value))
					default:
						break
					}
				} else {
					value := f.Value
					switch f.Operator {
					case "contains":
						queryFilter = append(queryFilter, db.Contains(f.Field, value))
					case "eq":
						queryFilter = append(queryFilter, db.Eq(f.Field, value))
					default:
						break
					}

				}
				if len(queryFilter) > 0 {
					if parm.Filter.Logic == "and" {
						query = append(query, db.And(queryFilter...))
					} else {
						query = append(query, db.Or(queryFilter...))
					}
				}
			} else {
				MultipleFilter := []*db.Filter{}
				for _, x := range parm.Filter.Filters {
					for _, y := range x.Filters {
						if y.Type == "date" {
							tmpdate, _ := time.Parse("20060102150405", y.Value)
							tmpdate = tmpdate.UTC()
							value := tmpdate
							switch y.Operator {
							case "eq":
								MultipleFilter = append(MultipleFilter, db.Eq(y.Field, value))
							case "gte":
								MultipleFilter = append(MultipleFilter, db.Gte(y.Field, value))
							case "lte":
								MultipleFilter = append(MultipleFilter, db.Lte(y.Field, value))
							default:
								break
							}
						} else {
							value := y.Value
							switch y.Operator {
							case "contains":
								MultipleFilter = append(MultipleFilter, db.Contains(y.Field, value))
							case "eq":
								MultipleFilter = append(MultipleFilter, db.Eq(y.Field, value))
							default:
								break
							}

						}
					}
				}
				if len(MultipleFilter) > 0 {
					if parm.Filter.Logic == "and" {
						query = append(query, db.And(MultipleFilter...))
					} else {
						query = append(query, db.Or(MultipleFilter...))
					}
				}
			}

		}
	}

	s := new(acl.Session)
	// Get Sort Value
	if len(parm.Sort) > 0 {
		for _, x := range parm.Sort {
			if x.Dir == "asc" {
				sort = append(sort, x.Field)
			} else {
				sort = append(sort, "-"+x.Field)
			}
		}
	} else {
		sort = append(sort, "-created")
	}
	csr, e = c.AclCtx.Connection.NewQuery().From(s.TableName()).Where(db.And(query...)).Skip(parm.Skip).Take(parm.Take).Order(sort...).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	returnvalue := struct {
		Data  []tk.M
		Total int
	}{
		Data:  result,
		Total: csr.Count(),
	}

	return c.SetResultInfo(false, "", returnvalue)
}

func (c *AclController) GetUserLoginActivityReferences(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	result := tk.M{}
	result.Set("LoginActivity", new(acl.Session))

	// Get Country List
	CountryList := []RegionModel{}
	csr, err := c.Ctx.Connection.NewQuery().From(NewRegion().TableName()).Cursor(nil)
	err = csr.Fetch(&CountryList, 0, false)
	csr.Close()
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	result.Set("CountryList", CountryList)

	LoginIDList := []string{}
	FullNameList := []string{}
	UserList := []acl.User{}
	csr, err = c.AclCtx.Connection.NewQuery().From(new(acl.User).TableName()).Cursor(nil)
	err = csr.Fetch(&UserList, 0, false)
	csr.Close()
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	for _, x := range UserList {
		LoginIDList = append(LoginIDList, x.LoginID)
		FullNameList = append(FullNameList, x.FullName)
	}
	result.Set("LoginIDList", LoginIDList)
	result.Set("FullNameList", FullNameList)

	UserRolesList := []acl.Group{}
	csr, err = c.AclCtx.Connection.NewQuery().From(new(acl.Group).TableName()).Cursor(nil)
	err = csr.Fetch(&UserRolesList, 0, false)
	csr.Close()
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	result.Set("UserRolesList", UserRolesList)

	ActionList := []tk.M{}
	csr, err = c.AclCtx.Connection.NewQuery().From("acl_action").Where(db.Eq("app", c.APP_NAME)).Cursor(nil)
	err = csr.Fetch(&ActionList, 0, false)
	csr.Close()
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	result.Set("ActionList", ActionList)

	//select distinct from duplicate data
	encountered := map[string]bool{}
	ModuleList := []string{}

	for _, x := range ActionList {
		encountered[x.GetString("modules")] = true
	}
	for key, _ := range encountered {
		ModuleList = append(ModuleList, key)
	}

	result.Set("ModuleList", ModuleList)

	return c.SetResultInfo(false, "", result)
}

func (c *AclController) GetUserLoginActivityList2(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	type Filters struct {
		Field    string
		Operator string
		Value    string
		Type     string
	}
	type Filter struct {
		Filters []Filters
		Logic   string
	}
	type Sort struct {
		Field string
		Dir   string
	}
	p := struct {
		Skip     int
		Take     int
		Page     int
		PageSize int
		Filter   Filter
		Sort     []Sort
	}{}
	e := k.GetPayload(&p)

	resultuser := []tk.M{}

	d := new(acl.User)
	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.Startwith("groups", "CB_")).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&resultuser, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	logins := []interface{}{}
	// tk.Println("resultuser", resultuser)
	for _, i := range resultuser {
		login := i.Get("loginid")
		// tk.Println("login", logins)
		logins = append(logins, login)
	}

	result := []tk.M{}
	s := new(acl.Session)
	var arrsort []string

	query := []*db.Filter{}
	if p.Filter.Logic == "and" {
		if len(p.Filter.Filters) > 0 {
			// tk.Println("LEN", len(p.Filter.Filters))
			if len(p.Filter.Filters) == 1 {
				// tk.Println("LENS1", len(p.Filter.Filters))
				if p.Filter.Filters[0].Operator == "eq" {
					if p.Filter.Filters[0].Type == "date" {
						dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
						date := dates.UTC()
						// tk.Println("TIME", date)
						query = append(query, db.Eq(p.Filter.Filters[0].Field, date))
					} else {
						query = append(query, db.Eq(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					}
				}

				if p.Filter.Filters[0].Operator == "contains" {
					query = append(query, db.Contains(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))

				}
				if p.Filter.Filters[0].Operator == "neq" {
					query = append(query, db.Ne(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
				}
				if p.Filter.Filters[0].Operator == "startswith" {
					query = append(query, db.Startwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
				}

				if p.Filter.Filters[0].Operator == "lte" {
					dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
					date := dates.UTC()
					query = append(query, db.Lte(p.Filter.Filters[0].Field, date))
				}

				if p.Filter.Filters[0].Operator == "gte" {
					dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
					date := dates.UTC()
					query = append(query, db.Gte(p.Filter.Filters[0].Field, date))
				}

				if p.Filter.Filters[0].Operator == "endswith" {
					query = append(query, db.Endwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
				}
			}
			if len(p.Filter.Filters) == 2 {
				// tk.Println("LENS2", len(p.Filter.Filters))
				if p.Filter.Filters[1].Operator == "Eq" {
					if p.Filter.Filters[1].Type == "date" {
						dates, _ := time.Parse("20060102150405", p.Filter.Filters[1].Value)
						date := dates.UTC()
						// tk.Println("TIME", date)
						query = append(query, db.Eq(p.Filter.Filters[1].Field, date))
					} else {
						query = append(query, db.Eq(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
						query = append(query, db.Eq(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value))
					}
				}
				if p.Filter.Filters[1].Operator == "contains" {
					query = append(query, db.Contains(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					query = append(query, db.Contains(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value))
				}
				if p.Filter.Filters[1].Operator == "neq" {
					query = append(query, db.Ne(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					query = append(query, db.Ne(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value))
				}
				if p.Filter.Filters[1].Operator == "startswith" {
					query = append(query, db.Startwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					query = append(query, db.Startwith(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value))
				}

				if p.Filter.Filters[0].Operator == "lte" {
					dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
					date := dates.UTC()
					query = append(query, db.Lte(p.Filter.Filters[0].Field, date))
				}

				if p.Filter.Filters[1].Operator == "lte" {
					dates, _ := time.Parse("20060102150405", p.Filter.Filters[1].Value)
					date := dates.UTC()
					query = append(query, db.Lte(p.Filter.Filters[1].Field, date))
				}

				if p.Filter.Filters[0].Operator == "gte" {
					dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
					date := dates.UTC()
					query = append(query, db.Gte(p.Filter.Filters[0].Field, date))
					// tk.Println("DATE", p.Filter.Filters[0].Value)
					// tk.Println("DATE GTE", date)
				}
				if p.Filter.Filters[1].Operator == "gte" {
					dates, _ := time.Parse("20060102150405", p.Filter.Filters[1].Value)
					date := dates.UTC()
					query = append(query, db.Gte(p.Filter.Filters[1].Field, date))
					// tk.Println("DATE", p.Filter.Filters[0].Value)
					// tk.Println("DATE GTE", date)
				}

				if p.Filter.Filters[1].Operator == "endswith" {
					query = append(query, db.Endwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					query = append(query, db.Endwith(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value))
				}
			}
		}

	}

	if p.Filter.Logic == "or" {
		if len(p.Filter.Filters) > 0 {
			if len(p.Filter.Filters) == 1 {
				if p.Filter.Filters[0].Operator == "eq" {
					if p.Filter.Filters[0].Type == "date" {
						dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
						date := dates.UTC()
						// tk.Println("TIME", date)
						query = append(query, db.Eq(p.Filter.Filters[0].Field, date))
					} else {
						query = append(query, db.Eq(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					}

				}
				if p.Filter.Filters[0].Operator == "contains" {
					query = append(query, db.Contains(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
				}
				if p.Filter.Filters[0].Operator == "neq" {
					query = append(query, db.Ne(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
				}
				if p.Filter.Filters[0].Operator == "startswith" {
					if p.Filter.Filters[0].Type == "date" {
						dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
						date := dates.UTC()
						// tk.Println("TIME", date)
						query = append(query, db.Gte(p.Filter.Filters[0].Field, date))
					} else {
						query = append(query, db.Startwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					}
				}
				if p.Filter.Filters[0].Operator == "endswith" {
					if p.Filter.Filters[0].Type == "date" {
						dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
						date := dates.UTC()
						// tk.Println("TIME", date)
						query = append(query, db.Lte(p.Filter.Filters[0].Field, date))
					} else {
						query = append(query, db.Endwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					}
				}
			}
			if len(p.Filter.Filters) == 2 {
				// tk.Println("LENS2", len(p.Filter.Filters))
				if p.Filter.Filters[1].Operator == "Eq" {
					if p.Filter.Filters[1].Type == "date" {
						dates, _ := time.Parse("20060102150405", p.Filter.Filters[1].Value)
						date := dates.UTC()
						// tk.Println("TIME", date)
						query = append(query, db.Eq(p.Filter.Filters[1].Field, date))
					} else {
						query = append(query, db.Or(db.Eq(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value), db.Eq(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value)))
					}
				}
				if p.Filter.Filters[1].Operator == "contains" {
					query = append(query, db.Or(db.Contains(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value), db.Contains(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value)))
					// tk.Println("QUERY", query)
				}
				if p.Filter.Filters[1].Operator == "neq" {
					query = append(query, db.Or(db.Ne(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value), db.Ne(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value)))
				}
				if p.Filter.Filters[1].Operator == "startswith" {
					query = append(query, db.Or(db.Startwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value), db.Startwith(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value)))
				}

				if p.Filter.Filters[0].Operator == "lte" {
					dates0, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
					date0 := dates0.UTC()

					dates1, _ := time.Parse("20060102150405", p.Filter.Filters[1].Value)
					date1 := dates1.UTC()

					query = append(query, db.Or(db.Lte(p.Filter.Filters[0].Field, date0), db.Lte(p.Filter.Filters[1].Field, date1)))

				}

				if p.Filter.Filters[0].Operator == "gte" {
					dates0, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
					date0 := dates0.UTC()

					dates1, _ := time.Parse("20060102150405", p.Filter.Filters[1].Value)
					date1 := dates1.UTC()

					query = append(query, db.Or(db.Gte(p.Filter.Filters[0].Field, date0), db.Gte(p.Filter.Filters[1].Field, date1)))
				}

				if p.Filter.Filters[1].Operator == "endswith" {
					query = append(query, db.Or(db.Endwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value), db.Endwith(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value)))

				}
			}
		}

	}

	if len(p.Sort) > 0 {
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+p.Sort[0].Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(p.Sort[0].Field))
			}
		}
	}

	arrsort2 := strings.Join(arrsort, ",")

	// return arrsort2
	csr, e = c.AclCtx.Connection.NewQuery().From(s.TableName()).Where(db.In("loginid", logins...)).Skip(p.Skip).Take(p.Take).Order("-created").Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	if arrsort2 != "" {
		csr, e = c.AclCtx.Connection.NewQuery().From(s.TableName()).Where(db.In("loginid", logins...)).Skip(p.Skip).Take(p.Take).Order(arrsort2, "-created").Cursor(nil)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
	}

	if len(p.Filter.Filters) > 0 {
		query = append(query, db.In("loginid", logins...))
		csr, e = c.AclCtx.Connection.NewQuery().From(s.TableName()).Where(query...).Skip(p.Skip).Take(p.Take).Cursor(nil)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
		if arrsort2 != "" {
			csr, e = c.AclCtx.Connection.NewQuery().From(s.TableName()).Where(query...).Skip(p.Skip).Take(p.Take).Order(arrsort2, "-created").Cursor(nil)
			if e != nil {
				return c.ErrorResultInfo(e.Error(), nil)
			}
		}
	}

	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	// results2 := result

	// csr, err := c.AclCtx.Find(new(Access), query.Set("AGGR", "$sum"))
	// csr, err := c.AclCtx.Connection.NewQuery().From(s.TableName()).Cursor(nil)
	// defer csr.Close()
	// if err != nil {
	// 	return err.Error()
	// }

	data := struct {
		Data  []tk.M
		Total int
	}{
		Data:  result,
		Total: csr.Count(),
	}
	// fmt.Println(data)
	return data

}
func (c *AclController) GetUserLoginActivityList2Old(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	type Filters struct {
		Field    string
		Operator string
		Value    string
		Type     string
	}
	type Filter struct {
		Filters []Filters
		Logic   string
	}
	type Sort struct {
		Field string
		Dir   string
	}
	p := struct {
		Skip     int
		Take     int
		Page     int
		PageSize int
		Filter   Filter
		Sort     []Sort
	}{}
	e := k.GetPayload(&p)

	resultuser := []tk.M{}

	d := new(acl.User)
	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.Startwith("groups", "CB_")).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&resultuser, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	logins := []interface{}{}
	// tk.Println("resultuser", resultuser)
	for _, i := range resultuser {
		login := i.Get("loginid")
		// tk.Println("login", logins)
		logins = append(logins, login)
	}

	result := []tk.M{}
	s := new(acl.Session)
	var arrsort []string

	query := []*db.Filter{}
	if p.Filter.Logic == "and" {
		if len(p.Filter.Filters) > 0 {
			// tk.Println("LEN", len(p.Filter.Filters))
			if len(p.Filter.Filters) == 1 {
				// tk.Println("LENS1", len(p.Filter.Filters))
				if p.Filter.Filters[0].Operator == "eq" {
					if p.Filter.Filters[0].Type == "date" {
						dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
						date := dates.UTC()
						// tk.Println("TIME", date)
						query = append(query, db.Eq(p.Filter.Filters[0].Field, date))
					} else {
						query = append(query, db.Eq(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					}
				}

				if p.Filter.Filters[0].Operator == "contains" {
					query = append(query, db.Contains(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))

				}
				if p.Filter.Filters[0].Operator == "neq" {
					query = append(query, db.Ne(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
				}
				if p.Filter.Filters[0].Operator == "startswith" {
					query = append(query, db.Startwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
				}

				if p.Filter.Filters[0].Operator == "lte" {
					dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
					date := dates.UTC()
					query = append(query, db.Lte(p.Filter.Filters[0].Field, date))
				}

				if p.Filter.Filters[0].Operator == "gte" {
					dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
					date := dates.UTC()
					query = append(query, db.Gte(p.Filter.Filters[0].Field, date))
				}

				if p.Filter.Filters[0].Operator == "endswith" {
					query = append(query, db.Endwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
				}
			}
			if len(p.Filter.Filters) == 2 {
				// tk.Println("LENS2", len(p.Filter.Filters))
				if p.Filter.Filters[1].Operator == "Eq" {
					if p.Filter.Filters[1].Type == "date" {
						dates, _ := time.Parse("20060102150405", p.Filter.Filters[1].Value)
						date := dates.UTC()
						// tk.Println("TIME", date)
						query = append(query, db.Eq(p.Filter.Filters[1].Field, date))
					} else {
						query = append(query, db.Eq(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
						query = append(query, db.Eq(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value))
					}
				}
				if p.Filter.Filters[1].Operator == "contains" {
					query = append(query, db.Contains(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					query = append(query, db.Contains(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value))
				}
				if p.Filter.Filters[1].Operator == "neq" {
					query = append(query, db.Ne(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					query = append(query, db.Ne(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value))
				}
				if p.Filter.Filters[1].Operator == "startswith" {
					query = append(query, db.Startwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					query = append(query, db.Startwith(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value))
				}

				if p.Filter.Filters[0].Operator == "lte" {
					dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
					date := dates.UTC()
					query = append(query, db.Lte(p.Filter.Filters[0].Field, date))
				}

				if p.Filter.Filters[1].Operator == "lte" {
					dates, _ := time.Parse("20060102150405", p.Filter.Filters[1].Value)
					date := dates.UTC()
					query = append(query, db.Lte(p.Filter.Filters[1].Field, date))
				}

				if p.Filter.Filters[0].Operator == "gte" {
					dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
					date := dates.UTC()
					query = append(query, db.Gte(p.Filter.Filters[0].Field, date))
				}
				if p.Filter.Filters[1].Operator == "gte" {
					dates, _ := time.Parse("20060102150405", p.Filter.Filters[1].Value)
					date := dates.UTC()
					query = append(query, db.Gte(p.Filter.Filters[1].Field, date))
				}

				if p.Filter.Filters[1].Operator == "endswith" {
					query = append(query, db.Endwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					query = append(query, db.Endwith(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value))
				}
			}
		}

	}

	if p.Filter.Logic == "or" {
		if len(p.Filter.Filters) > 0 {
			if len(p.Filter.Filters) == 1 {
				if p.Filter.Filters[0].Operator == "eq" {
					if p.Filter.Filters[0].Type == "date" {
						dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
						date := dates.UTC()
						// tk.Println("TIME", date)
						query = append(query, db.Eq(p.Filter.Filters[0].Field, date))
					} else {
						query = append(query, db.Eq(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					}

				}
				if p.Filter.Filters[0].Operator == "contains" {
					query = append(query, db.Contains(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
				}
				if p.Filter.Filters[0].Operator == "neq" {
					query = append(query, db.Ne(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
				}
				if p.Filter.Filters[0].Operator == "startswith" {
					if p.Filter.Filters[0].Type == "date" {
						dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
						date := dates.UTC()
						// tk.Println("TIME", date)
						query = append(query, db.Gte(p.Filter.Filters[0].Field, date))
					} else {
						query = append(query, db.Startwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					}
				}
				if p.Filter.Filters[0].Operator == "endswith" {
					if p.Filter.Filters[0].Type == "date" {
						dates, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
						date := dates.UTC()
						// tk.Println("TIME", date)
						query = append(query, db.Lte(p.Filter.Filters[0].Field, date))
					} else {
						query = append(query, db.Endwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
					}
				}
			}
			if len(p.Filter.Filters) == 2 {
				// tk.Println("LENS2", len(p.Filter.Filters))
				if p.Filter.Filters[1].Operator == "Eq" {
					if p.Filter.Filters[1].Type == "date" {
						dates, _ := time.Parse("20060102150405", p.Filter.Filters[1].Value)
						date := dates.UTC()
						// tk.Println("TIME", date)
						query = append(query, db.Eq(p.Filter.Filters[1].Field, date))
					} else {
						query = append(query, db.Or(db.Eq(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value), db.Eq(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value)))
					}
				}
				if p.Filter.Filters[1].Operator == "contains" {
					query = append(query, db.Or(db.Contains(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value), db.Contains(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value)))
					// tk.Println("QUERY", query)
				}
				if p.Filter.Filters[1].Operator == "neq" {
					query = append(query, db.Or(db.Ne(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value), db.Ne(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value)))
				}
				if p.Filter.Filters[1].Operator == "startswith" {
					query = append(query, db.Or(db.Startwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value), db.Startwith(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value)))
				}

				if p.Filter.Filters[0].Operator == "lte" {
					dates0, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
					date0 := dates0.UTC()

					dates1, _ := time.Parse("20060102150405", p.Filter.Filters[1].Value)
					date1 := dates1.UTC()

					query = append(query, db.Or(db.Lte(p.Filter.Filters[0].Field, date0), db.Lte(p.Filter.Filters[1].Field, date1)))

				}

				if p.Filter.Filters[0].Operator == "gte" {
					dates0, _ := time.Parse("20060102150405", p.Filter.Filters[0].Value)
					date0 := dates0.UTC()

					dates1, _ := time.Parse("20060102150405", p.Filter.Filters[1].Value)
					date1 := dates1.UTC()

					query = append(query, db.Or(db.Gte(p.Filter.Filters[0].Field, date0), db.Gte(p.Filter.Filters[1].Field, date1)))
				}

				if p.Filter.Filters[1].Operator == "endswith" {
					query = append(query, db.Or(db.Endwith(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value), db.Endwith(p.Filter.Filters[1].Field, p.Filter.Filters[1].Value)))

				}
			}
		}

	}

	if len(p.Sort) > 0 {
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+p.Sort[0].Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(p.Sort[0].Field))
			}
		}
	}

	arrsort2 := strings.Join(arrsort, ",")

	// return arrsort2
	csr, e = c.AclCtx.Connection.NewQuery().From(s.TableName()).Where(db.In("loginid", logins...)).Skip(p.Skip).Take(p.Take).Order("-created").Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	if arrsort2 != "" {
		csr, e = c.AclCtx.Connection.NewQuery().From(s.TableName()).Where(db.In("loginid", logins...)).Skip(p.Skip).Take(p.Take).Order(arrsort2, "-created").Cursor(nil)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
	}

	if len(p.Filter.Filters) > 0 {
		query = append(query, db.In("loginid", logins...))
		csr, e = c.AclCtx.Connection.NewQuery().From(s.TableName()).Where(query...).Skip(p.Skip).Take(p.Take).Cursor(nil)
		if e != nil {
			return c.ErrorResultInfo(e.Error(), nil)
		}
		if arrsort2 != "" {
			csr, e = c.AclCtx.Connection.NewQuery().From(s.TableName()).Where(query...).Skip(p.Skip).Take(p.Take).Order(arrsort2, "-created").Cursor(nil)
			if e != nil {
				return c.ErrorResultInfo(e.Error(), nil)
			}
		}
	}

	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	results2 := result

	// csr, err := c.AclCtx.Find(new(Access), query.Set("AGGR", "$sum"))
	csr, err := c.AclCtx.Connection.NewQuery().From(s.TableName()).Cursor(nil)
	defer csr.Close()
	if err != nil {
		return err.Error()
	}

	data := struct {
		Data  []tk.M
		Total int
	}{
		Data:  results2,
		Total: csr.Count(),
	}
	// fmt.Println(data)
	return data

}

func (c *AclController) GetDetailLogBasedOnSession(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	type Filters struct {
		Field    string
		Operator string
		Value    string
	}
	type Filter struct {
		Filters []Filters
		Logic   string
	}
	// type Sort struct {
	// 	Field string
	// 	Dir   string
	// }
	p := struct {
		Skip     int
		Take     int
		Page     int
		PageSize int
		Filter   Filter
		// Sort []Sort
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	// csr, e := c.AclCtx.Find(new(LogModel), toolkit.M{})
	query := tk.M{}
	if p.Filter.Logic == "and" {
		if len(p.Filter.Filters) > 0 {
			if p.Filter.Filters[0].Operator == "eq" {
				query.Set("where", db.Eq(p.Filter.Filters[0].Field, p.Filter.Filters[0].Value))
			}
		}
	}

	csr, e := c.AclCtx.Find(new(LogModel), query.Set("order", []string{"DateAccess"}).Set("skip", p.Skip).Set("limit", p.Take))
	defer csr.Close()

	// if len(p.Sort) > 0 {
	// 	var arrsort []string
	// 	for _, val := range p.Sort {
	// 		if val.Dir == "desc" {
	// 			arrsort = append(arrsort, strings.ToLower("-"+p.Sort[0].Field))
	// 		} else {
	// 			arrsort = append(arrsort, strings.ToLower(p.Sort[0].Field))
	// 		}
	// 	}
	// 	csr, e = c.AclCtx.Find(new(LogModel), toolkit.M{}.Set("order", arrsort).Set("skip", p.Skip).Set("limit", p.Take))
	// 	defer csr.Close()

	// 	// fmt.Println(p.Sort[0].Field, p.Sort[0].Dir)
	// }
	results := make([]LogModel, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error()
	}
	results2 := results

	csr, err := c.AclCtx.Find(new(LogModel), query.Set("AGGR", "$sum"))
	defer csr.Close()
	if err != nil {
		return err.Error()
	}

	data := struct {
		Data  []LogModel
		Total int
	}{
		Data:  results2,
		Total: csr.Count(),
	}
	// fmt.Println(data)
	return data

	// result := tk.M{}
	// result.Set("LoginActivity", new(acl.Session))
	// return c.SetResultInfo(false, "", data)
}

func (c *AclController) ExportXLSUserLoginActivityLog(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		tk.Printf(err.Error())
	}

	k.Config.OutputType = knot.OutputJson

	type SubFilters struct {
		Field    string
		Operator string
		Value    string
		Type     string
	}
	type Filters struct {
		Field    string
		Operator string
		Value    string
		Type     string
		Filters  []SubFilters
	}
	type Filter struct {
		Filters []Filters
		Logic   string
	}
	type Sort struct {
		Field string
		Dir   string
	}

	p := struct {
		Filter           Filter
		Id               string
		LoginId          string
		FullName         string
		Country          string
		GroupDescription string
		Created          string
		Expired          string
		Modules          []string
		Sort             []Sort
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	result := []tk.M{}
	sort := []string{}
	query := []*db.Filter{}
	queryFilter := []*db.Filter{}

	// query = append(query, db.Eq("sessionid", p.Id))

	if len(p.Filter.Filters) > 0 {
		for _, f := range p.Filter.Filters {

			if len(f.Filters) == 0 {
				if f.Type == "date" {
					tmpdate, _ := time.Parse("20060102150405", f.Value)
					tmpdate = tmpdate.UTC()
					value := tmpdate
					switch f.Operator {
					case "eq":
						queryFilter = append(queryFilter, db.Eq(f.Field, value))
					case "gte":
						queryFilter = append(queryFilter, db.Gte(f.Field, value))
					case "lte":
						queryFilter = append(queryFilter, db.Lte(f.Field, value))
					default:
						break
					}
				} else {
					value := f.Value
					switch f.Operator {
					case "contains":
						queryFilter = append(queryFilter, db.Contains(f.Field, value))
					case "eq":
						queryFilter = append(queryFilter, db.Eq(f.Field, value))
					case "startswith":
						queryFilter = append(queryFilter, db.Startwith(f.Field, value))
					case "endswith":
						queryFilter = append(queryFilter, db.Endwith(f.Field, value))
					default:
						break
					}

				}
				if len(queryFilter) > 0 {
					if p.Filter.Logic == "and" {
						query = append(query, db.And(queryFilter...))
					} else {
						query = append(query, db.Or(queryFilter...))
					}
				}

			} else {
				// tk.Println("MASUK MULTIPLE")
				MultipleFilter := []*db.Filter{}
				for _, x := range p.Filter.Filters {
					// tk.Println("Filters", x.Filters)
					for _, y := range x.Filters {
						// tk.Println("ObjFilter", y)
						if y.Type == "date" {
							tmpdate, _ := time.Parse("20060102150405", y.Value)
							tmpdate = tmpdate.UTC()
							value := tmpdate
							switch y.Operator {
							case "eq":
								MultipleFilter = append(MultipleFilter, db.Eq(y.Field, value))
							case "gte":
								MultipleFilter = append(MultipleFilter, db.Gte(y.Field, value))
							case "lte":
								MultipleFilter = append(MultipleFilter, db.Lte(y.Field, value))
							default:
								break
							}
						} else {
							value := y.Value
							switch y.Operator {
							case "contains":
								MultipleFilter = append(MultipleFilter, db.Contains(y.Field, value))
							case "eq":
								MultipleFilter = append(MultipleFilter, db.Eq(y.Field, value))
							case "startswith":
								MultipleFilter = append(MultipleFilter, db.Startwith(y.Field, value))
							case "endswith":
								MultipleFilter = append(MultipleFilter, db.Endwith(y.Field, value))
							default:
								break
							}

						}
					}

				}
				if len(MultipleFilter) > 0 {
					if p.Filter.Logic == "and" {
						query = append(query, db.And(MultipleFilter...))
					} else {
						query = append(query, db.Or(MultipleFilter...))
					}
				}
			}
		}
	}

	if len(p.Sort) > 0 {
		for _, x := range p.Sort {
			if x.Dir == "asc" {
				sort = append(sort, x.Field)
			} else {
				sort = append(sort, "-"+x.Field)
			}
		}
	} else {
		sort = append(sort, "-dateaccess")
	}

	if len(p.Modules) > 0 {
		for _, x := range p.Modules {
			query = append(query, db.Eq("module", x))
		}

	}

	query = append(query, db.Eq("sessionid", p.Id))

	csr, e := c.AclCtx.Connection.NewQuery().From("acl_log").Where(db.And(query...)).Order(sort...).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	es := csr.Fetch(&result, 0, false)
	csr.Close()
	if es != nil {
		return c.ErrorResultInfo(es.Error(), nil)
	}

	created, _ := time.Parse("20060102150405", p.Created)
	createddate := created.Format("Monday, Jan _2 2006, 15:04:05")

	expired, _ := time.Parse("20060102150405", p.Expired)
	expireddate := expired.Format("Monday, Jan _2 2006, 15:04:05")

	font := xlsx.NewFont(11, "Calibri")
	style := xlsx.NewStyle()
	style.Font = *font

	fontHdr := xlsx.NewFont(11, "Calibri")
	fontHdr.Bold = true
	styleHdr := xlsx.NewStyle()
	styleHdr.Font = *fontHdr

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Session Id"

	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = p.Id

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Login Id"

	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = p.LoginId

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Full Name"

	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = p.FullName

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Country"

	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = p.Country

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "User Role(s)"

	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = p.GroupDescription

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Last Login"

	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = createddate

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Expired Date"

	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = expireddate

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = ""

	cell = row.AddCell()
	cell.SetStyle(style)
	cell.Value = ""

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Action"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "What Changed"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Old Value"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "New Value"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Action Time"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Resource Url"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Sources"

	for _, data := range result {

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = strings.Title(data.GetString("do"))

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("whatchanged")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = strings.Title(data.GetString("oldvalue"))

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = strings.Title(data.GetString("newvalue"))

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.Get("dateaccess").(time.Time).Format("Monday, Jan _2 2006, 15:04:05")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("requesturi")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("sources")

	}
	ExcelFilename := "Application Audit Trails Details " + time.Now().Format("20060102150405") + ".xlsx"
	err = file.Save(c.DownloadPath + "/" + ExcelFilename)

	if err != nil {
		tk.Printf(err.Error())
	}

	return c.SetResultInfo(false, "", ExcelFilename)
}

func (c *AclController) ExportXLSUserLoginActivity(k *knot.WebContext) interface{} {
	// c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		tk.Printf(err.Error())
	}

	resultuser := []tk.M{}

	// Get PayLoad
	type Filters struct {
		Field    string
		Operator string
		Value    string
		Type     string
	}
	type Filter struct {
		Filters []Filters
		Logic   string
	}
	type Sort struct {
		Field string
		Dir   string
	}
	parm := struct {
		Filter Filter
		Sort   []Sort
	}{}
	e := k.GetPayload(&parm)

	// Get LoginID Data
	d := new(acl.User)
	csr, e := c.AclCtx.Connection.NewQuery().From(d.TableName()).Where(db.Startwith("groups", "CB_")).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&resultuser, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	logins := []interface{}{}
	for _, i := range resultuser {
		login := i.Get("loginid")
		logins = append(logins, login)
	}

	result := []tk.M{}
	query := []*db.Filter{}
	sort := []string{}
	// Do some shit here
	query = append(query, db.In("loginid", logins...))
	queryFilter := []*db.Filter{}
	if len(parm.Filter.Filters) > 0 {
		for _, f := range parm.Filter.Filters {
			if f.Type == "date" {
				tmpdate, _ := time.Parse("20060102150405", f.Value)
				tmpdate = tmpdate.UTC()
				value := tmpdate
				switch f.Operator {
				case "eq":
					queryFilter = append(queryFilter, db.Eq(f.Field, value))
				case "gte":
					queryFilter = append(queryFilter, db.Gte(f.Field, value))
				case "lte":
					queryFilter = append(queryFilter, db.Lte(f.Field, value))
				default:
					break
				}
			} else {
				value := f.Value
				switch f.Operator {
				case "contains":
					queryFilter = append(queryFilter, db.Contains(f.Field, value))
				case "eq":
					queryFilter = append(queryFilter, db.Eq(f.Field, value))
				default:
					break
				}

			}
		}
	}
	if len(queryFilter) > 0 {
		if parm.Filter.Logic == "and" {
			query = append(query, db.And(queryFilter...))
		} else {
			query = append(query, db.Or(queryFilter...))
		}
	}

	s := new(acl.Session)
	// Get Sort Value
	if len(parm.Sort) > 0 {
		for _, x := range parm.Sort {
			if x.Dir == "asc" {
				sort = append(sort, x.Field)
			} else {
				sort = append(sort, "-"+x.Field)
			}
		}
	} else {
		sort = append(sort, "-created")
	}
	csr, e = c.AclCtx.Connection.NewQuery().From(s.TableName()).Where(db.And(query...)).Order(sort...).Cursor(nil)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	e = csr.Fetch(&result, 0, false)
	csr.Close()
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	font := xlsx.NewFont(11, "Calibri")
	style := xlsx.NewStyle()
	style.Font = *font

	fontHdr := xlsx.NewFont(11, "Calibri")
	fontHdr.Bold = true
	styleHdr := xlsx.NewStyle()
	styleHdr.Font = *fontHdr

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Session Id"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Login Id"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Last Login"

	cell = row.AddCell()
	cell.SetStyle(styleHdr)
	cell.Value = "Expired Date"

	for _, data := range result {
		created := data.Get("created").(time.Time).UTC().Format("Monday, Jan _2 2006, 15:04:05")
		expired := data.Get("expired").(time.Time).UTC().Format("Monday, Jan _2 2006, 15:04:05")
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("id")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = data.GetString("loginid")

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = created

		cell = row.AddCell()
		cell.SetStyle(style)
		cell.Value = expired

	}
	ExcelFilename := "Application Audit Trail " + time.Now().Format("20060102150405") + ".xlsx"

	err = file.Save(c.DownloadPath + "/" + ExcelFilename)

	if err != nil {
		tk.Printf(err.Error())
	}

	return c.SetResultInfo(false, "", ExcelFilename)
}
