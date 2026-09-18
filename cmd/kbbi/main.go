package main

import (
	"context"

	"github.com/raf555/kbbi-api/internal/cmdfx"
	"github.com/raf555/kbbi-api/internal/dictionary/dictionaryfx"
	"github.com/raf555/kbbi-api/internal/home/homefx"
	httpfx "github.com/raf555/kbbi-api/internal/http/fx"
	"github.com/raf555/kbbi-api/internal/swagger/swaggerfx"
	"github.com/raf555/kbbi-api/internal/ui/uifx"
)

func main() {
	err := cmdfx.Run(context.TODO(),
		dictionaryfx.Module,
		homefx.Module,
		swaggerfx.Module,
		uifx.Module,
		httpfx.ServerInvoker,
	)

	if err != nil {
		panic("main: " + err.Error())
	}
}
