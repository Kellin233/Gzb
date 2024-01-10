package news

import (
	fybDatabase "FybBackend/database"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"gorm.io/gorm"
)

func NewsDetail(e *gin.Engine, db *gorm.DB) {
	e.GET("/v1/frontend/news/newsDetail", func(context *gin.Context) {
		var result *multierror.Error
		mp := make(map[string]interface{})
		mp["id"] = context.DefaultQuery("id", "")
		news, _, err1 := fybDatabase.SelectSingleNewsByCondition(db, mp)
		result = multierror.Append(result, err1)

		if result.ErrorOrNil() == nil {
			context.JSON(200, gin.H{
				"code":    200,
				"message": "get all newsInfo success!",
				"data":    news,
			})
		} else {
			code := 404
			context.JSON(code, gin.H{
				"code":    code,
				"message": result.Error(),
			})
		}
	})
}
