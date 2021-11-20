# build stage
FROM registry.access.redhat.com/ubi8/ubi:8.5 AS build-env
RUN dnf install -y golang
ADD . /src
RUN cd /src && go build -o ff-tab-mqtt

# final stage
FROM registry.access.redhat.com/ubi8/ubi-minimal:8.5
WORKDIR /app
COPY --from=build-env /src/ff-tab-mqtt /app/
ENTRYPOINT ./ff-tab-mqtt
