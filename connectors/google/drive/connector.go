package drive

import (
	"context"
	"fmt"
	"time"

	"connector-demo/utils"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// DriveConnector Google Drive API封装
type DriveConnector struct {
	tokenManager *utils.TokenManager
}

// Owner 文件所有者
type Owner struct {
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

// Permission 文件权限
type Permission struct {
	ID           string `json:"id"`
	Type         string `json:"type"`         // user, group, domain, anyone
	Role         string `json:"role"`         // owner, writer, reader
	EmailAddress string `json:"emailAddress"` // 仅对 user/group 类型有效
}

// File Google Drive 文件结构（扩展版）
type File struct {
	// 基本标识与内容来源
	ID             string            `json:"id"`             // 文件唯一ID
	Name           string            `json:"name"`           // 文件名
	MimeType       string            `json:"mimeType"`       // 文件类型
	WebViewLink    string            `json:"webViewLink"`    // 浏览器访问链接
	WebContentLink string            `json:"webContentLink"` // 文件下载链接
	ExportLinks    map[string]string `json:"exportLinks"`    // Docs/Sheets 可导出链接

	// 内容和文本
	IndexableText string `json:"indexableText"` // contentHints.indexableText

	// 时间和版本信息
	CreatedTime  time.Time `json:"createdTime"`  // 文件创建时间
	ModifiedTime time.Time `json:"modifiedTime"` // 文件最后修改时间
	Version      int64     `json:"version"`      // 文件版本号

	// 权限和共享信息
	Owners          []Owner      `json:"owners"`          // 文件所有者列表
	Shared          bool         `json:"shared"`          // 是否共享
	Permissions     []Permission `json:"permissions"`     // 具体权限
	ViewedByMe      bool         `json:"viewedByMe"`      // 是否查看过
	WritersCanShare bool         `json:"writersCanShare"` // 是否可分享

	// 文件元数据
	Parents           []string `json:"parents"`           // 父目录ID列表
	Size              int64    `json:"size"`              // 文件大小（字节）
	Description       string   `json:"description"`       // 文件描述
	ThumbnailLink     string   `json:"thumbnailLink"`     // 缩略图
	OriginalFilename  string   `json:"originalFilename"`  // 原始文件名
	FullFileExtension string   `json:"fullFileExtension"` // 完整扩展名
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

	// 获取文件列表，扩展字段
	filesResp, err := service.Files.List().
		PageSize(maxResults).
		OrderBy("modifiedTime desc").
		Fields("files(id, name, mimeType, createdTime, modifiedTime, size, webViewLink, webContentLink, thumbnailLink, parents, owners, exportLinks, contentHints/indexableText, description, fullFileExtension, version, shared, viewedByMe, writersCanShare, permissions)").
		Do()
	if err != nil {
		return nil, fmt.Errorf("获取文件列表失败: %v", err)
	}

	var result []File
	for _, f := range filesResp.Files {
		file := mapDriveFileToFile(f)
		result = append(result, *file)
	}

	return result, nil
}

// GetFile 获取单个文件详情
func (dc *DriveConnector) GetFile(userID string, fileID string) (*File, error) {
	service, err := dc.GetService(userID)
	if err != nil {
		return nil, err
	}

	f, err := service.Files.Get(fileID).
		Fields("id, name, mimeType, createdTime, modifiedTime, size, webViewLink, webContentLink, thumbnailLink, parents, owners, exportLinks, contentHints/indexableText, description, fullFileExtension, version, shared, viewedByMe, writersCanShare, permissions").
		Do()
	if err != nil {
		return nil, fmt.Errorf("获取文件详情失败: %v", err)
	}

	file := mapDriveFileToFile(f)
	return file, nil
}

// mapDriveFileToFile 将 drive.File 转换为 File
func mapDriveFileToFile(f *drive.File) *File {
	var owners []Owner
	for _, o := range f.Owners {
		owners = append(owners, Owner{
			DisplayName:  o.DisplayName,
			EmailAddress: o.EmailAddress,
		})
	}

	var permissions []Permission
	for _, p := range f.Permissions {
		permissions = append(permissions, Permission{
			ID:           p.Id,
			Type:         p.Type,
			Role:         p.Role,
			EmailAddress: p.EmailAddress,
		})
	}

	createdTime, _ := time.Parse(time.RFC3339, f.CreatedTime)
	modifiedTime, _ := time.Parse(time.RFC3339, f.ModifiedTime)

	return &File{
		ID:             f.Id,
		Name:           f.Name,
		MimeType:       f.MimeType,
		WebViewLink:    f.WebViewLink,
		WebContentLink: f.WebContentLink,
		ExportLinks:    f.ExportLinks,
		IndexableText: func() string {
			if f.ContentHints != nil {
				return f.ContentHints.IndexableText
			}
			return ""
		}(),
		CreatedTime:       createdTime,
		ModifiedTime:      modifiedTime,
		Version:           f.Version,
		Owners:            owners,
		Shared:            f.Shared,
		ViewedByMe:        f.ViewedByMe,
		WritersCanShare:   f.WritersCanShare,
		Permissions:       permissions,
		Parents:           f.Parents,
		Size:              f.Size,
		Description:       f.Description,
		ThumbnailLink:     f.ThumbnailLink,
		OriginalFilename:  f.OriginalFilename,
		FullFileExtension: f.FullFileExtension,
	}
}
