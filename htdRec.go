package main

import (
	"golib/modules/config"
	"golib/modules/gormdb"
	"golib/modules/run"
	"prodPmpCld/global"
	"htdRec/reconc"
	"golib/modules/logr"
	"os"
	"fmt"
	"flag"
	"runtime/pprof"
	"runtime"
)

type TaskList struct {
	Name   string
	Action run.IRun
}

var g_taskList []TaskList = []TaskList{

	{"CRTFILE", &reconc.CrtFile{}},
}

var cpuprofile = flag.String("C", "cpu.prof", "write cpu profile `file`")
var memprofile = flag.String("M", "mem.prof", "write memory profile to `file")

func main() {

	flag.Parse()

	args := os.Args //获取用户输入的所有参数
	if args == nil || len(args) < 2 || len(args[1]) != 8 {
		fmt.Println(`请带一个格式为: [20161119]的查询日期参数！`)
		return
	}

	initParam := run.InitParams{}

	err := config.InitModuleByParams(global.CONFIGFILE)
	if err != nil {
		fmt.Println("读取配置文件失败", global.CONFIGFILE, err)
		//logr.Info("读取配置文件失败", global.CONFIGFILE, err)
		return
	}

	//logr.Info("开始初始化日志")
	err = logr.InitModules()
	if err != nil {
		logr.Info("初始化日志失败", err)
		return
	}
	logr.Info("初始化日志成功")

	err = gormdb.InitModule()
	if err != nil {
		logr.Info("初始化数据库失败", err)
		return
	}
	pp := config.BoolDefault("prof", false)
	if pp {
		if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
				logr.Error("could not create CPU profile: ", err)
			}
			if err = pprof.StartCPUProfile(f); err != nil {
				logr.Error("could not Start CPU profile: ", err)
			}
			defer pprof.StopCPUProfile()
		}
		if *memprofile != "" {
			f, err := os.Create(*memprofile)
			if err != nil {
				logr.Error("could not create memory profile: ", err)
			}
			runtime.GC()
			if err = pprof.WriteHeapProfile(f); err != nil {
				logr.Error("could not  write memory profile: ", err)
			}

			f.Close()
		}
	}


	for _, task := range g_taskList {
		ac := task.Action
		task.Name = args[1] //清算日期
		//task.Name = "20160815"
		err = ac.Init(initParam, task.Name)
		if err != nil {
			logr.Info("初始化失败: ", task.Name, err)
			return
		}
		//fmt.Println(task.Name + "初始化成功")
	}

	for _, task := range g_taskList {
		ac := task.Action
		ac.Run()
	}

	logr.Info("程序启动成功")

	logr.Info("----main end!------")
	//runtime.Goexit()

}
