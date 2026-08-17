-- flowershop 数据库初始化脚本
-- 说明：docker compose 启动时后端通过 GORM AutoMigrate 自动建表，
-- 此文件为手动初始化 PostgreSQL 时使用（执行前请先 CREATE DATABASE flowershop_db）。
-- 完整建表语句见 backend/migrations/001_init.sql。
\i ../backend/migrations/001_init.sql
