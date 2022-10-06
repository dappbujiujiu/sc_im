package main

import(
	server "github.com/dappbujiujiu/sc_im/module"
)
const (
	HOST = "127.0.0.1"
	PORT = 8888
)

func main() {
	//因为server属于同包内，所以不用import
	Server := server.NewServer(HOST, PORT)
	Server.Start()
}