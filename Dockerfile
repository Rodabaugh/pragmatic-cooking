FROM debian:stable-slim
WORKDIR /app
COPY pragmatic-cooking /app/
COPY static /app/static/
RUN chmod +x /app/pragmatic-cooking

RUN apt-get update && apt-get install -y ca-certificates
RUN update-ca-certificates

CMD ["./pragmatic-cooking"]
