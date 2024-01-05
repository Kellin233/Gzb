package userInfo

import (
	fybDatabase "FybBackend/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"gorm.io/gorm"
)

func AccountSecurity(e *gin.Engine, db *gorm.DB) {
	e.POST("/v1/frontend/user/accountSecurity", func(context *gin.Context) {
		var result *multierror.Error
		mp1 := make(map[string]interface{})
		mp2 := make(map[string]interface{})
		b, err1 := context.GetRawData()
		err2 := json.Unmarshal(b, &mp1)
		mp2["id"] = mp1["id"]
		delete(mp1, "id")
		_, err3 := fybDatabase.UpdateSingleUserByCondition(db, mp2, mp1)
		result = multierror.Append(result, err1, err2, err3)

		if result.ErrorOrNil() == nil {
			context.JSON(200, gin.H{
				"code":    200,
				"message": "重置用户密码和邮箱成功！",
			})
		} else {
			context.JSON(404, gin.H{
				"code":    404,
				"message": result.Error(),
			})
		}
	})
}
