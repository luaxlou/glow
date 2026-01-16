package manager

import (
	"github.com/luaxlou/glow/internal/configmanager"
	"github.com/luaxlou/glow/pkg/api"
)

func ListResources() ([]api.ResourceRef, error) {
	var resources []api.ResourceRef

	// Check MySQL
	var mysqlConfig api.MySQLConfig
	if err := configmanager.GetSystemConfigJSON("mysql_info", &mysqlConfig); err == nil && mysqlConfig.Host != "" {
		resources = append(resources, api.ResourceRef{
			Kind: "mysql",
			Name: "mysql-local",
			Port: mysqlConfig.Port,
		})
	}

	// Check Redis
	var redisConfig api.RedisConfig
	if err := configmanager.GetSystemConfigJSON("redis_info", &redisConfig); err == nil && redisConfig.Host != "" {
		resources = append(resources, api.ResourceRef{
			Kind: "redis",
			Name: "redis-local",
			Port: redisConfig.Port,
		})
	}

	// Check Nginx
	var nginxConfig api.NginxSystemConfig
	if err := configmanager.GetSystemConfigJSON("nginx_info", &nginxConfig); err == nil && nginxConfig.BinaryPath != "" {
		resources = append(resources, api.ResourceRef{
			Kind: "nginx",
			Name: "nginx-system",
		})
	}

	return resources, nil
}
