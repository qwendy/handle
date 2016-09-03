package handleFeedback

import "testing"

var path = "YiMei_Report.xml"

// func TestMapParseXML(t *testing.T) {
// 	ym := NewYiMeiXML(1)
// 	if err := ym.Parse(path); err != nil {
// 		t.Error("parse Err: ", err)
// 	}
// 	// ym.PrintOriginData()
// 	ym.StorePhones()
// 	ym.PrintStatisticData()
// 	// ym.EchoErrorPhones()
// 	// ym.EchoStatisticData()
// 	// t.Error("stop!")
// }

func TestRedisParseXML(t *testing.T) {
	ym := NewYiMeiXML(2)
	if err := ym.Parse(path); err != nil {
		t.Error("parse Err: ", err)
	}
	// ym.PrintOriginData()
	ym.StorePhones()
	ym.PrintStatisticData()
	// ym.EchoErrorPhones()
	// ym.EchoStatisticData()
	// t.Error("stop!")
}
