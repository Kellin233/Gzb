package database

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// Enterprise

func SearchEnterpriseByRule(db *gorm.DB, where map[string]interface{}) (error, []Enterprise, int64) {
	var result *multierror.Error
	var enterprise []Enterprise
	var count int64
	err := db.Table("enterprise").Where(where).Find(&enterprise).Count(&count).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	return result, enterprise, count
}

func SearchEnterpriseByName(db *gorm.DB, name string) (error, []Enterprise, int64) {
	var enterprise []Enterprise
	var result error
	name = name + "%"
	err := db.Table("enterprise").Where("name like ?", name).Find(&enterprise).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	var count int64
	err2 := db.Table("enterprise").Where("name like ?", name).Find(&enterprise).Count(&count).Error
	if err2 != nil {
		result = multierror.Append(result, err2)
	}
	return result, enterprise, count
}

func DeleteEnterprise(db *gorm.DB, where map[string]interface{}) (int64, error) {
	var count int64 = 0
	var enterprise Enterprise
	err := db.Where(where).Find(&enterprise).Count(&count).Error
	if count == 0 && err == nil {
		return 0, errors.New("要删除的记录不存在")
	}
	err = db.Delete(&enterprise).Error
	return count, err
}

func UpdateSingleEnterpriseByCondition(db *gorm.DB, where map[string]interface{}, update map[string]interface{}) (int64, error) {
	var count int64 = 0
	err := db.Table("enterprise").Where(where).Count(&count).Error
	if count == 0 && err == nil {
		return 0, errors.New("要修改的记录不存在")
	}
	err = db.Table("enterprise").Where(where).Updates(update).Count(&count).Error
	return count, err
}

func AddEnterprise(db *gorm.DB, values map[string]interface{}) (int64, error) {
	var count int64 = 0
	err := db.Table("enterprise").Where("account = ?", values["account"]).Count(&count).Error
	fmt.Println(values)
	if count > 0 {
		return 0, errors.New("用户已存在")
	}
	err = db.Table("enterprise").Create(values).Count(&count).Error
	return count, err
}

func SelectAllEnterpriseByPage(db *gorm.DB, query string, pageNum int64, pageSize int64) ([]Enterprise, int64, error) {
	var count int64 = 0
	var enterprises []Enterprise
	var err error
	if query != "" {
		query = query + "%"
		db = db.Table("enterprise").Where("enterprise.name like ?", query).Order("id").Find(&enterprises).Count(&count)
	} else {
		db = db.Table("enterprise").Order("id").Find(&enterprises).Count(&count)
	}
	err = db.Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&enterprises).Error
	if count == 0 && err == nil {
		return enterprises, 0, errors.New("要查询的记录不存在")
	}
	return enterprises, count, err
}

func SelectSingleEnterpriseByCondition(db *gorm.DB, where map[string]interface{}) (Enterprise, int64, error) {
	var count int64 = 0
	var enterprise Enterprise
	err := db.Table("enterprise").Where(where).Find(&enterprise).Count(&count).Error
	if count == 0 {
		return enterprise, 0, errors.New("查询的记录不存在")
	}
	return enterprise, count, err
}

// Job

func SelectJobByCondition(db *gorm.DB, where map[string]interface{}) ([]Job, int64, error) {
	var jobs []Job
	var count int64 = 0
	err := db.Table("job").InnerJoins("Enterprise").Where(where).Find(&jobs).Count(&count).Error
	if count == 0 && err == nil {
		return jobs, count, nil
	}
	return jobs, count, err
}

func SearchJobByName(db *gorm.DB, name string) (error, []Job, int64) {
	var jobs []Job
	var result error
	var count int64
	name = name + "%"
	err := db.Table("job").InnerJoins("Enterprise").Where("job.name like ?", name).Find(&jobs).Count(&count).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	return result, jobs, count
}

func SearchSingleJobByName(db *gorm.DB, name string) (error, Job, int64) {
	var job Job
	var count int64
	err2 := db.Where("job.name = ?", name).InnerJoins("Enterprise").Find(&job).Count(&count).Error
	if count == 0 && err2 == nil {
		return nil, job, 0
	}
	return err2, job, count
}

func DeleteJob(db *gorm.DB, where map[string]interface{}) (int64, error) {
	var count int64 = 0
	var job Job
	err := db.Where(where).Find(&job).Count(&count).Error
	if count == 0 && err == nil {
		return 0, errors.New("要删除的记录不存在")
	}
	err = db.Delete(&job).Error
	return count, err
}

func UpdateSingleJobByCondition(db *gorm.DB, where map[string]interface{}, update map[string]interface{}) (int64, error) {
	var count int64 = 0
	err := db.Table("job").Where(where).Count(&count).Error
	if count == 0 && err == nil {
		return 0, errors.New("要修改的记录不存在")
	}
	err = db.Table("job").Where(where).Updates(update).Count(&count).Error
	return count, err
}

func AddJob(db *gorm.DB, values map[string]interface{}) (int64, error) {
	var count int64 = 0
	err := db.Table("job").Create(values).Count(&count).Error
	return count, err
}

func SelectAllJobByPage(db *gorm.DB, query string, pageNum int64, pageSize int64) ([]Job, int64, error) {
	var count int64 = 0
	var jobs []Job
	var err error
	if query != "" {
		query = query + "%"
		db = db.Table("job").InnerJoins("Enterprise").
			Where("enterprise.id = ?", query).Order("job.id").Find(&jobs).Count(&count)
	} else {
		db = db.Table("job").InnerJoins("Enterprise").
			Order("job.id").Find(&jobs).Count(&count)
	}
	err = db.Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&jobs).Error
	if count == 0 && err == nil {
		return jobs, 0, errors.New("要查询的记录不存在")
	}
	return jobs, count, err
}
func SelectAllJobByCondition(db *gorm.DB, where map[string]interface{}) ([]Job, int64, error) {
	var count int64 = 0
	var jobs []Job
	var err error
	db = db.Table("job").InnerJoins("Enterprise").Where("enterprise.name like ?", where).Order("id").Find(&jobs).Count(&count)
	if count == 0 && err == nil {
		return jobs, 0, errors.New("要查询的记录不存在")
	}
	return jobs, count, err
}

func SelectSingleJobByCondition(db *gorm.DB, where map[string]interface{}) (Job, int64, error) {
	var count int64 = 0
	var job Job
	err := db.Table("job").InnerJoins("Enterprise").Where("job.id = ?", where["id"]).Find(&job).Count(&count).Error
	if count == 0 {
		return job, 0, errors.New("查询的记录不存在")
	}
	return job, count, err
}

// Comment

func DeleteComment(db *gorm.DB, where map[string]interface{}) (int64, error) {
	var count int64 = 0
	var comment Comment
	err := db.Where(where).Find(&comment).Count(&count).Error
	if count == 0 && err == nil {
		return 0, errors.New("要删除的记录不存在")
	}
	err = db.Delete(&comment).Error
	return count, err
}

func UpdateSingleCommentByCondition(db *gorm.DB, where map[string]interface{}, update map[string]interface{}) (int64, error) {
	var count int64 = 0
	err := db.Table("comment").Where(where).Count(&count).Error
	if count == 0 && err == nil {
		return 0, errors.New("要修改的记录不存在")
	}
	err = db.Table("comment").Where(where).Updates(update).Count(&count).Error
	return count, err
}

func SelectSingleCommentByCondition(db *gorm.DB, where map[string]interface{}) (Comment, int64, error) {
	var count int64 = 0
	var comment Comment
	err := db.Joins("Author").Where(where).First(&comment).Count(&count).Error
	if count == 0 {
		return comment, 0, errors.New("查询的记录不存在")
	}
	return comment, count, err
}

func SelectAllCommentByPage(db *gorm.DB, query string, pageNum int64, pageSize int64) ([]Comment, int64, error) {
	var count int64 = 0
	var comments []Comment
	if query != "" {
		db = db.Table("comment").InnerJoins("Author").
			Where("comment.targetPost = ?", query).Order("state asc, id").Find(&comments).Count(&count)
	} else {
		db = db.Table("comment").InnerJoins("Author").
			Order("comment.state asc, id").Find(&comments).Count(&count)
	}
	err := db.Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&comments).Error
	if count == 0 && err == nil {
		return comments, 0, errors.New("要查询的记录不存在")
	}
	return comments, count, err
}
func AddComment(db *gorm.DB, values map[string]interface{}) (int64, error) {
	mp := make(map[string]interface{})
	mp["id"] = values["userID"]
	_, count, _ := SelectSingleUserByCondition(db, mp)
	if count == 0 {
		return 0, errors.New("要插入的记录有误")
	}
	values["publishTime"] = time.Now()
	err := db.Table("comment").Create(values).Count(&count).Error
	return count, err
}

func AddComments(db *gorm.DB, values map[string]interface{}) (int64, error) {
	var result *multierror.Error
	mp := make(map[string]interface{})
	mp["ID"] = values["userId"]
	delete(values, "account")
	user, count, err := SelectSingleUserByCondition(db, mp)
	result = multierror.Append(result, err)
	if count == 0 {
		result = multierror.Append(result, errors.New("要插入的记录有误，插入的用户不存在"))
	}
	values["userID"] = user.ID
	delete(mp, "ID")
	if _, ok := values["postId"]; ok {
		mp["ID"] = values["postId"]

	} else if _, ok := values["queId"]; ok {
		mp["ID"] = values["queId"]
	}

	post, count, err := SelectSingleQueByCondition(db, mp)
	result = multierror.Append(result, err)
	if count == 0 {
		result = multierror.Append(result, errors.New("要插入的记录有误，插入的帖子不存在"))
	}

	if _, ok := values["postId"]; ok {
		values["content"] = values["comment"]
		delete(values, "comment")
		delete(values, "postId")
	} else if _, ok := values["queId"]; ok {
		values["content"] = values["answer"]
		delete(values, "answer")
		delete(values, "queId")
	}
	values["targetPost"] = post.ID

	delete(values, "userId")
	t1 := time.Now().Year()
	t2 := time.Now().Month()
	t3 := time.Now().Day()
	t4 := time.Now().Hour()
	t5 := time.Now().Minute()
	t6 := time.Now().Second()
	t7 := time.Now().Nanosecond()
	currentTimeData := time.Date(t1, t2, t3, t4, t5, t6, t7, time.Local) //获取当前时间，返回当前时间Time
	values["publishTime"] = currentTimeData
	values["state"] = 0
	err = db.Table("comment").Create(values).Count(&count).Error
	result = multierror.Append(result, err)
	return count, result
}

func SearchCommentByQueId(db *gorm.DB, queId int64) (int64, []Comment, error) {
	var result *multierror.Error
	var comments []Comment
	var whereMap map[string]interface{} = make(map[string]interface{})
	whereMap["ID"] = queId
	whereMap["partID"] = 2
	_, count, err := SelectSingleQueByCondition(db, whereMap)
	if count == 0 {
		result = multierror.Append(result, errors.New("帖子id不是问题！"))
	}
	if err != nil {
		result = multierror.Append(result, err)
	}
	err = db.Preload("Author").Where("targetPost = ? ", queId).Order("publishTime DESC").Find(&comments).Count(&count).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	return count, comments, result
}

func GetAdoptedAnswerByQueId(db *gorm.DB, queId int64) (string, error) {
	var result *multierror.Error
	var adoptRecord AdoptRecord
	var count int64 = 0
	err := db.Where("postId = ? ", queId).Find(&adoptRecord).Count(&count).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	return adoptRecord.CommentId, result
}

func AddAdoptRecord(db *gorm.DB, queId int64, answerId int64) (bool, error) {
	tx := db.Begin()

	adoptRecord := AdoptRecord{
		PostId:    strconv.FormatInt(queId, 10),
		CommentId: strconv.FormatInt(answerId, 10),
	}
	err := tx.Create(&adoptRecord).Error
	if err != nil {
		tx.Rollback()
		return false, err
	}

	var post Post
	err = tx.Where("ID = ?", queId).First(&post).Error
	if err != nil {
		tx.Rollback()
		return false, err
	}

	reward := post.Reward

	var comment Comment
	err = tx.Where("ID = ?", answerId).First(&comment).Error
	if err != nil {
		tx.Rollback()
		return false, err
	}

	var user User
	err = tx.Where("ID = ?", comment.UserID).First(&user).Error
	if err != nil {
		tx.Rollback()
		return false, err
	}

	err = tx.Model(&Post{}).Where("ID = ?", queId).Update("Reward", 0).Error
	if err != nil {
		tx.Rollback()
		return false, err
	}

	err = tx.Model(&User{}).Where("ID = ?", comment.UserID).Update("Balance", user.Balance+reward).Error
	if err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}

// News

func DeleteNews(db *gorm.DB, where map[string]interface{}) (int64, error) {
	var count int64 = 0
	var news News
	err := db.Where(where).Find(&news).Count(&count).Error
	if count == 0 && err == nil {
		return 0, errors.New("要删除的记录不存在")
	}
	err = db.Delete(&news).Error
	return count, err
}

func SelectSingleNewsByCondition(db *gorm.DB, where map[string]interface{}) (News, int64, error) {
	var count int64 = 0
	var news News
	err := db.Where(where).First(&news).Count(&count).Error
	if count == 0 {
		return news, 0, errors.New("查询的记录不存在")
	}
	return news, count, err
}

func SelectAllNewsByCondition(db *gorm.DB, where map[string]interface{}) ([]News, int64, error) {
	var count int64 = 0
	var newses []News
	err := db.Where(where).Find(&newses).Count(&count).Error
	return newses, count, err
}

func SelectAllNewsByPage(db *gorm.DB, query string, pageNum int64, pageSize int64) ([]News, int64, error) {
	var count int64 = 0
	var newses []News

	if query != "" {
		query = query + "%"
		db = db.Table("news").Where("author like ?", query).Count(&count)
	} else {
		db = db.Table("news").Count(&count)
	}
	err := db.Order("id asc").Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&newses).Error
	if count == 0 && err == nil {
		return newses, 0, errors.New("查询的记录不存在")
	}
	return newses, count, err
}

func UpdateSingleNewsByCondition(db *gorm.DB, where map[string]interface{}, update map[string]interface{}) (int64, error) {
	var count int64 = 0
	err := db.Table("news").Where(where).Count(&count).Error
	if count == 0 && err == nil {
		return 0, errors.New("要修改的记录不存在")
	}
	err = db.Table("news").Where(where).Updates(update).Count(&count).Error
	return count, err
}

func AddNews(db *gorm.DB, values map[string]interface{}) (int64, error) {
	var count int64 = 0
	err := db.Table("news").Create(values).Count(&count).Error
	return count, err
}

// Post ------------------------------------------------------------

func AddPost(db *gorm.DB, values map[string]interface{}) (int64, error) {
	mp := make(map[string]interface{})
	mp["account"] = values["account"]
	delete(values, "account")
	user, count, _ := SelectSingleUserByCondition(db, mp)
	if count == 0 {
		return 0, errors.New("要插入的记录有误")
	}
	values["authorID"] = user.ID
	values["favorite"] = 0
	values["like"] = 0
	values["publishTime"] = time.Now()
	err := db.Table("post").Create(values).Count(&count).Error
	return count, err
}

func DeletePost(db *gorm.DB, where map[string]interface{}) (int64, error) {
	var count int64 = 0
	var post Post
	err := db.Where(where).Find(&post).Count(&count).Error
	if count == 0 && err == nil {
		return 0, errors.New("要删除的记录不存在")
	}
	err = db.Delete(&post).Error
	return count, err
}

func SelectSinglePostByCondition(db *gorm.DB, where map[string]interface{}) (Post, int64, error) {
	var count int64 = 0
	var post Post
	err := db.Table("post").InnerJoins("Author").InnerJoins("Part").Where("post.id = ?", where["id"]).Find(&post).Count(&count).Error
	if count != 0 {
		return post, count, err
	}
	return post, 0, errors.New("查询的记录不存在")
}

func SelectSingleQueByCondition(db *gorm.DB, where map[string]interface{}) (Post, int64, error) {
	var count int64 = 0
	var post Post
	err := db.Where(where).First(&post).Count(&count).Error
	if count == 0 {
		return post, 0, errors.New("查询的记录不存在")
	}
	return post, count, err
}

func SelectAllPostByCondition(db *gorm.DB, where map[string]interface{}) ([]Post, int64, error) {
	var count int64 = 0
	var posts []Post
	err := db.InnerJoins("Author").InnerJoins("Part").Where(where).Count(&count).Error
	return posts, count, err
}

func SelectAllPostsByAuthorId(db *gorm.DB, authorId int64) ([]Post, int64, error) {
	var count int64 = 0
	var posts []Post
	err := db.Table("post").Where("authorId = ? ", authorId).Find(&posts).Count(&count).Error
	return posts, count, err
}

func SelectAllPostByPage(db *gorm.DB, query string, pageNum int64, pageSize int64) ([]Post, int64, error) {
	var count int64 = 0
	var posts []Post
	var err error
	if query != "" {
		query = query + "%"
		db = db.Table("post").InnerJoins("Author").InnerJoins("Part").
			Where("account like ?", query).Order("state asc, id").Find(&posts).Count(&count)
	} else {
		db = db.Table("post").Joins("Author").Joins("Part").
			Order("state asc, id").Find(&posts).Count(&count)
	}
	err = db.Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&posts).Error
	if count == 0 && err == nil {
		return posts, 0, errors.New("要查询的记录不存在")
	}
	return posts, count, err
}

func UpdateSinglePostByCondition(db *gorm.DB, where map[string]interface{}, update map[string]interface{}) (int64, error) {
	var count int64 = 0
	err := db.Table("post").Where(where).Count(&count).Error
	if count == 0 && err == nil {
		return 0, errors.New("要修改的记录不存在")
	}
	err = db.Table("post").Where(where).Updates(update).Count(&count).Error
	return count, err
}

func SearchAllNewInfo(db *gorm.DB, userId int64) ([]Post, error) {
	var result *multierror.Error
	var posts []Post
	var posts1 []Post
	var postResult []Post
	err := db.Preload("Author").Where("partID = ? && AuthorId = ? ", 1, userId).Order("publishTime DESC").Find(&posts).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	err = db.Preload("Author").Where("partID = ? && AuthorId != ? ", 1, userId).Order("publishTime DESC").Find(&posts1).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	postResult = append(posts, posts1...)
	return postResult, err
}

func SearchNewInfoComment(db *gorm.DB, postId string) (error, []Comment) {
	var result *multierror.Error
	var comments []Comment
	err := db.Preload("Author").Where("targetPost = ? ", postId).Order("publishTime DESC").Find(&comments).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	return result, comments
}

func SearchNewInfoDetails(db *gorm.DB, postId int64) (error, []Post) {
	var result *multierror.Error
	var posts []Post
	err := db.Preload("Author").Where("Id = ? && partID = ? ", postId, 1).Find(&posts).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	return result, posts
}

func IsExistedLikeRecord(db *gorm.DB, postId int64, userId int64) (error, bool) {
	var result *multierror.Error
	var likeRecord []LikeRecord
	var isExistedLikeRecord bool
	var count int64
	err := db.Preload("Post").Preload("Author").Where("postId = ? && userId = ? ", postId, userId).Find(&likeRecord).Count(&count).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	if count > 0 {
		isExistedLikeRecord = true
	} else {
		isExistedLikeRecord = false
	}
	return result, isExistedLikeRecord
}

func GetLikeNumByPostId(db *gorm.DB, postId int64) (error, int64) {
	var result *multierror.Error
	var likeRecord []LikeRecord
	var count int64
	err := db.Preload("Post").Preload("Author").Where("postId = ? ", postId).Find(&likeRecord).Count(&count).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	return result, count
}

func IsExistedFavoriteRecord(db *gorm.DB, postId int64, userId int64) (error, bool) {
	var result *multierror.Error
	var favoriteRecord []FavoriteRecord
	var isExistedFavoriteRecord bool
	var count int64
	err := db.Preload("Article").Preload("Author").Where("articleID = ? && userID = ? ", postId, userId).Find(&favoriteRecord).Count(&count).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	if count > 0 {
		isExistedFavoriteRecord = true
	} else {
		isExistedFavoriteRecord = false
	}
	return result, isExistedFavoriteRecord
}

func IsExistedAdoptRecord(db *gorm.DB, postId int64) (bool, error) {
	var count int64
	err := db.Model(&AdoptRecord{}).Where("postId = ?", postId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func AddPostFrontend(db *gorm.DB, where map[string]interface{}) (bool, error) {
	var result *multierror.Error
	var resultSign bool

	userId := int64(where["userId"].(float64))
	title := where["title"].(string)
	content := where["content"].(string)
	partId := int64(where["type"].(float64))
	summary := where["summary"].(string)
	coverUrl := where["img"].(string)
	reward := int64(where["reward"].(float64))

	err := db.Transaction(func(tx *gorm.DB) error {
		// 创建帖子记录
		post := Post{
			AuthorID:    userId,
			Title:       title,
			Content:     content,
			PartID:      partId,
			Summary:     summary,
			CoverUrl:    coverUrl,
			PublishTime: time.Now(),
			Reward:      reward,
		}
		if err := tx.Create(&post).Error; err != nil {
			return err
		}

		// 更新用户余额
		user := User{}
		if err := tx.Where("ID = ?", userId).First(&user).Error; err != nil {
			return err
		}
		user.Balance -= reward
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		result = multierror.Append(result, err)
		resultSign = false
	} else {
		resultSign = true
	}

	return resultSign, result
}

// User ------------------------------------------------------------

func SelectUserForLogin(db *gorm.DB, account string, password string) (User, int64, error) {
	var count int64 = 0
	var user User
	err := db.Where("account = ?", account).First(&user).Count(&count).Error
	if count == 0 {
		return user, 0, errors.New("账户不存在")
	} else if password != user.Password {
		return user, 0, errors.New("密码错误！")
	}
	return user, count, err
}
func SelectSingleUserByCondition(db *gorm.DB, where map[string]interface{}) (User, int64, error) {
	var count int64 = 0
	var user User
	err := db.Where(where).First(&user).Count(&count).Error
	if count == 0 {
		return user, 0, errors.New("要查询的记录不存在")
	}
	return user, count, err
}
func SelectAllUserByCondition(db *gorm.DB, where map[string]interface{}) ([]User, int64, error) {
	var count int64 = 0
	var users []User
	err := db.Where(where).Find(&users).Count(&count).Error
	if count == 0 {
		return users, 0, errors.New("要查询的记录不存在")
	}
	return users, count, err
}

func SelectAllUserByPage(db *gorm.DB, query string, pageNum int64, pageSize int64) ([]User, int64, error) {
	var count int64 = 0
	var users []User
	if query != "" {
		query = query + "%"
		db = db.Table("user").Where("account like ?", query)
	}
	db.Table("user").Count(&count)
	err := db.Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&users).Error
	if count == 0 && err == nil {
		return users, 0, errors.New("要查询的记录不存在")
	}
	return users, count, nil
}

func AddUser(db *gorm.DB, values map[string]interface{}) (int64, error) {
	var count int64 = 0
	err := db.Table("user").Where("account = ? ", values["account"]).Count(&count).Error
	if count > 0 && err == nil {
		return 0, errors.New("记录已存在")
	}
	err = db.Table("user").Create(values).Count(&count).Error
	return count, err
}

func DeleteUser(db *gorm.DB, where map[string]interface{}) (int64, error) {
	var count int64 = 0
	var user User
	err := db.Where(where).Find(&user).Count(&count).Error
	if count == 0 && err == nil {
		return 0, errors.New("要删除的记录不存在")
	}
	err = db.Delete(&user).Error
	return count, err
}

func AddBalance(db *gorm.DB, userID int64, sum int64) (bool, error) {
	err := db.Model(&User{}).Where("id = ?", userID).Update("balance", gorm.Expr("balance + ?", sum)).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func DecreaseBalance(db *gorm.DB, userID int64, amount int64) (string, bool, error) {
	var user User
	err := db.First(&user, userID).Error
	if err != nil {
		return "", false, err
	}

	if user.Balance < amount {
		return "", false, errors.New("学币余额不足")
	}

	user.Balance -= amount
	err = db.Save(&user).Error
	if err != nil {
		return "", false, err
	}

	return user.Account, true, nil
}

// Feedback

func AddFeedback(db *gorm.DB, values map[string]interface{}) (int64, error) {
	var count int64 = 0
	err := db.Table("feedback").Create(values).Count(&count).Error
	return count, err
}

func SelectAllFeedbackByPage(db *gorm.DB, query string, pageNum int64, pageSize int64) ([]Feedback, int64, error) {
	var count int64 = 0
	var feedbacks []Feedback
	if query != "" {
		query = query + "%"
		db = db.Table("feedback").InnerJoins("Author").Where("account like ?", query).Order("state, id asc").Find(&feedbacks).Count(&count)
	} else {
		db = db.Table("feedback").InnerJoins("Author").Order("state, id asc").Find(&feedbacks).Count(&count)
	}
	err := db.Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&feedbacks).Error
	if count == 0 && err == nil {
		return feedbacks, 0, errors.New("查询的记录不存在")
	}
	return feedbacks, count, err
}

func UpdateSingleFeedbackByCondition(db *gorm.DB, where map[string]interface{}, update map[string]interface{}) (int64, error) {
	var count int64 = 0
	err := db.Table("feedback").Where(where).Count(&count).Error
	if count == 0 && err == nil {
		return 0, errors.New("要修改的记录不存在")
	}
	err = db.Table("feedback").Where(where).Updates(update).Count(&count).Error
	return count, err
}

func SearchFeedbackByUserId(db *gorm.DB, userId int64) (int64, []Feedback, error) {
	var result *multierror.Error
	var feedbacks []Feedback
	var count int64
	err := db.Preload("Author").Where("userID = ? ", userId).Find(&feedbacks).Count(&count).Error
	if count == 0 {
		result = multierror.Append(result, errors.New("找不到该收藏记录！"))
	}
	if err != nil {
		result = multierror.Append(result, err)
	}
	return count, feedbacks, result
}

// Admin ------------------------------------------------------------

func SelectSingleAdminByCondition(db *gorm.DB, where map[string]interface{}) (Admin, int64, error) {
	var count int64 = 0
	var admin Admin
	err := db.Where(where).First(&admin).Count(&count).Error
	if count == 0 && err == nil {
		return admin, 0, errors.New("查询的记录不存在！")
	}
	return admin, count, err
}
func UpdateSingleAdminByCondition(db *gorm.DB, where map[string]interface{}, values map[string]interface{}) (int64, error) {
	var count int64 = 0
	err := db.Table("admin").Where(where).Updates(values).Count(&count).Error
	return count, err
}

func UpdateSingleUserByCondition(db *gorm.DB, where map[string]interface{}, values map[string]interface{}) (int64, error) {
	var count int64 = 0
	err := db.Table("user").Where(where).Count(&count).Error
	if count == 0 && err == nil {
		return 0, errors.New("要修改的记录不存在！")
	}
	err = db.Table("user").Where(where).Updates(values).Count(&count).Error
	return count, err
}

// Que ------------------------------------------------------------

func SearchAllQue(db *gorm.DB, userId int64) (int64, []Post, error) {
	var result *multierror.Error
	var posts []Post
	var count1 int64
	var count2 int64
	err := db.Preload("Author").Where("partID = ? && authorId = ? ", 2, userId).Order("publishTime DESC").Find(&posts).Count(&count1).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	var posts1 []Post
	err = db.Preload("Author").Where("partID = ? && authorId != ? ", 2, userId).Order("publishTime DESC").Find(&posts1).Count(&count2).Error
	if err != nil {
		result = multierror.Append(result, err)
	}
	postsResult := append(posts, posts1...)
	return count1 + count2, postsResult, result
}

func SearchQueByQueId(db *gorm.DB, queId int64) (int64, []Post, error) {
	var result *multierror.Error
	var posts []Post
	var count int64
	err := db.Preload("Author").Where("ID = ? and partID = 2", queId).Find(&posts).Count(&count).Error
	if count == 0 {
		result = multierror.Append(result, errors.New("找不到该问题！"))
	}
	if err != nil {
		result = multierror.Append(result, err)
	}
	return count, posts, result
}

// Favorite ------------------------------------------------------------

func SearchFavoriteByUserId(db *gorm.DB, userId int64) (int64, []FavoriteRecord, error) {
	var result *multierror.Error
	var favorites []FavoriteRecord
	var count int64
	err := db.Preload("Article").Preload("Author").Where("userID = ? ", userId).Find(&favorites).Count(&count).Error
	if count == 0 {
		result = multierror.Append(result, errors.New("找不到该收藏记录！"))
	}
	if err != nil {
		result = multierror.Append(result, err)
	}
	return count, favorites, result
}

func DeleteFavoriteByCondition(db *gorm.DB, where map[string]interface{}) (int64, error) {
	var count int64
	var favoriteRecord FavoriteRecord

	err := db.Where(where).Find(&favoriteRecord).Count(&count).Error
	fmt.Println(favoriteRecord)
	if count == 0 && err == nil {
		return 0, errors.New("要删除的记录不存在")
	}
	err = db.Delete(&favoriteRecord).Error
	return count, err
}

func UpdateFavoriteByCondition(db *gorm.DB, where map[string]interface{}) (bool, error) {
	var result *multierror.Error
	var count int64
	var resultSign bool
	var favoriteRecord FavoriteRecord
	if where["isCollected"] == "false" {
		delete(where, "isCollected")
		err := db.Where(where).Find(&favoriteRecord).Count(&count).Error
		if count == 0 && err == nil {
			result = multierror.Append(result, errors.New("要删除的记录不存在"))
			result = multierror.Append(result, err)
		}
		err = db.Delete(&favoriteRecord).Error
		if err != nil {
			result = multierror.Append(result, err)
			resultSign = false
		} else {
			resultSign = true
			err = db.Table("post").Where("id = ?", where["articleID"]).UpdateColumn("favorite", gorm.Expr("favorite - ?", 1)).Error
			if err != nil {
				result = multierror.Append(result, err)
				resultSign = false
			}
		}

	} else if where["isCollected"] == "true" {
		fr := FavoriteRecord{
			UserID:    where["userID"].(string),
			ArticleID: where["articleID"].(string),
		}
		err := db.Create(&fr).Error
		if err != nil {
			result = multierror.Append(result, err)
			resultSign = false
		} else {
			resultSign = true
			err = db.Table("post").Where("id = ?", where["articleID"]).UpdateColumn("favorite", gorm.Expr("favorite + ?", 1)).Error
			if err != nil {
				result = multierror.Append(result, err)
				resultSign = false
			}
		}

	}
	return resultSign, result
}

// Like ------------------------------------------------------------

func UpdateLikeByCondition(db *gorm.DB, where map[string]interface{}) (bool, error) {
	var result *multierror.Error
	var count int64
	var resultSign bool
	var likeRecord LikeRecord
	if where["isLiked"] == "false" {
		delete(where, "isLiked")
		err := db.Where(where).Find(&likeRecord).Count(&count).Error
		if count == 0 && err == nil {
			result = multierror.Append(result, errors.New("要删除的记录不存在"))
			result = multierror.Append(result, err)
		}
		err = db.Delete(&likeRecord).Error
		if err != nil {
			result = multierror.Append(result, err)
			resultSign = false
		} else {
			resultSign = true
			err = db.Table("post").Where("id = ?", where["postId"]).UpdateColumn("`like`", gorm.Expr("`like` - ?", 1)).Error
			if err != nil {
				result = multierror.Append(result, err)
				resultSign = false
			}
		}

	} else if where["isLiked"] == "true" {
		lr := LikeRecord{
			UserId: where["userId"].(string),
			PostId: where["postId"].(string),
		}
		err := db.Create(&lr).Error
		if err != nil {
			result = multierror.Append(result, err)
			resultSign = false
		} else {
			resultSign = true
			err = db.Table("post").Where("id = ?", where["postId"]).UpdateColumn("`like`", gorm.Expr("`like` + ?", 1)).Error
			if err != nil {
				result = multierror.Append(result, err)
				resultSign = false
			}
		}

	}
	return resultSign, result
}
