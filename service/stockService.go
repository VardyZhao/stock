package service

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"stock/db"
	"stock/model"
	"stock/util"
	"strings"
	"time"
)

type StockService struct {
	Url         string
	Referer     string
	OrgCode     string
	QueryFlag   bool
	CollectDate time.Time
	DB          *gorm.DB
}

func (ds *StockService) Run() {

	ds.DB = db.Conn.GetDB("default")

	sl := make([]model.Stock, 0)
	page := 1
	for ds.QueryFlag {
		res := ds.request(page)
		stks, err := ds.parseHTML(res)
		if err != nil {
			ds.QueryFlag = false
			log.Fatal(err)
			return
		}
		sl = append(sl, stks...)
		page++
		fmt.Println(ds.OrgCode)
		fmt.Println(page)
		if !ds.QueryFlag {
			break
		}
	}
	if len(sl) > 0 {
		ds.bulkUpsert(sl)
	}
}

func (ds *StockService) request(page int) string {
	url := fmt.Sprintf(ds.Url, ds.OrgCode, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("创建请求失败: %v", err)
	}

	cookie := "searchGuide=sg; __utma=156575163.1042276024.1727017372.1727017372.1727017372.1; __utmz=156575163.1727017372.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); spversion=20130314; historystock=002166%7C*%7C603833%7C*%7C688036%7C*%7C300652; log=; Hm_lvt_722143063e4892925903024537075d0d=1727017317,1728918035; Hm_lpvt_722143063e4892925903024537075d0d=1728918035; HMACCOUNT=B30027AD0448A688; Hm_lvt_929f8b362150b1f77b477230541dbbc2=1727017317,1728918035; Hm_lpvt_929f8b362150b1f77b477230541dbbc2=1728918035; Hm_lvt_78c58f01938e4d85eaf619eae71b4ed1=1727017318,1728918035; refreshStat=off; Hm_lvt_60bad21af9c824a4a0530d5dbf4357ca=1727017396,1728918054; Hm_lvt_f79b64788a4e377c608617fba4c736e2=1727017395,1728918054; Hm_lpvt_60bad21af9c824a4a0530d5dbf4357ca=1728918342; Hm_lpvt_78c58f01938e4d85eaf619eae71b4ed1=1728918342; Hm_lpvt_f79b64788a4e377c608617fba4c736e2=1728918342; v=A6Pn321XwIdGDIyz9s8rbMCJMuxImDiOcSJ-E9UB_O95Rc2SXWjHKoH8C1Lm"
	// 设置请求头
	headers := map[string]string{
		"Accept":             "text/html, */*; q=0.01",
		"Accept-Language":    "zh-CN,zh;q=0.9",
		"Connection":         "keep-alive",
		"Cookie":             cookie,
		"Referer":            ds.Referer,
		"Sec-Fetch-Dest":     "empty",
		"Sec-Fetch-Mode":     "cors",
		"Sec-Fetch-Site":     "same-origin",
		"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36",
		"X-Requested-With":   "XMLHttpRequest",
		"hexin-v":            "A6Pn321XwIdGDIyz9s8rbMCJMuxImDiOcSJ-E9UB_O95Rc2SXWjHKoH8C1Lm",
		"sec-ch-ua":          `"Google Chrome";v="129", "Not=A?Brand";v="8", "Chromium";v="129"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"Windows"`,
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("请求发送失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("请求失败：%v", resp.Status)
	}

	body, err := util.ReadBodyWithCharset(resp)
	if err != nil {
		log.Fatalf("读取响应失败: %v", err)
	}
	ds.Referer = url
	return body
}

// 解析 HTML 表格并提取数据
func (ds *StockService) parseHTML(html string) ([]model.Stock, error) {
	stocks := make([]model.Stock, 0)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return stocks, fmt.Errorf("解析 HTML 失败: %v", err)
	}

	doc.Find("table.m-table tbody tr").Each(func(i int, row *goquery.Selection) {
		var stk model.Stock

		// 提取每一列的数据，去除多余空格
		stk.OnListDate = strings.TrimSpace(row.Find("td").Eq(0).Text())
		t2, _ := time.Parse("2006-01-02", stk.OnListDate)
		if t2.Before(ds.CollectDate) {
			ds.QueryFlag = false
			return
		}
		nameCell := row.Find("td").Eq(1)
		if linkTag := nameCell.Find("a"); linkTag.Length() > 0 {
			stk.Url, _ = linkTag.Attr("href")
			stk.Name = linkTag.Text()
		}
		stk.OnListReason = strings.TrimSpace(row.Find("td").Eq(2).Text())
		stk.PriceChange = util.ConvertFloatStrToInt(strings.TrimSpace(row.Find("td").Eq(3).Text()))
		stk.AmountBought = util.ConvertFloatStrToInt(strings.TrimSpace(row.Find("td").Eq(4).Text()))
		stk.AmountSold = util.ConvertFloatStrToInt(strings.TrimSpace(row.Find("td").Eq(5).Text()))
		stk.NetTradingAmount = util.ConvertFloatStrToInt(strings.TrimSpace(row.Find("td").Eq(6).Text()))
		stk.Sector = strings.TrimSpace(row.Find("td").Eq(7).Text())
		stk.OrgCode = ds.OrgCode

		stocks = append(stocks, stk)
	})
	if len(stocks) == 0 {
		ds.QueryFlag = false
	}

	return stocks, nil
}

// 批量入库
func (ds *StockService) bulkUpsert(stocks []model.Stock) {
	result := ds.DB.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}, {Name: "on_list_date"}, {Name: "org_code"}},
			DoUpdates: clause.AssignmentColumns([]string{"url", "on_list_date", "on_list_reason", "price_change", "amount_bought", "amount_sold", "net_trading_amount", "sector"}),
		},
	).CreateInBatches(stocks, 100)

	if result.Error != nil {
		fmt.Println("批量插入/更新失败:", result.Error)
	} else {
		fmt.Printf("成功插入/更新 %d 条记录\n", result.RowsAffected)
	}
}
