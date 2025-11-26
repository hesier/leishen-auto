package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PauseRequest 暂停请求结构体
type PauseRequest struct {
	AccountToken string `json:"account_token"`
	Lang         string `json:"lang"`
}

// PauseResponse 暂停响应结构体
type PauseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// UserInfoResponse 用户信息响应结构体
type UserInfoResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		UserPauseTime        int    `json:"user_pause_time"`
		Nickname             string `json:"nickname"`
		Mobile               string `json:"mobile"`
		PauseStatus          string `json:"pause_status"`
		PauseStatusID        int    `json:"pause_status_id"`
		LastPauseTime        string `json:"last_pause_time"`
		VipLevel             string `json:"vip_level"`
		PackageTitle         string `json:"package_title"`
		StopedRemaining      string `json:"stoped_remaining"`
		ExpiryTime           string `json:"expiry_time"`
		IsPayUser            string `json:"is_pay_user"`
		MobilePauseStatus    int    `json:"mobile_pause_status"`
		MobileExpiryTime     string `json:"mobile_expiry_time"`
		MobileExpiryTimeSamp int    `json:"mobile_expiry_time_samp"`
		NowDate              string `json:"now_date"`
		NowTimeSamp          int    `json:"now_time_samp"`
		UserEarnMinutes      string `json:"user_earn_minutes"`
	} `json:"data"`
}

// HTTPClient HTTP客户端接口
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client 雷神加速器客户端
type Client struct {
	BaseURL    string
	HTTPClient HTTPClient
	Timeout    time.Duration
}

// NewClient 创建新的雷神客户端
func NewClient() *Client {
	return &Client{
		BaseURL: "https://webapi.leigod.com",
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		Timeout: 5 * time.Second,
	}
}

// Pause 暂停加速器
func (c *Client) Pause(accountToken, lang string) (*PauseResponse, error) {
	pauseReq := PauseRequest{
		AccountToken: accountToken,
		Lang:         lang,
	}

	jsonData, err := json.Marshal(pauseReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	url := c.BaseURL + "/api/user/pause"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var pauseResp PauseResponse
	if err := json.Unmarshal(body, &pauseResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &pauseResp, nil
}

// GetUserInfo 获取用户信息
func (c *Client) GetUserInfo(accountToken, lang string) (*UserInfoResponse, error) {
	userInfoReq := PauseRequest{
		AccountToken: accountToken,
		Lang:         lang,
	}

	jsonData, err := json.Marshal(userInfoReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	url := c.BaseURL + "/api/user/info"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var userInfoResp UserInfoResponse
	if err := json.Unmarshal(body, &userInfoResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &userInfoResp, nil
}
