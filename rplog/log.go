package rplog

import "time"

type Log struct {
	Id            string
	LogTime       time.Time
	LogMsg        string
	BinaryContent BinaryContent `bson:"binary_content"`
	TestItemRef   string        `bson:"last_modified"`
	LastModified  time.Time
	Level         LogLevel
}

type BinaryContent struct {
	Id          string
	ThumbnailId string `bson:"thumbnail_id"`
	ContentType string `bson:"content_type"`
}

type LogLevel struct {
	Level int `bson:"log_level"`
}

func main() {
}
