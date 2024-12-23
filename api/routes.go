package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/vivek-344/AdRouter/db/sqlc"
)

type deliveryRequest struct {
	AppID   string `binding:"required" form:"app"`
	Country string `binding:"required" form:"country"`
	Os      string `binding:"required" form:"os"`
}

func (s *Server) delivery(ctx *gin.Context) {
	var req deliveryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := s.store.Delivery(ctx.Request.Context(), db.DeliveryParams{
		AppID:   req.AppID,
		Country: req.Country,
		Os:      req.Os,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(response) == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{"error": "no campaign available"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

type createCampaignRequest struct {
	Cid         string `binding:"required" json:"cid"`
	Name        string `binding:"required,min=6,max=32" json:"name"`
	Img         string `binding:"required" json:"img"`
	Cta         string `binding:"required" json:"cta"`
	AppID       string `json:"app"`
	AppRule     string `binding:"omitempty,oneof=include exclude" json:"app_rule"`
	Country     string `json:"country"`
	CountryRule string `binding:"omitempty,oneof=include exclude" json:"country_rule"`
	Os          string `json:"os"`
	OsRule      string `binding:"omitempty,oneof=include exclude" json:"os_rule"`
}

func (s *Server) createCampaign(ctx *gin.Context) {
	var req createCampaignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.AppID != "" && req.AppRule == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "AppRule field is empty"})
		return
	}
	if req.Country != "" && req.CountryRule == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "CountryRule field is empty"})
		return
	}
	if req.Os != "" && req.OsRule == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "OsRule field is empty"})
		return
	}

	campaign, err := s.store.CreateCampaign(ctx.Request.Context(), db.CreateCampaignParams{
		Cid:         req.Cid,
		Name:        req.Name,
		Img:         req.Img,
		Cta:         req.Cta,
		AppID:       req.AppID,
		AppRule:     db.RuleType(req.AppRule),
		Country:     req.Country,
		CountryRule: db.RuleType(req.CountryRule),
		Os:          req.Os,
		OsRule:      db.RuleType(req.OsRule),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, campaign)
}

type getCampaignRequest struct {
	Cid string `binding:"required" uri:"cid"`
}

func (s *Server) getCampaign(ctx *gin.Context) {
	var req getCampaignRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign, err := s.store.ReadCampaign(ctx.Request.Context(), req.Cid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, campaign)
}

type addCampaignRequest struct {
	Cid  string `binding:"required" json:"cid"`
	Name string `binding:"required" json:"name"`
	Img  string `binding:"required" json:"img"`
	Cta  string `binding:"required" json:"cta"`
}

func (s *Server) addCampaign(ctx *gin.Context) {
	var req addCampaignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign, err := s.store.CreateCampaign(ctx.Request.Context(), db.CreateCampaignParams{
		Cid:  req.Cid,
		Name: req.Name,
		Img:  req.Img,
		Cta:  req.Cta,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, campaign)
}

type addTargetAppRequest struct {
	Cid   string `binding:"required" json:"cid"`
	AppID string `binding:"required" json:"app"`
	Rule  string `binding:"required,oneof=include exclude" json:"rule"`
}

func (s *Server) addTargetApp(ctx *gin.Context) {
	var req addTargetAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target_app, err := s.store.AddTargetApp(ctx.Request.Context(), db.AddTargetAppParams{
		Cid:   req.Cid,
		AppID: req.AppID,
		Rule:  db.RuleType(req.Rule),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, target_app)
}

type addTargetCountryRequest struct {
	Cid     string `binding:"required" json:"cid"`
	Country string `binding:"required" json:"country"`
	Rule    string `binding:"required,oneof=include exclude" json:"rule"`
}

func (s *Server) addTargetCountry(ctx *gin.Context) {
	var req addTargetCountryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target_country, err := s.store.AddTargetCountry(ctx.Request.Context(), db.AddTargetCountryParams{
		Cid:     req.Cid,
		Country: req.Country,
		Rule:    db.RuleType(req.Rule),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, target_country)
}

type addTargetOsRequest struct {
	Cid  string `binding:"required" json:"cid"`
	Os   string `binding:"required" json:"os"`
	Rule string `binding:"required,oneof=include exclude" json:"rule"`
}

func (s *Server) addTargetOs(ctx *gin.Context) {
	var req addTargetOsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target_os, err := s.store.AddTargetOs(ctx.Request.Context(), db.AddTargetOsParams{
		Cid:  req.Cid,
		Os:   req.Os,
		Rule: db.RuleType(req.Rule),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, target_os)
}

type deleteCampaignRequest struct {
	Cid string `binding:"required" uri:"cid"`
}

func (s *Server) deleteCampaign(ctx *gin.Context) {
	var req deleteCampaignRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.store.DeleteCampaign(ctx.Request.Context(), req.Cid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "campaign deleted successfully"})
}

type deleteTargetAppRequest struct {
	Cid string `binding:"required" uri:"cid"`
}

func (s *Server) deleteTargetApp(ctx *gin.Context) {
	var req deleteTargetAppRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.store.DeleteTargetApp(ctx.Request.Context(), req.Cid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

type deleteTargetCountryRequest struct {
	Cid string `binding:"required" uri:"cid"`
}

func (s *Server) deleteTargetCountry(ctx *gin.Context) {
	var req deleteTargetCountryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.store.DeleteTargetCountry(ctx.Request.Context(), req.Cid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

type deleteTargetOsRequest struct {
	Cid string `binding:"required" uri:"cid"`
}

func (s *Server) deleteTargetOs(ctx *gin.Context) {
	var req deleteTargetOsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.store.DeleteTargetOs(ctx.Request.Context(), req.Cid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "resource deleted successfully"})
}

type toggleStatusRequest struct {
	Cid string `binding:"required" uri:"cid"`
}

func (s *Server) toggleStatus(ctx *gin.Context) {
	var req toggleStatusRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.store.ToggleStatus(ctx.Request.Context(), req.Cid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	campaign, err := s.store.GetCampaign(ctx.Request.Context(), req.Cid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "campaign " + campaign.Cid + " status toggled to " + string(campaign.Status) + " successfully"})
}

type updateCampaignNameRequest struct {
	Cid  string `binding:"required" json:"cid"`
	Name string `binding:"required,min=6,max=32" json:"name"`
}

func (s *Server) updateCampaignName(ctx *gin.Context) {
	var req updateCampaignNameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign, err := s.store.UpdateCampaignName(ctx.Request.Context(), db.UpdateCampaignNameParams{
		Cid:  req.Cid,
		Name: req.Name,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "campaign " + campaign.Cid + " name changed to " + campaign.Name + " successfully"})
}

type updateCampaignImageRequest struct {
	Cid string `binding:"required" json:"cid"`
	Img string `binding:"required" json:"img"`
}

func (s *Server) updateCampaignImage(ctx *gin.Context) {
	var req updateCampaignImageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign, err := s.store.UpdateCampaignImage(ctx.Request.Context(), db.UpdateCampaignImageParams{
		Cid: req.Cid,
		Img: req.Img,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "campaign " + campaign.Cid + " image changed to " + campaign.Img + " successfully"})
}

type updateCampaignCtaRequest struct {
	Cid string `binding:"required" json:"cid"`
	Cta string `binding:"required" json:"cta"`
}

func (s *Server) updateCampaignCta(ctx *gin.Context) {
	var req updateCampaignCtaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign, err := s.store.UpdateCampaignCta(ctx.Request.Context(), db.UpdateCampaignCtaParams{
		Cid: req.Cid,
		Cta: req.Cta,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "campaign " + campaign.Cid + " cta changed to " + campaign.Cta + " successfully"})
}

type updateTargetAppRequest struct {
	Cid     string `binding:"required" json:"cid"`
	AppID   string `binding:"required" json:"app"`
	AppRule string `binding:"required,oneof=include exclude" json:"rule"`
}

func (s *Server) updateTargetApp(ctx *gin.Context) {
	var req updateTargetAppRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target_app, err := s.store.AddTargetApp(ctx.Request.Context(), db.AddTargetAppParams{
		Cid:   req.Cid,
		AppID: req.AppID,
		Rule:  db.RuleType(req.AppRule),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, target_app)
}

type updateTargetCountryRequest struct {
	Cid         string `binding:"required" json:"cid"`
	Country     string `binding:"required" json:"country_id"`
	CountryRule string `binding:"required,oneof=include exclude" json:"rule"`
}

func (s *Server) updateTargetCountry(ctx *gin.Context) {
	var req updateTargetCountryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target_country, err := s.store.AddTargetCountry(ctx.Request.Context(), db.AddTargetCountryParams{
		Cid:     req.Cid,
		Country: req.Country,
		Rule:    db.RuleType(req.CountryRule),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, target_country)
}

type updateTargetOsRequest struct {
	Cid    string `binding:"required" json:"cid"`
	Os     string `binding:"required" json:"os"`
	OsRule string `binding:"required,oneof=include exclude" json:"rule"`
}

func (s *Server) updateTargetOs(ctx *gin.Context) {
	var req updateTargetOsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	target_os, err := s.store.AddTargetOs(ctx.Request.Context(), db.AddTargetOsParams{
		Cid:  req.Cid,
		Os:   req.Os,
		Rule: db.RuleType(req.OsRule),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, target_os)
}
