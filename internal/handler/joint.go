package handler

import (
	"chow/internal/model"
	"chow/internal/service"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JointHandler struct {
	jointService *service.JointService
}

func NewJointHandler(jointService *service.JointService) *JointHandler {
	return &JointHandler{
		jointService: jointService,
	}
}

// CreateJoint godoc
// @Summary Create new joint
// @Description Create a new food joint
// @Tags joints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateJointReq true "Joint details"
// @Success 201 {object} model.SuccessResponse{data=model.Joint} "Joint added successfully"
// @Failure 401 {object} model.ErrorResponse "User not authenticated"
// @Failure 409 {object} model.ErrorResponse "Joint already exists"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /joints [post]
func (h *JointHandler) CreateJoint(c *gin.Context) {
	// validate data
	var req model.CreateJointReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Failed to validate data", Detail: err.Error()})
		return
	}
	// get user info from auth context
	user, ok := GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "User not authenticated"})
		return
	}

	// save file if uploaded

	joint, err := h.jointService.CreateJoint(c.Request.Context(), &model.Joint{Name: req.Name, Latitude: req.Latitude, Longitude: req.Longitude, Description: req.Description, IsApproved: false, CreatorID: user.ID, PhotoURL: nil})
	if err != nil {
		if errors.Is(err, service.ErrAlreadyExist) {
			c.JSON(http.StatusConflict, model.ErrorResponse{Message: err.Error()})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to add new joint"})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{Message: "Joint added successfully", Data: joint})
}

// SearchJoints godoc
// @Summary Search joints
// @Description Search for joints by name or description
// @Tags joints
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} model.SuccessResponse{data=[]model.Joint} "Joints retrieved successfully"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /joints/search [get]
func (h *JointHandler) SearchJoints(c *gin.Context) {
	var pagination model.PaginationQuery
	if err := c.ShouldBind(&pagination); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Failed to validate pagination params", Detail: err.Error()})
		return
	}
	q := c.Query("q")
	offset, limit := pagination.GetOffsetAndLimit()

	joints, err := h.jointService.SearchForJointByNameOrDescription(c.Request.Context(), q, offset, limit)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to search joints"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Joints retrieved successfully", Data: joints})
}

// GetNearByJoints godoc
// @Summary Find nearby joints
// @Description Get food joints near specified coordinates within radius
// @Tags joints
// @Accept json
// @Produce json
// @Param request query model.NearbyJointsQuery true "Search parameters"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} model.SuccessResponse{data=[]model.Joint} "Nearby Joints retrieved successfully"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /joints/nearby [get]
func (h *JointHandler) GetNearByJoints(c *gin.Context) {
	var query struct {
		model.PaginationQuery
		model.NearbyJointsQuery
	}
	if err := c.ShouldBind(&query); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Failed to validate pagination params", Detail: err.Error()})
		return
	}
	offset, limit := query.GetOffsetAndLimit()

	joints, err := h.jointService.GetNearbyJoints(c.Request.Context(), model.Coordinate{Latitude: query.Latitude, Longitude: query.Longitude}, query.Radius, offset, limit)
	if err != nil {
		if errors.Is(err, service.ErrMaxSearchRadiusExceeded) {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to retrieve nearby joints"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Nearby Joints retrieved successfully", Data: joints})
}

// VoteJoint godoc
// @Summary Vote on a joint
// @Description Upvote or downvote a food joint
// @Tags joints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Joint ID"
// @Param request body model.VoteJointReq true "Vote direction"
// @Success 200 {object} model.SuccessResponse{data=model.Joint} "Vote recorded successfully"
// @Success 200 {object} model.SuccessResponse "Vote already recorded"
// @Failure 401 {object} model.ErrorResponse "User not authenticated"
// @Failure 404 {object} model.ErrorResponse "Joint not found"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /joints/{id}/vote [post]
func (h *JointHandler) VoteJoint(c *gin.Context) {
	// validate params and data
	var param model.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Failed to validate joint ID", Detail: err.Error()})
		return
	}

	var vote model.VoteJointReq
	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Invalid vote direction", Detail: err.Error()})
		return
	}

	// get current user details
	user, ok := GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: "User not authenticated"})
		return
	}

	joint, err := h.jointService.VoteForJoint(c.Request.Context(), param.GetID(), &model.Vote{
		UserID:    user.ID,
		JointID:   param.GetID(),
		Direction: vote.Direction,
	})

	if err != nil {
		if errors.Is(err, service.ErrJointNotFound) {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Message: "Joint not found"})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to process vote"})
		return
	}

	// if joint is null at this point, the user tried voting in the same direction again so a success response is simply returned
	if joint == nil {
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Vote already recorded"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Vote recorded successfully", Data: joint})
}

// GetAllJoints godoc
// @Summary Get all joints
// @Description Get all approved food joints with pagination
// @Tags joints
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} model.SuccessResponse{data=[]model.Joint} "Joints retrieved successfully"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /joints [get]
func (h *JointHandler) GetAllJoints(c *gin.Context) {
	var pagination model.PaginationQuery
	if err := c.ShouldBind(&pagination); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Failed to validate pagination params", Detail: err.Error()})
		return
	}
	offset, limit := pagination.GetOffsetAndLimit()

	joints, err := h.jointService.GetAllJoints(c.Request.Context(), offset, limit)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to retrieve joints"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Joints retrieved successfully", Data: joints})
}

// GetJoint godoc
// @Summary Get one joint
// @Description Get a single joint by ID
// @Tags joints
// @Accept json
// @Produce json
// @Param id path string true "Joint ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} model.SuccessResponse{data=model.Joint} "Joint retrieved successfully"
// @Failure 404 {object} model.ErrorResponse "Joint not found"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /joints/{id} [get]
func (h *JointHandler) GetJoint(c *gin.Context) {
	var param model.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Failed to validate joint ID", Detail: err.Error()})
		return
	}

	joint, err := h.jointService.GetJointByID(c.Request.Context(), param.GetID())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to retrieve joint"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Joint retrieved successfully", Data: joint})
}

// UpdateJoint godoc
// @Summary Update joint
// @Description Update an existing joint's details
// @Tags joints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Joint ID"
// @Param request body model.CreateJointReq true "Updated joint details"
// @Success 200 {object} model.SuccessResponse{data=model.Joint} "Joint updated successfully"
// @Failure 401 {object} model.ErrorResponse "User not authenticated"
// @Failure 403 {object} model.ErrorResponse "User not authorized"
// @Failure 404 {object} model.ErrorResponse "Joint not found"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /joints/{id} [patch]
func (h *JointHandler) UpdateJoint(c *gin.Context) {
	var param model.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Failed to validate joint ID", Detail: err.Error()})
		return
	}

	var req model.CreateJointReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Failed to validate data", Detail: err.Error()})
		return
	}

	// get current user. updates are restricted to only creator of joint and admin
	user, ok := GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: "User not authenticated"})
		return
	}

	// get existing joint to verify ownership
	existingJoint, err := h.jointService.GetJointByID(c.Request.Context(), param.GetID())
	if err != nil {
		if errors.Is(err, service.ErrJointNotFound) {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Message: "Joint not found"})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to update joint"})
		return
	}

	// verify owner details
	if existingJoint.CreatorID != user.ID && user.Role != model.Admin {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Message: "User not authorized to update this joint"})
		return
	}

	// update only name, location details and description
	joint, err := h.jointService.UpdateJointByID(c.Request.Context(), param.GetID(), &model.Joint{
		Name:        req.Name,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Description: req.Description,
		IsApproved:  existingJoint.IsApproved,
		CreatorID:   existingJoint.CreatorID,
		PhotoURL:    existingJoint.PhotoURL,
	})

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to update joint"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Joint updated successfully", Data: joint})
}

// DeleteJoint godoc
// @Summary Delete joint
// @Description Delete an existing joint (admin only)
// @Tags joints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Joint ID"
// @Success 200 {object} model.SuccessResponse "Joint deleted successfully"
// @Failure 401 {object} model.ErrorResponse "User not authenticated"
// @Failure 403 {object} model.ErrorResponse "User not authorized"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /joints/{id} [delete]
func (h *JointHandler) DeleteJoint(c *gin.Context) {
	var param model.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Invalid joint ID", Detail: err.Error()})
		return
	}

	// get admin info
	_, ok := GetCurrentAdmin(c)
	if !ok {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Message: "User not authorized to perform this action"})
		return
	}

	if err := h.jointService.DeleteJointByID(c.Request.Context(), param.GetID()); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to delete joint"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Joint deleted successfully"})
}
