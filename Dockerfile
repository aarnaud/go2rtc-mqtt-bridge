FROM golang:alpine AS builderimage
WORKDIR /go/src/go2rtc-mqtt-bridge
COPY . .
RUN go build -o go2rtc-mqtt-bridge main.go


###################################################################

FROM alpine
COPY --from=builderimage /go/src/go2rtc-mqtt-bridge/go2rtc-mqtt-bridge /app/
WORKDIR /app
CMD ["./go2rtc-mqtt-bridge"]
