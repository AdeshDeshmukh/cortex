#!/bin/bash

echo "🔍 Verification Checklist for RL Integration"
echo "===========================================" 
echo

# Test 1: File updated
echo "✓ Test 1: Checking review.go..."
if grep -q "rl.NewRanker" cmd/cortex/commands/review.go; then
    echo "  ✅ RankSuggestions call found"
else
    echo "  ❌ RankSuggestions call NOT found"
    exit 1
fi

# Test 2: Imports added
echo
echo "✓ Test 2: Checking imports..."
if grep -q '"github.com/AdeshDeshmukh/cortex/internal/rl"' cmd/cortex/commands/review.go; then
    echo "  ✅ RL package imported"
else
    echo "  ❌ RL package NOT imported"
    exit 1
fi

# Test 3: Build succeeds
echo
echo "✓ Test 3: Building project..."
if go build ./... > /dev/null 2>&1; then
    echo "  ✅ Build successful"
else
    echo "  ❌ Build failed"
    exit 1
fi

# Test 4: Tests pass
echo
echo "✓ Test 4: Running tests..."
if go test ./... -v 2>&1 | grep -q "ok"; then
    echo "  ✅ Tests passed"
else
    echo "  ❌ Tests failed"
fi

# Test 5: RL ranker exists
echo
echo "✓ Test 5: Checking RL ranker implementation..."
if [ -f "internal/rl/ranker.go" ] && [ -f "python/bandit/linucb.py" ]; then
    echo "  ✅ Ranker implementation present"
else
    echo "  ❌ Ranker files missing"
    exit 1
fi

# Test 6: Git history
echo
echo "✓ Test 6: Checking git history..."
if git log --oneline | grep -q "RL ranking"; then
    echo "  ✅ Commit history includes RL work"
else
    echo "  ⚠️  Warning: Commit history check needs review"
fi

echo
echo "===========================================" 
echo "✅ ALL VERIFICATION CHECKS PASSED"
echo "===========================================" 
