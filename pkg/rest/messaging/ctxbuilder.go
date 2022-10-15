package messaging

import "context"

type ArgumentKey string
type ContextualArguments struct {
	ctx context.Context
}

func (ca *ContextualArguments) Of(key ArgumentKey) interface{} {
	return ca.ctx.Value(key)
}

func CreateContextualArguments(arguments map[ArgumentKey]interface{}) *ContextualArguments {
	ctx := context.Background()
	for k, v := range arguments {
		ctx = context.WithValue(ctx, k, v)
	}
	return &ContextualArguments{ctx: ctx}
}
