FROM node:19.6-alpine

RUN apk update && apk add bash

WORKDIR /poeticmetric

# copy only package definition files
COPY package.json .
COPY package-lock.json .

# install dependencies
RUN npm install

COPY @types @types
COPY blog blog
COPY components components
COPY contexts contexts
COPY docs docs
COPY helpers helpers
COPY hooks hooks
COPY lib lib
COPY pages pages
COPY public public
COPY styles styles
COPY next.config.js .
COPY next-sitemap.config.js .
COPY sentry.client.config.js .
COPY sentry.server.config.js .
COPY tsconfig.json .

COPY docker-entrypoint.development.sh /usr/local/bin/docker-entrypoint.sh
COPY scripts/generate-config.sh /usr/local/bin/generate-config.sh

EXPOSE 80

CMD ["docker-entrypoint.sh"]
