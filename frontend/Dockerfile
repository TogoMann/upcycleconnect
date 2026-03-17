FROM oven/bun:alpine

WORKDIR /app

COPY package.json ./
COPY bun.lock ./
COPY . .

RUN bun install

EXPOSE 5173

CMD [ "bun", "dev", "--", "--host"]
