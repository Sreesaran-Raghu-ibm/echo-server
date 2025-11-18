# Docker Deployment Guide

This guide explains how to build and run the WebSocket Echo Server using Docker.

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Build and start the container
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the container
docker-compose down
```

### Using Docker CLI

```bash
# Build the image
docker build -t websocket-echo-server .

# Run the container
docker run -d -p 8080:8080 --name echo-server websocket-echo-server

# View logs
docker logs -f echo-server

# Stop and remove the container
docker stop echo-server && docker rm echo-server
```

## Accessing the Server

Once running, access the server at:
- **Web Interface**: http://localhost:8080
- **WebSocket Endpoint**: ws://localhost:8080/ws

## Docker Image Details

### Multi-Stage Build

The Dockerfile uses a multi-stage build process:

1. **Builder Stage** (golang:1.21-alpine)
   - Initializes Go module
   - Downloads dependencies (gorilla/websocket)
   - Compiles static binary with optimizations

2. **Runtime Stage** (alpine:latest)
   - Minimal Alpine Linux base (~5MB)
   - Non-root user for security
   - Health check enabled
   - Final image size: ~15-20MB

### Security Features

- ✅ Non-root user (appuser:1000)
- ✅ Static binary (no runtime dependencies)
- ✅ Minimal attack surface
- ✅ Health checks enabled
- ✅ CA certificates included

## Configuration

### Environment Variables

- `PORT`: Server port (default: 8080)

Example with custom port:
```bash
docker run -d -p 9000:9000 -e PORT=9000 websocket-echo-server
```

### Custom Build Arguments

Build for different architectures:
```bash
# For ARM64
docker build --platform linux/arm64 -t websocket-echo-server:arm64 .

# For AMD64
docker build --platform linux/amd64 -t websocket-echo-server:amd64 .
```

## Health Checks

The container includes a health check that:
- Runs every 30 seconds
- Times out after 3 seconds
- Retries 3 times before marking unhealthy
- Starts checking after 5 seconds

Check container health:
```bash
docker inspect --format='{{.State.Health.Status}}' echo-server
```

## Troubleshooting

### View Container Logs
```bash
docker logs echo-server
# or with docker-compose
docker-compose logs echo-server
```

### Interactive Shell Access
```bash
docker exec -it echo-server sh
```

### Rebuild Without Cache
```bash
docker build --no-cache -t websocket-echo-server .
# or with docker-compose
docker-compose build --no-cache
```

### Check Container Stats
```bash
docker stats echo-server
```

## Production Deployment

### Using Docker Swarm

```bash
# Initialize swarm
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.yml echo-stack

# Scale service
docker service scale echo-stack_echo-server=3

# Remove stack
docker stack rm echo-stack
```

### Using Kubernetes

Create a deployment:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: echo-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: echo-server
  template:
    metadata:
      labels:
        app: echo-server
    spec:
      containers:
      - name: echo-server
        image: websocket-echo-server:latest
        ports:
        - containerPort: 8080
        livenessProbe:
          httpGet:
            path: /
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 30
---
apiVersion: v1
kind: Service
metadata:
  name: echo-server
spec:
  selector:
    app: echo-server
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

## Image Management

### Tag and Push to Registry

```bash
# Tag image
docker tag websocket-echo-server:latest your-registry/websocket-echo-server:v1.0.0

# Push to registry
docker push your-registry/websocket-echo-server:v1.0.0
```

### Clean Up

```bash
# Remove container
docker rm -f echo-server

# Remove image
docker rmi websocket-echo-server

# Clean up all unused resources
docker system prune -a
```

## Performance Tips

1. **Use BuildKit** for faster builds:
   ```bash
   DOCKER_BUILDKIT=1 docker build -t websocket-echo-server .
   ```

2. **Multi-platform builds**:
   ```bash
   docker buildx build --platform linux/amd64,linux/arm64 -t websocket-echo-server .
   ```

3. **Resource limits**:
   ```bash
   docker run -d \
     --memory="256m" \
     --cpus="0.5" \
     -p 8080:8080 \
     websocket-echo-server
   ```

## Testing the Deployment

### Using curl
```bash
# Test HTTP endpoint
curl http://localhost:8080

# Test WebSocket (requires websocat)
websocat ws://localhost:8080/ws
```

### Using Docker Network
```bash
# Create custom network
docker network create echo-net

# Run with custom network
docker run -d --network echo-net --name echo-server websocket-echo-server

# Test from another container
docker run --rm --network echo-net alpine/curl curl http://echo-server:8080
```

## Support

For issues or questions:
- Check container logs: `docker logs echo-server`
- Verify health status: `docker inspect echo-server`
- Review Dockerfile and docker-compose.yml configurations