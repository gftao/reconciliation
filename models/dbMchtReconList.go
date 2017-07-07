package models

type Tbl_mcht_recon_list struct {
	MCHT_CD    string
	USER       string
	PASSWD     string
	HOST       string
	PORT       string
	REMOTE_DIR string
}

func (t Tbl_mcht_recon_list) TableName() string {
	return "tbl_mcht_recon_list"
}