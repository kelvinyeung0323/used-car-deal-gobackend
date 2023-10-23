package security

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"time"
	"used-car-deal-gobackend/base/logger"
	"used-car-deal-gobackend/base/web"
	"used-car-deal-gobackend/model"
)

// TODO:缓存登录，权限角色等信息
//var authCache = cache.New(4*time.Minute, 10*time.Minute)

type Authorization struct {
	*model.User
	RoleAuthCodes []string        `json:"roleAuthCodes"`
	Res           []model.Res     `json:"res"`
	Permissions   map[string]bool `json:"permissions"`
}

var log = logger.GetInstance()

type AuthorizationService interface {
	GetAuthorization(ctx *gin.Context) *Authorization
	GetAuthKey(ctx *gin.Context) string
}

type AuthHelper struct {
	authCache   *cache.Cache
	authService AuthorizationService
}

func (h *AuthHelper) GetAuthorization(ctx *gin.Context) (*Authorization, error) {
	//获取userId
	authKey := h.authService.GetAuthKey(ctx)
	if authKey == "" {
		return nil, fmt.Errorf("获取权限ID")
	}
	//从cache中获取权限信息
	var auth *Authorization
	authInCache, ok := h.authCache.Get(authKey)
	if !ok {
		auth = h.authService.GetAuthorization(ctx)
		if auth == nil {
			return nil, fmt.Errorf("获取不到权限信息")
		}
		h.authCache.Set(authKey, auth, cache.DefaultExpiration)
		return auth, nil
	} else {
		auth, ok = authInCache.(*Authorization)
		if ok {
			return auth, nil
		} else {
			return nil, fmt.Errorf("缓存转换错误")
		}
	}
}
func (h *AuthHelper) RefreshAuthorization(ctx *gin.Context) (*Authorization, error) {
	//获取userId
	authKey := h.authService.GetAuthKey(ctx)
	if authKey == "" {
		return nil, fmt.Errorf("获取权限ID")
	}
	auth := h.authService.GetAuthorization(ctx)
	if auth == nil {
		return nil, fmt.Errorf("获取不到权限信息")
	}
	h.authCache.Set(authKey, auth, cache.DefaultExpiration)
	return auth, nil

}

func NewAuthHelper(authService AuthorizationService) *AuthHelper {
	helper := &AuthHelper{authService: authService}
	authCache := cache.New(4*time.Minute, 10*time.Minute)
	helper.authCache = authCache
	return helper
}
func (h *AuthHelper) HasPermission(perms string) gin.HandlerFunc {

	return func(ctx *gin.Context) {

		auth, err := h.GetAuthorization(ctx)
		if err != nil {
			log.Warnf("authorize error:%v", err)
			web.Err(web.FORBIDDEN)
			ctx.Abort()
			return
		}
		if auth.UserName == "admin" {
			//如果是超级管理员
			ctx.Next()
			return
		}

		if !auth.Permissions[perms] {
			web.Err(web.FORBIDDEN)
			ctx.Abort()
		}

		ctx.Next()
	}
}
