package models

import (
"reflect"
"strings"
)

type GYFileStrt  struct {
	FileHead     string
	GYFileInfo []GYFileHeadInfo
}

type GYFileHeadInfo struct {
	STLM_DATE   string //清算日期
	TRAN_CD   string //交易代码
	JC_BIZ_CD string //缴存编号
	MCHT_CD    string //商户号
	TERM_ID		string//终端号
	PAY_AMT		string//交易金额
	SYS_ORDER_ID		string//扣款流水
	ZZJ_ORDER_ID	string//住建局流水
	ACC_AMOUNT	string//银行卡号
	TRAN_TM	string//交易时间
	JC_KIND	string//缴存方式
	JC_TYPE 	string//缴存类型
	BANK_CD	string//银行编号
	SUP_AMOUNT	string//监管账户
}

func (fs *GYFileStrt) Init() {
	fs.FileHead = "清算日期|交易代码|缴存编号|商户号|终端号|交易金额（分）|扣款流水|住建局流水号|银行卡号|交易时间|缴交方式|缴交类型|银行编号|监管账户|"
}

func (fs GYFileHeadInfo) HToString() string {
	t := reflect.TypeOf(fs)
	v := reflect.ValueOf(fs)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		strs = append(strs, v.Field(i).String())
	}
	s := strings.Join(strs, "|")
	return s
}


