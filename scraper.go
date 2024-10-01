package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Page struct {
	Id int `json:"id"`
}

type Book struct {
	Pages []Page `json:"pages"`
}

type BookInfo struct {
	Name    string
	Type    int
	Pages   []int
	PageNum int
}

func modifyUrl(originalUrl string, newParamKey string, newParamValue string) (string, error) {
	// 解析 URL
	parsedUrl, err := url.Parse(originalUrl)
	if err != nil {
		return "", err
	}

	// 设置新的查询参数
	query := url.Values{}
	query.Set(newParamKey, newParamValue)

	// 更新 URL 查询部分
	parsedUrl.RawQuery = query.Encode()
	return parsedUrl.String(), nil
}

func parseDocview(doc *goquery.Document) Book {
	var book Book // 在此声明 book 变量

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if strings.HasPrefix(text, "initDocview") {
			// 创建一个正则表达式，用来匹配括号内的内容
			re := regexp.MustCompile(`\((.*?)\)`)

			// 找到括号内的内容
			matches := re.FindStringSubmatch(text)

			if len(matches) > 1 {
				text = matches[1]
				// 使用 json.Unmarshal 来解析 JSON 字符串
				err := json.Unmarshal([]byte(text), &book)
				if err != nil {
					fmt.Println("JSON 解析错误:", err)
					return
				}

			}
		}
	})

	return book // 返回解析后的 book
}

func processDocuments(doc *goquery.Document, book *BookInfo) {
	find := doc.Find(".nodes-list").Find("tbody").Find("tr")

	if len(find.Nodes) == 0 {
		fmt.Println("文档类型：I型文档")
		book.Type = 1
		return
	}

	fmt.Println("文档类型：II型文档")
	book.Type = 2

	var start, end int

	val, exists := find.First().Find("td").Last().Find("a").Attr("href")
	if exists {
		resp, err2 := http.Get("https://docs.historyrussia.org" + val)
		if err2 != nil {
			log.Fatal(err2)
		}
		defer resp.Body.Close()
		reader, err2 := goquery.NewDocumentFromReader(resp.Body)
		if err2 != nil {
			log.Fatal(err2)
		}
		docview := parseDocview(reader)
		start = docview.Pages[0].Id
	}

	val2, exists2 := find.Last().Find("td").Last().Find("a").Attr("href")
	if exists2 {
		resp, err := http.Get("https://docs.historyrussia.org" + val2)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		reader, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		docview := parseDocview(reader)
		end = docview.Pages[len(docview.Pages)-1].Id
	}

	//fmt.Printf("起始页%d，终止页%d，共计%d页。\n", start, end, end-start+1)
	for i := start; i <= end; i++ {
		book.Pages = append(book.Pages, i)
	}

}

func processDocumentsHard(doc *goquery.Document, book *BookInfo) {
	find := doc.Find(".nodes-list").Find("tbody").Find("tr")

	if len(find.Nodes) == 0 {
		fmt.Println("文档类型：I型文档")
		book.Type = 1
		return
	}

	fmt.Println("文档类型：II型文档")
	book.Type = 2

	find.Each(func(i int, s *goquery.Selection) {
		val, exists := s.Find("td").Last().Find("a").Attr("href")
		if exists {
			resp, err := http.Get("https://docs.historyrussia.org" + val)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			reader, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			docview := parseDocview(reader)
			for page := range docview.Pages {
				book.Pages = append(book.Pages, page)
			}
		}
	})

	//val, exists := find.First().Find("td").Last().Find("a").Attr("href")
	//if exists {
	//	resp, err2 := http.Get("https://docs.historyrussia.org" + val)
	//	if err2 != nil {
	//		log.Fatal(err2)
	//	}
	//	defer resp.Body.Close()
	//	reader, err2 := goquery.NewDocumentFromReader(resp.Body)
	//	if err2 != nil {
	//		log.Fatal(err2)
	//	}
	//	docview := parseDocview(reader)
	//	start = docview.Pages[0].Id
	//}
	//
	//val2, exists2 := find.Last().Find("td").Last().Find("a").Attr("href")
	//if exists2 {
	//	resp, err := http.Get("https://docs.historyrussia.org" + val2)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	defer resp.Body.Close()
	//	reader, err := goquery.NewDocumentFromReader(resp.Body)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	docview := parseDocview(reader)
	//	end = docview.Pages[len(docview.Pages)-1].Id
	//}
	//
	////fmt.Printf("起始页%d，终止页%d，共计%d页。\n", start, end, end-start+1)
	//for i := start; i <= end; i++ {
	//	book.Pages = append(book.Pages, i)
	//}

}

func Scrape(url string, mode int) *BookInfo {
	s, err := modifyUrl(url, "per_page", "1000")
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.Get(s)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	book := BookInfo{Name: s, Type: 1, Pages: make([]int, 0)}

	// title
	name := doc.Find("h1").Text()
	if name != "" {
		book.Name = name
	}

	// 编译一个正则表达式，匹配字符串开头的数字序列
	re := regexp.MustCompile(`^\d+`)
	numStr := re.FindString(strings.TrimSpace(doc.Find(".value_of_type_11").Text()))
	pageNum, _ := strconv.Atoi(numStr)
	fmt.Println("页数：", pageNum)
	book.PageNum = pageNum

	if mode == 0 {
		processDocuments(doc, &book)
	} else {
		processDocumentsHard(doc, &book)
	}
	docview := parseDocview(doc)
	//fmt.Printf("起始页%d，结束页%d，共计%d页。\n", docview.Pages[0].Id, docview.Pages[len(docview.Pages)-1].Id, len(docview.Pages))
	for _, page := range docview.Pages {
		book.Pages = append(book.Pages, page.Id)
	}
	fmt.Println(len(book.Pages))
	sort.Ints(book.Pages)
	return &book
}

//func main() {
//	Scrape("https://docs.historyrussia.org/ru/nodes/425712-tovarisch-komsomol-1918-locale-nil-1968-dok-sezdov-konferentsiy-i-tsk-vlksm-t-2-1941-locale-nil-1968")
//	Scrape("https://docs.historyrussia.org/ru/nodes/425703-locale-nil-kogda-my-byli-molodye-locale-nil-dok-dilogiya-kn-2-studencheskie-otryady-prikamya-1965-1992-gg")
//	Scrape("https://docs.historyrussia.org/ru/nodes/405045-sovety-severo-vostoka-sssr-ch-2-1941-1961-gg")
//	Scrape("https://docs.historyrussia.org/ru/nodes/405044-sovety-severo-vostoka-sssr-ch-1-1928-1940-gg")
//Scrape("https://docs.historyrussia.org/ru/nodes/425713-boevoy-vosemnadtsatyy-god", 0)
//Scrape("https://docs.historyrussia.org/ru/nodes/404939-politbyuro-i-delo-beriya", 0)
//	Scrape("https://docs.historyrussia.org/ru/nodes/405054-zaharchenko-a-v-repinetskiy-a-i-soldatova-o-n-nkvd-i-ekonomika-v-gody-velikoy-otechestvennoy-voyny")
//Scrape("https://docs.historyrussia.org/ru/nodes/354231-politbyuro-i-delo-viktora-abakumova", 0)
//}
