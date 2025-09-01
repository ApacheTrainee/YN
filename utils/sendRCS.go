package utils

import (
	"YN/config"
	"YN/log"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// RCS请求参数定义
type fieldFunctionSend struct {
	Data struct {
		FunctionName string `json:"function_name"`
		Value        bool   `json:"value"`
	} `json:"data"`
}

// RCS接口响应参数定义
type fieldFunctionResponse struct {
	ErrInfo struct {
		ErrMsg  string `json:"err_msg"`
		ErrCode int    `json:"err_code"`
	} `json:"err_info"`
}

// 调用RCS接口更新Field-function值
func SendRCS(functionName string, value bool) error {
	request := fieldFunctionSend{}
	request.Data.FunctionName = functionName
	request.Data.Value = value

	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}
	log.Logger.Infof("send to RCS: request body is %v", string(jsonData))

	// 创建HTTP请求
	url := fmt.Sprintf("http://%s:%s/rbrainrobot/set_field_function", config.Config.RcsIP, config.Config.RcsPort)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 解析响应
	var response fieldFunctionResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	log.Logger.Infof("RCS response body: %+v", response)

	// 检查错误码
	if response.ErrInfo.ErrCode != 0 {
		return fmt.Errorf("RCS error: %v (code: %v)", response.ErrInfo.ErrMsg, response.ErrInfo.ErrCode)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("RCS returned status: %d", resp.StatusCode)
	}

	return nil
}
