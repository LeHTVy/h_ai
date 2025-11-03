#!/bin/bash

# H-AI API Test Script
# Usage: ./test_api.sh

BASE_URL="http://localhost:8888"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🧪 Testing H-AI API...${NC}"
echo ""

# Test 1: Health Check
echo -e "${YELLOW}1. Health Check...${NC}"
health=$(curl -s "$BASE_URL/health")
echo "$health" | jq . 2>/dev/null || echo "$health"
echo ""

# Test 2: Analyze Target
echo -e "${YELLOW}2. Analyze Target...${NC}"
analyze=$(curl -s -X POST "$BASE_URL/api/intelligence/analyze-target" \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com", "analysis_type": "comprehensive"}')
echo "$analyze" | jq . 2>/dev/null || echo "$analyze"
echo ""

# Test 3: Select Tools
echo -e "${YELLOW}3. Select Tools...${NC}"
tools=$(curl -s -X POST "$BASE_URL/api/intelligence/select-tools" \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com", "target_type": "web_application"}')
echo "$tools" | jq . 2>/dev/null || echo "$tools"
echo ""

# Test 4: Optimize Parameters
echo -e "${YELLOW}4. Optimize Parameters...${NC}"
params=$(curl -s -X POST "$BASE_URL/api/intelligence/optimize-parameters" \
  -H "Content-Type: application/json" \
  -d '{
    "tool": "nmap",
    "parameters": {"target": "example.com"},
    "context": {"target": "example.com"}
  }')
echo "$params" | jq . 2>/dev/null || echo "$params"
echo ""

# Test 5: Create Attack Chain
echo -e "${YELLOW}5. Create Attack Chain...${NC}"
chain=$(curl -s -X POST "$BASE_URL/api/intelligence/create-attack-chain" \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com", "analysis_type": "comprehensive"}')
echo "$chain" | jq . 2>/dev/null || echo "$chain"
echo ""

# Test 6: Process Dashboard
echo -e "${YELLOW}6. Process Dashboard...${NC}"
dashboard=$(curl -s "$BASE_URL/api/processes/dashboard")
echo "$dashboard" | jq . 2>/dev/null || echo "$dashboard"
echo ""

# Test 7: Cache Stats
echo -e "${YELLOW}7. Cache Stats...${NC}"
cache=$(curl -s "$BASE_URL/api/cache/stats")
echo "$cache" | jq . 2>/dev/null || echo "$cache"
echo ""

echo -e "${GREEN}✅ Tests completed!${NC}"
echo ""
echo "💡 To test security tools (nmap, gobuster, etc.), ensure tools are installed:"
echo "   sudo apt install -y nmap gobuster nuclei sqlmap"
