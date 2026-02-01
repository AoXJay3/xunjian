package main

import (
	"net/http"
	"net/url"
	"strings"
)

// 访问登录页，让服务端下发JSESSIONID，使用http.NewRequest的http.MethodGet的MethodGet
func visitLoginPage(client *http.Client) error {
	req, _ := http.NewRequest(
		http.MethodGet,
		baseURL+"/goLogin.do",
		nil,
	)

	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Referer", baseURL+"/admin/goAdmin.do") // 设置请求来源，模拟服务器访问，告诉服务端这是从/goAdmin.do跳转过来的
	resp, err := client.Do(req)
	// 如果有错误，返回错误内容，否则返回nil
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// 看Cookie fmt.Println(resp.Cookies())
	// for _, c := range resp.Cookies() {
	// 	fmt.Printf("cookies: %s=%s\n", c.Name, c.Value)
	// }

	return nil
}

// 登录，使用http.NewRequest的http.MethodPost
func login(client *http.Client) error {
	data := url.Values{}
	data.Set("username", "wlw_aoxiangjun")
	data.Set("password", "WLW@aoxiangjun123")
	data.Set("local", "zh-CN")
	// fmt.Println("data:", data)
	// fmt.Println("data.Encode:", data.Encode())
	// fmt.Println("string.NewReader:", strings.NewReader(data.Encode()))

	req, _ := http.NewRequest(
		http.MethodPost,
		baseURL+"/passport/login.do",
		strings.NewReader(data.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Referer", baseURL+"/goLogin.do") // 设置请求来源，模拟从goLogin跳转过来

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
