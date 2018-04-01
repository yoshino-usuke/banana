package helper

import (
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type BotData struct {
	SecretChannel string
	ChannelKey    int
}

const EndPoint = "https://trialbot-api.line.me"
const EventType = "138311608800106203"

func GetBotData(project string) (*BotData, error) {
	path := fmt.Sprintf(CONF_FILE, project)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var d BotData
	err = yaml.Unmarshal(file, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}
