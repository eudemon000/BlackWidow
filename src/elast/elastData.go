package elast

import (
	"fmt"
	"net/http"
	"BlackWidow/src/Config"
	_"strings"
	"bytes"
	"encoding/json"
	logUtils "BlackWidow/src/logPackage"
	"strings"
	"io/ioutil"
)

type ElastType struct {
	protocol		string
	host			string
	port			string
	index			string
	elastType		string
	search			string
	analysisApi		string
	analysisName	string
}

type AllSearch struct {
	Title			string	`json:"title"`
	Content			string	`json:"content"`
	Description		string	`json:"description"`
	Keyword			string	`json:"keyword"`
	Url				string	`json:"url"`
	Md5				string	`json:"md5"`
}

type analyzer struct {
	Tokenizer		string	`json:"tokenizer"`
	Text			string	`json:"text"`
}

var e *ElastType

var elUrl string

var analyzerUrl string

func init() {
	e = new(ElastType)
	con := new(Config.Config)
	con.InitConfig()
	e.protocol = con.Read("elast", "protocol")
	e.host = con.Read("elast", "host")
	e.port = con.Read("elast", "port")
	e.index = con.Read("elast", "index")
	e.elastType = con.Read("elast", "type")
	e.search = con.Read("elast", "search")
	e.analysisApi = con.Read("elast", "analysisApi")
	e.analysisName = con.Read("elast", "analysisName")
	elUrl = formatUrl()
	analyzerUrl = formatAnalysis()
}

func formatUrl() string {
	buf := bytes.Buffer{}
	buf.WriteString(e.protocol)
	buf.WriteString("://")
	buf.WriteString(e.host)
	buf.WriteString(":")
	buf.WriteString(e.port)
	buf.WriteString("/")
	buf.WriteString(e.index)
	buf.WriteString("/")
	buf.WriteString(e.elastType)
	buf.WriteString("/")
	result := buf.String()
	fmt.Println(result)
	return result
}

func formatAnalysis() string {
	buf := bytes.Buffer{}
	buf.WriteString(e.protocol)
	buf.WriteString("://")
	buf.WriteString(e.host)
	buf.WriteString(":")
	buf.WriteString(e.port)
	buf.WriteString("/")
	buf.WriteString(e.index)
	buf.WriteString("/")
	buf.WriteString(e.analysisApi)
	result := buf.String()
	fmt.Println(result)
	return result
}

func SendElast(item *AllSearch, text string) {
	params, err := json.Marshal(item)
	fmt.Println("提交的格式===>", string(params))
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
	}
	resp, err := http.Post(elUrl + item.Md5, "application/JSON", strings.NewReader(string(params)))
	if err != nil {
		logUtils.Msg(logUtils.Error, e)
		return
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		return
	}
	fmt.Println("调用结果为:", string(result))
	defer resp.Body.Close()

	//分词
	an := new(analyzer)
	an.Tokenizer = e.analysisName
	an.Text = text
	anJson, err := json.Marshal(an)
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		return
	}

	fmt.Println("提交格式===>", string(anJson))

	anResp, err := http.Post(analyzerUrl, "application/JSON", strings.NewReader(string(anJson)))
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		return
	}

	anResult, err := ioutil.ReadAll(anResp.Body)
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		return
	}
	fmt.Println("分词结果===>", string(anResult))
	defer anResp.Body.Close()
}
