package handleFeedback

type Container interface {
	Push(args interface{})    // 将数据储存到容器中
	Print()                   // 输出存储的信息到终端上
	EchoErrorPhones() error   // 输出错误号码到特定的文件中
	EchoStatisticData() error // 统计错误号码数，并将结果输出到文件中
}
