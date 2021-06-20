package do

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/LxAday/goerror"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

//hosts文件路径
const (
	//hosts文件路径
	hostsPath = `C:\Windows\System32\drivers\etc\hosts`
	//ip正则
	regIp = `((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}`
)

//查询地址
var url = map[string]string{
	"github.com":                   "https://github.com.ipaddress.com",
	"github.global.ssl.fastly.net": "https://fastly.net.ipaddress.com/github.global.ssl.fastly.net",
}

type Do struct {
	req      *http.Client
	writeMap map[string]string
}

func New() *Do {
	return &Do{
		req: &http.Client{
			Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
			Timeout:   time.Second * 30,
		},
		writeMap: map[string]string{},
	}
}

//Run 运行
func (d *Do) Run() error {
	//先清空已有配置
	if e := d.clearConf(); e != nil {
		return goerror.Wrap(e, "Run")
	}
	for k, v := range url {
		if e := d.query(v, k); e != nil {
			fmt.Println(goerror.Wrap(e, "Run2"))
		}
	}
	if len(d.writeMap) > 0 {
		return d.write()
	}
	return nil
}

//write 修改hosts
func (d *Do) write() error {
	file, e := os.ReadFile(hostsPath)
	if e != nil {
		return goerror.Wrap(e, "write")
	}
	contentArr := strings.Split(string(file), "\n")
	var newContentArr []string
look:
	for _, v := range contentArr {
		if !strings.HasPrefix(v, "#") {
			for key := range d.writeMap {
				if strings.Contains(v, key) {
					continue look
				}
			}
		}
		newContentArr = append(newContentArr, v)
	}
	for k, v := range d.writeMap {
		newContentArr = append(newContentArr, fmt.Sprintf(`%s	%s`, v, k))
	}
	return os.WriteFile(hostsPath, []byte(strings.Join(newContentArr, "\n")), os.ModePerm)
}

//query 请求查询
//url 请求地址
//host 配置对应域名
func (d *Do) query(url, host string) error {
	b, e := d.ping(host)
	if e != nil {
		return goerror.Wrap(e, "query")
	}
	if !b {
		if e = d.request(url, host); e != nil {
			return goerror.Wrap(e, "query1")
		}
	}
	return nil
}

//clearConf 清除已有配置
func (d *Do) clearConf() error {
	d.writeMap = map[string]string{} //清空配置数据
	return d.write()
}

//ping 通过ping获取最近cdn ip
func (d *Do) ping(host string) (bool, error) {
	output, _ := exec.Command("cmd", "/C", fmt.Sprintf("ping -t %s -a -n 1", host)).Output()
	str := string(bytes.TrimSpace(output))
	if str == "" {
		return false, nil
	}
	//正则匹配ip
	ip := regexp.MustCompile(regIp).FindString(str)
	if ip == "" {
		return false, nil
	}
	d.writeMap[host] = ip
	return true, nil
}

//request 请求查询ip
func (d *Do) request(url, host string) error {
	request, e := http.NewRequest("GET", url, nil)
	if e != nil {
		return goerror.Wrap(e, "request")
	}
	request.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	request.Header.Add("cache-control", "no-cache")
	request.Header.Add("pragma", "no-cache")
	request.Header.Add("referer", "https://www.ipaddress.com/")
	request.Header.Add("upgrade-insecure-requests", "1")
	request.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	do, e := d.req.Do(request) //请求
	if e != nil {
		return goerror.Wrap(e, "request1")
	}
	defer do.Body.Close()
	buffer := &bytes.Buffer{}
	if _, e = io.Copy(buffer, do.Body); e != nil {
		return goerror.Wrap(e, "request2")
	}
	ip := regexp.MustCompile(`<tr><th>IP Address</th><td><ul class="comma-separated"><li>` + regIp + `</li></ul>`).FindString(buffer.String())
	if ip != "" {
		//正则匹配ip
		ip = regexp.MustCompile(regIp).FindString(ip)
		d.writeMap[host] = ip
	}
	return nil
}
