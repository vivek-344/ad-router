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

	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/ping", func(c *gin.Context) {
		start := time.Now()
		c.JSON(http.StatusOK, gin.H{
			"ping": time.Since(start).String(),
		})
	})

	router.GET("/v1/delivery", server.delivery)
	router.GET("/v1/get_campaign/:cid", server.getCampaign)
	router.POST("/v1/create_campaign", server.createCampaign)
	router.POST("/v1/add_campaign", server.addCampaign)
	router.POST("/v1/add_target_app", server.addTargetApp)
	router.POST("/v1/add_target_country", server.addTargetCountry)
	router.POST("/v1/add_target_os", server.addTargetOs)
	router.PATCH("/v1/toggle_status/:cid", server.toggleStatus)
	router.PATCH("/v1/update_campaign_name", server.updateCampaignName)
	router.PATCH("/v1/update_campaign_image", server.updateCampaignImage)
	router.PATCH("/v1/update_campaign_cta", server.updateCampaignCta)
	router.PATCH("/v1/update_target_app", server.updateTargetApp)
	router.PATCH("/v1/update_target_country", server.updateTargetCountry)
	router.PATCH("/v1/update_target_os", server.updateTargetOs)
	router.DELETE("/v1/delete_campaign/:cid", server.deleteCampaign)
	router.DELETE("/v1/delete_target_app/:cid", server.deleteTargetApp)
	router.DELETE("/v1/delete_target_country/:cid", server.deleteTargetCountry)
	router.DELETE("/v1/delete_target_os/:cid", server.deleteTargetOs)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
