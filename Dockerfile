FROM alpine
ADD future.web.basic /future.web.basic
ADD appconfig.yml /appconfig.yml
ENTRYPOINT [ "/future.web.basic" ]
