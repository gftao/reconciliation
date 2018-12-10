package reconc

import (
	"golib/gerror"
	"golib/modules/gormdb"
	"github.com/jinzhu/gorm"
	"htdRec/models"
	"golib/modules/config"
	"bufio"
	"strconv"
	"htdRec/myfstp"
	"bytes"
	"strings"

	"golib/modules/logr"
	"database/sql"
	"htdRec/myftp"
	"golib/security"
	"time"
	"os"
)

type CrtFile struct {
	FileName             string                     //对账文件名
	FilePath             string                     //路径
	SysDate              string                     //当前时间
	MCHT_CD              string                     //当前机构号
	STLM_DATE            string                     //清算日期
	Ins_id_cd            []string                   //全部机构号
	FileStrt             models.FileStrt            //文件结构
	Tbl_Clear_Data       []models.Tbl_clear_txn     //当前机构号数据
	Tbl_Tfr_his_log_Data models.Tbl_tfr_his_trn_log //当前机构号对应交易日志
	dbtype               string
	dbstr                string
	MCHT_TP              map[string]string
	sendto               bool
}

func (cf *CrtFile) Init(chainName string, mc string) gerror.IError {
	cf.FileName = "C_"

	cf.FilePath = config.StringDefault("filePath", "")
	cf.sendto = config.BoolDefault("sendto", true)
	logr.Info("是否发送ftp：", cf.sendto)

	cf.SysDate = time.Now().Format("20060102") //今天
	cf.FileStrt = models.FileStrt{}
	cf.STLM_DATE = chainName //清算日期
	cf.MCHT_TP = make(map[string]string, 1)
	cf.FileStrt.Init()
	//查表获取需要生产队长文件机构的机构号,新增加的表
	cf.MCHT_CD = mc
	err := cf.InitMCHTCd(mc)
	if err != nil {
		return err
	}
	cf.indb()
	return nil
}

func (cf *CrtFile) indb() {
	//"prodPmpCld:prodPmpCld@tcp(192.168.20.60:3306)/prodPmpCld?charset=utf8&parseTime=True&loc=Local"
	config.SetSection("db1")
	dbtype := config.StringDefault("db.type", "mysql")
	dbhost := config.StringDefault("db.host", "127.0.0.1")
	dbport := config.StringDefault("db.port", "3306")
	dbname := config.StringDefault("db.dbname", "prodPmpCld")
	dbuser := config.StringDefault("db.user", "root")
	dbpasswd := config.StringDefault("db.passwd", "")
	connStr := dbuser + ":" + dbpasswd + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname +
		"?charset=utf8&parseTime=True&loc=Local"

	cf.dbtype = dbtype
	cf.dbstr = connStr
	logr.Info("[db1]：", cf.dbtype, cf.dbstr)
}

func (cf *CrtFile) Run() {

	gerr := cf.SaveToFile()
	if gerr != nil {
		logr.Errorf("%s", gerr)
	}
	return
}
func (cf *CrtFile) SaveToFile() gerror.IError {
	//读数据
	gerr := cf.ReadDate()
	if gerr != nil {
		logr.Info("未读到清算数据,不创建对账文件MCHT_CD:", cf.MCHT_CD)
		return gerr
	}
	//创建文件
	fn := cf.geneFile()
	logr.Info("对账文件路径：", fn)
	f, err := os.Create(fn)
	defer f.Close()
	if err != nil {
		return gerror.NewR(1001, err, "创建文件失败:%s", fn)
	}
	buf := []byte{}
	b := bytes.NewBuffer(buf)
	b.WriteString(cf.FileStrt.FileHead + "\r\n")
	b.WriteString(cf.FileStrt.HToString())
	b.WriteString("\r\n")

	b.WriteString(cf.FileStrt.FileBody) //文件体 标识
	b.WriteString("\r\n")
	for _, tc := range cf.FileStrt.FileBodys {
		Str := tc.BToString()
		b.WriteString(Str)
		b.WriteString("\r\n")
	}
	//logr.Infof("--读取的数据1---[%s]", b)
	rb := b.Bytes()
	logr.Infof("--读取的数据2---[%s]", string(rb))

	cf.postToSftp(cf.FileName, rb)

	w := bufio.NewWriter(f) //创建新的 Writer 对象
	var n int64
	for {
		c, err := b.WriteTo(w)
		if err != nil {
			break
		}
		if n == int64(b.Len()) {
			break
		}
		n = n + c
	}
	//fmt.Println("---读取的buf数据--:", buf)
	w.Flush()
	f.Sync()
	return nil
}

func (cf *CrtFile) postToSftp(fileName string, fileData []byte) {

	dbc := gormdb.GetInstance()

	rows, err := dbc.Raw("SELECT * FROM tbl_mcht_recon_list WHERE MCHT_CD = ?", cf.MCHT_CD).Rows()
	if err != nil {
		logr.Info("dbc find failed:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		tmr := &models.Tbl_mcht_recon_list{}
		dbc.ScanRows(rows, tmr)
		if tmr.USER == "" || tmr.PASSWD == "" || tmr.HOST == "" {
			logr.Infof("读配置错误:[%s][%s][%s]", tmr.USER, tmr.PASSWD, tmr.HOST)
			return
		}
		user := tmr.USER
		password := tmr.PASSWD
		host := tmr.HOST
		port := tmr.PORT
		rmtDir := tmr.REMOTE_DIR
		trans_ty := tmr.Transp_ty
		AESKEY := tmr.EXT2
		rmtDir = strings.Replace(rmtDir, "\\", "//", -1)
		logr.Infof("->:[%s][%s][%s][%s][%s][%s][%s]", trans_ty, user, password, host, port, fileName, rmtDir)
		logr.Info("Aes密钥：", AESKEY)
		if cf.sendto {
			if AESKEY != "" {
				logr.Info("Aes密钥：", AESKEY)
				cipherdata, err := security.AesEcbEncrypt(fileData, []byte(AESKEY))
				if err != nil {
					logr.Info(cf.FileName, err)
					return
				}
				AESinfo := security.EncodeBase64(cipherdata)
				fileData = AESinfo
				logr.Infof("--Aes加密后的数据---[%s]", fileData)
			}

			switch trans_ty {
			case "0":
				logr.Infof("SFTP:")
				myfstp.PosByteSftp(user, password, host, port, fileName, rmtDir, fileData)
			case "1":
				logr.Infof("FTP with TLS:")
				err = myftp.MyftpTSL(user, password, host, port, fileName, rmtDir, fileData)
				if err != nil {
					logr.Error(err)
				}
			case "2":
				logr.Infof("FTP without TLS:")
				err = myftp.Myftp(user, password, host, port, fileName, rmtDir, fileData)
				if err != nil {
					logr.Error(err)
				}
			default:

			}
		}
	}
}

func (cf *CrtFile) Finish() {
	return
}

func (cf *CrtFile) ReadDate() gerror.IError {
	dbc := gormdb.GetInstance()
	var rows *sql.Rows
	var err error
	logr.Info(cf.MCHT_CD, "->", cf.MCHT_TP[cf.MCHT_CD])
	if cf.MCHT_TP[cf.MCHT_CD] == "0" {
		rows, err = dbc.Raw("SELECT * FROM tbl_clear_txn WHERE MCHT_CD = ? and STLM_DATE = ?", cf.MCHT_CD, cf.STLM_DATE).Rows()
	} else if cf.MCHT_TP[cf.MCHT_CD] == "1" {
		rows, err = dbc.Raw("SELECT * FROM tbl_clear_txn WHERE JT_MCHT_CD = ? and STLM_DATE = ?", cf.MCHT_CD, cf.STLM_DATE).Rows()
	} else {
		return nil
	}
	defer rows.Close()
	if err == gorm.ErrRecordNotFound {
		return gerror.NewR(0000, err, "MCHT_CD[%s]记录不存在", cf.MCHT_CD)
	}
	if err != nil {
		logr.Info("tbl_clear_txn find fail: ", err)
		return gerror.NewR(1000, err, " ")
	}
	//cf.Tbl_Clear_Data = []models.Tbl_clear_txn{}

	for rows.Next() {
		tc := models.Tbl_clear_txn{}
		dbc.ScanRows(rows, &tc)
		cf.Tbl_Clear_Data = append(cf.Tbl_Clear_Data, tc)
	}
	//处理文件头
	if len(cf.Tbl_Clear_Data) > 0 {
		cf.FileStrt.FileHeadInfo.INS_ID_CD = cf.Tbl_Clear_Data[0].INS_ID_CD
	} else {
		return gerror.NewR(0000, err, "MCHT_CD[%s]记录不存在", cf.MCHT_CD)
	}

	gerr := cf.saveDatatoFStru()
	if gerr != nil {
		return gerr
	}
	return nil
}

func (cf *CrtFile) saveDatatoFStru() gerror.IError {
	//cf.FileStrt.FileBodys = make([]models.Body,0)
	cf.FileStrt.FileBodys = []models.Body{}
	dbc := gormdb.GetInstance()
	dbt, err := gorm.Open(cf.dbtype, cf.dbstr)
	if err != nil {
		logr.Info(err)
		return gerror.NewR(1000, err, "打开数据库链接失败")
	}
	defer dbt.Close()
	dbt = dbt.Set("gorm:table_options", "ENGINE=InnoDB")
	dbt.DB().Ping()
	record := 0           //交易总笔数
	trans_amt_T := 0.0    //清算金额
	true_fee_mod_T := 0.0 //清算手续费
	trnrecont_T := 0.0    //结算总金额
	for _, tc := range cf.Tbl_Clear_Data {
		b := models.Body{}
		tfr := models.Tbl_tfr_his_trn_log{}
		tran := models.Tran_logs{}
		err := dbc.Where("KEY_RSP = ?", tc.KEY_RSP).Find(&tfr).Error
		if err != nil {
			logr.Infof("dbc find failed, KEY_RSP = %s, err = %s", tc.KEY_RSP, err)
			continue
		}

		record ++
		a, _ := strconv.ParseFloat(tc.TRANS_AMT, 64)
		f, _ := strconv.ParseFloat(tc.MCHT_FEE, 64)
		m, _ := strconv.ParseFloat(tc.MCHT_SET_AMT, 64)

		trans_amt_T += a
		true_fee_mod_T += f
		trnrecont_T += m

		b.MCHT_CD = tc.MCHT_CD
		b.TRANS_DATE = tfr.TRANS_DT
		b.TRANS_TIME = tfr.TRANS_MT
		b.STLM_DATE = cf.STLM_DATE
		b.TERM_ID = tc.TERM_ID
		b.TRANS_KIND = tc.TXN_DESC
		b.KEY_RSP = tc.KEY_RSP
		b.PAN = tc.PAN
		b.CARD_KIND_DIS = tc.CARD_KIND_DIS
		b.TRANS_AMT = tc.TRANS_AMT
		b.TRUE_FEE_MOD = tc.MCHT_FEE
		b.MCHT_SET_AMT = tc.MCHT_SET_AMT
		b.ERR_FEE_IN = "0"
		b.ERR_FEE_OUT = "0"
		logr.Info("prod_cd:", tfr.PROD_CD)
		if tfr.PROD_CD == "1151" {
			b.SYS_ID = tfr.INDUSTRY_ADDN_INF
		} else {
			if tfr.RETRI_REF_NO[:1] == cf.STLM_DATE[3:4] {
				b.SYS_ID = cf.STLM_DATE[:3] + tfr.RETRI_REF_NO
			} else {
				b.SYS_ID = cf.STLM_DATE[:4] + tfr.RETRI_REF_NO
			}
		}
		b.INS_IN = "0"
		b.INS_REAL_IN = "0"
		b.INS_OUT = "0"
		b.PROXY_CD = "0"
		b.MEMBER_ID = "0"
		b.DUES = tc.HZJG_FEE
		b.PROD_CD = tfr.PROD_CD
		b.TRAND_CD = tfr.MA_TRANS_CD
		b.BIZ_CD = tfr.BIZ_CD
		//"2017050415054098157697"
		err = dbt.Where("sys_order_id = ?", b.SYS_ID).Find(&tran).Error
		if err != nil {
			logr.Info("db tran_logs find sys_order_id failed:%s\n", err)
		}
		logr.Infof("sys_order_id=%s, cust_order_id=%s", b.SYS_ID, tran.CUST_ORDER_ID)
		b.CUST_ORDER_ID = tran.CUST_ORDER_ID
		b.EXT_FLD = tran.Ext_fld7
		cf.FileStrt.FileBodys = append(cf.FileStrt.FileBodys, b)
	}

	//cf.FileStrt.FileHeadInfo.TrnSucCount = strconv.Itoa(record)
	cf.FileStrt.FileHeadInfo.Stlm_date = cf.STLM_DATE
	cf.FileStrt.FileHeadInfo.TrnSucAm = strconv.FormatFloat(trans_amt_T, 'f', 2, 64)
	cf.FileStrt.FileHeadInfo.TrnFeeT = strconv.FormatFloat(true_fee_mod_T, 'f', 2, 64)
	cf.FileStrt.FileHeadInfo.TrnReconT = strconv.FormatFloat(trnrecont_T, 'f', 2, 64)
	logr.Info("成功总笔数:", record)

	cf.FileStrt.FileHeadInfo.TrnSucCount = strconv.Itoa(record)
	return nil
}

func (cf *CrtFile) GetInsIdCd() (string, bool) {
	l := len(cf.Ins_id_cd)
	if l == 0 {
		return "", false
	}

	cf.MCHT_CD = cf.Ins_id_cd[0]
	cf.Ins_id_cd = cf.Ins_id_cd[1:]
	logr.Infof("取机构号：%s; 剩余机构号:%v", cf.MCHT_CD, cf.Ins_id_cd)

	return cf.MCHT_CD, true
}

func (cf *CrtFile) geneFile() string {
	cd, ok := cf.GetInsIdCd()
	if ok {
		cf.FileName = cf.FileName[:2] + cd
	} else {
		return ""
	}

	cf.FileName = cf.FileName + "_" + cf.STLM_DATE + ".txt"
	logr.Info("生成对账文件名称：", cf.FileName)
	p := cf.FilePath + cf.FileName
	return p
}

func (cf *CrtFile) InitMCHTCd(mc string) gerror.IError {

	//商户号
	dbc := gormdb.GetInstance()
	//rows, err := dbc.Raw("SELECT distinct MCHT_CD FROM tbl_mcht_recon_list").Rows()
	//rows, err := dbc.Raw("SELECT distinct MCHT_CD, mcht_ty FROM tbl_mcht_recon_list").Rows()
	tbrec := models.Tbl_mcht_recon_list{}
	err := dbc.Where("  MCHT_CD = ?", mc).Find(&tbrec).Error

	if err == gorm.ErrRecordNotFound {
		logr.Infof("商户[%s]对账信息未配置:%s\n", mc, err)
		return gerror.NewR(1000, err, "商户[%s]对账信息未配置:%s", mc, err)
	}
	if err != nil {
		logr.Infof("商户[%s]对账信息查询失败:%s\n", mc, err)
		return gerror.NewR(1000, err, "商户[%s]对账信息查询失败:%s", mc, err)
	}
	logr.Info("对账配置表:%+v\n", tbrec)
	if mc != "" {
		cf.Ins_id_cd = append(cf.Ins_id_cd, mc)
		cf.MCHT_TP[mc] = tbrec.Mcht_ty
	}

	logr.Info("初始化商户号:", cf.MCHT_TP)
	return nil
}
