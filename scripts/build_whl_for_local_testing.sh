# 1. Compile the Go binary for your Mac into the Python package bin folder
mkdir -p python/src/legacy_ssh_proxy/bin
CGO_ENABLED=0 go build -ldflags="-s -w" \
  -o python/src/legacy_ssh_proxy/bin/legacy-ssh-proxy \
  ./cmd/legacy-ssh-proxy

# 2. Build and tag the wheel
python3 -m venv ./tmp/test-venv
./tmp/test-venv/bin/pip install --upgrade pip build wheel twine
./tmp/test-venv/bin/python -m build --wheel --outdir python/dist-raw python/
./tmp/test-venv/bin/wheel tags --remove --platform-tag macosx_11_0_arm64 python/dist-raw/*.whl
mkdir -p python/dist && mv python/dist-raw/*.whl python/dist/

# 3. Validate wheel metadata (checks what PyPI checks before accepting an upload)
./tmp/test-venv/bin/twine check python/dist/*.whl

# 4. Install the wheel and verify both CLI and Python import work
./tmp/test-venv/bin/pip install python/dist/*.whl
./tmp/test-venv/bin/legacy-ssh-proxy -h
./tmp/test-venv/bin/python -c "from legacy_ssh_proxy import legacy_ssh_proxy, get_binary_path; print('Binary:', get_binary_path())"
