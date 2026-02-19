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
	msg := c.Query("msg")
	c.HTML(http.StatusOK, "index.html", gin.H{"Manuais": data.GetManuais(), "Msg": msg})
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
	c.Redirect(http.StatusSeeOther, "/?msg=criado")
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
	c.Redirect(http.StatusFound, "/?msg=deletado")
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
func RecebeUpdateById(c *gin.Context){
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
	form, err := c.MultipartForm()
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	if strings.TrimSpace(inputTitulo) == "" || strings.TrimSpace(inputConteudo) == "" || strings.TrimSpace(inputSecao) == ""{
		c.Redirect(http.StatusFound, "/novo")
		return
	}
	

	for _, arquivo := range form.File["arquivos"]{
		destino := "./uploads/" + arquivo.Filename
		
		if err := c.SaveUploadedFile(arquivo, destino); err != nil{
			log.Println("Erro ao salvar o anexo:", err)
			continue
		}
		novoAnexo := models.Anexo{
            Nome:          arquivo.Filename,
            Tamanho_bytes: int(arquivo.Size),
            Caminho:       destino,
            Tipo_arquivo:  arquivo.Header.Get("Content-Type"),
            Manual_id:     id, 
	}
		err := data.InsertAnexo(novoAnexo)
        if err != nil {
            log.Println("Erro ao inserir anexo no banco:", err)
        }
	}
	copiaManual.ID = id
	copiaManual.Titulo = inputTitulo
	copiaManual.Conteudo = inputConteudo
	copiaManual.Secao = inputSecao
	
	err = data.UpdateManual(copiaManual)
	if err != nil {
        log.Println("Erro ao atualizar manual:", err) 
        c.Redirect(http.StatusSeeOther, "/")
        return
    }
	c.Redirect(http.StatusSeeOther, "/?msg=atualizado")
	}


func DeleteAnexoHandler(c *gin.Context){
	idAnexo := c.Param("id")
	id, err := strconv.Atoi(idAnexo)

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"mensagem": "erro ao converter id"})
		return
	}
	err = data.DeleteAnexo(id)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"mensagem": "erro ao chamar a funcao na data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"output": "Deletado com sucesso!"})
}

func ReordenarAnexo(c *gin.Context){
	idAnexo := c.Param("id")
	direcao := c.Param("direcao")

	id, err := strconv.Atoi(idAnexo)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"mensagem": "erro ao converter id"})
		return
	}

	err = data.ReordenarAnexoData(id, direcao)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"mensagem": "erro ao chamar a funcao na data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output":"Reordenado!"})
	
}
