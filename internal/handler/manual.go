package handler

import (
	"log"
	"manual/internal/data"
	"manual/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ListarManuais(c *gin.Context){
	c.HTML(http.StatusOK, "index.html", gin.H{"Manuais": data.GetManuais()})
}

func ExibirManualPorId(c *gin.Context){
	idManual := c.Param("id")
	id, err := strconv.Atoi(idManual)
	var copiaManual models.Manual
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	copiaManual, err = data.GetManualByID(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	c.HTML(http.StatusOK, "visualizar.html", copiaManual)
}

func CriarManual(c *gin.Context){
	//novoId := len(data.Manuais) + 1
	novoId := 0
	inputTitulo := c.PostForm("titulo")
	inputConteudo := c.PostForm("conteudo")
	inputSecao := c.PostForm("secao") 
	form, err := c.MultipartForm()
	var novoManual models.Manual
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	if strings.TrimSpace(inputTitulo) == "" || strings.TrimSpace(inputConteudo) == "" || strings.TrimSpace(inputSecao) == ""{
		c.Redirect(http.StatusFound, "/novo")
		return
	}
	var listaAnexos []models.Anexo
	for _, arquivo := range form.File["arquivos"]{
		destino := "./uploads/" + arquivo.Filename
		
		if err := c.SaveUploadedFile(arquivo, destino); err != nil{
			log.Println("Erro ao salvar o arquivo:", err)
			continue
		}

		novoAnexo := models.Anexo{
			Nome: arquivo.Filename,
			Tamanho_bytes: int(arquivo.Size),
			Caminho: destino,
			Tipo_arquivo: arquivo.Header.Get("Content-Type"),
		}
		listaAnexos = append(listaAnexos, novoAnexo)
	}
	
	novoManual = models.Manual{ID: novoId, Titulo: inputTitulo, Conteudo: inputConteudo, Secao: inputSecao, Arquivos: listaAnexos}
	data.InsertManual(novoManual)
	c.Redirect(http.StatusSeeOther, "/")
}

func DeleteManualById(c *gin.Context){
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	err = data.DeleteManual(id)

	if err != nil{
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	c.Redirect(http.StatusFound, "/")
}

func UpdateManualById(c *gin.Context){
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil{
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	var copiaManual models.Manual

	copiaManual, err = data.GetManualByID(id)
	if err != nil{
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	c.HTML(http.StatusOK, "editar.html", copiaManual)
	
}
func RecebeUpdateById(c * gin.Context){
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil{
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	var copiaManual models.Manual

	inputTitulo := c.PostForm("titulo")
	inputConteudo := c.PostForm("conteudo")
	inputSecao := c.PostForm("secao")

	copiaManual.ID = id
	copiaManual.Titulo = inputTitulo
	copiaManual.Conteudo = inputConteudo
	copiaManual.Secao = inputSecao

	err = data.UpdateManual(copiaManual)
	if err != nil {
        log.Println("Erro ao atualizar manual:", err) // <--- O Dedo-Duro
        c.Redirect(http.StatusSeeOther, "/")
        return
    }
	c.Redirect(http.StatusSeeOther, "/")
}