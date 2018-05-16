package reconc

import (
	"golib/gerror"
	"golib/modules/run"
	"golib/modules/logr"
	"golib/modules/gormdb"
	"golib/modules/config"
)

type IMYRun interface {
	Init(string, string) gerror.IError
	Run()
}

type GenFile struct {
	STLM_DATE  string //清算日期
	MCHT_CD    string
	FileName   string            //对账文件名
	MCHT_CDS   []string          //全部机构号
	MCHT_RECTY map[string]string //1-集团商户
	Action     IMYRun
}

func (g *GenFile) Init(initParams run.InitParams, chainName string) gerror.IError {
	g.STLM_DATE = chainName //清算日期
	g.MCHT_RECTY = make(map[string]string, 1)
	err := g.InitMCHTCd()
	if err != nil {
		return gerror.NewR(1001, err, "无记录")
	}

	return nil
}

func (g *GenFile) Run() {
	for {
		cd, ok := g.GetMCHTCd()
		if !ok {
			return
		}
		logr.Infof("根据商户号注册方法:%s", g.MCHT_RECTY[cd])
		switch g.MCHT_RECTY[cd] {
		case "1"://银川热力专用
			g.Action = &CrtFunc1{}
		default:
			g.Action = &CrtFile{}
		}

		g.Action.Init(g.STLM_DATE, cd)
		g.Action.Run()
	}

}

func (g *GenFile) Finish() {

}

func (g *GenFile) InitMCHTCd() error {

	//商户号
	mc, ok := config.String("MCHT_CD")
	if ok && mc != "" {
		g.MCHT_CDS = append(g.MCHT_CDS, mc)
		mt := config.StringDefault("MCHT_RECTY", "0")
		logr.Infof("读配置文件%s,%s", mc, mt)
		g.MCHT_RECTY[mc] = mt
		return nil
	}

	dbc := gormdb.GetInstance()
	rows, err := dbc.Raw("SELECT distinct MCHT_CD, EXT1 FROM tbl_mcht_recon_list").Rows()
	if err != nil {
		logr.Info("dbc.Raw fail:%s\n", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		mc, mt := "", ""
		rows.Scan(&mc, &mt)
		if mc != "" {
			g.MCHT_CDS = append(g.MCHT_CDS, mc)
			g.MCHT_RECTY[mc] = mt
		}
	}

	logr.Infof("初始化商户号:+v", g.MCHT_RECTY)
	return nil
}

func (g *GenFile) GetMCHTCd() (string, bool) {
	l := len(g.MCHT_CDS)
	if l == 0 {
		return "", false
	}

	g.MCHT_CD = g.MCHT_CDS[0]
	g.MCHT_CDS = g.MCHT_CDS [1:]
	logr.Infof("取机构号：%s; 剩余机构号:%v", g.MCHT_CD, g.MCHT_CDS)

	return g.MCHT_CD, true
}
