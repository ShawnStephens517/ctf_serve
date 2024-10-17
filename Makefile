# Makefile for building ctf_serve for multiple platforms

.PHONY: all linux macos macos_universal windows clean

all: linux macos macos_universal windows

linux:
	@echo "Building Linux binary..."
	GOOS=linux GOARCH=amd64 go build -o bin/ctf_serve_linux_amd64

macos:
	@echo "Building macOS binaries (Intel and ARM)..."
	GOOS=darwin GOARCH=amd64 go build -o bin/ctf_serve_macos_amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/ctf_serve_macos_arm64

macos_universal: macos
	@echo "Creating macOS universal binary..."
	lipo -create -output bin/ctf_serve_macos_universal bin/ctf_serve_macos_amd64 bin/ctf_serve_macos_arm64

windows:
	@echo "Building Windows binary..."
	GOOS=windows GOARCH=amd64 go build -o bin/ctf_serve_windows_amd64.exe

clean:
	@echo "Cleaning up..."
	rm -rf bin

help:
	@echo "Makefile for building ctf_serve for multiple platforms."
	@echo ""
	@echo "Usage:"
	@echo "  make linux            Build the Linux binary"
	@echo "  make macos            Build the macOS binaries (Intel and ARM)"
	@echo "  make macos_universal  Create a macOS universal binary"
	@echo "  make windows          Build the Windows binary"
	@echo "  make all              Build all binaries (Linux, macOS, Windows)"
	@echo "  make clean            Clean up built binaries"