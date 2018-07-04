package etcdng

import (
	"crypto/tls"
	"crypto/x509"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

func NewTLSConfig(opt EtcdOptions) *tls.Config {
	var cfg *tls.Config = nil

	if opt.CertFile != "" && opt.KeyFile != "" {
		var rpool *x509.CertPool = nil
		if opt.CaFile != "" {
			if pemBytes, err := ioutil.ReadFile(opt.CaFile); err == nil {
				rpool = x509.NewCertPool()
				rpool.AppendCertsFromPEM(pemBytes)
			} else {
				log.Errorf("Error reading Etcd Cert CA File: %v", err)
			}
		}

		if tlsCert, err := tls.LoadX509KeyPair(opt.CertFile, opt.KeyFile); err == nil {
			cfg = &tls.Config{
				RootCAs:            rpool,
				Certificates:       []tls.Certificate{tlsCert},
				InsecureSkipVerify: opt.InsecureSkipVerify,
			}
		} else {
			log.Errorf("Error loading KeyPair for TLS client: %v", err)
		}
	}

	// If InsecureSkipVerify is provided, assume TLS
	if (opt.EnableTLS || opt.InsecureSkipVerify) && cfg == nil {
		cfg = &tls.Config{
			InsecureSkipVerify: opt.InsecureSkipVerify,
		}
	}
	return cfg
}
