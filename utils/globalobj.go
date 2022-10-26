package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

//存储一切有关zinx的全局参数 供其他模块使用
//一些参数可以通过配置文件由用户自定义

type GlobalObj struct {

	//Server
	TcpServer ziface.IServer //当前Zinx全局的Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器监听的端口
	Name      string         //当前服务器的名称

	//	Zinx
	Version        string //Zinx版本
	MaxConn        int    //最大连接数量
	MaxPackageSize uint32 //当前Zinx框架数据包的最大尺寸
}

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &g)
	if err != nil {
		panic(err)
	}
}

//定义一个全局的对外GlobalObj对象
var GlobalConfig *GlobalObj

//提供一个Init方法初始化当前的全局对象

func ConfigInit() {
	//如果配置文件没有加载就是默认值
	GlobalConfig = &GlobalObj{
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "ZinxServerApp",
		Version:        "V0.4",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}

	//应该尝试从Config/zinx中的用户自定义的json文件中读取
	GlobalConfig.Reload()
}
