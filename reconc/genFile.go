package reconc

import (
	"golib/gerror"
	"golib/modules/config"
	"golib/modules/gormdb"
	"golib/modules/logr"
	"golib/modules/run"
	"strings"
	"sync"
	"time"
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
	MCHT_TIMER map[string]string //定时器
	Action     IMYRun
}

func (g *GenFile) Init(initParams run.InitParams, chainName string) gerror.IError {
	g.STLM_DATE = chainName //清算日期
	g.MCHT_RECTY = make(map[string]string)
	g.MCHT_TIMER = make(map[string]string)
	err := g.InitMCHTCd()
	if err != nil {
		return gerror.NewR(1001, err, "无记录")
	}

	return nil
}

func (g *GenFile) Run() {
	var wg sync.WaitGroup

	for {
		cd, ok := g.GetMCHTCd()
		if !ok {
			break
		}
		logr.Infof("根据商户号注册方法:%v", g.MCHT_RECTY[cd])
		switch g.MCHT_RECTY[cd] {
		case "0":
			g.Action = &CrtFile2{} //default划付日期版
		case "1": //银川热力专用
			g.Action = &CrtFunc1{}
		case "2": //【生态圈】POS对账文件说明
			g.Action = &Ecosph{}
		case "3": //中法供水
			g.Action = &SinoFrench{}
		case "4": //福州房管局
			g.Action = &FuzhouFJFile{}
		case "5": //昆山住建局
			g.Action = &KunshanZJFile{}
		case "6": //湖州住建局
			g.Action = &HZFile{}
		case "7": //柳州住建局
			g.Action = &LiuzhouZJFile{}
		case "8": //南通住建局
			g.Action = &NantongZJFile{}
		case "9": //昆明住建局
			g.Action = &KunmingZJFile{}
		case "10": //衢州住建局
			g.Action = &QuzhouZJFile{}
			/*		case "11": //成都住建局
					g.Action = &ChengduZJFile{}*/
		case "12": //石家庄住建局
			g.Action = &ShijiazhuangZJFile{}
		case "13": //贵阳住建局
			g.Action = &GYFile{}
			/*		case "14": //上海科技馆
					g.Action = &SHKJGFile{}*/
		case "15": //济南住建局
			g.Action = &JinanZJFile{}
		case "16": //贵阳二手房
			g.Action = &GYESFile{}
		default:
			g.Action = &CrtFile{}
		}

		if md, ok := g.MCHT_TIMER[cd]; ok {
			logr.Infof("根据商户号等待时间:%v", g.MCHT_TIMER[cd])
			wg.Add(1)
			d, _ := time.ParseDuration(md)
			go func(a IMYRun, duration time.Duration) {
				defer wg.Done()
				//fmt.Println("duration=", duration)
				time.Sleep(duration)
				gerr := a.Init(g.STLM_DATE, cd)
				if gerr == nil {
					a.Run()
				} else {
					logr.Error(gerr)
				}

			}(g.Action, d)
		} else {
			gerr := g.Action.Init(g.STLM_DATE, cd)
			if gerr == nil {
				g.Action.Run()
			} else {
				logr.Error(gerr)
			}
		}

	}

	wg.Wait()
	//os.Exit(0)
}

func (g *GenFile) Finish() {

}

func (g *GenFile) InitMCHTCd() error {
	mc, ok := config.String("MCHT_CD")
	if ok && mc != "" {
		g.MCHT_CDS = append(g.MCHT_CDS, mc)
		mt := config.StringDefault("MCHT_RECTY", "0")
		logr.Infof("读配置文件:%s,%s", mc, mt)
		g.MCHT_RECTY[mc] = mt
		tm := config.StringDefault("MCHT_TIMER", "0")
		g.MCHT_TIMER[mc] = tm
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
			mts := strings.Split(mt, "|")
			//fmt.Println(mts)
			for i, _ := range mts {
				switch i {
				case 0:
					g.MCHT_RECTY[mc] = mts[i]
				case 1:
					g.MCHT_TIMER[mc] = mts[i]
				}
			}
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
	g.MCHT_CDS = g.MCHT_CDS[1:]
	logr.Infof("取机构号:%s; 剩余机构号:%v", g.MCHT_CD, g.MCHT_CDS)

	return g.MCHT_CD, true
}
