APP_NAME=blog-api
RELEASE=v1
APP_PATH=release/${RELEASE}/

build:
	make clear
	@mkdir ${APP_PATH}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o ${APP_PATH}${APP_NAME} -installsuffix cgo . 
	ln -s ${APP_PATH} app

clear:
	go clean
	rm -rf ${APP_PATH}

run:
	clear
	go run .

docker:
	docker build --rm -f "Dockerfile"  -t ${APP_NAME}:latest .

docker-run:
	-docker stop ${APP_NAME}
	-docker rm ${APP_NAME}
	docker run -d --restart=always --name ${APP_NAME} --network app --ip 172.18.0.16 -v /home/so/Documents/conf/blog-api:/.blog ${APP_NAME}
