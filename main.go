package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var asciiArt = `

___  ____  ____      ____  ____  ____  _  _  ____ 
/ __)(_  _)(  __)    / ___)(  __)(  _ \/ )( \(  __)
( (__   )(   ) _)____ \___ \ ) _)  )   /\ \/ / ) _) 
\___) (__) (__)(____)(____/(____)(__\_) \__/ (____)

`

var rootCmd = &cobra.Command{
	Use:   "ctf_serve",
	Short: "A lightweight HTTP server for serving files",
	Long: `ctf_serve is a lightweight CLI tool to serve files from a directory over HTTP.

Environment Variables:
  - DIRECTORY: The directory to serve files from.
  - IP: The IP address to bind the server to.
  - PORT: The port number to use for the server.

Flags:
  -d, --directory: Specify the directory to serve files from (default is prompted).
  -i, --ip: Specify the IP address to bind to (default is prompted).
  -p, --port: Specify the port number to use (default is 8080).

Examples:
  ctf_serve --directory /var/www/html --port 8080
  ctf_serve -d /tmp`,

	Run: func(cmd *cobra.Command, args []string) {
		// Print colorful ASCII art in rainbow colors
		rainbowColors := []string{"\033[31m", "\033[33m", "\033[32m", "\033[36m", "\033[34m", "\033[35m"}
		for i, char := range asciiArt {
			color := rainbowColors[i%len(rainbowColors)]
			fmt.Printf("%s%c\033[0m", color, char)
		}
		fmt.Println()

		// Get configuration values from flags or environment variables
		directory := viper.GetString("directory")
		if directory == "" {
			directory = selectDirectory()
		}

		// Check if the directory exists
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			log.Fatal("Directory does not exist.")
		}

		ipAddr := viper.GetString("ip")
		if ipAddr == "0.0.0.0" || ipAddr == "" {
			ipAddr = selectNetworkInterface()
		}

		port := viper.GetInt("port")
		if port == 8080 { // Check if default port value is used to prompt for selection
			port = selectPort()
		}

		// Start the HTTP server
		addr := fmt.Sprintf("%s:%d", ipAddr, port)
		fmt.Println("------------------------------------------------------")

		message := fmt.Sprintf("\nServing %s on http://%s\n", directory, addr)
		greenColor := "\033[32m"
		for _, char := range message {
			fmt.Printf("%s%c\033[0m", greenColor, char)
		}
		/*
			for i, char := range message {
				color := rainbowColors[i%len(rainbowColors)]
				fmt.Printf("%s%c\033[0m", color, char)
			}
		*/
		fmt.Println()

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			serveFile(w, r, directory)
		})

		server := &http.Server{Addr: addr}
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Server error: %v", err)
			}
		}()

		// Wait for a stop signal
		fmt.Println("\nPress ENTER to stop the server...")
		bufio.NewReader(os.Stdin).ReadString('\n')

		// Shutdown the server gracefully
		if err := server.Shutdown(nil); err != nil {
			log.Fatalf("Server Shutdown Failed: %v", err)
		}
		fmt.Println("Server stopped gracefully.")
	},
}

func init() {
	rootCmd.Flags().BoolP("help", "h", false, "Display help information")
	// Define flags for the CLI
	rootCmd.Flags().StringP("directory", "d", "", "Directory to serve files from")
	rootCmd.Flags().StringP("ip", "i", "", "IP address to bind to")
	rootCmd.Flags().IntP("port", "p", 8080, "Port number to use")

	// Bind flags to Viper
	viper.BindPFlag("directory", rootCmd.Flags().Lookup("directory"))
	viper.BindPFlag("ip", rootCmd.Flags().Lookup("ip"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
}

func main() {
	http.Handle("/icons/", http.StripPrefix("/icons/", http.FileServer(http.Dir("./icons"))))
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}

func selectDirectory() string {
	Home, _ := os.UserHomeDir()
	fmt.Println("Home Directory:", Home)

	fmt.Println("Select a directory to serve files from:")
	fmt.Println("1) /var/www/html")
	fmt.Println("2) /home/" + Home)
	fmt.Println("3) /tmp")
	fmt.Println("4) Custom directory")
	fmt.Print("Enter your choice [1-4]: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		return "/var/www/html"
	case "2":
		return filepath.Join(Home)
	case "3":
		return "/tmp"
	case "4":
		return promptUser("Please enter a custom directory: ")
	default:
		fmt.Println("Invalid choice. Defaulting to /tmp.")
		return "/tmp"
	}
}

func selectNetworkInterface() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal("Failed to get network interfaces: ", err)
	}

	var ifaceList []net.Interface
	var ipList []net.IP
	fmt.Println("------------------------------------------------------")
	fmt.Println("Select a network interface to serve from:")
	idx := 1
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.To4() != nil {
				fmt.Printf("%d) %s (%s)\n", idx, iface.Name, ipNet.IP.String())
				ifaceList = append(ifaceList, iface)
				ipList = append(ipList, ipNet.IP)
				idx++
			}
		}
	}

	if len(ifaceList) == 0 {
		log.Fatal("No suitable network interfaces found.")
	}

	fmt.Printf("Enter your choice [1-%d]: ", len(ifaceList))
	reader := bufio.NewReader(os.Stdin)
	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)
	choice, err := strconv.Atoi(choiceStr)
	if err != nil || choice < 1 || choice > len(ifaceList) {
		log.Fatal("Invalid choice.")
	}

	return ipList[choice-1].String()
}

func selectPort() int {
	fmt.Println("------------------------------------------------------")
	fmt.Println("Select a port to serve from:")
	fmt.Println("1) 80 (HTTP)")
	fmt.Println("2) 443 (HTTPS)")
	fmt.Println("3) 8080 (Alternative HTTP)")
	fmt.Println("4) Custom port")
	fmt.Print("Enter your choice [1-4]: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		return 80
	case "2":
		return 443
	case "3":
		return 8080
	case "4":
		portStr := promptUser("Please enter a custom port: ")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			log.Fatal("Invalid port number.")
		}
		return port
	default:
		fmt.Println("Invalid choice. Defaulting to port 8080.")
		return 8080
	}
}

func promptUser(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func serveFile(w http.ResponseWriter, r *http.Request, directory string) {
	// Construct the requested path
	requestedPath := r.URL.Path
	decodedPath, err := filepath.Abs(filepath.Join(directory, filepath.FromSlash(strings.TrimPrefix(requestedPath, "/"))))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Prevent path traversal attacks by ensuring the path stays within the served directory
	if !strings.HasPrefix(decodedPath, directory) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Check if the path exists
	fileInfo, err := os.Stat(decodedPath)
	if os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If the path is a directory, render the directory listing
	if fileInfo.IsDir() {
		files, err := os.ReadDir(decodedPath)
		if err != nil {
			http.Error(w, "Unable to read directory", http.StatusInternalServerError)
			return
		}

		// Load the index.html template from the templates directory
		templatePath := filepath.Join("templates", "index.html")
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			http.Error(w, "Unable to load template", http.StatusInternalServerError)
			return
		}

		// Create a list of file entries for use in the template
		var fileEntries []map[string]string
		imageExtensions := []string{".jpg", ".png", ".gif", ".svg", ".webp", ".jpeg"}
		excelExtensions := []string{".xlsx", ".xls", ".csv"}
		docExtensions := []string{".doc", ".docx"}
		archiveExtensions := []string{".tar", ".tgz", ".zip", ".7zip", ".pkzip", ".gzip"}
		virtualizationExtensions := []string{".vmdk", ".ova", ".ovf", ".qcow2", ".qcow"}

		for _, file := range files {
			entry := map[string]string{
				"Name": file.Name(),
				"URL":  filepath.Join(r.URL.Path, file.Name()),
			}
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if file.IsDir() {
				entry["Icon"] = "folder-open-regular.svg"
			} else if strings.Contains(file.Name(), "exploit") {
				entry["Icon"] = "exploit_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".json") {
				entry["Icon"] = "json_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".txt") {
				entry["Icon"] = "file-lines-regular.svg"
			} else if strings.HasSuffix(file.Name(), ".go") {
				entry["Icon"] = "golang-brands-solid.svg"
			} else if strings.HasSuffix(file.Name(), ".py") {
				entry["Icon"] = "python_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".rs") {
				entry["Icon"] = "rust-brands-solid.svg"
			} else if strings.HasSuffix(file.Name(), ".sh") {
				entry["Icon"] = "sh_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".jar") {
				entry["Icon"] = "java_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".yaml") {
				entry["Icon"] = "filetype_yml_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".yml") {
				entry["Icon"] = "filetype_yml_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".js") {
				entry["Icon"] = "javascript_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".ts") {
				entry["Icon"] = "javascript_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".cs") {
				entry["Icon"] = "c_sharp_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".ppt") {
				entry["Icon"] = "office365_powerpoint_icon.svg"
			} else if slices.Contains(docExtensions, ext) {
				entry["Icon"] = "office365_word_icon.svg"
			} else if slices.Contains(excelExtensions, ext) {
				entry["Icon"] = "office365_excel_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".pdf") {
				entry["Icon"] = "pdf_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".iso") {
				entry["Icon"] = "iso_icon.png"
			} else if slices.Contains(archiveExtensions, ext) {
				entry["Icon"] = "archive_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".exe") {
				entry["Icon"] = "exe_icon.svg"
			} else if strings.HasSuffix(file.Name(), ".dmg") {
				entry["Icon"] = "dmg_icon.svg"
			} else if slices.Contains(virtualizationExtensions, ext) {
				entry["Icon"] = "vm_icon.svg"
			} else {
				if slices.Contains(imageExtensions, ext) {
					entry["Icon"] = "image_icon.svg"
				} else if strings.HasSuffix(file.Name(), ".HEIC") {
					entry["Icon"] = "image_icon.svg"
				} else {
					entry["Icon"] = "file_empty_icon.svg"
				}
			}
			fileEntries = append(fileEntries, entry)
		}

		// Execute the template with the file list data
		data := map[string]interface{}{
			"Directory": decodedPath,
			"FileList":  fileEntries,
		}
		w.Header().Set("Content-Type", "text/html")
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Unable to render template", http.StatusInternalServerError)
		}
		return
	}

	// If the path is a file, serve it
	contentType := "application/octet-stream"
	if ext := filepath.Ext(decodedPath); ext != "" {
		contentType = mimeTypeByExtension(ext)
	}

	w.Header().Set("Content-Type", contentType)
	http.ServeFile(w, r, decodedPath)
}

func mimeTypeByExtension(ext string) string {
	// Basic MIME type mapping
	switch ext {
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".json":
		return "application/json"
	default:
		return "application/octet-stream"
	}
}
