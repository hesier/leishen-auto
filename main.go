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

	// 首先获取用户状态
	fmt.Println("🔍正在检查用户状态...")
	userInfo, err := client.GetUserInfo(cfg.AccountToken, cfg.Lang)
	if err != nil {
		errorMsg := fmt.Sprintf("获取用户状态失败: %v", err)
		fmt.Printf("❌%s\n", errorMsg)
		// 发送错误通知
		if notifyErr := barkNotifier.SendErrorNotification(errorMsg); notifyErr != nil {
			fmt.Printf("❌发送通知失败: %v\n", notifyErr)
		}
		os.Exit(1)
	}

	if userInfo.Code != 0 {
		errorMsg := fmt.Sprintf("获取用户状态失败: %d - %s", userInfo.Code, userInfo.Msg)
		fmt.Printf("❌%s\n", errorMsg)
		// 发送错误通知
		if notifyErr := barkNotifier.SendErrorNotification(errorMsg); notifyErr != nil {
			fmt.Printf("❌发送通知失败: %v\n", notifyErr)
		}
		os.Exit(1)
	}

	// 检查暂停状态
	if userInfo.Data.PauseStatusID == 1 {
		// 已暂停状态
		msg := fmt.Sprintf("用户已经在暂停状态，暂停状态：%s，最后暂停时间：%s",
			userInfo.Data.PauseStatus, userInfo.Data.LastPauseTime)
		fmt.Printf("👌%s\n", msg)
		// 发送状态通知
		if notifyErr := barkNotifier.SendNotification("雷神加速器状态", msg); notifyErr != nil {
			fmt.Printf("❌发送通知失败: %v\n", notifyErr)
		}
		fmt.Println("⌛️结束运行")
		return
	}

	fmt.Printf("✅当前状态：%s，可以进行暂停操作\n", userInfo.Data.PauseStatus)

	// 未暂停状态，执行暂停操作
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
