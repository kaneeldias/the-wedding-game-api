package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"the-wedding-game-api/middleware"
	"the-wedding-game-api/middleware/validators"
	"the-wedding-game-api/models"
	"the-wedding-game-api/types"
)

func GetChallengeById(c *gin.Context) {
	id, err := validators.ValidateGetChallengeByIdRequest(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	challenge, err := models.GetChallengeByID(id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := types.ChallengeCreatedResponse{
		Id:          challenge.ID,
		Name:        challenge.Name,
		Description: challenge.Description,
		Points:      challenge.Points,
		Image:       challenge.Image,
		Status:      challenge.Status,
		Type:        challenge.Type,
	}
	c.IndentedJSON(http.StatusOK, response)
	return
}

func CreateChallenge(c *gin.Context) {
	if err := middleware.CheckIsAdmin(c); err != nil {
		_ = c.Error(err)
		return
	}

	challengeRequest, err := validators.ValidateCreateChallengeRequest(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	createdChallenge, err := models.CreateNewChallenge(challengeRequest)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := types.ChallengeCreatedResponse{
		Id:          createdChallenge.ID,
		Name:        createdChallenge.Name,
		Description: createdChallenge.Description,
		Points:      createdChallenge.Points,
		Image:       createdChallenge.Image,
		Status:      createdChallenge.Status,
		Type:        createdChallenge.Type,
	}
	c.IndentedJSON(http.StatusCreated, response)
	return
}

func GetAllChallenges(c *gin.Context) {
	challengesArr, err := models.GetAllChallenges()
	if err != nil {
		_ = c.Error(err)
		return
	}

	var response types.GetChallengesResponse
	for _, challenge := range challengesArr {
		response.Challenges = append(response.Challenges, types.ChallengeCreatedResponse{
			Id:          challenge.ID,
			Name:        challenge.Name,
			Description: challenge.Description,
			Points:      challenge.Points,
			Image:       challenge.Image,
			Status:      challenge.Status,
			Type:        challenge.Type,
		})
	}

	c.IndentedJSON(http.StatusOK, response)
	return
}

func VerifyAnswer(c *gin.Context) {
	if err := middleware.CheckIsLoggedIn(c); err != nil {
		_ = c.Error(err)
		return
	}

	challengeId, verifyAnswerRequest, err := validators.ValidateVerifyAnswerRequest(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	correct, err := models.VerifyAnswer(challengeId, verifyAnswerRequest.Answer)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !correct {
		response := types.VerifyAnswerResponse{Correct: false}
		c.IndentedJSON(http.StatusOK, response)
		return
	}

	user, err := middleware.GetCurrentUser(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	submission := models.NewSubmission(user.ID, challengeId, verifyAnswerRequest.Answer)
	_, err = submission.Save()
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := types.VerifyAnswerResponse{Correct: true}
	c.IndentedJSON(http.StatusOK, response)
	return
}
