# Nōto
A simple notes taking app made as a reference, showing how to structure a go codebase, build and test it
while obeying the sonar cloud Quality gates

## Development

### Dependencies
* Backend
    * Go [1.23.2 or higher](https://go.dev/doc/)
    * fiber [http server](https://gofiber.io/)
    * sqlite3 [DB](https://www.sqlite.org/index.html)
    * gorm [db ORM](https://gorm.io/)
    * [air](https://github.com/air-verse/air) (live reload)
* Frontend
    * Typescript v5 or higher
    * Bun
    * React with [TanStack Router](https://tanstack.com/router)
    * [Vite](https://vite.dev/)
    * [Tailwind CSS](https://tailwindcss.com/) v4

### Setup

#### Making RSA CERTs required for the JWTs
```bash
mkdir secrets
cd secrets
openssl genrsa -out private_key.pem 2048
openssl rsa -in private_key.pem -outform PEM -pubout -out public_key.pem.pub
```

#### Install dependencies
```bash
make install
```

### Running

#### Start both backend and frontend (recommended)
```bash
make dev
```

#### Start backend only (air live reload on :8001)
```bash
make backend
```

#### Start frontend only (vite dev server)
```bash
make frontend
```

#### Stop all services
```bash
make stop
```

#### Build backend binary
```bash
make build-backend
```

