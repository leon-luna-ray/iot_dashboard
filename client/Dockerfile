# Frontend development Dockerfile
FROM node:lts-alpine

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

RUN npm install -g pnpm

COPY package.json .
COPY pnpm-lock.yaml .

RUN pnpm install

EXPOSE 5173

CMD ["pnpm", "run", "dev", "--host"]
