# Build Stage
FROM node:20.17.0 AS build

WORKDIR /app

COPY package*.json ./

RUN npm install

# Dev Stage
FROM build AS dev

ENV NODE_ENV=development
ENV PATH /app/node_modules/.bin:$PATH

COPY . .

EXPOSE 5173

WORKDIR /app

RUN npm install vite

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
