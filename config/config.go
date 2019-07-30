package config

import (
	"os"
)

// ServiceName name of service existing in app service
var ServiceName = "qr"

// Port port of server to open and listen to connections
var Port = GetEnv("PORT_TEST", "10000")

// AuthToken for internal requests
//var AuthToken = utils.GetEnv("INTERNAL_REQUEST_AUTHENTICATION_TOKEN", "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJuYW1laWQiOiJtZXNzYWdpbmctc2VydmljZSIsImlkIjoiNTk4MWExZTQxZDQxYzg0Y2FlOTA0ZmQyIiwidW5pcXVlX25hbWUiOiJtZXNzYWdpbmctc2VydmljZSIsInN1YiI6Im1lc3NhZ2luZy1zZXJ2aWNlIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdC8iLCJhdWQiOiJiOWRjNzEyYzk1MmI0YWFmYjQ4MWFiZWRlMGZlYzRkOCIsImV4cCI6OTk5OTk5OTk5OSwibmJmIjoxNDk3MTc4MjQ1LCJyb2xlIjpbInpldXMiXX0.MB6XL2wfqO84KM7geXy3foQiq9uH0nDIh3LUk_VqDU0")
var AuthToken = GetEnv("INTERNAL_REQUEST_AUTHENTICATION_TOKEN", "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJuYW1laWQiOiJhY2NvdW50cy1zZXJ2aWNlIiwiaWQiOiI1OTgxYTFlNDFkNDFjODRjYWU5MDRmZDMiLCJ1bmlxdWVfbmFtZSI6ImFjY291bnRzLXNlcnZpY2UiLCJzdWIiOiJhY2NvdW50cy1zZXJ2aWNlIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdC8iLCJyb2xlIjpbInpldXMiLCJyb3N0YW0iXSwiYXVkIjoiYjlkYzcxMmM5NTJiNGFhZmI0ODFhYmVkZTBmZWM0ZDgiLCJleHAiOjk5OTk5OTk5OTksIm5iZiI6MTQ5NzE3ODI0NSwiYXBwIjoiNWE5NTYzZjM4NDllMDY3NzEwNWRmNTI5In0.fyU5e4KXpZilnDcxhKRkbYw0paAX15RNGXpifgWvHbY")

// API Key
var APIKey = GetEnv("APIKey", "5aa7e856ae7fbc00016ac5a0ede56b6989e14706a6215f4207a40996")

// ClientID of current application
var ClientID = GetEnv("CLIENT_ID", "b9dc712c952b4aafb481abede0fec4d8")

// ClientSecret of current application
var ClientSecret = GetEnv("CLIENT_SECRET", "J0RYjUcIZHgm41GyPt4wEWUqKzOPXCQAY7n2/ZkQ7WE=")

//MongoHost url
var MongoHost = GetEnv("MongoHost", "localhost")

// var MongoURI = GetEnv("MongoURI", "mongodb://192.168.41.221:32017")
var MongoURI = GetEnv("MongoURI", "mongodb://localhost:27017")

var DatabaseName = GetEnv("DBNAME", "hermes_rc")

//MongoDBName
var MongoDBName = GetEnv("MongoDBName", "hermes_rc")

var ApplicationServiceURL = GetEnv("APPLICATION_SERVICE_URL", "https://api.paygear.ir/application/v3")

// Club Base URL
var ClubBaseURL = GetEnv("ClubBaseURL", "https://api.paygear.ir/club")

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
