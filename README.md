# Abfun - Bilibili Clone (Microservices)

This project is a microservices-based clone of Bilibili, designed for low-resource environments.

## Architecture

- **Services**: User, Video, Danmaku, Interaction
- **Database**: MySQL 8.0 (Optimized for low memory)
- **Cache**: Redis 7.0
- **Message Queue**: RabbitMQ 3.x
- **Gateway**: Nginx (Reverse Proxy + SSL)
- **Storage**: Local Filesystem (Mapped via Nginx)

## Features

- **User System**: Register, Login (JWT)
- **Video System**: Upload (Local Storage), Feed, Tags
- **Danmaku System**: WebSocket Real-time Danmaku
- **Interaction System**: Likes (RabbitMQ Async), Favorites

## Setup & Deployment

### Prerequisites
- Docker (with Compose plugin)
- Domain pointing to this server (e.g., `api.abfun.me`)

### 1. Initialize SSL Certificates
Run the initialization script to set up Let's Encrypt certificates.
```bash
chmod +x init-ssl.sh
./init-ssl.sh
```

### 2. Start Services
Start all services in the background.
```bash
docker compose up -d
```

### 3. Verify
Check if all containers are running healthy.
```bash
docker compose ps
```

## Configuration

- **Nginx**: `deploy/nginx/nginx.conf`
- **Docker**: `docker-compose.yml`
- **Environment**: Configured via `environment` variables in `docker-compose.yml`.

## API Endpoints

- **User**: `/api/v1/user/register`, `/api/v1/user/login`
- **Video**: `/api/v1/videos` (POST/GET), `/api/v1/tags`
- **Danmaku**: `/api/v1/danmaku/ws` (WebSocket), `/api/v1/danmaku/list`
- **Interaction**: `/api/v1/interaction/like`, `/api/v1/interaction/favorites`

## Notes
- MySQL is configured to use minimal memory (~512MB limit).
- Video uploads are stored in `./uploads` on the host machine.
