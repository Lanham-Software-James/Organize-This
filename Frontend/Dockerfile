# Build Stage
FROM node:20.17.0 AS build

ENV NODE_ENV=development

WORKDIR /app

COPY package*.json ./

RUN npm install

# Dev Stage
FROM build AS dev

COPY . .

EXPOSE 5173

WORKDIR /app

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]