FROM golang as builder

RUN go get github.com/0x111/servicenow-instance-wakeup && cd /go/src/github.com/0x111/servicenow-instance-wakeup && make

#-------
FROM chromedp/headless-shell
ENV HEADLESS "false"
ENV DEBUG "false"
ENV USERNAME ""
ENV PASSWORD ""
WORKDIR /app
COPY --from=builder /go/src/github.com/0x111/servicenow-instance-wakeup/build/servicenow-instance-wakeup-linux-amd64 /app/servicenow-instance-wakeup
RUN chmod a+x /app/servicenow-instance-wakeup
ENTRYPOINT /app/servicenow-instance-wakeup -headless=$HEADLESS -username=$USERNAME -password=$PASSWORD -debug=$DEBUG
