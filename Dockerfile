FROM mhart/alpine-node:6

RUN apk add --no-cache make gcc g++ python

COPY package.json package.json
COPY webpack.config.production.js webpack.config.production.js

RUN npm install --production

COPY public/img public/img 
COPY src src
COPY server server
COPY entry.jsx entry.jsx
COPY .babelrc .babelrc 
RUN npm run build

CMD npm run start

EXPOSE 3300