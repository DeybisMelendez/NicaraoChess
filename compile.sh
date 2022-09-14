GOOS=windows GOARCH=amd64 go build -o bin/nicarao-amd64-windows.exe .
GOOS=darwin GOARCH=amd64 go build -o bin/nicarao-amd64-macos .
GOOS=linux GOARCH=amd64 go build -o bin/nicarao-amd64-linux .