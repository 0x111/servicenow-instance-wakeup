FROM golang as builder
ADD . /src
RUN cd /src && make

#-------
FROM chromedp/headless-shell
ENV HEADLESS "true"
ENV DEBUG "false"
ENV USERNAME ""
ENV PASSWORD ""
WORKDIR /app
COPY --from=builder /src/build/servicenow-instance-wakeup-linux-amd64 /app/servicenow-instance-wakeup
RUN chmod a+x /app/servicenow-instance-wakeup
ENTRYPOINT /app/servicenow-instance-wakeup -headless=$HEADLESS -username=$USERNAME -password=$PASSWORD -debug=$DEBUG