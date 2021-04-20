package models

import (
	"reflect"
	"strings"
)

type KMFileStrt struct {
	FileHead       string
	KMFileHeadInfo KMFileHeadInfo

	FileBody  string
	FileBodys []KMBody
}

type KMFileHeadInfo struct {
	Area_CD     string //服务商地区代码
	Stlm_date   string //清算日期
	TrnSucCount string //成功总笔数
	TrnReconT   string //交易总结算额
}

type KMBody struct {
	MCHT_CD       string //商户号
	GF_BIZ_CD     string //购房业务编号
	TRANS_KIND    string //交易类型
	TRANS_DATE    string //业务发生时间  交易日期 Tbl_tfr_his_trn_log->TRANS_DT
	STLM_DATE     string //清算日期
	MCHT_SET_AMT  string //交易结算资金
	SYS_ID        string //系统流水号  INDUSTRY_ADDN_INF(扫码)RETRI_REF_NO(收单)
	CUST_ORDER_ID string //第三方订单号//机构上送订单号
	EXT_FLD1      string //备注1
	EXT_FLD2      string //备注2
	EXT_FLD3      string //备注3
	EXT_FLD4      string //备注4----对应昆明交易金额
	EXT_FLD5      string //备注5
	EXT_FLD6      string //备注6----对应昆明手续费
	EXT_FLD7      string //备注7
	EXT_FLD8      string //备注8
	EXT_FLD9      string //备注9
}

func (fs *KMFileStrt) Init() {
	fs.FileHead = "服务商地区代码,清算日期,交易总笔数,清算总金额"
	fs.FileBody = "商户号,购房业务编号,交易方式,业务发生时间,清算日期,交易清算金额,交易流水号,第三方订单号,备注1,备注2,备注3," +
		"备注4,备注5,备注6,备注7,备注8,备注9"
}

func (fs KMFileStrt) HToString() string {
	t := reflect.TypeOf(fs.KMFileHeadInfo)
	v := reflect.ValueOf(fs.KMFileHeadInfo)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		strs = append(strs, v.Field(i).String())
	}
	s := strings.Join(strs, ",")
	return s
}

func (by KMBody) BToString() string {
	t := reflect.TypeOf(by)
	v := reflect.ValueOf(by)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		//fmt.Printf("%s---%v \n", t.Field(i).Name,v.Field(i).Interface())
		//str = append(str, v.Field(i).Interface().(string))
		strs = append(strs, v.Field(i).String())

	}

	s := strings.Join(strs, ",")

	return s
}
