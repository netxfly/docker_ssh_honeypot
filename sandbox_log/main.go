package main

import (
	"cloud-honeypot/sandbox_log/settings"
	"cloud-honeypot/sandbox_log/logger"
	"cloud-honeypot/sandbox_log/util"

	"time"
	"os"
)

func main() {
	logs := make(map[string]string)
	logs["vsftpd-server"] = settings.VsftpLog
	logs["openssh-server"] = settings.SshLog
	logs["history"] = settings.HistLog
	logs["rsync"] = settings.RsyncLog
	logs["mysql-server"] = settings.MysqlLog
	logs["redis-server"] = settings.RedisLog
	logs["mongodb-server"] = settings.MongodbLog

	logger.Log.Infoln("Log config info:", settings.APIURL, settings.KEY, settings.MODE)

	for tag, fileName := range logs {
		logger.Log.Printf("Tag: %v, path: %v, Truncate: %v", tag, fileName, os.Truncate(fileName, 0))

		if settings.MODE == "syslog" {
			// go util.MonitorLog(tag, fileName)
		} else {
			go util.CheckLog(tag, fileName)
		}
	}

	for {
		time.Sleep(3 * time.Second)
	}
}
