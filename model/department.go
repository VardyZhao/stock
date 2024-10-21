package model

import "gorm.io/gorm"

type Department struct {
	gorm.Model
	Name              string `gorm:"type:varchar(50);not null;comment:营业部名称"`
	OrgCode           string `gorm:"type:varchar(50);not null;comment:营业部标识;unique"`
	Url               string `gorm:"type:varchar(1000);not null;comment:营业部链接;"`
	Appearances       int    `gorm:"not null;default:0;comment:上榜次数"`
	FundsUsed         string `gorm:"type:varchar(50);not null;default:'';comment:合计动用资金"`
	AnnualAppearances int    `gorm:"not null;default:0;comment:年内上榜次数"`
	AnnualStocks      int    `gorm:"not null;default:0;comment:年内买入股票只数"`
	SuccessRate       string `gorm:"not null;default:'';comment:年内3日跟买成功率"`
}
