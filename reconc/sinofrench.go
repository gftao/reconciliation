package reconc

import (
	"htdRec/models"
	"golib/modules/logr"
	"golib/modules/config"
	"github.com/jinzhu/gorm"
	"golib/gerror"
	"os"
	"bytes"
	"bufio"
	"golib/modules/gormdb"
	"database/sql"
	"strconv"
	"strings"
	"golib/security"
	"htdRec/myfstp"
	"htdRec/myftp"
)

type SinoFrench struct {
	FileName  string //对账文件名
	FilePath  string //路径
	SysDate   string //当前时间
	MCHT_CD   string //当前机构号
	STLM_DATE string //清算日期
	//Ins_id_cd            []string                   //全部机构号
	FileStrt             models.FileSinoFrench      //文件结构
	Tbl_Clear_Data       []models.Tbl_clear_txn     //当前机构号数据
	Tbl_Tfr_his_log_Data models.Tbl_tfr_his_trn_log //当前机构号对应交易日志
	dbtype               string
	dbstr                string
	MCHT_TP              string // map[string]string
	sendto               bool
}

func (cf *SinoFrench) Init(chainName string, mc string) gerror.IError {
	cf.FileName = "spdb02_"
	cf.FilePath = config.StringDefault("filePath", "")
	cf.sendto = config.BoolDefault("sendto", true)
	logr.Info("是否发送ftp：", cf.sendto)
	cf.STLM_DATE = chainName //清算日期
	cf.MCHT_CD = mc
	err := cf.InitMCHTCd(mc)
	if err != nil {
		return err
	}
	cf.indb()
	return nil
}

func (cf *SinoFrench) indb() {
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
func (cf *SinoFrench) Run() {
	gerr := cf.SaveToFile()
	if gerr != nil {
		logr.Errorf("%s", gerr)
	}
	return
}
func (cf *SinoFrench) InitMCHTCd(mc string) gerror.IError {
	//商户号
	dbc := gormdb.GetInstance()
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
	logr.Infof("商户对账配置表:%+v\n", tbrec)

	cf.MCHT_TP = tbrec.Mcht_ty

	logr.Infof("初始化商户号:[%+v][%s]", mc, cf.MCHT_TP)
	return nil
}
func (cf *SinoFrench) SaveToFile() gerror.IError {
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
	b.WriteString(cf.FileStrt.HToString())
	b.WriteString("\r\n")
	for _, tc := range cf.FileStrt.FileBodys {
		Str := tc.BToString()
		b.WriteString(Str)
		b.WriteString("\r\n")
	}
	logr.Infof("--读取的数据---[%s]", b)
	rb := b.Bytes()
	rb = bytes.TrimSpace(rb)
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
	w.Flush()
	f.Sync()
	return nil
}
func (cf *SinoFrench) ReadDate() gerror.IError {
	dbc := gormdb.GetInstance()
	var rows *sql.Rows
	var err error
	logr.Info("中法户号:", cf.MCHT_CD, "->", cf.MCHT_TP)
	if cf.MCHT_TP == "0" {
		rows, err = dbc.Raw("SELECT * FROM tbl_clear_txn WHERE MCHT_CD = ? and STLM_DATE = ?", cf.MCHT_CD, cf.STLM_DATE).Rows()
	} else if cf.MCHT_TP == "1" {
		rows, err = dbc.Raw("SELECT * FROM tbl_clear_txn WHERE JT_MCHT_CD = ? and STLM_DATE = ?", cf.MCHT_CD, cf.STLM_DATE).Rows()
	} else {
		return gerror.NewR(1000, err, "商户[%s]类型错误[%s]", cf.MCHT_CD, cf.MCHT_TP)
	}
	defer rows.Close()

	if err == gorm.ErrRecordNotFound {
		logr.Info("商户[%s]没有对账数据", cf.MCHT_CD)
		return gerror.NewR(1000, err, "商户[%s]没有对账数据", cf.MCHT_CD, err)
	}
	if err != nil {
		logr.Info("查商户[%s]对账数据失败:%s", cf.MCHT_CD, err)
		return gerror.NewR(1000, err, "查商户[%s]对账数据失败:%s", cf.MCHT_CD, err)
	}

	for rows.Next() {
		tc := models.Tbl_clear_txn{}
		dbc.ScanRows(rows, &tc)
		cf.Tbl_Clear_Data = append(cf.Tbl_Clear_Data, tc)
	}
	if len(cf.Tbl_Clear_Data) > 0 {
		cf.FileStrt.FileHeadInfo.INS_ID_CD = cf.Tbl_Clear_Data[0].INS_ID_CD
	}
	cf.FileStrt.FileHeadInfo.Stlm_date = cf.STLM_DATE
	cf.FileStrt.FileHeadInfo.TimeB = cf.STLM_DATE[:4] + "-" + cf.STLM_DATE[4:6] + "-" + cf.STLM_DATE[6:] + " 00:00:00"
	cf.FileStrt.FileHeadInfo.TimeE = cf.STLM_DATE[:4] + "-" + cf.STLM_DATE[4:6] + "-" + cf.STLM_DATE[6:] + " 23:59:59"

	gerr := cf.saveDatatoFStru()
	if gerr != nil {
		return gerr
	}
	return nil
}

func (cf *SinoFrench) geneFile() string {
	cf.FileName = cf.FileName + cf.STLM_DATE + ".txt"
	logr.Info("生成对账文件名称：", cf.FileName)
	p := cf.FilePath + cf.FileName
	return p
}
func (cf *SinoFrench) saveDatatoFStru() gerror.IError {
	//cf.FileStrt =  models.FileSinoFrench{}
	dbc := gormdb.GetInstance()
	dbt, err := gorm.Open(cf.dbtype, cf.dbstr)
	if err != nil {
		logr.Info(err)
		return gerror.NewR(1000, err, "打开数据库链接失败")
	}
	defer dbt.Close()
	dbt = dbt.Set("gorm:table_options", "ENGINE=InnoDB")
	dbt.DB().Ping()
	record := 0        //交易总笔数
	trans_amt_T := 0.0 //清算金额

	for _, tc := range cf.Tbl_Clear_Data {
		b := models.SFBody{}
		tfr := models.Tbl_tfr_his_trn_log{}
		tran := models.Tran_logs{}
		err := dbc.Where("KEY_RSP = ?", tc.KEY_RSP).Find(&tfr).Error
		if err != nil {
			logr.Infof("dbc find failed, KEY_RSP = %s, err = %s", tc.KEY_RSP, err)
			continue
		}
		record ++
		a, _ := strconv.ParseFloat(tc.TRANS_AMT, 64)

		trans_amt_T += a
		b.TERM_ID = tc.TERM_ID
		logr.Info("prod_cd:", tfr.PROD_CD)
		var sysId string
		if tfr.PROD_CD == "1151" {
			sysId = tfr.INDUSTRY_ADDN_INF

		} else {
			if tfr.RETRI_REF_NO[:1] == cf.STLM_DATE[3:4] {
				sysId = cf.STLM_DATE[:3] + tfr.RETRI_REF_NO
			} else {
				sysId = cf.STLM_DATE[:4] + tfr.RETRI_REF_NO
			}
		}
		b.TRANS_AMT = tc.TRANS_AMT

		err = dbt.Where("sys_order_id = ?", sysId).Find(&tran).Error
		if err != nil {
			logr.Info("查询sys_order_id失败:%s\n", err)
		}
		logr.Infof("sys_order_id=%s, cust_order_id=%s", sysId, tran.CUST_ORDER_ID)
		b.CUST_ORDER_ID = tran.CUST_ORDER_ID

		//b.Regist_Meter_No = "注册号"
		//b.Amount = "水量"

		cf.FileStrt.FileBodys = append(cf.FileStrt.FileBodys, b)
	}
	cf.FileStrt.FileHeadInfo.TrnSucCount = strconv.Itoa(record)
	cf.FileStrt.FileHeadInfo.TrnReconT = strconv.FormatFloat(trans_amt_T, 'f', 2, 64)

	return nil
}
func (cf *SinoFrench) postToSftp(fileName string, fileData []byte) {

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
