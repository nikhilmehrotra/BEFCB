package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type AccessibilityModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            bson.ObjectId `bson:"_id" , json:"_id" `
	RoleID        string
	AccessID      string
	Title         string
	Url           string
	ParentID      string
	AllowStatus   bool
	Global        AccessibilityValueModel
	Region        AccessibilityValueModel
	Country       AccessibilityValueModel
}
type AccessibilityValueModel struct {
	Create  bool
	Read    bool
	Update  bool
	Delete  bool
	Owned   bool
	Curtain bool
	Upload  bool
}

func NewAccessibility() *AccessibilityModel {
	m := new(AccessibilityModel)
	m.ID = bson.NewObjectId()
	return m
}

func (e *AccessibilityModel) RecordID() interface{} {
	return e.ID
}

func (m *AccessibilityModel) TableName() string {
	return "acl_accessibility"
}
