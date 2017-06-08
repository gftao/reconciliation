package main

import (
	"golib/modules/gormdb"
	"fmt"
	"prodPmpCld/global"
	"golib/modules/config"
	"golib/modules/run"

	"Reconciliation/reconc"
	//"runtime"
	"os"
)

type TaskList struct {
	Name   string
	Action run.IRun
}

var g_taskList []TaskList = []TaskList{
	//{"ORDERCLI", &orderCli.OrderCliComm{}},
	{"CRTFILE", &reconc.CrtFile{}},
}

func main() {

	args := os.Args //获取用户输入的所有参数
	if args == nil || len(args) < 2{
		fmt.Println("请带一个日期参数")
		return
	}

	initParam := run.InitParams{}

	err := config.InitModuleByParams(global.CONFIGFILE)
	if err != nil {
		fmt.Println("读取配置文件失败", global.CONFIGFILE, err)
		return
	}
	err = gormdb.InitModule()
	if err != nil {
		fmt.Println("初始化数据库失败", err)
		return
	}

	for _, task := range g_taskList {
		ac := task.Action
		task.Name = args[1]  //清算日期
		//task.Name = "20161119"
		err = ac.Init(initParam, task.Name)
		if err != nil {
			fmt.Println("初始化失败: ", task.Name, err)
			return
		}
		fmt.Println(task.Name + "初始化成功")
	}

	for _, task := range g_taskList {
		ac := task.Action
		ac.Run()
	}

	fmt.Println("程序启动成功")

	fmt.Println("----main end!------")
	//runtime.Goexit()

}
