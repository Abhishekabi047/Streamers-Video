package repo

import (
	"fmt"
	"log"
	"time"
	"video/pkg/models"
	"video/pkg/repo/interfaces"

	"gorm.io/gorm"
)

type VideoRepo struct{
	DB *gorm.DB
}

func NewVideoRepo(db *gorm.DB) interfaces.VideoRepo{
	return &VideoRepo{
		DB: db,
	}
}

func (c *VideoRepo) CreateVideoId(input models.Video)(bool,error) {
	video:= &models.Video{
		S3_path: input.S3_path,
		Title: input.Title,
		Discription: input.Discription,
		Category: input.Category,
		Views: 0,
		UserId: input.UserId,
		Video_id: input.Video_id,
	}

	if err :=c.DB.Create(video).Error;err != nil{
		return false,err
	}
	return true,nil
}

func(c *VideoRepo) FetchUserVideos(userid int) ([]*models.Video,error) {
	var data []*models.Video
	if err := c.DB.Where("user_id=? AND blocked=?",userid,false).Find(&data).Error;err != nil{
		return nil,err
	}
	if len(data) == 0{
		log.Println("fetching empty array")
		return []*models.Video{},nil
	}
	return data,nil

}

func(c *VideoRepo) FetchAllVideos() ([]*models.Video,error) {
	var data []*models.Video

	if err :=c.DB.Model(&models.Video{}).Where("archived = ? AND blocked = ?",false,false).Find(&data).Error;err != nil{
		return nil,err
	}
	if len(data) == 0{
		log.Println("fetching empty array")
		return []*models.Video{},nil
	}
	return data,nil
}

func (c *VideoRepo) FindArchivedVideos(userId int) ([]*models.Video,error) {
	var data []*models.Video

	if err :=c.DB.Model(&models.Video{}).Where("user_id=? AND archived=?",userId,true).Find(&data).Error;err!= nil{
		return nil,err
	}
	if len(data) == 0{
		log.Println("fetching empty array")
		return []*models.Video{},nil
	}
	return data, nil
}

func (c *VideoRepo) ArchivedVideos(videoId string) (bool,error) {
	var video models.Video
	if err :=c.DB.Where("video_id=?",videoId).First(&video).Error;err != nil{
		return false,nil
	}
	video.Archived=!video.Archived
	if err :=c.DB.Save(&video).Error;err != nil{
		return false,err
	}
	return true, nil
}

func (c *VideoRepo) GetVideoById(videoId string) (*models.Video,error) {
	var video models.Video
	if err:=c.DB.Where("video_id=?",videoId).First(&video).Error;err != nil{
		if err == gorm.ErrRecordNotFound{
			return nil,fmt.Errorf("Video not found")
		}
		return nil,err
	}
	video.Views++
	if err :=c.DB.Save(&video).Error;err != nil{
		log.Println("error in incrementing views")
		return nil,err
	}
	view:=models.Viewer{
		VideoID: videoId,
		UserId: video.UserId,
		Timestamp: time.Now(),
	}
	if err := c.DB.Create(&view).Error;err != nil{
		return nil,err
	}
	return &video,nil
}

func (c *VideoRepo) CreateClipId(input models.Clip)(bool,error) {
	clip:= &models.Clip{
		S3_path: input.S3_path,
		Title: input.Title,
		Category: input.Category,
		Views: 0,
		UserId: input.UserId,
		Clip_id: input.Clip_id,
	}

	if err :=c.DB.Create(clip).Error;err != nil{
		return false,err
	}
	return true,nil
}


func(c *VideoRepo) FetchUserClips(userid int) ([]*models.Clip,error) {
	var data []*models.Clip
	if err := c.DB.Where("user_id=? AND blocked=?",userid,false).Find(&data).Error;err != nil{
		return nil,err
	}
	if len(data) == 0{
		log.Println("fetching empty array")
		return []*models.Clip{},nil
	}
	return data,nil

}


func(c *VideoRepo) FetchAllClips() ([]*models.Clip,error) {
	var data []*models.Clip

	if err :=c.DB.Model(&models.Clip{}).Where("archived = ? AND blocked = ?",false,false).Find(&data).Error;err != nil{
		return nil,err
	}
	if len(data) == 0{
		log.Println("fetching empty array")
		return []*models.Clip{},nil
	}
	return data,nil
}


func (c *VideoRepo) FindArchivedClips(userId int) ([]*models.Clip,error) {
	var data []*models.Clip

	if err :=c.DB.Model(&models.Clip{}).Where("user_id=? AND archived=?",userId,true).Find(&data).Error;err!= nil{
		return nil,err
	}
	if len(data) == 0{
		log.Println("fetching empty array")
		return []*models.Clip{},nil
	}
	return data, nil
}

func (c *VideoRepo) ArchivedClip(clipId string) (bool,error) {
	var clip models.Clip
	if err :=c.DB.Where("clip_id=?",clipId).First(&clip).Error;err != nil{
		return false,nil
	}
	clip.Archived=!clip.Archived
	if err :=c.DB.Save(&clip).Error;err != nil{
		return false,err
	}
	return true, nil
}

func (c *VideoRepo) GetClipById(clipId string) (*models.Clip,error) {
	var clip models.Clip
	if err:=c.DB.Where("clip_id=?",clipId).First(&clip).Error;err != nil{
		if err == gorm.ErrRecordNotFound{
			return nil,fmt.Errorf("clip not found")
		}
		return nil,err
	}
	clip.Views++
	if err :=c.DB.Save(&clip).Error;err != nil{
		log.Println("error in incrementing views")
		return nil,err
	}
	view:=models.ClipViewer{
		ClipId: clipId,
		UserId: clip.UserId,
		Timestamp: time.Now(),
	}
	if err := c.DB.Create(&view).Error;err != nil{
		return nil,err
	}
	return &clip,nil
}