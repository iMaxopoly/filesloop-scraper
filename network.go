package main

import (
	"time"

	"github.com/valyala/fasthttp"
)

func requestGetFollowRedirect(link string, timeout uint) (body []byte, err error) {
	_, body, err = fasthttp.GetTimeout(nil, link, time.Duration(timeout)*time.Second)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func requestGet(link string, timeout uint) (body []byte, header fasthttp.ResponseHeader, err error) {
	redirCount := 0
redir:
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURI(link)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.SetMethod("GET")

	err = fasthttp.DoTimeout(req, resp, time.Duration(timeout)*time.Second)
	if err != nil {
		return []byte{}, fasthttp.ResponseHeader{}, err
	}

	body = resp.Body()
	header = resp.Header

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)

	location := string(header.Peek("Location"))
	if location != "" && location != link && redirCount < 15 {
		debugLog("redir anticipated", link, "->", location)
		link = location
		redirCount++
		goto redir
	}

	return body, header, nil
}

func requestPost(link string, requestBody []byte, timeout uint) (body []byte, header fasthttp.ResponseHeader, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(link)
	req.Header.Add("Requested-With-AngularJS", "true")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	req.Header.SetMethod("POST")
	req.SetBody(requestBody)

	if err = fasthttp.DoTimeout(req, resp, time.Duration(timeout)*time.Second); err != nil {
		return []byte{}, fasthttp.ResponseHeader{}, err
	}

	body = resp.Body()
	header = resp.Header
	return body, header, nil
}
