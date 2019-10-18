FROM alpine
ADD future.web.basic /future.web.basic
ADD appconfig.json /appconfig.json
ENTRYPOINT [ "/future.web.basic" ]
