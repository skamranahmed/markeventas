# build stage
FROM golang:1.17-alpine AS builder

# setup the working directory inside the container
WORKDIR /app

# copy the files from the host computer to the container
COPY . .

# build go binary
RUN go build -o main main.go

# run stage
FROM alpine:3.15

# setup the working directory inside the container
WORKDIR /app

# copy the binary build from the builder to the current stage container
COPY --from=builder /app/main .

# define the commands that are to executed when the container starts
CMD [ "/app/main" ]