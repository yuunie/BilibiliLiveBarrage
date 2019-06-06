package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"net/http"
	"net/url"
	"io/ioutil"
	"fmt"
	"time"
	"encoding/json"
	"reflect"
	"strconv"
)



var (
	mw *walk.MainWindow
	msgUrl string = "https://api.live.bilibili.com/ajax/msg"
	Status bool = false
	infoBox *walk.TextEdit
	roomId *walk.LineEdit
	startBtn *walk.PushButton
)
type Song struct {
	PlayStatus	bool
	User		string
	Name		string
	Singer		string
}

func start() bool {

	// 检测是否输入房间号
	var id string = roomId.Text();
	if id == "" {
		fmt.Println("没有获得RoomId")
		walk.MsgBox(mw, "信息", "请输入直播间ID", walk.MsgBoxIconInformation)
		return false
	}

	// 切换按钮文本和状态
	startBtn.SetText("停止")
	Status = true

	go listenMsg(id)

	return true
}

func listenMsg(id string) bool {
	
	for {
		if !Status {
			break;
		}
		
		// test
		res, err := http.PostForm(msgUrl, url.Values{"roomid": {id}})
		if err != nil {
			fmt.Println(err)
			return false
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer res.Body.Close()

		// fmt.Println(string(body))
		msgJson := string(body)

		// infoBox.SetText(msgJson)

		// 解析JSON到MAP
		m := make(map[string]interface{})
		err = json.Unmarshal([]byte(msgJson), &m)
		
		if err != nil {
			fmt.Println(err)
			return false
		}
		// fmt.Println(m["data"].(map[string]interface{})["admin"].([]interface{})[1].(map[string]interface{})["text"])
		msgStatus 	:= m["code"].(float64)
		strmsgStatus := strconv.FormatFloat(float64(msgStatus), 'f', 0, 64)
		
		if strmsgStatus != "0" {
			walk.MsgBox(mw, "错误", "获取数据错误", walk.MsgBoxIconInformation)
		}
		msgAdminMap := m["data"].(map[string]interface{})["admin"].([]interface{})
		msgroomMap := m["data"].(map[string]interface{})["room"].([]interface{})

		_, _ = msgAdminMap, msgroomMap

		var text string;
		// for _, v := range msgAdminMap {
		// 	text += v.(map[string]interface{})["text"].(string) + "\r\n"
		// }
		
		for _, v := range msgroomMap {
			content		:= v.(map[string]interface{})["text"].(string)
			nickname 	:= v.(map[string]interface{})["nickname"].(string)
			uid 		:= v.(map[string]interface{})["uid"].(float64)
			timeline 	:= v.(map[string]interface{})["timeline"].(string)
			isadmin 	:= v.(map[string]interface{})["isadmin"].(float64)
			vip 		:= v.(map[string]interface{})["vip"].(float64)
			svip 		:= v.(map[string]interface{})["svip"].(float64)
			// _, _, _, _, _, _, _ = content, nickname, uid, timeline, isadmin, vip, svip
			// text += timeline + " " + nickname + " (" + uid + ") "
			fmt.Println(reflect.TypeOf(uid))
			// float转string保留0位小数
			struid := strconv.FormatFloat(float64(uid), 'f', 0, 64)
			strisadmin := strconv.FormatFloat(float64(isadmin), 'f', 0, 64)
			strvip := strconv.FormatFloat(float64(vip), 'f', 0, 64)
			strsvip := strconv.FormatFloat(float64(svip), 'f', 0, 64)

			text += nickname + " (" + struid + ") "
			if strisadmin == "1" {
				text += " [管理员] "
			}
			if strvip == "1" {
				text += " [VIP] "
			}
			if strsvip == "1" {
				text += " [SVIP] "
			}
			text += "[" + timeline +"]"
			text += ":\r\n"

			text += content + "\r\n\r\n\r\n"
		}
		if text != "" {
			infoBox.SetText(text)
		} else {
			// infoBox.SetText("当前直播间没有弹幕或直播间ID错误")
		}
		
		// fmt.Println("for")
		// 1纳秒 = 0.000 000 001秒
		time.Sleep(1000000000)
	}


	return true
}

func stop() {
	infoBox.SetText("")
	// 切换按钮文本和状态
	startBtn.SetText("开始")
	Status = false
}

func main() {

	MainWindow{
		AssignTo: &mw,
		Title:   "BilibiliLive - V0.1Bate",
		MinSize: Size{500, 600},
		Layout:  VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					Label{
						Text: "直播间ID:",
					},
					LineEdit{
						AssignTo: &roomId,
					},	
					PushButton{
						AssignTo: &startBtn,
						Text: "开始",
						OnClicked: func() {
							if Status == false {
								start()
							} else {
								stop()
							}
						},
					},
					PushButton{
						Text: "关于",
						OnClicked: func() {
							walk.MsgBox(mw, "关于", "BilibiliLive - v0.1\r\nAuthor: Unie Yu", walk.MsgBoxIconInformation)
						},
					},
				},
			},
			TextEdit{
				AssignTo: &infoBox,
				ReadOnly: true,
				Text: "输入直播间ID,点击开始按钮开始获取弹幕\r\n\r\n如 https://live.bilibili.com/6619197 的直播间ID为 6619197\r\n\r\n暂不支持短ID",
			},
			// TableView{
			// 	CheckBoxes:       true,
			// 	ColumnsOrderable: true,
			// 	MultiSelection:   true,
			// 	Columns: []TableViewColumn{
			// 		{Title: "状态"},
			// 		{Title: "点歌人"},
			// 		{Title: "歌曲"},
			// 		{Title: "歌手"},
			// 	},
			// },
			
		},
	}.Run()
}
