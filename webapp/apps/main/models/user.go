package models

import (
	"errors"
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/orm/v1"
	"github.com/eaciit/toolkit"
)

type UserModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string `json:"_id",bson:"_id"`
	LoginID       string
	FullName      string
	Email         string
	Password      string
	Groups        []string
	IsImportant   bool

	Enable    bool
	Grants    []acl.AccessGrant
	LoginType acl.LoginTypeEnum
	LoginConf toolkit.M
}

func (u *UserModel) TableName() string {
	return "acl_users"
}

func (u *UserModel) RecordID() interface{} {
	return u.ID
}

func (u *UserModel) PreSave() error {
	if u.ID == "" && acl.IsUserExist(u.LoginID) {
		return errors.New("acl user is exist")
	}

	if u.ID == "" {
		u.ID = toolkit.RandomString(32)
	}
	return nil
}
