#!/bin/bash

# Test script for debate generation endpoints
BASE_URL="http://localhost:8080/api"

echo "Testing Debate Generation Endpoints"
echo "=================================="

# Test 0: Health Check
echo -e "\n0. Testing GET /debates/health"
echo "--------------------------------"
curl -s "$BASE_URL/debates/health" | jq '.'

# Test 1: Generate AI Prompt (GET)
echo -e "\n1. Testing GET /debates/generate (AI Prompt only)"
echo "------------------------------------------------"
curl -s "$BASE_URL/debates/generate?match_id=1321727&type=pre_match" | jq '.'

# Test 2: Generate Complete Debate (POST)
echo -e "\n2. Testing POST /debates/generate (Complete debate)"
echo "----------------------------------------------------"
curl -s -X POST "$BASE_URL/debates/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "match_id": "1321727",
    "debate_type": "pre_match",
    "force_regenerate": false
  }' | jq '.'

# Test 3: Generate Post-Match Debate
echo -e "\n3. Testing POST /debates/generate (Post-match)"
echo "------------------------------------------------"
curl -s -X POST "$BASE_URL/debates/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "match_id": "1321727",
    "debate_type": "post_match",
    "force_regenerate": false
  }' | jq '.'

# Test 4: Get Debates by Match
echo -e "\n4. Testing GET /debates/match"
echo "-----------------------------"
curl -s "$BASE_URL/debates/match?match_id=1321727" | jq '.'

# Test 5: Error handling - Invalid match ID
echo -e "\n5. Testing Error Handling - Invalid match ID"
echo "----------------------------------------------"
curl -s -X POST "$BASE_URL/debates/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "match_id": "invalid",
    "debate_type": "pre_match"
  }' | jq '.'

# Test 6: Error handling - Invalid debate type
echo -e "\n6. Testing Error Handling - Invalid debate type"
echo "------------------------------------------------"
curl -s -X POST "$BASE_URL/debates/generate" \
  -H "Content-Type: application/json" \
  -d '{
    "match_id": "1321727",
    "debate_type": "invalid_type"
  }' | jq '.'

echo -e "\nTests completed!" 