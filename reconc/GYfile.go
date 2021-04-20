package reconc

import (
	"github.com/jinzhu/gorm"
	"golib/gerror"
	"golib/modules/config"
	"golib/modules/gormdb"
	"golib/modules/logr"
	"htdRec/models"
	"htdRec/myfstp"
	"htdRec/myftp"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type GYFile struct {
	FileName             string                     //对账文件名
	FilePath             string                     //路径
	SysDate              string                     //当前时间
	MCHT_CD              string                     //当前机构号
	STLM_DATE            string                     //清算日期
	Ins_id_cd            []string                   //全部机构号
	FileStrt             models.GYFileStrt     //文件结构
	Tbl_Clear_Data       []models.Tbl_clear_txn     //当前机构号数据
	Tbl_Tfr_his_log_Data models.Tbl_tfr_his_trn_log //当前机构号对应交易日志
	dbtype               string
	dbstr                string
	MCHT_TP              string // map[string]string
	sendto               bool
	dbt                  *gorm.DB
}

func (cf *GYFile) Init(chainName string, mc string) gerror.IError {
	cf.FileName = "gazjj_ysk_pfyh_"
	cf.FilePath = config.StringDefault("filePath", "")
	cf.sendto = config.BoolDefault("sendto", true)
	logr.Info("是否发送ftp：", cf.sendto)
	cf.STLM_DATE = chainName //清算日期
	cf.FileStrt.Init()
	//cf.MCHT_TP = make(map[string]string, 1)
	cf.MCHT_CD = mc
	err := cf.InitMCHTCd(mc)
	if err != nil {
		return err
	}
	cf.indb()
	return nil
}

func (cf *GYFile) indb() {
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

	dbt, err := gorm.Open(cf.dbtype, cf.dbstr)
	if err != nil {
		logr.Info("打开数据库链接失败", err)
		return
	}
	//defer dbt.Close()
	dbt = dbt.Set("gorm:table_options", "ENGINE=InnoDB")
	cf.dbt = dbt
	cf.dbt.DB().Ping()
}

func (cf *GYFile) InitMCHTCd(mc string) gerror.IError {

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

	for _, c := range strings.Split(tbrec.EXT4, ",") {
		cf.Ins_id_cd = append(cf.Ins_id_cd, c)

	}
	cf.MCHT_TP = tbrec.Mcht_ty
	logr.Infof("初始化商户号:[%+v][%s]", cf.Ins_id_cd, cf.MCHT_TP)
	return nil
}

func (cf *GYFile) Run() {
	defer cf.dbt.Close()

	gerr := cf.SaveToFile()
	if gerr != nil {
		logr.Errorf("%s", gerr)
	}
	return
}

func (cf *GYFile) SaveToFile() gerror.IError {

	//创建文件
	fn := cf.geneFile()
	logr.Info("对账文件路径：", fn)
	f, err := os.Create(fn)

	if err != nil {
		return gerror.NewR(1001, err, "创建文件失败:%s", fn)
	}
	//读数据
	gerr := cf.ReadDate(f)
	if gerr != nil {
		logr.Info("未读到清算数据,不创建对账文件MCHT_CD:", cf.MCHT_CD)
		return gerr
	}
	f.Sync()
	//f.Close()
	//fp, err := os.Open(fn)
	//if err != nil {
	//	return gerror.NewR(1004, err, "打开文件失败:%s", fn)
	//}
	//defer fp.Close()

	f.Seek(0, 0)
	cf.postToSftp(cf.FileName, f)
	f.Close()
	return nil
}

func (cf *GYFile) ReadDate(fp *os.File) gerror.IError {
	dbc := gormdb.GetInstance()

	rows, err := dbc.Raw("SELECT * FROM tbl_clear_txn  WHERE STLM_DATE = ? and MCHT_CD in (?) ", cf.STLM_DATE, cf.Ins_id_cd).Rows()
	defer rows.Close()
	if err == gorm.ErrRecordNotFound {
		return gerror.NewR(1000, err, "查对账数据失败:%s", err)
	}
	if err != nil {
		return gerror.NewR(1000, err, "查对账数据失败:%s", err)
	}
	//.Println(cf.FileStrt.FileHead)
	fp.WriteString(cf.FileStrt.FileHead + "\r\n")
	logr.Infof("--读取贵阳住建数据---")
	for rows.Next() {
		tc := models.Tbl_clear_txn{}
		dbc.ScanRows(rows, &tc)
		if tc.KEY_RSP == "" {
			continue
		}
		b, gerr := cf.saveDatatoFStru(&tc)
		if gerr != nil {
			return gerr
		}
		if b == nil {
			continue
		}
		//logr.Infof("KEY_RSP=%s", tc.KEY_RSP)
		//logr.Infof("%s", b.ToString())
		logr.Infof("[%s]",b.HToString())
		fp.WriteString(b.HToString() + "|\n")
	}
	l, _ := fp.Seek(-1, 2)
	fp.Truncate(l)

	return nil
}

func (cf *GYFile) saveDatatoFStru(tc *models.Tbl_clear_txn) (*models.GYFileHeadInfo, gerror.IError) {
	cf.FileStrt.GYFileInfo = []models.GYFileHeadInfo{}
	dbc := gormdb.GetInstance()
	dbt, err := gorm.Open(cf.dbtype, cf.dbstr)
	if err != nil {
		logr.Info(err)
		return nil,gerror.NewR(1000, err, "打开数据库链接失败")
	}
	defer dbt.Close()
	dbt = dbt.Set("gorm:table_options", "ENGINE=InnoDB")
	dbt.DB().Ping()

	b := models.GYFileHeadInfo{}
	tfr := models.Tbl_tfr_his_trn_log{}
	tran := models.Tran_logs{}
	tran_noti := models.Tran_logs{}
	mcht := models.Tbl_mcht_bankaccount{}
	//bank := models.Tbl_bank_bin_inf{}
	/*err := dbc.Where("KEY_RSP = ?", tc.KEY_RSP).Find(&tfr).Error
	if err != nil {
		logr.Infof("dbc find failed, KEY_RSP = %s, err = %s", tc.KEY_RSP, err)
		continue
	}*/
	err = dbc.Where("owner_cd = ?", tc.MCHT_CD).Find(&mcht).Error
	if err != nil {
		logr.Infof("dbc find failed, owner_cd = %s, err = %s", tc.MCHT_CD, err)
	}
	var SYS_ID string
	err = dbc.Where("KEY_RSP = ?", tc.KEY_RSP).Find(&tfr).Error
	if err != nil {
		logr.Infof("dbc find failed, KEY_RSP = %s, err = %s", tc.KEY_RSP, err)
	}
	logr.Info("prod_cd:", tfr.PROD_CD)
	if tfr.PROD_CD == "1151" {
		SYS_ID = tfr.INDUSTRY_ADDN_INF
	} else {
		if tfr.RETRI_REF_NO[:1] == cf.STLM_DATE[3:4] {
			SYS_ID = cf.STLM_DATE[:3] + tfr.RETRI_REF_NO
		} else {
			SYS_ID = cf.STLM_DATE[:4] + tfr.RETRI_REF_NO
		}
	}
	//fmt.Println(SYS_ID)
	err = dbt.Where("sys_order_id = ?", SYS_ID).Find(&tran).Error
	if err != nil {
		logr.Info("db tran_logs find sys_order_id failed:%s\n", err)
	}
	tran_cd := "6201"
	err = dbt.Where("orig_sys_order_id = ? and tran_cd = ?", SYS_ID, tran_cd).Find(&tran_noti).Error
	if err != nil {
		logr.Info("db tran_logs find orig_sys_order_id failed:%s\n", err)
	}

	b.STLM_DATE = cf.STLM_DATE
	b.TRAN_CD = tran.TRAN_CD
	b.JC_BIZ_CD = tran.CUST_ORDER_ID
	b.MCHT_CD = tran.MCHT_CD
	b.TERM_ID = tran.TERM_ID
	set_amt := strings.Split(tc.MCHT_SET_AMT,".")
	for _,v := range set_amt{
		b.PAY_AMT += v
	}
	//set_amt_str := set_amt[0]+set_amt[1]
/*	set_amt,_ := strconv.ParseFloat(tc.MCHT_SET_AMT,64)
	set_amt_str := strconv.FormatFloat(set_amt*100,'f',-1,64)*/
	//amt,_ := strconv.ParseFloat(tran.TRAN_AMT,64)
	//b.PAY_AMT = fmt.Sprintf("%.2f",amt/100)
	b.SYS_ORDER_ID = tran.SYS_ORDER_ID
	b.ZZJ_ORDER_ID = tran_noti.INS_ORDER_ID
	b.ACC_AMOUNT = tran.PRI_ACCT_NO
	b.TRAN_TM = tran.TRAN_DT_TM[8:]
	b.JC_KIND = "1"
	b.JC_TYPE = "2"
	b.BANK_CD = ""
	b.SUP_AMOUNT = ""
	return &b, nil
}

func (cf *GYFile) geneFile() string {
	cf.FileName = cf.FileName + cf.STLM_DATE+".txt"
	logr.Info("生成对账文件名称：", cf.FileName)
	p := cf.FilePath + cf.FileName
	return p
}

func (cf *GYFile) GetInsIdCd() (string, bool) {
	l := len(cf.Ins_id_cd)
	if l == 0 {
		return "", false
	}

	cf.MCHT_CD = cf.Ins_id_cd[0]
	cf.Ins_id_cd = cf.Ins_id_cd[1:]
	logr.Infof("取机构号：%s; 剩余机构号:%v", cf.MCHT_CD, cf.Ins_id_cd)

	return cf.MCHT_CD, true
}

func (cf *GYFile) postToSftp(fileName string, fileData io.Reader) {

	data, err := ioutil.ReadAll(fileData)
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

			switch trans_ty {
			case "0":
				logr.Infof("SFTP:")
				err = myfstp.PosIOSftp(user, password, host, port, fileName, rmtDir, fileData)
				if err != nil {
					logr.Error(err)
				}
				err = myfstp.PosByteSftp(user, password, host, port, fileName+".finish", rmtDir, []byte{})
				if err != nil {
					logr.Error(err)
				}
			case "1":
				logr.Infof("FTP with TLS:")
				err = myftp.MyftpTSL(user, password, host, port, fileName, rmtDir, data)
				if err != nil {
					logr.Error(err)
				}
				err = myftp.MyftpTSL(user, password, host, port, fileName+".finish", rmtDir, []byte{})
				if err != nil {
					logr.Error(err)
				}
			case "2":
				logr.Infof("FTP without TLS:")
				err = myftp.Myftp(user, password, host, port, fileName, rmtDir, data)
				if err != nil {
					logr.Error(err)
				}
				err = myftp.Myftp(user, password, host, port, fileName+".finish", rmtDir, []byte{})
				if err != nil {
					logr.Error(err)
				}
			default:
				logr.Infof("default FTP is not support")
			}
		}
	}
}

