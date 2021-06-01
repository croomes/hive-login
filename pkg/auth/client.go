package auth

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// httpClient return an HTTP client which trusts the provided root CAs.
func httpClient(rootCAs string) (*http.Client, error) {
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if rootCAs != "" {
		tlsConfig := tls.Config{RootCAs: x509.NewCertPool()}
		rootCABytes, err := ioutil.ReadFile(rootCAs)
		if err != nil {
			return nil, fmt.Errorf("failed to read root-ca: %v", err)
		}
		if !tlsConfig.RootCAs.AppendCertsFromPEM(rootCABytes) {
			return nil, fmt.Errorf("no certs found in root CA file %q", rootCAs)
		}
		transport.TLSClientConfig = &tlsConfig
	}

	return &http.Client{
		Transport: &transport,
	}, nil
}
