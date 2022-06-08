package models

import (
	"reflect"
	"strings"
)

type GYESFileStrt struct {
	FileHead     string
	GYESFileInfo []GYESFileHeadInfo
}

type GYESFileHeadInfo struct {
	Zjjgcmno           string //合同监管号
	PaymentSequenceNo  string //交易流水号
	PaymentAmount      string //交易金额
	PaymentTime        string //交易时间
	PaymentAccount     string //付款账号
	PaymentAccountName string //付款人姓名
	CommodityNo        string //商户号
	TerminalNo         string //终端号
	Nothing            string
}

func (fs *GYESFileStrt) Init() {
	fs.FileHead = ""
}

func (fs GYESFileHeadInfo) HToString() string {
	t := reflect.TypeOf(fs)
	v := reflect.ValueOf(fs)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		strs = append(strs, v.Field(i).String())
	}
	s := strings.Join(strs, "^?")
	return s
}
