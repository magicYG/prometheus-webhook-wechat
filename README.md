# prometheus-webhook-wechat

[![Go Report Card](https://goreportcard.com/badge/github.com/SecurityNeo/prometheus-webhook-wechat)](https://goreportcard.com/badge/github.com/SecurityNeo/prometheus-webhook-wechat)

## 使用
```
usage: prometheus-webhook-wechat [<flags>]


Flags:
  -h, --[no-]help               Show context-sensitive help (also try --help-long and --help-man).
      --listen-port=":80"       The port to listen for API interface,default: :80
      --config.file=config.yml  The full path of config file,default: config.yml
      --template.file=template.tmpl
                                The full path of template file,default: template.tmpl
      --log.level=info          Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt       Output format of log messages. One of: [logfmt, json]
      --[no-]version            Show application version.

```

## 配置
```yaml
targets:
  webhook1:
    url: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxx
  webhook2:
    url: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxYY
```

## 告警模版
```
{{ define "__subject" }}
[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}]
{{ end }}


{{ define "__alert_list" }}{{ range . }}
---
>**告警主题**: {{ .Annotations.summary }}
>**告警类型**: {{ .Labels.alertname }}
>**告警级别**: {{ .Labels.severity }}
>**告警实例**: {{ .Labels.instance }}
>**告警内容**: {{ index .Annotations "description" }}
>**告警来源**: [{{ .GeneratorURL }}]({{ .GeneratorURL }})
>**告警时间**: {{ .StartsAt.Local.Format "2006-01-02 15:04:05" }}
{{ end }}{{ end }}

{{ define "__resolved_list" }}{{ range . }}
---
>**告警主题**: {{ .Annotations.summary }}
>**告警类型**: {{ .Labels.alertname }}
>**告警级别**: {{ .Labels.severity }}
>**告警实例**: {{ .Labels.instance }}
>**告警内容**: {{ index .Annotations "description" }}
>**告警来源**: [{{ .GeneratorURL }}]({{ .GeneratorURL }})
>**告警时间**: {{ .StartsAt.Local.Format "2006-01-02 15:04:05" }}
>**恢复时间**: {{ .EndsAt.Local.Format "2006-01-02 15:04:05" }}
{{ end }}{{ end }}


{{ define "msg.title" }}
{{ template "__subject" . }}
{{ end }}

{{ define "msg.content" }}
{{ if gt (len .Alerts.Firing) 0 }}
**产生<font color="warning">{{ .Alerts.Firing | len  }}</font>个故障**
{{ template "__alert_list" .Alerts.Firing }}
{{ end }}

{{ if gt (len .Alerts.Resolved) 0 }}
**恢复<font color="info">{{ .Alerts.Resolved | len  }}</font>个故障**
{{ template "__resolved_list" .Alerts.Resolved }}
{{ end }}
{{ end }}

{{ template "msg.title" . }}
{{ template "msg.content" . }}
```
***注意***：
1. 当前告警内容仅支持markdown
2. 制作模版时需参考企业微信说明文档，查看其当前能支持哪些语法，[机器人配置说明](https://developer.work.weixin.qq.com/document/path/91770)
