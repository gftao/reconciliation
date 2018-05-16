package models

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
