FROM golang:1.25.5-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/user_bot ./cmd/user_bot

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /out/user_bot /user_bot
USER nonroot:nonroot
ENTRYPOINT ["/user_bot"]
