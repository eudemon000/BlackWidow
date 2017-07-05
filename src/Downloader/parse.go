package Downloader

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/henrylee2cn/pholcus/common/mahonia"
	"BlackWidow/src/utils"
	"BlackWidow/src/Database"
	"strings"
	_"golang.org/x/net/html/charset"
	_"io/ioutil"
	_"bytes"
	"regexp"
	logUtils "BlackWidow/src/logPackage"
	"fmt"
	_"container/list"
)

var cUrl string

func Parser(url string) []string {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		return []string{}
	}
	cUrl = doc.Url.String()
	var webCharset string
	headTag := doc.Find("head")
	metaTag := headTag.Find("meta")
	webCharset = checkCharset(metaTag)
	//获取页面的关键词，根据编码进行编码转换，并保存到数据库
	keyword := checkTag(metaTag, webCharset)
	fmt.Println("转码后===》", keyword)
	md5, _ := utils.Md5(url)
	urlId := Database.InsertUrls(url, md5, "YES", "1")
	fmt.Println("url表id", urlId)
	if urlId > 0 {
		//从待爬取表里删除
		Database.RemoveWaitUrl(url)
	}

	//获取接下来需要爬取的URL，放放入队列中
	bodyTag := doc.Find("body")
	//var resultUrl list.List
	var resultUrl []string
	bodyTag.Each(func(i int, bodySelect *goquery.Selection) {
		resultUrl = findUrls(bodySelect)
	})
	return resultUrl

}

//检查页面的编码类型
func checkCharset(sele *goquery.Selection) (webCharset string) {
	//var webCharset string
	defer func() {
		if err := recover(); err != nil {
			//fmt.Println(err)
			logUtils.Msg(logUtils.Error, err)
		}
	}()
	sele.Each(func(i int, m *goquery.Selection) {
		var wOK bool
		webCharset, wOK = m.Attr("charset")
		if !wOK {
			httpEquiv, hOk := m.Attr("http-equiv")
			if hOk && httpEquiv == "Content-Type" {
				content, _ := m.Attr("content")
				webCharset = content
				//fmt.Println(webCharset)
				panic(webCharset)
			}
		} else {
			panic(webCharset)
		}


	})
	//fmt.Println("charset===>", webCharset)
	return
}

//检查meta信息
func checkTag(sele *goquery.Selection, webCharset string) string {
	var tag string
	sele.Each(func(i int, m *goquery.Selection) {

		result, ok := m.Attr("name")
		if ok {
			if result == "keywords" || result =="Keywords" || result == "description" || result == "Description" {
				content, _ := m.Attr("content")
				//fmt.Println(content)
				tag = formatStr(content, webCharset)
				//tag = content
				fmt.Println("content===>", content)
				/*if content != "" {
					content = formatStr(content, webCharset)
					err := sqlConn.InsertTag(content, url)
					if err != nil {
						fmt.Println(err)
					}
					//fmt.Println(content, err)
				}*/
				//return tag
			}
		}

	})

	return tag
}

func formatStr(str, setCharset string) string {
	setCharset = strings.ToLower(setCharset)
	if strings.Contains(setCharset, "gbk") {
		de := mahonia.NewDecoder("gbk")
		result := de.ConvertString(str)
		//result := Decode(str, "gbk")
		return result
	} else if strings.Contains(setCharset, "gb2312") {
		de := mahonia.NewDecoder("gb2312")
		result := de.ConvertString(str)
		//result := Decode(str, "gb2312")
		return result

	}
	return str

}

//获取页面上所有的URL
func findUrls(bodySelect *goquery.Selection) []string {
	//var array list.List = list.List{}
	var array []string = make([]string, 0)
	//var in int = 0
	aTag := bodySelect.Find("a")
	aTag.Each(func(index int, node *goquery.Selection) {
		tempUrl, ok := node.Attr("href")
		if ok {
			//此处暂时判断链接以http开头，未来需要判断相对地址，暂时不做处理
			/*if strings.Index(tempUrl, "http") != -1 {
				array = append(array, tempUrl)
			}*/
			result, ok := Format(tempUrl)
			if ok {
				array = append(array, result)
				//array.PushBack(result)
			}
			//manage.PushData(tempUrl)
		}
	})
	return array
}

func Format(str string) (result string, ok bool) {
	//fmt.Println("接口方法调用", str)
	//首先判断是不是是不是javascript，#或*开头的,如果是代表不是合法URL
	ok, err := regexp.MatchString("^javascript|^#|^\\*", str)
	if err !=nil {
		logUtils.Msg(logUtils.Error, err)
		return "", false
	}
	if ok {
		return "", false
	}

	//判断是不是http开头的，http和https均可判断
	ok, err = regexp.MatchString("^http", str)
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		return "", false
	}
	if ok {
		/*lastIndex := strings.LastIndex(str, "/")
		if lastIndex != -1 && lastIndex == len(str) {
			str = str[:lastIndex - 1]
		}*/
		return str, true
	}

	//还要一种是相对路径，分两种情况，1、"/"开头；2、非"/"开头
	ok, err = regexp.MatchString("^/{1}[a-zA-Z0-9]{1,}?", str)
	if ok {
		//需要找路径根
		strs := strings.Split(cUrl, "/")
		//fmt.Println("当前路径", cUrl, strs)
		re := strs[0] + "//" + strs[2] + str
		return re, false
	}

	ok, err = regexp.MatchString("[a-zA-Z0-9]{1,}?", str)
	if err != nil {
		logUtils.Msg(logUtils.Error, err)
		return "", false
	}
	if ok {
		postion := strings.LastIndex(cUrl, "/")
		postion += 1
		a := cUrl[0:postion]
		re := a + str
		return re, true
	}
	return "", false
}
