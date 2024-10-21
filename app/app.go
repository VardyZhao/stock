package app

import (
	"stock/config"
	"stock/db"
	"stock/service"
	"stock/util"
	"time"
)

func Init() {
	// 加载环境变量
	config.LoadEnv()
	// 加载配置
	config.LoadConfig(config.Env.RootDir + config.Env.Separate + "config.yaml")
	// 设置日志级别
	util.BuildLogger(config.Conf.GetString("log_level"))
	// 连接数据库
	db.Load()
}

func Run() {
	ds := service.DepartmentService{
		MinPage: 1,
		MaxPage: 10,
		Url:     "https://data.10jqka.com.cn/ifmarket/lhbyyb/type/1/tab/sbcs/field/sbcs/sort/desc/page/%d/",
	}
	ds.Run()
	departments := ds.GetAll()
	if len(departments) > 0 {
		for _, dept := range departments {
			currentDate := time.Now()
			t2, _ := time.Parse("2006-01-02", "2024-10-16")
			currentDate = t2
			ss := service.StockService{
				Url:         "http://data.10jqka.com.cn/ifmarket/lhbhistory/orgcode/%s/field/ENDDATE/order/desc/page/%d/",
				Referer:     dept.Url,
				OrgCode:     dept.OrgCode,
				QueryFlag:   true,
				CollectDate: currentDate,
			}
			ss.Run()
		}
	}
}
