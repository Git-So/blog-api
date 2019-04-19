APP_NAME=app
RELEASE=v1
APP_PATH=release/${RELEASE}/

build:
	make clear
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o ${APP_PATH}${APP_NAME}
	ln -s ${APP_PATH} app

clear:
	rm -rf ${APP_PATH}

run:
	clear
	go run .