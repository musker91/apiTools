package utils

import "testing"

func TestCheckProtocolHttp(t *testing.T) {
	bool := CheckProtocolHttp("218.26.178.22:8080", "www.baidu.com")
	if !bool {
		t.Error("fail")
		return
	}
	t.Log("success")
}

func TestGetRedirectUrl(t *testing.T) {
	url := "https://aweme.snssdk.com/aweme/v1/play/?s_vid=93f1b41336a8b7a442dbf1c29c6bbc568f42debb6d4be55bfbe1731ca905dc268a68662e85f9efb793d4f61489c94c8183b223f772bda1702ee2c3cd1cf4f292&line=0"

	header := map[string]string{"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) CriOS/56.0.2924.75 Mobile/14E5239e Safari/602.1"}

	redirectUrl, err := GetRedirectUrl(url, "", header)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(redirectUrl)
}
