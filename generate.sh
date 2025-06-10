#!/bin/bash

## Данный скрипт создает TOTAL баннеров

URL="http://localhost:8080/banners"
TOTAL=150

echo "Creating $TOTAL banners"

for ((i = 1; i <= TOTAL; i++)); do
  NAME="Banner_$(date +%s%N | cut -b10-19)_$RANDOM"
  PAYLOAD="{\"name\": \"$NAME\"}"

  response=$(curl -s -w "\n%{http_code}" -X POST "$URL" \
    -H "Content-Type: application/json" \
    -d "$PAYLOAD")

  body=$(echo "$response" | head -n1)
  code=$(echo "$response" | tail -n1)

  if [ "$code" -eq 201 ] || [ "$code" -eq 200 ]; then
    echo "✅ {$code} Banner created: $NAME"
  else
    echo "❌ {$code} Failed to create banner: $NAME"
    echo "Server response: $body"
  fi
done

echo "Done."