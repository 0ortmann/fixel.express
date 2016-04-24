FROM node:5

COPY index.html index.html
COPY expressServer.js expressServer.js
COPY package.json package.json
COPY webpack.config.production.js webpack.config.production.js

RUN npm install --production

COPY public/img public/img 
COPY src src
COPY entry.jsx entry.jsx
RUN npm run build
RUN rm -rf src

CMD npm run start

EXPOSE 3300