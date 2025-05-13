#!/bin/sh

# This script debugs Next.js build issues
echo "=== ENVIRONMENT ==="
pwd
node -v
npm -v

echo "\n=== PACKAGE.JSON SCRIPTS ==="
cat package.json | grep -A 10 '"scripts"'

echo "\n=== BUILD ATTEMPT ==="
npm run build

echo "\n=== BUILD OUTPUT ==="
if [ -d ".next" ]; then
  echo ".next directory exists"
  ls -la .next
  
  echo "\n=== CHECKING FOR BUILD ID ==="
  if [ -f ".next/BUILD_ID" ]; then
    echo "BUILD_ID exists:"
    cat .next/BUILD_ID
  else
    echo "BUILD_ID is missing!"
    find .next -type f | grep -i id
  fi
  
  echo "\n=== SERVER DIRECTORY ==="
  if [ -d ".next/server" ]; then
    ls -la .next/server
  else
    echo ".next/server directory is missing!"
  fi
else
  echo ".next directory is missing completely!"
fi

echo "\n=== END OF DEBUG INFO ==="