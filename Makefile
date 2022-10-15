.PHONY: all build install clean

all: build install

build:
        @echo "Building..."
        go build -ldflags="-extldflags=-static" -o ./bin/nabaz ./cmd/nabaz

install: build
        @echo "Installing..." 
        cp -f bin/* /usr/local/bin/
        

clean:
        @echo "Cleaning up..."
        rm *.txt
