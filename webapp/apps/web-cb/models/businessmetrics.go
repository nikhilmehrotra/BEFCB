/*
* @Author: Ainur
* @Date:   2016-11-08 12:28:44
* @Last Modified by:   Ainur
* @Last Modified time: 2017-05-31 01:39:05
 */

package models

import (
	"github.com/eaciit/orm"
	"time"
)

type BusinessMetric struct {
	orm.ModelBase           `bson:"-",json:"-"`
	Id                      string
	DataPoint               string ` bson:"DataPoint" , json:"DataPoint" `
	DecimalFormat           string ` bson:"DecimalFormat" , json:"DecimalFormat" `
	MetricType              string ` bson:"MetricType" , json:"MetricType" `
	BDId                    string ` bson:"bdid" , json:"bdid" `
	Description             string
	ValueType               int // 1:Higher is better, 0:Lower Is Better
	MetricDirection         int
	Type                    string //value : 'cumulative','spot'
	TargetPeriod            time.Time
	TargetPeriodStr         string
	TargetValue             float64
	BaseLinePeriod          time.Time
	BaseLinePeriodStr       string
	BaseLineValue           float64
	CurrentPeriod           time.Time
	CurrentPeriodStr        string
	CurrentValue            float64
	CurrentYTDValue         float64
	CurrentYTDValueVsBudget float64
	Display                 string
	NABaseline              bool
	NAActual                bool
	NATarget                bool
	NABudget                bool
	RegionalData            []MetricPartialData
	CountryData             []MetricPartialData
	ActualData              []ActualValue
	MetricFiles             []MetricFile ` bson:"MetricFiles" , json:"MetricFiles" `
	OrderIndex              int
	PreviousValue           float64
	UpdatedDate             time.Time
}

type ActualValue struct {
	Period    time.Time
	PeriodStr string
	Value     float64
	Flag      string
}
type MetricPartialData struct {
	Name                    string
	Baseline                float64
	CurrentYTDValue         float64
	CurrentYTDValueVsBudget float64
	Target                  float64
	NABaseline              bool
	NAActual                bool
	NABudget                bool
	NATarget                bool
	Rag                     string
}

func NewMetricPartialData() *MetricPartialData {
	m := new(MetricPartialData)
	return m
}

func NewBusinessMetric() *BusinessMetric {
	m := new(BusinessMetric)
	return m
}
