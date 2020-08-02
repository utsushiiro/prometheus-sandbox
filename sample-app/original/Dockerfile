FROM golang:1.14-alpine as build

ENV APP_NAME sample-app
ENV WORK_DIR /work_dir
RUN mkdir -p ${WORK_DIR}
WORKDIR ${WORK_DIR}
COPY . .
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o ${APP_NAME}

FROM alpine:3.12
ENV APP_NAME sample-app
ENV WORK_DIR /work_dir
COPY --from=build ${WORK_DIR}/${APP_NAME} /usr/local/bin/${APP_NAME}

# https://stackoverflow.com/questions/37904682/how-do-i-use-docker-environment-variable-in-entrypoint-array
ENTRYPOINT ["sh", "-c", "/usr/local/bin/${APP_NAME}"]
