package main

import (
	"flag"
	"fmt"
	"golib/modules/config"
	"golib/modules/gormdb"
	"golib/modules/logr"
	"golib/modules/run"
	"htdRec/reconc"
	"os"
	"prodPmpCld/global"
)

type TaskList struct {
	Name   string
	Action run.IRun
}

var g_taskList []TaskList = []TaskList{
	//{"CRTFILE", &reconc.FuzhouFJFile{}},
	{"YCRL", &reconc.GenFileTemp{}},
}

func main() {

	flag.Parse()
	args := os.Args //获取用户输入的所有参数
	//args = []string{"1", "20200106"}
	if args == nil || len(args) < 2 || len(args[1]) != 8 {
		fmt.Println(`请带一个格式为[20161119]的查询日期参数！`)
		return
	}

	initParam := run.InitParams{}

	err := config.InitModuleByParams(global.CONFIGFILE)
	if err != nil {
		fmt.Println("读取配置文件失败", global.CONFIGFILE, err)
		return
	}

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

	for _, task := range g_taskList {
		ac := task.Action
		task.Name = args[1] //清算日期
		//task.Name = "20190121"

		err = ac.Init(initParam, task.Name)
		if err != nil {
			logr.Info("初始化失败: ", task.Name, err)
			return
		}
	}

	for _, task := range g_taskList {
		ac := task.Action
		ac.Run()
	}

	logr.Info("程序启动成功")

	logr.Info("----main end!------")
	//runtime.Goexit()
}
