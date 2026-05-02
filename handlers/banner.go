package handlers

import (
	"Zhooze/usecase"
	"Zhooze/utils/models"
	"Zhooze/utils/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddBanner(c *gin.Context) {
	var banner models.BannerRequest
	banner.Title1 = c.PostForm("title1")
	banner.Title2 = c.PostForm("title2")
	banner.Subtitle1 = c.PostForm("subtitle1")
	banner.Subtitle2 = c.PostForm("subtitle2")
	banner.Link = c.PostForm("link")

	file, err := c.FormFile("image")
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "image is required", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	res, err := usecase.AddBanner(banner, file)
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to add banner", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	// Save file to disk
	if err := c.SaveUploadedFile(file, res.Image); err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to save image", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusCreated, "Banner added successfully", res, nil)
	c.JSON(http.StatusCreated, successRes)
}

func GetBannersAdmin(c *gin.Context) {
	banners, err := usecase.GetBannersForAdmin()
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch banners", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Banners fetched successfully", banners, nil)
	c.JSON(http.StatusOK, successRes)
}

func GetBannersUser(c *gin.Context) {
	banners, err := usecase.GetBannersForUser()
	if err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to fetch banners", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}
	successRes := response.ClientResponse(http.StatusOK, "Banners fetched successfully", banners, nil)
	c.JSON(http.StatusOK, successRes)
}

func DeleteBanner(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	if err := usecase.DeleteBanner(id); err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to delete banner", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Banner deleted successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}

func ToggleBannerStatus(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errRes := response.ClientResponse(http.StatusBadRequest, "invalid id", nil, err.Error())
		c.JSON(http.StatusBadRequest, errRes)
		return
	}

	if err := usecase.ToggleBannerStatus(id); err != nil {
		errRes := response.ClientResponse(http.StatusInternalServerError, "failed to toggle status", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errRes)
		return
	}

	successRes := response.ClientResponse(http.StatusOK, "Banner status toggled successfully", nil, nil)
	c.JSON(http.StatusOK, successRes)
}
