package models

type Tbl_clear_txn struct {
	COMPANY_CD          string `gorm:"column:COMPANY_CD"`
	INS_ID_CD           string `gorm:"column:INS_ID_CD"`       //机构号
	ACQ_INS_ID_CD       string `gorm:"column:ACQ_INS_ID_CD"`
	FWD_INS_ID_CD       string `gorm:"column:FWD_INS_ID_CD"`
	MCHT_CD             string `gorm:"column:MCHT_CD"`         //商户号
	MCHT_NAME           string `gorm:"column:MCHT_NAME"`
	MCHT_SHORT_NAME     string `gorm:"column:MCHT_SHORT_NAME"`
	MCC_CD              string `gorm:"column:MCC_CD"`
	MCC_CD_42           string `gorm:"column:MCC_CD_42"`
	MCC_DESC            string `gorm:"column:MCC_DESC"`
	TRANS_DATE_TIME     string `gorm:"column:TRANS_DATE_TIME"` //交易日期
	STLM_DATE           string `gorm:"column:STLM_DATE"`       //清算日期
	TRANS_KIND          string `gorm:"column:TRANS_KIND"`      //交易类型
	TXN_DESC            string `gorm:"column:TXN_DESC"`
	TRANS_STATE         string `gorm:"column:TRANS_STATE"`
	STLM_FLG            string `gorm:"column:STLM_FLG"`
	TRANS_AMT           string `gorm:"column:TRANS_AMT"`       //交易本金 清算金额
	CREDITCARDLIMIT     string `gorm:"column:CREDITCARDLIMIT"`
	CUP_SSN             string `gorm:"column:CUP_SSN"`
	AUTHR_ID_RESP       string `gorm:"column:AUTHR_ID_RESP"`
	PAN                 string `gorm:"column:PAN"`             //交易卡号
	CARD_KIND_DIS       string `gorm:"column:CARD_KIND_DIS"`   //卡类型
	BANK_CODE           string `gorm:"column:BANK_CODE"`
	BANK_NAME           string `gorm:"column:BANK_NAME"`
	BRANCH_CD           string `gorm:"column:BRANCH_CD"`
	BRANCH_NM           string `gorm:"column:BRANCH_NM"`
	TERM_ID             string `gorm:"column:TERM_ID"`         //终端编号
	ORG_TRANS_DATE_TIME string `gorm:"column:ORG_TRANS_DATE_TIME"`
	ORG_CUP_SSN         string `gorm:"column:ORG_CUP_SSN"`
	POS_ENTRY_MODE      string `gorm:"column:POS_ENTRY_MODE"`
	RSP_CODE            string `gorm:"column:RSP_CODE"`
	TRUE_FEE_MOD        string `gorm:"column:TRUE_FEE_MOD"`    //交易手续费 错的
	TRUE_FEE_BI         string `gorm:"column:TRUE_FEE_BI"`
	TRUE_FEE_FD         string `gorm:"column:TRUE_FEE_FD"`
	TRUE_FEE_FFD        string `gorm:"column:TRUE_FEE_FFD"`
	VAR_1               string `gorm:"column:VAR_1"`
	VAR_2               string `gorm:"column:VAR_2"`
	VAR_3               string `gorm:"column:VAR_3"`
	VAR_4               string `gorm:"column:VAR_4"`
	VIR_FEE_MOD         string `gorm:"column:VIR_FEE_MOD"`
	VIR_FEE_BI          string `gorm:"column:VIR_FEE_BI"`
	VIR_FEE_BD          string `gorm:"column:VIR_FEE_BD"`
	VIR_FEE_FD          string `gorm:"column:VIR_FEE_FD"`
	MCHT_FEE            string `gorm:"column:MCHT_FEE"`		//交易手续费
	VAR_5               string `gorm:"column:VAR_5"`
	MCHT_VIR_FEE        string `gorm:"column:MCHT_VIR_FEE"`
	STAND_BANK_FEE      string `gorm:"column:STAND_BANK_FEE"`
	BANK_FEE            string `gorm:"column:BANK_FEE"`
	HZJG_FEE            string `gorm:"column:HZJG_FEE"` //应付费用
	JGSY                string `gorm:"column:JGSY"`
	AIP_FEE             string `gorm:"column:AIP_FEE"`
	MCHT_SET_AMT        string `gorm:"column:MCHT_SET_AMT"`    //交易结算资金
	HZJGYFPPFWF         string `gorm:"column:HZJGYFPPFWF"`
	JGYFPPFWF           string `gorm:"column:JGYFPPFWF"`
	AIPYFPPFWF          string `gorm:"column:AIPYFPPFWF"`
	ERR_FEE_IN          string `gorm:"column:ERR_FEE_IN"`
	ERR_FEE_OUT         string `gorm:"column:ERR_FEE_OUT"`
	ERR_CODE            string `gorm:"column:ERR_CODE"`
	JT_MCHT_CD          string `gorm:"column:JT_MCHT_CD"`
	EXPAND_ORG_CD       string `gorm:"column:EXPAND_ORG_CD"`
	SPE_SERV_INST       string `gorm:"column:SPE_SERV_INST"`
	PROP_INS            string `gorm:"column:PROP_INS"`
	EXPAND_ORG_FEE      string `gorm:"column:EXPAND_ORG_FEE"`
	SPE_SERV_FEE        string `gorm:"column:SPE_SERV_FEE"`
	PROP_INS_FEE        string `gorm:"column:PROP_INS_FEE"`
	EXPAND_ORG_PP       string `gorm:"column:EXPAND_ORG_PP"`
	SPE_SERV_PP         string `gorm:"column:SPE_SERV_PP"`
	PROP_INS_PP         string `gorm:"column:PROP_INS_PP"`
	EXPAND_FEE_IN       string `gorm:"column:EXPAND_FEE_IN"`
	EXPAND_FEE_OUT      string `gorm:"column:EXPAND_FEE_OUT"`
	CUP_IFINSIDE_SIGN   string `gorm:"column:CUP_IFINSIDE_SIGN"`
	SP_CHARG_TYPE       string `gorm:"column:SP_CHARG_TYPE"`
	SP_CHARG_LEV        string `gorm:"column:SP_CHARG_LEV"`
	TERM_SSN            string `gorm:"column:TERM_SSN"`
	SN_SSN              string `gorm:"column:SN_SSN"`
	UP_CHL_ID           string `gorm:"column:UP_CHL_ID"`
	CONV_MCHT_CD        string `gorm:"column:CONV_MCHT_CD"`
	CONV_TERM_ID        string `gorm:"column:CONV_TERM_ID"`
	CHL_TRUE_FEE        string `gorm:"column:CHL_TRUE_FEE"`
	CHL_STD_FEE         string `gorm:"column:CHL_STD_FEE"`
	CHL_FEE_PRE_FLG     string `gorm:"column:CHL_FEE_PRE_FLG"`
	SYS_SER             string `gorm:"column:SYS_SER"`         //系统流水号
	VAR_6               string `gorm:"column:VAR_6"`
	QUDAO_FEE           string `gorm:"column:QUDAO_FEE"`
	QUDAO_FEE_MIN       string `gorm:"column:QUDAO_FEE_MIN"`
	QUDAO_FEE_MIX       string `gorm:"column:QUDAO_FEE_MIX"`
	QUDAO_FEE_FD        string `gorm:"column:QUDAO_FEE_FD"`
	INS_FEE             string `gorm:"column:INS_FEE"`
	INS_MY_FEE          string `gorm:"column:INS_MY_FEE"`
	INS_COST_FEE        string `gorm:"column:INS_COST_FEE"`
	INS_MY_FEE_AMT      string `gorm:"column:INS_MY_FEE_AMT"`
	INS_SPLIT_FEE       string `gorm:"column:INS_SPLIT_FEE"`
	INS_RES_FEE         string `gorm:"column:INS_RES_FEE"`
	PINP_FEE            string `gorm:"column:PINP_FEE"`
	PINP_FEE_INF        string `gorm:"column:PINP_FEE_INF"`
	PINP_FEE_TOP        string `gorm:"column:PINP_FEE_TOP"`
	PINP_STAT           string `gorm:"column:PINP_STAT"`
	T0_STAT             string `gorm:"column:T0_STAT"`
	KEY_RSP             string `gorm:"column:KEY_RSP"`         //交易流水号
	REMARK              string `gorm:"column:REMARK"`
	REMARK1             string `gorm:"column:REMARK1"`
	REMARK2             string `gorm:"column:REMARK2"`
	REMARK3             string `gorm:"column:REMARK3"`
	REMARK4             string `gorm:"column:REMARK4"`
	REMARK5             string `gorm:"column:REMARK5"`
}

func (t *Tbl_clear_txn) TableName() string {
	return "tbl_clear_txn"
}