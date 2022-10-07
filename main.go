package main

import(
	server "github.com/dappbujiujiu/sc_im/module"
	lib "github.com/dappbujiujiu/sc_im/lib"
)

func main() {
	conf := lib.GetConf()	//获取server配置
	//因为server属于同包内，所以不用import
	Server := server.NewServer(conf.Server.Host, conf.Server.Port)
	Server.Start()
}