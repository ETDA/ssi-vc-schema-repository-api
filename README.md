# VC Schema API

## Introduction
\
This repository contain 2 services
- vc-schema-api
- validator-api

It's easier to maintain because these 2 services always communicate with each other and easier to develop the services.

However, we run the services separately, by define the environment when start the service.

## Step to start the service
- Copy file `.env.sample` to `.env`
- run `docker-compose up -d`
- you can access the service via `http:localhost:{port}`
- port will be specified in `docker-compose.yml`
- services defualt port:
    - vc-schema-api : `8080`
    - validator-api : `8082`