FROM node:16-alpine as build

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

ARG REACT_APP_API_URL
ENV REACT_APP_API_URL=$REACT_APP_API_URL
ENV GENERATE_SOURCEMAP=false

RUN npm run build

FROM nginx:alpine

COPY nginx.http.conf /etc/nginx/conf.d/default.conf

COPY --from=build /app/build /usr/share/nginx/html

COPY docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
