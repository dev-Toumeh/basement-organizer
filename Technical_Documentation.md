# Technical Documentation
## Start Development 

Run build:
    cd hello
    go run .

# Decisions
- project's main focus is learning `Go` in backend and frontend in general (starting with `HTMX` as a lightweight library)
- solve a real world problem
- other goals are simplicity, fast and simple development, desktop and mobile support
    - only browser based, no native apps
    - focus on necessary core functionalities
    - avoiding features that are nice to have
    - less frameworks/dependencies
- have more time for relevant problems and learning
    - use Github as main platform (central place for information)
        - version control
        - documentation
        - issues 
        - no infrastructure needed
        - less searching
        - ability to use CI/CD (Github Actions) at a later time if necessary

## Overview

## Technologies used
- Go backend
- HTMX frontend
- 

## Development Environment Configuration
### Tree
- the environment tree look like the following:


    .
    └── basment-organizer/
        ├── docker/
        │   └── Dockerfile
        ├── docker-compose.yml
        └── .env

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



