package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/vivek-344/AdRouter/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func (server *Server) Router() *gin.Engine {
	return server.router
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		start := time.Now()
		c.JSON(http.StatusOK, gin.H{
			"ping": time.Since(start).String(),
		})
	})

	router.GET("/get_campaign/:cid", server.getCampaign)
	router.POST("/create_campaign", server.createCampaign)
	router.POST("/add_campaign", server.addCampaign)
	router.POST("/add_target_app", server.addTargetApp)
	router.POST("/add_target_country", server.addTargetCountry)
	router.POST("/add_target_os", server.addTargetOs)
	router.PATCH("/toggle_status/:cid", server.toggleStatus)
	router.PATCH("/update_campaign_name", server.updateCampaignName)
	router.PATCH("/update_campaign_image", server.updateCampaignImage)
	router.PATCH("/update_campaign_cta", server.updateCampaignCta)
	router.PATCH("/update_target_app", server.updateTargetApp)
	router.PATCH("/update_target_country", server.updateTargetCountry)
	router.PATCH("/update_target_os", server.updateTargetOs)
	router.DELETE("/delete_campaign/:cid", server.deleteCampaign)
	router.DELETE("/delete_target_app/:cid", server.deleteTargetApp)
	router.DELETE("/delete_target_country/:cid", server.deleteTargetCountry)
	router.DELETE("/delete_target_os/:cid", server.deleteTargetOs)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
