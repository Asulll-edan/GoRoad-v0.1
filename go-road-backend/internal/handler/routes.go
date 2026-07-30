package handler

import (
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"

	"go-road-backend/internal/config"
	"go-road-backend/internal/middleware"
	"go-road-backend/internal/repository/redis"
	"go-road-backend/internal/repository/postgres"
	domain "go-road-backend/internal/domain/auth"
)

func SetupRoutes(
	app *fiber.App,
	cfg *config.Config,
	cache redis.CacheRepository,
	db *postgres.Database,
	logger *zap.Logger,
	authService domain.Service,
) {
	authMw := middleware.AuthMiddleware(cfg.JWTSecret)
	adminMw := middleware.AdminMiddleware()
	rateLimitMw := middleware.RateLimitMiddleware(cache)

	api := app.Group("/api/v1", rateLimitMw)

	// Store services in context
	api.Use(func(c fiber.Ctx) error {
		c.Locals("config", cfg)
		c.Locals("auth_service", authService)
		c.Locals("logger", logger)
		c.Locals("cache", cache)
		c.Locals("db", db)
		return c.Next()
	})

	// ============================================
	// AUTH — Public
	// ============================================
	auth := api.Group("/auth")
	auth.Post("/register", handleRegister)
	auth.Post("/login", handleLogin)
	auth.Post("/refresh", handleRefreshToken)
	auth.Post("/forgot-password", handleForgotPassword)
	auth.Post("/reset-password", handleResetPassword)

	// ============================================
	// AUTH — Protected
	// ============================================
	authProtected := api.Group("/auth", authMw)
	authProtected.Get("/me", handleGetProfile)
	authProtected.Put("/me", handleUpdateProfile)
	authProtected.Put("/me/avatar", handleUploadAvatar)
	authProtected.Put("/me/password", handleChangePassword)
	authProtected.Post("/devices", handleRegisterDevice)
	authProtected.Delete("/devices/:id", handleUnregisterDevice)

	// ============================================
	// RIDERS
	// ============================================
	riders := api.Group("/riders", authMw)
	riders.Get("/:id/profile", handleGetRiderProfile)
	riders.Get("/:id/badges", handleGetRiderBadges)
	riders.Get("/:id/stats", handleGetRiderStats)
	riders.Get("/:id/motors", handleGetRiderMotors)
	riders.Get("/me/followers", handleGetFollowers)
	riders.Get("/me/following", handleGetFollowing)
	riders.Post("/:id/follow", handleFollowUser)
	riders.Delete("/:id/follow", handleUnfollowUser)
	riders.Post("/:id/block", handleBlockUser)
	riders.Delete("/:id/block", handleUnblockUser)

	// ============================================
	// ROOMS
	// ============================================
	rooms := api.Group("/rooms", authMw)
	rooms.Get("/", handleListRooms)
	rooms.Get("/discover", handleDiscoverRooms)
	rooms.Post("/", handleCreateRoom)
	rooms.Get("/:id", handleGetRoom)
	rooms.Put("/:id", handleUpdateRoom)
	rooms.Delete("/:id", handleDeleteRoom)
	rooms.Post("/:id/join", handleJoinRoom)
	rooms.Post("/:id/leave", handleLeaveRoom)
	rooms.Get("/:id/members", handleListRoomMembers)
	rooms.Put("/:id/members/:userId/role", handleUpdateMemberRole)
	rooms.Get("/:id/settings", handleGetRoomSettings)
	rooms.Put("/:id/settings", handleUpdateRoomSettings)
	rooms.Post("/:id/start", handleStartTouring)
	rooms.Post("/:id/pause", handlePauseTouring)
	rooms.Post("/:id/complete", handleCompleteTouring)
	rooms.Post("/:id/cancel", handleCancelTouring)

	// ============================================
	// CONVOY
	// ============================================
	convoy := api.Group("/rooms/:id/convoy", authMw)
	convoy.Post("/formation", handleCreateFormation)
	convoy.Put("/formation/:formationId", handleUpdateFormation)
	convoy.Get("/formation", handleGetActiveFormation)
	convoy.Post("/location", handleUpdateLocation)
	convoy.Get("/locations", handleGetLocations)
	convoy.Get("/tracking", handleGetTracking)

	// ============================================
	// ROUTES
	// ============================================
	routes := api.Group("/routes", authMw)
	routes.Post("/", handleCreateRoute)
	routes.Get("/:id", handleGetRoute)
	routes.Put("/:id", handleUpdateRoute)
	routes.Delete("/:id", handleDeleteRoute)
	routes.Post("/:id/waypoints", handleAddWaypoint)
	routes.Get("/:id/waypoints", handleListWaypoints)
	routes.Post("/import/gpx", handleImportGPX)
	routes.Get("/:id/export/gpx", handleExportGPX)
	routes.Post("/:id/activate", handleActivateRoute)

	// ============================================
	// CHAT
	// ============================================
	chat := api.Group("/rooms/:id/messages", authMw)
	chat.Get("/", handleListMessages)
	chat.Post("/", handleSendMessage)
	chat.Get("/:messageId", handleGetMessage)
	chat.Put("/:messageId", handleEditMessage)
	chat.Delete("/:messageId", handleDeleteMessage)
	chat.Post("/:messageId/pin", handlePinMessage)
	chat.Post("/:messageId/read", handleMarkRead)

	// ============================================
	// EMERGENCY
	// ============================================
	emergency := api.Group("/emergency", authMw)
	emergency.Post("/", handleReportEmergency)
	emergency.Get("/", handleListEmergencies)
	emergency.Get("/:id", handleGetEmergency)
	emergency.Post("/:id/acknowledge", handleAcknowledgeEmergency)
	emergency.Post("/:id/resolve", handleResolveEmergency)
	emergency.Post("/sos", handleTriggerSOS)
	emergency.Post("/sos/:id/dismiss", handleDismissSOS)

	// ============================================
	// VOTING
	// ============================================
	voting := api.Group("/rooms/:id/votings", authMw)
	voting.Post("/", handleCreateVoting)
	voting.Get("/", handleListVotings)
	voting.Get("/:votingId", handleGetVoting)
	voting.Post("/:votingId/vote", handleSubmitVote)
	voting.Post("/:votingId/close", handleCloseVoting)
	voting.Get("/:votingId/results", handleGetVotingResults)

	// ============================================
	// WEATHER
	// ============================================
	weather := api.Group("/weather", authMw)
	weather.Get("/current", handleGetCurrentWeather)
	weather.Get("/forecast", handleGetWeatherForecast)
	weather.Get("/alerts", handleGetWeatherAlerts)
	weather.Get("/route", handleGetRouteWeather)

	// ============================================
	// AI
	// ============================================
	ai := api.Group("/ai", authMw)
	ai.Post("/chat", handleAIChat)
	ai.Post("/itinerary", handleGenerateItinerary)
	ai.Post("/cost", handleEstimateCost)
	ai.Post("/route-advice", handleRouteAdvice)
	ai.Post("/packing-list", handlePackingList)
	ai.Post("/safety", handleSafetyAdvice)

	// ============================================
	// MOTORS
	// ============================================
	motors := api.Group("/motors", authMw)
	motors.Post("/", handleCreateMotor)
	motors.Get("/", handleListMotors)
	motors.Get("/:id", handleGetMotor)
	motors.Put("/:id", handleUpdateMotor)
	motors.Delete("/:id", handleDeleteMotor)
	motors.Post("/:id/primary", handleSetPrimaryMotor)

	// ============================================
	// FUEL LOGS
	// ============================================
	fuel := api.Group("/fuel-logs", authMw)
	fuel.Post("/", handleCreateFuelLog)
	fuel.Get("/", handleListFuelLogs)
	fuel.Get("/:id", handleGetFuelLog)
	fuel.Put("/:id", handleUpdateFuelLog)
	fuel.Delete("/:id", handleDeleteFuelLog)
	fuel.Get("/analytics", handleFuelAnalytics)

	// ============================================
	// EXPENSES
	// ============================================
	expenses := api.Group("/expenses", authMw)
	expenses.Post("/", handleCreateExpense)
	expenses.Get("/", handleListExpenses)
	expenses.Get("/:id", handleGetExpense)
	expenses.Put("/:id", handleUpdateExpense)
	expenses.Delete("/:id", handleDeleteExpense)
	expenses.Get("/summary", handleExpenseSummary)

	// ============================================
	// NOTIFICATIONS
	// ============================================
	notif := api.Group("/notifications", authMw)
	notif.Get("/", handleListNotifications)
	notif.Post("/:id/read", handleMarkNotificationRead)
	notif.Post("/read-all", handleMarkAllRead)
	notif.Get("/unread-count", handleUnreadCount)
	notif.Put("/preferences", handleUpdateNotificationPreferences)

	// ============================================
	// POI
	// ============================================
	poi := api.Group("/poi", authMw)
	poi.Get("/nearby", handleGetNearbyPOI)
	poi.Get("/:id", handleGetPOIDetail)
	poi.Post("/report", handleReportPOI)
	poi.Get("/categories", handleGetPOICategories)

	// ============================================
	// BADGES
	// ============================================
	badges := api.Group("/badges", authMw)
	badges.Get("/", handleListBadges)
	badges.Get("/me", handleMyBadges)
	badges.Get("/progress", handleBadgeProgress)

	// ============================================
	// SOCIAL / FEED
	// ============================================
	social := api.Group("/social", authMw)
	social.Post("/posts", handleCreatePost)
	social.Get("/posts/:id", handleGetPost)
	social.Delete("/posts/:id", handleDeletePost)
	social.Post("/posts/:id/like", handleLikePost)
	social.Delete("/posts/:id/like", handleUnlikePost)
	social.Get("/posts/:id/comments", handleListComments)
	social.Post("/posts/:id/comments", handleCreateComment)
	social.Delete("/posts/:id/comments/:commentId", handleDeleteComment)
	social.Post("/report", handleReportContent)

	// ============================================
	// FEED
	// ============================================
	feed := api.Group("/feed", authMw)
	feed.Get("/", handleGetFeed)
	feed.Get("/explore", handleExploreFeed)

	// ============================================
	// CHECKLISTS
	// ============================================
	checklist := api.Group("/checklists", authMw)
	checklist.Get("/templates", handleListTemplates)
	checklist.Post("/templates", handleCreateTemplate)
	checklist.Get("/templates/:id", handleGetTemplate)
	checklist.Post("/rooms/:roomId/items", handleCreateTouringChecklist)
	checklist.Get("/rooms/:roomId/items", handleGetTouringChecklist)
	checklist.Put("/rooms/:roomId/items/:itemId", handleToggleChecklistItem)

	// ============================================
	// QR
	// ============================================
	qr := api.Group("/qr", authMw)
	qr.Get("/me", handleGetMyQRCard)
	qr.Post("/regenerate", handleRegenerateQR)
	qr.Get("/scan/:code", handleScanQR)

	// ============================================
	// SERVICE REMINDERS
	// ============================================
	reminders := api.Group("/service-reminders", authMw)
	reminders.Post("/", handleCreateReminder)
	reminders.Get("/", handleListReminders)
	reminders.Put("/:id", handleUpdateReminder)
	reminders.Delete("/:id", handleDeleteReminder)
	reminders.Post("/:id/complete", handleCompleteReminder)

	// ============================================
	// UPLOAD
	// ============================================
	upload := api.Group("/upload", authMw)
	upload.Post("/", handleUploadFile)
	upload.Post("/avatar", handleUploadAvatar)
	upload.Post("/photo", handleUploadPhoto)
	upload.Delete("/:id", handleDeleteFile)

	// ============================================
	// VOICE (LiveKit Token)
	// ============================================
	voice := api.Group("/voice", authMw)
	voice.Post("/token", handleVoiceToken)

	// ============================================
	// ADMIN (Protected + Admin-only)
	// ============================================
	admin := api.Group("/admin", authMw, adminMw)
	admin.Get("/dashboard", handleAdminDashboard)
	admin.Get("/users", handleAdminListUsers)
	admin.Get("/users/:id", handleAdminGetUser)
	admin.Post("/users/:id/ban", handleAdminBanUser)
	admin.Post("/users/:id/unban", handleAdminUnbanUser)
	admin.Get("/rooms", handleAdminListRooms)
	admin.Get("/rooms/:id", handleAdminGetRoom)
	admin.Get("/reports", handleAdminListReports)
	admin.Post("/reports/:id/review", handleAdminReviewReport)
	admin.Get("/emergency", handleAdminListEmergency)
	admin.Get("/analytics", handleAdminAnalytics)
	admin.Get("/logs", handleAdminLogs)

	logger.Info("Routes registered successfully")
}
