package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/manzil-infinity180/backend-devops/pkg/db/models"
	"log"
	"net/http"
	"strconv"
)

func StartAPI(pgdb *pg.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger, middleware.WithValue("DB", pgdb))
	r.Route("/comment", func(r chi.Router) {
		r.Post("/", createComment)
		r.Get("/", getComments)
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("up and running"))
	})
	r.Get("/{commentID}", getCommentByID)
	r.Put("/{commentID}", updateCommentByID)
	r.Delete("/{commentID}", deleteCommentByID)
	return r
}

type CreateCommentRequest struct {
	Comment string `json:"comment"`
	UserID  int64  `json:"user_id"`
}

type CommentsResponse struct {
	Error    string            `json:"error"`
	Comments []*models.Comment `json:"comments"`
	Success  bool              `json:"success"`
}

type CommentResponse struct {
	Error   string          `json:"error"`
	Comment *models.Comment `json:"comment"`
	Success bool            `json:"success"`
}

func handleErr(w http.ResponseWriter, err error) {
	res := &CommentResponse{
		Success: false,
		Error:   err.Error(),
		Comment: nil,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error sending response %v\n", err)
	}
	w.WriteHeader(http.StatusBadRequest)
}
func handleDBFromContextErr(w http.ResponseWriter) {
	res := &CommentResponse{
		Success: false,
		Error:   "could not get the DB from context",
		Comment: nil,
	}
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("error sending response %v\n", err)
	}
	w.WriteHeader(http.StatusBadRequest)

}
func createComment(w http.ResponseWriter, r *http.Request) {
	req := &CreateCommentRequest{}
	err := json.NewDecoder(r.Body).Decode(req)

	if err != nil {
		handleErr(w, err)
		return
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		handleDBFromContextErr(w)
		return
	}
	comment, err := models.CreateComment(pgdb, &models.Comment{
		Comment: req.Comment,
		UserID:  req.UserID,
	})
	if err != nil {
		handleErr(w, err)
		return
	}
	res := &CommentResponse{
		Success: true,
		Error:   "",
		Comment: comment,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func getComments(w http.ResponseWriter, r *http.Request) {
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		handleDBFromContextErr(w)
		return
	}
	comments, err := models.GetComments(pgdb)
	if err != nil {
		handleErr(w, err)
		return
	}
	res := &CommentsResponse{
		Success:  true,
		Error:    "",
		Comments: comments,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

func getCommentByID(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "commentID")
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		handleDBFromContextErr(w)
		return
	}
	comment, err := models.GetComment(pgdb, commentID)
	if err != nil {
		handleErr(w, err)
		return
	}
	res := &CommentResponse{
		Success: true,
		Error:   "",
		Comment: comment,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

type CommentRequest struct {
	Comment string `json:"comment"`
	UserID  int64  `json:"user_id"`
}

func updateCommentByID(w http.ResponseWriter, r *http.Request) {
	req := &CommentRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		handleErr(w, err)
		return
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		handleDBFromContextErr(w)
		return
	}
	//get the commentID to know what comment to modify
	commentID := chi.URLParam(r, "commentID")
	//we get a string but we need to send an int so we convert it
	intCommentID, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		handleErr(w, err)
		return
	}
	comment, err := models.UpdateComment(pgdb, &models.Comment{
		ID:      intCommentID,
		Comment: req.Comment,
		UserID:  req.UserID,
	})
	if err != nil {
		handleErr(w, err)
	}
	res := &CommentResponse{
		Success: true,
		Error:   "",
		Comment: comment,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("error encoding response: %v", err)
	}

}

func deleteCommentByID(w http.ResponseWriter, r *http.Request) {
	req := CommentRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		handleErr(w, err)
		return
	}
	pgdb, ok := r.Context().Value("DB").(*pg.DB)
	if !ok {
		handleDBFromContextErr(w)
		return
	}
	commentID := chi.URLParam(r, "commentID")
	intCommentID, err := strconv.ParseInt(commentID, 10, 64)
	if err != nil {
		handleErr(w, err)
		return
	}
	err = models.DeleteComment(pgdb, intCommentID)
	if err != nil {
		handleErr(w, err)
	}
	res := &CommentResponse{
		Success: true,
		Error:   "",
		Comment: nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
