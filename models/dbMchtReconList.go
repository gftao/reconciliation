package models

type Tbl_mcht_recon_list struct {
	MCHT_CD    string `gorm:"clumn:MCHT_CD"`
	USER       string `gorm:"clumn:user"`
	PASSWD     string `gorm:"clumn:passwd"`
	HOST       string `gorm:"clumn:host"`
	PORT       string `gorm:"clumn:port"`
	REMOTE_DIR string `gorm:"clumn:remote_dir"`
	MCHT_TP    string `gorm:"clumn:mcht_tp"`
	Transp_ty  string `gorm:"clumn:transp_ty"`
	EXT1       string `gorm:"clumn:EXT1"`
	EXT2       string `gorm:"clumn:EXT2"`
	EXT3       string `gorm:"clumn:EXT3"`
	EXT4       string `gorm:"clumn:EXT4"`
	EXT5       string `gorm:"clumn:EXT5"`
}

func (t Tbl_mcht_recon_list) TableName() string {
	return "tbl_mcht_recon_list"
}
