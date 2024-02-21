package common

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type FilePath struct {
	// 初始化文件夹路径
	InitFile string
	// 配置文件路径
	ConfigFile string
	// 日志文件路径
	LogFile string
	// 输出文件夹路径
	OutputDir string
	// templates 文件夹路径
	TemplatesDir string
}

func InitConfig() (*Config, *FilePath, error) {
	// 初始化文件路径
	filePath := InitPathAndLogger()

	//NewLogger(logFile)
	logrus.Info("App Init")
	// 配置文件路径
	logrus.Infof("Config File Path: %s", filePath.ConfigFile)

	// 配置文件
	if !Exists(filePath.ConfigFile) {
		// 创建配置文件
		_, err := os.Create(filePath.ConfigFile)
		if err != nil {
			return nil, nil, err
		}

		// 写入默认配置
		DefaultConfig := &Config{}
		err = WriteConfig(filePath.ConfigFile, DefaultConfig)

		// 调转到设置页面
		logrus.Error("请先进行参数配置")
		return nil, nil, err
	}

	// 读取配置文件
	config, err := ReadConfig(filePath.ConfigFile)
	if err != nil {
		return nil, nil, err
	}

	// 判断输出文件夹是否存在
	if !IsDir(filePath.OutputDir) {
		// 创建文件夹
		_ = os.Mkdir(filePath.OutputDir, os.ModePerm)
	}

	if !IsDir(filePath.TemplatesDir) {
		// 创建文件夹
		_ = os.Mkdir(filePath.TemplatesDir, os.ModePerm)
	}

	return config, filePath, nil
}

/*
InitPathAndLogger
初始化文件路径和日志
*/
func InitPathAndLogger() *FilePath {
	// 初始化文件路径
	filePath := &FilePath{}
	initDir := GetConfigDir()
	filePath.InitFile = initDir

	// 日志
	logFile := filepath.Join(initDir, AppLogFile)
	filePath.LogFile = logFile

	// 配置文件
	configFile := filepath.Join(initDir, ConfigFile)
	filePath.ConfigFile = configFile

	return filePath
}
