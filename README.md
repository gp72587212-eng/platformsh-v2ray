# upsun-v2ray

Relais XHTTP pour V2Ray/Xray, déployable sur Upsun (ex‑Platform.sh).

## Configuration requise

- Un VPS avec Xray configuré (inbound VLESS + XHTTP)
- Un compte Upsun (région fr-3 recommandée)
- Un domaine public (ex: reprise.orange-business.com) pour le SNI

## Variables d'environnement

| Variable               | Description                     | Exemple                    |
|------------------------|---------------------------------|----------------------------|
| `REMOTE_HOST`          | IP du serveur Xray              | 164.92.226.103             |
| `REMOTE_PORT`          | Port du serveur Xray            | 2083                       |
| `REMOTE_HOST_HEADER`   | Host header attendu par Xray   | upx.mjsd.ir                |

## Déploiement

1. Créer un projet sur Upsun (région fr-3)
2. Importer ce dépôt
3. Définir les variables d'environnement
4. Déployer

L'URL Upsun obtenue (ex: `main-xxxxx.fr-3.platformsh.site`) sera utilisée comme `host=` dans l'URL VLESS.

## Exemple d'URL client