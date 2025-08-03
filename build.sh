#!/bin/bash

set -e  # Exit on any error

echo "Starting build process"

mkdir -p output
rm -rf output/*

echo "Building backend"
go build -o output/elliptical .

if [ -d "frontend" ]; then
    cd frontend
    
    echo "Installing frontend dependencies"
    npm install
    
    echo "Building frontend"
    npm run build
    cd ..
    
    echo "Frontend build completed!"
else
    echo "Frontend directory not found, skipping frontend build..."
fi