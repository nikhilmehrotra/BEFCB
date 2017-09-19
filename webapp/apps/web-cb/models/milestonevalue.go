package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MilestoneValueModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" json:"_id" `
	InitiativeId  string
	MetricId      string
	MetricName    string
	Denomination  string
	Description   string
	ActualValue   bool
	TotalValue    bool
	RAGValue      bool

	Country     string
	CountryCode string
	Year        int
	Jan         float64
	Feb         float64
	Mar         float64
	Apr         float64
	May         float64
	Jun         float64
	Jul         float64
	Aug         float64
	Sep         float64
	Oct         float64
	Nov         float64
	Dec         float64
	NAJan       bool
	NAFeb       bool
	NAMar       bool
	NAApr       bool
	NAMay       bool
	NAJun       bool
	NAJul       bool
	NAAug       bool
	NASep       bool
	NAOct       bool
	NANov       bool
	NADec       bool
	TotalJan    float64
	TotalFeb    float64
	TotalMar    float64
	TotalApr    float64
	TotalMay    float64
	TotalJun    float64
	TotalJul    float64
	TotalAug    float64
	TotalSep    float64
	TotalOct    float64
	TotalNov    float64
	TotalDec    float64
	NATotalJan  bool
	NATotalFeb  bool
	NATotalMar  bool
	NATotalApr  bool
	NATotalMay  bool
	NATotalJun  bool
	NATotalJul  bool
	NATotalAug  bool
	NATotalSep  bool
	NATotalOct  bool
	NATotalNov  bool
	NATotalDec  bool
	RAGJan      string
	RAGFeb      string
	RAGMar      string
	RAGApr      string
	RAGMay      string
	RAGJun      string
	RAGJul      string
	RAGAug      string
	RAGSep      string
	RAGOct      string
	RAGNov      string
	RAGDec      string

	LatestRAG string
	PeriodRAG string

	Created_By   string
	Created_Date time.Time
	Updated_By   string
	Updated_Date time.Time
}

func NewMilestoneValue() *MilestoneValueModel {
	m := new(MilestoneValueModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *MilestoneValueModel) RecordID() interface{} {
	return e.Id
}

func (m *MilestoneValueModel) TableName() string {
	return "MilestoneValue"
}
