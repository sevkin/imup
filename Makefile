all: test

test: # mock # unit tests
	go test ./...

mock: # https://github.com/vektra/mockery
	[ -x `which mockery` ] && \
		mockery -all

# dev stuff below

LISTEN := "localhost:3000"
API := http://$(LISTEN)/api/v1
IMG := testdata/image.jpg

run:
	go run main.go -listen $(LISTEN)

form:
	curl -v -F "image=@$(IMG)" $(API)/upload/form

json:
	(echo -n '{"image": "'; base64 $(IMG); echo '"}') \
	| curl -v -H "Content-Type: application/json" -d @- $(API)/upload/json
