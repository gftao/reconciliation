package models

import (
	"reflect"
	"strings"
)

type FuzhouStrt struct {
	FileHead     string
	FileHeadInfo FuzhouHeadInfo

	FileBody  string
	FileBodys []FuzhouBody
}

type FuzhouHeadInfo struct {
	Area_CD     string //服务商地区代码
	Stlm_date   string //清算日期
	TrnSucCount string //成功总笔数
	//TrnSucAm    string //成功总金额
	//TrnFeeT     string //交易总手续费
	TrnReconT string //交易总结算额
}

type FuzhouBody struct {
	MCHT_CD       string //商户号
	GF_BIZ_CD     string //购房业务编号
	TRANS_KIND    string //交易类型
	TRANS_DATE    string //业务发生时间  交易日期 Tbl_tfr_his_trn_log->TRANS_DT
	STLM_DATE     string //清算日期
	MCHT_SET_AMT  string //交易结算资金
	SYS_ID        string //系统流水号  INDUSTRY_ADDN_INF(扫码)RETRI_REF_NO(收单)
	CUST_ORDER_ID string //第三方订单号//机构上送订单号
	EXT_FLD1      string //备注
	EXT_FLD2      string //备注
	EXT_FLD3      string //备注
}

func (fs *FuzhouStrt) Init() {
	fs.FileHead = "服务商地区代码,清算日期,交易总笔数,清算总金额"
	fs.FileBody = "商户号,购房业务编号,交易方式,业务发生时间,清算日期,交易清算金额,交易流水号,第三方订单号,备注1,备注2,备注3"
}

func (fs FuzhouStrt) HToString() string {
	t := reflect.TypeOf(fs.FileHeadInfo)
	v := reflect.ValueOf(fs.FileHeadInfo)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		strs = append(strs, v.Field(i).String())
	}
	s := strings.Join(strs, ",")
	return s
}

func (by FuzhouBody) BToString() string {
	t := reflect.TypeOf(by)
	v := reflect.ValueOf(by)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		strs = append(strs, v.Field(i).String())
	}

	s := strings.Join(strs, ",")

	return s
}
