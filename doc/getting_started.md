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
Using nodemon: https://www.npmjs.com/package/nodemon

**Install**:
```bash
npm install -g nodemon # or using yarn: yarn global add nodemon
```

**Run**:
```bash
nodemon
```


This will auto start and restart the go main server.
By default it uses `nodemon.json` config.

Running without config:
```bash
nodemon -e 'go,html' --signal SIGTERM --exec 'go' run main.go
```

#### Manual
Run `go run .`

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

