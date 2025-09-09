package drive

import "fmt"

// Service Google Drive数据处理接口
type Service struct {
	connector *DriveConnector
}

// NewService 创建新的Drive服务
func NewService(connector *DriveConnector) *Service {
	return &Service{
		connector: connector,
	}
}

// GetFiles 获取文件列表
func (s *Service) GetFiles(userID string, limit int64) ([]File, error) {
	files, err := s.connector.ListFiles(userID, limit)
	if err != nil {
		return nil, fmt.Errorf("获取文件列表失败: %v", err)
	}
	return files, nil
}

// GetFileDetail 获取文件详情
func (s *Service) GetFileDetail(userID string, fileID string) (*File, error) {
	file, err := s.connector.GetFile(userID, fileID)
	if err != nil {
		return nil, fmt.Errorf("获取文件详情失败: %v", err)
	}
	return file, nil
}

// GetRecentFiles 获取最近修改的文件
func (s *Service) GetRecentFiles(userID string, limit int64) ([]File, error) {
	return s.GetFiles(userID, limit)
}

// GetFilesByType 按类型获取文件
func (s *Service) GetFilesByType(userID string, mimeType string, limit int64) ([]File, error) {
	// 这里可以实现按类型筛选文件的逻辑
	// 暂时返回所有文件，后续可以扩展
	return s.GetFiles(userID, limit)
}