package api

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// BarkNotifier Bark通知器
type BarkNotifier struct {
	Token string
	Client *http.Client
}

// NewBarkNotifier 创建新的Bark通知器
func NewBarkNotifier(token string) *BarkNotifier {
	return &BarkNotifier{
		Token: token,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendNotification 发送通知
func (b *BarkNotifier) SendNotification(title, body string) error {
	if b.Token == "" {
		return nil // 如果没有配置token，直接返回，不发送通知
	}

	barkURL := fmt.Sprintf("https://api.day.app/%s/%s/%s", b.Token, url.QueryEscape(title), url.QueryEscape(body))

	resp, err := b.Client.Get(barkURL)
	if err != nil {
		return fmt.Errorf("发送Bark通知失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bark通知发送失败，状态码: %d", resp.StatusCode)
	}

	return nil
}

// SendSuccessNotification 发送成功通知
func (b *BarkNotifier) SendSuccessNotification(msg string) error {
	return b.SendNotification("雷神加速器暂停成功", msg)
}

// SendErrorNotification 发送错误通知
func (b *BarkNotifier) SendErrorNotification(errMsg string) error {
	return b.SendNotification("雷神加速器暂停失败", errMsg)
}