FROM alpine
COPY ./dist/ /dist
WORKDIR /dist
ENTRYPOINT ["kubecli"]