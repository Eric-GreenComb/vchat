package main

import (
	"fmt"
	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/menu"
	"github.com/pelletier/go-toml"
)

var TomlTree *toml.TomlTree
var CorpId, CorpSecret string

var AccessTokenServer *mp.DefaultAccessTokenServer
var MpClient *mp.Client

func init() {
	TomlTree, _ = toml.LoadFile("config.toml")

	CorpId = TomlTree.Get("vchat.corpId").(string)
	CorpSecret = TomlTree.Get("vchat.corpSecret").(string)

	AccessTokenServer = mp.NewDefaultAccessTokenServer(CorpId, CorpSecret, nil) // 一個應用只能有一個實例
	MpClient = mp.NewClient(AccessTokenServer, nil)
}

func main() {
	var subButtons1 = make([]menu.Button, 3)
	subButtons1[0].SetAsViewButton("每日新文", "http://green-comb.com/")
	subButtons1[1].SetAsViewButton("往日新文", "http://www.baidu.com/")
	subButtons1[2].SetAsClickButton("关于我们", "company")

	var subButtons2 = make([]menu.Button, 5)
	subButtons2[0].SetAsViewButton("This推荐", "http://green-comb.com/")
	subButtons2[1].SetAsViewButton("This推荐2", "http://www.baidu.com/")
	subButtons2[2].SetAsViewButton("推荐3", "http://www.baidu.com/")
	subButtons2[3].SetAsViewButton("推荐4", "http://www.baidu.com/")
	subButtons2[4].SetAsViewButton("推荐5", "http://www.baidu.com/")

	var subButtons3 = make([]menu.Button, 3)
	subButtons3[0].SetAsViewButton("greencomb", "http://green-comb.com/")
	subButtons3[1].SetAsViewButton("搜索", "http://www.baidu.com/")
	subButtons3[2].SetAsViewButton("APP", "http://ssskdjf.com/download/aaa.apk")

	var mn menu.Menu
	mn.Buttons = make([]menu.Button, 3)
	mn.Buttons[0].SetAsSubMenuButton("水调新文", subButtons1)
	mn.Buttons[1].SetAsSubMenuButton("水调推荐", subButtons2)
	mn.Buttons[2].SetAsSubMenuButton("Customer", subButtons3)

	menuClient := (*menu.Client)(MpClient)
	if err := menuClient.CreateMenu(mn); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("set menu success")
}
