package util

import (
	"log"
	"github.com/hpcloud/tail"
)

func MonitorLog(tag, logName string) {
	t, err := tail.TailFile(logName, tail.Config{Follow: true})
	// fmt.Println(t, err)
	if err == nil {
		for line := range t.Lines {
			//log.Println(settings.RsyslogProto, settings.Rsyslog, syslog.LOG_ERR, tag)
			//l3, err := syslog.Dial(settings.RsyslogProto, settings.Rsyslog, syslog.LOG_ERR, tag)
			// l3, err := syslog.New(syslog.LOG_ERR, tag)
			//if err != nil {
			//	log.Fatal(err)
			//}
			//log.SetOutput(l3)
			log.Println(line.Text)
		}
	}
}
