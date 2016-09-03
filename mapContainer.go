package handleFeedback

import (
	"fmt"
	"os"
	"strings"
)

type acMap map[int]phoneStatusMap
type phoneStatusMap map[string][]string
type mapContainer struct {
	c map[int]acMap
}

// NewmapContainer 使用map作为存储用户活动返回的手机号码状态的容器
func NewmapContainer() Container {
	return &mapContainer{
		c: make(map[int]acMap),
	}
}

//Push 将数据储存到容器中
func (mc *mapContainer) Push(arg interface{}) {
	switch message := arg.(type) {
	case YiMeiMessage:
		mc.pushYiMeiXMLData(message)
		return
	default:
		return
	}
}

func (mc *mapContainer) pushYiMeiXMLData(message YiMeiMessage) {
	activityID := message.Seqid >> 16
	userID := message.Seqid - (activityID << 16)
	// fmt.Printf("userID:%d, activityID:%d\n", userID, activityID)
	if _, ok := mc.c[userID]; !ok {
		mc.c[userID] = make(acMap)
		// mc.c[userID][activityID] = make(phoneStatusMap)
	}
	if _, ok := mc.c[userID][activityID]; !ok {
		mc.c[userID][activityID] = make(phoneStatusMap)
	}

	if _, ok := mc.c[userID][activityID][message.State]; ok {
		mc.c[userID][activityID][message.State] = append(mc.c[userID][activityID][message.State], message.PhoneNum)
	} else {
		mc.c[userID][activityID][message.State] = make([]string, 1)
		mc.c[userID][activityID][message.State][0] = message.PhoneNum
	}
}

// Print 输出存储的信息到终端上
func (mc *mapContainer) Print() {
	// fmt.Printf("map: %v\n", mc.c)
	for userID, am := range mc.c {
		fmt.Printf("userID:%d \n", userID)
		for activityID, sm := range am {
			fmt.Printf("activityID:%d	", activityID)
			for status, phones := range sm {
				fmt.Printf("status:%s, phones:%d	", status, len(phones))
			}
		}
		fmt.Printf("------\n")
	}
}

// EchoErrorPhones 将状态不是DELIVRD的号码全部写入特定的文件中
func (mc *mapContainer) EchoErrorPhones() error {
	for userID, am := range mc.c {
		for activityID, sm := range am {
			for status, phones := range sm {
				if !strings.EqualFold(status, "DELIVRD") {
					path := fmt.Sprintf("%d", userID)
					fileName := fmt.Sprintf("%d/%d.txt", userID, activityID)
					if err := os.MkdirAll(path, 0777); err != nil {
						return err
					}
					file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
					if err != nil {
						return err
					}
					p := strings.Join(phones, "\r\n")
					file.Write([]byte(p))
					file.Close()
				}
			}
		}
	}
	return nil
}

func (mc *mapContainer) EchoStatisticData() error {
	for userID, am := range mc.c {
		for activityID, sm := range am {
			errNum := 0
			for status, phones := range sm {
				if !strings.EqualFold(status, "DELIVRD") {
					errNum = errNum + len(phones)
				}
			}
			if errNum == 0 {
				continue
			}
			fileName := fmt.Sprintf("statistic.txt")
			file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			defer file.Close()
			if err != nil {
				return err
			}
			content := fmt.Sprintf("userID: %d , activityID: %d , errNum: %d , \r\n", userID, activityID, errNum)
			file.Write([]byte(content))
		}
	}
	return nil
}
