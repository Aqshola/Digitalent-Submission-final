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

type CreateSocialReq struct {
	Name             string `json:"name" validate:"required,max=50"`
	Social_media_url string `json:"social_media_url" validate:"required,max=191"`
}

func CreateSocial(ctx *gin.Context) {
	var createReq CreateSocialReq
	var dataSocial models.SocialMedia
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

	dataSocial = models.SocialMedia{
		Name:             createReq.Name,
		Social_Media_Url: createReq.Social_media_url,
		Created_at:       time.Now(),
		User_Id:          uint(userData["id"].(float64)),
	}

	errCreate := library.GetDB().Table("social_media").Create(&dataSocial).Error

	if errCreate != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error Create Comment",
			"message": errCreate,
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":               dataSocial.Id,
		"message":          dataSocial.Name,
		"social_media_url": dataSocial.Social_Media_Url,
		"user_id":          dataSocial.User_Id,
		"created_at":       dataSocial.Created_at,
	})

}

type SocialRes struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Social_Media_Url string `json:"social_media_url"`
	Created_at       string `json:"created_at"`
	Updated_at       string `json:"updated_at"`
	User_Id          string `json:"-"`
	User             *User  `json:"user"`
}

func GetSocials(ctx *gin.Context) {
	var listSocial []SocialRes

	errGet := library.GetDB().Table("social_media").Preload("User").Find(&listSocial).Error

	if errGet != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error Get Social",
			"message": errGet.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, listSocial)

}

func UpdateSocial(ctx *gin.Context) {
	var updateReq CreateSocialReq
	var detailSocial models.SocialMedia
	socialMediaId := ctx.Param("socialMediaId")

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

	errDetail := library.GetDB().Table("social_media").Where("id = ?", socialMediaId).Take(&detailSocial).Error

	if errDetail != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Error Update",
			"message": "Unable to find Photo",
		})
		return
	}

	library.GetDB().Table("social_media").Model(&detailSocial).Where("id = ?", socialMediaId).Updates(models.SocialMedia{
		Name:             updateReq.Name,
		Social_Media_Url: updateReq.Social_media_url,
		Updated_at:       time.Now(),
	})

	ctx.JSON(http.StatusOK, gin.H{
		"id":               detailSocial.Id,
		"name":             detailSocial.Name,
		"social_media_url": detailSocial.Social_Media_Url,
		"user_id":          detailSocial.User_Id,
		"updated_at":       detailSocial.Updated_at,
	})

}

func DeleteSocial(ctx *gin.Context) {
	socialMediaId := ctx.Param("socialMediaId")

	errDelete := library.GetDB().Table("social_media").Where("id = ?", socialMediaId).Delete(models.SocialMedia{}).Error

	if errDelete != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Error Delete Social",
			"message": "Cannot find Social Media",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your social media has been successfully deleted",
	})
}
