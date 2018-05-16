package models

import (
	"reflect"
	"strings"
	"encoding/xml"
	"golib/gerror"
	"golib/security"
	"golib/defs"
	"crypto/rsa"
	"fmt"
	"github.com/jinzhu/gorm"
	"golib/modules/logr"
)

type FileStrt1 struct {
	FileHead     string
	FileHeadInfo FileHeadInfo1

	FileBody  string
	FileBodys []Body1
}

type FileHeadInfo1 struct {
	TrnSucCount string //总笔数
	TrnSucAm    string //成功总金额
}

type Body1 struct {
	USER_NUMBER    string //用户编号
	CHARGE_YEAR_ID string //年度id
	TRACE_NO       string //流水号
	DATE           string //缴费日期
	AMT            string //缴费金额
	TYPE           string //收费方式
	POS_ID         string //终端号
}

func (fs FileStrt1) HToString() string {
	t := reflect.TypeOf(fs.FileHeadInfo)
	v := reflect.ValueOf(fs.FileHeadInfo)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		strs = append(strs, v.Field(i).String())
	}
	s := strings.Join(strs, "|")
	return s
}

func (by Body1) BToString() string {
	t := reflect.TypeOf(by)
	v := reflect.ValueOf(by)
	strs := []string{}
	for i := 0; i < t.NumField(); i++ {
		strs = append(strs, v.Field(i).String())

	}
	s := strings.Join(strs, "|")

	return s
}

type F587 struct {
	XMLName   xml.Name  `xml:"UTILITY_PAYMENT" json:"-"`
	SIGN_TYPE string    `xml:"SIGN_TYPE" json:"sign_type"`
	SIGN      string    `xml:"SIGN" json:"sign"`
	DATA_AREA F587_AREA `json:"data_area"`
}

type F587_AREA struct {
	XMLName   xml.Name `xml:"DATA_AREA" json:"-"`
	TRAN_CODE string   `xml:"TRAN_CODE" json:"tran_code"`
	POS_ID    string   `xml:"POS_ID" json:"POS_ID"`
	DATE      string   `xml:"DATE"`
	FILE_NAME string   `xml:"FILE_NAME"`
	BANKNO    string   `xml:"BANKNO"`
}

func (f *F587) BuilMsg(m interface{}) ([]byte, error) {
	b, err := xml.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (f *F587) SingUp(dbtype, dbstr string, content []byte) (string, error) {

	//取出密钥信息
	keyHandle := KeyHandleInfo{}
	gerr := GetProdKeyHandle(&keyHandle, dbtype, dbstr, "1130", "0000037", "5201")
	if gerr != nil {
		return "", gerr
	}

	rsaSign, err := security.RsaSignSha1Base64(keyHandle.ServerPriKey, content)
	if err != nil {
		return "", gerror.New(30070, defs.TRN_SYS_ERROR, err, "生成签名失败")
	}
	return rsaSign, nil
}

func BuilSoap(FuncName, xl string) ([]byte, error) {
	if FuncName == "" {
		return nil, fmt.Errorf("funcName is nil")
	}
	x := xml.Header + xl
	h := `<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope" xmlns:heat="http://heat.sheshu.cn/"><soap:Header/><soap:Body><heat:` +
		FuncName + `><xmlParam><![CDATA[` + x + `]]></xmlParam></heat:` +
		FuncName + `> </soap:Body></soap:Envelope>`

	return []byte(h), nil
}

func LoadResp(RspBuf []byte) (string, error) {
	Rsp := string(RspBuf)
	lbegin := strings.Index(string(Rsp), `<result>`)
	lend := strings.Index(string(Rsp), `</result>`)
	if lbegin == -1 || lend == -1 {
		return "", fmt.Errorf("应答内容格式错误<result>未找到")
	}
	r := string(Rsp)[lbegin+len(`<result>`):lend]

	r = strings.Replace(r, `&lt;`, "<", -1)
	r = strings.Replace(r, `&gt;`, ">", -1)
	lbegin = strings.Index(string(r), `<DATA_AREA>`)
	lend = strings.Index(string(r), `</DATA_AREA>`)
	if lbegin == -1 || lend == -1 {
		return "", fmt.Errorf("应答内容格式错误<DATA_AREA>未找到")
	}
	signed := r[lbegin:lend+len(`</DATA_AREA>`)]
	lbegin = strings.Index(string(r), `<SIGN>`)
	lend = strings.Index(string(r), `</SIGN>`)
	if lbegin == -1 || lend == -1 {
		return "", fmt.Errorf("应答内容格式错误<SIGN>未找到")
	}
	Signature := r[lbegin+len(`<SIGN>`):lend]

	//ok, gerr := RsaVerify(signed, Signature)
	//if !ok || gerr != nil {
	//	return "", gerr
	//}
	logr.Debugf("验证签名：%s,%s", signed, Signature)

	return r, nil
}

type KeyHandleInfo struct {
	Term_key     string
	TermPubKey   *rsa.PublicKey
	TermPriKey   *rsa.PrivateKey
	ServerPubKey *rsa.PublicKey
	ServerPriKey *rsa.PrivateKey
}

func GetProdKeyHandle(k *KeyHandleInfo, dbtype, dbstr, prod_cd, biz_cd, trans_cd string) error {

	//取出终端密钥信息
	prodInfo := DbOrderInfo{}
	dbc, err := gorm.Open(dbtype, dbstr)
	if err != nil {
		logr.Info("open db err:", err)
		return err
	}
	defer dbc.Close()
	dbc = dbc.Set("gorm:table_options", "ENGINE=InnoDB")
	dbc.DB().Ping()
	err = dbc.Where("prod_cd = ? and biz_cd = ? and trans_cd =?",
		prod_cd, biz_cd, trans_cd).Find(&prodInfo).Error
	if err != nil {
		return gerror.New(10040, defs.TRN_SYS_ERROR, err, "取产品密钥信息失败")
	}

	prikey, err := security.GetRsaPrivateKeyByString(prodInfo.ServerPriKey)
	if err != nil {
		return gerror.New(10050, defs.TRN_SYS_ERROR, err, "取产品私钥信息失败")
	}
	k.ServerPriKey = prikey
	pubkey, err := security.GetRsaPublicKeyByString(prodInfo.ServerPubKey)
	if err != nil {
		return gerror.New(10060, defs.TRN_SYS_ERROR, err, "取产品公钥信息失败")
	}
	k.TermPubKey = pubkey
	gconf.HttpUrl = prodInfo.Http_url
	return nil
}
