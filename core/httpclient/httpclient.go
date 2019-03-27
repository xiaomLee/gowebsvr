package httpclient

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type HttpClientContext struct {
	ConnTimeoutMs int64
	ReadTimeoutMs int64
	Url           string
	Method        string
	Param         map[string]string
	ArrayParam    map[string][]string
	Header        map[string]string
	TryTimes      int
	Body          string
	requestParams string
	client        *http.Client
}

func (ctx *HttpClientContext) String() string {
	return fmt.Sprintf("url=%v||request=%v||body=%v||header=%v",
		ctx.Url,
		ctx.requestParams,
		ctx.Body,
		ctx.Header,
	)
}

type HttpClientResponse struct {
	Begin    time.Time
	HttpCode int
	Latency  time.Duration
	RespBody []byte
}

func (resp *HttpClientResponse) String() string {
	return fmt.Sprintf("cspanid=%v||HttpCode=%v||response=%v||proc_time=%.2f",
		resp.Begin.UnixNano(),
		resp.HttpCode,
		strings.Replace(string(resp.RespBody), "\n", ";", -1),
		resp.Latency.Seconds()*1000,
	)
}

func NewHttpClientContext() (ctx *HttpClientContext) {
	ctx = &HttpClientContext{
		ConnTimeoutMs: 50,
		ReadTimeoutMs: 200,
		Param:         make(map[string]string),
		ArrayParam:    make(map[string][]string),
		Header:        make(map[string]string),
	}
	//add common headers
	ctx.Header["request-id"] = strconv.FormatInt(time.Now().Unix(), 10)
	return
}

func (ctx *HttpClientContext) NewHTTPTimeoutClient() {
	if ctx.client != nil {
		return
	}
	connTimeout := time.Duration(time.Duration(ctx.ConnTimeoutMs) * time.Millisecond)
	readTimeout := time.Duration(time.Duration(ctx.ReadTimeoutMs) * time.Millisecond)
	ctx.client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, connTimeout)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(readTimeout))
				return conn, nil
			},
			ResponseHeaderTimeout: readTimeout,
		},
	}
	return
}

func (ctx *HttpClientContext) NewHTTPBasicAuthClient(username, password string) {
	if ctx.client != nil {
		return
	}
	ctx.client = &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				req.SetBasicAuth(username, password)
				return nil, nil
			},
		},
	}
	return
}

func (ctx *HttpClientContext) DoHttpCall() (httpResp *HttpClientResponse, err error) {

	var request *http.Request
	var resp *http.Response

	httpResp = &HttpClientResponse{
		Begin: time.Now(),
	}

	ctx.Header["Didi-Header-Spanid"] = fmt.Sprintf("%v", httpResp.Begin.UnixNano())
	defer func() {
		ctx.TryTimes++
		httpResp.Latency = time.Since(httpResp.Begin)
		if err != nil {
			fmt.Println("http_failure", "ctx", ctx, "httpResp", httpResp, "err", err)
		} else {
			fmt.Println("http_success", "ctx", ctx, "httpResp", httpResp, "err", err)
		}
	}()

	ctx.NewHTTPTimeoutClient()
	if ctx.client == nil {
		err = fmt.Errorf("OOM:NewHTTPTimeoutClient()")
		return
	}

	var reader io.Reader = nil
	//if ctx.Method != "GET" && ctx.Method != "HEAD" {
	if len(ctx.Body) > 0 {
		reader = strings.NewReader(ctx.Body)
	} else {
		values := url.Values{}
		for k, v := range ctx.Param {
			values.Set(k, v)
		}
		for k, l := range ctx.ArrayParam {
			for _, v := range l {
				values.Add(k, v)
			}
		}
		ctx.requestParams = values.Encode()
		reader = strings.NewReader(ctx.requestParams)
	}
	//}
	if ctx.Method == "GET" && len(ctx.requestParams) > 0 {
		ctx.Url = ctx.Url + "?" + ctx.requestParams
		reader = nil
	}

	request, err = http.NewRequest(ctx.Method, ctx.Url, reader)
	if err != nil {
		fmt.Println("http_failure", "ctx", ctx, "http.NewRequest(),err", err)
		return
	}
	for k, v := range ctx.Header {
		request.Header.Set(k, v)
	}

	resp, err = ctx.client.Do(request)
	if err != nil {
		fmt.Println("http_failure", "ctx", ctx, "client.Do(),err", err)
		return
	}

	httpResp.HttpCode = resp.StatusCode
	httpResp.RespBody, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}
