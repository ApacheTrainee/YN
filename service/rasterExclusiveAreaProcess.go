package service

import (
	"YN/global"
	"YN/log"
	"YN/utils"
)

func RasterExclusiveAreaProcess() {
	go func() {
		for {
			select {
			case boolInfo := <-global.RasterExclusiveAreaChan1:
				if boolInfo != global.RasterExclusiveArea1 {
					// 开始进入独占区触发
					if boolInfo == true {
						if err := utils.SendRCS("RasterExclusiveArea1", true); err != nil {
							log.WebLogger.Errorf("SendRCS RasterExclusiveArea1 err: %v", err)
						}
					}

					// 离开完独占区后出发
					if boolInfo == false {
						log.WebLogger.Infof("已经离开独占区1")
					}

					global.RasterExclusiveArea1 = boolInfo
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case boolInfo := <-global.RasterExclusiveAreaChan2:
				if boolInfo != global.RasterExclusiveArea2 {
					// 开始进入独占区触发
					if boolInfo == true {
						if err := utils.SendRCS("RasterExclusiveArea2", true); err != nil {
							log.WebLogger.Errorf("SendRCS RasterExclusiveArea2 err: %v", err)
						}
					}

					// 离开完独占区后出发
					if boolInfo == false {
						log.WebLogger.Infof("已经离开独占区2")
					}

					global.RasterExclusiveArea2 = boolInfo
				}
			}
		}
	}()
}
