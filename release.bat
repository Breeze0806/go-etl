set GO111MODULE=on
cd %GOPATH%\src\github.com\Breeze0806\go-etl
go mod download
go mod vendor   
go run tools\datax\build\main.go -i %IGNORE_PACKAGES%
cd cmd\datax
go build 
cd ../..
go run tools/datax/release/main.go