package github

import (
	"PoCopilot/backend/common"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

type ActionSender struct {
	githubService *Service
}

func NewGithubActionSender(owner, repoName, token string) (*ActionSender, error) {

	gas := &ActionSender{}
	githubConfig := &Config{
		Owner:       owner,
		RepoName:    repoName,
		GithubToken: token,
	}
	githubService, err := NewGithubService(githubConfig)
	if err != nil {
		return nil, err
	}

	gas.githubService = githubService
	return gas, nil
}

// Send 发送 github action 操作
func (a *ActionSender) Send(rawData string, targetList []string) error {

	// 创建 action 文件
	ac, err := a.CreateActionFile(&targetList, strings.NewReader(rawData))
	if err != nil {
		logrus.Errorf("Create Action File Error: %s", err)
		return err
	}
	logrus.Info("Create Action File Success")

	// 创建 workflows 文件
	err = a.githubService.CreateFile(a.githubService.config.Owner, a.githubService.config.RepoName, []byte(*ac))
	if err != nil {
		logrus.Errorf("Create Workflows File Error: %s", err)
		return err
	}
	logrus.Info("Create Workflows File Success")

	return nil
}

/*
CreateActionFile
生成 action 文件
*/
func (a *ActionSender) CreateActionFile(targets *[]string, rawHttp io.Reader) (*string, error) {

	// 生成6位随机字符串
	randStr := common.GenerateRandomString(6)
	logrus.Infof("Create Action File KeyWord: %s", randStr)
	actionContent := fmt.Sprintf(`name: autoAction
on: [push]
jobs:
  Test-GitHub-Actions-%s:
    runs-on: ubuntu-latest
    steps:
`, randStr)

	command, u, err := common.HttpRaw2Curl((*targets)[0], rawHttp)
	if err != nil {
		return nil, err
	}

	for _, host := range *targets {
		// TODO 判断 host 正确性
		targetURL := host + u.RequestURI()
		addCommand := append(command[:len(command)-1], fmt.Sprintf("$'%s'", targetURL))
		actionContent += "      - name: curl\n"
		actionContent += "        continue-on-error: true\n"
		actionContent += "        run: |\n"
		actionContent += fmt.Sprintf("          %s\n", "echo -e '\\n-------------------------------------------------\\n'")
		actionContent += fmt.Sprintf("          echo -e '\\n%s\\n'\n", targetURL)
		actionContent += fmt.Sprintf("          %s\n", strings.Join(addCommand, " "))
		actionContent += fmt.Sprintf("          %s\n", "echo -e '\\n-------------------------------------------------\\n'")
	}

	return &actionContent, nil
}
