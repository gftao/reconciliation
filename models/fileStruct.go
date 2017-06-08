package models


type FileStrt struct {
	FileHead     string
	FileHeadInfo FileHeadInfo
	//FileHeadInfo
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
			     //机构上送订单号
	SYS_ID        string //系统流水号  INDUSTRY_ADDN_INF(扫码)RETRI_REF_NO(收单)
	INS_IN        string //机构基准收入 0
	INS_REAL_IN   string //机构实际收入 0
	INS_OUT       string //机构营销返佣 0
	PROXY_CD      string //代理编码 0
	MEMBER_ID     string //会员号 0
}

func (fs *FileStrt) Init() {
	fs.FileHead = "机构代码,清算日期,交易总笔数,清算金额,清算手续费,结算总金额"
	fs.FileBody = "商户号,交易日期,交易时间,清算日期,终端编号,交易类型,交易流水号,交易卡号,卡类型,交易本金,交易手续费,交易结算资金,应收差错费用,应付差错费用,系统流水号,机构基准收入,机构实际收入,机构营销返佣,代理编码,会员号"

}

func (fs FileStrt) HToString() string {
	return fs.FileHeadInfo.INS_ID_CD + "," + fs.FileHeadInfo.Stlm_date + "," + fs.FileHeadInfo.TrnSucCount + "," + fs.FileHeadInfo.TrnSucAm + "," + fs.FileHeadInfo.TrnFeeT + "," + fs.FileHeadInfo.TrnReconT
}

func (by Body) BToString() string {
	return by.MCHT_CD + ","		 +
	by.TRANS_DATE + "," +
	by.TRANS_TIME + "," +
	by.STLM_DATE + "," +
	by.TERM_ID + "," +
	by.TRANS_KIND + "," +
	by.KEY_RSP + "," +
	by.PAN + "," +
	by.CARD_KIND_DIS + "," +
	by.TRANS_AMT + "," +
	by.TRUE_FEE_MOD + "," +
	by.MCHT_SET_AMT + "," +
	by.ERR_FEE_IN + "," +
	by.ERR_FEE_OUT + "," +
	by.SYS_ID + "," +
	by.INS_IN + "," +
	by.INS_REAL_IN + "," +
	by.INS_OUT + "," +
	by.PROXY_CD + "," +
	by.MEMBER_ID
}