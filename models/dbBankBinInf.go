package models

type Tbl_bank_bin_inf struct {
	INS_ID_CD          string `gorm:"column:INS_ID_CD"`
	INS_ID_NM           string `gorm:"column:INS_ID_NM"`
	ACC_LEN           string `gorm:"column:ACC_LEN"`
	BIN_LEN           string `gorm:"column:BIN_LEN"`
	CARD_BIN           string `gorm:"column:CARD_BIN"`
	CARD_TP           string `gorm:"column:CARD_TP"`
	CARD_DIS           string `gorm:"column:CARD_DIS"`
	REC_OPR_ID           string `gorm:"column:REC_OPR_ID"`
	REC_UPD_OPR           string `gorm:"column:REC_UPD_OPR"`
	REC_CRT_TS           string `gorm:"column:REC_CRT_TS"`
	REC_UPD_TS           string `gorm:"column:REC_UPD_TS"`
}

func (t *Tbl_bank_bin_inf) TableName() string {
	return "tbl_bank_bin_inf"
}