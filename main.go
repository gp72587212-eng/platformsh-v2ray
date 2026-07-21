package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	remoteHost       string
	remotePort       string
	remoteHostHeader string
)

func init() {
	remoteHost = getEnv("REMOTE_HOST", "62.171.180.164")
	remotePort = getEnv("REMOTE_PORT", "80")
	remoteHostHeader = getEnv("REMOTE_HOST_HEADER", "")

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

const landingPage = `<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Proxy XHTTP</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            margin: 0;
            padding: 0;
            background: #f6f8fa;
            color: #24292e;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
        }
        .container {
            max-width: 600px;
            padding: 40px;
            text-align: center;
            background: white;
            border-radius: 16px;
            box-shadow: 0 4px 12px rgba(0,0,0,0.1);
        }
        h1 { font-size: 28px; font-weight: 600; color: #0366d6; }
        p { color: #586069; line-height: 1.6; margin: 20px 0; }
        .status {
            display: inline-block;
            padding: 6px 16px;
            background: #dcffe4;
            color: #28a745;
            border-radius: 20px;
            font-size: 14px;
            font-weight: 600;
        }
        .badge {
            display: inline-block;
            padding: 4px 12px;
            background: #f1f8ff;
            color: #0366d6;
            border-radius: 12px;
            font-size: 13px;
            margin: 4px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Proxy XHTTP</h1>
        <p>Ce service est un relais sécurisé pour les connexions VLESS via le protocole XHTTP.</p>
        <div class="status">● Service Operational</div>
        <br><br>
        <div>
            <span class="badge">VLESS</span>
            <span class="badge">XHTTP</span>
            <span class="badge">TLS</span>
        </div>
        <p style="font-size: 14px; color: #6a737d; margin-top: 30px;">
            Ce proxy est optimisé pour les connexions avec <strong>Orange Cameroun</strong> 🇨🇲
        </p>
    </div>
</body>
</html>`

func main() {
	target := fmt.Sprintf("http://%s:%s", remoteHost, remotePort)
	remoteURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("URL cible invalide : %v", err)
	}

	transport := &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		DisableCompression:    false,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	proxy := httputil.NewSingleHostReverseProxy(remoteURL)
	proxy.Transport = transport
	proxy.FlushInterval = 100 * time.Millisecond

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = remoteURL.Scheme
		req.URL.Host = remoteURL.Host
		req.URL.Path = req.URL.Path
		req.URL.RawQuery = req.URL.RawQuery
		if remoteHostHeader != "" {
			req.Host = remoteHostHeader
		}
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Erreur proxy : %v", err)
		http.Error(w, fmt.Sprintf("Erreur proxy : %v", err), http.StatusBadGateway)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")

		// Détection navigateur
		if strings.Contains(userAgent, "Mozilla") || strings.Contains(userAgent, "Chrome") || strings.Contains(userAgent, "Safari") || strings.Contains(userAgent, "Edg") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, landingPage)
			log.Printf("🌐 Landing page affichée pour %s", r.RemoteAddr)
			return
		}

		// Requête VLESS/XHTTP → proxy
		log.Printf("=> %s %s (Host: %s)", r.Method, r.URL.Path, r.Host)
		proxy.ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	log.Printf("Relais XHTTP démarré sur %s", addr)
	log.Printf("Cible : %s", target)
	if remoteHostHeader != "" {
		log.Printf("Host forcé : %s", remoteHostHeader)
	}

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Erreur serveur : %v", err)
	}
}