
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
	packr2

build:
	go build -o readimension *.go

clean:
	rm -rf books covers uploads readimension.db
	
compile: bundle build

preset: preset-sass

preset-sass:
	mkdir -p .bin
	curl -fkLo .bin/sass.tar.gz https://github.com/sass/dart-sass/releases/download/1.29.0/dart-sass-1.29.0-${detected_sass_OS}-x64.tar.gz
	tar xvzf .bin/sass.tar.gz -C .bin
	chmod +x .bin/dart-sass/sass
	rm -rf .bin/sass.tar.gz