FROM mhart/alpine-node:6

RUN apk add --no-cache make gcc g++ python

COPY site-front/package.json package.json
COPY site-front/webpack.config.production.js webpack.config.production.js

RUN npm install --production

COPY site-front/public public
COPY site-front/src src
COPY site-front/server server
COPY site-front/entry.jsx entry.jsx
COPY site-front/.babelrc .babelrc
COPY site-front/routes.js routes.js
COPY site-front/routes.config.development.js routes.config.development.js
COPY site-front/routes.config.production.js routes.config.production.js

RUN npm run build

CMD npm run start

EXPOSE 3300