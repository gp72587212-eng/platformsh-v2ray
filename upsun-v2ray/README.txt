=============================================
  UPSUN V2RAY PROXY — Format .upsun/config.yaml
=============================================

STRUCTURE DU PROJET :
  .upsun/
      config.yaml      ← config app + routes (nouveau format Upsun)
  router.php           ← proxy V2Ray
  index.html           ← page d'accueil
  README.txt

ETAPE 1 — Modifier router.php
------------------------------
Ligne 7 : remplace TON_IP_VPS par l'IP de ton VPS
  define('VPS_HOST', '123.456.789.0');

ETAPE 2 — Installer le CLI Upsun
----------------------------------
  curl -fsSL https://raw.githubusercontent.com/platformsh/cli/main/installer.sh | bash

  OU pour le CLI Upsun natif :
  curl -fsSL https://raw.githubusercontent.com/upsun/cli/main/installer.sh | bash

ETAPE 3 — Se connecter et déployer
------------------------------------
  upsun login

  git init
  git add .
  git commit -m "init"

  upsun project:create --title "v2ray-proxy" --region fr-3
  upsun push

ETAPE 4 — URL générée
----------------------
  main-XXXXXXX-YYYYYYY.fr-3.platformsh.site

ETAPE 5 — Config V2Ray client
------------------------------
"xhttpSettings": {
  "host": "main-XXXXXXX-YYYYYYY.fr-3.platformsh.site",
  "mode": "auto",
  "path": "/",
  "extra": {
    "xPaddingBytes": "100-1000",
    "scMaxEachPostBytes": "1000000"
  }
},
"vnext": [{
  "address": "TON_IP_VPS",
  "port": 443
}]

=============================================
