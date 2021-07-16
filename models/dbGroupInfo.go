package models

type Tbl_group_info struct {
	SHOP_ACCOUNT_ID string `gorm:"column:SHOP_ACCOUNT_ID"`
	MCHT_CD         string `gorm:"column:MCHT_CD"`
	SHOP_ACCOUNT    string `gorm:"column:SHOP_ACCOUNT"`
	GROUP_NAME      string `gorm:"column:GROUP_NAME"`
	TERM_ID         string `gorm:"column:TERM_ID"`
	REC_OPR_ID      string `gorm:"column:REC_OPR_ID"`
	REC_UPD_OPR     string `gorm:"column:REC_UPD_OPR"`
	REC_CRT_TS      string `gorm:"column:REC_CRT_TS"`
	REC_UPD_TS      string `gorm:"column:REC_UPD_TS"`
	SHOP_NAME       string `gorm:"column:SHOP_NAME"`
}

func (t *Tbl_group_info) TableName() string {
	return "tbl_group_info"
}
