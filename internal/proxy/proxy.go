package proxy

import (
	"Proxy/internal/usecase"
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
)

const portNum string = ":8080"

type ProxyHandler struct {
	uc usecase.UsecaseI
}

func (p *ProxyHandler) ProxyHTTPS(w http.ResponseWriter, r *http.Request) {
	p.handleTunneling(w, r)
}

func (p *ProxyHandler) ProxyHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Proxy-Connection")
	r.RequestURI = "/"

	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	reqId, err := p.uc.SaveRequest(r)
	if err != nil {
		log.Printf("Error save:  %v", err)
	}

	_, err = p.uc.SaveResponse(reqId, resp)
	if err != nil {
		log.Printf("Error save:  %v", err)
	}
}

func (p *ProxyHandler) handleTunneling(w http.ResponseWriter, r *http.Request) {
	log.Printf("CONNECT requested to %v (from %v)", r.Host, r.RemoteAddr)

	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Fatal("http server doesn't support hijacking connection")
	}

	clientConn, _, err := hj.Hijack()
	if err != nil {
		log.Fatal("http hijacking failed")
	}

	host, _, err := net.SplitHostPort(r.Host)
	err = genCert(host)
	if err != nil {
		log.Fatal("generate cert failed")
	}

	tlsCert, err := tls.LoadX509KeyPair("/certs/certificate.crt", "/certs/private_key.key")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := clientConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n")); err != nil {
		log.Fatal("error writing status to client:", err)
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:               tls.VersionTLS13,
		Certificates:             []tls.Certificate{tlsCert},
	}

	tlsConn := tls.Server(clientConn, tlsConfig)
	defer tlsConn.Close()

	connReader := bufio.NewReader(tlsConn)

	req, err := http.ReadRequest(connReader)
	if err != nil {
		log.Fatal(err)
	}

	if b, err := httputil.DumpRequest(r, false); err == nil {
		log.Printf("incoming request:\n%s\n", string(b))
	}

	changeRequestToTarget(req, r.Host)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("error sending request to target:", err)
	}
	if b, err := httputil.DumpResponse(resp, false); err == nil {
		log.Printf("target response:\n%s\n", string(b))
	}
	defer resp.Body.Close()

	if err := resp.Write(tlsConn); err != nil {
		log.Println("error writing response back:", err)
	}

	reqId, err := p.uc.SaveRequest(req)
	if err != nil {
		log.Printf("Error save:  %v", err)
	}

	_, err = p.uc.SaveResponse(reqId, resp)
	if err != nil {
		log.Printf("Error save:  %v", err)
	}

}

func changeRequestToTarget(req *http.Request, targetHost string) {
	targetUrl := addrToUrl(targetHost)
	targetUrl.Path = req.URL.Path
	targetUrl.RawQuery = req.URL.RawQuery
	req.URL = targetUrl
	req.RequestURI = ""
}

func addrToUrl(addr string) *url.URL {
	if !strings.HasPrefix(addr, "https") {
		addr = "https://" + addr
	}
	u, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}
	return u
}

func genCert(host string) error {
	log.Println(host)
	cmd := exec.Command("/certs/gen_cert.sh", host, strconv.Itoa(rand.Int()))

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		log.Fatal(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	return nil
}

func (ph *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		ph.ProxyHTTPS(w, r)
	} else {
		ph.ProxyHTTP(w, r)
	}
}

func Run(u usecase.UsecaseI) {
	proxyHandler := &ProxyHandler{uc: u}
	http.ListenAndServe(portNum, proxyHandler)
}
