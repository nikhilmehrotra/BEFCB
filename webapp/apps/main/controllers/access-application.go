package controllers

import (
	"eaciit/scb-apps/webapp/apps/main/models"
	"eaciit/scb-apps/webapp/helper"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"os"
	"path/filepath"
	"strings"
)

func (c *AccessController) GetApplication(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	appsMap := tk.M{}

	type ForgetMe struct{}
	basePath := filepath.Join(helper.GetAppBasePath(ForgetMe{}), "..")
	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		frontPath := strings.Trim(strings.Replace(path, basePath, "", -1), "/ ")
		frontPathFolderName := strings.Split(frontPath, "/")[0]
		if frontPathFolderName == "main" {
			return nil
		}
		if strings.TrimSpace(frontPathFolderName) == "" {
			return nil
		}
		if strings.HasPrefix(frontPathFolderName, ".") {
			return nil
		}

		if !appsMap.Has(frontPathFolderName) {
			appsMap.Set(frontPathFolderName, true)
		}
		return nil
	})
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	// ====== get saved data

	csr, err := c.Ctx.Find(new(models.ApplicationModel), nil)
	defer csr.Close()
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	data := make([]models.ApplicationModel, 0)
	err = csr.Fetch(&data, 0, false)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	for appName := range appsMap {
		isExists := false
		for _, each := range data {
			if each.ID == appName && !isExists {
				isExists = true
			}
		}

		if !isExists {
			row := models.ApplicationModel{ID: appName}
			data = append(data, row)
		}
	}

	return c.SetResultOK(data)
}

func (c *AccessController) SaveApplication(k *knot.WebContext) interface{} {
	c.SetResponseTypeAJAX(k)

	payload := new([]models.ApplicationModel)
	err := k.GetPayload(payload)
	if err != nil {
		return c.SetResultError(err.Error(), nil)
	}

	for _, each := range *payload {
		if each.ID == "" {
			each.ID = bson.NewObjectId().Hex()
		}

		err = c.Ctx.Connection.
			NewQuery().
			From(new(models.ApplicationModel).TableName()).
			Where(dbox.Eq("_id", each.ID)).
			Save().
			Exec(tk.M{"data": each})
		if err != nil {
			return c.SetResultError(err.Error(), nil)
		}
	}

	return c.SetResultOK(nil)
}
