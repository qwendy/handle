package handleFeedback

import (
	"fmt"

	"gopkg.in/redis.v4"
)

type redisContainer struct {
	client *redis.Client
}

func NewRidesContainer() Container {
	rc := &redisContainer{}
	rc.init()
	return rc
}

func (rc *redisContainer) init() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rc.client = client

}
func (rc *redisContainer) ping() {
	pong, err := rc.client.Ping().Result()
	fmt.Println(pong, err)
}
func (rc *redisContainer) Push(arg interface{}) {
	switch message := arg.(type) {
	case YiMeiMessage:
		rc.pushYiMeiXMLData(message)
		return
	default:
		return
	}

}

func (rc *redisContainer) pushYiMeiXMLData(message YiMeiMessage) {
	activityID := message.Seqid >> 16
	userID := message.Seqid - (activityID << 16)
	// 用户的集合。设置接收到的用户的集合
	rc.client.SAdd("user", userID)
	// 设置某个用户的活动集合
	userActivitySet := fmt.Sprintf("user:%d", userID)
	rc.client.SAdd(userActivitySet, activityID)
	// 活动的手机号码状态集合
	activityPhoneSet := fmt.Sprintf("user_%d_activity_%d", userID, activityID)
	// 手机号码和状态
	phoneStatus := fmt.Sprintf("%s,%s", message.PhoneNum, message.State)
	rc.client.SAdd(activityPhoneSet, phoneStatus)
}

func (rc *redisContainer) Print() {
	for _, userID := range rc.client.SMembers("user").Val() {
		fmt.Printf("user:%s \n", userID)
		userActivitySet := fmt.Sprintf("user:%s", userID)
		for _, activityID := range rc.client.SMembers(userActivitySet).Val() {
			fmt.Printf("activityID:%s  \n", activityID)
			activityPhoneSet := fmt.Sprintf("user_%s_activity_%s", userID, activityID)
			for _, phoneStatus := range rc.client.SMembers(activityPhoneSet).Val() {
				fmt.Printf("phoneStatus:%s  ", phoneStatus)
			}
			fmt.Printf("\n")
		}

	}
}

func (rc *redisContainer) EchoErrorPhones() error {
	return nil
}

func (rc *redisContainer) EchoStatisticData() error {
	return nil
}
