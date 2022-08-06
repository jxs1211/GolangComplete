## main step:

0 go generate ./...
1 go build -buildmode=plugin -o ./plugin/plugin.so ./plugin/plugin.go
2 go build -o main ./main.go
3 ./main