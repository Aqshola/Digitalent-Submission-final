package controllers

import (
	"final-project/library"
	"final-project/models"
	"net/http"
	"time"

	"final-project/helpers"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RegisterReq struct {
	Age      int    `json:"age" validate:"required,lte=100,gte=8"`
	Email    string `json:"email" validate:"required,email,max=191"`
	Password string `json:"password" validate:"required,min=6"`
	Username string `json:"username" validate:"required,max=10"`
}

func Register(ctx *gin.Context) {
	var dataRegister models.User
	var registReq RegisterReq

	ctx.ShouldBindJSON(&registReq)
	valid, trans := helpers.Valid()
	err := valid.Struct(registReq)

	if err != nil {
		errs := err.(validator.ValidationErrors)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Error Validation",
			"message": errs.Translate(trans),
		})
		return
	}

	dataRegister = models.User{
		Age:        registReq.Age,
		Username:   registReq.Username,
		Email:      registReq.Email,
		Password:   helpers.HashPass(registReq.Password),
		Created_at: time.Now(),
	}

	errRegister := library.GetDB().Table("users").Create(&dataRegister).Error
	if errRegister != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err":     "Unable to create user",
			"message": errRegister,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"age":      dataRegister.Age,
		"email":    dataRegister.Email,
		"id":       dataRegister.Id,
		"username": dataRegister.Username,
	})
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func Login(ctx *gin.Context) {
	var loginReq LoginReq
	var dataLogin models.User

	ctx.ShouldBindJSON(&loginReq)

	valid, trans := helpers.Valid()
	err := valid.Struct(loginReq)

	if err != nil {
		errs := err.(validator.ValidationErrors)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Error Validation",
			"message": errs.Translate(trans),
		})
		return
	}

	errEmail := library.GetDB().Table("users").Where("email = ?", loginReq.Email).Take(&dataLogin).Error
	if errEmail != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"err":      "Error Login",
			"messsage": "Unable to find email",
		})
		return
	}

	comparePass := helpers.ComparePass(loginReq.Password, dataLogin.Password)
	if !comparePass {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"err":     "Error Login",
			"message": "Wrong Password",
		})
		return
	}

	_token := helpers.GenerateJWT(dataLogin.Id, dataLogin.Email)

	ctx.JSON(http.StatusOK, gin.H{
		"token": _token,
	})

}

type UpdateReq struct {
	Email    string `json:"email" validate:"required,email,max=191"`
	Username string `json:"username" validate:"required,max=10"`
}

func UpdateUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	var updateReq UpdateReq
	var detailUser models.User

	ctx.ShouldBindJSON(&updateReq)
	valid, trans := helpers.Valid()
	err := valid.Struct(updateReq)

	if err != nil {
		errs := err.(validator.ValidationErrors)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"err": errs.Translate(trans),
		})
		return
	}

	errDetail := library.GetDB().Table("users").Where("id = ?", userId).Take(&detailUser).Error

	if errDetail != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Error Update",
			"message": "Unable to find account",
		})
		return
	}

	library.GetDB().Table("users").Model(&detailUser).Where("id = ?", userId).Updates(models.User{
		Email:      updateReq.Email,
		Username:   updateReq.Username,
		Updated_at: time.Now(),
	})

	ctx.JSON(http.StatusOK, gin.H{
		"id":         detailUser.Id,
		"email":      detailUser.Email,
		"username":   detailUser.Username,
		"age":        detailUser.Age,
		"updated_at": detailUser.Updated_at,
	})

}

func DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("userId")

	errDelete := library.GetDB().Table("users").Where("id = ?", userId).Delete(models.User{}).Error

	if errDelete != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"err":     "Error Delete",
			"message": "Account not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Your account has been successfully deleted",
	})
}
