package models

import (
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type BusinessMetricsDataTempModel struct {
	orm.ModelBase             `bson:"-",json:"-"`
	Id                        bson.ObjectId `bson:"_id" , json:"_id" `
	Period                    time.Time
	Year                      int
	BusinessName              string
	SCId                      int
	ScorecardCategory         string
	BMId                      string
	BusinessMetric            string
	BusinessMetricDescription string
	MajorRegion               string
	Region                    string
	Country                   string
	CountryCode               string
	Baseline                  float64
	Actual                    float64
	ActualYTD                 float64
	FullYearForecast          float64
	Target                    float64
	CreatedDate               time.Time
	CreatedBy                 string
	UpdatedDate               time.Time
	UpdatedBy                 string
	IsCurrent                 int
	SourceUID                 string
	Source                    string
	NABaseline                bool
	NAActual                  bool
	NATarget                  bool
	RAG                       string
	NABudget                  bool
	Budget                    float64
	RemainingBudget           float64
	RemainingBudgetOpposite   float64
}

func NewBusinessMetricsDataTempModel() *BusinessMetricsDataTempModel {
	m := new(BusinessMetricsDataTempModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *BusinessMetricsDataTempModel) RecordID() interface{} {
	return e.Id
}

func (s *BusinessMetricsDataTempModel) TableName() string {
	return "BusinessMetricsDataTemp"
}
