package util

import (
	"github.com/hpcloud/tail"
	"cloud-honeypot/sandbox_log/util/plugins"
)

// Send log to remote server
func CheckLog(tag, logName string) (err error) {
	t, err := tail.TailFile(logName, tail.Config{Follow: true})
	if err == nil {
		for line := range t.Lines {
			plugins.SaveResult(plugins.Check(tag, line.Text))
		}
	}
	return err
}
