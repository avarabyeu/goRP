package db

import (
	"time"
	"gopkg.in/mgo.v2"
)

//Log represents log collection in mongodb
type Log struct {
	ID            string `bson:"id"`
	LogTime       time.Time `bson:"logTime"`
	LogMsg        string `bson:"logMsg"`
	BinaryContent BinaryContent `bson:"binary_content"`
	TestItemRef   string    `bson:"testItemRef"`
	LastModified  time.Time `bson:"last_modified"`
	Level         LogLevel `bson:"level"`
}

//BinaryContent represents Binary Content sub-documer in RP mongodb
type BinaryContent struct {
	ID          string
	ThumbnailID string `bson:"thumbnail_id"`
	ContentType string `bson:"content_type"`
}

//LogLevel is int representation of Log's level
type LogLevel struct {
	Level int `bson:"log_level"`
}

//LogRepo is Repository pattern implementation
type LogRepo struct {
	session *mgo.Session
	db      string
}

//NewLogRepo creates new instance of LogRepo object
func NewLogRepo(s *mgo.Session, db string) *LogRepo {
	return &LogRepo{
		session: s,
		db:      db,

	}
}

//FindAll returns all entities of collection
func (r *LogRepo) FindAll() ([]*Log, error) {
	s := r.session.Clone()
	defer s.Close()

	var logs []*Log
	e := s.DB(r.db).C("log").Find(nil).Limit(10).All(&logs)
	return logs, e
}

//Save saves entity to collection
func (r *LogRepo) Save(log *Log) error {
	s := r.session.Clone()
	defer s.Close()
	return s.DB(r.db).C("log").Insert(log)
}
