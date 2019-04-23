/**
 *
 * By So http://sooo.site
 * -----
 *    Don't panic.
 * -----
 *
 */

package routers

import (
	v1 "github.com/Git-So/blog-api/api/v1"
	"github.com/Git-So/blog-api/middleware/apicache"
	"github.com/Git-So/blog-api/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func v1router(r *gin.Engine) {
	// v1路由组
	base := r.Group("v1")

	// 公共路由组
	pubilc := base.Group("")
	{
		// 文章
		article := pubilc.Group("article")
		article.Use(apicache.Cache()) // api 处缓存
		{
			article.GET("info/:ID", v1.ArticleInfo)
			article.GET("hot/:PageNum", v1.HotArticleList)
			article.GET("list/:PageNum", v1.ArticleList)
		}

		// 评论
		comment := pubilc.Group("comment")
		{
			comment.GET("list/:PageNum", v1.CommentList)
			comment.POST("info", v1.CreateComment)
		}

		// 验证码
		captcha := pubilc.Group("captcha")
		{
			captcha.GET("", v1.GetCaptcha)
			captcha.POST("", v1.VerifyCaptcha)
		}

		// 登入
		login := pubilc.Group("login")
		{
			login.GET("/code", v1.GetLoginCode)
		}

		// 标签
		tag := pubilc.Group("tag")
		tag.Use(apicache.Cache()) // api 处缓存
		{
			tag.GET("list/:PageNum", v1.TagList)
		}

		// 专题
		subject := pubilc.Group("subject")
		subject.Use(apicache.Cache()) // api 处缓存
		{
			subject.GET("list/:PageNum", v1.SubjectList)
			subject.GET("info/:ID", v1.SubjectInfo)
		}

		// 友链
		link := pubilc.Group("link")
		link.Use(apicache.Cache()) // api 处缓存
		{
			link.GET("list/:PageNum", v1.LinkList)
		}

		// 配置
		config := pubilc.Group("config")
		{
			config.GET("list/me", v1.MeInfo)
		}

	}

	// 私有路由组
	private := base.Group("")
	private.Use(jwt.Auth())
	{
		// 文章
		article := private.Group("article")
		{
			article.POST("info", v1.CreateArticle)
			article.PUT("info", v1.UpdateArticle)
			article.DELETE("info", v1.DeleteArticle)
		}

		// 评论
		comment := private.Group("comment")
		{
			comment.PUT("info", v1.UpdateComment)
			comment.DELETE("info", v1.DeleteComment)
		}

		// 标签
		tag := private.Group("tag")
		{
			tag.POST("info", v1.CreateTag)
			tag.DELETE("info", v1.DeleteTag)
		}

		// 专题
		subject := private.Group("subject")
		{
			subject.POST("info", v1.CreateSubject)
			subject.PUT("info", v1.UpdateSubject)
			subject.DELETE("info", v1.DeleteSubject)
		}

		// 友链
		link := private.Group("link")
		{
			link.POST("info", v1.CreateLink)
			link.PUT("info", v1.UpdateLink)
			link.DELETE("info", v1.DeleteLink)
		}

		// 配置
		config := private.Group("config")
		{
			config.POST("info", v1.UpdateConfig)
		}
	}

	return
}
