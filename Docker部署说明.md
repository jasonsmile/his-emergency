# Docker 部署说明

## 1. 构建服务镜像

```bash
docker build -t emergency-his-server:0.1.0 ./server
```

## 2. 导出镜像

```bash
docker save emergency-his-server:0.1.0 -o emergency-his-server-0.1.0.tar
```

## 3. 导入镜像

```bash
docker load -i emergency-his-server-0.1.0.tar
```

## 4. 部署目录结构

```text
/opt/emergency-his/
├── deploy/
│   └── docker-compose.yml
└── server/
    └── config/
        └── config.yaml
```

## 5. 启动服务

```bash
cd /opt/emergency-his
sudo docker compose -f deploy/docker-compose.yml up -d
```
