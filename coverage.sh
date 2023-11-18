mkdir -p ./coverage
go test ./... -coverprofile=./coverage/coverage.out
go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html
rm ./coverage/coverage.out
