package units

//go:generate go test ./...

//go:generate go test -bench=./...

//go:generate go test -bench=./... -benchmem

//go:generate golangci-lint run ./...

//go:generate sh coverage.sh

//go:generate go test -bench=./... -benchmem -cpuprofile profile.out
