docker run --rm -v "$PWD":/go/src/github.com/jsfehler/s3uploadserver -w /go/src/github.com/jsfehler/s3uploadserver iron/go:dev go build -o s3uploadserver
docker build -t jsfehler/s3uploadserver:latest .
