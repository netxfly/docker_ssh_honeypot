package settings

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"path/filepath"
)

var (
	Cfg    *ini.File
	APIURL string
	// Rsyslog      string
	// RsyslogProto string
	KEY  string
	MODE string

	VsftpLog   string
	SshLog     string
	HistLog    string
	RsyncLog   string
	RedisLog   string
	MysqlLog   string
	MongodbLog string
)

func init() {
	log.SetPrefix("[honeypot-sandbox-log]")
	var err error
	source := "conf/app.ini"
	Cfg, err = ini.Load(source)
	//log.Println(Cfg, err)
	if err != nil {
		log.Panicln(err)
	}

	APIURL = Cfg.Section("").Key("api_url").MustString("")
	// Rsyslog = Cfg.Section("").Key("rsyslog_addr").MustString("127.0.0.1:514")
	// RsyslogProto = Cfg.Section("").Key("rsyslog_protocol").MustString("tcp")
	KEY = Cfg.Section("").Key("key").MustString("")
	MODE = Cfg.Section("").Key("mode").MustString("syslog")

	curDir, _ := filepath.Abs("..")

	secLogs := Cfg.Section("logs")
	VsftpLog = secLogs.Key("vsftpd_log").MustString(fmt.Sprintf("%v/logs/vsftpd/syslog", curDir))
	SshLog = secLogs.Key("sshd_log").MustString(fmt.Sprintf("%v/logs/openssh/auth.log", curDir))
	HistLog = secLogs.Key("history_log").MustString(fmt.Sprintf("%v/logs/openssh/syslog", curDir))
	RsyncLog = secLogs.Key("rsync_log").MustString(fmt.Sprintf("%v/logs/rsync/rsyncd.log", curDir))
	MysqlLog = secLogs.Key("mysql_log").MustString(fmt.Sprintf("%v/logs/mysql/mysql.log", curDir))
	RedisLog = secLogs.Key("redis_log").MustString(fmt.Sprintf("%v/logs/redis/syslog", curDir))
	MongodbLog = secLogs.Key("mongodb_log").MustString(fmt.Sprintf("%v/logs/mongodb/syslog", curDir))
}
