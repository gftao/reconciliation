package reconc

import (
	"golib/gerror"
	"golib/modules/run"
	"golib/modules/gormdb"
	"prodPmpCld/utils/sutil"
	"fmt"

	"github.com/jinzhu/gorm"
	"os"
	"Reconciliation/models"
	"golib/modules/config"
	"bufio"
	"strconv"

	"strings"
)

type CrtFile struct {
	FileName             string                     //对账文件名
	FilePath             string                     //路径
	SysDate              string                     //当前时间
	InsIdCd              string                     //当前机构号
	STLM_DATE            string                     //清算日期
	Ins_id_cd            []string                   //全部机构号
	FileStrt             models.FileStrt            //文件结构
	Tbl_Clear_Data       []models.Tbl_clear_txn     //当前机构号数据
	Tbl_Tfr_his_log_Data models.Tbl_tfr_his_trn_log //当前机构号对应交易日志
}

type InsIdCd struct {
	Ins_id_cd string
}

func (cf *CrtFile)Init(initParams run.InitParams, chainName string) gerror.IError {
	cf.FileName = "C_"
	cf.FilePath = config.StringDefault("filePath", "./")
	cf.SysDate = sutil.GetSysDate()//今天
	cf.FileStrt = models.FileStrt{}
	cf.STLM_DATE = chainName        //清算日期

	cf.FileStrt.Init()
	cf.InitInsIdCd()//查表取机构号

	return nil
}

func (cf *CrtFile) Run() {
	for {
		if len(cf.Ins_id_cd) == 0 {
			return
		}
		fp := cf.geneFile()
		fmt.Println(fp)
		f, err := os.Create(fp)

		if err != nil {
			fmt.Println(cf.FileName, err)
			return
		}

		//读数据
		cf.ReadDate()

		w := bufio.NewWriter(f)  //创建新的 Writer 对象
		//fmt.Printf("头标识：%s\n",cf.FileStrt.FileHead)
		w.WriteString(cf.FileStrt.FileHead)        //文件头 标识
		w.WriteString("\r\n")
		w.WriteString(cf.FileStrt.HToString())
		w.WriteString("\r\n")

		w.WriteString(cf.FileStrt.FileBody)//文件体 标识
		w.WriteString("\r\n")

		for _, tc := range cf.FileStrt.FileBodys {
			tcStr := tc.BToString()
			//if tc.INS_ID_CD == "62510000" {
			//	fmt.Printf("%v\n",tc)
			//}
			w.WriteString(tcStr)
			w.WriteString("\r\n")
			//w.Flush()
		}

		w.Flush()
		f.Sync()
		f.Close()
	}

	return
}
func (cf *CrtFile) Finish() {
	return
}
func (cf *CrtFile) StoreReconc() {

}

func (cf *CrtFile) ReadDate() {
	dbc := gormdb.GetInstance()

	rows, err := dbc.Raw("SELECT * FROM tbl_clear_txn WHERE INS_ID_CD = ? and STLM_DATE = ?", cf.InsIdCd, cf.STLM_DATE).Rows() // (*sql.Rows, error)
	defer rows.Close()
	if err == gorm.ErrRecordNotFound {
		fmt.Printf("dbc.Raw fail:", err)
		return
	}
	if err != nil {
		fmt.Printf("dbc.Raw fail:", err)
		return
	}
	cf.Tbl_Clear_Data = make([]models.Tbl_clear_txn, 0)
	record := 0        //交易总笔数
	trans_amt_T := 0.0   //清算金额
	true_fee_mod_T := 0.0  //清算手续费
	trnrecont_T := 0.0  //结算总金额
	for rows.Next() {
		record ++
		tc := models.Tbl_clear_txn{}
		dbc.ScanRows(rows, &tc)

		a, _ := strconv.ParseFloat(tc.TRANS_AMT, 64)
		f, _ := strconv.ParseFloat(tc.TRUE_FEE_MOD, 64)
		m, _ := strconv.ParseFloat(tc.MCHT_SET_AMT, 64)
		fmt.Println(a, f, m)
		trans_amt_T += a
		true_fee_mod_T += f
		trnrecont_T += m
		cf.Tbl_Clear_Data = append(cf.Tbl_Clear_Data, tc)
		fmt.Println(a, f, m)
	}
	fmt.Println(trans_amt_T, true_fee_mod_T, trnrecont_T)
	//处理文件头
	cf.FileStrt.FileHeadInfo.INS_ID_CD = cf.InsIdCd
	cf.FileStrt.FileHeadInfo.TrnSucCount = strconv.Itoa(record)
	cf.FileStrt.FileHeadInfo.Stlm_date = cf.STLM_DATE
	cf.FileStrt.FileHeadInfo.TrnSucAm = strconv.FormatFloat(trans_amt_T,'f', 2, 64)
	cf.FileStrt.FileHeadInfo.TrnFeeT = strconv.FormatFloat(true_fee_mod_T,'f', 2, 64)
	cf.FileStrt.FileHeadInfo.TrnReconT = strconv.FormatFloat(trnrecont_T,'f', 2, 64)
	fmt.Printf("成功总笔数:%d\n", record)

	cf.saveDatatoFStru()


}

func (cf *CrtFile) saveDatatoFStru() {
	cf.FileStrt.FileBodys = make([]models.Body,0)
	for _, tc := range cf.Tbl_Clear_Data {
		b := models.Body{}
		tl := models.Tbl_tfr_his_trn_log{}
		dbc := gormdb.GetInstance()
		err := dbc.Where("KEY_RSP = ?", tc.KEY_RSP).Find(&tl).Error
		if err != nil {
			fmt.Printf("saveDatatoFStru db find failed:",err)
			//return
		}

		b.MCHT_CD  	= tc.MCHT_CD
		b.TRANS_DATE    = tl.TRANS_DT
		b.TRANS_TIME    = tl.TRANS_MT
		b.STLM_DATE     = cf.STLM_DATE
		b.TERM_ID       = tc.TERM_ID
		b.TRANS_KIND    = tc.TRANS_KIND
		b.KEY_RSP       = tc.KEY_RSP
		b.PAN           = tc.PAN
		b.CARD_KIND_DIS = tc.CARD_KIND_DIS
		b.TRANS_AMT     = tc.TRANS_AMT
		b.TRUE_FEE_MOD  = tc.TRUE_FEE_MOD
		b.MCHT_SET_AMT  = tc.MCHT_SET_AMT
		b.ERR_FEE_IN    = "0"
		b.ERR_FEE_OUT	= "0"
		if strings.EqualFold(tl.PROD_CD, "1151"){
			b.SYS_ID        = tl.INDUSTRY_ADDN_INF
		} else {
			b.SYS_ID        = tl.INDUSTRY_ADDN_INF
		}

		b.SYS_ID        = tl.RETRI_REF_NO
		b.INS_IN        = "0"
		b.INS_REAL_IN   = "0"
		b.INS_OUT       = "0"
		b.PROXY_CD      = "0"

		cf.FileStrt.FileBodys = append(cf.FileStrt.FileBodys,b)
	}
}

func (cf *CrtFile) GetInsIdCd() (string, bool) {
	l := len(cf.Ins_id_cd)
	if l == 0 {
		return "", false
	}

	cf.InsIdCd = cf.Ins_id_cd[0]
	cf.Ins_id_cd = cf.Ins_id_cd[1:]
	fmt.Printf("取机构号：%s; 剩余机构号:%v\n", cf.InsIdCd, cf.Ins_id_cd)

	return cf.InsIdCd, true
}

func (cf *CrtFile) geneFile() string {
	inscd, ok := cf.GetInsIdCd()
	if ok {
		cf.FileName = "C_"
		cf.FileName = cf.FileName + inscd
	} else {
		return ""
	}

	cf.FileName = cf.FileName + "_" + cf.STLM_DATE + ".txt"
	fmt.Printf("生成对账文件名称：%s\n", cf.FileName)

	return cf.FilePath + cf.FileName
}

func (cf *CrtFile) InitInsIdCd() {

	dbc := gormdb.GetInstance()
	//rows, err := dbc.DB().Query("select nextval('INS_ID_CD')from " + tl.TableName())
	//rows, err := dbc.Select("distinct INS_ID_CD").Find(&tl).Rows()
	//rows, err := dbc.Table(tl.TableName()).Select("INS_ID_CD").Rows()
	//dbc.Find(&tl, "INS_ID_CD = ?", "62510000")
	/*
	dbc = gormdb.GetInstance()
	tc := models.Tbl_clear_txn{}
	//err := dbc.Where(" key_rsp = ? ", "20160815155530959915").Find(&tc).Error
	err := dbc.First(&tc).Error

	if err == gorm.ErrRecordNotFound {
		fmt.Printf("dbc.where fail:",err)
		return
	}
	if err != nil{
		fmt.Printf("dbc.where fail:",err)
		return
	}
	fmt.Printf("tl:%v\n",tc)
	*/

	rows, err := dbc.Raw("SELECT distinct INS_ID_CD FROM tbl_ins_reconciliation").Rows() // (*sql.Rows, error)
	defer rows.Close()
	if err == gorm.ErrRecordNotFound {
		fmt.Printf("dbc.Raw fail:", err)
		return
	}
	if err != nil {
		fmt.Printf("dbc.Raw fail:", err)
		return
	}

	for rows.Next() {
		in := ""
		rows.Scan(&in)
		if in != "" {
			cf.Ins_id_cd = append(cf.Ins_id_cd, in)
		}
	}

	fmt.Printf("初始化机构号:%v\n", cf.Ins_id_cd)

}