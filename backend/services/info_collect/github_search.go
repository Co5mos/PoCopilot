package info_collect

import (
	"PoCopilot/backend/pkg/github"
	"PoCopilot/backend/services"
	"net/url"
	"sort"
	"time"

	github2 "github.com/google/go-github/v61/github"
)

type githubSearchResult struct {
	FullName string `json:"full_name"`
	HtmlURL  string `json:"html_url"`
	PushedAt string `json:"pushed_at"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
}

func GithubSearchPoc(githubToken, keyword string) services.Msg {

	results := make([]githubSearchResult, 0)

	githubConfig := &github.Config{
		GithubToken: githubToken,
	}
	gs, err := github.NewGithubService(githubConfig)
	if err != nil {
		return services.Msg{
			Code: 500,
			Msg:  err.Error(),
		}
	}

	// 查询代码
	searchResults, err := gs.GetSearchCode(keyword)
	if err != nil {
		return services.Msg{
			Code: 500,
			Msg:  err.Error(),
		}
	}
	if len(searchResults) > 0 {
		for _, searchResult := range searchResults {
			for _, code := range searchResult.CodeResults {

				var pushedAt string
				path := *code.Path
				commits, err := gs.GetFileCommits(*code.Repository.Owner.Login, *code.Repository.Name, path)

				// 检查 commits 列表是否为空
				if err != nil || len(commits) == 0 {
					pushedAt = ""
				} else {
					// 获取最后一个 commit
					commit := commits[len(commits)-1]
					pushedAt = timeToString(commit.Commit.Committer.Date)
				}

				result := githubSearchResult{
					FullName: *code.Repository.FullName,
					HtmlURL:  *code.HTMLURL,
					PushedAt: pushedAt,
					FileName: *code.Name,
				}
				results = append(results, result)
			}
		}
	}

	// 查询仓库
	repos, err := gs.GetRepos(keyword)
	if err != nil {
		return services.Msg{
			Code: 500,
			Msg:  err.Error(),
		}
	}
	if *repos.Total > 0 {
		for _, repo := range repos.Repositories {
			result := githubSearchResult{
				FullName: *repo.FullName,
				HtmlURL:  *repo.HTMLURL,
				PushedAt: timeToString(repo.PushedAt),
				FileName: "",
			}
			results = append(results, result)
		}
	}

	// 对 results 按 PushedAt 时间倒序排序
	sort.Slice(results, func(i, j int) bool {
		ti, _ := time.Parse("2006-01-02 15:04:05", results[i].PushedAt)
		tj, _ := time.Parse("2006-01-02 15:04:05", results[j].PushedAt)
		return ti.After(tj)
	})

	return services.Msg{
		Code: 200,
		Msg:  results,
	}
}

func GithubSearchCode(githubToken, htmlURL string) services.Msg {
	githubConfig := &github.Config{
		GithubToken: githubToken,
	}
	gs, err := github.NewGithubService(githubConfig)
	if err != nil {
		return services.Msg{
			Code: 500,
			Msg:  err.Error(),
		}
	}

	// 将 htmlURL 转换为 owner, repo, path
	decodedURL, err := url.QueryUnescape(htmlURL)
	if err != nil {
		return services.Msg{
			Code: 500,
			Msg:  err.Error(),
		}
	}
	owner, repo, path := gs.ParseHtmlURL(decodedURL)
	file, err := gs.GetFileCode(owner, repo, path)
	if err != nil {
		return services.Msg{
			Code: 500,
			Msg:  err.Error(),
		}
	}

	if file != nil {
		return services.Msg{
			Code: 200,
			Msg:  *file.Content,
		}
	} else {
		return services.Msg{
			Code: 500,
			Msg:  "获取文件代码失败",
		}
	}
}

func timeToString(t *github2.Timestamp) string {
	if t != nil {
		// 将时间转换为东八区
		loc, _ := time.LoadLocation("Asia/Shanghai")
		eastEightZoneTime := t.In(loc)
		// 返回东八区的时间字符串表示，例如："2006-01-02 15:04:05"
		return eastEightZoneTime.Format("2006-01-02 15:04:05")
	}
	return ""
}
