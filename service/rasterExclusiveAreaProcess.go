package service

import (
	"YN/global"
	"YN/log"
	"YN/utils"
)

func RasterExclusiveAreaProcess() {
	go func() {
		for boolInfo := range global.RasterExclusiveAreaChan1 { // 遍历完也会阻塞等待。除非发起方关闭管道
			if boolInfo == global.RasterExclusiveArea1 {
				continue
			}

			// 开始进入独占区触发
			if boolInfo == true {
				// 先发请求给all part关闭光栅
				if err := utils.SendAllPartRaster("5M01", 1); err != nil {
					log.WebLogger.Errorf("SendAllPartRaster RasterExclusiveArea1 err: %v", err)
				}

				// 后发请求给RCS
				if err := utils.SendRCS("5M01", true); err != nil {
					log.WebLogger.Errorf("SendRCS RasterExclusiveArea1 err: %v", err)
				}
			} else { // 离开完独占区后触发
				// 发请求给all part开启光栅
				if err := utils.SendAllPartRaster("5M01", 2); err != nil {
					log.WebLogger.Errorf("SendAllPartRaster RasterExclusiveArea1 err: %v", err)
				}

				log.WebLogger.Infof("AGV已经离开独占区1: 5M01")
			}

			global.RasterExclusiveArea1 = boolInfo
		}
	}()

	go func() {
		for boolInfo := range global.RasterExclusiveAreaChan2 { // 遍历完也会阻塞等待。除非发起方关闭管道
			if boolInfo == global.RasterExclusiveArea2 {
				continue
			}

			// 开始进入独占区触发
			if boolInfo == true {
				// 先发请求给all part关闭光栅
				if err := utils.SendAllPartRaster("5M02", 1); err != nil {
					log.WebLogger.Errorf("SendAllPartRaster RasterExclusiveArea2 err: %v", err)
				}

				if err := utils.SendRCS("5M02", true); err != nil {
					log.WebLogger.Errorf("SendRCS RasterExclusiveArea2 err: %v", err)
				}
			} else { // 离开完独占区后触发
				// 先发请求给all part关闭光栅
				if err := utils.SendAllPartRaster("5M02", 2); err != nil {
					log.WebLogger.Errorf("SendAllPartRaster RasterExclusiveArea2 err: %v", err)
				}

				log.WebLogger.Infof("AGV已经离开独占区2: 5M02")
			}

			global.RasterExclusiveArea2 = boolInfo
		}
	}()
}
