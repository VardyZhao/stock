package app

import (
	"stock/config"
	"stock/db"
	"stock/service"
	"stock/util"
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
		MaxPage: 13,
		Url:     "https://data.10jqka.com.cn/ifmarket/lhbyyb/type/1/tab/sbcs/field/sbcs/sort/desc/page/%d/",
	}
	ds.Run()
}
