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

type CreateReq struct {
	Title     string `json:"title" validate:"required,max=50"`
	Caption   string `json:"caption" validate:"required,max=150"`
	Photo_url string `json:"photo_url" validate:"required"`
}

func CreatePhoto(ctx *gin.Context) {
	var createReq CreateReq
	var dataPhoto models.Photo

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

	dataPhoto = models.Photo{
		Title:      createReq.Title,
		Caption:    createReq.Caption,
		Photo_url:  createReq.Photo_url,
		Created_at: time.Now(),
		User_Id:    uint(userData["id"].(float64)),
	}

	errCreate := library.GetDB().Table("photos").Create(&dataPhoto).Error
	if errCreate != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err":     "Error Create Photo",
			"message": errCreate,
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         dataPhoto.Id,
		"title":      dataPhoto.Title,
		"caption":    dataPhoto.Caption,
		"user_id":    dataPhoto.User_Id,
		"created_at": dataPhoto.Created_at,
		"photo_url":  dataPhoto.Photo_url,
	})
}

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username" `
	Email    string `json:"email" `
}

type GetPhotoRes struct {
	Id         uint      `json:"id"`
	Title      string    `json:"title"`
	Caption    string    `json:"caption"`
	User_Id    uint      `json:"user_id"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
	User       *User     `json:"user"`
}

func GetPhoto(ctx *gin.Context) {
	var listPhoto []GetPhotoRes

	err := library.GetDB().Debug().Table("photos").Preload("User").Find(&listPhoto).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"result":  nil,
			"error":   "Error Get Photo",
			"message": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": listPhoto,
	})
}

func UpdatePhoto(ctx *gin.Context) {
	photoId := ctx.Param("photoId")

	var updateReq CreateReq
	var detailPhoto models.Photo
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

	errDetail := library.GetDB().Table("photos").Where("id = ?", photoId).Take(&detailPhoto).Error

	if errDetail != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Error Update",
			"message": "Unable to find Photo",
		})
		return
	}

	library.GetDB().Table("photos").Model(&detailPhoto).Where("id = ?", photoId).Updates(models.Photo{
		Title:      updateReq.Title,
		Caption:    updateReq.Caption,
		Photo_url:  updateReq.Photo_url,
		Updated_at: time.Now(),
	})

	ctx.JSON(http.StatusOK, gin.H{
		"id":         detailPhoto.Id,
		"title":      detailPhoto.Title,
		"caption":    detailPhoto.Caption,
		"photo_url":  detailPhoto.Photo_url,
		"user_id":    detailPhoto.User_Id,
		"updated_at": detailPhoto.Updated_at,
	})
}

func DeletePhoto(ctx *gin.Context) {
	photoId := ctx.Param("photoId")
	errDelete := library.GetDB().Table("photos").Where("id = ?", photoId).Delete(models.Photo{}).Error

	if errDelete != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error Delete",
			"message": "Photo Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted!",
	})

}
