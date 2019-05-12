LISTEN := "localhost:3000"
API := http://$(LISTEN)/api/v1
IMG := testdata/image.jpg

all: run

run:
	go run main.go -listen $(LISTEN)

form:
	curl -v -F "image=@$(IMG)" $(API)/upload/form
