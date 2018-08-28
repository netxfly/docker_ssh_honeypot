/*

Copyright (c) 2018 sec.lu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package plugins

import (
	"regexp"
	"time"
	"os"
	"net/http"
	"net/url"
	"encoding/json"

	"cloud-honeypot/sandbox_log/models"
	"cloud-honeypot/sandbox_log/settings"
	"cloud-honeypot/sandbox_log/misc"
	"cloud-honeypot/sandbox_log/logger"
)

var (
	RegexpSSH *regexp.Regexp
	err       error
)

func init() {
	RegexpSSH, err = regexp.Compile(`Honeypot: Username: (.+?) Password: (.+?), from: (.+?), result: (.+?)`)
}

func Check(tag, content string) (checkResult models.CheckResult, result bool) {
	switch tag {
	case "openssh-server":
		checkResult, result = CheckSSH(content, tag)
	}
	if result {
		logger.Log.Infof("from ip: %v, user: %v, password: %v, result: %v", checkResult.Ip,
			checkResult.Username, checkResult.Password, checkResult.Status)
	}
	return checkResult, result
}

func SaveResult(checkResult models.CheckResult, result bool) {
	if result {
		content, _ := json.Marshal(checkResult)
		t := time.Now().Format("2006-01-02 15:04:05")
		apiData := models.APIDATA{}
		apiData.Tag = checkResult.Tag
		apiData.Content = string(content)
		apiData.Hostname, _ = os.Hostname()
		data, err := json.Marshal(apiData)
		if err == nil {
			_, err = http.PostForm(settings.APIURL, url.Values{"timestamp": {t},
				"secureKey": {misc.MakeSign(t, settings.KEY)}, "data": {string(data)}})
		}
	}
}
