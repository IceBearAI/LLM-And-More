package auth

type Middleware func(Service) Service
