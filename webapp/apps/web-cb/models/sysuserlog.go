package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type SysUserLogModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id"`
	Userid        bson.ObjectId
	Username      string
	Logintime     time.Time
}

func NewSysUserLogModel() *SysUserLogModel {
	m := new(SysUserLogModel)
	m.Id = bson.NewObjectId()
	m.Logintime = time.Now()
	return m
}

func (s *SysUserLogModel) RecordID() interface{} {
	return s.Id
}

func (s *SysUserLogModel) TableName() string {
	return "SysUsersLog"
}
