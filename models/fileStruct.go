package models

import (
	"reflect"
	"strings"
)

type FileStrt struct {
	FileHead     string
	FileHeadInfo FileHeadInfo

	FileBody     string
	FileBodys    []Body
}

type FileHeadInfo struct {
	INS_ID_CD   string //机构代码
	Stlm_date   string //清算日期
	TrnSucCount string //成功总笔数
	TrnSucAm    string //成功总金额
	TrnFeeT     string //交易总手续费
	TrnReconT   string //交易总结算额
}

type Body struct {
	MCHT_CD       string //商户号
	TRANS_DATE    string //交易日期 Tbl_tfr_his_trn_log->TRANS_DT
	TRANS_TIME    string //交易时间 Tbl_tfr_his_trn_log->TRANS_MT
	STLM_DATE     string //清算日期
	TERM_ID       string //终端编号
	TRANS_KIND    string //交易类型
	KEY_RSP       string //交易流水号
	PAN           string //交易卡号
	CARD_KIND_DIS string //卡类型
	TRANS_AMT     string //交易本金
	TRUE_FEE_MOD  string //交易手续费
	MCHT_SET_AMT  string //交易结算资金
	ERR_FEE_IN    string //应收差错费用 0
	ERR_FEE_OUT   string //应付差错费用 0
	SYS_ID        string //系统流水号  INDUSTRY_ADDN_INF(扫码)RETRI_REF_NO(收单)
	INS_IN        string //机构基准收入 0
	INS_REAL_IN   string //机构实际收入 0
	INS_OUT       string //机构营销返佣 0
	PROXY_CD      string //代理编码 0
	MEMBER_ID     string //会员号 0

	DUES          string //应付费用
	PROD_CD       string //产品码
	TRAND_CD      string //交易码
	BIZ_CD        string //业务码
	CUST_ORDER_ID string //第三方订单号//机构上送订单号
}

func (fs *FileStrt) Init() {
	fs.FileHead = "机构代码,清算日期,交易总笔数,清算金额,清算手续费,结算总金额"
	fs.FileBody = "商户号,交易日期,交易时间,清算日期,终端编号,交易类型,交易流水号," +
		"交易卡号,卡类型,交易本金,交易手续费,交易结算资金,应收差错费用,应付差错费用," +
		"系统流水号,机构基准收入,机构实际收入,机构营销返佣,代理编码,会员号," +
		"应付费用,产品码,交易码,业务码,第三方订单号"

}

func (fs FileStrt) HToString() string {
	t := reflect.TypeOf(fs.FileHeadInfo)
	v := reflect.ValueOf(fs.FileHeadInfo)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		strs = append(strs, v.Field(i).String())
	}
	s := strings.Join(strs, ",")
	return s
}

func (by Body) BToString() string {
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