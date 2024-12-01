// service/model/GameInfo.go
// KifuResponseに使用する対局情報

package model

import (
	"fmt"
	"strings"
)

// ------------------------------------------------------------
type TimeRuleString string

func (t TimeRuleString) ToGameInfo() GameInfo {
	gameInfo := GameInfo{}
	timeRules := strings.Split(string(t), "+")
	if len(timeRules) > 0 && timeRules[0] != "0" {
		gameInfo["持ち時間"] = fmt.Sprintf("%s秒", timeRules[0])
	}
	if len(timeRules) > 1 && timeRules[1] != "0" {
		gameInfo["秒読み"] = fmt.Sprintf("%s秒", timeRules[1])
	}
	if len(timeRules) > 2 && timeRules[2] != "0" {
		gameInfo["秒加算"] = fmt.Sprintf("%s秒", timeRules[2])
	}

	return gameInfo
}

// ------------------------------------------------------------
type GameInfo map[string]string

func (t GameInfo) Merge(gameInfo GameInfo) {
	for k, v := range gameInfo {
		t[k] = v
	}
}

func (t *Kifu) buildSummaryGameInfo() GameInfo {
	gameInfo := GameInfo{}
	if t.BlackPlayer != nil {
		gameInfo["先手"] = *t.BlackPlayer
	}
	if t.WhitePlayer != nil {
		gameInfo["後手"] = *t.WhitePlayer
	}
	if t.StartedAt != nil {
		gameInfo["対局日時"] = t.StartedAt.Format("2006-01-02 03:04")
	}
	return gameInfo
}

func (t *Kifu) buildGameInfo(options []*KifuOption) map[string]string {
	gameInfo := t.buildSummaryGameInfo()
	if t.TimeRule != nil {
		gameInfo.Merge(t.TimeRule.ToGameInfo())
	}
	for _, option := range options {
		gameInfo[option.Name] = option.Value
	}
	return gameInfo
}

// ------------------------------------------------------------

func (t *Kifu) buildTags(kifuTags []*KifuTag) []string {
	tags := []string{}
	for _, kifuTag := range kifuTags {
		tags = append(tags, kifuTag.Name)
	}
	return tags
}
