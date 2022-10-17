package controllers

import (
	"final-project/helpers"
	"final-project/library"
	"final-project/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

type CreateCommentReq struct {
	Message  string `json:"message" validate:"required,max=191"`
	Photo_Id uint   `json:"photo_id" validate:"required"`
}

func CreateComment(ctx *gin.Context) {
	var createReq CreateCommentReq
	var dataComment models.Comment
	userData := ctx.MustGet("userData").(jwt.MapClaims)

	ctx.ShouldBindJSON(&createReq)

	valid, trans := helpers.Valid()
	err := valid.Struct(createReq)

	if err != nil {
		errs := err.(validator.ValidationErrors)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Error Validation",
			"message": errs.Translate(trans),
		})
		return
	}

	dataComment = models.Comment{
		Message:    createReq.Message,
		Photo_Id:   createReq.Photo_Id,
		Created_at: time.Now(),
		User_Id:    uint(userData["id"].(float64)),
	}

	errCreate := library.GetDB().Table("comments").Create(&dataComment).Error

	if errCreate != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error Create Comment",
			"message": errCreate,
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         dataComment.Id,
		"message":    dataComment.Message,
		"photo_id":   dataComment.Photo_Id,
		"created_at": dataComment.Created_at,
	})

}

type Photo struct {
	Id        uint   `json:"id"`
	Title     string `json:"title"`
	Caption   string `json:"caption"`
	Photo_url string `json:"photo_url"`
	User_id   uint   `json:"user_id"`
}
type CommentRes struct {
	Id         string    `json:"id"`
	Message    string    `json:"message"`
	Photo_id   string    `json:"photo_id"`
	User_id    string    `json:"user_id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	User       *User     `json:"user"`
	Photo      *Photo    `json:"photo"`
}

func GetComments(ctx *gin.Context) {
	var listComment []CommentRes

	errGet := library.GetDB().Debug().Table("comments").Preload("User").Preload("Photo").Find(&listComment).Error

	if errGet != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error Get Data",
			"message": errGet,
		})

		return
	}

	ctx.JSON(http.StatusOK, listComment)

}

type UpdateCommentReq struct {
	Message string `json:"message"`
}

func UpdateComment(ctx *gin.Context) {
	commentId := ctx.Param("commentId")

	var updateReq UpdateCommentReq
	var detailComment models.Comment

	ctx.ShouldBindJSON(&updateReq)

	valid, trans := helpers.Valid()
	err := valid.Struct(updateReq)

	if err != nil {
		errs := err.(validator.ValidationErrors)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Error Validation",
			"message": errs.Translate(trans),
		})
		return
	}

	errDetail := library.GetDB().Table("comments").Where("id = ?", commentId).Take(&detailComment).Error

	if errDetail != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Error Update",
			"message": "Unable to find Comment",
		})
		return
	}

	library.GetDB().Table("comments").Model(&detailComment).Where("id = ?", commentId).Updates(models.Comment{
		Message:    updateReq.Message,
		Updated_at: time.Now(),
	})

	ctx.JSON(http.StatusOK, gin.H{
		"id":         detailComment.Id,
		"photo_id":   detailComment.Photo_Id,
		"message":    detailComment.Message,
		"user_id":    detailComment.User_Id,
		"updated_at": detailComment.Updated_at,
	})

}

func DeleteComment(ctx *gin.Context) {
	commentId := ctx.Param("commentId")

	errDelete := library.GetDB().Table("comments").Where("id = ?", commentId).Delete(models.Comment{}).Error

	if errDelete != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error Delete",
			"message": "Comment Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}
