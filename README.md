# user-management-api

## Requirement
1. PostgreSQL 13.20
2. Go 1.23.2
3. Postman

## Setup Instruction
1. run script untuk create database pada file migration.sql di direktori setup (ubah nama database jika diperlukan)
2. run script untuk create table pada file migration.sql di direktori setup
3. copy dan paste file .env.example di dalam root folder, rename jadi .env
4. sesuaikan value di .env dengan koneksi ke database yang digunakan
5. untuk menjalankan aplikasi run commend di root folder: go run .\cmd\main.go 

## API documentation
1. import User Management API Ridzqy.postman_collection.json ke Postman
2. running request untuk case positive dan negative

## Running Test 
1. untuk coverage test jalankan command berikut pada root folder: go test -v -cover ./...
2. untuk melakukan unit testing jalankan command berikut pada root folder: go test ./... -v