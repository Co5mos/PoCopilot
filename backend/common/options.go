package common

import (
	"github.com/projectdiscovery/goflags"
)

type Options struct {
	// CVEs contains a list of CVEs to run nuclei on
	CVEs goflags.StringSlice `json:"cves"`
	// URLs contains a single URLs to run nuclei on
	URLs goflags.StringSlice `json:"urls"`
	// Number of cve to add from nuclei table
	NucleiNumber          int    `json:"nucleiNumber"`
	NucleiTemplateRunPath string `json:"nucleiTemplateRunPath"`

	// Debug is a flag for enabling debugging output
	Debug bool `json:"debug"`
	Proxy bool `json:"proxy"`

	// nuclei
	Nuclei int `json:"nuclei"`

	// xpoc
	Xpoc       bool `json:"xpoc"`       // xpoc
	XpocNumber int  `json:"xpocNumber"` // xpoc 编号

	// AI
	InternalChatGPT bool `json:"internalChatGPT"` // 内部搭建的 ChatGPT
	OpenAIAPI       bool `json:"openAIAPI"`       // 使用 OpenAI 官方 API
	Claude2         bool `json:"claude2"`         // 使用 Claude2

	// CNNVD site
	CnnvdNewSite bool `json:"cnnvdNewSite"`
}
