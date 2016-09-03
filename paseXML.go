package handleFeedback

type XMLParser interface {
	PrintOriginData()
	Parse(path string) error
	StorePhones() error
	PrintStatisticData()
	EchoErrorPhones() error
	EchoStatisticData() error
}
