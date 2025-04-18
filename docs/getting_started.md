# Getting Started
<!--toc:start-->
- [Getting Started](#getting-started)
    - [Start Server](#start-server)
      - [Automatic](#automatic)
      - [Manual](#manual)
  - [Development Environment Configuration](#development-environment-configuration)
    - [Tree](#tree)
    - [Modifications Process](#modifications-process)
    - [General commands](#general-commands)
<!--toc:end-->

### Start Server
#### Automatic
Using npm:

**Install**:
```bash
npm install
```

**Run**:
```bash
npm start
```

This will auto start and restart the Go main server using nodemon. By default,
it uses the nodemon.json config.

Running without config:
```bash
nodemon -e 'go,html' --signal SIGTERM --exec 'go' run main.go
```

#### Manual
```bash
curl -o internal/static/js https://unpkg.com/htmx.org@2.0.1/dist/htmx.min.js
Run `go run .`
```
### Generate Template constants
Manually generate with:
```bash
cd ./tools
go run ./template_constants_generator.go
```

Or automatically on html change with:
```bash
cd ./tools
nodemon
```

## Testing
### run all the tests 
open the Project root directory and run the following 

```
go test  ./...
```
#### options
-add -v flag to print any fmt logs
-add -count=1 to prevent using the cache 

## Development Environment Configuration
### Tree
- the environment tree look like the following:

```
.
└── basment-organizer/
    ├── docker/
    │   └── Dockerfile
    ├── docker-compose.yml
    └── .env
```

### Modifications Process
- To modify the development environment, first switch to the **setup-maintenance-dev-env** branch, commit your changes,
  and then merge those changes into the main branch.
- These process need to be followed every time modifications to the development environment are required.

### General commands

 - docker compose  run
```bash
 docker-compose up -d
```
- docker compose stop
```bash
 docker-compose stop
```
- docker-compose stop and remove the Services
```bash
 docker-compose down
```
- start Terminal inside docker-compose service
```bash
docker-compose exec -it <service-name> /bin/bash
```
- start Terminal inside docker-compose service as root user
```bash
docker-compose exec -it -u 0 <service-name> /bin/bash
```
- build the service again after you made changes in dockerfile
```bash
 docker-compose build <service-name>
```
- build the service again after you made changes in dockerfile without cache
```bash
 docker-compose build <service-name> --no-cache
```
- run and build the service after you made changes in dockerfile
```bash
 docker-compose up -d <service-name> --build
```

