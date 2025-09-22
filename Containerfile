FROM docker.io/library/archlinux:latest

RUN pacman -Syu --noconfirm \
    && pacman -S --noconfirm go \
    && pacman -Scc --noconfirm

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/bin:$PATH

COPY . /app
WORKDIR /app
RUN go mod tidy
RUN go build -o /app/cmd/stucon /app/cmd/main.go
