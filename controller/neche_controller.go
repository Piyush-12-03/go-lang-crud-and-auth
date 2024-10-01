package controller

import (
	"net/http"
	"strconv"

	"example.com/go-project/data/request"
	"example.com/go-project/data/response"
	"example.com/go-project/helper"
	"example.com/go-project/services"
	"github.com/gin-gonic/gin"
)

type NecheController struct {
	necheService services.NecheService
}

func NewNecheController(service services.NecheService) *NecheController {
	return &NecheController{
		necheService: service,
	}
}

// Create Neche
func (controller *NecheController) Create(ctx *gin.Context) {
	createNecheRequest := request.CreateNecheRequest{}
	err := ctx.ShouldBindJSON(&createNecheRequest)
	if err != nil {
		helper.ErrorPanic(err)
		return
	}

	if createNecheRequest.TagID == 0 {
		helper.ErrorPanic(err)
		return
	}

	err = controller.necheService.Create(createNecheRequest)
	if err != nil {
		helper.ErrorPanic(err)
		return
	}

	webresponse := response.Response{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   createNecheRequest,
		Msg:    "Neche Created Successfully",
	}
	ctx.JSON(http.StatusOK, webresponse)
}

// Find all Neches
func (controller *NecheController) FindAll(ctx *gin.Context) {
	neches, err := controller.necheService.FindAll()
	if err != nil {
		helper.ErrorPanic(err)
		return
	}

	webresponse := response.Response{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   neches,
		Msg:    "Neches Fetched Successfully.",
	}
	ctx.JSON(http.StatusOK, webresponse)
}

// Find Neche by ID
func (controller *NecheController) FindById(ctx *gin.Context) {
	necheIdStr := ctx.Param("necheId")
	necheId, err := strconv.Atoi(necheIdStr)
	if err != nil {
		helper.ErrorPanic(err)
		return
	}

	neche, err := controller.necheService.FindById(necheId)
	if err != nil {
		helper.ErrorPanic(err)
		return
	}

	webresponse := response.Response{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   neche,
	}
	ctx.JSON(http.StatusOK, webresponse)
}

// Delete Neche by ID
func (controller *NecheController) Delete(ctx *gin.Context) {
	necheIdStr := ctx.Param("necheId")
	necheId, err := strconv.Atoi(necheIdStr)
	if err != nil {
		helper.ErrorPanic(err)
		return
	}

	err = controller.necheService.Delete(necheId)
	if err != nil {
		helper.ErrorPanic(err)
		return
	}

	webresponse := response.Response{
		Code:   http.StatusOK,
		Status: "ok",
		Msg:    "Neche Deleted Successfully",
	}
	ctx.JSON(http.StatusOK, webresponse)
}
