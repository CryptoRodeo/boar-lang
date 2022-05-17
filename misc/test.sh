subjects=(parser lexer ast token evaluator object)
for subject in "${subjects[@]}"; do /usr/local/go/bin/go test "./$subject"; done