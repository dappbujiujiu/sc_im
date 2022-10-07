# sc_im
golang im server &amp; client (for study)

# GUIDE 
## EDIT: 
go.mod require github.com/dappbujiujiu/sc_im v1.0.5

## COMMAND: 
go mod download 

## In Your Server Project:
``` 
import scImServer "github.com/dappbujiujiu/sc_im/module"

server := scImServer.NewServer(Host, Port)
```

## Client Tool
cd \`go env GOMODCACHE\`/github.com/dappbujiujiu/sc_im@v1.0.5

sudo go build -o ./im_client client/c.go

./im_client -host 127.0.0.1 -port 8888