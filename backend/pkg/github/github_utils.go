package github

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v61/github"
	"github.com/sirupsen/logrus"
)

/**
github 信息获取
*/

/*
GetGithubRepoInfo
查询仓库信息
*/
func (g *Service) GetGithubRepoInfo(owner string, repoName string) *github.Repository {

	rs, resp, err := g.client.Repositories.Get(*g.ctx, owner, repoName)
	logrus.Infof("Github repo by owner resp status code: %d", resp.StatusCode)

	if err != nil {
		logrus.Error(err.Error())
	}

	logrus.Infof("Get repo from github... %s/%s", owner, repoName)
	return rs
}

/*
GetRepoInfoByID
通过仓库 id 获取仓库信息
*/
func (g *Service) GetRepoInfoByID(repoID int64) (bool, *github.Repository) {
	repo, resp, err := g.client.Repositories.GetByID(*g.ctx, repoID)
	logrus.Infof("Github repo by id resp status code...%d", resp.StatusCode)

	if err != nil {
		logrus.Error(err.Error())
		return false, nil
	} else {
		return true, repo
	}
}

/*
GetLastCommitDatetime
获取最新一次master分支的commit时间
*/
func (g *Service) GetLastCommitDatetime(owner, repo string) time.Time {
	branch, resp, err := g.client.Repositories.GetBranch(*g.ctx, owner, repo, "master", 1)
	if err != nil {
		logrus.Errorf("Get branch failed...%s/%s", owner, repo)
	}

	logrus.Infof("Get branch success...%d", resp.StatusCode)
	datetime := branch.GetCommit().Commit.Author.GetDate()
	newTime := datetime.Add(1 * time.Second)

	return newTime
}

/*
GetGithubRepoPushedData
获取 GitHub repo 最新的 push 数据
*/
func (g *Service) GetGithubRepoPushedData(owner, repo, pushedAt string) []*string {

	// 获取最新一次push的commits
	since, err := time.Parse("2006-01-02 15:04:05", pushedAt)
	if err != nil {
		logrus.Error(err.Error())
	}
	commits, resp, err := g.client.Repositories.ListCommits(*g.ctx, owner, repo, &github.CommitsListOptions{Since: since})
	if err != nil {
		logrus.Error(err.Error())
	}
	logrus.Infof("Github commits resp status code...%d", resp.StatusCode)

	var addedFiles []*string
	for _, commit := range commits {
		commitSHA := commit.SHA

		// 获取 commit files
		opt := &github.ListOptions{}
		repoCommit, resp1, err1 := g.client.Repositories.GetCommit(*g.ctx, owner, repo, *commitSHA, opt)
		if err1 != nil {
			logrus.Error(err1.Error())
		}
		logrus.Infof("Github commit resp status code...%d", resp1.StatusCode)
		repoCommitFiles := repoCommit.Files

		for _, f := range repoCommitFiles {
			// 添加文件
			if *f.Status == "added" {
				logrus.Infof("New pushed file...%s", *f.Filename)
				addedFiles = append(addedFiles, f.Filename)
			}
		}
	}

	return addedFiles
}

/*
GetRepos
检索仓库
*/
func (g *Service) GetRepos(keyword string) (*github.RepositoriesSearchResult, error) {
	result, _, err := g.client.Search.Repositories(*g.ctx, keyword, &github.SearchOptions{Sort: "updated"})
	if err != nil {
		return nil, err
	}
	return result, nil
}

/*
GetSearchCode
搜索代码
*/
func (g *Service) searchCodeByExtension(keyword, extension string) (*github.CodeSearchResult, error) {
	query := fmt.Sprintf("%s in:file extension:%s", keyword, extension)
	result, _, err := g.client.Search.Code(*g.ctx, query, &github.SearchOptions{Sort: "indexed"})
	if err != nil {
		return nil, err
	}

	if *result.Total > 0 {
		return result, nil
	}

	return nil, nil
}

func (g *Service) GetSearchCode(keyword string) ([]*github.CodeSearchResult, error) {
	result := make([]*github.CodeSearchResult, 0)

	extensions := []string{"yaml", "py"}

	for _, extension := range extensions {
		csr, err := g.searchCodeByExtension(keyword, extension)
		if err != nil {
			return nil, err
		}

		if csr != nil {
			result = append(result, csr)
		}
	}

	return result, nil
}

/*
GetFileCode
获取文件代码
*/
func (g *Service) GetFileCode(owner, repo, path string) (*github.RepositoryContent, error) {
	file, _, _, err := g.client.Repositories.GetContents(*g.ctx, owner, repo, path, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, err
	}
	return file, nil
}

/*
GetFileCommits
获取文件 commits
*/
func (g *Service) GetFileCommits(owner, repo, path string) ([]*github.RepositoryCommit, error) {
	commits, _, err := g.client.Repositories.ListCommits(*g.ctx, owner, repo, &github.CommitsListOptions{Path: path})
	if err != nil {
		return nil, err
	}
	return commits, nil
}

func (g *Service) ParseHtmlURL(htmlURL string) (string, string, string) {
	// 移除URL前缀 "https://github.com/"
	trimURL := strings.TrimPrefix(htmlURL, "https://github.com/")

	// 使用 "/" 分割字符串
	parts := strings.Split(trimURL, "/")

	// 检查分割后的部分是否满足基本的长度要求
	if len(parts) < 5 {
		return "", "", "" // 返回空字符串，如果不符合预期格式
	}

	owner := parts[0] // 第一部分是 owner
	repo := parts[1]  // 第二部分是 repo

	// 由于路径可能包含多个部分，我们需要将剩余的部分再组合起来
	path := strings.Join(parts[4:], "/")

	return owner, repo, path
}
