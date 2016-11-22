package main

import (
	"bufio"
	"configuration"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"server"
	"strings"
	"text/template"
)

var staticPages = populateStatisPages()

func serveWeb() {
	gorillaRoute := mux.NewRouter()
	gorillaRoute.HandleFunc("/", serveContent)
	gorillaRoute.HandleFunc("/{page_alias}", serveContent)

	http.HandleFunc("/images/", serveResourse)
	http.HandleFunc("/stylesheets/", serveResourse)
	http.HandleFunc("/js/", serveResourse)
	http.HandleFunc("/fonts/", serveResourse)

	http.Handle("/", gorillaRoute)
	http.ListenAndServe(":8080", nil)
}

func serveContent(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	page_alias := urlParams["page_alias"]
	if page_alias == "" {
		page_alias = "index"
	}
	staticPage := staticPages.Lookup(page_alias + ".html")
	if staticPage == nil {
		staticPage = staticPages.Lookup("404.html")
		w.WriteHeader(404)
	}
	staticPage.Execute(w, nil)

}

func populateStatisPages() *template.Template {
	result := template.New("templates")
	templatePaths := new([]string)

	basePath := "app/pages"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()
	templatePathsRaw, _ := templateFolder.Readdir(-1)
	for _, pathInfo := range templatePathsRaw {
		log.Println(pathInfo.Name())
		*templatePaths = append(*templatePaths, basePath+"/"+pathInfo.Name())
	}
	result.ParseFiles(*templatePaths...)
	return result
}

func serveResourse(w http.ResponseWriter, r *http.Request) {
	path := "app/" + r.URL.Path
	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css; charset=utf-8"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png; charset=utf-8"
	} else if strings.HasSuffix(path, ".jpg") {
		contentType = "image/jpg; charset=utf-8"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "application/javascript; charset=utf-8"
	} else {
		contentType = "text/plain; charset=utf-8"
	}
	log.Println(path)
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		w.Header().Add("Content-Type", contentType)
		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}

func main() {
	switch configuration.Config.HttpsUsage {
	case "AdminOnly":
		//checkHttpsCertificates()
		httpRouter := httptreemux.New()
		httpsRouter := httptreemux.New()
		// Blog and pages as http
		server.InitializeBlog(httpRouter)
		server.InitializePages(httpRouter)
		// Blog and pages as https
		server.InitializeBlog(httpsRouter)
		server.InitializePages(httpsRouter)
		// Admin as https and http redirect
		// Add redirection to http router
		httpRouter.GET("/admin/", httpsRedirect)
		httpRouter.GET("/admin/*path", httpsRedirect)
		// Add routes to https router
		server.InitializeAdmin(httpsRouter)
		// Start https server
		log.Println("Starting https server on port " + httpsPort + "...")
		go func() {
			err := http.ListenAndServeTLS(httpsPort, filenames.HttpsCertFilename, filenames.HttpsKeyFilename, httpsRouter)
			if err != nil {
				log.Fatal("Error: Couldn't start the HTTPS server:", err)
			}
		}()
		// Start http server
		log.Println("Starting http server on port " + httpPort + "...")
		err := http.ListenAndServe(httpPort, httpRouter)
		if err != nil {
			log.Fatal("Error: Couldn't start the HTTP server:", err)
		}
	case "All":
		checkHttpsCertificates()
		httpsRouter := httptreemux.New()
		httpRouter := httptreemux.New()
		// Blog and pages as https
		server.InitializeBlog(httpsRouter)
		server.InitializePages(httpsRouter)
		// Admin as https
		server.InitializeAdmin(httpsRouter)
		// Add redirection to http router
		httpRouter.GET("/", httpsRedirect)
		httpRouter.GET("/*path", httpsRedirect)
		// Start https server
		log.Println("Starting https server on port " + httpsPort + "...")
		go func() {
			err := http.ListenAndServeTLS(httpsPort, filenames.HttpsCertFilename, filenames.HttpsKeyFilename, httpsRouter)
			if err != nil {
				log.Fatal("Error: Couldn't start the HTTPS server:", err)
			}
		}()
		// Start http server
		log.Println("Starting http server on port " + httpPort + "...")
		err := http.ListenAndServe(httpPort, httpRouter)
		if err != nil {
			log.Fatal("Error: Couldn't start the HTTP server:", err)
		}
	default: // This is configuration.HttpsUsage == "None"
		httpRouter := httptreemux.New()
		// Blog and pages as http
		server.InitializeBlog(httpRouter)
		server.InitializePages(httpRouter)
		// Admin as http
		server.InitializeAdmin(httpRouter)
		// Start http server
		log.Println("Starting server without HTTPS support. Please enable HTTPS in " + filenames.ConfigFilename + " to improve security.")
		log.Println("Starting http server on port " + httpPort + "...")
		err := http.ListenAndServe(httpPort, httpRouter)
		if err != nil {
			log.Fatal("Error: Couldn't start the HTTP server:", err)
		}
	}
	serveWeb()
}
