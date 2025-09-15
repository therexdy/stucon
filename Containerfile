FROM archlinux:latest

RUN pacman -Syu --noconfirm \
    && pacman -S --noconfirm go \
    && pacman -Scc --noconfirm

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/bin:$PATH

