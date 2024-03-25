package interfaces

import "video/pkg/models"

type VideoRepo interface {
	CreateVideoId(models.Video) (bool, error)
	FetchAllVideos() ([]*models.Video, error)
	FetchUserVideos(int) ([]*models.Video, error)
	GetVideoById( string) (*models.Video,error)
	ArchivedVideos( string) (bool,error)
	FindArchivedVideos( int) ([]*models.Video,error)
}
