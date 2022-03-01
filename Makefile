
ifeq ($(shell uname), Darwin)
	detected_sass_OS := macos
else
	detected_sass_OS := linux
endif

.PHONY	:	run build bundle sass
all: sass bundle run
run:
	go run main.go app.go --env=development

sass:
	.bin/dart-sass/sass --no-source-map public/styles/style.scss public/styles/style.css

bundle:
	pkger -o route

build:
	go build -o readimension *.go

clean:
	rm -rf books covers uploads readimension.db
	
compile: bundle build

preset: preset-sass

preset-sass:
	go install github.com/markbates/pkger/cmd/pkger@latest
	mkdir -p .bin
	curl -fkLo .bin/sass.tar.gz https://github.com/sass/dart-sass/releases/download/1.49.9/dart-sass-1.49.9-${detected_sass_OS}-x64.tar.gz
	tar xvzf .bin/sass.tar.gz -C .bin
	chmod +x .bin/dart-sass/sass
	rm -rf .bin/sass.tar.gz