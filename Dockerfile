FROM scratch
WORKDIR /srv

# Static files
COPY static /srv/static

# Binary
COPY app /srv/app
CMD ["./app"]