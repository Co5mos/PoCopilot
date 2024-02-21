package github

import (
	"PoCopilot/backend/common"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-github/v54/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Service struct {
	client      *github.Client
	ctx         *context.Context
	config      *Config
	EndProgress chan struct{}
	TargetNum   int
}

type Config struct {
	owner       string
	repoName    string
	githubToken string
}

func NewGithubService(owner, repoName, token string) *Service {

	ctx := context.Background()

	s := &Service{
		ctx: &ctx,
		config: &Config{
			owner:       owner,
			repoName:    repoName,
			githubToken: token,
		},
		EndProgress: make(chan struct{}),
		TargetNum:   0,
	}
	s.GetGithubClient()
	return s
}

/*
GetGithubClient
获取 github client
*/
func (g *Service) GetGithubClient() {
	//ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.config.githubToken},
	)
	tc := oauth2.NewClient(*g.ctx, ts)
	client := github.NewClient(tc)

	g.client = client
	//g.ctx = &ctx
}

/*
GetRepo
检索仓库
*/
func (g *Service) GetRepo(owner, repoName string) *github.Repository {
	repo, resp, err := g.client.Repositories.Get(*g.ctx, owner, repoName)
	if err != nil {
		logrus.Error(err.Error())
	}
	logrus.Infof("Github repos resp status code: %d", resp.StatusCode)
	return repo
}

/*
GetAction
获取仓库 action
*/
func (g *Service) GetAction(owner, repoName string) *github.ArtifactList {
	artifacts, resp, err := g.client.Actions.ListArtifacts(*g.ctx, owner, repoName, &github.ListOptions{})
	if err != nil {
		logrus.Error(err.Error())
	}

	logrus.Infof("Github ArtifactList resp status code: %d", resp.StatusCode)

	return artifacts
}

/*
CreateFile
创建 .github/workflows/xxxx.yaml 文件
*/
func (g *Service) CreateFile(owner, repoName string, content []byte) error {

	filepath := ".github/workflows/"

	// 获取文件
	contentGetOpts := &github.RepositoryContentGetOptions{
		Ref: "main",
	}
	_, directoryContent, getResp, err := g.client.Repositories.GetContents(
		*g.ctx, owner, repoName, filepath, contentGetOpts)
	if err != nil {
		return err
	}
	logrus.Infof("Github GetContents resp status code: %d", getResp.StatusCode)

	// 删除旧文件，不然action会全部执行一遍
	for _, file := range directoryContent {

		deleteFileOpts := &github.RepositoryContentFileOptions{
			Message: github.String("Delete test action file"),
			SHA:     file.SHA,
			Branch:  github.String("main"),
		}

		_, deleteResp, err := g.client.Repositories.DeleteFile(
			*g.ctx, owner, repoName, *file.Path, deleteFileOpts)
		if err != nil {
			return err
		}
		logrus.Infof("Github DeleteFile resp status code: %d", deleteResp.StatusCode)
		logrus.Infof("Github DeleteFile: %s", *file.Path)
	}

	createFilePath := filepath + "test.yaml"

	// 创建新文件
	// 生成6位随机字符串
	randStr := common.GenerateRandomString(6)
	logrus.Infof("Workflow Run Key: %s", randStr)
	contentFileOpts := &github.RepositoryContentFileOptions{
		Message: github.String(fmt.Sprintf("Action test message %s", randStr)),
		Content: content,
		Branch:  github.String("main"),
	}
	_, resp, err := g.client.Repositories.CreateFile(*g.ctx, owner, repoName, createFilePath, contentFileOpts)
	if err != nil {
		// logrus.Info(err)
		return err
	}

	logrus.Infof("Github CreateFile resp status code: %d", resp.StatusCode)
	logrus.Infof("Github CreateFile: %s", createFilePath)

	return nil
}

/*
GetJobs
获取 workflow jobs
*/
func (g *Service) GetJobs(owner, repoName string) ([]*github.WorkflowJob, error) {

	// 获取指定 workflow
	workflow, resp, err := g.client.Actions.GetWorkflowByFileName(*g.ctx, owner, repoName, "test.yaml")
	if err != nil {
		// logrus.Info(err)
		return nil, err
	}
	logrus.Infof("Github get specific workflow resp status code: %d", resp.StatusCode)
	if resp.StatusCode == 401 {
		err := errors.New("github get specific workflow resp 401")
		return nil, err
	}

	// 获取最新一条 run 信息
	listOpts := &github.ListOptions{
		Page:    1,
		PerPage: 1,
	}
	opts := &github.ListWorkflowRunsOptions{
		ListOptions: *listOpts,
	}

	// TODO 502 Server Error
	workflowRuns, resp2, err := g.client.Actions.ListWorkflowRunsByID(
		*g.ctx, owner, repoName, workflow.GetID(), opts)
	if err != nil {
		// logrus.Info(err)
		return nil, err
	}
	logrus.Infof("Github get last workflow run resp status code: %d", resp2.StatusCode)
	if resp2.StatusCode == 401 {
		err := errors.New("github get last workflow run resp 401")
		return nil, err
	}

	// 获取 job log
	if *workflowRuns.TotalCount == 0 {
		// logrus.Info("There is no workflow run")
		err := errors.New("there is no workflow run")
		logrus.Error(err.Error())
		return nil, err
	}
	run := workflowRuns.WorkflowRuns[0]

	logrus.Infof("workflow run name: %s", *run.Name)
	jobs, _, err := g.client.Actions.ListWorkflowJobs(*g.ctx, owner, repoName, run.GetID(), &github.ListWorkflowJobsOptions{})
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	if *jobs.TotalCount == 0 {
		// logrus.Info("There is no workflow run")
		err := errors.New("there is no job run")
		logrus.Error(err.Error())
		return nil, err
	}

	// 这种发包方式默认只有一个 job
	return jobs.Jobs, nil
}

/*
GetJobLog
获取 workflow job log
*/
func (g *Service) GetJobLog(owner, repoName string, job *github.WorkflowJob) (*string, error) {
	// TODO 日志量太大，内存占用过大
	client := resty.New()

	jobLogURL, _, err := g.client.Actions.GetWorkflowJobLogs(*g.ctx, owner, repoName, *job.ID, true)
	if err != nil {
		logrus.Errorf("GetWorkflowJobLogs error: %s", err.Error())
		return nil, err
	}

	logResp, err := client.R().Get(jobLogURL.String())
	if err != nil {
		logrus.Errorf("Something wrong: %s", err)
		return nil, err
	}

	logString := string(logResp.Body())

	logrus.Infof("Get latest job status: %s", job.GetStatus())
	if job.GetStatus() == "completed" {
		return &logString, nil
	}
	logrus.Infof("Github wait for job completed: %s", time.Now())

	return &logString, nil
}

/*
GetCompletedStepNum
获取 job step 状态为 completed 的数量
*/
func (g *Service) GetCompletedStepNum(owner, repoName string, job *github.WorkflowJob) (int, error) {

	newJob, jobResp, err := g.client.Actions.GetWorkflowJobByID(*g.ctx, owner, repoName, *job.ID)
	if err != nil {
		// logrus.Info(err)
		return 0, err
	}
	logrus.Infof("Github get last job resp status code: %d", jobResp.StatusCode)

	completedNum := 0

	logrus.Infof("Running job: %s, %s", newJob.GetName(), newJob.StartedAt)
	logrus.Infof("Running job steps len: %d", len(newJob.Steps))
	for _, step := range newJob.Steps {
		logrus.Infof("Step Number: %d, Step Name: %s, Step Status: %s, Step Conclusion: %s",
			step.GetNumber(),
			step.GetName(),
			step.GetStatus(),
			step.GetConclusion(),
		)

		if step.GetName() == "Set up job" || step.GetName() == "Complete job" {
			continue
		}

		if step.GetStatus() == "completed" {
			completedNum += 1
		}
	}

	logrus.Infof("Completed Step Num: %d", completedNum)
	return completedNum, nil
}

// ActionsStopProgress 停止操作
func (g *Service) ActionsStopProgress() {
	logrus.Info("Stop Progress")
	g.EndProgress <- struct{}{}
}
