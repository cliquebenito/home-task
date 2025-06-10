#!/bin/bash

## Данный скрипт увеличивает счетчик банера - BANNER_ID на количество AMOUNT

BANNER_ID=6
AMOUNT=100
URL="http://localhost:8080/counter/$BANNER_ID"

for ((i = 1; i <= AMOUNT; i++)); do
  curl -s -o /dev/null -w "[%{http_code}] " "$URL"
done

echo