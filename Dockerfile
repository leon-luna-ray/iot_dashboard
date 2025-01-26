# Stage 1: Build the client
FROM node:lts-alpine as client

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

RUN npm install -g pnpm

COPY client/package.json client/pnpm-lock.yaml ./

RUN pnpm install

COPY client/ ./

# Add logging to verify the build step
RUN echo "🚧 Building client..." && pnpm run build && echo "✅ Client build complete."

# Stage 2: Build the server and serve client build
FROM python:3.12-slim-bullseye

WORKDIR /app

ENV PATH="${PATH}:/root/.local/bin"
ENV PYTHONPATH=.

RUN pip install --upgrade pip

COPY server/requirements.txt ./
COPY server/src/ ./src/

RUN pip install -r requirements.txt

EXPOSE 9090

CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "9090", "--reload"]