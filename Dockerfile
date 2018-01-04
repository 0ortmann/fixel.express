FROM node:9-alpine

COPY package.json package.json

RUN npm install --production

COPY public public
COPY src src
COPY server server
COPY .babelrc .babelrc

COPY webpack.config.production.js webpack.config.production.js
COPY config/config.production.js config/config.production.js

RUN npm run build

CMD npm run start

EXPOSE 3300