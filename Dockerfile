FROM golang:1.19-bullseye AS build-stage
WORKDIR /project
COPY go.mod go.sum ./
RUN  --mount=type=cache,target=$GOPATH/pkg/mod \
    go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o keeper.say

# Deploy the application binary into a lean image
FROM ubuntu:22.04 AS build-release-stage
WORKDIR /project
COPY --from=build-stage /project/keeper.say keeper.say
EXPOSE 8080
EXPOSE 5050
RUN useradd app
USER app:app
CMD ["./keeper.say"]
