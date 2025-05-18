# ┌────────────── Frontend build ──────────────┐
FROM node:18 AS frontend
WORKDIR /app/frontend
# install deps & build
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# ┌────────────── Backend build ──────────────┐
FROM golang:1.24.2-bookworm AS builder
WORKDIR /app
# pull in module files and download deps
COPY go.* ./
COPY go.mod go.sum ./
RUN go mod download
# copy everything else
COPY . .
RUN echo "Contents of /app in builder stage:" && ls -R /app
RUN echo "Contents of /app/assets in builder stage:" && ls -R /app/assets || echo "/app/assets not found or empty in builder"

RUN CGO_ENABLED=0 go build -o main main.go

# ┌────────────── Final image ──────────────┐
FROM debian:bookworm
# (optional: add a non-root user here)
WORKDIR /app
# copy in the compiled binary
COPY --from=builder /app/bin/api .
# copy in the game data assets
COPY --from=builder /app/assets ./assets
RUN echo "Contents of /app/assets in final image:" && ls -R /app/assets || echo "/app/assets not found or empty in final image"
# copy in the static assets
COPY --from=frontend /app/frontend/dist ./frontend/dist

# if your Go HTTP server is set up to serve "./frontend/build" for static files:
EXPOSE 8080
# Run the API binary on container start
CMD ['./main']
