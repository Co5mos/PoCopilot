package common

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ReadConfig 读取 yaml 文件
func ReadConfig(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}
	return &config, nil
}

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) // os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// GetConfigDir ...
func GetConfigDir() string {
	userPath, err := user.Current()
	if err != nil {
		panic("获取应用配置目录失败: " + err.Error())
	}

	configPath := filepath.Join(userPath.HomeDir, fmt.Sprintf(".%s", strings.ToLower(App)))

	if !Exists(configPath) {
		err = os.Mkdir(configPath, os.ModePerm)
		if err != nil {
			panic("创建应用配置目录失败: " + err.Error())
		}
	}
	return configPath
}

// PrettyLog 切分日志
func PrettyLog(acLog string) *string {

	// .* -------------------------------------------------[\s\S]*? -------------------------------------------------
	// (\d+-\d+-\d+T\d+:\d+:\d+\.\d+Z )(.*)
	var ret []string

	reg1 := regexp.MustCompile(`.* -------------------------------------------------[\s\S]*? -------------------------------------------------`)
	result1 := reg1.FindAllStringSubmatch(acLog, -1)

	for _, v := range result1 {
		reg2 := regexp.MustCompile(`(\d+-\d+-\d+T\d+:\d+:\d+\.\d+Z )(.*)`)
		result2 := reg2.FindAllStringSubmatch(v[0], -1)

		for _, v2 := range result2 {
			ret = append(ret, v2[2])
		}
	}

	prettyAcLog := strings.Join(ret, "\n")
	// 剔除前后的 -------------------------------------------------
	//prettyAcLog = strings.Trim(prettyAcLog, " -------------------------------------------------")
	return &prettyAcLog
}

/*
WriteConfig
写入配置文件
*/
func WriteConfig(filename string, config *Config) error {
	buf, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, buf, 0644)
	if err != nil {
		return err
	}

	return nil
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

/*
HttpRaw2Curl
将 http raw 包 转换成 curl 命令
*/
// TODO 配置 header 可选
func HttpRaw2Curl(host string, rawHttp io.Reader) ([]string, *url.URL, error) {

	// raw to http.request
	buf := bufio.NewReader(rawHttp)

	req, err := http.ReadRequest(buf)
	if err != nil {
		return nil, nil, err
	}

	reqUrl, err := url.Parse(host + req.RequestURI)
	if err != nil {
		return nil, nil, err
	}

	req.RequestURI = ""
	req.URL = reqUrl

	command, err := GetCurlCommand(req)
	if err != nil {
		return nil, nil, err
	}
	return command.Slice, reqUrl, nil
}
