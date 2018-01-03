FROM node:9-alpine

COPY package.json package.json

RUN npm install --production

COPY public public
COPY src src
COPY server server
COPY entry.jsx entry.jsx
COPY .babelrc .babelrc
COPY routes.js routes.js

COPY webpack.config.production.js webpack.config.production.js
COPY config.development.js config.development.js
COPY config.production.js config.production.js

RUN npm run build

CMD npm run start

EXPOSE 3300