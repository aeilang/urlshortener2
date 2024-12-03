package api

import (
	"context"
	"net/http"

	"github.com/aeilang/urlshortener/internal/models"
	"github.com/labstack/echo/v4"
)

type URLService interface {
	CreateURL(ctx context.Context, req models.CreateURLRequest) (*models.CreateURLResponse, error)

	GetOrignalURL(ctx context.Context, shortCode string) (string, error)
}

type URLHandler struct {
	urlService URLService
}

func NewURLHandler(urlService URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

// POST /api/url 接收orignal_url, custom_code, duration 返回short_url
func (h *URLHandler) CreateURL(c echo.Context) error {
	// 从request body 里提取出CreateURLReqeust
	var req models.CreateURLRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// 验证这个CreateURLRequest的数据格式
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 调用业务函数, 业务函数范围short_url
	res, err := h.urlService.CreateURL(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// 把short_url作为响应进行返回
	return c.JSON(http.StatusCreated, res)
}

// GET /:code 重定向到orignail_url
func (h *URLHandler) RedirectURL(c echo.Context) error {
	// 取出路径参数
	shortCode := c.Param("code")
	//
	orignalURL, err := h.urlService.GetOrignalURL(c.Request().Context(), shortCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.Redirect(http.StatusPermanentRedirect, orignalURL)
}
