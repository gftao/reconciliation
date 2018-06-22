package reconc

import (
	"golib/gerror"
	"golib/modules/gormdb"
	"github.com/jinzhu/gorm"
	"os"
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
	"fmt"
)

type CrtFunc1 struct {
	FileName string //对账文件名
	FilePath string //路径
	//SysDate              string                     //当前时间
	MCHT_CD              string                     //当前机构号
	STLM_DATE            string                     //清算日期
	Ins_id_cd            []string                   //全部机构号
	FileStrt             models.FileStrt1           //文件结构
	Tbl_Clear_Data       []models.Tbl_clear_txn     //清算表数据
	Tran_logs            []models.Tran_logs         //当前交易表数据
	Tbl_Tfr_his_log_Data models.Tbl_tfr_his_trn_log //当前机构号对应交易日志
	dbtype               string
	dbstr                string
	MCHT_TP              map[string]string
	sendto               bool
	recon                models.Tbl_mcht_recon_list
	TERM_ID              string
}

func (cf *CrtFunc1) Init(chainName string, mc string) gerror.IError {
	cf.FileName = "spdb_"

	cf.FilePath = config.StringDefault("filePath", "")
	cf.sendto = config.BoolDefault("sendto", true)
	logr.Info("是否发送ftp：", cf.sendto)

	cf.FileStrt = models.FileStrt1{}
	cf.STLM_DATE = chainName //清算日期
	cf.MCHT_TP = make(map[string]string, 1)
	//查表获取需要生产队长文件机构的机构号,新增加的表
	cf.MCHT_CD = mc
	cf.InitMCHTCd(mc)
	cf.indb()
	return nil
}

func (cf *CrtFunc1) indb() {
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

func (cf *CrtFunc1) Run() {
	err := cf.SaveToFile()
	if err != nil {
		logr.Error(err)
		return
	}
	err = cf.DoF587()
	if err != nil {
		logr.Error(err)
		return
	}
	return
}

func (cf *CrtFunc1) SaveToFile() error {
	//读数据
	err := cf.ReadDate()
	if err != nil {
		return err
	}
	//读取数据成功，创建文件
	fp := cf.geneFile()
	logr.Info("对账文件路径：", fp)
	f, err := os.Create(fp)
	defer f.Close()
	if err != nil {
		logr.Info(cf.FileName, err)
		return err
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
	rb := b.Bytes()
	logr.Infof("<<读取的数据>>[%s]", string(rb))

	err = cf.postToSftp(cf.FileName, rb)
	if err != nil {
		logr.Error(err)
	}

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

func (cf *CrtFunc1) postToSftp(fileName string, fileData []byte) error {

	dbc := gormdb.GetInstance()

	rows, err := dbc.Raw("SELECT * FROM tbl_mcht_recon_list WHERE MCHT_CD = ?", cf.MCHT_CD).Rows()
	if err != nil {
		logr.Info("dbc find failed:", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		tmr := &models.Tbl_mcht_recon_list{}
		dbc.ScanRows(rows, tmr)
		if tmr.USER == "" || tmr.PASSWD == "" || tmr.HOST == "" {
			logr.Infof("读配置错误:[%s][%s][%s]", tmr.USER, tmr.PASSWD, tmr.HOST)
			return fmt.Errorf("读配置错误:[%s][%s][%s]", tmr.USER, tmr.PASSWD, tmr.HOST)
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
					return err
				}
				AESinfo := security.EncodeBase64(cipherdata)
				fileData = AESinfo
				logr.Infof("--Aes加密后的数据---[%s]", fileData)
			}

			switch trans_ty {
			case "0":
				logr.Infof("SFTP:")
				err = myfstp.PosByteSftp(user, password, host, port, fileName, rmtDir, fileData)
				if err != nil {
					return err
				}
			case "1":
				logr.Infof("FTP with TLS:")
				err = myftp.MyftpTSL(user, password, host, port, fileName, rmtDir, fileData)
				if err != nil {
					logr.Error(err)
					return err
				}
			case "2":
				logr.Infof("FTP without TLS:")
				err = myftp.Myftp(user, password, host, port, fileName, rmtDir, fileData)
				if err != nil {
					logr.Error(err)
					return err
				}
			}
		}
	}
	return nil
}

func (cf *CrtFunc1) Finish() {
	return
}

func (cf *CrtFunc1) ReadDate() error {
	dbc := gormdb.GetInstance()
	var rows *sql.Rows
	var err error
	logr.Info(cf.MCHT_CD, "->", cf.MCHT_TP[cf.MCHT_CD])
	if cf.MCHT_TP[cf.MCHT_CD] == "0" {
		rows, err = dbc.Raw("SELECT * FROM tbl_clear_txn WHERE MCHT_CD = ? and STLM_DATE = ?", cf.MCHT_CD, cf.STLM_DATE).Rows()
	} else if cf.MCHT_TP[cf.MCHT_CD] == "1" {
		rows, err = dbc.Raw("SELECT * FROM tbl_clear_txn WHERE JT_MCHT_CD = ? and STLM_DATE = ?", cf.MCHT_CD, cf.STLM_DATE).Rows()
	} else {
		return fmt.Errorf("商户类型不支持")
	}
	defer rows.Close()
	if err == gorm.ErrRecordNotFound {
		logr.Info("tbl_clear_txn not find: ", err)
		return err
	} else if err != nil {
		logr.Info("tbl_clear_txn find fail: ", err)
		return err
	}
	for rows.Next() {
		tc := models.Tbl_clear_txn{}
		dbc.ScanRows(rows, &tc)

		cf.Tbl_Clear_Data = append(cf.Tbl_Clear_Data, tc)
	}
	if len(cf.Tbl_Clear_Data) <= 0 {
		return fmt.Errorf("MCHT_CD[%s]记录不存在", cf.MCHT_CD)
	}
	err = cf.saveDatatoFStru()
	if err != nil {
		return err
	}
	return nil
}

func (cf *CrtFunc1) saveDatatoFStru() error {
	cf.FileStrt.FileBodys = []models.Body1{}
	dbc := gormdb.GetInstance()
	dbt, err := gorm.Open(cf.dbtype, cf.dbstr)
	if err != nil {
		logr.Info("open db err:", err)
		return err
	}
	defer dbt.Close()
	dbt = dbt.Set("gorm:table_options", "ENGINE=InnoDB")
	dbt.DB().Ping()
	record := 0        //交易总笔数
	trans_amt_T := 0.0 //清算金额

	for _, tc := range cf.Tbl_Clear_Data {
		b := models.Body1{}
		tfr := models.Tbl_tfr_his_trn_log{}
		tran := models.Tran_logs{}
		oritran := models.Tran_logs{}
		err := dbc.Where("KEY_RSP = ?", tc.KEY_RSP).Find(&tfr).Error
		if err != nil {
			logr.Infof("dbc find failed, KEY_RSP = %s, err = %s", tc.KEY_RSP, err)
			continue
		}

		record ++
		a, _ := strconv.ParseFloat(tc.TRANS_AMT, 64)
		trans_amt_T += a

		//cf.TERM_ID = tc.TERM_ID
		b.AMT = tc.TRANS_AMT
		b.TYPE = tc.TXN_DESC
		b.DATE = tfr.TRANS_DT[:4] + "-" + tfr.TRANS_DT[4:6] + "-" + tfr.TRANS_DT[6:]
		b.POS_ID = tc.TERM_ID

		SYS_ID := ""
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

		err = dbt.Where("sys_order_id = ?", SYS_ID).Find(&tran).Error
		if err != nil {
			logr.Info("tran_logs查询原始交易失败sys_order_id = %s\n", err)
		}
		logr.Infof("sys_order_id=%s, cust_order_id=%s", SYS_ID, tran.CUST_ORDER_ID)

		b.TRACE_NO = tran.CUST_ORDER_ID
		if tran.TRAN_CD == "3011" || tran.TRAN_CD == "2011" {
			err = dbt.Where("sys_order_id = ?", tran.ORIG_SYS_ORDER_ID).Find(&oritran).Error
			if err != nil {
				logr.Info("tran_logs查询消费原始交易失败sys_order_id = %s\n", err)
			}
			u := strings.Split(strings.TrimSpace(oritran.Ext_fld7), "|")
			if len(u) == 2 {
				b.USER_NUMBER = u[0]
				b.CHARGE_YEAR_ID = u[1]
			}
		} else {
			u := strings.Split(strings.TrimSpace(tran.Ext_fld7), "|")
			if len(u) == 2 {
				b.USER_NUMBER = u[0]
				b.CHARGE_YEAR_ID = u[1]
			}
		}

		cf.FileStrt.FileBodys = append(cf.FileStrt.FileBodys, b)
	}

	/*
	for _, tl := range cf.Tran_logs {
		b := models.Body1{}
		record ++
		a, _ := strconv.ParseFloat(tl.TRAN_AMT, 64)
		if tl.TRAN_CD == "1151" {
			a /= 100
		} else if tl.TRAN_CD == "3151" {
			a /= 100
			a *= (-1)
			for _, tll := range cf.Tran_logs {
				if tl.ORIG_SYS_ORDER_ID == tll.SYS_ORDER_ID {
					tl.Ext_fld7 = tll.Ext_fld7
				}
			}
		}
		trans_amt_T += a

		b.TRACE_NO = tl.CUST_ORDER_ID
		b.AMT = strconv.FormatFloat(a, 'f', 2, 64)
		b.TYPE = tl.TRAN_NM
		b.DATE = tl.TRANS_DT[:4] + "-" + tl.TRANS_DT[4:6] + "-" + tl.TRANS_DT[6:]
		b.POS_ID = tl.TERM_ID
		u := strings.Split(strings.TrimSpace(tl.Ext_fld7), "|")
		if len(u) == 2 {
			b.USER_NUMBER = u[0]
			b.CHARGE_YEAR_ID = u[1]
		}
		cf.FileStrt.FileBodys = append(cf.FileStrt.FileBodys, b)
	}
	*/

	cf.FileStrt.FileHeadInfo.TrnSucCount = strconv.Itoa(record)
	cf.FileStrt.FileHeadInfo.TrnSucAm = strconv.FormatFloat(trans_amt_T, 'f', 2, 64)
	logr.Info("成功总笔数:", record)

	return nil
}

func (cf *CrtFunc1) GetInsIdCd() (string, bool) {
	logr.Infof("取机构号：%s", cf.MCHT_CD)

	return cf.MCHT_CD, true
}

func (cf *CrtFunc1) geneFile() string {

	cf.FileName = cf.FileName + cf.STLM_DATE + ".txt"
	logr.Info("生成对账文件名称：", cf.FileName)
	p := cf.FilePath + cf.FileName
	return p
}

func (cf *CrtFunc1) InitMCHTCd(mc string) {

	//商户号
	dbc := gormdb.GetInstance()
	err := dbc.Where("  MCHT_CD = ?", mc).Find(&cf.recon).Error

	if err == gorm.ErrRecordNotFound {
		logr.Info("dbc.Raw fail:%s\n", err)
		return
	}
	if err != nil {
		logr.Info("dbc.Raw fail:%s\n", err)
		return
	}
	logr.Info("对账配置表:%+v\n", cf.recon)
	if mc != "" {
		cf.Ins_id_cd = append(cf.Ins_id_cd, mc)
		cf.MCHT_TP[mc] = cf.recon.Mcht_ty
	}

	logr.Info("初始化商户号:", cf.MCHT_TP)

}

func (cf *CrtFunc1) DoF587() error {
	fn := "checkAccountF587"
	req := &models.F587{}
	req.SIGN_TYPE = "01"
	area := models.F587_AREA{}
	area.TRAN_CODE = "F587"
	area.POS_ID = "0"
	area.DATE = cf.STLM_DATE
	area.FILE_NAME = cf.FileName
	area.BANKNO = "01"
	req.DATA_AREA = area
	sig, err := req.BuilMsg(req.DATA_AREA)
	if err != nil {
		return err
	}
	logr.Infof("xml=%s", sig)
	req.SIGN, err = req.SingUp(cf.dbtype, cf.dbstr, sig)
	if err != nil {
		return err
	}
	//logr.Debugf("F587=%+v", *req)
	reqMsg, err := req.BuilMsg(req)
	if err != nil {
		return err
	}
	logr.Debugf("F587=%+v", string(reqMsg))
	b, err := models.BuilSoap(fn, string(reqMsg))
	if err != nil {
		return err
	}
	//logr.Infof("BuilSoap xml=%s", b)
	rsp, err := models.Comm(b)
	if err != nil {
		return err
	}
	//logr.Infof("BuilSoap  rsp xml=%s", rsp)
	rs, err := models.LoadResp(rsp)
	if err != nil {
		return err
	}
	logr.Debugf("Resp=%s", rs)

	return nil
}
