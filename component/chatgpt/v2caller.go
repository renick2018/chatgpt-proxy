package chatgpt

import (
	"bytes"
	"chatgpt-proxy/lib/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func callServer(url string, data map[string]interface{}) (map[string]interface{}, error) {
	bs, _ := json.Marshal(data) // POST 请求的数据

	logger.Info(fmt.Sprintf("%s request: %s", url, string(bs)))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	if err != nil {
		logger.Warning(fmt.Sprintf("Error creating HTTP request: %+v", err))
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json") // 设置请求头

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req) // 发送请求
	if err != nil {
		logger.Warning(fmt.Sprintf("Error sending HTTP request: %+v", err))
		return nil, err
	}

	defer resp.Body.Close()

	logger.Info(fmt.Sprintf("HTTP Response Status: %+v", resp.Status))

	// 读取响应体
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	logger.Info(fmt.Sprintf("HTTP Response Body:: %+v", buf.String()))

	var rsp = make(map[string]interface{})
	err = json.Unmarshal(buf.Bytes(), &rsp)
	return rsp, err
}
