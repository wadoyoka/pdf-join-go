#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(git rev-parse --show-toplevel)"
cd "$REPO_ROOT"

GO_LICENSES="$(go env GOPATH)/bin/go-licenses"

if [ ! -x "$GO_LICENSES" ]; then
  echo "go-licenses not found. Installing..."
  go install github.com/google/go-licenses@latest
fi

# Regenerate THIRD_PARTY_LICENSES/
rm -rf THIRD_PARTY_LICENSES
"$GO_LICENSES" save . --save_path=THIRD_PARTY_LICENSES

# Regenerate CREDITS.txt and cmd/credits.txt
python3 -c "
import os

base = 'THIRD_PARTY_LICENSES'
output = []

for root, dirs, files in sorted(os.walk(base)):
    for f in sorted(files):
        filepath = os.path.join(root, f)
        rel = os.path.relpath(filepath, base)
        if rel.startswith('github.com/wadoyoka'):
            continue
        lib = '/'.join(rel.split('/')[:-1])
        output.append('=' * 60)
        output.append(lib)
        output.append('=' * 60)
        with open(filepath) as fh:
            output.append(fh.read().strip())
        output.append('')

print('\n'.join(output))
" > CREDITS.txt

cp CREDITS.txt cmd/credits.txt

echo "License files updated."
