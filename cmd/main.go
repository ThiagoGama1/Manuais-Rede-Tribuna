package main

import (
	"manual/internal/data"
	"manual/internal/handler"

	"github.com/gin-gonic/gin"
)

//go mod init manual
func main(){
	data.LoadDatabase()
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/uploads", "./uploads")
	router.GET("/", handler.ListarManuais)
	router.GET("/manual/:id", handler.ExibirManualPorId)
	router.GET("/novo", func(c *gin.Context) {
    	c.HTML(200, "form.html", nil)
	})
	router.POST("/criar", handler.CriarManual)
	router.POST("/deletar/:id", handler.DeleteManualById)
	router.GET("/editar/:id", handler.UpdateManualById)
	router.POST("/atualizar/:id", handler.RecebeUpdateById)
	router.POST("/deletar-anexo/:id", handler.DeleteAnexoHandler)
	router.POST("/reordenar-anexo/:id/:direcao", handler.ReordenarAnexo)
	router.Run(":3609")
}