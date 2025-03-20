# Configuration Terraform pour UQAM Grade Notifier

Ce répertoire contient la configuration Terraform pour déployer l'application UQAM Grade Notifier sur Google Cloud Platform.

## Prérequis

1. Installer [Terraform](https://www.terraform.io/downloads.html)
2. Installer [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
3. Avoir un compte Google Cloud avec un projet créé
4. Avoir une clé SSH générée

## Configuration

1. Copier le fichier d'exemple de variables :

   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Modifier le fichier `terraform.tfvars` avec vos valeurs :

   - `project_id` : ID de votre projet GCP
   - `region` : Région GCP (par défaut : us-central1)
   - `zone` : Zone GCP (par défaut : us-central1-a)
   - `ssh_user` : Votre nom d'utilisateur SSH
   - `ssh_pub_key_path` : Chemin vers votre clé publique SSH

3. Authentifier avec Google Cloud :
   ```bash
   gcloud auth application-default login
   ```

## Déploiement

1. Initialiser Terraform :

   ```bash
   terraform init
   ```

2. Vérifier le plan de déploiement :

   ```bash
   terraform plan
   ```

3. Appliquer la configuration :
   ```bash
   terraform apply
   ```

## Infrastructure créée

- Une instance e2-micro (gratuite) sur GCP
- Une règle de pare-feu pour autoriser le trafic HTTP/HTTPS
- Installation automatique de Docker
- Déploiement automatique de l'application via Docker Compose

## Nettoyage

Pour supprimer l'infrastructure :

```bash
terraform destroy
```

## Notes

- L'instance e2-micro est gratuite dans la limite des quotas de GCP
- Une IP publique éphémère est attribuée à l'instance
- L'application est déployée automatiquement via Docker Compose
- Les ports 80 et 443 sont ouverts pour le trafic HTTP/HTTPS
