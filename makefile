run:
	go run main.go
tidy:
	go mod tidy
ssh:
	ssh-keygen -f "/home/codespace/.ssh/known_hosts" -R "[localhost]:22222"
	ssh localhost -p 22222 < main.go