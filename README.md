# go2rtc MQTT Bridge

## Features:

- [X] MQTT integration for Home Assistant
- [X] Play on speaker

## Config:

Exemple `go2rtc-mqtt-bridge.yaml`

```yaml
AUDIO_FILES_PATH: "/PATH_TO_YOUR_AUDIO_FILES"
MQTT_BROKER_HOST: "192.168.1.245" # require if MQTT_ENABLED
MQTT_BROKER_PORT: 1883
MQTT_CLIENT_ID: "go2rtc"
MQTT_BASE_TOPIC: "go2rtc"
MQTT_USERNAME: "go2rtc"
MQTT_PASSWORD: "CHANGEME"
GO2RTC_USERNAME: "admin"    # optional
GO2RTC_PASSWORD: "CHANGEME" # optional
```


## Example with mosquitto 

```bash
mosquitto_pub -h 192.168.1.123 -u go2rtc -P CHANGEME -t "go2rtc/Camera1/playonspeaker"  -m "siren.mp3"
```

## Example with Home Assistant

```yaml
        - service: mqtt.publish
          data:
            topic: go2rtc/Camera1/playonspeaker
            payload: "siren.mp3"
```