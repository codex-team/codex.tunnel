go build -o bin/client-macos .
env GOOS=linux GOARCH=386 go build -o bin/client-linux-i386 .
env GOOS=linux GOARCH=386 go build -o bin/client-linux-amd64 .
env GOOS=windows GOARCH=386 go build -o bin/client-windows-i386.exe .
env GOOS=windows GOARCH=amd64 go build -o bin/client-windows-amd64.exe .