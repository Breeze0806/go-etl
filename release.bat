set GO111MODULE=on
go mod download
go mod vendor   
go generate ./...
cd cmd\datax
go build -ldflags="-s -w"
cd ../..
go run tools/datax/release/main.go