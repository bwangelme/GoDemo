package oauth

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	url = `https://{授权页链接}?response_type=code&client_id={应用client_id}&redirect_uri={client_id对应的回调地址}&state={自定义参数}`
	url = `https://fuwu.pinduoduo.com/service-market/auth`
)

// Response struct 用于解析 OAuth token 响应
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// 获取 access_token 的函数
func getAccessToken(clientID, clientSecret, tokenURL string) (string, error) {
	// 构造 POST 请求体
	formData := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
	}

	// 发送 POST 请求
	resp, err := http.PostForm(tokenURL, formData)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// 检查是否请求成功
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status: %s, body: %s", resp.Status, string(body))
	}

	// 解析 JSON 响应
	var tokenResp TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// 返回 access_token
	return tokenResp.AccessToken, nil
}

func main() {
	// OAuth 2.0 的 token endpoint
	tokenURL := "https://your-oauth-server.com/oauth/token" // 替换成实际的 OAuth 服务器的 token URL
	clientID := "your_client_id"                            // 替换成实际的 client_id
	clientSecret := "your_client_secret"                    // 替换成实际的 client_secret

	// 获取 access_token
	accessToken, err := getAccessToken(clientID, clientSecret, tokenURL)
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
	}

	// 输出获取到的 access_token
	fmt.Printf("Access Token: %s\n", accessToken)
}
