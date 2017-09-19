package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type BusinessMetricsNotificationModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	BMId          string
	BMName        string
	BMData        []BusinessMetricsDataTempModel
	Period        time.Time //optional. Owner : Period Upload,Finance : none
	Type          string    // "owner/finance"
	Updated_Date  time.Time
	Source        string
	Updated_By    string
	HasOpen_By    []string
}

func NewBusinessMetricsNotificationModel() *BusinessMetricsNotificationModel {
	m := new(BusinessMetricsNotificationModel)
	m.Id = bson.NewObjectId()
	return m

}

func (e *BusinessMetricsNotificationModel) RecordID() interface{} {
	return e.Id
}

func (m *BusinessMetricsNotificationModel) TableName() string {
	return "BusinessMetricsNotification"
}
