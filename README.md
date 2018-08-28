## 概述

实现一个高交互的SSH蜜罐有方式有多种：

- 可以使用现成的开源系统，如`kippo`、`cowrie`等，
- 可以利用一些SSH库做2次开发，如`https://github.com/gliderlabs/ssh`
- 也可以用Docker定制

市面上已经有这么多开源的SSH蜜罐了，我为什么还要再造个轮子呢，理由如下：

- 部署这些系统要安装一堆的依赖，运维部署成本较高
- 想要的功能没法添加，不想要的功能也没法删减
- 如果要支持多种高交互的服务，必须用一堆开源的系统拼凑出一堆服务来，每种系统的后端DB不同，数据结构也各不相同，无法统一处理
- 自研的系统部署简单，服务可自由扩展，数据格式可自由定制和统一，方便运维与运营

## 技术架构

笔者在之前的文章[《自制攻击欺骗防御系统》](https://zhuanlan.zhihu.com/p/23535920)中介绍过完整的架构，整个系统由以下几个模块组成：

- Agent端，部署于服务器中的Agent，用于实时获取用户的访问日志并传递到检测端Server中，如果是恶意攻击，则会将流量重定向到沙盒中。目前支持的服务有：

    - WEB
    - FTP
    - SSH
    - Rsync
    - Mysql
    - Redis
    - Mongodb

- Server端，攻击检测服务器，实时检测Agent传递过来的日志并判断是否为攻击者，并为Agent动态、实时地维护了一份攻击者的来源IP策略
- Mamager端，策略管理服务器，有为Agent和server提供策略、攻击log统计、查看的功能
- 高交互蜜罐系统及守护进程，高交互蜜罐系统由docker封装的一些服务组成，守护进程负责把仿真系统中产生的LOG的数据格式化后再传给Server端进行攻击检测与入库

本文只讲SSH高交互蜜罐的实现。

### SSH高交互蜜罐的Docker实现

docker中默认的openssh服务没有记录bash命令及ssh密码的功能，需要修改bash及openssh的源码并重新编译到Docker中，以下为`Dockerfile`的内容：
```bash
FROM debian:jessie

RUN export DEBIAN_FRONTEND='noninteractive' && \
    apt-get update -qq && \
    apt-get install -qqy --no-install-recommends openssh-server rsyslog wget patch make gcc curl libc6-dev net-tools vim && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* && \
    mkdir -p /softs/bash-4.3.30-active-syslog/ && mkdir -p /home/test

COPY softs/ /softs/

RUN chmod +x /softs/bash-4.3.30-active-syslog/install-bash-syslog.sh
RUN chmod +x /softs/install_ssh.sh

WORKDIR /softs/bash-4.3.30-active-syslog/
RUN ["/bin/bash", "./install-bash-syslog.sh"]

WORKDIR /softs/
RUN ["/bin/bash", "./install_ssh.sh"]

EXPOSE 22/tcp

COPY entrypoint.sh /etc/systemd/system/

COPY Shanghai /etc/localtime

RUN chmod +x /etc/systemd/system/entrypoint.sh

ENTRYPOINT ["/etc/systemd/system/entrypoint.sh"]
```

- `install-bash-syslog.sh`为记录详细执行命令记录的bash安装脚本；
- `/install_ssh.sh`为记录openssh密码的安装脚本。

完成后的Docker镜像封装脚本的地址为[https://github.com/netxfly/docker_ssh_honeypot/blob/master/install.sh](https://github.com/netxfly/docker_ssh_honeypot/blob/master/install.sh)。

执行sh -x ./install.sh后经过一段时间的等待后，会直接封装一个Openssh镜像并启动。

### 测试截图

openssh的密码破解记录如下：

![](http://docs.xsec.io/images/docker_ssh_honeypot/openssh.jpg)

bash的命令执行记录如下：

![](http://docs.xsec.io/images/docker_ssh_honeypot/bash.jpg)

如果允许攻击者登录到蜜罐中，还需要修改蜜罐的IP地址和主机名，防止攻击者一登录成功就识别出来；其次要做好ACL，防止攻击者通过蜜罐作为跳板攻击内网中其他系统。


### 仿真系统的守护进程实现

仿真系统的守护进程的作用是实时监控Docker服务产生的数据，并格式化后发送到server端中进行检测，通过插件的形式，可以支持多种服务的日志格式化，然后将格式化的日志通过http协议发送到server端。

```go
func main() {
	logs := make(map[string]string)
	logs["vsftpd-server"] = settings.VsftpLog
	logs["openssh-server"] = settings.SshLog
	logs["history"] = settings.HistLog
	logs["rsync"] = settings.RsyncLog
	logs["mysql-server"] = settings.MysqlLog
	logs["redis-server"] = settings.RedisLog
	logs["mongodb-server"] = settings.MongodbLog

	// logger.Log.Infoln("Log config info:", settings.APIURL, settings.KEY, settings.MODE)

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
```

SSH密码破解检测插件的实现如下：

1. 通过正则匹配中破解密码的记录并格式化
2. 确定为破解行为的记录，通过http接口传到server中

```go

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

```

server端接到数据后做相应的处理后入库，并在管理端中展示，例如：

![](http://docs.xsec.io/images/docker_ssh_honeypot/data.jpg)

