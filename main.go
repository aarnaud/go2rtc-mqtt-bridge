package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	mqtt_client "go2rtc-mqtt-bridge/mqtt-client"
	"go2rtc-mqtt-bridge/utils"
	"io"
	"net/http"
	"path"
	"regexp"
	"time"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	config := utils.GetConfig()
	hostRegex := regexp.MustCompile(fmt.Sprintf(`^%s/(.+?)/`, config.MQTT.BaseTopic))
	cli := mqtt_client.NewMQTT(config)

	go func() {
		for {
			cli.PublishAvailability()
			time.Sleep(time.Second * 30)
		}
	}()
	cli.WatchTopicPlayOnSpeaker(func(client mqtt.Client, message mqtt.Message) {
		matches := hostRegex.FindStringSubmatch(message.Topic())
		if matches == nil && len(matches) < 2 {
			log.Info().Msgf("failed to extract camera host from topic: %s", message.Topic())
			return
		}
		cameraName := matches[1]
		audioFilePath := path.Join(config.AudioFilesPath, string(message.Payload()))
		srcArgs := fmt.Sprintf("ffmpeg:%s#audio=pcma#input=file", audioFilePath)

		go2rtcUrl := config.Go2rtcURL.JoinPath("api/streams")

		queryValues := go2rtcUrl.Query()

		queryValues.Add("dst", cameraName)
		queryValues.Add("src", srcArgs)

		httpCli := &http.Client{
			Timeout: time.Second * 30,
		}
		req, err := http.NewRequest(http.MethodPost, go2rtcUrl.String(), nil)
		if config.Go2rtcUsername != "" && config.Go2rtcPassword != "" {
			req.SetBasicAuth(config.Go2rtcUsername, config.Go2rtcPassword)
		}
		req.URL.RawQuery = queryValues.Encode()
		response, err := httpCli.Do(req)
		if err != nil {
			log.Error().Msgf("failed to send %s to %s with error %s", message.Payload(), cameraName, err)
			return
		}
		if response.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(response.Body)
			log.Error().Msgf("failed to send %s to %s with status %d: %s",
				message.Payload(), cameraName, response.StatusCode, body)
			return
		}
		log.Info().Msgf("play %s on %s", message.Payload(), cameraName)
	})
}
