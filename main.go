package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)


type MessageSender interface {
	Send(roomId string, msg string) error
}

type DingtalkMessageSender struct {
	url string
}

func NewDingtalkMessageSender (dingtalkUrl string) (*DingtalkMessageSender, error) {
	var s = DingtalkMessageSender{url: dingtalkUrl}
	return &s, nil
}

var _ MessageSender = (*DingtalkMessageSender)(nil)

func (d *DingtalkMessageSender) Send(roomId string, msg string) error {
	fmt.Println("hello,world,start send msg to DingDingRobot")
	client := &http.Client{}
	//拼接json参数
	data1:=make(map[string]interface{})
	data1["msgtype"]="text"
	data2:=make(map[string]interface{})
	data2["isAtAll"]=false
	data1["at"]=data2
	data3:=make(map[string]interface{})
	data3["content"]=msg
	data1["text"] = data3

	bytesData, err :=json.Marshal(data1)
	if err != nil {
		return err
	}

	req,_:=http.NewRequest(http.MethodPost,d.url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")
	resp,_:=client.Do(req)
	defer resp.Body.Close()
	body,_:= ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
	return nil
}



type WeatherReporter interface {
	See() error
}

type WeatherReporterImpl struct {

	s MessageSender
}

var _ WeatherReporter = (*WeatherReporterImpl)(nil)


func NewWeatherReporterImpl(s MessageSender) (*WeatherReporterImpl, error) {
	return &WeatherReporterImpl{
		s: s,
	}, nil
}

func (i *WeatherReporterImpl) See() error {
	// todo see
	msg := "提醒 天要下雨"
	err := i.s.Send("mock", msg)
	if err != nil {
		return err
	}
	return nil
}


func main() {

	sender, err := NewDingtalkMessageSender("https://oapi.dingtalk.com/robot/send?access_token=ffaabe93a835ff732b8053c0cd54c1e8315a8f906ddc0cc722dad5e833ff281c")
	if err != nil {
		log.Fatalln("sender init failed: " + err.Error())
	}

	r, err := NewWeatherReporterImpl(sender)
	if err != nil {
		log.Fatalln("123")
	}

	err = r.See()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return
	//消息URL

}