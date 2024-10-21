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
)

type DepartmentService struct {
	MinPage int
	MaxPage int
	Url     string
	DB      *gorm.DB
}

func (ds *DepartmentService) Run() {

	ds.DB = db.Conn.GetDB("default")

	dl := make([]model.Department, 0)
	for i := ds.MinPage; i < ds.MaxPage; i++ {
		res := ds.request(i)
		deps, err := ds.parseHTML(res)
		if err != nil {
			log.Fatal(err)
			return
		}
		dl = append(dl, deps...)
	}
	if len(dl) > 0 {
		ds.bulkUpsert(dl)
	}
}

func (ds *DepartmentService) request(page int) string {
	url := fmt.Sprintf(ds.Url, page)
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
		"Referer":            "https://data.10jqka.com.cn/market/longhu/",
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
	return body
}

// 解析 HTML 表格并提取数据
func (ds *DepartmentService) parseHTML(html string) ([]model.Department, error) {
	departments := make([]model.Department, 0)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return departments, fmt.Errorf("解析 HTML 失败: %v", err)
	}

	doc.Find("table.m-table tbody tr").Each(func(i int, row *goquery.Selection) {
		var dept model.Department

		// 提取每一列的数据，去除多余空格
		//dept.Name = strings.TrimSpace(row.Find("td").Eq(1).Text())
		nameCell := row.Find("td").Eq(1)
		if linkTag := nameCell.Find("a"); linkTag.Length() > 0 {
			dept.Url, _ = linkTag.Attr("href")
			orgCode, _ := util.ExtractOrgCode(dept.Url)
			dept.OrgCode = orgCode
			dept.Name, _ = linkTag.Attr("title")
		}
		dept.Appearances = util.ParseInt(strings.TrimSpace(row.Find("td").Eq(2).Text()))
		dept.FundsUsed = strings.TrimSpace(row.Find("td").Eq(3).Text())
		dept.AnnualAppearances = util.ParseInt(strings.TrimSpace(row.Find("td").Eq(4).Text()))
		dept.AnnualStocks = util.ParseInt(strings.TrimSpace(row.Find("td").Eq(5).Text()))
		dept.SuccessRate = strings.TrimSpace(row.Find("td").Eq(6).Text())

		departments = append(departments, dept)
	})

	return departments, nil
}

// 批量入库
func (ds *DepartmentService) bulkUpsert(departments []model.Department) {
	result := ds.DB.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "org_code"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "appearances", "funds_used", "annual_appearances", "annual_stocks", "success_rate"}),
		},
	).CreateInBatches(departments, 100)

	if result.Error != nil {
		fmt.Println("批量插入/更新失败:", result.Error)
	} else {
		fmt.Printf("成功插入/更新 %d 条记录\n", result.RowsAffected)
	}
}

func (ds *DepartmentService) GetAll() []model.Department {
	var departments []model.Department
	result := ds.DB.Find(&departments)
	if result.Error != nil {
		log.Fatalln(result.Error)
	}
	return departments
}
