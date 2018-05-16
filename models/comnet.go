package models

import (
	"net"
	"time"
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"bytes"
	"golib/modules/logr"
	"io/ioutil"
)

type G_conf struct {
	HttpUrl string
}

var gconf G_conf

func Comm(ReqBuf []byte) ([]byte, error) {
	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 60}
	body := bytes.NewBuffer(ReqBuf)

	req, _ := http.NewRequest("POST", gconf.HttpUrl, body)
	req.Header.Set("Content-Type", "application/soap+xml;charset=UTF-8")

	out, _ := httputil.DumpRequestOut(req, true)
	logr.Infof("http 请求包:\n-------------------\n [%s]\n-------------------\n ", out)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logr.Infof("http.Status: %s is not success!", resp.Status)
	}
	RspBuf, _ := ioutil.ReadAll(resp.Body)
	logr.Infof("响应报文：\n-------------------\n[%s]\n-------------------\n", RspBuf)

	return RspBuf, nil
}
