package helper

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"crypto/md5"
	"encoding/hex"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/csv"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"gopkg.in/gomail.v2"
)

var (
	DebugMode bool
)

var config_system = func() string {
	d, _ := os.Getwd()
	d += "/conf/confsystem.json"
	return d
}()

func GetPathConfig() (result map[string]interface{}) {
	result = make(map[string]interface{})

	ci := &dbox.ConnectionInfo{config_system, "", "", "", nil}
	conn, e := dbox.NewConnection("json", ci)
	if e != nil {
		return
	}

	e = conn.Connect()
	defer conn.Close()
	csr, e := conn.NewQuery().Select("*").Cursor(nil)
	if e != nil {
		return
	}
	defer csr.Close()
	data := []toolkit.M{}
	e = csr.Fetch(&data, 0, false)
	if e != nil {
		return
	}
	result["folder-path"] = data[0].GetString("folder-path")
	result["restore-path"] = data[0].GetString("restore-path")
	result["folder-img"] = data[0].GetString("folder-img")
	return
}

func CreateResult(success bool, data interface{}, message string) map[string]interface{} {
	if !success {
		fmt.Println("ERROR! ", message)
		if DebugMode {
			panic(message)
		}
	}

	return map[string]interface{}{
		"data":    data,
		"success": success,
		"message": message,
	}
}

func UploadHandler(r *knot.WebContext, filename, dstpath string, existCount int) (error, string) {
	file, handler, err := r.Request.FormFile(filename)
	if err != nil {
		return err, ""
	}
	defer file.Close()

	dstSource := ""
	uploadfilename := ""
	filecounted := ""
	if existCount > 0 {
		filecounted = toolkit.ToString(existCount) + " - " + handler.Filename
	} else {
		filecounted = handler.Filename
	}
	dstSource = dstpath + toolkit.PathSeparator + filecounted
	uploadfilename = filecounted

	f, err := os.OpenFile(dstSource, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err, ""
	}
	defer f.Close()
	io.Copy(f, file)

	return nil, uploadfilename
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func OcirSendEmail(_to, _mailmsg string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "admin.support@eaciit.com")
	m.SetHeader("To", _to)

	m.SetHeader("Subject", "[no-reply] Self password reset")
	m.SetBody("text/html", _mailmsg)

	d := gomail.NewPlainDialer("smtp.office365.com", 587, "admin.support@eaciit.com", "B920Support")
	err := d.DialAndSend(m)

	return err
}

//Reference,Time,Description,SessionId,LoginId,RequestAddr,LoadingTimes,_id,Action
func PrepareConnectionLogFile(sfile string) (dbox.IConnection, error) {
	var config = map[string]interface{}{"useheader": true,
		"delimiter": ",",
		"newfile":   true,
		"mapheader": []toolkit.M{toolkit.M{}.Set("SessionId", "string"),
			toolkit.M{}.Set("LoginId", "string"),
			toolkit.M{}.Set("Action", "string"),
			toolkit.M{}.Set("Reference", "string"),
			toolkit.M{}.Set("RequestAddr", "string"),
			toolkit.M{}.Set("Time", "string"),
			toolkit.M{}.Set("LoadingTimes", "string"),
			toolkit.M{}.Set("Description", "string"),
		}}
	ci := &dbox.ConnectionInfo{sfile, "", "", "", config}
	c, e := dbox.NewConnection("csv", ci)
	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}

	return c, nil
}

func ReadConfig() map[string]string {
	wd, _ := os.Getwd()
	ret := make(map[string]string)
	file, err := os.Open(filepath.Join(wd, "bef/conf/app.conf"))
	if err == nil {
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}

			sval := strings.Split(string(line), "=")
			if len(sval) > 2 {
				var tempval = ""
				for xi, x := range sval {
					if xi != 0 {
						if xi == 1 {
							tempval += x
						} else {
							tempval += "=" + x
						}
					}
				}
				ret[sval[0]] = tempval
			} else if len(sval) > 0 {
				ret[sval[0]] = sval[1]
			}
		}
	} else {
		log.Println(err.Error())
	}

	return ret
}

func PrepareConnection() (dbox.IConnection, error) {
	config := ReadConfig()
	ci := &dbox.ConnectionInfo{config["host"], config["database"], config["username"], config["password"], toolkit.M{}.Set("timeout", 10)}
	c, e := dbox.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}

	return c, nil
}
