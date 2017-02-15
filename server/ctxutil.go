package server

import "golang.org/x/net/context"

//cannot use string as context key so we need to introduce own type
type ctxKey string

const userKey = ctxKey("user")

func setUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

//GetUser obtains value from the Context
func GetUser(ctx context.Context) *User {
	return ctx.Value(userKey).(*User)
}
