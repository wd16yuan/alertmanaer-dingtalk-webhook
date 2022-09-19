## Alertmanager Dingtalk Webhook

A simple webhook  service support send Prometheus alert message to Dingtalk.

## How To Use

```
# cd alertmanaer-dingtalk-webhook
# go build
# webhook -h
  -convert
        token or secert convert to ciphertext
  -defaultUrl string
        global dingtalk robot webhook (default "https://oapi.dingtalk.com/robot/send")
  -h    help
  -secret string
        dingtalk robot webhook secret
  -token string
        dingtalk robot webhook token
 
# webhook -convert -token=xxxx -secret=xxx    // get token secret ciphertext
# webhook -token="ciphertext format" -secret="ciphertext format"  // add the above ciphertext to start the service
```

Or you can overwrite by add annotations to Prometheus alertrule to special the dingtalk webhook for each alert rule.

```
groups:
- name: hostStatsAlert
  rules:
  - alert: hostCpuUsageAlert
    expr: sum(avg without (cpu)(irate(node_cpu{mode!='idle'}[5m]))) by (instance) > 0.85
    for: 1m
    labels:
      severity: page
    annotations:
      summary: "Instance {{ $labels.instance }} CPU usgae high"
      description: "{{ $labels.instance }} CPU usage above 85% (current value: {{ $value }})"
      dingtalkUrl: "https://oapi.dingtalk.com/robot/send"
      dingtalkToken: "ciphertext format"
      dingtalkSecret: "ciphertext format"
```

Build Docker image

```shell
# docker build -t dingtalk-Webhook:v1 .
# docker run --rm -p 8002:8080 dingtalk-Webhook:v1
# curl http://127.0.0.1:8002
Welcome to dingtalk alarm sending api!
```

