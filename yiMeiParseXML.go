package handleFeedback

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// YiMeiXMLParser 解析亿美回执的结构体
type YiMeiXMLParser struct {
	data      YiMeiXMLRoot
	container Container
}
type YiMeiXMLRoot struct {
	Root []YiMeiXMLData `xml:"response"`
}
type YiMeiXMLData struct {
	Messages  []YiMeiMessage `xml:"message"`
	ErrorCode int            `xml:"error"`
}

type YiMeiMessage struct {
	PhoneNum string `xml:"srctermid"`
	State    string `xml:"state"`
	Seqid    int    `xml:"seqid"`
}

// NewYiMeiXML 生成解析亿美的xml回执文件。method参数用于选择保存解析好的数据的容器。
func NewYiMeiXML(method int) XMLParser {
	var container Container
	switch method {
	case 1:
		container = NewMapContainer()
	case 2:
		container = NewRidesContainer()
	}
	return &YiMeiXMLParser{
		// data:      &YiMeiXMLData{},
		container: container,
	}
}

// PrintOriginData 将解析的初始化数据打印出来
func (ym *YiMeiXMLParser) PrintOriginData() {
	fmt.Println(ym.data)
}

// Parse 将数据简单的解析xml文件。数据存在内存中。
func (ym *YiMeiXMLParser) Parse(path string) error {
	file, err := os.Open(path)
	var errMsg string
	if err != nil {
		errMsg = fmt.Sprintf("open file error: %v", err)
		return errors.New(errMsg)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		errMsg = fmt.Sprintf("read file error: %v", err)
		return errors.New(errMsg)
	}
	err = xml.Unmarshal(data, &ym.data)
	if err != nil {
		errMsg = fmt.Sprintf("unmarshal data error: %v", err)
		return errors.New(errMsg)
	}
	return nil
}

// StorePhones 存储并统计手机号码的状态，并将结果存到容器中.重复的内容会替换
func (ym *YiMeiXMLParser) StorePhones() error {
	var errMsg string
	for _, data := range ym.data.Root {
		if data.ErrorCode != 0 {
			errMsg = fmt.Sprintf("xml文件中的code错误，code为：%d", data.ErrorCode)
			return errors.New(errMsg)
		}
		for _, message := range data.Messages {
			// fmt.Printf("push: %v \n", message)
			ym.container.Push(message)
		}
	}

	return nil
}

// PrintStatisticData 强统计好的数据输出出来
func (ym *YiMeiXMLParser) PrintStatisticData() {
	ym.container.Print()
}

// EchoStatisticData 将XML文件的统计结果输出到文件中结果。
func (ym *YiMeiXMLParser) EchoStatisticData() error {
	return ym.container.EchoStatisticData()
}

//EchoErrorPhones 输出手机号码到特定的目录
func (ym *YiMeiXMLParser) EchoErrorPhones() error {
	return ym.container.EchoErrorPhones()
}
