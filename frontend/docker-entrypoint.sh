#!/bin/sh

if [ "$ENV" = "dev" ]; then
  echo "Using development Nginx configuration..."
  cp /etc/nginx/templates/nginx.dev.conf /etc/nginx/conf.d/default.conf
else
  echo "Using production Nginx configuration..."
  cp /etc/nginx/templates/nginx.prod.conf /etc/nginx/conf.d/default.conf
fi

exec "$@"
