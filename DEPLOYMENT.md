# Panduan Deploy Building Report Backend ke VPS

## Prerequisites
- VPS dengan Ubuntu 20.04+ atau CentOS 7+
- Minimal 2GB RAM, 2 vCPU, 20GB Storage
- Domain yang sudah diarahkan ke IP VPS
- Akses root atau sudo ke VPS

## 1. Setup VPS Awal

### Update sistem
```bash
sudo apt update && sudo apt upgrade -y
```

### Install Docker & Docker Compose
```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Logout dan login kembali untuk apply group changes
```

### Install tools tambahan
```bash
sudo apt install -y git vim htop unzip certbot
```

## 2. Clone dan Setup Project

### Clone repository
```bash
cd /opt
sudo git clone https://github.com/your-repo/building-report-backend.git
sudo chown -R $USER:$USER building-report-backend
cd building-report-backend
```

### Setup environment variables
```bash
# Copy dan edit file environment
cp .env.example .env
vim .env
```

**Konfigurasi .env untuk production:**
```env
# Application
APP_ENV=production
APP_PORT=8080
APP_ALLOWED_ORIGINS=https://your-domain.com

# Database - GANTI PASSWORD!
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-super-secure-password
DB_NAME=building_reports
DB_SSL_MODE=disable

# Redis - GANTI PASSWORD!
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password
REDIS_DB=0

# MinIO - GANTI CREDENTIALS!
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=your-minio-access-key
MINIO_SECRET_KEY=your-minio-secret-key
MINIO_USE_SSL=false
MINIO_BUCKET_NAME=reports
MINIO_PUBLIC_URL=https://storage.your-domain.com

# JWT - GENERATE SECURE SECRET!
JWT_SECRET=your-super-secure-jwt-secret-min-32-chars
JWT_EXPIRY_HOURS=24
```

## 3. Setup SSL Certificate

### Install Certbot dan generate certificate
```bash
# Stop nginx jika running
sudo docker-compose down

# Generate certificate
sudo certbot certonly --standalone -d your-domain.com -d www.your-domain.com -d minio.your-domain.com -d storage.your-domain.com

# Copy certificates ke project directory
sudo mkdir -p nginx/ssl
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem nginx/ssl/
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem nginx/ssl/
sudo chown -R $USER:$USER nginx/ssl/
```

### Setup auto-renewal
```bash
# Edit crontab
sudo crontab -e

# Tambahkan line berikut untuk renewal otomatis setiap 2 bulan
0 12 * */2 * /usr/bin/certbot renew --quiet && /usr/local/bin/docker-compose -f /opt/building-report-backend/docker-compose.yml restart nginx
```

## 4. Konfigurasi Nginx

Edit file nginx configuration dan ganti domain:
```bash
vim nginx/sites-available/building-report.conf
```

Ganti `your-domain.com` dengan domain Anda di semua tempat.

## 5. Deploy Application

### Build dan start services
```bash
# Build dan start semua services
docker-compose up -d --build

# Cek status services
docker-compose ps

# Cek logs jika ada error
docker-compose logs -f app
docker-compose logs -f nginx
```

### Setup database (jika diperlukan migration)
```bash
# Masuk ke container app untuk run migration
docker-compose exec app /bin/sh

# Atau run migration dari host jika ada script
# docker-compose exec app ./migrate up
```

## 6. Monitoring & Maintenance

### Cek health services
```bash
# Health check endpoint
curl https://your-domain.com/health

# Cek status containers
docker-compose ps

# Monitoring logs
docker-compose logs -f --tail=100
```

### Backup Database
```bash
# Create backup script
cat > backup-db.sh << 'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker-compose exec -T postgres pg_dump -U postgres building_reports > backup_${DATE}.sql
EOF

chmod +x backup-db.sh

# Setup cron untuk backup otomatis (setiap hari jam 2 pagi)
crontab -e
# Tambahkan: 0 2 * * * /opt/building-report-backend/backup-db.sh
```

### Update Application
```bash
# Pull latest code
git pull origin main

# Rebuild dan restart
docker-compose down
docker-compose up -d --build

# Cek logs
docker-compose logs -f app
```

## 7. Firewall Setup

```bash
# Setup UFW firewall
sudo ufw enable
sudo ufw allow ssh
sudo ufw allow 80
sudo ufw allow 443

# Optional: Allow specific IPs only
# sudo ufw allow from YOUR_IP_ADDRESS to any port 22
```

## 8. Performance Tuning

### System optimizations
```bash
# Increase file limits
echo "* soft nofile 65536" | sudo tee -a /etc/security/limits.conf
echo "* hard nofile 65536" | sudo tee -a /etc/security/limits.conf

# Optimize kernel parameters
cat << 'EOF' | sudo tee -a /etc/sysctl.conf
net.core.somaxconn = 1024
net.ipv4.tcp_max_syn_backlog = 1024
vm.swappiness = 10
EOF

sudo sysctl -p
```

### Docker optimizations
```bash
# Edit docker daemon config
cat << 'EOF' | sudo tee /etc/docker/daemon.json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "storage-driver": "overlay2"
}
EOF

sudo systemctl restart docker
```

## 9. Security Checklist

- ✅ Firewall dikonfigurasi dengan benar
- ✅ SSL/TLS certificate installed
- ✅ Environment variables menggunakan password yang kuat
- ✅ Database dan Redis menggunakan authentication
- ✅ Nginx security headers dikonfigurasi
- ✅ Docker containers berjalan sebagai non-root user
- ✅ Regular backup database
- ✅ Log monitoring setup

## 10. Troubleshooting

### Common Issues

**Port already in use:**
```bash
sudo lsof -i :80
sudo lsof -i :443
# Kill process yang menggunakan port tersebut
```

**SSL Certificate issues:**
```bash
# Cek certificate expiry
openssl x509 -in nginx/ssl/fullchain.pem -text -noout | grep "Not After"

# Renew certificate manually
sudo certbot renew
```

**Database connection issues:**
```bash
# Cek database logs
docker-compose logs postgres

# Connect ke database untuk debug
docker-compose exec postgres psql -U postgres -d building_reports
```

**High memory usage:**
```bash
# Cek memory usage per container
docker stats

# Restart services jika perlu
docker-compose restart
```

### Useful Commands
```bash
# Cek semua containers
docker ps -a

# Cek resource usage
docker stats

# Cleanup unused images
docker system prune -f

# View logs dengan timestamp
docker-compose logs -f -t app

# Execute command in container
docker-compose exec app /bin/sh
```

## Support

Jika ada masalah dengan deployment, cek:
1. Logs aplikasi: `docker-compose logs -f app`
2. Logs nginx: `docker-compose logs -f nginx`
3. System logs: `sudo journalctl -fu docker`
4. Resource usage: `htop` atau `docker stats`