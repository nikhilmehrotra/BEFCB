package controllers

import (
	"crypto/md5"
	"eaciit/scb-apps/webapp/apps/main/models"
	"eaciit/scb-apps/webapp/helper"
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"io"
	"sort"
	"strings"
)

func (c *AccessController) GetUser(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	csr, err := acl.Find(new(models.UserModel), nil, nil)
	defer csr.Close()
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	data := make([]models.UserModel, 0)
	err = csr.Fetch(&data, 0, false)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(data)
}

func (c *AccessController) SelectUser(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := new(models.UserModel)
	err := k.GetPayload(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	err = acl.FindByID(payload, payload.ID)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(payload)
}

func (c *AccessController) DeleteUser(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := new(models.UserModel)
	err := k.GetPayload(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	if payload.IsImportant {
		return c.SetResultError("Cannot delete this user", nil)
	}

	err = acl.Delete(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	return c.SetResultOK(nil)
}

func (c *AccessController) SaveUser(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := new(models.UserModel)
	err := k.GetPayload(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	if payload.IsImportant {
		return c.SetResultError("Cannot update this user", nil)
	}

	payload.Password = strings.TrimSpace(payload.Password)
	if payload.ID == "" {
		payload.ID = bson.NewObjectId().Hex()
	}

	err = acl.Save(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	if payload.ID == "" {
		payload.Enable = true
		err = acl.ChangePassword(payload.ID, payload.Password)
		if err != nil {
			return c.SetResultError(err.Error(), nil)
		}
	} else {
		// get old user info
		oldUserInfo := new(models.UserModel)
		err = acl.FindByID(oldUserInfo, payload.ID)
		if err != nil {
			return c.SetResultError(err.Error(), nil)
		}

		// generate hashed password
		hasher := md5.New()
		io.WriteString(hasher, payload.Password)
		hashedPassword := tk.Sprintf("%x", hasher.Sum(nil))

		// update password if not match with previous password
		if (payload.Password != "") && (payload.Password != oldUserInfo.Password) && (hashedPassword != oldUserInfo.Password) {
			err = acl.ChangePassword(payload.ID, payload.Password)
			if err != nil {
				return c.SetResultError(err.Error(), nil)
			}
		}
	}

	return c.SetResultOK(payload)
}

// ==== user log

type Sessions []acl.Session

func (a Sessions) Len() int           { return len(a) }
func (a Sessions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Sessions) Less(i, j int) bool { return helper.IsTimeAfter(a[i].Created, a[j].Created) }

func (c *AccessController) GetUserLog(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	csr, err := acl.Find(new(acl.Session), nil, nil)
	defer csr.Close()
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	data := make([]acl.Session, 0)
	err = csr.Fetch(&data, 0, false)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	max := 100
	if len(data) < max {
		max = len(data)
	}

	sort.Sort(Sessions(data))
	return c.SetResultOK(data[0:max])
}
