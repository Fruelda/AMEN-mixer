package events

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func Emit(ctx context.Context, event string, data any) {

	runtime.EventsEmit(
		ctx,
		event,
		data,
	)

}
