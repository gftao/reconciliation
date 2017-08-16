package models

type Tbl_mcht_recon_list struct {
	MCHT_CD    string `gorm:"clumn:MCHT_CD"`
	USER       string `gorm:"clumn:user"`
	PASSWD     string `gorm:"clumn:passwd"`
	HOST       string `gorm:"clumn:host"`
	PORT       string `gorm:"clumn:port"`
	REMOTE_DIR string `gorm:"clumn:remote_dir"`
	MCHT_TP    string `gorm:"clumn:mcht_tp"`
}

func (t Tbl_mcht_recon_list) TableName() string {
	return "tbl_mcht_recon_list"
}