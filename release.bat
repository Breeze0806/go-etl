set GO111MODULE=on
git clone https://github.com/ibmdb/go_ibm_db %GOPATH%\src\github.com\ibmdb\go_ibm_db
cd %GOPATH%\src\github.com\ibmdb\go_ibm_db\installer  && go run setup.go
set path=%path%;%GOPATH%\src\github.com\ibmdb\clidriver\bin
cd %GOPATH%\src\github.com\Breeze0806\go-etl
go mod download
go mod vendor   
go generate ./...
cd cmd\datax
go build 
cd ../..