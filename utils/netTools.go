package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

// http proxy get方法
func HttpProxyGet(dataUrl, proxyIp string, headers map[string]string) (data []byte, response *http.Response, err error) {
	transport := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, //ssl证书报错问题
		DisableKeepAlives: false,                                 //关闭连接复用，因为后台连接过多最后会造成端口耗尽
		MaxIdleConns:      100,                                   //最大空闲连接数量
		IdleConnTimeout:   time.Duration(5 * time.Second),        //空闲连接超时时间
	}
	if proxyIp != "" { // 设置代理
		proxyUrl, _ := url.Parse("http://" + proxyIp)
		transport.Proxy = http.ProxyURL(proxyUrl)
	}
	client := &http.Client{
		Timeout:   time.Duration(30 * time.Second),
		Transport: transport,
	}

	request, err := http.NewRequest("GET", dataUrl, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.70 Safari/537.36")
	if headers != nil {
		for key, value := range headers {
			request.Header.Set(key, value)
		}
	}
	if err != nil {
		return
	}
	// 请求数据
	resp, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("request %s, proxyIp: (%s),err: %v", dataUrl, proxyIp, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("request %s fail status not is 200 ok, status: %s", dataUrl, resp.Status)
		return
	}
	// 读取数据
	buf := make([]byte, 128)
	data = make([]byte, 0, 2048)
	for {
		n, err := resp.Body.Read(buf)
		data = append(data, buf[:n]...)

		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
	}
	response = resp
	return
}

// http proxy post方法
func HttpProxyPost(dataUrl string, reqData interface{}, proxyIp string, headers map[string]string) (data []byte, response *http.Response, err error) {
	transport := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, //ssl证书报错问题
		DisableKeepAlives: false,                                 //关闭连接复用，因为后台连接过多最后会造成端口耗尽
		MaxIdleConns:      100,                                   //最大空闲连接数量
		IdleConnTimeout:   time.Duration(5 * time.Second),        //空闲连接超时时间
	}
	if proxyIp != "" { // 设置代理
		proxyUrl, _ := url.Parse("http://" + proxyIp)
		transport.Proxy = http.ProxyURL(proxyUrl)
	}
	client := &http.Client{
		Timeout:   time.Duration(30 * time.Second),
		Transport: transport,
	}
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST", dataUrl, bytes.NewReader(reqBody))
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.70 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
	if headers != nil {
		for key, value := range headers {
			request.Header.Set(key, value)
		}
	}
	if err != nil {
		return
	}
	// 请求数据
	resp, err := client.Do(request)
	if err != nil {
		err = fmt.Errorf("request %s, proxyIp: (%s),err: %v", dataUrl, proxyIp, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("request %s fail status not is 200 ok, status: %s", dataUrl, resp.Status)
		return
	}
	// 读取数据
	buf := make([]byte, 128)
	data = make([]byte, 0, 2048)
	for {
		n, err := resp.Body.Read(buf)
		data = append(data, buf[:n]...)

		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
	}
	response = resp
	return
}

// 校验代理 http和http协议协议代理
func CheckProtocolHttp(proxyAddr, checkToUrl string) bool {
	httpClient := &http.Client{
		Timeout: time.Duration(10 * time.Second), //客户端设置10秒超时
	}
	httpClient.Transport = &http.Transport{
		DisableKeepAlives: false,                          //关闭连接复用，因为后台连接过多最后会造成端口耗尽
		MaxIdleConns:      100,                            //最大空闲连接数量
		IdleConnTimeout:   time.Duration(5 * time.Second), //空闲连接超时时间
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   proxyAddr,
		}),                                                //设置http代理地址
	}
	resp, err := httpClient.Get(fmt.Sprintf("http://%s", checkToUrl))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

// 获取302重定向后的url地址
func GetRedirectUrl(sourceUrl string, proxyIp string, headers map[string]string) (redirectUrl string, response *http.Response, err error) {
	transport := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true}, //ssl证书报错问题
		DisableKeepAlives: false,                                 //关闭连接复用，因为后台连接过多最后会造成端口耗尽
		MaxIdleConns:      100,                                   //最大空闲连接数量
		IdleConnTimeout:   time.Duration(5 * time.Second),        //空闲连接超时时间
	}

	if proxyIp != "" { // 设置代理
		proxyUrl, _ := url.Parse("http://" + proxyIp)
		transport.Proxy = http.ProxyURL(proxyUrl)
	}

	client := &http.Client{
		Timeout:   time.Duration(time.Second * 30),
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	request, err := http.NewRequest("GET", sourceUrl, nil)
	if err != nil {
		return
	}

	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.70 Safari/537.36")

	if headers != nil {
		for key := range headers {
			request.Header.Set(key, headers[key])
		}
	}

	rsq, err := client.Do(request)
	if err != nil {
		return
	}

	response = rsq

	statusCode := rsq.StatusCode
	if statusCode == 302 || statusCode == 301 {
		redirectUrl = rsq.Header.Get("Location")
		return
	} else {
		err = errors.New("request status not redirect")
		return
	}
}
