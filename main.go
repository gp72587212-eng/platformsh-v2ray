package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

var (
	remoteHost       string
	remotePort       string
	remoteHostHeader string
)

func init() {
	remoteHost = getEnv("REMOTE_HOST", "164.92.226.103")      // ⚠️ À remplacer par l'IP de votre VPS
	remotePort = getEnv("REMOTE_PORT", "80")                // ⚠️ Port du serveur Xray
	remoteHostHeader = getEnv("REMOTE_HOST_HEADER", "") // ⚠️ Host attendu par Xray

	flag.StringVar(&remoteHost, "host", remoteHost, "IP du serveur Xray")
	flag.StringVar(&remotePort, "port", remotePort, "Port du serveur Xray")
	flag.StringVar(&remoteHostHeader, "header", remoteHostHeader, "Valeur de l'en-tête Host")
	flag.Parse()
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	target := fmt.Sprintf("http://%s:%s", remoteHost, remotePort)
	remoteURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("URL cible invalide : %v", err)
	}

	// Transport HTTP avec des paramètres optimisés
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// Création du reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(remoteURL)
	proxy.Transport = transport
	proxy.FlushInterval = 100 * time.Millisecond // Permet le streaming

	// Director personnalisé pour forcer l'en-tête Host
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = remoteURL.Scheme
		req.URL.Host = remoteURL.Host
		req.URL.Path = req.URL.Path   // conserve le chemin d'origine
		req.URL.RawQuery = req.URL.RawQuery
		req.Host = remoteHostHeader   // ← FORCE l'en-tête Host
	}

	// Gestion des erreurs du proxy
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Erreur proxy : %v", err)
		http.Error(w, fmt.Sprintf("Erreur proxy : %v", err), http.StatusBadGateway)
	}

	// Route principale
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("=> %s %s (Host: %s)", r.Method, r.URL.Path, r.Host)
		proxy.ServeHTTP(w, r)
	})

	// Détermination du port d'écoute (Upsun fournit $PORT)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	log.Printf("Relais XHTTP démarré sur %s", addr)
	log.Printf("Cible : %s", target)
	log.Printf("Host forcé : %s", remoteHostHeader)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Erreur du serveur : %v", err)
	}
}