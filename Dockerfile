FROM golang:1.19 as builder

WORKDIR /builder
COPY go.mod .
COPY go.sum .
RUN go mod tidy
COPY . .
RUN go build -o bin/pronoea cmd/main.go

FROM node:16.18.1 as nodebuilder

WORKDIR /builder
COPY ./ui .
RUN npm install
RUN npm run build:stage


FROM golang:1.19

WORKDIR /app

COPY --from=builder /builder/bin/pronoea .
COPY --from=nodebuilder /builder/dist/ ./html/
COPY ./internal/config/config.yaml ./
COPY ./email-templates ./templates

EXPOSE 8080

CMD [ "/app/pronoea" ] 
