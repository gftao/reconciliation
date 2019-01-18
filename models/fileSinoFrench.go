package models

import (
	"reflect"
	"strings"
)

type FileSinoFrench struct {
	FileHead     string
	FileHeadInfo FileHeadSinoFrench

	FileBody  string
	FileBodys []SFBody
}
type FileHeadSinoFrench struct {
	INS_ID_CD   string //机构代码
	Stlm_date   string //交易日期
	TimeB       string //开始时间
	TimeE       string //结束时间
	TrnSucCount string //笔数
	TrnReconT   string //汇总金额
}

func (fs *FileSinoFrench) Init() {
	fs.FileHead = "机构代码|交易日期|开始时间|结束时间|笔数|汇总金额"
	fs.FileBody = "流水号|注册号+卡号|水量|金额|终端编号|网点编号|分公司营业所编号||"
}
func (fs FileSinoFrench) HToString() string {
	t := reflect.TypeOf(fs.FileHeadInfo)
	v := reflect.ValueOf(fs.FileHeadInfo)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		strs = append(strs, v.Field(i).String())
	}
	s := strings.Join(strs, "|")
	return s + "||"
}

type SFBody struct {
	CUST_ORDER_ID   string //流水号
	Regist_Meter_No string // |注册号+卡号
	Amount          string // |水量
	TRANS_AMT       string // |金额
	TERM_ID         string // |终端编号
	Net_No          string // |网点编号
	Mcht_NO         string // |分公司营业所编号
}

func (by SFBody) BToString() string {
	t := reflect.TypeOf(by)
	v := reflect.ValueOf(by)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		strs = append(strs, v.Field(i).String())
	}
	s := strings.Join(strs, "|")

	return s + "||"
}
