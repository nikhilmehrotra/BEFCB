package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UserAuditTrailModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	UserID        string
	SessionID     string
	ActionTime    time.Time
	TypeOfChange  string
	UserIDChanged string
	FieldChanged  string
	OldValue      string
	NewValue      string
}

func NewUserAuditTrailModel() *UserAuditTrailModel {
	m := new(UserAuditTrailModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *UserAuditTrailModel) RecordID() interface{} {
	return e.Id
}

func (s *UserAuditTrailModel) TableName() string {
	return "acl_useraudittrail"
}
