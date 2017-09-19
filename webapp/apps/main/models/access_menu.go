package models

import (
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/orm/v1"
	"github.com/eaciit/toolkit"
)

type AccessMenuModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string `json:"_id",bson:"_id"`
	ParentId      string
	Url           string
	Title         string
	Icon          string
	Index         int
	ApplicationID string

	Category       acl.AccessCategoryEnum
	Group1         string
	Group2         string
	Group3         string
	Enable         bool
	SpecialAccess1 string
	SpecialAccess2 string
	SpecialAccess3 string
	SpecialAccess4 string
}

func (a *AccessMenuModel) TableName() string {
	return "acl_access"
}

func (a *AccessMenuModel) RecordID() interface{} {
	return a.ID
}

func (a *AccessMenuModel) PreSave() error {
	if a.ID == "" {
		a.ID = toolkit.RandomString(32)
	}
	return nil
}
