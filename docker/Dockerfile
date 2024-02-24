FROM golang:1.22

# User variabels, you can set them in ./env
ARG APP_CONTAINER_USERNAME
ARG APP_CONTAINER_UID
ARG APP_CONTAINER_GID

WORKDIR /usr/src/app

# install programms
RUN apt-get update && apt-get install -y vim

# to keep the container running until we run the server
CMD ["tail", "-f", "/dev/null"]

# create new User inside the Conrainer
RUN groupadd -g ${APP_CONTAINER_GID} ${APP_CONTAINER_USERNAME} \
    && useradd -m -r -u ${APP_:qqCONTAINER_UID} -g ${APP_CONTAINER_USERNAME} ${APP_CONTAINER_USERNAME}