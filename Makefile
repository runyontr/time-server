

binary:
	CGO_ENABLED=0 GOOS=linux go build -o app *.go

docker: binary
	docker build -t runyonsolutions/time-server .
