#!/bin/bash

set -e

MODEL_DIR="models/base"
MODEL_FILE="deepseek-coder-1.3b-instruct.Q4_K_M.gguf"
MODEL_PATH="$MODEL_DIR/$MODEL_FILE"
MODEL_URL="https://huggingface.co/TheBloke/deepseek-coder-1.3b-instruct-GGUF/resolve/main/deepseek-coder-1.3b-instruct.Q4_K_M.gguf"

echo "🤖 Cortex Model Downloader"
echo ""

mkdir -p "$MODEL_DIR"

if [ -f "$MODEL_PATH" ]; then
    echo "✅ Model already exists at $MODEL_PATH"
    ls -lh "$MODEL_PATH"
    exit 0
fi

echo "📥 Downloading DeepSeek-Coder 1.3B (1.2GB)..."
echo "   This will take 2-10 minutes depending on connection"
echo ""

curl -L \
    --progress-bar \
    --output "$MODEL_PATH.tmp" \
    "$MODEL_URL"

mv "$MODEL_PATH.tmp" "$MODEL_PATH"

FILE_SIZE=$(stat -f%z "$MODEL_PATH" 2>/dev/null || stat -c%s "$MODEL_PATH" 2>/dev/null)
MIN_SIZE=1000000000

if [ "$FILE_SIZE" -lt "$MIN_SIZE" ]; then
    echo "❌ Download failed: File too small ($FILE_SIZE bytes)"
    rm "$MODEL_PATH"
    exit 1
fi

chmod 444 "$MODEL_PATH"

echo ""
echo "✅ Model downloaded successfully"
ls -lh "$MODEL_PATH"
echo ""
echo "🚀 Ready to use! Run: cortex review"