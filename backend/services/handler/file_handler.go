package handler

import (
	"PoCopilot/backend/common"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FileLoader struct {
	http.Handler
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

func (f *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	requestedFilePath := strings.TrimPrefix(req.URL.Path, "/")

	// 读取请求的文件数据
	initDir := common.GetConfigDir()
	filePath := filepath.Join(initDir, common.TmpDir, requestedFilePath)

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(fmt.Sprintf("Could not load file %s", requestedFilePath)))
		return
	}

	// 根据文件扩展名设置响应头的 Content-Type
	if strings.HasSuffix(requestedFilePath, ".png") {
		res.Header().Set("Content-Type", "image/png")
	} else if strings.HasSuffix(requestedFilePath, ".jpg") || strings.HasSuffix(requestedFilePath, ".jpeg") {
		res.Header().Set("Content-Type", "image/jpeg")
	} // ... 其他图像格式

	res.Write(fileData)
}
