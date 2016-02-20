package main

import (
	"flag"
	"fmt"
	"github.com/banerwai/vchat/util"
	"github.com/chanxuehong/wechat/mp"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"net/url"
)

var AppDir string
var _app_dir string
var TomlTree *toml.TomlTree
var CorpId, CorpSecret string

var AccessTokenServer *mp.DefaultAccessTokenServer
var MpClient *mp.Client

func init() {
	flag.StringVar(&_app_dir, "d", "default", "input app dir")
	flag.Parse()
	AppDir = fmt.Sprintf("./%s/", _app_dir)

	TomlTree, _ = toml.LoadFile(AppDir + "config.toml")

	CorpId = TomlTree.Get("vchat.corpId").(string)
	CorpSecret = TomlTree.Get("vchat.corpSecret").(string)

	AccessTokenServer = mp.NewDefaultAccessTokenServer(CorpId, CorpSecret, nil) // 一個應用只能有一個實例
	MpClient = mp.NewClient(AccessTokenServer, nil)
}

func main() {
	if err := CreateMenu(AppDir + "menu.json"); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(_app_dir + " set menu success")
}

func CreateMenu(menu_file string) (err error) {
	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/menu/create?access_token="

	_f, _ := ioutil.ReadFile(menu_file)

	_token, err := MpClient.Token()
	if err != nil {
		return
	}

	finalURL := incompleteURL + url.QueryEscape(_token)

	if err = util.PostRawJson(finalURL, _f, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
