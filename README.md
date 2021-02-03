# Demo

### Prepare
- [Golang](https://golang.org/doc/install) 
- [Docker](https://docs.docker.com/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Quick Start
- Build project image
```
docker-compose build
```

- Deploy project
```
docker-compose up -d
```

- Watch the running status
```
docker logs demo
```

### Operation
- New device
```
curl -X POST http://localhost:8080/v1/device -H 'content-type: application/json' -d '{"model": "Pro","color": "White","version": "1.0"}'
```

- Get device list
```
curl -X GET 'http://localhost:8080/v1/device?page=1&number=2'
```