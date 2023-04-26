\*\* Use this Command to build the Docker Image

# Push Docker Image with Builder

```bash
docker buildx build --platform linux/amd64,linux/arm64 -t monzim/uiudiscordbot:1.0.2 -f docker/Dockerfile . --push
```

# Create Multiplatfrom Build

```bash
docker buildx create --use --name mybuilder
```

# Builder

````bash
docker buildx build --platform linux/amd64,linux/arm64 -t monzim/uiudiscordbot:1.0.2 -f docker/Dockerfile . --push
```

# Build Docker Image with Builder with Tag

```bash
docker buildx build --builder mybuilder -t monzim/uiudiscordbot:latest -t monzim/uiudiscordbot:latest -f docker/Dockerfile .
````
