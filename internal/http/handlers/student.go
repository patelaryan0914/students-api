package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/patelaryan0914/students-api/internal/storage"
	"github.com/patelaryan0914/students-api/internal/types"
	"github.com/patelaryan0914/students-api/internal/utils/response"
)

func New(storage  storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating Students")

		var student types.Student
		err:= json.NewDecoder(r.Body).Decode(&student)
		fmt.Print(student)
		if errors.Is(err,io.EOF){
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(fmt.Errorf("empty body")))
			return 
		}
		if err != nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
		}

		// request Validation
		if err:=validator.New().Struct(student);err!=nil{
			validateErrs:=err.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest,response.ValidationError(validateErrs))
			return 
		}
		lastId,err:=storage.CreateStudent(student.Name,student.Email,student.Age)
		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return
		}
		slog.Info("User created Successfully",slog.String("userId",fmt.Sprint(lastId)))
		response.WriteJson(w,http.StatusCreated,map[string]int64{
			"id":lastId,
		})
   }
}
func GetByID(storage  storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		id:=r.PathValue("id")
		slog.Info("Fetching Students",slog.String("id",id))
		intId, err := strconv.ParseInt(id,10,64)
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
		}
		student,err:=storage.GetStudentById(intId)
		if err!=nil{
			slog.Error("Error getting User",slog.String("id",id))
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return
		}
		response.WriteJson(w,http.StatusOK,student)
   }
}