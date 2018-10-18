.PHONY	:	run build bundle sass
all: sass bundle run
run:
	go run main.go app.go --env=development

sass:
	sass --no-source-map public/styles/style.scss public/styles/style.css

bundle:
	packr

build:
	go build -o readimension *.go

compile: bundle build