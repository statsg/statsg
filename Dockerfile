FROM lacion/docker-alpine:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/ChristianWitts/statg-server"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/statg-server/bin

WORKDIR /opt/statg-server/bin

COPY bin/statg-server /opt/statg-server/bin/
RUN chmod +x /opt/statg-server/bin/statg-server

CMD /opt/statg-server/bin/statg-server