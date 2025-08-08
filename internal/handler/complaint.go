package handler

import (
	"chow/internal/model"
	"chow/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ComplaintHandler struct {
	complaintService *service.ComplaintService
}

func NewComplaintHandler(complaintService *service.ComplaintService) *ComplaintHandler {
	return &ComplaintHandler{
		complaintService: complaintService,
	}
}

// GetAllComplaints godoc
// @Summary Get all complaints
// @Description Get all complaints with pagination. Admins only Endpoint
// @Tags complaints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} model.SuccessResponse{data=[]model.Complaint} "Complaints retrieved successfully"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /complaints [get]
func (h *ComplaintHandler) GetAllComplaints(c *gin.Context) {
	var pagination model.PaginationQuery
	if err := c.ShouldBind(&pagination); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Failed to validate pagination params", Detail: err.Error()})
		return
	}
	offset, limit := pagination.GetOffsetAndLimit()

	// verify info from auth context
	_, ok := h.getAuthAdminOrModerator(c)
	if !ok {
		return
	}

	complaints, err := h.complaintService.GetAllComplaints(c.Request.Context(), offset, limit)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to retrieve complaints"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Complaints retrieved successfully", Data: complaints})
}

// GetComplaint godoc
// @Summary Get one complaint
// @Description Get a single complaint by ID
// @Tags complaints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Complaint ID"
// @Success 200 {object} model.SuccessResponse{data=model.Complaint} "Complaint retrieved successfully"
// @Failure 404 {object} model.ErrorResponse "Complaint not found"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /complaints/{id} [get]
func (h *ComplaintHandler) GetComplaint(c *gin.Context) {
	var param model.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Failed to validate complaint ID", Detail: err.Error()})
		return
	}

	// verify info from auth context
	if _, ok := h.getAuthUser(c); !ok {
		return
	}

	complaint, err := h.complaintService.GetComplaintByID(c.Request.Context(), param.GetID())
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to retrieve complaint"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Complaint retrieved successfully", Data: complaint})
}

// GetUserComplaints godoc
// @Summary Get user complaints
// @Description Get all complaints made by the authenticated user
// @Tags complaints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {object} model.SuccessResponse{data=[]model.Complaint} "User complaints retrieved successfully"
// @Failure 401 {object} model.ErrorResponse "User not authenticated"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /complaints/me [get]
func (h *ComplaintHandler) GetUserComplaints(c *gin.Context) {
	var pagination model.PaginationQuery
	if err := c.ShouldBind(&pagination); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{
			Message: "Failed to validate pagination params",
			Detail:  err.Error(),
		})
		return
	}

	// Get authenticated user
	user, ok := h.getAuthUser(c)
	if !ok {
		return
	}

	offset, limit := pagination.GetOffsetAndLimit()
	complaints, err := h.complaintService.GetUserComplaints(c.Request.Context(), user.ID, offset, limit)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: "Failed to retrieve user complaints",
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Message: "User complaints retrieved successfully",
		Data:    complaints,
	})
}

// ResolveComplaint godoc
// @Summary Resolve complaint
// @Description Resolves a specific complaint. Admins-only Endpoint
// @Tags complaints
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Complaint ID"
// @Success 200 {object} model.SuccessResponse{data=model.Complaint} "Complaint updated successfully"
// @Failure 401 {object} model.ErrorResponse "User not authenticated"
// @Failure 403 {object} model.ErrorResponse "User not authorized"
// @Failure 404 {object} model.ErrorResponse "Complaint not found"
// @Failure 422 {object} model.ErrorResponse "Validation error"
// @Failure 500 {object} model.ErrorResponse "Server error"
// @Router /complaints/{id}/resolve [patch]
func (h *ComplaintHandler) ResolveComplaint(c *gin.Context) {
	var param model.IDParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.ErrorResponse{Message: "Failed to validate complaint ID", Detail: err.Error()})
		return
	}

	// verify info from auth context
	if _, ok := h.getAuthAdminOrModerator(c); !ok {
		return
	}

	complaint, err := h.complaintService.UpdateComplaintStatusByID(c.Request.Context(), param.GetID(), model.ResolvedComplaint)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to resolve complaint"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Complaint updated successfully", Data: complaint})
}

// getAuthUser retrieves the authenticated user or return an unauthorized error if user is not present
func (h *ComplaintHandler) getAuthUser(c *gin.Context) (*model.AuthenticatedUser, bool) {
	// get user info from auth context
	user, ok := GetCurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: "User not authenticated"})
		return nil, false
	}
	return &user, true
}

// getAuthAdmin retrieves the authenticated user or return an unauthorized error if user is not present
func (h *ComplaintHandler) getAuthAdminOrModerator(c *gin.Context) (*model.AuthenticatedUser, bool) {
	// get admin info from auth context
	user, ok := GetCurrentAdmin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Message: "User is not authorized"})
		return nil, false
	}
	if user.Role != model.Moderator && user.Role != model.Admin {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Message: "User is not an admin"})
		return nil, false
	}
	return &user, true
}
