package models

import "time"

type MetricFile struct {
	FileName         string    `bson:"filename",json:"filename"`
	OriginalFileName string    `bson:"originalfilename",json:"originalfilename"`
	MonthYear        int       `bson:"monthyear",json:"monthyear"`
	UploadedAt       time.Time `bson:"uploadedat",json:"uploadedat"`
	UploaderName     string    `bson:"uploadername",json:"uploadername"`
}

func NewMetricFile() *MetricFile {
	return new(MetricFile)
}
