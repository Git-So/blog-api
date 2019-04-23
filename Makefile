APP_NAME=blog-api
RELEASE=v1
APP_PATH=release/${RELEASE}/

build:
	make clear
	@mkdir ${APP_PATH}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o ${APP_PATH}${APP_NAME} -installsuffix cgo . 
	ln -s ${APP_PATH} app

clear:
	rm -rf ${APP_PATH}

run:
	clear
	go run .

docker:
	make build
	docker build --rm -f "Dockerfile"  -t ${APP_NAME}:latest .
	docker stop blog-api
	docker rm blog-api
	docker run -d --restart=always --name blog-api --network app --ip 172.18.0.16 -v /home/so/Documents/conf/blog-api:/.blog blog-api
	