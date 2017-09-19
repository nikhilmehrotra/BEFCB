package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ScorecardDetailModel struct {
	orm.ModelBase      `bson:"-",json:"-"`
	Id                 bson.ObjectId `bson:"_id" json:"_id" `
	SCId               string
	SCName             string
	SCDetailCategoryId string
	SCDetailCategory   string
	Name               string
	Type               string // Val : cumulative/spot
	Denomination       string // Val : DOLLAR/PERCENTAGE\NUMERIC
	ValueType          int    // val : 1(Higher is better)/0(Lower is better)
	DecimalFormat      string // 0/1/2
	Description        string

	// Breakdown
	MajorRegion string
	Region      string
	Country     string
	CountryCode string

	// Period Value
	Year      int
	Jan       float64
	Feb       float64
	Mar       float64
	Apr       float64
	May       float64
	Jun       float64
	Jul       float64
	Aug       float64
	Sep       float64
	Oct       float64
	Nov       float64
	Dec       float64
	Baseline  float64
	Benchmark float64
	Target    float64

	RAGJan string
	RAGFeb string
	RAGMar string
	RAGApr string
	RAGMay string
	RAGJun string
	RAGJul string
	RAGAug string
	RAGSep string
	RAGOct string
	RAGNov string
	RAGDec string

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
	NABaseline  bool
	NABenchmark bool
	NATarget    bool

	NARAGJan bool
	NARAGFeb bool
	NARAGMar bool
	NARAGApr bool
	NARAGMay bool
	NARAGJun bool
	NARAGJul bool
	NARAGAug bool
	NARAGSep bool
	NARAGOct bool
	NARAGNov bool
	NARAGDec bool

	MetricReference string
	Updated_Date    time.Time
	Updated_By      string

	// Exclude :
	// RAG                string
	// NARag       bool
}

func NewScorecardDetail() *ScorecardDetailModel {
	m := new(ScorecardDetailModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *ScorecardDetailModel) RecordID() interface{} {
	return e.Id
}

func (m *ScorecardDetailModel) TableName() string {
	return "ScorecardDetail"
}
