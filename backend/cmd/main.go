package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"receptor/backend/internal/config"
	"receptor/backend/internal/handler"
	"receptor/backend/internal/middleware"
	"receptor/backend/internal/repository"
	"receptor/backend/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	usersRepo := repository.NewUsersRepository(pool)
	recipesRepo := repository.NewRecipesRepository(pool)
	ingredientsRepo := repository.NewIngredientsRepository(pool)
	favoritesRepo := repository.NewFavoritesRepository(pool)

	authSvc := service.NewAuthService(usersRepo, cfg.JWTSecret)
	recipeSvc := service.NewRecipeService(recipesRepo)
	favoriteSvc := service.NewFavoriteService(favoritesRepo)
	matchSvc := service.NewMatchService(pool, ingredientsRepo)

	authH := handler.NewAuthHandler(authSvc)
	recipeH := handler.NewRecipeHandler(recipeSvc)
	favoriteH := handler.NewFavoriteHandler(favoriteSvc)
	matchH := handler.NewMatchHandler(matchSvc)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.POST("/auth/register", authH.Register)
	r.POST("/auth/login", authH.Login)
	r.GET("/auth/me", middleware.JWT(cfg.JWTSecret), authH.Me)

	r.GET("/recipes", recipeH.List)
	r.GET("/recipes/:id", recipeH.GetByID)

	fav := r.Group("/favorites", middleware.JWT(cfg.JWTSecret))
	{
		fav.POST("/:id", favoriteH.Add)
		fav.DELETE("/:id", favoriteH.Remove)
		fav.GET("", favoriteH.List)
	}

	r.POST("/match/by-ingredients", matchH.Match)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
