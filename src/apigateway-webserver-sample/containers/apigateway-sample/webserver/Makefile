
build:
	docker build . -t sandbox.io/apigateway-sample/webserver:1.0.0

clean:
	docker rmi sandbox.io/apigateway-sample/webserver:1.0.0

run:
	docker run --name apigateway-sample -d -p 8080:8080 sandbox.io/apigateway-sample/webserver:1.0.0
