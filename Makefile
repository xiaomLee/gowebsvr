all: install

fmt:
	gofmt -l -w -s ./
govendor:
	govendor install
install: cleanvendor clean fmt govendor
	install -d output/conf/ output/logs/ output/bin
	GO15VENDOREXPERIMENT=1 go build -o output/bin/gowebsvr ./
	cp -r config/dev/* output/conf/
	cp control.sh output/

cleanvendor:
	rm -rf vendor/
clean:
	go clean -i ./
	rm -rf output
