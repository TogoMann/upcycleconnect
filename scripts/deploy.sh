
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
DEPLOY_DIR="${DEPLOY_DIR:-/opt/upcycleconnect}"
ENV_FILE="${ENV_FILE:-$REPO_ROOT/.env}"

if [ ! -f "$ENV_FILE" ]; then
    echo "Fichier d'environnement introuvable : $ENV_FILE"
    echo "Copiez .env.example vers .env et renseignez les valeurs de production avant de relancer ce script."
    exit 1
fi

set -a

source "$ENV_FILE"
set +a

REQUIRED_VARS=(DB_NAME DB_USER DB_PASSWORD DB_HOST DB_PORT JWT_SECRET FRONTEND_URL STRIPE_SECRET_KEY STRIPE_WEBHOOK_SECRET)
for var in "${REQUIRED_VARS[@]}"; do
    if [ -z "${!var:-}" ]; then
        echo "Variable d'environnement manquante : $var"
        exit 1
    fi
done

mkdir -p "$DEPLOY_DIR"

rsync -a --delete \
    --exclude 'node_modules' \
    --exclude '.git' \
    --exclude 'frontend/dist' \
    --exclude 'backend/tmp' \
    "$REPO_ROOT"/backend "$REPO_ROOT"/frontend "$REPO_ROOT"/scripts \
    "$REPO_ROOT"/docker-compose.prod.yml \
    "$DEPLOY_DIR"/

cp "$ENV_FILE" "$DEPLOY_DIR/.env"

cd "$DEPLOY_DIR"

docker compose -f docker-compose.prod.yml --env-file .env pull --ignore-pull-failures || true
docker compose -f docker-compose.prod.yml --env-file .env up -d --build

docker compose -f docker-compose.prod.yml --env-file .env ps

echo "Déploiement terminé. Répertoire de déploiement : $DEPLOY_DIR"
