build:
	go build -o bin/main main.go

run:
	go run main.go

deploy: build ./main

clean:
	rm -rf bin/main