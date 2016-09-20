package main

import (
	"flag"
	"fmt"
	"github.com/chanxuehong/wechat/mp"
	"github.com/chanxuehong/wechat/mp/message/request"
	"github.com/chanxuehong/wechat/mp/message/response"
	"github.com/chanxuehong/wechat/util"
	"github.com/codegangsta/negroni"
	"github.com/pelletier/go-toml"
	"log"
	"net/http"
	"time"
)

var (
	// TomlTree toml file tree
	TomlTree *toml.TomlTree
	// AesKey vchat api aeskey
	AesKey string
	// OriID vchat api oriid
	OriID string
	// Token vchat api token
	Token string
	// AppID vchat api appid
	AppID string
	// ListenPort vchat service port
	ListenPort string
)

func init() {
	flag.StringVar(&ListenPort, "p", "3000", "negroni listen port")
	flag.Parse()

	TomlTree, _ = toml.LoadFile("config.toml")

	AesKey = TomlTree.Get("vchat.aesKey").(string)
	OriID = TomlTree.Get("vchat.oriId").(string)
	Token = TomlTree.Get("vchat.token").(string)
	AppID = TomlTree.Get("vchat.appId").(string)
}

func main() {
	mux := http.NewServeMux()

	SetupBaseRouter(mux)
	SetupWeChatRouter(mux)

	n := negroni.Classic()
	n.UseHandler(mux)
	_port := ":" + ListenPort
	fmt.Println("negroni listen port : " + ListenPort)
	http.ListenAndServe(_port, n)
}

// SetupBaseRouter base router
func SetupBaseRouter(m *http.ServeMux) {
	m.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "OK!")
	})
	m.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pong")
	})
}

// SetupWeChatRouter vchat router
func SetupWeChatRouter(m *http.ServeMux) {
	m.Handle("/vxapi", getMpServerFrontend())
}

func getMpServerFrontend() *mp.ServerFrontend {
	_aesKey, err := util.AESKeyDecode(AesKey) // 这里 encodedAESKey 改成你自己的参数
	if err != nil {
		panic(err)
	}

	messageServeMux := mp.NewMessageServeMux()
	messageServeMux.MessageHandleFunc(request.MsgTypeText, TextMessageHandler)   // 注册文本处理 Handler
	messageServeMux.MessageHandleFunc(request.MsgTypeVoice, VoiceMessageHandler) // 注册处理Voice Handler

	messageServeMux.EventHandleFunc(request.EventTypeSubscribe, EventSubscribeHandler)

	messageServeMux.DefaultEventHandleFunc(DefaultEventHandler)

	// 下面函数的几个参数设置成你自己的参数: oriId, token, appId
	mpServer := mp.NewDefaultServer(OriID, Token, AppID, _aesKey, messageServeMux)

	mpServerFrontend := mp.NewServerFrontend(mpServer, mp.ErrorHandlerFunc(ErrorHandler), nil)
	return mpServerFrontend
}

// ErrorHandler error handler
func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Println(err.Error())
}

// TextMessageHandler 文本消息的 Handler
func TextMessageHandler(w http.ResponseWriter, r *mp.Request) {
	text := request.GetText(r.MixedMsg) // 可以省略, 直接从 r.MixedMsg 取值

	var resp *response.Text
	var kf *response.TransferToCustomerService
	switch text.Content {
	case "hello":
		resp = response.NewText(text.FromUserName, text.ToUserName, time.Now().Unix(), "Hi")
	case "auto":
		kf = response.NewTransferToCustomerService(text.FromUserName, text.ToUserName, time.Now().Unix(), "")
		mp.WriteRawResponse(w, r, kf)
		return
	default:
		resp = response.NewText(text.FromUserName, text.ToUserName, time.Now().Unix(), "Hello,World")
	}

	mp.WriteRawResponse(w, r, resp) // 明文模式
	//mp.WriteAESResponse(w, r, resp) // 安全模式
}

// VoiceMessageHandler Voice消息的 Handler
func VoiceMessageHandler(w http.ResponseWriter, r *mp.Request) {
	text := request.GetText(r.MixedMsg) // 可以省略, 直接从 r.MixedMsg 取值
	resp := response.NewText(text.FromUserName, text.ToUserName, time.Now().Unix(), "Hello, Voice")
	mp.WriteRawResponse(w, r, resp)
}

// EventSubscribeHandler Subscribe Event 的 Handler
func EventSubscribeHandler(w http.ResponseWriter, r *mp.Request) {
	text := request.GetText(r.MixedMsg) // 可以省略, 直接从 r.MixedMsg 取值
	// resp := response.NewText(text.FromUserName, text.ToUserName, time.Now().Unix(), "Welcome ...")
	var _articles []response.Article

	var _article response.Article
	_article.Title = "Title"
	_article.Description = "This is a DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription"
	_article.URL = "http://green-comb.com"
	_article.PicURL = "http://k.sinaimg.cn/n/sports/transform/20160201/iIJ7-fxnzpkx5758643.jpg/w57012f.jpg"

	_articles = append(_articles, _article)

	resp := response.NewNews(text.FromUserName, text.ToUserName, time.Now().Unix(), _articles)
	mp.WriteRawResponse(w, r, resp)
}

// DefaultEventHandler Subscribe Event 的 Handler
func DefaultEventHandler(w http.ResponseWriter, r *mp.Request) {
	text := request.GetText(r.MixedMsg)

	switch r.MixedMsg.Event {
	case "CLICK":
		switch r.MixedMsg.EventKey {
		case "lastrss":
			resp := response.NewText(text.FromUserName, text.ToUserName, time.Now().Unix(), "lastrss\n\nThis is lastrss\n\ntest")
			mp.WriteRawResponse(w, r, resp)
		case "lastwork":
			resp := response.NewText(text.FromUserName, text.ToUserName, time.Now().Unix(), "lastwork\n\nThis is lastwork\n\ntest")
			mp.WriteRawResponse(w, r, resp)
		case "lastsearchtag":
			resp := response.NewText(text.FromUserName, text.ToUserName, time.Now().Unix(), "lastsearchtag\n\nThis is lastsearchtag\n\ntest")
			mp.WriteRawResponse(w, r, resp)
		case "cooperation":
			resp := response.NewText(text.FromUserName, text.ToUserName, time.Now().Unix(), "cooperation\n\nThis is cooperation\n\ntest")
			mp.WriteRawResponse(w, r, resp)
		case "aboutus":

			var _articles []response.Article

			var _article response.Article
			_article.Title = "Title"
			_article.Description = "This is a DescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescriptionDescription"
			_article.URL = "http://green-comb.com"
			_article.PicURL = "http://k.sinaimg.cn/n/sports/transform/20160201/iIJ7-fxnzpkx5758643.jpg/w57012f.jpg"

			_articles = append(_articles, _article)

			resp := response.NewNews(text.FromUserName, text.ToUserName, time.Now().Unix(), _articles)
			mp.WriteRawResponse(w, r, resp)
		case "myopenid":
			resp := response.NewText(text.FromUserName, text.ToUserName, time.Now().Unix(), text.FromUserName)
			mp.WriteRawResponse(w, r, resp) // 明文模式
		default:

		}
	}
}
