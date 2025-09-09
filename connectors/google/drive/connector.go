package drive

import (
	"context"
	"fmt"

	"connector-demo/utils"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// DriveConnector Google Drive API封装
type DriveConnector struct {
	tokenManager *utils.TokenManager
}

// File Google Drive文件信息
type File struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	MimeType      string `json:"mimeType"`
	ModifiedTime  string `json:"modifiedTime"`
	Size          int64  `json:"size,omitempty"`
	WebViewLink   string `json:"webViewLink,omitempty"`
	ThumbnailLink string `json:"thumbnailLink,omitempty"`
	Parents       []string `json:"parents,omitempty"`
}

// NewDriveConnector 创建新的Drive连接器
func NewDriveConnector(tm *utils.TokenManager) *DriveConnector {
	return &DriveConnector{
		tokenManager: tm,
	}
}

// GetService 获取Drive服务客户端
func (dc *DriveConnector) GetService(userID string) (*drive.Service, error) {
	tokenInfo, exists := dc.tokenManager.GetToken(userID, "google")
	if !exists {
		return nil, fmt.Errorf("未找到Google访问令牌")
	}

	// 创建OAuth2客户端
	client := utils.CreateOAuth2Client(tokenInfo.AccessToken)
	
	// 创建Drive服务
	service, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("创建Drive服务失败: %v", err)
	}

	return service, nil
}

// ListFiles 获取文件列表
func (dc *DriveConnector) ListFiles(userID string, maxResults int64) ([]File, error) {
	service, err := dc.GetService(userID)
	if err != nil {
		return nil, err
	}

	// 获取文件列表
	files, err := service.Files.List().
		PageSize(maxResults).
		OrderBy("modifiedTime desc").
		Fields("files(id, name, mimeType, modifiedTime, size, webViewLink, thumbnailLink, parents)").
		Do()
	if err != nil {
		return nil, fmt.Errorf("获取文件列表失败: %v", err)
	}

	var result []File
	for _, file := range files.Files {
		result = append(result, File{
			ID:            file.Id,
			Name:          file.Name,
			MimeType:      file.MimeType,
			ModifiedTime:  file.ModifiedTime,
			Size:          file.Size,
			WebViewLink:   file.WebViewLink,
			ThumbnailLink: file.ThumbnailLink,
			Parents:       file.Parents,
		})
	}

	return result, nil
}

// GetFile 获取单个文件详情
func (dc *DriveConnector) GetFile(userID string, fileID string) (*File, error) {
	service, err := dc.GetService(userID)
	if err != nil {
		return nil, err
	}

	file, err := service.Files.Get(fileID).
		Fields("id, name, mimeType, modifiedTime, size, webViewLink, thumbnailLink, parents").
		Do()
	if err != nil {
		return nil, fmt.Errorf("获取文件详情失败: %v", err)
	}

	return &File{
		ID:            file.Id,
		Name:          file.Name,
		MimeType:      file.MimeType,
		ModifiedTime:  file.ModifiedTime,
		Size:          file.Size,
		WebViewLink:   file.WebViewLink,
		ThumbnailLink: file.ThumbnailLink,
		Parents:       file.Parents,
	}, nil
}