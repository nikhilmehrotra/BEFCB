package models

import(
	"time"
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type InitiativeDataModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id bson.ObjectId `bson:"_id" , json:"_id" `
	PortfolioID	string `bson:"PortfolioID", json:"PortfolioID"`
	PortfolioName string `bson:"PortfolioName", json:"PortfolioName"`
	PortfolioPMOName string `bson:"PortfolioPMOName", json:"PortfolioPMOName"`
	PortfolioPMOBankID string `bson:"PortfolioPMOBankID", json:"PortfolioPMOBankID"`
	ProgramID string `bson:"ProgramID", json:"ProgramID"`
	ProgramName string `bson:"ProgramName", json:"ProgramName"`
	ProgramManagerName string `bson:"ProgramManagerName", json:"ProgramManagerName"`
	ProgramManagerBankID string `bson:"ProgramManagerBankID", json:"ProgramManagerBankID"`
	ProgramSponsorName string `bson:"ProgramSponsorName", json:"ProgramSponsorName"`
	ProgramSponsorBankID string `bson:"ProgramSponsorBankID", json:"ProgramSponsorBankID"`
	ProgramAEName string `bson:"ProgramAEName", json:"ProgramAEName"`
	ProgramAEBankID string `bson:"ProgramAEBankID", json:"ProgramAEBankID"`
	ProgramPMOName string `bson:"ProgramPMOName", json:"ProgramPMOName"`
	ProgramPMOBankID string `bson:"ProgramPMOBankID", json:"ProgramPMOBankID"`
	ProjectInternalID string `bson:"ProjectInternalID", json:"ProjectInternalID"`
	ProjectID string `bson:"ProjectID", json:"ProjectID"`
	ProjectName string `bson:"ProjectName", json:"ProjectName"`
	ProjectManagerName string `bson:"ProjectManagerName", json:"ProjectManagerName"`
	ProjectManagerBankID	string	`bson:"ProjectManagerBankID", json:"ProjectManagerBankID"`
	ProjectAEName string `bson:"ProjectAEName", json:"ProjectAEName"`
	ProjectAEBankID string `bson:"ProjectAEBankID", json:"ProjectAEBankID"`
	ProjectSponsorName string `bson:"ProjectSponsorName", json:"ProjectSponsorName"`
	ProjectSponsorBankID string `bson:"ProjectSponsorBankID", json:"ProjectSponsorBankID"`
	ProjectPMOName string `bson:"ProjectPMOName", json:"ProjectPMOName"`
	ProjectPMOBankID string `bson:"ProjectPMOBankID", json:"ProjectPMOBankID"`
	ProjectClasification string `bson:"ProjectClasification", json:"ProjectClasification"`
	ProjectDriver string `bson:"ProjectDriver", json:"ProjectDriver"`
	ProjectCategory string `bson:"ProjectCategory", json:"ProjectCategory"`
	ProjectType string `bson:"ProjectType", json:"ProjectType"`
	ProjectStartDate time.Time `bson:"ProjectStartDate", json:"ProjectStartDate"`
	ProjectFinishDate time.Time `bson:"ProjectFinishDate", json:"ProjectFinishDate"`
	BaselineStartDate time.Time `bson:"BaselineStartDate", json:"BaselineStartDate"`
	BaselineFinishDate time.Time `bson:"BaselineFinishDate", json:"BaselineFinishDate"`
	ProjectFinancialDepartment string `bson:"ProjectFinancialDepartment", json:"ProjectFinancialDepartment"`
	ProjectFinancialLocation string `bson:"ProjectFinancialLocation", json:"ProjectFinancialLocation"`
	ProjectObjective string `bson:"ProjectObjective", json:"ProjectObjective"`
	TotalCapexApproved  float64 `bson:"TotalCapexApproved", json:"TotalCapexApproved"`
	TotalOpexApproved float64 `bson:"TotalOpexApproved", json:"TotalOpexApproved"`
	TotalCapexSpend float64 `bson:"TotalCapexSpend", json:"TotalCapexSpend"`
	TotalOpexSpend float64 `bson:"TotalOpexSpend", json:"TotalOpexSpend"`
	OpporProbStatement string `bson:"OpporProbStatement", json:"OpporProbStatement"`
	ProposedSolution string `bson:"ProposedSolution", json:"ProposedSolution"`
	Assumptions string `bson:"Assumptions", json:"Assumptions"`
	KeyDeliverables string `bson:"KeyDeliverables", json:"KeyDeliverables"`
	OverallStatus string `bson:"OverallStatus", json:"OverallStatus"`
	ScheduleStatus string `bson:"ScheduleStatus", json:"ScheduleStatus"`
	SCBBenefitStatus string `bson:"SCBBenefitStatus", json:"SCBBenefitStatus"`
	CostStatus string `bson:"CostStatus", json:"CostStatus"`
}

func NewInitiativeDataModel() *InitiativeDataModel {
	m := new(InitiativeDataModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *InitiativeDataModel) RecordID() interface{} {
	return e.Id
}

func (s *InitiativeDataModel) TableName() string {
	return "InitiativeData"
}