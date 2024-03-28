package models

import "time"

type Video struct {
	ID          uint `gorm:"primarykey"`
	Video_id    string
	S3_path     string
	Title       string
	Discription string
	Views       uint
	UserId      int
	UserName    string
	Blocked     bool `gorm:"default:false"`
	Archived    bool `gorm:"default:false"`
	Category    string
}

type Viewer struct {
	ID        uint `gorm:"primarykey"`
	VideoID   string
	UserId    int
	Timestamp time.Time
}

type Clip struct {
	ID       uint `gorm:"primarykey"`
	Clip_id  string
	S3_path  string
	Title    string
	Views    uint
	UserId   int
	UserName string
	Blocked  bool `gorm:"default:false"`
	Archived bool `gorm:"default:false"`
	Category string
}

type ClipViewer struct {
	ID        uint `gorm:"primarykey"`
	ClipId   string
	UserId    int
	Timestamp time.Time
}