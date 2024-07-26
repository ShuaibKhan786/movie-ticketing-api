package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/config"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/models"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/database"
	redisdb "github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/redis"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/services/security"
	"github.com/ShuaibKhan786/movie-ticketing-api/pkg/utils"
)

//TODO:
//	1.check if hall register or not
//	2.if hall registered continue the movie registration
//	3. registered the movie : movie_id
//	4. registered the movie_show using movie_id , hall_id
//	5. registered the movie_show_timing using movie_show_id
//	6. registered the actor, actress , director, producer : thier id
//	7. registered movie_actrees, movie_actor
//	8. 			movie_director, movie_producer : thier id and movie_id

//REDIS KEY FORMAT (for hall registration) hall:registered:admin:admin_id

func ShowRegister(w http.ResponseWriter, r *http.Request) {
	claims,  ok := r.Context().Value(config.ClaimsContextKey).(security.Claims)
	if !ok {
		utils.JSONResponse(&w, "invalid claims", http.StatusBadRequest)
		return
	}
	
	hallId, err := isHallRegistered(claims)
	if err != nil {
		utils.JSONResponse(&w, "hall not registered", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.JSONResponse(&w, "failed to read the body", http.StatusInternalServerError)
		return
	}

	if !utils.IsValidJson(body) {
		utils.JSONResponse(&w, "invalid payload", http.StatusBadRequest)
		return
	}

	var show models.Show
	if err := utils.DecodeJson(body, &show); err != nil{
		utils.JSONResponse(&w, "failed to encode the json payload", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := database.RegisterShow(ctx, show, hallId); err!= nil {
		utils.JSONResponse(&w, "failed to registered the show in database "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JSONResponse(&w, "successfully registered the show", http.StatusCreated)
}


func isHallRegistered(claims security.Claims) (int64 , error) {
	var hallId int64
	
	redisCtx, redisCancel := context.WithTimeout(context.Background(), 500 * time.Millisecond) 
	defer redisCancel()
	
	redisKey := fmt.Sprintf("hall:registered:%s:%d",claims.Role,claims.Id)

	redisValue, err := redisdb.Get(redisCtx, redisKey)
	if err != nil {
		return hallId ,fmt.Errorf("redis error : %w", err)
	}

	if redisValue == "" {
		hallId, err = database.GetId("hall", "admin_id", claims.Id)
		if err != nil {
			return hallId, fmt.Errorf("database error : %w", err)
		}

		redisCtx, redisCancel = context.WithTimeout(context.Background(), 500 * time.Millisecond) 
		defer redisCancel()

		if err := redisdb.Set(redisCtx, redisKey, hallId, config.RedisZeroExpirationTime); err != nil {
			return hallId, 	nil
		}

		return hallId, nil	
	}
	
	hallId, err = strconv.ParseInt(redisValue, 10, 64)
	if err != nil {
		return hallId, fmt.Errorf("conversion error : %w",err)
	}
	return hallId, nil
}
