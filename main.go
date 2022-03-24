package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)
var wg sync.WaitGroup
func init() {
	file := "./" +"message"+"_"+string(time.Now().Format("200601021504"))+ ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[spider]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	return
}

func InitConfig(dirName string) {
	workDir, _ := os.Getwd()
	viper.SetConfigName("app")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

}
func main() {
	now := time.Now()
	dirName := string(time.Now().Format("200601021504"))
	InitConfig(dirName)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	genUrl := GenUrl()
	wg.Add(len(genUrl))
	for i :=0;i<len(genUrl);i++ {
		//log.Printf("url的列表%s",genUrl[i])
		go func(i int) {
			defer wg.Done()
			keyword := strings.Split(genUrl[i],"=")
			keyword = strings.Split(keyword[1],"&")
			keywordone := keyword[0]
			ParseHtml(genUrl[i],keywordone)
		}(i)
	}
	wg.Wait()
	fmt.Println("耗时:", time.Since(now))
}
func ParseHtml(urlstr string,keyword string)  {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Fatal(err)
	}
	c := colly.NewCollector()
	// 超时设定
	c.SetRequestTimeout(100 * time.Second)
	// 指定Agent信息
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"
	c.OnRequest(func(r *colly.Request) {
		// Request头部设定
		r.Headers.Set("Host", u.Host)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", u.Host)
		r.Headers.Set("Referer", urlstr)
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")
	})
	c.OnResponse(func(resp *colly.Response) {

		// goquery直接读取resp.Body的内容
		htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
		// 读取url再传给goquery，访问url读取内容，此处不建议使用
		// htmlDoc, err := goquery.NewDocument(resp.Request.URL.String())
		if err != nil {
			log.Fatal(err)
		}

		//<div class="result c-container xpath-log new-pmd的mu的属性值"
		htmlDoc.Find("div.result.c-container").Each(func(i int, s *goquery.Selection) {
			band, _ := s.Attr("mu")
			if band != "" {
				log.Printf("关键字:%s,结果:%s",keyword,band)
			}

		})

	})
	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})
	err = c.Visit(urlstr)


}

func GenUrl() ([]string) {
	var UrlArr []string

	keywordAll := viper.GetString("keyword")
	keywordArr := strings.Split(keywordAll,",")
	pages := viper.GetInt("pages")
	for _,keyword:= range keywordArr {
		for i:=1;i<=pages;i++{
			urlstrone := fmt.Sprintf("https://www.baidu.com/s?wd=%s&pn=%v",keyword,(i-1)*10)
			UrlArr=append(UrlArr,urlstrone)

		}
	}
	return UrlArr
}