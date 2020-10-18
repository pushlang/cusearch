// -build OMIT

package contxt

import (
	"context"
)

type key int

const cxKey key = 0

func NewContext(ctx context.Context, cx string) context.Context {
	return context.WithValue(ctx, cxKey, cx)
}

func FromContext(ctx context.Context) (string, bool) {
	cx, ok := ctx.Value(cxKey).(string)
	return cx, ok
}
