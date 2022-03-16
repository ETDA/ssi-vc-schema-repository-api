FROM node:14.17.0-alpine3.13

WORKDIR /app

ADD package.json /app
ADD yarn.lock /app
RUN yarn

ADD /src /app/src

CMD [ "yarn", "dev" ]