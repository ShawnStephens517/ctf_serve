# CTF Serve

ctf_serve is a lightweight command-line tool designed for educational and Capture The Flag (CTF) purposes. It allows you to quickly serve files from a specified directory over HTTP, making it ideal for sharing files in a local network or during CTF competitions. It is a simple replacement for python's `python -m http.server"` with a bit of html/css to be easy on the eyes. The tool can be run on Linux, MacOS, and Windows platforms.

## ⚠️ Disclaimer

This tool is intended for educational and Capture The Flag (CTF) purposes only. Unauthorized use on systems without explicit permission is prohibited. By using this tool, you acknowledge that you are responsible for any actions taken using ctf_serve. The author and contributors hold no liability for any misuse or damages that may arise.

## Features

- Serve files from a specified directory over HTTP.

- Cross-platform compatibility: Linux, MacOS (Intel, ARM, Universal), and Windows.

- Simple, secure file server for CTF competitions and educational purposes.

## Downloading a Release

You can download pre-built binaries for Linux, MacOS, and Windows from the Releases page on GitHub.

### Available Binaries:

```
Linux (ctf_serve_linux_amd64)

MacOS (Intel, ARM, Universal)

ctf_serve_macos_amd64

ctf_serve_macos_arm64

ctf_serve_macos_universal

Windows 64-bit (ctf_serve_windows_amd64.exe)
```

# Building the Application

If you'd prefer to build the application yourself, you can use the provided Makefile to create binaries for your target platform.

**Prerequisites:**

- Go 1.23 or later

- make (for using the Makefile)

## Build Commands

To build the binaries, simply run the following commands based on your desired platform:

### Build all binaries

```
make all
```

### Build Linux binary

```
make linux
```

### Build macOS binaries (Intel and ARM)

```
make macos
```

### Create macOS universal binary

```
make macos_universal
```

### Build Windows binary

```
make windows
```

### Clean up built binaries

```
make clean
```

## Building Manually

If you prefer to manually build the application without using make, you can use the following commands:

### Linux (64-bit)

```
GOOS=linux GOARCH=amd64 go build -o bin/ctf_serve_linux_amd64
```

### MacOS (Intel)

```
GOOS=darwin GOARCH=amd64 go build -o bin/ctf_serve_macos_amd64
```

### MacOS (ARM)

```
GOOS=darwin GOARCH=arm64 go build -o bin/ctf_serve_macos_arm64
```

### MacOS (Universal)

```
lipo -create -output bin/ctf_serve_macos_universal bin/ctf_serve_macos_amd64 bin/ctf_serve_macos_arm64
```

# Windows (64-bit)

```
GOOS=windows GOARCH=amd64 go build -o bin/ctf_serve_windows_amd64.exe
```

# Running the Application

Once you have built the application or downloaded a release, you can run it as follows:

`./ctf_serve_linux_amd64 --directory /path/to/serve --port 8080`

```
Command-Line Options

--directory or -d: Specify the directory to serve files from.

--ip or -i: Specify the IP address to bind the server to.

--port or -p: Specify the port number to use (default is 8080).
```

**Example:**

```
./ctf_serve_linux_amd64 -d /var/www/html -p 8080
```

Help

You can also display a help guide with detailed information on usage and flags:

```
./ctf_serve_linux_amd64 --help
```

# Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any bugs or features you'd like to add.

# License

This project is licensed under the MIT License. See the LICENSE file for details.

# Contact

For questions or feedback, please open an issue on GitHub or contact the repository maintainer
