# Documentation Technique - Projet UpcycleConnect (UC)
**Auteur :** Équipe AlterNext  
**Version :** 1.0 (Projet Annuel 2025-2026)  
**Contexte :** Refonte complète du SI et déploiement de la plateforme d'économie circulaire.

---

## 1. Architecture Réseau Globale (RNUC)

Le système repose sur le **Routing Network Upcycle Connect (RNUC)**, simulé via EVE-NG sur un VPS dédié.

### 1.1. Plan d'Adressage et VLANs (Site Paris - Cœur)
| VLAN ID | Nom | Description | Accès |
| :--- | :--- | :--- | :--- |
| 10 | DIRECTION | Postes de la direction (S. Levy, etc.) | LAN |
| 20 | MARKETING | Équipe marketing (L. Terracid) | LAN |
| 30 | COMMERCIAL | Force de vente (R. Peuplus) | LAN |
| 40 | RH | Ressources Humaines (P. Chabrier) | LAN |
| 50 | INFORMATIQUE | DSI (N. Thavaud) et Admins | LAN |
| 60 | HMANAGER | Happyness Manager (A. Maclair) | LAN |
| 70 | REGIONAL | Directeur Régional (F. Molas) | LAN |
| 100 | DMZ | Mission 1 & Serveur Mail | DMZ / Externe |

### 1.2. Sécurité et Interconnexion
* **Firewalling :** * **Paris & Montreuil :** Cluster HA de 2 firewalls **OPNsense**.
    * **Suisse & Datacenter :** Cluster HA de 2 firewalls **pfSense**.
* **Protocoles de Routage :** Mise en œuvre de **OSPF** ou **RIP** pour le maillage RNUC.
* **Tunnels VPN :** * **Site-to-Site :** IPSec Over GRE entre Paris et les sites distants.
    * **Client-to-Site :** OpenVPN/Wireguard pour le télétravail et les techniciens à Montreuil.

---

## 2. Mission 1 : Développement et Dockerisation

L'application UpcycleConnect est déployée via une architecture de micro-services conteneurisés.

### 2.1. Stack Applicative
* **Backend :** API REST en **Go (Gin Framework)**.
* **Frontend :** 4 applications **Vue.js 3** (Public, Pro, Staff, Admin).
* **Base de données :** PostgreSQL hébergé sur la baie de stockage centrale.

### 2.2. Orchestration Docker
Le fichier `docker-compose.yml` à la racine pilote les conteneurs suivants sur le serveur **Docker Host** :
* `uc-api-backend`
* `uc-frontend-public`
* `uc-frontend-pro`
* `uc-frontend-staff`
* `uc-frontend-admin`
* `uc-db-postgres`

---

## 3. Mission 2 : Services Système et GPO

### 3.1. Active Directory (Paris)
* **OS :** Windows Server 2019/2022.
* **Rôles :** AD DS, DNS, DHCP.
* **GPO de base :** Fond d'écran imposé, mappage automatique des lecteurs réseaux, interdiction d'installation de logiciels.

### 3.2. Spécificités Suisse (Sécurité Stricte)
Conformément aux normes hors Schengen, les 4 postes suisses reçoivent des GPO critiques :
* Désactivation des ports USB.
* Blocage de l'invite de commande (CMD).
* Accès interdit aux paramètres du pare-feu Windows.
* Installation automatisée de l'agent de supervision.

### 3.3. Stockage et Sauvegarde
* **Baie de stockage :** Serveur Linux/TrueNAS centralisant toutes les données.
* **Sauvegarde :** Solution **VEEAM Backup** configurée pour des sauvegardes complètes (hebdomadaires) et incrémentielles (quotidiennes).

---

## 4. Services de Communication et Supervision

### 4.1. Téléphonie (Sites 11ème & 13ème)
* Déploiement d'un serveur **3CX** IPBX.
* Utilisation de softphones pour les salles de conférence.

### 4.2. Monitoring et Ticketing
* **GLPI :** Inventaire automatisé et gestion des incidents.
* **Zabbix :** Supervision en temps réel de l'état des serveurs, des firewalls et de la charge réseau.

---

## 5. Livrables Financiers

Tous les prix sont exprimés en Euros et soumis à une **TVA de 20%**.

* **Devis 001 :** Prestations de services (Coût Homme/Heure).
* **Devis 002 :** Infrastructure Système (Serveurs, Licences Windows).
* **Devis 003 :** Infrastructure Réseau (VPS EVE-NG, Équipements Firewall).
