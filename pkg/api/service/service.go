package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	"video/pkg/models"
	"video/pkg/pb/video"

	"video/pkg/repo/interfaces"
	"video/pkg/utils"

	"github.com/google/uuid"
)

type VideoServer struct {
	Repo interfaces.VideoRepo
	video.VideoServiceServer
}

func NewVideoServer(repo interfaces.VideoRepo) video.VideoServiceServer {
	return &VideoServer{
		Repo: repo,
	}
}

func (c *VideoServer) UploadVideo(stream video.VideoService_UploadVideoServer) error {
	var req models.Video
	var buffer bytes.Buffer

	fileUID := uuid.New()
	fileName := fileUID.String()
	s3Path := "streamers/" + fileName + ".mp4"

	for {
		uploadData, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		videoId := utils.GenerateUniqueString()

		req = models.Video{
			Title:       uploadData.Title,
			Discription: uploadData.Description,
			Category:    uploadData.Category,
			UserId:      int(uploadData.UserId),
			Video_id:    videoId,
		}
		_, err = buffer.Write(uploadData.Data)
		if err != nil {
			return err
		}
	}
	err := utils.UploadVideoToS3(buffer.Bytes(), s3Path)
	if err != nil {
		return err
	}
	req.S3_path = s3Path

	_, err = c.Repo.CreateVideoId(req)
	if err != nil {
		return err
	}
	return stream.SendAndClose(&video.UploadVideoResponse{
		Status:  http.StatusOK,
		Message: "Video succesfully uploaded",
		VideoId: "",
	})

}

func (c *VideoServer) FindUserVideo(ctx context.Context, input *video.FindUserVideoRequest) (*video.FindUserVideoResponse, error) {
	res, err := c.Repo.FetchUserVideos(int(input.Userid))
	if err != nil {
		return nil, err
	}
	data := make([]*video.FetchVideo, len(res))
	for i, v := range res {
		data[i] = &video.FetchVideo{
			VideoId:     v.Video_id,
			S3Path:      v.S3_path,
			OwnerId:     int32(v.UserId),
			Title:       v.Title,
			Discription: v.Discription,
			Views:       uint32(v.Views),
			Archived:    v.Archived,
			Blocked:     v.Blocked,
		}
	}
	resp := &video.FindUserVideoResponse{
		Videos: data,
	}
	return resp, err
}

func (c *VideoServer) FindAllVideo(ctx context.Context, input *video.FindAllVideoRequest) (*video.FindAllVideoResponse, error) {
	res, err := c.Repo.FetchAllVideos()
	if err != nil {
		return nil, err
	}

	data := make([]*video.FetchVideo, len(res))
	for i, v := range res {
		data[i] = &video.FetchVideo{
			VideoId:     v.Video_id,
			S3Path:      v.S3_path,
			OwnerId:     int32(v.UserId),
			Title:       v.Title,
			Discription: v.Discription,
			Category:    v.Category,
			Views:       uint32(v.Views),
			Archived:    v.Archived,
			Blocked:     v.Blocked,
		}
	}
	response := &video.FindAllVideoResponse{
		Videos: data,
	}
	return response, err
}

func (c *VideoServer) FindArchivedVideoByUserId(ctx context.Context, input *video.FindArchivedVideoByUserIdRequest) (*video.FindArchivedVideoByUserIdResponse, error) {
	res, err := c.Repo.FindArchivedVideos(int(input.Userid))
	if err != nil {
		return nil, err
	}
	data := make([]*video.FetchVideo, len(res))
	for i, v := range res {
		data[i] = &video.FetchVideo{
			VideoId:     v.Video_id,
			S3Path:      v.S3_path,
			Title:       v.Title,
			Discription: v.Discription,
			Archived:    v.Archived,
			Views:       uint32(v.Views),
			Blocked:     v.Blocked,
			Category:    v.Category,
			OwnerId:     int32(v.UserId),
		}
	}
	response := &video.FindArchivedVideoByUserIdResponse{
		Videos: data,
	}
	return response, err

}

func (c *VideoServer) ArchiveVideo(ctx context.Context, input *video.ArchiveVideoRequest) (*video.ArchiveVideoResponse, error) {
	fmt.Println("vide", input.VideoId)
	res, err := c.Repo.ArchivedVideos(input.VideoId)
	if err != nil {
		return nil, err
	}
	response := &video.ArchiveVideoResponse{
		Status: res,
	}
	return response, err
}

func (c *VideoServer) GetVideoById(ctx context.Context, input *video.GetVideoByIdRequest) (*video.GetVideoByIdResponse, error) {
	res, err := c.Repo.GetVideoById(input.VideoID)
	if err != nil {
		return nil, err
	}
	response := &video.GetVideoByIdResponse{
		VideoId:     res.Video_id,
		UserName:    res.UserName,
		Archived:    res.Archived,
		Blocked:     res.Blocked,
		Title:       res.Title,
		S3Path:      res.S3_path,
		Discription: res.Discription,
		Views:       uint32(res.Views),
		Category:    res.Category,
	}
	return response, nil
}

func (c *VideoServer) UploadClip(stream video.VideoService_UploadClipServer) error {
	var req models.Clip
	var buffer bytes.Buffer

	fileUID := uuid.New()
	fileName := fileUID.String()
	s3Path := "streamers_clip/" + fileName + ".mp4"
	// var totalDuration time.Duration

	for {
		uploadData, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		clipId := utils.GenerateUniqueString()

		req = models.Clip{
			Title:    uploadData.Title,
			Category: uploadData.Category,
			UserId:   int(uploadData.UserId),
			Clip_id:  clipId,
		}
		_, err = buffer.Write(uploadData.Data)
		if err != nil {
			return err
		}

		// chunkDuration,err:=utils.GetVideoDuration(uploadData.Data)
		// if err != nil{
		// 	return err
		// }
		// totalDuration += chunkDuration

		// if totalDuration > time.Minute{
		// 	return errors.New("video duration exceeds 1 minute")
		// }
	}
	chunkDuration, err := utils.GetVideoDuration(buffer.Bytes())
	if err != nil {
		return err
	}
	fmt.Println("Duration",chunkDuration)

	// Optionally, validate duration
	if chunkDuration > time.Minute {
		return errors.New("video duration exceeds 1 minute")
	}

	err1 := utils.UploadVideoToS3(buffer.Bytes(), s3Path)
	if err1 != nil {
		return err
	}
	req.S3_path = s3Path

	_, err = c.Repo.CreateClipId(req)
	if err != nil {
		return err
	}
	return stream.SendAndClose(&video.UploadClipResponse{
		Status:  http.StatusOK,
		Message: "Video succesfully uploaded",
		VideoId: req.Clip_id,
	})

}

func (c *VideoServer) FindUserClip(ctx context.Context, input *video.FindUserClipRequest) (*video.FindUserClipResponse, error) {
	res, err := c.Repo.FetchUserClips(int(input.Userid))
	if err != nil {
		return nil, err
	}
	data := make([]*video.FetchClip, len(res))
	for i, v := range res {
		data[i] = &video.FetchClip{
			ClipId:   v.Clip_id,
			S3Path:   v.S3_path,
			OwnerId:  int32(v.UserId),
			Title:    v.Title,
			Views:    uint32(v.Views),
			Archived: v.Archived,
			Blocked:  v.Blocked,
		}
	}
	resp := &video.FindUserClipResponse{
		Clips: data,
	}
	return resp, err
}

func (c *VideoServer) FindAllClip(ctx context.Context, input *video.FindAllClipRequest) (*video.FindAllClipResponse, error) {
	res, err := c.Repo.FetchAllClips()
	if err != nil {
		return nil, err
	}

	data := make([]*video.FetchClip, len(res))
	for i, v := range res {
		data[i] = &video.FetchClip{
			ClipId:   v.Clip_id,
			S3Path:   v.S3_path,
			OwnerId:  int32(v.UserId),
			Title:    v.Title,
			Category: v.Category,
			Views:    uint32(v.Views),
			Archived: v.Archived,
			Blocked:  v.Blocked,
		}
	}
	response := &video.FindAllClipResponse{
		Clips: data,
	}
	return response, err
}

func (c *VideoServer) FindArchivedClipByUserId(ctx context.Context, input *video.FindArchivedClipByUserIdRequest) (*video.FindArchivedClipByUserIdResponse, error) {
	res, err := c.Repo.FindArchivedClips(int(input.Userid))
	if err != nil {
		return nil, err
	}
	data := make([]*video.FetchClip, len(res))
	for i, v := range res {
		data[i] = &video.FetchClip{
			ClipId:   v.Clip_id,
			S3Path:   v.S3_path,
			Title:    v.Title,
			Archived: v.Archived,
			Views:    uint32(v.Views),
			Blocked:  v.Blocked,
			Category: v.Category,
			OwnerId:  int32(v.UserId),
		}
	}
	response := &video.FindArchivedClipByUserIdResponse{
		Clips: data,
	}
	return response, err

}

func (c *VideoServer) ArchiveClip(ctx context.Context, input *video.ArchiveClipRequest) (*video.ArchiveClipResponse, error) {

	res, err := c.Repo.ArchivedClip(input.ClipId)
	if err != nil {
		return nil, err
	}
	response := &video.ArchiveClipResponse{
		Status: res,
	}
	return response, err
}

func (c *VideoServer) GetClipById(ctx context.Context, input *video.GetClipByIdRequest) (*video.GetClipByIdResponse, error) {
	res, err := c.Repo.GetClipById(input.ClipId)
	if err != nil {
		return nil, err
	}
	response := &video.GetClipByIdResponse{
		ClipId:   res.Clip_id,
		UserName: res.UserName,
		Archived: res.Archived,
		Blocked:  res.Blocked,
		Title:    res.Title,
		S3Path:   res.S3_path,
		Views:    uint32(res.Views),
		Category: res.Category,
	}
	return response, nil
}
