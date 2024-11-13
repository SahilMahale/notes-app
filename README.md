# Notes-app
A simple notes taking app made as a reference, showing how to structure a go codebase, build and test it  
while obeying the sonar cloud Quality gates

## Development 

### Dependencies 
* Backend  
    * Go [1.23.2 or higher](https://go.dev/doc/)
    * fiber [http server](https://gofiber.io/)
    * sqlite3 [DB](https://www.sqlite.org/index.html) 
    * gorm [db ORM](https://gorm.io/)  
* Frontend  
    * npm [23.1.0 or higher](https://nodejs.org/en)
    * Typescript v5 or higher
    * npm [package manager](https://www.npmjs.com/)
    * vite [bundler/dev server](https://vitejs.dev/guide/why.html)
    * react

### Backend

#### Making RSA CERTs required for the JWTs
```bash
mkdir secrets
cd secrets
openssl genrsa -out private_key.pem 2048
openssl rsa -in private_key.pem -outform PEM -pubout -out public_key.pem.pub
```
#### Export the Absolute path to secrets
```bash
export APP_AUTH=<abs_path>/secrets
```

#### Starting backend for dev
```bash
cd notes-backend/
go mod tidy
go run cmd/main.go
```

