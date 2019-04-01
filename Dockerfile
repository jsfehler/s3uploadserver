FROM iron/go

WORKDIR /app

ADD s3uploadserver /app/

ENTRYPOINT ["./s3uploadserver"]
