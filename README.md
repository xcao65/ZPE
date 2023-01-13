# ZPE

ZPE User API

### How to Run

go to root of this repo

```bash
export GOPATH="$HOME/go"
export PATH=$PATH:$(go env GOPATH)/bin
go mod tidy
go run cmd/user/main.go
```

Server will listen on http://localhost:3000

Sample request
Create User
POST http://localhost:3000/users

```json
{
  "name": "jon",
  "email": "jon@email.com",
  "role": ["Admin"]
}
```


List Users 
GET http://localhost:3000/users
