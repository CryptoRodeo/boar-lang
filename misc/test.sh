subjects=(parser lexer ast token)
for subject in "${subjects[@]}"; do /usr/local/go/bin/go test "./$subject"; done