go build -o client-macos .
env GOOS=linux GOARCH=386 go build -o client-linux-i386 .
env GOOS=linux GOARCH=386 go build -o client-linux-amd64 .
env GOOS=windows GOARCH=386 go build -o client-windows-i386.exe .
env GOOS=windows GOARCH=amd64 go build -o client-windows-amd64.exe .