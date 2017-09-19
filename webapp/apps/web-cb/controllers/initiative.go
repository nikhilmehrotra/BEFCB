package controllers

import (
	"eaciit/scb-apps/webapp/apps/web-cb/helper"
	. "eaciit/scb-apps/webapp/apps/web-cb/models"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type InitiativeController struct {
	*BaseController
}

func (c *InitiativeController) Remove(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		Id string
	}{}
	e := k.GetPayload(&parm)
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}

	// Get MAP Reference
	SummaryBusinessDriver := []SummaryBusinessDriverModel{}
	csr, err := c.Ctx.Connection.NewQuery().From(new(SummaryBusinessDriverModel).TableName()).Cursor(nil)
	err = csr.Fetch(&SummaryBusinessDriver, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()
	BusinessDriverL1 := []BusinessDriverL1Model{}
	csr, err = c.Ctx.Connection.NewQuery().From(new(BusinessDriverL1Model).TableName()).Cursor(nil)
	err = csr.Fetch(&BusinessDriverL1, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()
	LifeCycles := []LifeCycleModel{}
	csr, err = c.Ctx.Connection.NewQuery().From(new(LifeCycleModel).TableName()).Cursor(nil)
	err = csr.Fetch(&LifeCycles, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()
	REFScorecard, REFLifeCycle, REFBusinessDriver, REFKeyMetrics := map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}
	for _, x := range SummaryBusinessDriver {
		REFScorecard[x.Parentid] = x.Parentname
		REFBusinessDriver[x.Idx] = x.Name
	}
	for _, x := range LifeCycles {
		REFLifeCycle[x.LifeCycleId] = x.Name
	}
	for _, x := range BusinessDriverL1 {
		for _, m := range x.BusinessMetric {
			REFKeyMetrics[m.Id] = m.Description
		}
	}

	eData := new(InitiativeModel)
	csr, err = c.Ctx.Connection.NewQuery().From(eData.TableName()).Where(dbox.Eq("_id", bson.ObjectIdHex(parm.Id))).Cursor(nil)
	err = csr.Fetch(eData, 1, false)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()
	eSetAsComplete := "NO"
	if eData.SetAsComplete {
		eSetAsComplete = "YES"
	}
	eProjectManager := strings.Join(eData.ProjectManager, ",")
	eAccountableExecutive := strings.Join(eData.AccountableExecutive, ",")
	eTechnologyLead := strings.Join(eData.TechnologyLead, ",")
	ePlannedCost := strconv.FormatFloat(eData.PlannedCost, 'E', -1, 64)
	eIsGlobal := "NO"
	if eData.IsGlobal {
		eIsGlobal = "YES"
	}
	LOG_NAME := "Remove Initiative"
	MODULES := "Inititative"
	c.Action(k, MODULES, LOG_NAME, "Project Name", eData.ProjectName, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Clarity ID", eData.InvestmentId, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Start Date", eData.StartDate.Format("02-01-2006"), "", "", "")
	c.Action(k, MODULES, LOG_NAME, "End Date", eData.FinishDate.Format("02-01-2006"), "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Set As Complete", eSetAsComplete, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Project Manager", eProjectManager, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Accountable Executive", eAccountableExecutive, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Technology Lead", eTechnologyLead, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "BEF Sponsor", eData.Sponsor, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Problem Statement", eData.ProblemStatement, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Project Description", eData.ProjectDescription, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "RAG", eData.DisplayProgress, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Planned Cost", ePlannedCost, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Project Driver", eData.ProjectDriver, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Project Classification", eData.ProjectClassification, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Is Global", eIsGlobal, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Region", strings.Join(eData.Region, ", "), "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Country", strings.Join(eData.Country, ", "), "", "", "")

	c.Action(k, MODULES, LOG_NAME, "Scorecard Category", REFScorecard[eData.SCCategory], "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Life Cycle Stage", REFLifeCycle[eData.LifeCycleId], "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Initaitive Type", eData.Type, "", "", "")
	c.Action(k, MODULES, LOG_NAME, "Impact on Business Driver", eData.BusinessImpact, "", "", "")
	e = c.Ctx.DeleteMany(new(InitiativeModel), dbox.Eq("_id", bson.ObjectIdHex(parm.Id)))
	if e != nil {
		return c.ErrorResultInfo(e.Error(), nil)
	}
	return c.SetResultInfo(false, "", nil)
}

type AttachmentSort []Attachment

func (p AttachmentSort) Len() int           { return len(p) }
func (p AttachmentSort) Less(i, j int) bool { return p[i].FileName < p[j].FileName }
func (p AttachmentSort) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (c *InitiativeController) attachmentChecker(bsonid bson.ObjectId) (error, []Attachment) {
	crs, err := c.Ctx.Connection.NewQuery().Where(dbox.Eq("_id", bsonid)).From(NewInitiativeModel().TableName()).Cursor(nil)
	defer crs.Close()
	if !tk.IsNilOrEmpty(err) {
		return err, nil
	}

	data := []InitiativeModel{}
	err = crs.Fetch(&data, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return err, nil
	}

	// sort.Sort(AttachmentSort(data[0].Attachments))
	return nil, data[0].Attachments
}

func (c *InitiativeController) Save(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var (
		err error
	)
	Attachments := []Attachment{}
	/*existingfiletotal, _ := strconv.Atoi(k.Request.FormValue("existingfiletotal"))
	for j := 0; j < existingfiletotal; j++ {
		attachmentFile := Attachment{}
		attachmentFile.FileName = k.Request.FormValue("ExistingAttachment" + tk.ToString(j) + "FileName")
		attachmentFile.Description = k.Request.FormValue("ExistingAttachment" + tk.ToString(j) + "Description")
		attachmentFile.Updated_Date, _ = time.Parse("20060102", k.Request.FormValue("ExistingAttachment"+tk.ToString(j)+"UpdatedDate"))
		attachmentFile.Updated_By = k.Request.FormValue("ExistingAttachment" + tk.ToString(j) + "UpdatedBy")
		Attachments = append(Attachments, attachmentFile)
	}
	tak komen by rangga*/
	// <<<<<<< .mine
	tempID := k.Request.FormValue("Id")
	attc := []Attachment{}
	if tempID != "" && tempID != "0" {

		/*buat folder base initiative id*/
		tk.Println("filepath", filepath.Join(c.UploadPath, tk.ToString(k.Request.FormValue("InitiativeID"))))
		if _, err := os.Stat(filepath.Join(c.UploadPath, tk.ToString(k.Request.FormValue("InitiativeID")))); !tk.IsNilOrEmpty(err) {
			if os.IsNotExist(err) {
				os.Mkdir(filepath.Join(c.UploadPath, tk.ToString(k.Request.FormValue("InitiativeID"))), 0755)
			}
		}
		/*buat folder base initiative id*/

		err, attc = c.attachmentChecker(bson.ObjectIdHex(tempID))
		if !tk.IsNilOrEmpty(err) {
			return c.ErrorResultInfo(err.Error(), nil)
		}
	}
	// =======
	Attachments = append(Attachments, attc...)
	existCount := 0
	filetotal, _ := strconv.Atoi(k.Request.FormValue("filetotal"))
	for j := 0; j < filetotal; j++ {
		/*cek kalo file exist*/
		file, handler, err := k.Request.FormFile("FileUpload" + tk.ToString(j))
		defer file.Close()
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
		uploadfilename := handler.Filename

		tempAttch := []Attachment{}
		for _, v := range attc {
			if strings.Contains(v.FileName, uploadfilename) {
				tempAttch = append(tempAttch, v)
				existCount++
			}
		}
		sort.Sort(sort.Reverse(AttachmentSort(tempAttch)))
		if existCount > 0 {
			existCount = 1
		}

		if len(tempAttch) > 0 {
			if strings.Contains(tempAttch[0].FileName, " - ") {
				lastExistFile := tk.ToInt(strings.Split(strings.Split(tempAttch[0].FileName, "+")[1], " - ")[0], tk.RoundingAuto)
				existCount = lastExistFile + 1
			}
		}
		/*cek kalo file exist*/

		err, filename := helper.UploadHandler(k, "FileUpload"+tk.ToString(j), filepath.Join(c.UploadPath, tk.ToString(k.Request.FormValue("InitiativeID"))), existCount)
		if err != nil {
			tk.Println(err)
			return helper.CreateResult(false, nil, err.Error())
		}

		attachmentFile := Attachment{}
		attachmentFile.FileName = filename
		attachmentFile.Description = k.Request.FormValue("FileDescription" + tk.ToString(j))
		attachmentFile.Updated_Date = time.Now()
		attachmentFile.Updated_By = k.Session("username").(string)
		Attachments = append(Attachments, attachmentFile)
	}

	Milestones := []Milestone{}
	milestonetotal, _ := strconv.Atoi(k.Request.FormValue("milestonetotal"))
	for j := 0; j < milestonetotal; j++ {
		milestone := Milestone{}
		milestone.Name = k.Request.FormValue("Milestone" + tk.ToString(j))
		milestone.StartDate, _ = time.Parse("20060102", k.Request.FormValue("MilestoneStartDate"+tk.ToString(j)))
		milestone.EndDate, _ = time.Parse("20060102", k.Request.FormValue("MilestoneEndDate"+tk.ToString(j)))
		milestone.Country = strings.Split(k.Request.FormValue("MilestoneCountry"+tk.ToString(j)), "|")

		milestone.Completed, _ = strconv.ParseBool(k.Request.FormValue("MilestoneCompleted" + tk.ToString(j)))
		milestone.CompletedDate, _ = time.Parse("20060102", k.Request.FormValue("MilestoneCompletedDate"+tk.ToString(j)))
		milestone.Seq, _ = strconv.ParseFloat(k.Request.FormValue("MilestoneSeq"+tk.ToString(j)), 64)
		if milestone.Name == "" || len(milestone.Country) == 0 {

		} else {
			Milestones = append(Milestones, milestone)
		}

	}

	maptotal, _ := strconv.Atoi(k.Request.FormValue("maptotal"))

	// Get MAP Reference
	SummaryBusinessDriver := []SummaryBusinessDriverModel{}
	csr, err := c.Ctx.Connection.NewQuery().From(new(SummaryBusinessDriverModel).TableName()).Cursor(nil)
	err = csr.Fetch(&SummaryBusinessDriver, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()
	BusinessDriverL1 := []BusinessDriverL1Model{}
	csr, err = c.Ctx.Connection.NewQuery().From(new(BusinessDriverL1Model).TableName()).Cursor(nil)
	err = csr.Fetch(&BusinessDriverL1, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()
	LifeCycles := []LifeCycleModel{}
	csr, err = c.Ctx.Connection.NewQuery().From(new(LifeCycleModel).TableName()).Cursor(nil)
	err = csr.Fetch(&LifeCycles, 0, false)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Close()
	REFScorecard, REFLifeCycle, REFBusinessDriver, REFKeyMetrics := map[string]string{}, map[string]string{}, map[string]string{}, map[string]string{}
	for _, x := range SummaryBusinessDriver {
		REFScorecard[x.Parentid] = x.Parentname
		REFBusinessDriver[x.Idx] = x.Name
	}
	for _, x := range LifeCycles {
		REFLifeCycle[x.LifeCycleId] = x.Name
	}
	for _, x := range BusinessDriverL1 {
		for _, m := range x.BusinessMetric {
			REFKeyMetrics[m.Id] = m.Description
		}
	}

	mode := k.Request.FormValue("mode")
	LOG_NAME := ""
	MODULE := "Inititative"
	switch mode {
	case "new":
		LOG_NAME = "Add New Initiative"
		break
	case "edit":
		LOG_NAME = "Update Initiative"
		break
	case "copyclone":
		LOG_NAME = "Create a Copy Of Initiative"
		break
	default:
		break
	}
	c.Action(k, "Inititative", LOG_NAME, "", "", "", "", "")
	d := new(InitiativeModel)
	for i := 0; i < maptotal; i++ {

		var iMap = strings.Split(k.Request.FormValue("Map"+tk.ToString(i)), "|")
		var KeyMetrics = []string{}
		if k.Request.FormValue("KeyMetrics"+tk.ToString(i)) != "" {
			KeyMetrics = strings.Split(k.Request.FormValue("KeyMetrics"+tk.ToString(i)), "|")
		}
		iID, e := c.GetNextIdSeq("initiative")
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		if mode == "copyclone" {
			d.InitiativeID = tk.ToString(k.Request.FormValue("InitiativeID"))
			d.Id = bson.NewObjectId()
		} else if mode == "edit" {
			d.InitiativeID = tk.ToString(k.Request.FormValue("InitiativeID"))
			d.Id = bson.ObjectIdHex(k.Request.FormValue("Id"))
		} else {
			// New
			d.InitiativeID = tk.ToString(iID)
			d.Id = bson.NewObjectId()
		}
		d.ProjectName = strings.Trim(k.Request.FormValue("ProjectName"), " ")
		d.LifeCycleId = iMap[0]
		d.SubLifeCycleId = iMap[1]

		d.BusinessDriverId = iMap[2]
		d.BusinessDriverImpact = iMap[3]
		d.Type = iMap[5]
		d.SCCategory = iMap[6]
		d.CBLedInitiatives, e = strconv.ParseBool(k.Request.FormValue("CBLedInitiatives"))
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		d.EX, e = strconv.ParseBool(k.Request.FormValue("EX"))
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		d.OE, e = strconv.ParseBool(k.Request.FormValue("OE"))
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		d.IsGlobal, e = strconv.ParseBool(k.Request.FormValue("IsGlobal"))
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		Region := []string{}
		regiontotal, _ := strconv.Atoi(k.Request.FormValue("regiontotal"))
		for j := 0; j < regiontotal; j++ {
			Region = append(Region, k.Request.FormValue("region"+tk.ToString(j)))
		}

		countrytotal, _ := strconv.Atoi(k.Request.FormValue("countrytotal"))
		Country := []string{}
		for j := 0; j < countrytotal; j++ {
			Country = append(Country, k.Request.FormValue("country"+tk.ToString(j)))
		}

		d.Region = Region
		d.Country = Country

		d.StartDate, _ = time.Parse("20060102", k.Request.FormValue("StartDate"))
		d.FinishDate, _ = time.Parse("20060102", k.Request.FormValue("FinishDate"))
		d.ProblemStatement = k.Request.FormValue("ProblemStatement")
		d.ProjectDescription = k.Request.FormValue("ProjectDescription")
		d.ProjectDriver = strings.Trim(k.Request.FormValue("ProjectDriver"), " ")
		ProjectManager := strings.Trim(k.Request.FormValue("ProjectManager"), " ")
		if ProjectManager != "" {
			d.ProjectManager = strings.Split(ProjectManager, ",")
		} else {
			d.ProjectManager = []string{}
		}
		d.BusinessImpact = iMap[4]
		d.InvestmentId = k.Request.FormValue("InvestmentId")
		AccountableExecutive := strings.Trim(k.Request.FormValue("AccountableExecutive"), " ")
		if AccountableExecutive != "" {
			d.AccountableExecutive = strings.Split(AccountableExecutive, ",")
		} else {
			d.AccountableExecutive = []string{}
		}
		TechnologyLead := strings.Trim(k.Request.FormValue("TechnologyLead"), " ")
		if TechnologyLead != "" {
			d.TechnologyLead = strings.Split(TechnologyLead, ",")
		} else {
			d.TechnologyLead = []string{}
		}

		d.ProgressCompletion, _ = strconv.ParseFloat(k.Request.FormValue("ProgressCompletion"), 0)
		d.PlannedCost, _ = strconv.ParseFloat(k.Request.FormValue("PlannedCost"), 0)
		d.ProjectClassification = k.Request.FormValue("ProjectClassification")

		d.SetAsComplete, _ = strconv.ParseBool(k.Request.FormValue("SetAsComplete"))
		d.CompletedDate, _ = time.Parse("20060102", k.Request.FormValue("CompletedDate"))
		d.DisplayProgress = k.Request.FormValue("DisplayProgress")
		d.Sponsor = k.Request.FormValue("Sponsor")

		d.IsInitiativeTracked, _ = strconv.ParseBool(k.Request.FormValue("IsInitiativeTracked"))
		d.MetricBenchmark = k.Request.FormValue("MetricBenchmark")
		d.AdoptionScoreDenomination = k.Request.FormValue("AdoptionScoreDenomination")
		d.UsefulResources = k.Request.FormValue("UsefulResources")

		d.InitiativeType = "KeyEnablers"
		d.Attachments = Attachments
		d.Milestones = Milestones
		d.KeyMetrics = []KeyMetric{}
		if len(KeyMetrics) > 0 {
			for _, key := range KeyMetrics {
				keyArr := strings.Split(key, ",")
				keyData := KeyMetric{}
				keyData.BMId = keyArr[0]
				keyData.DirectIndirect, _ = strconv.Atoi(keyArr[1])
				d.KeyMetrics = append(d.KeyMetrics, keyData)
			}
		}
		if d.OE {
			d.ImprovedEfficiencyCurrent, _ = strconv.ParseFloat(k.Request.FormValue("ImprovedEfficiencyCurrent"), 0)
			d.ImprovedEfficiencyTarget, _ = strconv.ParseFloat(k.Request.FormValue("ImprovedEfficiencyTarget"), 0)
			d.ClientExperienceCurrent, _ = strconv.ParseFloat(k.Request.FormValue("ClientExperienceCurrent"), 0)
			d.ClientExperienceTarget, _ = strconv.ParseFloat(k.Request.FormValue("ClientExperienceTarget"), 0)
			d.OperationalImprovementCurrent, _ = strconv.ParseFloat(k.Request.FormValue("OperationalImprovementCurrent"), 0)
			d.OperationalImprovementTarget, _ = strconv.ParseFloat(k.Request.FormValue("OperationalImprovementTarget"), 0)
			d.CSRIncreaseCurrent, _ = strconv.ParseFloat(k.Request.FormValue("CSRIncreaseCurrent"), 0)
			d.CSRIncreaseTarget, _ = strconv.ParseFloat(k.Request.FormValue("CSRIncreaseTarget"), 0)
			d.TurnAroundTimeCurrent, _ = strconv.ParseFloat(k.Request.FormValue("TurnAroundTimeCurrent"), 0)
			d.TurnAroundTimeTarget, _ = strconv.ParseFloat(k.Request.FormValue("TurnAroundTimeTarget"), 0)
		}

		// For Logging
		SetAsComplete := "NO"
		if d.SetAsComplete {
			SetAsComplete = "YES"
		}
		IsGlobal := "NO"
		if d.IsGlobal {
			IsGlobal = "YES"
		}
		IsInitiativeTracked := "No"
		if d.IsInitiativeTracked {
			IsInitiativeTracked = "YES"
		}
		PlannedCost := strconv.FormatFloat(d.PlannedCost, 'E', -1, 64)
		switch mode {
		case "new":
			c.Action(k, MODULE, LOG_NAME, "Project Name", "", d.ProjectName, "", "")
			c.Action(k, MODULE, LOG_NAME, "Clarity ID", "", d.InvestmentId, "", "")
			c.Action(k, MODULE, LOG_NAME, "Start Date", "", d.StartDate.Format("02-01-2006"), "", "")
			c.Action(k, MODULE, LOG_NAME, "End Date", "", d.FinishDate.Format("02-01-2006"), "", "")
			c.Action(k, MODULE, LOG_NAME, "Set As Complete", "", SetAsComplete, "", "")
			c.Action(k, MODULE, LOG_NAME, "Project Manager", "", ProjectManager, "", "")
			c.Action(k, MODULE, LOG_NAME, "Accountable Executive", "", AccountableExecutive, "", "")
			c.Action(k, MODULE, LOG_NAME, "Technology Lead", "", TechnologyLead, "", "")
			c.Action(k, MODULE, LOG_NAME, "BEF Sponsor", "", d.Sponsor, "", "")
			c.Action(k, MODULE, LOG_NAME, "Problem Statement", "", d.ProblemStatement, "", "")
			c.Action(k, MODULE, LOG_NAME, "Project Description", "", d.ProjectDescription, "", "")
			c.Action(k, MODULE, LOG_NAME, "RAG", "", d.DisplayProgress, "", "")
			c.Action(k, MODULE, LOG_NAME, "Planned Cost", "", PlannedCost, "", "")
			c.Action(k, MODULE, LOG_NAME, "Project Driver", "", d.ProjectDriver, "", "")
			c.Action(k, MODULE, LOG_NAME, "Project Classification", "", d.ProjectClassification, "", "")
			c.Action(k, MODULE, LOG_NAME, "Is Global", "", IsGlobal, "", "")
			c.Action(k, MODULE, LOG_NAME, "Region", "", strings.Join(Region, ", "), "", "")
			c.Action(k, MODULE, LOG_NAME, "Country", "", strings.Join(Country, ", "), "", "")

			c.Action(k, MODULE, LOG_NAME, "IsInitiativeTracked", "", IsInitiativeTracked, "", "")
			c.Action(k, MODULE, LOG_NAME, "MetricBenchmark", "", d.MetricBenchmark, "", "")
			c.Action(k, MODULE, LOG_NAME, "Adoption Score Denomination", "", d.AdoptionScoreDenomination, "", "")
			c.Action(k, MODULE, LOG_NAME, "Useful Resources", "", d.UsefulResources, "", "")

			c.Action(k, MODULE, LOG_NAME, "Scorecard Category", "", REFScorecard[d.SCCategory], "", "")
			c.Action(k, MODULE, LOG_NAME, "Life Cycle Stage", "", REFLifeCycle[d.LifeCycleId], "", "")
			c.Action(k, MODULE, LOG_NAME, "Initaitive Type", "", d.Type, "", "")
			c.Action(k, MODULE, LOG_NAME, "Impact on Business Driver", "", d.BusinessImpact, "", "")
			break
		case "edit":
			eData := new(InitiativeModel)
			csr, err = c.Ctx.Connection.NewQuery().From(eData.TableName()).Where(dbox.Eq("_id", d.Id)).Cursor(nil)
			err = csr.Fetch(eData, 1, false)
			if err != nil {
				return c.ErrorResultInfo(err.Error(), nil)
			}
			csr.Close()
			e = c.Ctx.DeleteMany(d, dbox.Eq("_id", d.Id))
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
			eSetAsComplete := "NO"
			if eData.SetAsComplete {
				eSetAsComplete = "YES"
			}
			eProjectManager := strings.Join(eData.ProjectManager, ",")
			eAccountableExecutive := strings.Join(eData.AccountableExecutive, ",")
			eTechnologyLead := strings.Join(eData.TechnologyLead, ",")
			ePlannedCost := strconv.FormatFloat(eData.PlannedCost, 'E', -1, 64)
			eIsGlobal := "NO"
			if eData.IsGlobal {
				eIsGlobal = "YES"
			}
			c.Action(k, MODULE, LOG_NAME, "Project Name", eData.ProjectName, d.ProjectName, "", "")
			c.Action(k, MODULE, LOG_NAME, "Clarity ID", eData.InvestmentId, d.InvestmentId, "", "")
			c.Action(k, MODULE, LOG_NAME, "Start Date", eData.StartDate.Format("02-01-2006"), d.StartDate.Format("02-01-2006"), "", "")
			c.Action(k, MODULE, LOG_NAME, "End Date", eData.FinishDate.Format("02-01-2006"), d.FinishDate.Format("02-01-2006"), "", "")
			c.Action(k, MODULE, LOG_NAME, "Set As Complete", eSetAsComplete, SetAsComplete, "", "")
			c.Action(k, MODULE, LOG_NAME, "Project Manager", eProjectManager, ProjectManager, "", "")
			c.Action(k, MODULE, LOG_NAME, "Accountable Executive", eAccountableExecutive, AccountableExecutive, "", "")
			c.Action(k, MODULE, LOG_NAME, "Technology Lead", eTechnologyLead, TechnologyLead, "", "")
			c.Action(k, MODULE, LOG_NAME, "BEF Sponsor", eData.Sponsor, d.Sponsor, "", "")
			c.Action(k, MODULE, LOG_NAME, "Problem Statement", eData.ProblemStatement, d.ProblemStatement, "", "")
			c.Action(k, MODULE, LOG_NAME, "Project Description", eData.ProjectDescription, d.ProjectDescription, "", "")
			c.Action(k, MODULE, LOG_NAME, "RAG", eData.DisplayProgress, d.DisplayProgress, "", "")
			c.Action(k, MODULE, LOG_NAME, "Planned Cost", ePlannedCost, PlannedCost, "", "")
			c.Action(k, MODULE, LOG_NAME, "Project Driver", eData.ProjectDriver, d.ProjectDriver, "", "")
			c.Action(k, MODULE, LOG_NAME, "Project Classification", eData.ProjectClassification, d.ProjectClassification, "", "")
			c.Action(k, MODULE, LOG_NAME, "Is Global", eIsGlobal, IsGlobal, "", "")
			c.Action(k, MODULE, LOG_NAME, "Region", strings.Join(eData.Region, ", "), strings.Join(Region, ", "), "", "")
			c.Action(k, MODULE, LOG_NAME, "Country", strings.Join(eData.Country, ", "), strings.Join(Country, ", "), "", "")

			c.Action(k, MODULE, LOG_NAME, "IsInitiativeTracked", "", IsInitiativeTracked, "", "")
			c.Action(k, MODULE, LOG_NAME, "MetricBenchmark", "", d.MetricBenchmark, "", "")
			c.Action(k, MODULE, LOG_NAME, "Adoption Score Denomination", "", d.AdoptionScoreDenomination, "", "")
			c.Action(k, MODULE, LOG_NAME, "Useful Resources", "", d.UsefulResources, "", "")

			c.Action(k, MODULE, LOG_NAME, "Scorecard Category", REFScorecard[eData.SCCategory], REFScorecard[d.SCCategory], "", "")
			c.Action(k, MODULE, LOG_NAME, "Life Cycle Stage", REFLifeCycle[eData.LifeCycleId], REFLifeCycle[d.LifeCycleId], "", "")
			c.Action(k, MODULE, LOG_NAME, "Initaitive Type", eData.Type, d.Type, "", "")
			c.Action(k, MODULE, LOG_NAME, "Impact on Business Driver", eData.BusinessImpact, d.BusinessImpact, "", "")
			break
		case "copyclone":
			c.Action(k, MODULE, LOG_NAME, "Project Name", "", d.ProjectName, "", "")
			c.Action(k, MODULE, LOG_NAME, "Clarity ID", "", d.InvestmentId, "", "")
			c.Action(k, MODULE, LOG_NAME, "Start Date", "", d.StartDate.Format("02-01-2006"), "", "")
			c.Action(k, MODULE, LOG_NAME, "End Date", "", d.FinishDate.Format("02-01-2006"), "", "")
			c.Action(k, MODULE, LOG_NAME, "Set As Complete", "", SetAsComplete, "", "")
			c.Action(k, MODULE, LOG_NAME, "Project Manager", "", ProjectManager, "", "")
			c.Action(k, MODULE, LOG_NAME, "Accountable Executive", "", AccountableExecutive, "", "")
			c.Action(k, MODULE, LOG_NAME, "Technology Lead", "", TechnologyLead, "", "")
			c.Action(k, MODULE, LOG_NAME, "BEF Sponsor", "", d.Sponsor, "", "")
			c.Action(k, MODULE, LOG_NAME, "Problem Statement", "", d.ProblemStatement, "", "")
			c.Action(k, MODULE, LOG_NAME, "Project Description", "", d.ProjectDescription, "", "")
			c.Action(k, MODULE, LOG_NAME, "RAG", "", d.DisplayProgress, "", "")
			c.Action(k, MODULE, LOG_NAME, "Planned Cost", "", PlannedCost, "", "")
			c.Action(k, MODULE, LOG_NAME, "Project Driver", "", d.ProjectDriver, "", "")
			c.Action(k, MODULE, LOG_NAME, "Project Classification", "", d.ProjectClassification, "", "")
			c.Action(k, MODULE, LOG_NAME, "Is Global", "", IsGlobal, "", "")
			c.Action(k, MODULE, LOG_NAME, "Region", "", strings.Join(Region, ", "), "", "")
			c.Action(k, MODULE, LOG_NAME, "Country", "", strings.Join(Country, ", "), "", "")

			c.Action(k, MODULE, LOG_NAME, "IsInitiativeTracked", "", IsInitiativeTracked, "", "")
			c.Action(k, MODULE, LOG_NAME, "MetricBenchmark", "", d.MetricBenchmark, "", "")
			c.Action(k, MODULE, LOG_NAME, "Adoption Score Denomination", "", d.AdoptionScoreDenomination, "", "")
			c.Action(k, MODULE, LOG_NAME, "Useful Resources", "", d.UsefulResources, "", "")

			c.Action(k, MODULE, LOG_NAME, "Scorecard Category", "", REFScorecard[d.SCCategory], "", "")
			c.Action(k, MODULE, LOG_NAME, "Life Cycle Stage", "", REFLifeCycle[d.LifeCycleId], "", "")
			c.Action(k, MODULE, LOG_NAME, "Initaitive Type", "", d.Type, "", "")
			c.Action(k, MODULE, LOG_NAME, "Impact on Business Driver", "", d.BusinessImpact, "", "")
			break
		default:
			break
		}
		// End of Logging

		d.Updated_Date = time.Now()
		d.Updated_By = k.Session("username").(string)
		e = c.Ctx.Save(d)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		if mode != "copyclone" {
			c.reMappingInitiative(d)
		}
	}

	return helper.CreateResult(true, d, "")
}

func (c *InitiativeController) reMappingInitiative(d *InitiativeModel) {
	initiativeList := make([]InitiativeModel, 0)
	csr, e := c.Ctx.Find(new(InitiativeModel), tk.M{}.Set("where", dbox.Eq("InitiativeID", d.InitiativeID)))
	e = csr.Fetch(&initiativeList, 0, false)
	if e != nil {
		tk.Println("ERR : ", e.Error())
	}
	csr.Close()
	if len(initiativeList) > 1 {
		for _, i := range initiativeList {
			d.Id = i.Id
			d.LifeCycleId = i.LifeCycleId
			d.SubLifeCycleId = i.SubLifeCycleId
			d.BusinessDriverId = i.BusinessDriverId
			d.BusinessDriverImpact = i.BusinessDriverImpact
			d.KeyMetrics = i.KeyMetrics
			d.Type = i.Type
			if d.BusinessDriverId == "BD6" {
				d.InitiativeType = "SupportingEnablers"
			} else {
				d.InitiativeType = "KeyEnablers"
			}
			d.BusinessImpact = i.BusinessImpact
			e = c.Ctx.Save(d)
			if e != nil {
				tk.Println("ERR : ", e.Error())
			}
		}
	}
}

func (c *InitiativeController) RemoveAttacment(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	res := ResultInfo{}
	frm := struct {
		Attachment      tk.M
		AttachmentIndex int
		InitiativeID    string
		Initiative_id   string
	}{}
	err := k.GetPayload(&frm)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	err, attch := c.attachmentChecker(bson.ObjectIdHex(frm.Initiative_id))
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	attachments := []Attachment{}
	filenametobedelete := ""
	for _, v := range attch {
		if v.FileName != frm.Initiative_id+"+"+frm.Attachment.GetString("filename") && v.FileName != frm.Attachment.GetString("filename") {
			attachments = append(attachments, v)
		} else {
			if strings.Contains(v.FileName, "+") {
				filenametobedelete = frm.Initiative_id + "+" + frm.Attachment.GetString("filename")
			} else {
				filenametobedelete = frm.Attachment.GetString("filename")
			}
		}
	}

	m := NewInitiativeModel()
	err = c.Ctx.GetById(m, bson.ObjectIdHex(frm.Initiative_id))
	m.Attachments = attachments
	err = c.Ctx.Save(m)
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	/// remove file
	err = os.Remove(filepath.Join(filepath.Join(c.UploadPath, frm.Initiative_id), filenametobedelete))
	if !tk.IsNilOrEmpty(err) {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	res.Data = m
	return res
}

func (c *InitiativeController) Get(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	res := NewInitiativeModel()
	parm := struct {
		Id string
	}{}
	err := k.GetPayload(&parm)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	csr, err := c.Ctx.Connection.NewQuery().From(res.TableName()).Where(dbox.Eq("_id", parm.Id)).Cursor(nil)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr.Fetch(&res, 1, false)
	csr.Close()
	result, err := tk.ToM(res)
	return c.SetResultInfo(false, "", result)
}

func (c *InitiativeController) GetSummary(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	parm := struct {
		NoCBLED     bool
		NoBankWide  bool
		NoCompleted bool
		NoRemaining bool
		NoHigh      bool
		NoMedium    bool
		NoLow       bool
	}{}
	err := k.GetPayload(&parm)
	pipes := []tk.M{}
	group := tk.M{}
	query := []tk.M{}

	initiatives := []InitiativeModel{}
	query1 := []*dbox.Filter{}
	if parm.NoCBLED {
		query1 = append(query1, dbox.Ne("type", "CBLED"))
	}
	if parm.NoBankWide {
		query1 = append(query1, dbox.Ne("type", "BANKWIDE"))
	}
	if parm.NoHigh {
		query1 = append(query1, dbox.Ne("BusinessImpact", "High"))
	}
	if parm.NoMedium {
		query1 = append(query1, dbox.Ne("BusinessImpact", "Medium"))
	}
	if parm.NoLow {
		query1 = append(query1, dbox.Ne("BusinessImpact", "Low"))
	}

	if len(query1) == 0 {
		csr, err := c.Ctx.Connection.NewQuery().From("Initiative").Cursor(nil)
		err = csr.Fetch(&initiatives, 0, false)
		csr.Close()
		if err != nil {
			initiatives = []InitiativeModel{}
			// return nil
		}
	} else {
		csr, err := c.Ctx.Connection.NewQuery().From("Initiative").Where(dbox.And(query1...)).Cursor(nil)
		err = csr.Fetch(&initiatives, 0, false)
		csr.Close()
		if err != nil {
			initiatives = []InitiativeModel{}
			// return nil
		}
	}

	complete := 0.0

	InitiativeIdCompleted := []string{}

	for _, o := range initiatives {
		progress := c.GetProgress(o)
		if progress >= 100 {
			InitiativeIdCompleted = append(InitiativeIdCompleted, o.InitiativeID)
			complete = complete + 1
		}
	}

	// tk.Println("InitiativeIdCompleted", InitiativeIdCompleted)

	result := tk.M{}

	// parm.NoBankWide = true

	if parm.NoCBLED {
		query = append(query, tk.M{"type": tk.M{"$ne": "CBLED"}})
	}
	if parm.NoBankWide {
		query = append(query, tk.M{"type": tk.M{"$ne": "BANKWIDE"}})
	}
	if parm.NoHigh {
		query = append(query, tk.M{"BusinessImpact": tk.M{"$ne": "High"}})
	}
	if parm.NoMedium {
		query = append(query, tk.M{"BusinessImpact": tk.M{"$ne": "Medium"}})
	}
	if parm.NoLow {
		query = append(query, tk.M{"BusinessImpact": tk.M{"$ne": "Low"}})
	}
	if parm.NoCompleted {
		query = append(query, tk.M{"InitiativeID": tk.M{"$in": InitiativeIdCompleted}})
	}
	if parm.NoRemaining {
		query = append(query, tk.M{"InitiativeID": tk.M{"$nin": InitiativeIdCompleted}})
	}
	if len(query) > 0 {
		pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	}

	err = helper.Deserialize(`
   {"$group" : 
       {
            "_id": { "$sum": 1 },
            "BusinessImpactHigh" : {"$sum":{ "$cond": [ { "$eq": [ "$BusinessImpact", "High" ] }, { "$sum": 1 }, null ] }},
            "BusinessImpactMedium" : {"$sum":{ "$cond": [ { "$eq": [ "$BusinessImpact", "Medium" ] }, { "$sum": 1 }, null ] }},
            "BusinessImpactLow" : {"$sum":{ "$cond": [ { "$eq": [ "$BusinessImpact", "Low" ] },{ "$sum": 1 } , null ] }},
            "CBLED" : {"$sum":{ "$cond": [ { "$eq": [ "$type", "CBLED" ] }, { "$sum": 1 }, null ] }},
            "BANKWIDE" : {"$sum":{ "$cond": [ { "$eq": [ "$type", "BANKWIDE" ] },{ "$sum": 1 } , null ] }},
            "Initiative" : { "$sum": 1 }
       }
    }
	`, &group)

	pipes = append(pipes, group)

	csr, err := c.Ctx.Connection.NewQuery().Command("pipe", pipes).
		From("Initiative").
		Cursor(nil)
	err = csr.Fetch(&result, 1, false)
	csr.Close()
	if err != nil {
		result.Set("BusinessImpactHigh", 0)
		result.Set("BusinessImpactMedium", 0)
		result.Set("BusinessImpactLow", 0)
		result.Set("CBLED", 0)
		result.Set("BANKWIDE", 0)
		result.Set("Initiative", 0)
		// return c.SetResultInfo(true, err.Error(), nil)
	}

	InitiativeResult := result.GetFloat64("Initiative")
	remaining := InitiativeResult - complete

	result.Set("Initiatives_YTDRemaining", remaining)
	result.Set("Initiatives_YTDCompleted", complete)

	return c.SetResultInfo(false, "", result)
}

func (c *InitiativeController) GetProgress(o InitiativeModel) float64 {
	progress := 0.0

	now := time.Now().UTC()
	sd := o.StartDate
	fd := o.FinishDate
	// tk.Println(o.InitiativeID)
	isAnyMilestones := false
	if o.Milestones != nil && len(o.Milestones) > 0 {
		totalMilestones := 0
		for _, m := range o.Milestones {
			if m.Name != "" {
				totalMilestones++
			}
		}
		if totalMilestones > 0 {
			isAnyMilestones = true
		}
	}
	if o.Milestones != nil && len(o.Milestones) > 0 && isAnyMilestones {

		TotalDays := 0.0
		for _, d := range o.Milestones {
			Days := 0.0
			sd := d.StartDate
			ed := d.EndDate
			Hours := ed.Sub(sd).Hours()
			Days = tk.Div(Hours, 24)
			TotalDays = TotalDays + Days
		}

		SumProgress := 0.0
		for _, d := range o.Milestones {
			Days := 0.0
			sd := d.StartDate
			ed := d.EndDate
			Hours := ed.Sub(sd).Hours()
			Days = tk.Div(Hours, 24)

			now := time.Now().UTC()
			DelivWeight := tk.Div(Days, TotalDays)
			// tk.Println("DelivWeight:", DelivWeight)
			HoursProgress := now.Sub(sd).Hours()
			DaysProgress := tk.Div(HoursProgress, 24)
			// tk.Println("DaysProgress:", DaysProgress)

			if ed.Before(now) {
				SumProgress += (DelivWeight)
				// tk.Println("InitiativeID", o.InitiativeID, "| COMPLETE")
				// tk.Println("SumProgress COMPLETE", SumProgress)
			} else {
				// tk.Println("InitiativeID", o.InitiativeID, "| REMAINING")
				TimeProgress := tk.Div(DaysProgress, Days)
				SumProgress += (DelivWeight * TimeProgress)
				// tk.Println("SumProgress REMAINING", SumProgress)
				// tk.Println("DelivWeight:", DelivWeight)
				// tk.Println("TimeProgress:", TimeProgress)
				// tk.Println("SumProgress_remaining:", (DelivWeight * TimeProgress))
			}
		}

		return SumProgress * 100

	} else {
		if fd.Before(now) || o.SetAsComplete == true {
			progress = 100
		} else {
			a := fd.Sub(sd).Hours()
			b := now.Sub(sd).Hours()
			daya := tk.Div(a, 24)
			dayb := tk.Div(b, 24)
			progress = tk.Div(dayb, daya) * 100
		}
		// tk.Println(o.InitiativeID)
		return progress
	}

}
