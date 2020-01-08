package models

import "time"

type Tbl_mcht_recon_list struct {
	MCHT_CD    string `gorm:"column:MCHT_CD"`
	USER       string `gorm:"column:user"`
	PASSWD     string `gorm:"column:passwd"`
	HOST       string `gorm:"column:host"`
	PORT       string `gorm:"column:port"`
	REMOTE_DIR string `gorm:"column:remote_dir"`
	Mcht_ty    string `gorm:"column:mcht_ty"`
	Transp_ty  string `gorm:"column:transp_ty"`
	EXT1       string `gorm:"column:EXT1"`
	EXT2       string `gorm:"column:EXT2"`
	EXT3       string `gorm:"column:EXT3"`
	EXT4       string `gorm:"column:EXT4"`
	EXT5       string `gorm:"column:EXT5"`
}

func (t Tbl_mcht_recon_list) TableName() string {
	return "tbl_mcht_recon_list"
}

type TBL_HOLI_INF struct {
	ID          int       `gorm:"column:ID"`
	START_DATE  string    `gorm:"column:START_DATE"`
	END_DATE    string    `gorm:"column:END_DATE"`
	UNION_FLAG  string    `gorm:"column:UNION_FLAG"`
	HOLIDAY_DSP string    `gorm:"column:HOLIDAY_DSP"`
	REC_OPR_ID  string    `gorm:"column:REC_OPR_ID"`
	REC_UPD_OPR string    `gorm:"column:REC_UPD_OPR"`
	REC_CRT_TS  time.Time `gorm:"column:REC_CRT_TS"`
	REC_UPD_TS  time.Time `gorm:"column:REC_UPD_TS"`
}

func (t TBL_HOLI_INF) TableName() string {
	return "tbl_holi_inf"
}
