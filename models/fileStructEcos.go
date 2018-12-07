package models

import (
	"reflect"
	"strings"
)

var TranCdConvert = map[string]string{
	"1011": "1101",
	"2011": "3101",
	"3011": "5151",

	"1131": "1001",
	"2131": "3001",
	"3131": "5141",
	"1031": "    ",
	"2031": "    ",
}

var CARDConvert = map[string]string{
	"0000007": "11",
	"0000008": "12",
	"0000000": "  ",
}

type FileStrtEchos struct {
	Stlm_date     string //清算日期
	MCHT_CD       string //收单商户号
	TERM_ID       string //收单终端号
	TRANS_TIME    string //交易传输时间
	PAN           string //卡号
	KEY_RSP       string //收单系统流水号
	TRAND_CD      string //收单内部交易码
	Resp_cd       string //交易应答码
	TRANS_AMT     string //交易金额
	CARD_KIND_DIS string //交易卡种
	CUST_ORDER_ID string //生态圈卡券交易唯一标识
	Stl_flag      string //对账标志
}

func (f FileStrtEchos) ToString() string {
	t := reflect.TypeOf(f)
	v := reflect.ValueOf(f)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		//if t.Field(i).Name == "PAN" {
		//	vi := v.Field(i).String()
		//	if vi == "0" {
		//		vi = ""
		//	}
		//	//strs = append(strs, fmt.Sprintf("%-#16s", vi))
		//	//strs = append(strs, v.Field(i).String())
		//} else {
		//	strs = append(strs, v.Field(i).String())
		//}
		strs = append(strs, v.Field(i).String())
	}
	s := strings.Join(strs, "||")
	return s
}
