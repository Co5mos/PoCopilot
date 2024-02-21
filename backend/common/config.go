package common

const (
	App        = "PoCopilot"   // app name
	AppLogFile = "app.log"     // app log 路径
	ConfigFile = "config.yaml" // config.yaml 路径
	OutputDir  = "output/"     // nuclei output 路径
	TmpDir     = "tmp/"        // tmp 路径
)

/*
Config
配置信息
*/
type Config struct {
	GithubToken string `yaml:"GithubToken"` // github token
	Owner       string `yaml:"Owner"`       // github 仓库拥有者
	RepoName    string `yaml:"RepoName"`    // github 仓库名称
}
