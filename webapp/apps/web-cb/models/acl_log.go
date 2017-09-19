package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type LogModel struct {
	orm.ModelBase    `bson:"-",json:"-"`
	Id               bson.ObjectId `bson:"_id" , json:"_id" `
	UserID           string
	FullName         string
	FirstName        string
	LastName         string
	Country          string
	CountryCode      string
	Group            string
	GroupDescription string
	Module           string
	Do               string
	WhatChanged      string
	OldValue         string
	NewValue         string
	DateAccess       time.Time
	SessionID        string
	LoginTime        time.Time
	ExpiredTime      time.Time
	SourceType       string
	Sources          string
	RequestURI       string
	IPAddress        string
}

func NewLogModel() *LogModel {
	m := new(LogModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *LogModel) RecordID() interface{} {
	return e.Id
}

func (s *LogModel) TableName() string {
	return "acl_log"
}
