package MqttClient

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	//"fmt"
	"log"
	//"os"
	//"path/filepath"
	"runtime"
	"testing"

	"github.com/gogf/gf/os/glog"

	"github.com/gogf/gf/frame/g"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func PathSep() string {
	if runtime.GOOS == "windows" {
		return `\`
	} else {
		return "/"
	}
}

func OnServerMsg(client mqtt.Client, msg mqtt.Message) {

}

func TestConnect(t *testing.T) {

	dir, errDir := filepath.Abs(filepath.Dir(os.Args[0]))
	if errDir != nil {
		log.Fatal(errDir)
	}

	var sPath = dir + PathSep() + "logs"
	fmt.Println("LogPath:" + sPath)

	sPath = "c:\\test.log"

	glog.SetPath(sPath)
	glog.SetStdoutPrint(true)
	glog.SetFile("log-{Ymd}.log")

	log.Print("test")

	g.Cfg().SetPath(`C:\Users\wangjinhui\go\src\github.com\JasonHonor\MqttClient`)
	g.Cfg().SetFileName("Config.ini")

	var mqttSubs map[string]mqtt.MessageHandler
	mqttSubs = make(map[string]mqtt.MessageHandler, 1)
	mqttSubs["/"] = OnServerMsg

	log.Printf("Before create. len=%d", len(mqttSubs))

	client := NewMqttClient(mqttSubs)

	msg := "test"

	client.Publish("test", &msg)

	time.Sleep(30 * time.Second)
}
