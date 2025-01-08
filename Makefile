build:
	go build -o ollama main.go

compile:
	echo "Compiling for linux and Windows"
	GOOS=linux GOARCH=386 go build -o bin/ollama-linux-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/ollama-windows-386 main.go