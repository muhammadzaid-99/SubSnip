// internal/models/task.go
package models

import "encoding/json"

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type SubtitleBox struct {
	TopLeft     Point `json:"top_left"`
	BottomRight Point `json:"bottom_right"`
}

type TaskRequest struct {
	TaskID       string      `json:"task_id"`
	VideoURL     string      `json:"video_url"`
	LanguageCode string      `json:"language_code"`
	StartTime    string      `json:"start_time"`
	EndTime      string      `json:"end_time"`
	FramesToSkip int         `json:"frames_to_skip"`
	SubtitleBox  SubtitleBox `json:"subtitle_box"`
}

func (t TaskRequest) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}
