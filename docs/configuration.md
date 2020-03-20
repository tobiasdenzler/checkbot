# Configuration

## Prometheus

Use the following snippet to scrape the checkbot metrics:
```
- job_name: checkbot
  scheme: https
  tls_config:
    ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    insecure_skip_verify: true
  static_configs:
    - targets: ['checkbot.checkbot.svc.cluster.local:4444']
```

### Lastrun

To check if your scripts have run successfully you can use the (internal) metric lastrun_info and lastresult_info. These metrics will provide information about the last run and result of each check:

```
checkbot_lastrun_info{interval="60",name="checkbot_missing_quota_on_project_total",offset="22",type="Gauge"} 1.57699768e+09
checkbot_lastrun_info{interval="60",name="checkbot_modified_scc_reconcile",offset="12",type="Gauge"} 1.576997641e+09

checkbot_lastresult_info{interval="60",name="checkbot_missing_quota_on_project_total",offset="22",type="Gauge"} 1
checkbot_lastresult_info{interval="60",name="checkbot_modified_scc_reconcile",offset="12",type="Gauge"} 0
```

Note:  Offset is the number of second that is used to randomly delay the execution of the script. To get the time of the next run you can add the interval and the offset to the current time.
