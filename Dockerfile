FROM golang as builder
MAINTAINER docker@ekito.fr

WORKDIR app

ADD . /src
RUN cd /src && make linux_amd64

FROM chromedp/headless-shell
COPY --from=builder /src/build/servicenow-instance-wakeup-linux-amd64 /app/servicenow-instance-wakeup
RUN chmod a+x /app/servicenow-instance-wakeup

RUN apt-get update && apt-get install -y \
    cron \
    golang \
    git

COPY servicenow /etc/cron.d/servicenow

RUN chmod 0744 /etc/cron.d/servicenow

COPY . /app

# Apply cron job
RUN crontab /etc/cron.d/servicenow

CMD ["cron", "-f"]
