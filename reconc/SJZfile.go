package reconc

import (
	"bufio"
	"bytes"
	"golib/gerror"
	"golib/modules/config"
	"golib/modules/gormdb"
	"htdRec/models"
	"htdRec/myfstp"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"

	"database/sql"
	"golib/modules/logr"
	"golib/security"
	"htdRec/myftp"
	"os"
	"time"
)

type ShijiazhuangZJFile struct {
	FileName             string                     //对账文件名
	FilePath             string                     //路径
	SysDate              string                     //当前时间
	MCHT_CD              string                     //当前机构号
	STLM_DATE            string                     //清算日期
	Pay_DATE_S           string                     //划付日期
	Pay_DATE_E           string                     //划付日期
	Ins_id_cd            []string                   //全部机构号
	FileStrt             models.SJZFileStrt         //文件结构
	Tbl_Clear_Data       []models.Tbl_clear_txn     //当前机构号数据
	Tbl_Tfr_his_log_Data models.Tbl_tfr_his_trn_log //当前机构号对应交易日志
	dbtype               string
	dbstr                string
	MCHT_TP              string
	sendto               bool
	empty                bool //1-创建空文件
	Area_cd              string
	Pay_DATE             bool //是否为划付日
}

func (cf *ShijiazhuangZJFile) Init(chainName string, mc string) gerror.IError {
	cf.FilePath = config.StringDefault("filePath", "")
	cf.sendto = config.BoolDefault("sendto", true)
	logr.Info("是否发送ftp：", cf.sendto)

	cf.SysDate = time.Now().Format("20060102") //今天
	cf.FileStrt = models.SJZFileStrt{}
	cf.STLM_DATE = chainName //清算日期
	//cf.MCHT_TP = make(map[string]string, 1)
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

func (cf *ShijiazhuangZJFile) indb() {
	//"prodPmpCld:prodPmpCld@tcp(192.168.20.60:3306)/prodPmpCld?charset=utf8&parseTime=True&loc=Local"
	config.SetSection("db1")
	dbtype := config.StringDefault("db.type", "mysql")
	dbhost := config.StringDefault("db.host", "127.0.0.1")
	dbport := config.StringDefault("db.port", "3306")
	dbname := config.StringDefault("db.dbname", "prodPmpCld")
	dbuser := config.StringDefault("db.user", "root")
	dbpasswd := config.StringDefault("db.passwd", "")
	//dbpasswd = models.DbpsDecrypt(dbpasswd)
	connStr := dbuser + ":" + dbpasswd + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname +
		"?charset=utf8&parseTime=True&loc=Local"

	cf.dbtype = dbtype
	cf.dbstr = connStr
	logr.Info("[db1]：", cf.dbtype, cf.dbstr)
}

func (cf *ShijiazhuangZJFile) Run() {

	gerr := cf.SaveToFile()
	if gerr != nil {
		logr.Errorf("%s", gerr)
	}
	return
}
func (cf *ShijiazhuangZJFile) SaveToFile() gerror.IError {
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
	l := len(cf.FileStrt.FileBodys) - 1
	for i, tc := range cf.FileStrt.FileBodys {
		Str := tc.BToString()
		b.WriteString(Str)
		if i < l {
			b.WriteString("\r\n")
		}
	}
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

func (cf *ShijiazhuangZJFile) postToSftp(fileName string, fileData []byte) {

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
				err = myfstp.PosByteSftp(user, password, host, port, fileName, rmtDir, fileData)
				if err != nil {
					logr.Error(err)
				}
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

func (cf *ShijiazhuangZJFile) Finish() {
	return
}

func (cf *ShijiazhuangZJFile) ReadDate() gerror.IError {
	dbc := gormdb.GetInstance()
	var rows *sql.Rows
	var err error
	//logr.Info(cf.MCHT_CD, "->", cf.MCHT_TP[cf.MCHT_CD])
	if cf.Pay_DATE {
		rows, err = dbc.Raw("SELECT * FROM tbl_clear_txn WHERE MCHT_CD in (?) and (STLM_DATE between ? and ? )", cf.Ins_id_cd, cf.Pay_DATE_S, cf.STLM_DATE).Rows()
	} else {
		rows, err = dbc.Raw("SELECT * FROM tbl_clear_txn WHERE MCHT_CD in (?) and STLM_DATE = ?", cf.Ins_id_cd, cf.STLM_DATE).Rows()
	}

	defer rows.Close()
	if err == gorm.ErrRecordNotFound && !cf.empty {
		return gerror.NewR(0000, err, "MCHT_CD[%s]记录不存在", cf.MCHT_CD)
	}
	if err != nil {
		logr.Info("tbl_clear_txn find fail: ", err)
		return gerror.NewR(1000, err, " ")
	}
	for rows.Next() {
		tc := models.Tbl_clear_txn{}
		dbc.ScanRows(rows, &tc)
		cf.Tbl_Clear_Data = append(cf.Tbl_Clear_Data, tc)
	}
	//fmt.Println(cf.Tbl_Clear_Data)
	if len(cf.Tbl_Clear_Data) < 0 {
		return gerror.NewR(0000, err, "MCHT_CD[%s]记录不存在", cf.MCHT_CD)
	}

	gerr := cf.saveDatatoFStru()
	if gerr != nil {
		return gerr
	}
	return nil
}

func (cf *ShijiazhuangZJFile) saveDatatoFStru() gerror.IError {
	//cf.FileStrt.FileBodys = make([]models.Body,0)
	cf.FileStrt.FileBodys = []models.SJZBody{}
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
	trnrecont_T := 0.0 //结算总金额
	for _, tc := range cf.Tbl_Clear_Data {
		b := models.SJZBody{}
		tfr := models.Tbl_tfr_his_trn_log{}
		tran := models.Tran_logs{}
		tran_noti := models.Tran_logs{}
		err := dbc.Where("KEY_RSP = ?", tc.KEY_RSP).Find(&tfr).Error
		if err != nil {
			logr.Infof("dbc find failed, KEY_RSP = %s, err = %s", tc.KEY_RSP, err)
			continue
		}
		m, _ := strconv.ParseFloat(tc.MCHT_SET_AMT, 64)
		trnrecont_T += m

		logr.Info("prod_cd:", tfr.PROD_CD)
		if tfr.PROD_CD == "1151" {
			b.SysId = tfr.INDUSTRY_ADDN_INF
		} else {
			if tfr.RETRI_REF_NO[:1] == cf.STLM_DATE[3:4] {
				b.SysId = cf.STLM_DATE[:3] + tfr.RETRI_REF_NO
			} else {
				b.SysId = cf.STLM_DATE[:3] + tfr.RETRI_REF_NO
			}
		}
		TRAND_CD := tfr.MA_TRANS_CD
		switch TRAND_CD[:1] {
		case "1", "2", "3":
			err = dbt.Where("sys_order_id = ?", b.SysId).Find(&tran).Error
			if err != nil {
				logr.Info("db tran_logs find sys_order_id failed:%s\n", err)
			}
			logr.Infof("sys_order_id=%s, cust_order_id=%s", b.SysId, tran.CUST_ORDER_ID)
			//continue
		}
		tran_cd := "6201"
		err = dbt.Where("orig_sys_order_id = ? and tran_cd = ?", b.SysId, tran_cd).Find(&tran_noti).Error
		if err != nil {
			logr.Info("db tran_logs find orig_sys_order_id failed:%s\n", err)
		}
		b.TranDate = tfr.TRANS_DT
		b.TranTime = tfr.TRANS_MT
		b.FrontDate = tfr.TRANS_DT
		b.FrontId = b.SysId
		b.Currency = "01"
		b.CurrFlag = "1"
		b.OutType = ""
		b.OutAccount = tran.PRI_ACCT_NO
		b.InType = ""
		b.InAccount = tran.CUST_ORDER_ID //监管账号
		b.TranAmt = tc.MCHT_SET_AMT
		b.ThirdId = tran_noti.INS_ORDER_ID //销账返回
		b.ThirdNo = ""
		b.PayPeriod = "0"
		b.TranTunnel = "03"
		b.UserAddress = ""
		b.Ext1 = ""
		b.Ext2 = ""
		b.Ext3 = ""
		b.Ext4 = ""
		b.Ext5 = ""

		record++
		cf.FileStrt.FileBodys = append(cf.FileStrt.FileBodys, b)
	}
	cf.FileStrt.SJZFileHeadInfo.Be_biz_cd = "450016002021051"
	cf.FileStrt.SJZFileHeadInfo.Pay_flag = "1"
	cf.FileStrt.SJZFileHeadInfo.Total_amt = strconv.FormatFloat(trnrecont_T, 'f', 2, 64) //交易总结算额
	logr.Info("成功总笔数:", record)
	cf.FileStrt.SJZFileHeadInfo.Total_tran = strconv.Itoa(record) //成功总笔数
	return nil
}

func (cf *ShijiazhuangZJFile) GetInsIdCd() (string, bool) {
	l := len(cf.Ins_id_cd)
	if l == 0 {
		return "", false
	}

	cf.MCHT_CD = cf.Ins_id_cd[0]
	cf.Ins_id_cd = cf.Ins_id_cd[1:]
	logr.Infof("取机构号：%s; 剩余机构号:%v", cf.MCHT_CD, cf.Ins_id_cd)

	return cf.MCHT_CD, true
}

func (cf *ShijiazhuangZJFile) geneFile() string {
	cd, ok := cf.GetInsIdCd()
	if ok {
		//cf.FileName = "YSFPOS_" + cd
		cf.FileName = cd //中间业务编码
	} else {
		return ""
	}
	logr.Info("商户号：", cd)
	cf.FileName = cf.FileName + "_" + cf.Pay_DATE_E + ".txt"
	logr.Info("生成对账文件名称：", cf.FileName)
	p := cf.FilePath + cf.FileName
	return p
}

func (cf *ShijiazhuangZJFile) InitMCHTCd(mc string) gerror.IError {
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
	logr.Info("对账配置表:%+v\n", tbrec)
	/*	if mc != "" {
		cf.Ins_id_cd = append(cf.Ins_id_cd, mc)
		cf.MCHT_TP[mc] = tbrec.Mcht_ty
	}*/
	if strings.ContainsAny(tbrec.EXT3, "[]") {
		tbrec.EXT3 = strings.Trim(tbrec.EXT3, "[]")
		v := strings.Split(tbrec.EXT3, ",")
		cf.empty = v[0] == "1"
		cf.Area_cd = v[1]
	}
	th := models.TBL_HOLI_INF{}
	od, _ := time.ParseDuration("24h")
	std, err := time.ParseInLocation("20060102", cf.STLM_DATE, time.Local)
	if err != nil {
		return gerror.NewR(1001, err, "\033[0;31m"+"假期解析失败"+"\033[0m \n")
	}
	py := std.Add(od).Format("20060102")

	err = dbc.Where("START_DATE <= ? and ? <= END_DATE", py, py).Find(&th).Error
	if err == gorm.ErrRecordNotFound {
		cf.Pay_DATE_S = py
		cf.Pay_DATE_E = py
		logr.Infof("商户[%s]划付日期为:%s", mc, cf.Pay_DATE_E)
	} else if err != nil {
		logr.Infof("商户[%s]划付日期查询失败:%s", mc, err)
		return gerror.NewR(1000, err, "商户[%s]划付日期查询失败:%s", mc, err)
	} else {
		//判断是否为划付日
		if py == th.END_DATE {
			cf.Pay_DATE = true
		} else {
			return gerror.NewR(1001, err, "商户[%s]未到划付日期:%s", mc, th.END_DATE)
		}
		td, err := time.ParseInLocation("20060102", th.START_DATE, time.Local)
		if err != nil {
			return gerror.NewR(1001, err, "\033[0;31m"+"假期解析失败"+"\033[0m \n")
		}

		dd, _ := time.ParseDuration("-24h")
		cf.Pay_DATE_S = td.Add(dd).Format("20060102")
		cf.Pay_DATE_E = py
		logr.Infof("商户[%s]假期[%s-%s][%s]合并,划付日期为:%v", mc, cf.Pay_DATE_S, cf.STLM_DATE, th.HOLIDAY_DSP, cf.Pay_DATE_E)
	}

	//logr.Info("初始化商户号:", cf.MCHT_TP)
	for _, c := range strings.Split(tbrec.EXT4, ",") {
		cf.Ins_id_cd = append(cf.Ins_id_cd, c)

	}
	cf.MCHT_TP = tbrec.Mcht_ty
	logr.Infof("初始化商户号:[%+v][%s]", cf.Ins_id_cd, cf.MCHT_TP)
	return nil
}
