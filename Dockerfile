FROM ubuntu:19.04
WORKDIR "/srv"
COPY . /srv
CMD ["./app"]