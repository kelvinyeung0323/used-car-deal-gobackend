package crawler

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"log"
	"time"
)

type DonCheDiCrawler struct {
}

func (crawler *DonCheDiCrawler) init() {

	//创建收集器
	c := colly.NewCollector(
		func(collector *colly.Collector) {
			extensions.RandomUserAgent(collector) //随机浏览器UA
		},
		colly.MaxDepth(2),
		colly.Async(true),                         //异步
		colly.AllowedDomains("www.dongchedi.com"), //设置允许请求的域名主机，可以多个，Visit只会发起这些域名下的请求
		colly.IgnoreRobotsTxt(),                   //忽略robots协议
		colly.CacheDir("cache-dir"),               //设置GET请求本地缓存文件夹
		colly.AllowURLRevisit())

	//并发请求设置
	c.Limit(&colly.LimitRule{DomainGlob: "*dongchedi.*", //匹配URL包含
		Parallelism: 10,               //并发请求 10
		RandomDelay: 5 * time.Second}) //设置发起请求随机延时0-5

	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("", func(element *colly.HTMLElement) {

	})

	err := c.Visit("")
	log.Printf("Error:%v", err)
}
