package main

import (
	"PoCopilot/backend/common"
	"PoCopilot/backend/services"
	"PoCopilot/backend/services/info_collect"
	"PoCopilot/backend/services/packet_sender"
	"PoCopilot/backend/services/settings"
	"context"

	"github.com/sirupsen/logrus"
)

// App struct
type App struct {
	ctx      context.Context
	Config   *common.Config
	FilePath *common.FilePath

	WebsocketService *common.WebsocketService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved,
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 初始化配置
	a.Config, a.FilePath, _ = common.InitConfig()

	// 初始化 websocket
	if a.WebsocketService == nil {
		logrus.Info("init websocket")
		a.WebsocketService = common.NewWebsocketService()
		go a.WebsocketService.Run()
	}
}

/*
配置相关
*/

// ReadConfig 读取配置文件
func (a *App) ReadConfig() *common.Config {
	return settings.ReadConfig(a.FilePath.ConfigFile)
}

// WriteConfig 写入配置文件
func (a *App) WriteConfig(config *common.Config) string {
	config, msg := settings.WriteConfig(a.FilePath.ConfigFile, config)
	a.Config = config
	return msg
}

/*
GitHub Action 相关
*/

// SendGithubAction 发送 GitHub action 操作
func (a *App) SendGithubAction(rawData string, targetList []string) services.Msg {
	return packet_sender.SendGithubAction(a.Config.Owner, a.Config.RepoName, a.Config.GithubToken, rawData, targetList)
}

// GetGithubActionLog 获取 GitHub action 日志
func (a *App) GetGithubActionLog(targetNum int) services.Msg {
	return packet_sender.GetGithubActionLog(a.Config.Owner, a.Config.RepoName, a.Config.GithubToken, targetNum, a.WebsocketService)
}

/*
信息收集
*/

// GithubSearchPoc GithubSearchPoc
func (a *App) GithubSearchPoc(keyword string) services.Msg {
	return info_collect.GithubSearchPoc(a.Config.GithubToken, keyword)
}

// GithubSearchCode GithubSearchCode
func (a *App) GithubSearchCode(htmlURL string) services.Msg {
	return info_collect.GithubSearchCode(a.Config.GithubToken, htmlURL)
}
