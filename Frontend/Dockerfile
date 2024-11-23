# Build Stage
FROM node:20.17.0-alpine AS build

WORKDIR /app

COPY package*.json ./

RUN npm ci

COPY . .

RUN npm run build

# Dev Stage
FROM build AS dev

ENV NODE_ENV=development
ENV PATH=/app/node_modules/.bin:$PATH

EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

# Production Stage
FROM node:20.17.0-alpine AS prd

ENV NODE_ENV=production

WORKDIR /app

COPY --from=build /app/package*.json ./
COPY --from=build /app/build ./build

ARG PUBLIC_API_URL
ARG API_URL

RUN echo "PUBLIC_API_URL=$PUBLIC_API_URL" > .env
RUN echo "API_URL=$API_URL" >> .env

RUN npm ci --omit dev

EXPOSE 3000

CMD ["node", "build"]
