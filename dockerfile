#Build distribution from golang:1.11
FROM golang:1.11 as build
ADD . /go/src/github.com/thomaspoignant/user-microservice
RUN go install /go/src/github.com/thomaspoignant/user-microservice

# Copy to distroless image to have a more secure container
FROM gcr.io/distroless/base
COPY --from=build /go/bin/user-microservice / 
CMD ["/user-microservice"]
EXPOSE 8080