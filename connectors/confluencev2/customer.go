// Package confluencev2
package confluencev2

// GenerateAuthorizationURL 生成授权链接
func GenerateAuthorizationURL(state string) string {
	baseAuthURL := "xxxxx" // 替换为实际的授权URL
	return baseAuthURL + "&state=" + state
}
