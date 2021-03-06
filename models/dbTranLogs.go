package models

type Tran_logs struct {
	REC_UPD_TS        string
	REC_CRT_TS        string
	TRAN_CD           string
	PROD_CD           string
	BIZ_CD            string
	MCHT_CD           string
	MCHT_NM           string
	TERM_ID           string
	TERM_SEQ          string
	TERM_BATCH        string
	ORDER_ID          string
	ORDER_DESC        string
	TRAN_DT_TM        string
	RESP_CD           string
	RESP_MSG          string
	SYS_ORDER_ID      string
	ORIG_SYS_ORDER_ID string
	ORIG_TRANS_DT     string
	ORDER_TIMEOUT     string
	PRE_AUTH_ID       string
	CURR_CD           string
	PRI_ACCT_NO       string
	TRAN_AMT          string
	SETT_DT           string
	CHECK_FLG         string
	ACQ_INS_ID_CD     string
	ISS_INS_ID_CD     string
	SIGN_IMG          string
	CERT_TP           string
	CERT_ID           string
	CUSTOMER_NM       string
	PHONE_NO          string
	SMS_CODE          string
	POS_ENTRY_CD      string
	TERM_ENTRY_CAP    string
	INSTAL_NUM        string
	INSTAL_RATE       string
	MCHT_FEE_SUBSIDY  string
	AUTH_CODE         string
	QR_TYPE           string
	TIME_OUT          string
	BUYER_USER        string
	INS_ORDER_ID      string
	GOODS_ID          string
	GOODS_NM          string
	GOODS_NUM         string
	GOODS_PRICE       string
	IP_ADDR           string
	GPS_ADDR          string
	FWD_INS_ID_CD     string
	CANCEL_FLG        string
	TRANS_ST          string
	TRAN_NM           string
	TRANS_IN_ACCT_NO  string
	TRANS_DT          string
	CLD_ORDER_ID      string
	DES_ACQ_INS_ID    string
	DES_MCHNT_CD      string
	Ext_fld1          string
	Ext_fld2          string
	Ext_fld3          string
	Ext_fld4          string
	Ext_fld5          string
	Ext_fld6          string
	Ext_fld7          string
	Ext_fld8          string
	Ext_fld9          string
	Ext_fld10         string
	Ext_fld11         string
	Ext_fld12         string
	Ext_fld13         string
	Ext_fld14         string
	Ext_fld15         string
	Ext_fld16         string
	Ext_fld17         string
	Ext_fld18         string
	Ext_fld19         string
	Ext_fld20         string
	CUST_ORDER_ID     string //
	ORDER_STAT        string
	ORDER_INFO        string
	BE_ORDER_ID       string
	ORDER_NOTIY_INFO  string
	BE_BIZ_CD         string
}

func (t Tran_logs) TableName() string {
	return "prodpmpcld.tran_logs" //prodPmpCld60
}