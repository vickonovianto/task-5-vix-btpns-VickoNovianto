package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"user-photo-api/database"
	"user-photo-api/helpers"
	"user-photo-api/models"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

const IDENTITY_KEY = "ID"

func GetJwtMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(os.Getenv("SECRET")),
		Timeout:     3 * 24 * time.Hour,
		MaxRefresh:  3 * 24 * time.Hour,
		IdentityKey: IDENTITY_KEY,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if user, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					IDENTITY_KEY: strconv.Itoa(user.ID),
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			tokenUserIdString := claims[IDENTITY_KEY].(string)
			c.Set("tokenUserIdString", tokenUserIdString)
			return &models.User{}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var req models.UserLoginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			req.Trim()
			_, err := govalidator.ValidateStruct(req)
			if err != nil {
				return "", err
			}
			user := new(models.User)
			copier.Copy(user, req)
			if err := database.Database().
				Where("email = ?", user.Email).
				First(user).Error; err != nil {
				return "", jwt.ErrFailedAuthentication
			}
			hashedPassword := user.Password
			if !helpers.IsPasswordMatch(hashedPassword, req.Password) {
				return "", jwt.ErrFailedAuthentication
			}
			c.Set("tokenUserIdString", strconv.Itoa(user.ID))
			return &models.User{
				ID: user.ID,
			}, nil
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			tokenUserIdAny, exists := c.Get("tokenUserIdString")
			if !exists {
				helpers.SendResponseError(c, http.StatusBadRequest, errors.New("ID not found in context"))
				return
			}
			tokenUserIdString := fmt.Sprintf("%v", tokenUserIdAny)
			tokenUserIdInt, err := strconv.Atoi(tokenUserIdString)
			if err != nil {
				helpers.SendResponseError(c, http.StatusBadRequest, err)
				return
			}

			userLoginResponse := new(models.UserLoginResponse)
			userLoginResponse.ID = tokenUserIdInt
			userLoginResponse.Token = token
			userLoginResponse.TokenExpire = expire.Format(time.RFC3339)
			helpers.SendResponseSuccess(c, userLoginResponse)
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			helpers.SendResponseError(c, code, errors.New(message))
		},
	})
}
