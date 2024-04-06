# eCar

![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

# 🎲 Service on Go(Gin) for for managing a car catalog 🎲

# 📞 Endpoints
```http
POST /api/add
- add new car
POST /api/all
- Get cars
DELETE /api/delete/:id
- delete car
PUT /api/update/:id
- update info about car

```

# 🏴‍☠️ Flags
```
a - ip for REST -a=host
cert - path to certificate -cert=path_to_certificate
key - path to private key -key=path_to_key
tls - enable or disable tls certificate -tls=false/true
d - connection string -d=connection string

```

# 💻 Config.env
```
HOST=0.0.0.0:8081
DATABASE_DSN=postgres://user:password@db:5432/dbname?sslmode=disable
TLS="false"
```

# 💎 Build
```

docker compose up --build

```
