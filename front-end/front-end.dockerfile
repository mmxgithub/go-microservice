FROM alpine:latest

RUN mkdir /app

# COPY --from=builder /app/brokerApp /app
COPY frontEnd /app

CMD [ "/app/frontEnd" ]