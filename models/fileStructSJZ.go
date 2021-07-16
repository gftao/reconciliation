package models

import (
	"reflect"
	"strings"
)

type SJZFileStrt  struct {
	FileHead     string
	SJZFileHeadInfo SJZFileHeadInfo

	FileBody  string
	FileBodys []SJZBody
}

type SJZFileHeadInfo struct {
	Be_biz_cd	string	//中间业务编码
	Pay_flag		string	//收付标志
	Total_amt	string	//交易总金额
	Total_tran	string	//交易总笔数
}

type SJZBody struct {
	TranDate    string //核心交易日期
	TranTime	string//前台时间
	SysId        string //核心流水
	FrontDate	string//前台日期
	FrontId		string//前台流水
	Currency		string//币种
	CurrFlag		string//钞汇标志
	OutType		string//转出类型
	OutAccount	string//转出账号
	InType		string//转入类型
	InAccount		string//转入账号
	TranAmt		string//交易金额
	ThirdId		string//第三方流水号
	ThirdNo		string//第三方号码
	PayPeriod		string//缴费账期
	TranTunnel			string//交易渠道
	UserAddress		string//用户地址
	Ext1		string//备注1
	Ext2		string//备注2
	Ext3		string//备注3
	Ext4		string//备注4
	Ext5		string//备注5

}

func (fs *SJZFileStrt) Init() {
	fs.FileHead = "中间业务编码|收付标志|交易总金额|总笔数"
	fs.FileBody = "交易日期|交易时间|核心流水|前台日期|前台流水|币种|钞汇标志|转出类型|转出帐号|转入类型|转入帐号|交易金额|" +
		"第三方流水号|第三方号码|缴费账期|交易渠道|用户地址|备注1|备注2|备注3|备注4|备注5"
}

func (fs SJZFileStrt) HToString() string {
	t := reflect.TypeOf(fs.SJZFileHeadInfo)
	v := reflect.ValueOf(fs.SJZFileHeadInfo)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		strs = append(strs, v.Field(i).String())
	}
	s := strings.Join(strs, "|")
	return s
}

func (by SJZBody) BToString() string {
	t := reflect.TypeOf(by)
	v := reflect.ValueOf(by)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		strs = append(strs, v.Field(i).String())
	}
	s := strings.Join(strs, "|")
	return s
}


