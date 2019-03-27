package httpclient

import (
	"encoding/json"
	"testing"
)

func TestHttpclient(t *testing.T) {
	httpctx := NewHttpClientContext()
	httpctx.Url = "http://127.0.0.1:8080/admin"
	httpctx.ConnTimeoutMs, httpctx.ReadTimeoutMs = 200, 200
	httpctx.Method = "POST"
	reqInfo := map[string]interface{}{
		"value": "123",
		"user":  "limingzhong",
	}

	b, _ := json.Marshal(reqInfo)
	httpctx.Body = string(b)
	httpctx.Header["Content-Type"] = "application/json"

	var httpResp *HttpClientResponse
	var err error
	if httpResp, err = httpctx.DoHttpCall(); err != nil {
		t.Errorf("请求失败||err=%v", err)
	}
	t.Logf("请求成功||resp.Bosy=%v", string(httpResp.RespBody))
}
