package main

import (
	"fmt"
	"os"
	"time"

	"leishen-auto/api"
	"leishen-auto/config"
)

func main() {
	fmt.Println("⌛️开始运行")

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("❌错误: %v\n", err)
		os.Exit(1)
	}

	// 创建Bark通知器
	barkNotifier := api.NewBarkNotifier(cfg.BarkToken)

	client := api.NewClient()

	resp, err := client.Pause(cfg.AccountToken, cfg.Lang)
	if err != nil {
		errorMsg := fmt.Sprintf("暂停失败: %v", err)
		fmt.Printf("❌%s\n", errorMsg)
		// 发送错误通知
		if notifyErr := barkNotifier.SendErrorNotification(errorMsg); notifyErr != nil {
			fmt.Printf("❌发送通知失败: %v\n", notifyErr)
		}
		os.Exit(1)
	}

	if resp.Code != 0 {
		if resp.Code == 400803 { // 400803 - 账号已经停止加速，请不要重复操作
			msg := fmt.Sprintf("已经暂停: %d - %s", resp.Code, resp.Msg)
			fmt.Printf("👌%s\n", msg)
			// 发送重复暂停通知
			if notifyErr := barkNotifier.SendNotification("雷神加速器状态", msg); notifyErr != nil {
				fmt.Printf("❌发送通知失败: %v\n", notifyErr)
			}
			fmt.Println("⌛️结束运行")
			return
		}
		errorMsg := fmt.Sprintf("暂停失败: %d - %s", resp.Code, resp.Msg)
		fmt.Printf("❌%s\n", errorMsg)
		// 发送错误通知
		if notifyErr := barkNotifier.SendErrorNotification(errorMsg); notifyErr != nil {
			fmt.Printf("❌发送通知失败: %v\n", notifyErr)
		}
		os.Exit(1)
	}

	// 发送成功通知
	successMsg := fmt.Sprintf("%s - %s", time.Now().Format("2006-01-02 15:04:05"), resp.Msg)
	if notifyErr := barkNotifier.SendSuccessNotification(successMsg); notifyErr != nil {
		fmt.Printf("❌发送通知失败: %v\n", notifyErr)
	}

	fmt.Printf("%d:%s\n", resp.Code, resp.Msg)
	fmt.Println("✔️暂停成功")
	fmt.Println("⌛️结束运行")
}
