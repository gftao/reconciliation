package models

import "time"

type DbOrderInfo struct {
	Prod_cd      string `gorm:"type:varchar(15);primary_key"`
	Biz_cd       string `gorm:"type:varchar(8);primary_key"`
	Trans_cd     string `gorm:"type:varchar(8);primary_key"`
	User_name    string `gorm:"type:varchar(100)"`
	User_passwd  string `gorm:"type:varchar(512)"`
	Prod_nm      string `gorm:"type:varchar(40)"`
	ServerPubKey string `gorm:"type:text(5000)"`
	ServerPriKey string `gorm:"type:text(5000)"`
	RecUpdTs     time.Time
	RecCrtTs     time.Time
	Active_flg   string `gorm:"type:varchar(1)"`
	Sign_flg     string `gorm:"type:varchar(1)"`
	Time_out     string `gorm:"type:varchar(10)"`
	Http_url     string `gorm:"type:varchar(100)"`
	Remark1      string `gorm:"type:varchar(100)"`
	Remark2      string `gorm:"type:varchar(100)"`
}

func (t DbOrderInfo) TableName() string {
	return "order_infos"
}
