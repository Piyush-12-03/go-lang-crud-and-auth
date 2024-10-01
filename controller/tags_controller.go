package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/go-project/data/request"
	"example.com/go-project/data/response"
	"example.com/go-project/helper"
	"example.com/go-project/services"
	"github.com/gin-gonic/gin"
)

type TagsController struct {
	tagsService services.TagsService
}

func NewTagsController(service services.TagsService) *TagsController {
	return &TagsController{
		tagsService: service,
	}
}

func (controller *TagsController) Create(ctx *gin.Context) {
	createTagsRequest := request.CreteTagsRequest{}
	err := ctx.ShouldBindJSON(&createTagsRequest)
	if err != nil {
		helper.ErrorPanic(err)
		return
	}

	err = controller.tagsService.Create(createTagsRequest)
	if err != nil {
		helper.ErrorPanic(err)
		return
	}

	webresponse := response.Response{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   createTagsRequest,
		Msg:    "Tag added successfully.",
	}
	ctx.JSON(http.StatusOK, webresponse)
}

func (controller *TagsController) Update(ctxhttp *gin.Context) {
	updateTagsRequest := request.UpdateTagsRequest{}
	err := ctxhttp.ShouldBindJSON(&updateTagsRequest)
	helper.ErrorPanic(err)

	tagId := ctxhttp.Param("tagId")
	id, err := strconv.Atoi(tagId)
	helper.ErrorPanic(err)
	updateTagsRequest.Id = id

	controller.tagsService.Update(updateTagsRequest)
	webresponse := response.Response{
		Code:   http.StatusOK,
		Status: "ok",
		Data: map[string]interface{}{
			"id":   id,
			"name": updateTagsRequest.Name,
		},
	}
	ctxhttp.Header("Content-type", "application/json")
	ctxhttp.JSON(http.StatusOK, webresponse)
}

func (controller *TagsController) Delete(ctxhttp *gin.Context) {
	tagId := ctxhttp.Param("tagId")
	id, err := strconv.Atoi(tagId)
	if err != nil {
		// Handle the error if the conversion fails
		webresponse := response.Response{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   fmt.Sprintf("Invalid tag ID: %s", tagId),
		}
		ctxhttp.JSON(http.StatusBadRequest, webresponse)
		return
	}

	// Call FindById and handle the two return values
	tag, err := controller.tagsService.FindById(id)
	if err != nil {
		// Tag not found
		webresponse := response.Response{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   fmt.Sprintf("Tag with id %d not found", id),
		}
		ctxhttp.JSON(http.StatusNotFound, webresponse)
		return
	}

	// Proceed to delete the tag
	err = controller.tagsService.Delete(id)
	if err != nil {
		// Handle potential error from the Delete method
		webresponse := response.Response{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   "Failed to delete the tag",
		}
		ctxhttp.JSON(http.StatusInternalServerError, webresponse)
		return
	}

	webresponse := response.Response{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   fmt.Sprintf("Deleted tag with id %d", tag.Id),
	}
	ctxhttp.JSON(http.StatusOK, webresponse)
}

func (controller *TagsController) FindById(ctxhttp *gin.Context) {
	tagId := ctxhttp.Param("tagId")
	id, err := strconv.Atoi(tagId)
	if err != nil {
		// Handle the error if the conversion fails
		webresponse := response.Response{
			Code:   http.StatusBadRequest,
			Status: "error",
			Data:   fmt.Sprintf("Invalid tag ID: %s", tagId),
		}
		ctxhttp.JSON(http.StatusBadRequest, webresponse)
		return
	}

	// Call FindById and handle the two return values
	tagResponse, err := controller.tagsService.FindById(id)
	if err != nil {
		// Tag not found
		webresponse := response.Response{
			Code:   http.StatusNotFound,
			Status: "error",
			Data:   fmt.Sprintf("Tag with id %d not found", id),
		}
		ctxhttp.JSON(http.StatusNotFound, webresponse)
		return
	}

	// Prepare successful response
	webresponse := response.Response{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   tagResponse,
	}
	ctxhttp.JSON(http.StatusOK, webresponse)
}

func (controller *TagsController) FindAll(ctx *gin.Context) {
	// Default pagination values if not provided
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("pageSize"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Call the service to fetch tags
	tagResponse, err := controller.tagsService.FindAll(pageSize, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.PaginatedResponse{
			Code:   http.StatusInternalServerError,
			Status: "error",
			Data:   nil,
			Limit:  pageSize,
			Offset: offset,
			Msg:    "Failed to fetch tags: " + err.Error(),
		})
		log.Println("Error fetching tags:", err)
		return
	}

	// Prepare the successful response
	webresponse := response.PaginatedResponse{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   tagResponse,
		Limit:  pageSize,
		Offset: offset,
		Msg:    "Fetched tags successfully.",
	}
	ctx.Header("Content-type", "application/json")
	ctx.JSON(http.StatusOK, webresponse)
	log.Println("Fetched tags successfully", webresponse)
}

// func (controller *TagsController) FindPaginated(ctxhttp *gin.Context) {
// 	page, _ := strconv.Atoi(ctxhttp.DefaultQuery("page", "1"))
// 	pageSize, _ := strconv.Atoi(ctxhttp.DefaultQuery("pageSize", "10"))

// 	tagResponse := controller.tagsService.FindPaginated(page, pageSize)
// 	webresponse := response.Response{
// 		Code:   http.StatusOK,
// 		Status: "ok",
// 		Data:   tagResponse,
// 		Msg:    "Fetched paginated tags successfully.",
// 	}
// 	ctxhttp.Header("Content-type", "application/json")
// 	ctxhttp.JSON(http.StatusOK, webresponse)
// }

// // Method for sorting tags by name
// func (controller *TagsController) FindAllSorted(ctxhttp *gin.Context) {
// 	sortBy := ctxhttp.DefaultQuery("sortBy", "id")
// 	order := ctxhttp.DefaultQuery("order", "asc")

// 	tagResponse := controller.tagsService.FindAllSorted(sortBy, order)
// 	webresponse := response.Response{
// 		Code:   http.StatusOK,
// 		Status: "ok",
// 		Data:   tagResponse,
// 		Msg:    "Fetched sorted tags successfully.",
// 	}
// 	ctxhttp.Header("Content-type", "application/json")
// 	ctxhttp.JSON(http.StatusOK, webresponse)
// }

// // Method for filtering tags by custom criteria (e.g., name starts with a specific letter)
// func (controller *TagsController) FindByCustomFilter(ctxhttp *gin.Context) {
// 	startsWith := ctxhttp.DefaultQuery("startsWith", "")

// 	tagResponse := controller.tagsService.FindByCustomFilter(startsWith)
// 	webresponse := response.Response{
// 		Code:   http.StatusOK,
// 		Status: "ok",
// 		Data:   tagResponse,
// 		Msg:    "Fetched filtered tags successfully.",
// 	}
// 	ctxhttp.Header("Content-type", "application/json")
// 	ctxhttp.JSON(http.StatusOK, webresponse)
// }
