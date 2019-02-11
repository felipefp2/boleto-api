package util

import (
	"crypto"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/mundipagg/goseq"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	s "github.com/fullsailor/pkcs7"
	"github.com/mundipagg/boleto-api/config"
)

var defaultDialer = &net.Dialer{Timeout: 16 * time.Second, KeepAlive: 16 * time.Second}

var cfg *tls.Config = &tls.Config{
	InsecureSkipVerify: true,
}
var client *http.Client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig:     cfg,
		Dial:                defaultDialer.Dial,
		TLSHandshakeTimeout: 16 * time.Second,
	},
}

// DefaultHTTPClient retorna um cliente http configurado para dar um skip na validação do certificado digital
func DefaultHTTPClient() *http.Client {
	return client
}

//Post faz um requisição POST para uma URL e retorna o response, status e erro
func Post(url, body, timeout string, header map[string]string) (string, int, error) {
	return doRequest("POST", url, body, timeout, header)
}

//Get faz um requisição GET para uma URL e retorna o response, status e erro
func Get(url, body, timeout string, header map[string]string) (string, int, error) {
	return doRequest("GET", url, body, timeout, header)
}

func doRequest(method, url, body, timeout string, header map[string]string) (string, int, error) {
	client := DefaultHTTPClient()
	client.Timeout = GetDurationTimeoutRequest(timeout) * time.Second
	message := strings.NewReader(body)
	req, err := http.NewRequest(method, url, message)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	resp, errResp := client.Do(req)
	if errResp != nil {
		return "", 0, errResp
	}
	defer resp.Body.Close()
	data, errResponse := ioutil.ReadAll(resp.Body)
	if errResponse != nil {
		return "", resp.StatusCode, errResponse
	}
	sData := string(data)
	return sData, resp.StatusCode, nil
}

//BuildTLSTransport creates a TLS Client Transport from crt, ca and key files
func BuildTLSTransport(crtPath string, keyPath string, caPath string) (*http.Transport, error) {
	cert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	return &http.Transport{TLSClientConfig: tlsConfig}, nil
}

//Sigs request
func SignRequest(request string, requestKey string) (string, error) {

	pkey, md5Pkey, err := parsePrivateKey()
	if err != nil {
		return "", err
	}

	chainCertificates, md5Chain, err := parseChainCertificates()
	if err != nil {
		return "", err
	}

	signedData, err := s.NewSignedData([]byte(request))
	if err != nil {
		return "", err
	}

	if err := signedData.AddSigner(chainCertificates, pkey, s.SignerInfoConfig{}); err != nil {
		return "", err
	}

	detachedSignature, err := signedData.Finish()
	if err != nil {
		return "", err
	}

	signedRequest := base64.StdEncoding.EncodeToString(detachedSignature)

	LogAssign(len(detachedSignature), signedRequest, md5Pkey, md5Chain, requestKey)

	return signedRequest, nil
}

func LogAssign(length int, signed string, md5Pkey string, md5Chain string, requestKey string) {
	seqLog, _ = goseq.GetLogger(config.Get().SEQUrl, config.Get().SEQAPIKey)
	prop := goseq.NewProperties()
	prop.AddProperty("Application", "BoletoOnline")
	prop.AddProperty("length DetachedSignature", length)
	prop.AddProperty("signedRequest", signed)
	prop.AddProperty("MD5 PrivateKey", md5Pkey)
	prop.AddProperty("MD5 ChainCA", md5Chain)
	prop.AddProperty("Operation", "Assigned")
	prop.AddProperty("RequestKey", requestKey)
	head := "[BoletoOnline: RegisterBoleto] - Assinatura request BradescoNetEmpresa - length DetachedSignature: " +strconv.Itoa(length)
	seqLog.Debug(head, prop)
}

//Read privatekey and parse to PKCS#1
func parsePrivateKey() (crypto.PrivateKey, string, error) {

	var md5PrivateKey string

	pkeyBytes, err := ioutil.ReadFile(config.Get().CertICP_PathPkey)
	if err != nil {
		return nil, "", err
	}

	if len(pkeyBytes) != 0{
		md5PrivateKeyByte := md5.Sum(pkeyBytes)
		md5PrivateKey = hex.EncodeToString(md5PrivateKeyByte[:])
	} else{
		md5PrivateKey = "Certificado não encontrado."
	}

	block, _ := pem.Decode(pkeyBytes)
	if block == nil {
		return nil, md5PrivateKey, errors.New("Key Not Found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, md5PrivateKey, err
		}
		return rsa, md5PrivateKey, nil
	default:
		return nil, md5PrivateKey, fmt.Errorf("SSH: Unsupported key type %q", block.Type)
	}

}

///Read chainCertificates and adapter to x509.Certificate
func parseChainCertificates() (*x509.Certificate, string, error) {
	var md5Chain string

	chainCertsBytes, err := ioutil.ReadFile(config.Get().CertICP_PathChainCertificates)
	if err != nil {
		return nil, "", err
	}

	if len(chainCertsBytes) != 0{
		md5ChainByte := md5.Sum(chainCertsBytes)
		md5Chain = hex.EncodeToString(md5ChainByte[:])
	} else{
		md5Chain = "Certificado não encontrado."
	}

	block, _ := pem.Decode(chainCertsBytes)
	if block == nil {
		return nil, md5Chain, errors.New("Key Not Found")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, md5Chain, err
	}

	return cert, md5Chain, nil
}

func doRequestTLS(method, url, body, timeout string, header map[string]string, transport *http.Transport) (string, int, error) {
	var client *http.Client = &http.Client{
		Transport: transport,
	}
	client.Timeout = GetDurationTimeoutRequest(timeout) * time.Second
	b := strings.NewReader(body)
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return "", 0, err
	}

	if header != nil {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	// Dump response
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	sData := string(data)
	return sData, resp.StatusCode, nil
}

func PostTLS(url, body, timeout string, header map[string]string, transport *http.Transport) (string, int, error) {
	return doRequestTLS("POST", url, body, timeout, header, transport)
}

//HeaderToMap converte um http Header para um dicionário string -> string
func HeaderToMap(h http.Header) map[string]string {
	m := make(map[string]string)
	for k, v := range h {
		m[k] = v[0]
	}
	return m
}

var (
	seqLog *goseq.Logger
)

