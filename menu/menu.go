package main

import (
	"flag"
	"fmt"
	api_http "github.com/banerwai/gommon/net/http"
	"github.com/chanxuehong/wechat/mp"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"net/url"
)

var (
	// AppDir vchat app dir
	AppDir string

	// TomlTree toml file tree
	TomlTree *toml.TomlTree

	// CorpID corpID
	CorpID string

	// CorpSecret CorpSecret
	CorpSecret string

	// AccessTokenServer mp DefaultAccessTokenServer
	AccessTokenServer *mp.DefaultAccessTokenServer

	// MpClient mp Client
	MpClient *mp.Client
)

func init() {
	var _appDir string
	flag.StringVar(&_appDir, "d", "default", "input app dir")
	flag.Parse()
	AppDir = fmt.Sprintf("./%s/", _appDir)

	TomlTree, _ = toml.LoadFile(AppDir + "config.toml")

	CorpID = TomlTree.Get("vchat.corpId").(string)
	CorpSecret = TomlTree.Get("vchat.corpSecret").(string)

	AccessTokenServer = mp.NewDefaultAccessTokenServer(CorpID, CorpSecret, nil) // 一個應用只能有一個實例
	MpClient = mp.NewClient(AccessTokenServer, nil)
}

func main() {
	if err := CreateMenu(AppDir + "menu.json"); err != nil {
		fmt.Println(err)
		return
	}
}

// CreateMenu gen vchat menu
func CreateMenu(menuFile string) (err error) {
	var result mp.Error

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/menu/create?access_token="

	_f, _ := ioutil.ReadFile(menuFile)

	_token, err := MpClient.Token()
	if err != nil {
		return
	}

	finalURL := incompleteURL + url.QueryEscape(_token)

	if err = api_http.PostRawJson(finalURL, _f, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result
		return
	}
	return
}
