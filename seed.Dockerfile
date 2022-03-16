FROM node:14.17.0-alpine3.13

WORKDIR /app

COPY ./package.json /app
COPY ./yarn.lock /app
RUN yarn
COPY ./knexfile.ts /app
COPY seeds /app/seeds
CMD yarn seed
