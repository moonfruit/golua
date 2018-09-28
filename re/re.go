package re

//go:generate go-bindata -pkg re -prefix ../lpeg/ ../lpeg/re.lua
//go:generate gofmt -w -s bindata.go
