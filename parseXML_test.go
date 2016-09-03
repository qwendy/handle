package handleFeedback

import "testing"

var path = "YiMei_Report.xml"

func TestParseXML(t *testing.T) {
	ym := NewYiMeiXML(1)
	if err := ym.Parse(path); err != nil {
		t.Error("parse Err: ", err)
	}
	// ym.PrintOriginData()
	ym.StorePhones()
	ym.PrintStatisticData()
	ym.EchoErrorPhones()
	ym.EchoStatisticData()
	t.Error("stop!")
}
