/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package service

import (
	"encoding/json"
	"strconv"

	"github.com/Git-So/blog-api/utils/conf"
	"github.com/Git-So/blog-api/utils/helper"

	"github.com/Git-So/blog-api/models"
	"github.com/Git-So/blog-api/utils/cache"
	"github.com/jinzhu/gorm"
	"github.com/wonderivan/logger"
)

// GetArticleInfoByID 通过文章序号获取文章信息
func (s *Service) GetArticleInfoByID(id uint, isAdmin bool) (*models.Article, error) {
	var cacheArticleInfo models.Article
	key := cache.GetKey(`GetArticleInfoByID`, id, isAdmin)

	// 获取缓存
	if s.IsCache {
		data, stat, err := cache.GetCacheData(key)
		if err == nil && stat {
			// 数据解析
			jsonData, err := helper.Debase64(data)
			if err == nil {
				json.Unmarshal(jsonData, &cacheArticleInfo)
				if cacheArticleInfo.State != 1 && !isAdmin {
					return nil, gorm.ErrRecordNotFound
				}
				return &cacheArticleInfo, nil
			}
			logger.Warn("缓存数据有误,无法解析：", key, data)
		}
	}

	// 查询数据
	cacheArticleInfo.ID = id
	err := cacheArticleInfo.Info(isAdmin)
	if isErrDB(err) {
		return nil, err
	}

	// 保存缓存
	dataString, err := json.Marshal(&cacheArticleInfo)
	if err != nil {
		return nil, err
	}
	cache.Get().SetEx(key, conf.Get().Cache.Expired, helper.Enbase64(dataString))

	return &cacheArticleInfo, nil
}

// GetHotArticleList 获取热门文章列表
func (s *Service) GetHotArticleList(isAdmin bool, pageNum, pageSize uint) (articleList []*models.Article, err error) {
	key := cache.GetKey(`GetHotArticleList`, pageNum, pageSize, isAdmin)

	// 获取缓存
	if s.IsCache {
		data, stat, err := cache.GetCacheData(key)
		if err == nil && stat {
			// 数据解析
			jsonData, err := helper.Debase64(data)
			if err == nil {
				json.Unmarshal(jsonData, &articleList)
				return articleList, nil
			}
			logger.Warn("缓存数据有误,无法解析：", key, data)
		}
	}

	// 查询数据
	var cacheArticle models.Article
	articleList, err = cacheArticle.HotList(isAdmin, pageNum, pageSize)
	if isErrDB(err) {
		return nil, err
	}

	// 保存缓存
	dataString, err := json.Marshal(&articleList)
	if err != nil {
		return nil, err
	}
	cache.Get().SetEx(key, conf.Get().Cache.Expired, helper.Enbase64(dataString))

	return
}

// DeleteArticle 删除文章
func (s *Service) DeleteArticle(id uint) (err error) {
	var ArticleInfo models.Article
	key := cache.GetKey(`DeleteArticle`, id)

	// 删除文章
	ArticleInfo.ID = id
	err = ArticleInfo.Delete()
	if isErrDB(err) {
		return
	}

	// 删除缓存
	cache.Get().Del(key)

	return
}

// ArticleTotal .
func (s *Service) ArticleTotal(isAdmin bool, where []interface{}) (count uint, err error) {
	var cacheArticleInfo models.Article
	key := cache.GetKey(`ArticleTotal`, where, isAdmin)

	// 获取缓存
	if s.IsCache {
		data, stat, err := cache.GetCacheData(key)
		if err == nil && stat {
			// 数据解析
			count, err := strconv.Atoi(data)
			if err == nil {
				return uint(count), nil
			}
			logger.Warn("缓存数据有误,无法解析：", key, data)
		}
	}

	// 查询数据
	count, err = cacheArticleInfo.Total(isAdmin, where)
	if isErrDB(err) {
		return 0, err
	}

	// 保存缓存
	if count > 0 {
		cache.Get().SetEx(key, conf.Get().Cache.Expired, count)
	}

	return count, nil
}

// isExistsArticle 是否存在文章
func (s *Service) isExistsArticle(isAdmin bool, where ...interface{}) (IsExists bool, err error) {
	var count uint
	count, err = s.ArticleTotal(isAdmin, where)

	if count > 0 {
		IsExists = true
	}
	return
}

// IsExistsArticleByID 。
func (s *Service) IsExistsArticleByID(isAdmin bool, id uint) (IsExists bool, err error) {
	return s.isExistsArticle(isAdmin, "id = ?", id)
}

// IsExistsArticleByTitle 。
func (s *Service) IsExistsArticleByTitle(isAdmin bool, title string) (IsExists bool, err error) {
	return s.isExistsArticle(isAdmin, "title = ?", title)
}

// UpdateArticle 。
func (s *Service) UpdateArticle(article *models.Article, tag []string, subjectID uint) (err error) {
	return article.Update(tag, subjectID)
}

// CreateArticle .
func (s *Service) CreateArticle(article *models.Article, tag []string, subjectID uint) (err error) {
	return article.Create(tag, subjectID)
}

// GetArticleList .
func (s *Service) GetArticleList(isAdmin bool, pageNum, pageSize uint, field []interface{}, order string, where []interface{}) (articleList []*models.Article, err error) {
	group := cache.GetKey("")
	key := cache.GetKey(append(where, `GetArticleList`, pageNum, pageSize, isAdmin)...)

	// 获取缓存
	if s.IsCache {
		data, stat, err := cache.GetCacheData(key)
		if err == nil && stat {
			// 数据解析
			jsonData, err := helper.Debase64(data)
			if err == nil {
				json.Unmarshal(jsonData, &articleList)
				return articleList, nil
			}
			logger.Warn("缓存数据有误,无法解析：", group, key, data)
		}
	}

	// 查询数据
	var cacheArticle models.Article
	articleList, err = cacheArticle.List(isAdmin, pageNum, pageSize, field, order, where)
	if isErrDB(err) {
		return nil, err
	}

	// 保存缓存
	dataString, err := json.Marshal(&articleList)
	if err != nil {
		return nil, err
	}
	cache.Get().SetEx(key, conf.Get().Cache.Expired, helper.Enbase64(dataString))

	return
}
