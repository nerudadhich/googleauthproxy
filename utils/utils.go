package utils

import (
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type HttpRequest struct {
	Method string
	Url string
	Body []byte
}

var requestMap = make(map[string]*HttpRequest)

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func CacheRequest(r *HttpRequest) string{
	rand.Seed(time.Now().UnixNano())
	bucket := RandString(8)
	requestMap[bucket] = r
	return bucket
}

func PopRequest(bucket string) *HttpRequest{
	_, ok := requestMap[bucket]
	var request *HttpRequest
	if ok {
		request = requestMap[bucket]
		delete(requestMap, bucket)
	}
	return request
}