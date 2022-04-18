package services

import (
	"encoding/json"
	"fmt"
	goisbn "github.com/abx123/go-isbn"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

// 极客API,用于获得中文书籍ISBN信息
var Jikeapikey string

type bookRetrieveResult struct {
	Name        string `json:"name"`        //书名
	Author      string `json:"author"`      //作者
	AuthorIntro string `json:"authorIntro"` //作者简介
	PhotoUrl    string `json:"photoUrl"`    //图片封面
	Publishing  string `json:"publishing"`  //出版社
	Published   string `json:"published"`   //出版时间
	Description string `json:"description"` //图书简介
	DoubanScore int    `json:"doubanScore"` //豆瓣评分
}
type bookRetriveHTTPResult struct {
	Ret  int                `json:"ret"`
	Msg  string             `json:"msg"`
	Data bookRetrieveResult `json:"Data"`
}

var plainProxyURL string = "127.0.0.1:6666"
var proxyURL, _ = url.Parse(plainProxyURL)

var myClient = &http.Client{
	Timeout:   10 * time.Second,
	Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		log.Printf("HTTP Request from %v Failure:%v\n", url, err.Error())
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		bodyString := string(bodyBytes)
		log.Printf("Get non-json reply from %v\nStatusCode:%v Body:%v", url, r.StatusCode, bodyString)
	}
	return err
}

func GetMetaDataByISBN(isbn string) (bookInfo Book, err error) {
	bookInfo = Book{}
	err = nil
	isbnRetriever := goisbn.NewGoISBN(goisbn.DEFAULT_PROVIDERS)
	rawBookInfo, rerr := isbnRetriever.Get(isbn)
	if rerr != nil {
		log.Printf("ISBN %v not found in Google Books and Open Library, Checking Jike\n", isbn)
		var resp bookRetriveHTTPResult
		// FIXME: Always Get 403
		ierr := getJson(fmt.Sprintf("https://api.jike.xyz/situ/book/isbn/%v?apikey=%v", isbn, Jikeapikey), &resp)
		if ierr != nil {
			return Book{}, ierr
		}
		log.Println(resp.Ret)
		log.Println(resp.Msg)
		log.Println(resp.Data)
		//if resp.Ret != 0 {
		//	return Book{}, errors.New("book cannot be retrieved, no result")
		//}
		bookInfo.Isbn = isbn
		bookInfo.Name = resp.Data.Name
		bookInfo.Author = resp.Data.Author
		bookInfo.Language = "Chinese"
	} else {
		bookInfo.Isbn = isbn
		var authors string
		for _, subAuthor := range rawBookInfo.Authors {
			authors += fmt.Sprintf("%v,", subAuthor)
		}
		bookInfo.Author = authors
		bookInfo.Language = rawBookInfo.Language
		bookInfo.Name = rawBookInfo.Title
	}
	return bookInfo, nil
}
