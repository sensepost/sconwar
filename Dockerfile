FROM golang:1 as build

LABEL maintainer="Leon Jacobs <leonja511@gmail.com>"

COPY . /src

WORKDIR /src
RUN make clean swagger-install deps swagger docker

# final image
FROM golang

COPY --from=build /src/sconwar /usr/local/bin

EXPOSE 8080

ENTRYPOINT ["sconwar"]
