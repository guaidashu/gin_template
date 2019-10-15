package controller

import (
	"gin_template/app/libs"
	"gin_template/app/models"
	"github.com/gin-gonic/gin"
)

func InitTables(ctx *gin.Context) {
	models.CreateTable()
	libs.Success(ctx, "create table")
}
