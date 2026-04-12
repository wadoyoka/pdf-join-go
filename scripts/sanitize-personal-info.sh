#!/usr/bin/env bash
set -euo pipefail

USERNAME="$(whoami)"
PC_NAME="$(hostname -s)"

found=0

for file in "$@"; do
  # Skip binary files
  if file "$file" | grep -q 'binary'; then
    continue
  fi

  if grep -qE "$PC_NAME|/Users/$USERNAME" "$file" 2>/dev/null; then
    echo "Sanitizing: $file"
    sed -i '' "s|$PC_NAME|pcName|g" "$file"
    sed -i '' "s|/Users/$USERNAME|/Users/userName|g" "$file"
    git add "$file"
    found=1
  fi
done

if [ "$found" -eq 1 ]; then
  echo "Personal info sanitized and re-staged."
fi
