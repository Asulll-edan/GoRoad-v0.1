package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"go-road-backend/internal/config"
	"go-road-backend/internal/event"
	"go-road-backend/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	fx.New(
		fx.Provide(
			config.LoadConfig,
			func(cfg *config.Config) (*nats.Conn, error) {
				return nats.Connect(cfg.NATSURL)
			},
			func(nc *nats.Conn) (*event.Publisher, error) {
				return event.NewPublisher(nc)
			},
			ws.NewHub,
			func(hub *ws.Hub, natsPub *event.Publisher, cfg *config.Config) *ws.MessageHandler {
				return ws.NewMessageHandler(hub, natsPub, nil, cfg.JWTSecret)
			},
			zap.NewProduction,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, hub *ws.Hub, handler *ws.MessageHandler, nc *nats.Conn, cfg *config.Config, logger *zap.Logger) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						mux := http.NewServeMux()
						mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
							tokenStr := r.URL.Query().Get("token")
							if tokenStr == "" {
								http.Error(w, "missing token", http.StatusUnauthorized)
								return
							}

							token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
								return []byte(cfg.JWTSecret), nil
							})
							if err != nil || !token.Valid {
								http.Error(w, "invalid token", http.StatusUnauthorized)
								return
							}

							claims, ok := token.Claims.(jwt.MapClaims)
							if !ok {
								http.Error(w, "invalid claims", http.StatusUnauthorized)
								return
							}

							userID, _ := claims["user_id"].(string)
							if userID == "" {
								http.Error(w, "invalid user", http.StatusUnauthorized)
								return
							}

							conn, err := upgrader.Upgrade(w, r, nil)
							if err != nil {
								log.Printf("ws upgrade error: %v", err)
								return
							}

							client := ws.NewWSClient(hub, conn, userID, handler)
							client.Start()
						})

						server := &http.Server{
							Addr:    ":" + cfg.WSPort,
							Handler: mux,
						}

						go server.ListenAndServe()
						go ws.StartPresenceChecker(hub, 30*time.Second)

						logger.Info("WebSocket server started", zap.String("port", cfg.WSPort))
						return nil
					},
					OnStop: func(ctx context.Context) error {
						logger.Info("WebSocket server stopped")
						return nil
					},
				})
			},
		),
	).Run()
}
