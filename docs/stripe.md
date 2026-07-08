# Paiements Stripe — Guide de mise en route

Ce document explique comment activer les paiements Stripe sur upcycleconnect, ce que vous devez faire de votre côté, et comment le système fonctionne.

## 1. Ce que le code fait déjà

Tout le code de paiement est en place et testé. Le système utilise **Stripe Checkout hébergé** : l'utilisateur est redirigé vers une page de paiement sécurisée hébergée par Stripe, puis renvoyé sur le site. Aucune donnée de carte ne transite ni n'est stockée par upcycleconnect (conformité PCI assurée par construction).

Trois cas d'usage sont couverts :

| Cas | Endpoint backend | Déclenché depuis |
| :-- | :-- | :-- |
| Abonnement Premium / Pro | `POST /subscriptions/checkout` | Page Plans (client) et Abonnements (pro) |
| Achat d'une annonce | `POST /listing-order/checkout` | Panier / paiement client |
| Encart publicitaire (mise en avant) | `POST /advertisement/{id}/checkout` | Page Publicités (pro) |

Deux endpoints complètent le flux :
- `GET /payments/verify?session_id=...` : vérifie le statut réel d'une session au retour de Stripe.
- `POST /webhooks/stripe` : reçoit la confirmation de paiement signée par Stripe et active réellement l'abonnement / la commande / la publicité.

Un abonnement, une commande ou une publicité ne passe à l'état « payé » **qu'après** confirmation du webhook, jamais à la simple création — c'est la garantie qu'aucun accès n'est accordé sans paiement réel.

## 2. Ce que vous devez faire (actions manuelles)

Ces étapes nécessitent votre compte Stripe personnel ; elles ne peuvent pas être automatisées.

### Étape 1 — Créer un compte Stripe
Rendez-vous sur https://dashboard.stripe.com/register et créez un compte. Restez en **mode test** (interrupteur en haut à droite du dashboard) tant que vous n'êtes pas prêt pour de vrais paiements.

### Étape 2 — Récupérer votre clé secrète
Dans https://dashboard.stripe.com/test/apikeys, copiez la **Clé secrète** (elle commence par `sk_test_...`).

Collez-la dans le fichier `.env` à la racine du projet :
```
STRIPE_SECRET_KEY=sk_test_51Xxxxxxxxxxxx
```

### Étape 3 — Configurer le webhook
Le webhook permet à Stripe de confirmer les paiements à votre backend.

**En développement local**, installez la CLI Stripe (https://stripe.com/docs/stripe-cli) puis lancez :
```
stripe login
stripe listen --forward-to localhost:8081/webhooks/stripe
```
La commande affiche un secret de signature (`whsec_...`). Copiez-le dans `.env` :
```
STRIPE_WEBHOOK_SECRET=whsec_xxxxxxxxxxxx
```
Laissez cette commande tournée pendant vos tests : elle relaie les événements Stripe vers votre machine.

**En production**, créez plutôt le webhook depuis https://dashboard.stripe.com/webhooks :
- URL de destination : `https://VOTRE-DOMAINE/webhooks/stripe`
- Événement à écouter : `checkout.session.completed`
- Copiez le secret de signature affiché dans la variable `STRIPE_WEBHOOK_SECRET`.

### Étape 4 — Redémarrer le backend
```
docker compose -f docker-compose.dev.yml up -d --build uc-api-dev
```
Le backend charge les nouvelles variables. Les avertissements « STRIPE_SECRET_KEY variable is not set » disparaissent.

### Étape 5 — Tester
Sur une page d'abonnement ou d'achat, lancez un paiement. Stripe affiche sa page de paiement de test. Utilisez la carte de test :
- Numéro : `4242 4242 4242 4242`
- Date : n'importe quelle date future — CVC : trois chiffres quelconques.

Après validation, vous êtes redirigé sur le site, le webhook confirme le paiement, et l'abonnement/commande devient actif.

## 3. Comment ça fonctionne (le flux complet)

```
Utilisateur clique "Payer"
        │
        ▼
Frontend  ── POST /…/checkout ──►  Backend
        │                          crée une Checkout Session Stripe
        │                          (montant, description, metadata)
        │  ◄──────  { url } ───────┘
        ▼
window.location.href = url   (redirection vers Stripe)
        │
        ▼
Page de paiement Stripe (carte saisie chez Stripe)
        │
        ├── paiement réussi ──►  Stripe redirige vers SuccessURL du site
        │                        (le frontend appelle /payments/verify
        │                         pour afficher le bon statut)
        │
        └── en parallèle ─────►  Stripe POST /webhooks/stripe (signé)
                                 Backend vérifie la signature,
                                 lit metadata.type et active réellement :
                                   • subscription → ChoosePlan
                                   • listing_order → CreatePaidOrder
                                   • advertisement → activation campagne
```

Les métadonnées (`type`, `user_id`, `plan_id`, etc.) sont attachées à la session Stripe à sa création et relues par le webhook : c'est ce qui permet au backend de savoir quoi activer sans faire confiance au frontend.

## 4. Points d'attention

- Les cartes de test (`4242…`) ne fonctionnent qu'en mode test. Pour de vrais paiements, activez votre compte Stripe (informations légales/bancaires) et remplacez les clés `sk_test_` / `whsec_` par leurs équivalents `sk_live_`.
- Le fichier `.env` n'est pas versionné (il est dans `.gitignore`) : vos clés restent locales et ne partent jamais sur GitHub.
- Si le webhook n'est pas configuré, le paiement s'effectue chez Stripe mais l'abonnement/commande ne s'active pas côté site : le webhook est la pièce indispensable du dispositif.
