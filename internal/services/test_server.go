package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vladimirbannikov/BIP-backend/internal/models"
	"github.com/vladimirbannikov/BIP-backend/internal/models/structs"
	"github.com/vladimirbannikov/BIP-backend/internal/utils/logger"
	"io"
	"net/http"
	"strconv"
)

// TODO: service layer
// todo: rating вернуть
// todo: get profile
// todo: get test
// todo: получить скор по тесту

type TestsModelManager interface {
	GetTests(ctx context.Context, limit int, offset int) ([]structs.TestSimple, error)
	GetFullTestByID(ctx context.Context, id int) (structs.TestFull, error)
	CheckUserAnswers(ctx context.Context, tokenStr string, testID int, userAnswers []structs.UserAnswer) (structs.TestResult, error)
	GetRating(ctx context.Context, category string, limit int, offset int) (structs.Rating, error)
}

type testsServer struct {
	t TestsModelManager
}

func (s *testsServer) GetAllTests(w http.ResponseWriter, req *http.Request) {
	data, status := s.getAllTests(req.Context())
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		return
	}
}

func (s *testsServer) getAllTests(ctx context.Context) ([]byte, int) {
	p, err := s.t.GetTests(ctx, 0, 0)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("GetTests: %v", err))
		return nil, http.StatusInternalServerError
	}

	tests := make([]GetAllTestsRespTest, 0)
	for _, t := range p {
		pic64 := base64.StdEncoding.EncodeToString(t.Picture)
		tests = append(tests, GetAllTestsRespTest{
			Id:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Category:    t.Category,
			DiffLevel:   t.DiffLevel,
			Picture:     pic64,
		})
	}

	articleJSON, _ := json.Marshal(
		GetAllTestsResp{
			Tests: tests,
		},
	)
	return articleJSON, http.StatusOK
}

func (s *testsServer) GetTest(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars[testIdKey])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, status := s.getTest(req.Context(), id)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func (s *testsServer) getTest(ctx context.Context, id int) ([]byte, int) {
	t, err := s.t.GetFullTestByID(ctx, id)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf(err.Error()))
		return nil, http.StatusInternalServerError
	}

	questions := make([]GetFullTestRespQuestion, 0)
	for _, q := range t.Questions {
		answers := make([]GetFullTestRespAnswer, 0)
		for _, a := range q.Answers {
			answers = append(answers, GetFullTestRespAnswer{
				Id:     a.ID,
				Answer: a.Answer,
			})
		}
		questions = append(questions, GetFullTestRespQuestion{
			Id:       q.ID,
			Question: q.Question,
			IsSong:   q.IsSong,
			Answers:  answers,
		})
	}

	articleJSON, _ := json.Marshal(
		GetFullTestResp{
			Id:          t.ID,
			Name:        t.Name,
			DiffLevel:   t.DiffLevel,
			Description: t.Description,
			Category:    t.Category,
			Questions:   questions,
		},
	)
	return articleJSON, http.StatusOK
}

func (s *testsServer) GetScore(w http.ResponseWriter, req *http.Request) {
	tokenStr := req.Header.Get(authHeaderStr)
	if tokenStr == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(req)
	testId, err := strconv.Atoi(vars[testIdKey])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("Service: GetScore: ReadAll error: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var unm GetUserScoreReq
	if err = json.Unmarshal(body, &unm); err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("Service: GetScore: json.Unmarshal error: %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, status := s.getScore(req.Context(), tokenStr, testId, unm)
	logger.Log(logger.InfoPrefix, fmt.Sprintf("Response: %v %s", status, data))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func (s *testsServer) getScore(ctx context.Context, tokenStr string, testID int, answers GetUserScoreReq) ([]byte, int) {
	userAnswers := make([]structs.UserAnswer, 0)
	for _, ans := range answers.Answers {
		userAnswers = append(userAnswers, structs.UserAnswer{
			QuestionID: ans.QuestionId,
			AnswerID:   ans.AnswerId,
		})
	}

	score, err := s.t.CheckUserAnswers(ctx, tokenStr, testID, userAnswers)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, http.StatusNotFound
		}
		if errors.Is(err, models.ErrBadCredentials) || errors.Is(err, models.ErrInvalidToken) {
			return nil, http.StatusUnauthorized
		}
		logger.Log(logger.ErrPrefix, fmt.Sprintf(err.Error()))
		return nil, http.StatusInternalServerError
	}

	articleJSON, _ := json.Marshal(
		GetUserScoreResp{
			UserScore: score.UserScore,
			Total:     score.Total,
		})
	return articleJSON, http.StatusOK
}

func (s *testsServer) GetRating(w http.ResponseWriter, req *http.Request) {
	data, status := s.getRating(req.Context())
	logger.Log(logger.InfoPrefix, fmt.Sprintf("Response: %v %s", status, data))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		return
	}
}

func (s *testsServer) getRating(ctx context.Context) ([]byte, int) {
	rating, err := s.t.GetRating(ctx, "", 0, 0)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf(err.Error()))
		return nil, http.StatusInternalServerError
	}

	ratingUnits := make([]GetRatingRespUnit, 0)
	for _, unit := range rating.Rating {
		ratingUnits = append(ratingUnits, GetRatingRespUnit{
			Login: unit.Login,
			Place: unit.Place,
			Score: unit.Score,
		})
	}

	articleJSON, _ := json.Marshal(
		GetRatingResp{
			Rating: ratingUnits,
		})
	return articleJSON, http.StatusOK
}
