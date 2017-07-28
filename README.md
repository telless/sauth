# Sauth

## Description

Caching auth service for WSO+SQL DB written in Golang 

## Build and start

For build just clone this repo and run
```
go build
```
This project has only one dependency:
```
go get github.com/go-sql-driver/mysql
```
After build just run binary
```
chmod +x sauth
./sauth
```
Tested with [pm2 Process Manager](https://github.com/Unitech/pm2)

```
pm2 start sauth -n auth_service --interpreter=none -x -- --sock=/path/to/sock
```
## Usage

By default server listen on /tmp/sauth.sock, you can override it if start with param
```
./sauth --sock=/path/to/unix/sock/file
```
All credentials for WSO and SQL DB stored in `parameters.json`, there is `parameters.json.example` for json struct example.

Expected input: auth token `<type> <token_string>` 
Example: `Bearer mF_9.B5f-4.1JqM`

Exprected output: json-serialized user data from DB 
Example: 
```
{
	"id":13,
	"name":"Test user",
	"email":"someone@test.local",
	"partnerId":42,
	"type":"Jedi",
	"role":127,
	"dataExpire":"2017-07-28 18:26:30",
	"tokenExpire":"2017-09-08 04:25:36"
}
```

## Bad news
Currently project hardly depends on db struct (user table), lowercased token types and so on, so probably it can be used only as example for implementation SOAP service, unix-socket server or something else.
