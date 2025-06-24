# ✏️Go backend server template

  
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24.4-blue)](https://golang.org/)


## 📌 About

Personal backend project thet emerged from Go language learining. 
Project's purpose is to be a functional template for future Go projects.   

## ✨ Features

- 🐘 **Postgres integration** – via go_pg
- 👥 **User authorization** – name and passwords are stored in postgres
- 🔐 **JWT authentication** – functionality to create and validate tokens
- 🛡️ **Restricted routes** – api available only for requests with valid tokens

## 🛠 Installation
```bash
# Cloning repo
git clone https://github.com/Obezyan0941/go_tests.git
cd go_tests

# Install dependencies
go mod download

# Configure .env file
cp .env.sample .env
# input variables data

# run project
docker compose up -d  # runs postgres in docker
go run ./cmd/server   # runs server locally