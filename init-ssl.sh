#!/bin/bash
set -euo pipefail

trap 'echo "[ERROR] Script interrupted" >&2; exit 1' INT TERM

if ! docker compose version >/dev/null 2>&1; then
  echo 'Error: docker compose is not available.' >&2
  exit 1
fi

log() {
  printf '[%s] %s\n' "$(date +"%Y-%m-%d %H:%M:%S")" "$1"
}

MIN_FREE_MEM=350   # MB
MAX_LOAD=250       # load * 100 (2.5)
WAIT_INTERVAL=30   # seconds
MAX_WAIT_CYCLES=10

wait_for_capacity() {
  local phase="$1"
  local cycles=0
  while true; do
    local mem_available
    mem_available=$(awk '/MemAvailable/ {printf "%d", $2/1024}' /proc/meminfo)
    local load_scaled
    load_scaled=$(awk '{printf "%d", $1*100}' /proc/loadavg)
    local load_display
    load_display=$(awk '{print $1}' /proc/loadavg)

    if (( mem_available >= MIN_FREE_MEM && load_scaled <= MAX_LOAD )); then
      log "Resource check OK (avail ${mem_available}MB, load ${load_display})"
      break
    fi

    cycles=$((cycles + 1))
    if (( cycles >= MAX_WAIT_CYCLES )); then
      log "Resource pressure persists before ${phase}, continuing cautiously."
      break
    fi

    log "Resources tight before ${phase} (avail ${mem_available}MB, load ${load_display}). Sleeping ${WAIT_INTERVAL}s ..."
    sleep ${WAIT_INTERVAL}
  done
}

pull_image_if_needed() {
  local image="$1"
  if docker image inspect "$image" >/dev/null 2>&1; then
    log "Image ${image} already present, skip pull"
    return
  fi
  wait_for_capacity "pulling ${image}"
  log "Pulling image ${image} ..."
  docker pull "$image"
}

ensure_tls_params() {
  local data_path="$1"
  if [ -e "$data_path/conf/options-ssl-nginx.conf" ] && [ -e "$data_path/conf/ssl-dhparams.pem" ]; then
    return
  fi
  log "Downloading recommended TLS parameters ..."
  mkdir -p "$data_path/conf"
  curl -s \
    "https://raw.githubusercontent.com/certbot/certbot/master/certbot-nginx/certbot_nginx/_internal/tls_configs/options-ssl-nginx.conf" \
    -o "$data_path/conf/options-ssl-nginx.conf"
  curl -s "https://raw.githubusercontent.com/certbot/certbot/master/certbot/certbot/ssl-dhparams.pem" \
    -o "$data_path/conf/ssl-dhparams.pem"
}

start_nginx_minimal() {
  log "Starting nginx only (no dependent services) ..."
  wait_for_capacity "starting nginx"
  docker compose up -d --no-deps nginx
  sleep 5
}

stop_nginx() {
  if docker compose ps nginx >/dev/null 2>&1; then
    log "Stopping nginx to free resources ..."
    docker compose stop nginx >/dev/null 2>&1 || true
  fi
}

create_dummy_certificate() {
  local domains=("$@")
  local rsa_key_size=4096
  local domain="${domains[0]}"
  local path="/etc/letsencrypt/live/${domain}"
  log "Creating dummy certificate for ${domains[*]}"
  mkdir -p "./certbot/conf/live/${domain}"
  docker compose run --rm --entrypoint /bin/sh certbot -c \
    "openssl req -x509 -nodes -newkey rsa:${rsa_key_size} -days 1 -keyout '${path}/privkey.pem' -out '${path}/fullchain.pem' -subj '/CN=localhost'"
}

cleanup_dummy_certificate() {
  local domain="$1"
  log "Deleting dummy certificate for ${domain}"
  docker compose run --rm --entrypoint /bin/sh certbot -c \
    "rm -Rf /etc/letsencrypt/live/${domain} /etc/letsencrypt/archive/${domain} /etc/letsencrypt/renewal/${domain}.conf"
}

request_certificate() {
  local domains=("$@")
  local rsa_key_size=4096

  local domain_args=""
  for domain in "${domains[@]}"; do
    domain_args+=" -d ${domain}"
  done

  local email_arg
  if [ -z "$email" ]; then
    email_arg="--register-unsafely-without-email"
  else
    email_arg="-m ${email}"
  fi

  local staging_arg=""
  if [ "$staging" != "0" ]; then
    staging_arg="--staging"
  fi

  log "Requesting Let's Encrypt certificate for ${domains[*]}"
  docker compose run --rm --entrypoint /bin/sh certbot -c \
    "certbot certonly --webroot -w /var/www/certbot ${staging_arg} ${email_arg} ${domain_args} --rsa-key-size ${rsa_key_size} --agree-tos --force-renewal"
}

reload_nginx() {
  log "Reloading nginx with new certificates ..."
  docker compose exec nginx nginx -s reload
}

# --- Script configuration ---
domains=(api.abfun.me)
data_path="./certbot"
email="admin@abfun.me"
staging=0

if [ -d "$data_path" ]; then
  read -p "Existing data found for ${domains[*]}. Continue and replace existing certificate? (y/N) " decision
  if [ "$decision" != "Y" ] && [ "$decision" != "y" ]; then
    exit 0
  fi
fi

ensure_tls_params "$data_path"

pull_image_if_needed "nginx:alpine"
pull_image_if_needed "certbot/certbot"

stop_nginx

create_dummy_certificate "${domains[@]}"

start_nginx_minimal

cleanup_dummy_certificate "${domains[0]}"

request_certificate "${domains[@]}"

reload_nginx

log "Certificate workflow completed. Consider running 'docker compose up -d' when system load is low to start all services."
