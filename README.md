# sc_im
golang im server &amp; client (for study)

# GUIDE 
EDIT: go.mod require github.com/dappbujiujiu/sc_im v1.0.2

COMMAND: go mod download 

In Your Project: 
import scImServer "github.com/dappbujiujiu/sc_im/module"
server := scImServer.NewServer(Host, Port)  //Host string, Port int

Client:
go run github.com/dappbujiujiu/sc_im/client/c.go -host 127.0.0.1 -port 8888