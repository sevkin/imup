all: test

test: # mock # unit tests
	go test ./...

mock: # https://github.com/vektra/mockery
	[ -x `which mockery` ] && \
		mockery -all

build:
	docker-compose build
up:
	docker-compose up
down:
	docker-compose down

clean:
	rm -rf tmp

# dev stuff below

LISTEN := "localhost:3000"
API := http://$(LISTEN)/api/v1
IMG := testdata/image.jpg

run:
	mkdir -p tmp
	go run main.go -listen $(LISTEN) -storage ./tmp -thumbcmd ./thumb100.sh

form:
	curl -v -F "image=@$(IMG)" $(API)/upload/form

json:
	(echo -n '{"image": "'; base64 $(IMG); echo '"}') \
	| curl -v -H "Content-Type: application/json" -d @- $(API)/upload/json

serve:
	python -m SimpleHTTPServer 5000
url:
	curl -v "$(API)/upload/url?image=http://localhost:5000/$(IMG)"