all: install

fmt:
	gofmt -l -w -s ./
glide_install:
	glide install
install: clean fmt glide_install
	install -d output/conf/ output/logs/ output/bin
	GO15VENDOREXPERIMENT=1 go build -o output/bin/gowebsvr
	cp -r conf/* output/conf/
	cp control.sh output/

dev: clean fmt glide_install
	install -d output/conf/ output/logs/ output/bin
	GO15VENDOREXPERIMENT=1 go build -o output/bin/gowebsvr
	cp -r conf/dev/* output/conf/
	cp control.sh output/

cleanvendor:
	rm -rf vendor/
clean:
	go clean -i ./
	rm -rf output
