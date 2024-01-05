package routers

import (
	"FybBackend/routers/v1/frontend/academy"
	"FybBackend/routers/v1/frontend/circle"
	"FybBackend/routers/v1/frontend/job"
	"FybBackend/routers/v1/frontend/news"
	"FybBackend/routers/v1/frontend/user/exchange"
	"FybBackend/routers/v1/frontend/user/login"
	"FybBackend/routers/v1/frontend/user/myFavorites"
	"FybBackend/routers/v1/frontend/user/myFeedbacks"
	"FybBackend/routers/v1/frontend/user/myPosts"
	"FybBackend/routers/v1/frontend/user/recharge"
	"FybBackend/routers/v1/frontend/user/register"
	"FybBackend/routers/v1/frontend/user/selectUsers"
	"FybBackend/routers/v1/frontend/user/userInfo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitFrontend(r *gin.Engine, db *gorm.DB) {

	//feedback
	myFeedbacks.SearchNewQue(r)

	//user
	selectUsers.SelectUsers(r, db)
	userInfo.BasicUserInfo(r, db)
	userInfo.Settings(r, db)
	userInfo.ResetPassword(r, db)
	userInfo.GetPassword(r, db)
	userInfo.AccountSecurity(r, db)
	userInfo.AddFeedback(r, db)
	userInfo.MyPoss(r, db)
	register.UserRegister(r, db)
	login.PasswordLogin(r, db)
	myPosts.SearchNewQue(r)
	myPosts.DeletePost(r, db)
	myFavorites.DeleteFavorite(r, db)
	register.SendEmail(r)
	register.ValidateEmailCode(r)
	recharge.Recharge(r)
	exchange.Exchange(r)

	//news
	news.NewsList(r, db)
	news.NewsDetail(r, db)

	//job
	job.SearchByRule(r, db)
	job.SearchByName(r, db)
	job.SelectByCode(r, db)

	//academy
	academy.SearchByName(r)
	academy.SelectAcademyByCode(r)
	academy.SelectScoreByTypeFirstSecondLevel(r)
	academy.SearchByRule(r)

	//circle
	circle.SearchByName(r)
	circle.SearchNewInfoComment(r)
	circle.SearchNewInfoDetails(r)
	circle.PostComments(r)
	circle.PostAnswers(r)
	circle.SearchNewQue(r)
	circle.SearchQueAnswer(r)
	circle.SearchQueDetails(r)
	circle.UpdatePostCollected(r, db)
	circle.UpdatePostLiked(r, db)
	circle.UploadPost(r, db)
	circle.PostAnswerStatus(r, db)

	//favorite
	myFavorites.SearchNewQue(r)
}
