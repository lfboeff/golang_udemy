package controllers

import (
	"api_mod/src/autenticacao"
	"api_mod/src/database"
	"api_mod/src/modelos"
	"api_mod/src/repositorios"
	"api_mod/src/respostas"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CriarPublicacao adiciona uma nova publicação no banco de dados
func CriarPublicacao(w http.ResponseWriter, r *http.Request) {

	usuarioIDNoToken, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		respostas.Erro(w, http.StatusUnauthorized, err)
		return
	}

	bodyRequest, err := io.ReadAll(r.Body)
	if err != nil {
		respostas.Erro(w, http.StatusUnprocessableEntity, err)
		return
	}

	var publicacao modelos.Publicacao

	if err = json.Unmarshal(bodyRequest, &publicacao); err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	publicacao.AutorID = usuarioIDNoToken

	if err = publicacao.Preparar(); err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Conectar()
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorioPublicacoes := repositorios.NovoRepositorioDePublicacoes(db)

	publicacao.ID, err = repositorioPublicacoes.Criar(publicacao)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusCreated, publicacao)
}

// BuscarPublicaco traz as publicações que aparecem no feed do usuário
func BuscarPublicacoes(w http.ResponseWriter, r *http.Request) {

	usuarioIDNoToken, err := autenticacao.ExtrairUsuarioID(r)
	if err != nil {
		respostas.Erro(w, http.StatusUnauthorized, err)
		return
	}

	db, err := database.Conectar()
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorioPublicacoes := repositorios.NovoRepositorioDePublicacoes(db)

	publicacoes, err := repositorioPublicacoes.Buscar(usuarioIDNoToken)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusCreated, publicacoes)

}

// BuscarPublicacao traz uma única publicação do feed do usuário
func BuscarPublicacao(w http.ResponseWriter, r *http.Request) {

	parameters := mux.Vars(r)
	publicacaoID, err := strconv.ParseUint(parameters["publicacaoId"], 10, 64)
	if err != nil {
		respostas.Erro(w, http.StatusBadRequest, err)
		return
	}

	db, err := database.Conectar()
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}
	defer db.Close()

	repositorioPublicacoes := repositorios.NovoRepositorioDePublicacoes(db)

	publicacao, err := repositorioPublicacoes.BuscarPorID(publicacaoID)
	if err != nil {
		respostas.Erro(w, http.StatusInternalServerError, err)
		return
	}

	respostas.JSON(w, http.StatusOK, publicacao)
}

// AtualizarPublicacao altera os dados de uma publicação
func AtualizarPublicacao(w http.ResponseWriter, r *http.Request) {

}

// DeletarPublicacao remove uma publicação
func DeletarPublicacao(w http.ResponseWriter, r *http.Request) {

}
