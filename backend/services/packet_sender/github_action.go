package packet_sender

import (
	"PoCopilot/backend/common"
	"PoCopilot/backend/pkg/github"
	"PoCopilot/backend/services"
	"time"

	"github.com/sirupsen/logrus"
)

func SendGithubAction(owner, repoName, token, rawData string, targetList []string) services.Msg {
	gas, err := github.NewGithubActionSender(owner, repoName, token)
	if err != nil {
		return services.Msg{
			Code: 500,
			Msg:  err.Error(),
		}
	}

	err = gas.Send(rawData, targetList)
	if err != nil {
		return services.Msg{
			Code: 500,
			Msg:  err.Error(),
		}
	} else {
		return services.Msg{
			Code: 200,
			Msg:  "发送成功",
		}
	}
}

func GetGithubActionLog(owner, repoName, token string, targetNum int, ws *common.WebsocketService) services.Msg {

	githubConfig := &github.Config{
		Owner:       owner,
		RepoName:    repoName,
		GithubToken: token,
	}

	gas, err := github.NewGithubService(githubConfig)
	if err != nil {
		return services.Msg{
			Code: 500,
			Msg:  err.Error(),
		}
	}

	gas.TargetNum = targetNum
	time.Sleep(2 * time.Second)

	jobs, err := gas.GetJobs(owner, repoName)
	if err != nil {
		logrus.Errorf("Error fetching jobs: %s", err)
		return services.Msg{
			Code: 500,
			Msg:  "获取 Job 失败",
		}
	}
	logrus.Infof("lasted job: %s, %s", jobs[0].GetName(), jobs[0].GetStartedAt())

	// TODO 解决死循环问题
	go func() {
		end := gas.EndProgress
		num := 0.0
		completedStepNum := 0

		for {
			time.Sleep(1 * time.Second)
			logrus.Info("Running...")

			completedStepNum, err = gas.GetCompletedStepNum(owner, repoName, jobs[0])
			if err != nil {
				logrus.Errorf("Error fetching completed step num: %s", err)
				continue
			}

			if completedStepNum < gas.TargetNum {
				continue
			}

			acLog, err := gas.GetJobLog(owner, repoName, jobs[0])
			if err != nil {
				logrus.Errorf("Error fetching job log: %s", err)
				continue
			}

			prettyAcLog := common.PrettyLog(*acLog)

			// Check if WebSocket is still active before sending
			if ws != nil && prettyAcLog != nil {
				logrus.Info("Send to client...")
				ws.BroadcastToClients(*prettyAcLog)
			}

			select {
			case <-end:
				return
			default:
			}

			num = float64(completedStepNum) / float64(gas.TargetNum)
			logrus.Infof("completedStepNum: %d, TargetNum: %d, Percent: %f", completedStepNum, gas.TargetNum, num)

			if completedStepNum == gas.TargetNum {
				break
			}
		}

		// TODO make sure this resets when we hide etc...
		logrus.Info("Github Action Completed")
		gas.ActionsStopProgress()
	}()

	return services.Msg{
		Code: 200,
		Msg:  "success",
	}
}
