package models

import (
	"reflect"
	"strings"
)

type HZFileStrt  struct {
	//INS_ID_CD	string
	FileHead     string
	HZFileInfo []HZFileHeadInfo
}

type HZFileHeadInfo struct {
	KEY_RSP   string //交易流水号
	MasterAcccount   string //主账户户号
	SubAccount string //子账户户号
	TranDate    string //交易发生日期
	TranTime	string//交易发生时间
	Currency		string//币种
	TranAmt		string//交易金额
	TranAccountName	string//对方账户户名
	TranBankName	string//对方银行名称
	EXT_FLD1	string//交易备注
	EXT_FLD2	string//交易备注2
}

func (fs *HZFileStrt) Init() {
	fs.FileHead = "交易流水号|主账户户号|子账户户号|交易发生日期|交易发生时间|币种|交易金额|对方账户户名|对方银行名称|交易备注|交易备注2"
}

func (fs HZFileHeadInfo) HToString() string {
	t := reflect.TypeOf(fs)
	v := reflect.ValueOf(fs)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		strs = append(strs, v.Field(i).String())
	}
	s := strings.Join(strs, "|")
	return s
}

