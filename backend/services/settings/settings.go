package settings

import (
	"PoCopilot/backend/common"

	"github.com/sirupsen/logrus"
)

func ReadConfig(filename string) *common.Config {
	config, err := common.ReadConfig(filename)
	if err != nil {
		logrus.Errorf("读取配置文件失败: %s", err)
	}
	return config
}

func WriteConfig(filename string, config *common.Config) (*common.Config, string) {
	err := common.WriteConfig(filename, config)
	if err != nil {
		errMsg := "写入配置文件失败: " + err.Error()
		logrus.Error(errMsg)
		return nil, errMsg
	}

	config, _, _ = common.InitConfig()
	return config, "写入配置文件成功"
}
