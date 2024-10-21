package model

import "gorm.io/gorm"

type Stock struct {
	gorm.Model
	Name             string `gorm:"type:varchar(50);not null;comment:股票名称;index:unique-name-date-org,unique;priority:1"`
	Url              string `gorm:"type:varchar(1000);not null;comment:股票链接;"`
	OnListDate       string `gorm:"type:varchar(50);not null;comment:上榜日期;index:unique-name-date-org,unique;priority:2"`
	OrgCode          string `gorm:"type:varchar(50);not null;comment:营业部标识;index:unique-name-date-org,unique;priority:3"`
	OnListReason     string `gorm:"type:varchar(255);not null;comment:上榜原因"`
	PriceChange      int    `gorm:"not null;default:0;comment:涨跌幅(%)"`
	AmountBought     int    `gorm:"not null;default:0;comment:买入额（万）"`
	AmountSold       int    `gorm:"not null;default:0;comment:卖出额（万）"`
	NetTradingAmount int    `gorm:"not null;default:0;comment:买卖净额（万）"`
	Sector           string `gorm:"type:varchar(50);not null;default:'';comment:所属板块"`
}
