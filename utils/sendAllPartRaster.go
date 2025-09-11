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
type reqAllPartRasterDTO struct {
	Position  string `json:"position"`
	TaskState int    `json:"task_state"`
}

// RCS接口响应参数定义
type responseAllPartRasterDTO struct {
	ReturnMessage string `json:"Return_message"`
}

// 调用RCS接口更新Field-function值
func SendAllPartRaster(position string, taskState int) error {
	request := reqAllPartRasterDTO{
		Position:  position,
		TaskState: taskState,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}
	log.Logger.Infof("send to all-part: request body is %v", string(jsonData))

	// 创建HTTP请求
	url := fmt.Sprintf("http://%s:%s/api/SmartWarehous/Update_task_state_Grating", config.Config.AllPartIP, config.Config.AllPartPort)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 解析响应
	var response responseAllPartRasterDTO
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	log.Logger.Infof("all-part response body: %+v", response)

	// 检查错误码
	if response.ReturnMessage != "OK" {
		return fmt.Errorf("all-part error: %v", response.ReturnMessage)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("all-part returned status: %d", resp.StatusCode)
	}

	return nil
}
