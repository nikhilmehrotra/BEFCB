package models

import (
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type InitiativeModel struct {
	orm.ModelBase         `bson:"-",json:"-"`
	Id                    bson.ObjectId `bson:"_id" , json:"_id" `
	InitiativeID          string        `bson:"InitiativeID" , json:"InitiativeID" `
	ProjectName           string        `bson:"ProjectName" , json:"ProjectName" `
	LifeCycleId           string        `bson:"LifeCycleId" , json:"LifeCycleId" `
	SubLifeCycleId        string        `bson:"SubLifeCycleId" , json:"SubLifeCycleId" `
	SCCategory            string        `bson:"SCCategory" , json:"SCCategory"`
	BusinessDriverId      string        `bson:"BusinessDriverId" , json:"BusinessDriverId" `
	BusinessDriverImpact  string        `bson:"BusinessDriverImpact" , json:"BusinessDriverImpact" `
	CBLedInitiatives      bool          `bson:"CBLedInitiatives" , json:"CBLedInitiatives" `
	EX                    bool          `bson:"EX" , json:"EX" `
	OE                    bool          `bson:"OE" , json:"OE" `
	StartDate             time.Time     `bson:"StartDate" , json:"StartDate" `
	FinishDate            time.Time     `bson:"FinishDate" , json:"FinishDate" `
	ProblemStatement      string        `bson:"ProblemStatement" , json:"ProblemStatement" `
	ProjectDescription    string        `bson:"ProjectDescription" , json:"ProjectDescription" `
	ProjectManager        []string      `bson:"ProjectManager" , json:"ProjectManager" `
	BusinessImpact        string        `bson:"BusinessImpact" , json:"BusinessImpact" `
	InvestmentId          string        `bson:"InvestmentId" , json:"InvestmentId" `
	AccountableExecutive  []string      `bson:"AccountableExecutive" , json:"AccountableExecutive" `
	TechnologyLead        []string      `bson:"TechnologyLead" , json:"TechnologyLead" `
	ProgressCompletion    float64       `bson:"ProgressCompletion" , json:"ProgressCompletion" `
	PlannedCost           float64       `bson:"PlannedCost" , json:"PlannedCost" `
	ProjectClassification string        `bson:"ProjectClassification" , json:"ProjectClassification" `
	Attachments           []Attachment  `bson:"Attachments" , json:"Attachments" `
	DeletedAttachments    []Attachment  `bson:"DeletedAttachments" , json:"DeletedAttachments" `
	InitiativeType        string        `bson:"InitiativeType" , json:"InitiativeType" `
	ProjectDriver         string        `bson:"ProjectDriver" , json:"ProjectDriver" `
	IsGlobal              bool          `bson:"IsGlobal" , json:"IsGlobal" `
	Region                []string      `bson:"Region" , json:"Region" `
	Country               []string      `bson:"Country" , json:"Country" `
	CommentList           []Comment     `bson:"CommentList" , json:"CommentList" `
	Type                  string
	SetAsComplete         bool      `bson:"SetAsComplete" , json:"SetAsComplete" `
	CompletedDate         time.Time `bson:"CompletedDate" , json:"CompletedDate" `
	DisplayProgress       string    `bson:"DisplayProgress" , json:"DisplayProgress" `
	Sponsor               string    `bson:"Sponsor" , json:"Sponsor" `

	IsInitiativeTracked       bool   `bson:"IsInitiativeTracked" , json:"IsInitiativeTracked" `
	MetricBenchmark           string `bson:"MetricBenchmark" , json:"MetricBenchmark" `
	AdoptionScoreDenomination string `bson:"AdoptionScoreDenomination" , json:"AdoptionScoreDenomination" `
	UsefulResources           string `bson:"UsefulResources" , json:"UsefulResources" `

	Milestones []Milestone `bson:"Milestones" , json:"Milestones" `
	KeyMetrics []KeyMetric `bson:"KeyMetrics" , json:"KeyMetrics" `
	// Metrics for OE
	ImprovedEfficiencyCurrent     float64 `bson:"ImprovedEfficiencyCurrent" , json:"ImprovedEfficiencyCurrent" `
	ImprovedEfficiencyTarget      float64 `bson:"ImprovedEfficiencyTarget" , json:"ImprovedEfficiencyTarget" `
	ClientExperienceCurrent       float64 `bson:"ClientExperienceCurrent" , json:"ClientExperienceCurrent" `
	ClientExperienceTarget        float64 `bson:"ClientExperienceTarget" , json:"ClientExperienceTarget" `
	OperationalImprovementCurrent float64 `bson:"OperationalImprovementCurrent" , json:"OperationalImprovementCurrent" `
	OperationalImprovementTarget  float64 `bson:"OperationalImprovementTarget" , json:"OperationalImprovementTarget" `
	CSRIncreaseCurrent            float64 `bson:"CSRIncreaseCurrent" , json:"CSRIncreaseCurrent" `
	CSRIncreaseTarget             float64 `bson:"CSRIncreaseTarget" , json:"CSRIncreaseTarget" `
	TurnAroundTimeCurrent         float64 `bson:"TurnAroundTimeCurrent" , json:"TurnAroundTimeCurrent" `
	TurnAroundTimeTarget          float64 `bson:"TurnAroundTimeTarget" , json:"TurnAroundTimeTarget" `
	Updated_Date                  time.Time
	Updated_By                    string
}

type KeyMetric struct {
	BMId           string `bson:"BMId" , json:"BMId" `
	Name           string `bson:"Name" , json:"Name" `
	DirectIndirect int    `bson:"DirectIndirect" , json:"DirectIndirect" `
}

func NewKeyMetric() *KeyMetric {
	m := new(KeyMetric)
	return m
}

type Milestone struct {
	Name          string    `bson:"Name" , json:"Name" `
	StartDate     time.Time `bson:"StartDate" , json:"StartDate" `
	EndDate       time.Time `bson:"EndDate" , json:"EndDate" `
	Country       []string  `bson:"Country" , json:"Country" `
	Completed     bool      `bson:"Completed" , json:"Completed" `
	CompletedDate time.Time `bson:"CompletedDate" , json:"CompletedDate" `
	Seq           float64   `bson:"Seq" , json:"Seq" `
}

func NewMilestone() *Milestone {
	m := new(Milestone)
	return m
}

type Attachment struct {
	Id           string
	FileName     string
	Description  string
	Updated_Date time.Time
	Updated_By   string
}

func NewAttachment() *Attachment {
	m := new(Attachment)
	return m
}

type Comment struct {
	Id           string
	Username     string `bson:"Username" , json:"Username" `
	DateInput    string `bson:"DateInput" , json:"DateInput" `
	Comment      string `bson:"Comment" , json:"Comment" `
	Updated_Date time.Time
	Updated_By   string
}

func NewInitiativeModel() *InitiativeModel {
	m := new(InitiativeModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *InitiativeModel) RecordID() interface{} {
	return e.Id
}

func (s *InitiativeModel) TableName() string {
	return "Initiative"
}
