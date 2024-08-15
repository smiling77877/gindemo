package main

import (
	"gindemo/webook/internal/events"
	"gindemo/webook/pkg/ginx"
	"gindemo/webook/pkg/grpcx"
)

type App struct {
	consumers   []events.Consumer
	server      *grpcx.Server
	adminServer *ginx.Server
}
