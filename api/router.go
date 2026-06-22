package router

import (
	"github.com/gorilla/mux"
	"github.com/temuka-api-service/internal/handler"
	"github.com/temuka-api-service/internal/publisher"
	"github.com/temuka-api-service/internal/repository"
	"github.com/temuka-api-service/internal/service"
	"github.com/temuka-api-service/middleware"
	database "github.com/temuka-api-service/util/database"
	fileStorage "github.com/temuka-api-service/util/file_storage"
	keyValueStore "github.com/temuka-api-service/util/key_value_store"
	"github.com/temuka-api-service/util/queue"
)

func Routes(db database.PostgresWrapper, redis keyValueStore.RedisWrapper, storage fileStorage.S3Wrapper, rmq queue.RabbitMQChannel) *mux.Router {
	router := mux.NewRouter()

	// Init repositories
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	communityRepo := repository.NewCommunityRepository(db)
	moderatorRepo := repository.NewModeratorRepository(db)
	reportRepo := repository.NewReportRepository(db)
	universityRepo := repository.NewUniversityRepository(db)
	majorRepo := repository.NewMajorRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	conversationRepo := repository.NewConversationRepository(db)

	// Init publishers
	searchIndexPublisher := publisher.NewSearchIndexPublisher(rmq)

	// Init services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)
	postService := service.NewPostService(postRepo, userRepo, commentRepo, notificationRepo, communityRepo, redis, searchIndexPublisher)
	notificationService := service.NewNotificationService(notificationRepo)
	commentService := service.NewCommentService(commentRepo, postRepo, notificationRepo, reportRepo)
	communityService := service.NewCommunityService(communityRepo)
	moderatorService := service.NewModeratorService(moderatorRepo, notificationRepo)
	reportService := service.NewReportService(reportRepo)
	universityService := service.NewUniversityService(universityRepo, reviewRepo)
	majorService := service.NewMajorService(majorRepo)
	locationService := service.NewLocationService(locationRepo)
	conversationService := service.NewConversationService(conversationRepo, userRepo)
	fileService := service.NewFileService(storage)

	// Init controllers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	communityHandler := handler.NewCommunityHandler(communityService)
	commentHandler := handler.NewCommentHandler(commentService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	moderatorHandler := handler.NewModeratorHandler(moderatorService)
	reportHandler := handler.NewReportHandler(reportService)
	universityHandler := handler.NewUniversityHandler(universityService)
	majorHandler := handler.NewMajorHandler(majorService)
	locationHandler := handler.NewLocationHandler(locationService)
	conversationHandler := handler.NewConversationHandler(conversationService)
	fileUploadHandler := handler.NewFileHandler(fileService)

	// Init routers
	authRouter := router.PathPrefix("/api/auth").Subrouter()
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRouter.HandleFunc("/resetPassword/{id}", authHandler.ResetPassword).Methods("POST")

	userRouter := router.PathPrefix("/api/user").Subrouter()
	userRouter.Use(middleware.CheckAuth)
	userRouter.HandleFunc("", userHandler.CreateUser).Methods("POST")
	userRouter.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
	userRouter.HandleFunc("/search", userHandler.SearchUsers).Methods("GET")
	userRouter.HandleFunc("/follow", userHandler.FollowUser).Methods("POST")
	userRouter.HandleFunc("/followers", userHandler.GetFollowers).Methods("GET")
	userRouter.HandleFunc("/{id}", userHandler.GetUserDetail).Methods("GET")

	postRouter := router.PathPrefix("/api/post").Subrouter()
	postRouter.Use(middleware.CheckAuth)
	postRouter.HandleFunc("", postHandler.CreatePost).Methods("POST")
	postRouter.HandleFunc("/{id}", postHandler.GetPostDetail).Methods("GET")
	postRouter.HandleFunc("/timeline/{user_id}", postHandler.GetTimelinePosts).Methods("GET")
	postRouter.HandleFunc("/user/{user_id}", postHandler.GetUserPosts).Methods("GET")
	postRouter.HandleFunc("/like/{id}", postHandler.LikePost).Methods("PUT")
	postRouter.HandleFunc("/{id}", postHandler.DeletePost).Methods("DELETE")
	postRouter.HandleFunc("/{id}", postHandler.UpdatePost).Methods("PUT")

	commentRouter := router.PathPrefix("/api/comment").Subrouter()
	commentRouter.Use(middleware.CheckAuth)
	commentRouter.HandleFunc("", commentHandler.AddComment).Methods("POST")
	commentRouter.HandleFunc("/replies", commentHandler.ShowReplies).Methods("GET")
	commentRouter.HandleFunc("/{commentId}", commentHandler.DeleteComment).Methods("DELETE")
	commentRouter.HandleFunc("/show", commentHandler.ShowCommentsByPost).Methods("GET")

	communityRouter := router.PathPrefix("/api/community").Subrouter()
	communityRouter.Use(middleware.CheckAuth)
	communityRouter.HandleFunc("", communityHandler.CreateCommunity).Methods("POST")
	communityRouter.HandleFunc("", communityHandler.GetCommunities).Methods("GET")
	communityRouter.HandleFunc("/join/{community_id}", communityHandler.JoinCommunity).Methods("POST")
	communityRouter.HandleFunc("/post/{id}", communityHandler.GetCommunityPosts).Methods("GET")
	communityRouter.HandleFunc("/user", communityHandler.GetUserJoinedCommunities).Methods("POST")
	communityRouter.HandleFunc("/{slug}", communityHandler.GetCommunityDetail).Methods("GET")
	communityRouter.HandleFunc("/{id}", communityHandler.DeleteCommunity).Methods("DELETE")
	communityRouter.HandleFunc("/{id}", communityHandler.UpdateCommunity).Methods("PUT")

	fileRouter := router.PathPrefix("/api/file").Subrouter()
	fileRouter.Use(middleware.CheckAuth)
	fileRouter.HandleFunc("", fileUploadHandler.Upload).Methods("POST")

	notificationRouter := router.PathPrefix("/api/notification").Subrouter()
	notificationRouter.HandleFunc("/list/{user_id}", notificationHandler.GetNotificationsByUser).Methods("GET")

	moderatorRouter := router.PathPrefix("/api/moderator").Subrouter()
	moderatorRouter.Use(middleware.CheckAuth)
	moderatorRouter.HandleFunc("/send", moderatorHandler.SendModeratorRequest).Methods("POST")
	moderatorRouter.HandleFunc("/{id}", moderatorHandler.RemoveModerator).Methods("DELETE")

	reportRouter := router.PathPrefix("/api/report").Subrouter()
	reportRouter.Use(middleware.CheckAuth)
	reportRouter.HandleFunc("", reportHandler.CreateReport).Methods("POST")
	reportRouter.HandleFunc("/{id}", reportHandler.DeleteReport).Methods("DELETE")

	universityRouter := router.PathPrefix("/api/university").Subrouter()
	universityRouter.Use(middleware.CheckAuth)
	universityRouter.HandleFunc("", universityHandler.AddUniversity).Methods("POST")
	universityRouter.HandleFunc("/{id}", universityHandler.UpdateUniversity).Methods("PUT")
	universityRouter.HandleFunc("/{slug}", universityHandler.GetUniversityDetail).Methods("GET")
	universityRouter.HandleFunc("", universityHandler.GetUniversities).Methods("GET")
	universityRouter.HandleFunc("/review", universityHandler.AddReview).Methods("POST")
	universityRouter.HandleFunc("/review/university_id", universityHandler.GetUniversityReviews).Methods("GET")

	majorRouter := router.PathPrefix("/api/major").Subrouter()
	majorRouter.Use(middleware.CheckAuth)
	majorRouter.HandleFunc("", majorHandler.AddMajor).Methods("POST")
	majorRouter.HandleFunc("/{id}", majorHandler.GetMajorDetail).Methods("GET")
	majorRouter.HandleFunc("", majorHandler.GetMajors).Methods("GET")
	majorRouter.HandleFunc("/review", majorHandler.AddMajorReview).Methods("POST")
	majorRouter.HandleFunc("/review/major_id", majorHandler.GetMajorReviews).Methods("GET")
	majorRouter.HandleFunc("/university/{university_id}", majorHandler.GetMajorsByUniversity).Methods("GET")

	locationRouter := router.PathPrefix("/api/location").Subrouter()
	locationRouter.Use(middleware.CheckAuth)
	locationRouter.HandleFunc("", locationHandler.AddLocation).Methods("POST")
	locationRouter.HandleFunc("", locationHandler.GetLocations).Methods("GET")
	locationRouter.HandleFunc("/{id}", locationHandler.UpdateLocation).Methods("PUT")

	conversationRouter := router.PathPrefix("/api/conversation").Subrouter()
	conversationRouter.Use(middleware.CheckAuth)
	conversationRouter.HandleFunc("", conversationHandler.AddConversation).Methods("POST")
	conversationRouter.HandleFunc("/{id}", conversationHandler.DeleteConversation).Methods("DELETE")
	conversationRouter.HandleFunc("/{id}", conversationHandler.GetConversationDetail).Methods("GET")
	conversationRouter.HandleFunc("/participant", conversationHandler.AddParticipant).Methods("POST")
	conversationRouter.HandleFunc("/message", conversationHandler.AddMessage).Methods("POST")
	conversationRouter.HandleFunc("/message/{conversation_id}", conversationHandler.RetrieveMessages).Methods("GET")
	conversationRouter.HandleFunc("/all/{user_id}", conversationHandler.GetConversationsByUserID).Methods("GET")

	return router
}
