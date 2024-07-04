BUILDDIR ?= $(CURDIR)/build

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILDDIR)/linux-amd64/go-tha-utxos


build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILDDIR)/darwin-amd64/go-tha-utxos


build-win64:
	GOOS=windows GOARCH=amd64 go build -o $(BUILDDIR)/win64/go-tha-utxos.exe


build:
	go build -o $(BUILDDIR)/go-tha-utxos

release: build-linux build-darwin-amd64 build-win64
	rm -rf $(BUILDDIR)/compressed
	mkdir -p $(BUILDDIR)/compressed
	zip -j $(BUILDDIR)/compressed/win64.zip $(BUILDDIR)/win64/go-tha-utxos.exe
	tar -czvf $(BUILDDIR)/compressed/darwin-amd64.tar.gz -C $(BUILDDIR)/darwin-amd64/ .
	tar -czvf $(BUILDDIR)/compressed/linux-amd64.tar.gz -C $(BUILDDIR)/linux-amd64/ .
