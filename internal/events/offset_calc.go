package events

import (
	"context"
	"strconv"

	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/internal/log"
	"github.com/kaleido-io/firefly/pkg/database"
	"github.com/kaleido-io/firefly/pkg/fftypes"
)

func calcFirstOffset(ctx context.Context, di database.Plugin, pfe *fftypes.SubOptsFirstEvent) (firstOffset int64, err error) {
	firstEvent := fftypes.SubOptsFirstEventNewest
	if pfe != nil {
		firstEvent = *pfe
	}
	firstOffset = -1
	var useNewest bool
	switch firstEvent {
	case "", fftypes.SubOptsFirstEventNewest:
		useNewest = true
	case fftypes.SubOptsFirstEventOldest:
		useNewest = false
	default:
		specificSequence, err := strconv.ParseInt(string(firstEvent), 10, 64)
		if err != nil {
			return -1, i18n.WrapError(ctx, err, i18n.MsgInvalidFirstEvent, firstEvent)
		}
		if specificSequence < -1 {
			return -1, i18n.NewError(ctx, i18n.MsgNumberMustBeGreaterEqual, -1)
		}
		firstOffset = specificSequence
		useNewest = false
	}
	if useNewest {
		f := database.EventQueryFactory.NewFilter(ctx).And().Sort("sequence").Descending().Limit(1)
		newestEvents, err := di.GetEvents(ctx, f)
		if err != nil {
			return firstOffset, err
		}
		if len(newestEvents) > 0 {
			return newestEvents[0].Sequence, nil
		}
	}
	log.L(ctx).Debugf("Event poller initial offest: %d (newest=%t)", firstOffset, useNewest)
	return firstOffset, err
}
