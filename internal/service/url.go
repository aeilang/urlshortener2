package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/aeilang/urlshortener/config"
	"github.com/aeilang/urlshortener/internal/models"
	"github.com/aeilang/urlshortener/internal/repo"
)

type ShortCodeGenerator interface {
	GenerateID() string
}

type Cacher interface {
	SetURL(ctx context.Context, url repo.Url) error
	GetURLByShortCode(ctx context.Context, shortCode string) (*repo.Url, error)
}

type URLService struct {
	querier            repo.Querier
	shortCodeGenerator ShortCodeGenerator
	cfg                *config.Config
	cache              Cacher
	db                 *sql.DB
}

func NewURLService(db *sql.DB, shortCodeGenerator ShortCodeGenerator, cfg *config.Config, cache Cacher) *URLService {
	return &URLService{
		querier:            repo.New(db),
		shortCodeGenerator: shortCodeGenerator,
		cfg:                cfg,
		cache:              cache,
		db:                 db,
	}
}

func (s *URLService) CreateURL(ctx context.Context, req models.CreateURLRequest) (*models.CreateURLResponse, error) {
	// 判断custom_code是不是提供了
	var (
		shortCode  string
		expired_at time.Time
		is_custom  bool
		err        error
	)

	if req.CustomCode != "" {
		// 别名存在
		isAvialble, err := s.querier.IsShortCodeAvaliable(ctx, req.CustomCode)
		if err != nil {
			return nil, err
		}

		if !isAvialble {
			return nil, errors.New("别名已存在")
		}

		shortCode = req.CustomCode
		is_custom = true
	} else {
		// 别名不存在, 需要生成别名
		shortCode, err = s.createShortCode(ctx, 0)
		if err != nil {
			return nil, err
		}
	}

	if req.Duration == nil {
		// duration 不存在
		expired_at = time.Now().Add(s.cfg.App.DefaultExpiration)
	} else {
		expired_at = time.Now().Add(time.Hour * time.Duration(*req.Duration))
	}

	// 插入数据
	url, err := s.querier.CreateURL(ctx, repo.CreateURLParams{
		OrignalUrl: req.OrignalURL,
		ShortCode:  shortCode,
		IsCustom:   is_custom,
		ExpiredAt:  expired_at,
	})

	if err != nil {
		return nil, err
	}

	// 把URL缓存
	if err := s.cache.SetURL(ctx, url); err != nil {
		return nil, err
	}

	return &models.CreateURLResponse{
		ShortURL:  s.cfg.App.BaseURL + "/" + url.ShortCode,
		ExpiredAt: url.ExpiredAt,
	}, nil
}

func (s *URLService) GetOrignalURL(ctx context.Context, shortCode string) (string, error) {
	// 首先从缓存里取
	url, err := s.cache.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if url != nil {
		return url.OrignalUrl, nil
	}

	// 缓存里没有数据，从postgres中取
	url2, err := s.querier.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return "", err
	}
	if url2.OrignalUrl == "" {
		return "", errors.New("url 为空")
	}

	// 把取出的url2存入缓存
	if err := s.cache.SetURL(ctx, url2); err != nil {
		return "", err
	}

	return url2.OrignalUrl, nil
}

func (s *URLService) Cleanup(ctx context.Context) error {
	return s.querier.DeleteExpiredURLs(ctx)
}

func (s *URLService) createShortCode(ctx context.Context, n int) (string, error) {
	if n > 4 {
		return "", errors.New("shortCode 生成重试次数耗尽")
	}

	shortCode := s.shortCodeGenerator.GenerateID()
	log.Println(shortCode)
	isAvaliable, err := s.querier.IsShortCodeAvaliable(ctx, shortCode)
	if err != nil {
		return "", err
	}

	if !isAvaliable {
		return s.createShortCode(ctx, n+1)
	}

	return shortCode, nil
}
