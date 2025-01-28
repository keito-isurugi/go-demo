package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func main() {
	// Collyのコレクターを作成
	c := colly.NewCollector(
		colly.AllowedDomains("example.com"), // 許可するドメインを指定
	)

	// ページを訪問した際に呼ばれるコールバック
	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("ページタイトル:", e.Text)
	})

	// リンクを見つけた際に呼ばれるコールバック
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		fmt.Println("リンク発見:", link)

		// 絶対URLなら、そのリンクも訪問
		err := c.Visit(e.Request.AbsoluteURL(link))
		if err != nil {
			log.Printf("リンク訪問中にエラー: %v\n", err)
		}
	})

	// エラーが発生した場合のコールバック
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("エラー発生:", err)
	})

	// スクレイピング開始
	startURL := "https://example.com/"
	fmt.Println("開始URL:", startURL)
	err := c.Visit(startURL)
	if err != nil {
		log.Fatalf("訪問に失敗: %v", err)
	}
}
