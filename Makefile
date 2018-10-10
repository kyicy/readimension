.PHONY	:	run build bundle
all: bundle run
run:
	go run main.go app.go --env=development

bundle:
	packr

build:
	go build -o wankel *.go

compile: bundle build