package models

import (
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type BusinessMetricsDataBudgetTempModel struct {
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
	Target                    float64
	CreatedDate               time.Time
	CreatedBy                 string
	UpdatedDate               time.Time
	UpdatedBy                 string
	// IsCurrent                 int
	SourceUID  string
	Source     string
	NABaseline bool
	NATarget   bool
	NABudget   bool
	Budget     float64
}

func NewBusinessMetricsDataBudgetTempModel() *BusinessMetricsDataBudgetTempModel {
	m := new(BusinessMetricsDataBudgetTempModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *BusinessMetricsDataBudgetTempModel) RecordID() interface{} {
	return e.Id
}

func (s *BusinessMetricsDataBudgetTempModel) TableName() string {
	return "BusinessMetricsDataBudgetTemp"
}
