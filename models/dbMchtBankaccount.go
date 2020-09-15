package models

type Tbl_mcht_bankaccount struct {
	OWNER_CD          string `gorm:"column:owner_cd"`
	NAME           string `gorm:"column:name"`
	ACCOUNT           string `gorm:"column:account"`
	UC_BC_CD           string `gorm:"column:uc_bc_cd"`
	BANK_CODE           string `gorm:"column:bank_code"`
	BANK_NAME           string `gorm:"column:bank_name"`
	REC_CRT_TS           string `gorm:"column:rec_crt_ts"`
	REC_UPD_TS           string `gorm:"column:rec_upd_ts"`
}

func (t *Tbl_mcht_bankaccount) TableName() string {
	return "tbl_mcht_bankaccount"
}