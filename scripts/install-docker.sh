#!/bin/bash
DIR="$(dirname "$(realpath "$0")")"
BASE_NAME="basement"

# If container with BASE_NAME exists, try appending 1, 2, or 3
if docker container inspect "$BASE_NAME" > /dev/null 2>&1; then
  for i in 1 2 3; do
    NEW_NAME="${BASE_NAME}${i}"
    if ! docker container inspect "$NEW_NAME" > /dev/null 2>&1; then
      BASE_NAME="$NEW_NAME"
      break
    fi
  done
fi

echo "Container name set to: $BASE_NAME"
# Run container, print a message before executing git commands in the script
docker run -it --name "$BASE_NAME" -v "$DIR/install-debian.sh:/app/install-debian.sh" -p 8101:8101 --entrypoint /bin/bash -w /app -e SUDO=" " debian:bookworm
