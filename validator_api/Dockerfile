FROM node:14.17.0-alpine3.13

WORKDIR /app

ADD package.json /app
ADD yarn.lock /app
ADD tsconfig.json /app
RUN yarn

ADD /src /app/src
RUN yarn build
CMD [ "node", "dist/index.js" ]
