#!/bin/bash

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
API_URL="$BASE_URL/api"

cleanup() {
  echo -e "\n=== Cleanup ==="
  if [ -n "$TOKEN" ]; then
    echo "Cleaning up test data..."
    for id in $TEACHER_IDS; do
      curl -s -X DELETE "$API_URL/teachers/$id" -H "Authorization: Bearer $TOKEN" >/dev/null || true
    done
    for id in $SUBJECT_IDS; do
      curl -s -X DELETE "$API_URL/subjects/$id" -H "Authorization: Bearer $TOKEN" >/dev/null || true
    done
    for id in $CLASSROOM_IDS; do
      curl -s -X DELETE "$API_URL/classrooms/$id" -H "Authorization: Bearer $TOKEN" >/dev/null || true
    done
    for id in $GROUP_IDS; do
      curl -s -X DELETE "$API_URL/groups/$id" -H "Authorization: Bearer $TOKEN" >/dev/null || true
    done
    echo "Cleanup done"
  fi
}

trap cleanup EXIT

echo "=== Testing Auth ==="
LOGIN_RESP=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"login":"admin","password":"admin123"}')

TOKEN=$(echo $LOGIN_RESP | jq -r '.token')
if [ "$TOKEN" = "null" ] || [ -z "$TOKEN" ]; then
  echo "FAIL: Login failed"
  echo $LOGIN_RESP | jq .
  exit 1
fi
echo "PASS: Login successful"
echo "Token: ${TOKEN:0:50}..."

echo -e "\n=== Testing GET Teachers ==="
TEACHERS=$(curl -s "$API_URL/teachers" -H "Authorization: Bearer $TOKEN")
COUNT=$(echo $TEACHERS | jq '. | length')
if [ "$COUNT" -gt 0 ]; then
  echo "PASS: Got $COUNT teachers"
else
  echo "FAIL: No teachers returned"
  exit 1
fi

echo -e "\n=== Testing CRUD Groups ==="
echo "Creating group..."
GROUP=$(curl -s -X POST "$API_URL/groups" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test-Group-001"}')
GROUP_ID=$(echo $GROUP | jq -r '.id')
if [ "$GROUP_ID" != "null" ] && [ -n "$GROUP_ID" ]; then
  echo "PASS: Created group with ID=$GROUP_ID"
  GROUP_IDS="$GROUP_IDS $GROUP_ID"
else
  echo "FAIL: Failed to create group"
  exit 1
fi

echo "Updating group..."
UPDATED=$(curl -s -X PUT "$API_URL/groups/$GROUP_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test-Group-001-Updated"}')
NAME=$(echo $UPDATED | jq -r '.name')
if [ "$NAME" = "Test-Group-001-Updated" ]; then
  echo "PASS: Group updated to '$NAME'"
else
  echo "FAIL: Group update failed"
  exit 1
fi

echo -e "\n=== Testing CRUD Subjects ==="
echo "Creating subject..."
SUBJECT=$(curl -s -X POST "$API_URL/subjects" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test-Math"}')
SUBJECT_ID=$(echo $SUBJECT | jq -r '.id')
if [ "$SUBJECT_ID" != "null" ] && [ -n "$SUBJECT_ID" ]; then
  echo "PASS: Created subject with ID=$SUBJECT_ID"
  SUBJECT_IDS="$SUBJECT_IDS $SUBJECT_ID"
else
  echo "FAIL: Failed to create subject"
  exit 1
fi

echo -e "\n=== Testing CRUD Classrooms ==="
echo "Creating classroom..."
CLASSROOM=$(curl -s -X POST "$API_URL/classrooms" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"number":"T-101"}')
CLASSROOM_ID=$(echo $CLASSROOM | jq -r '.id')
if [ "$CLASSROOM_ID" != "null" ] && [ -n "$CLASSROOM_ID" ]; then
  echo "PASS: Created classroom with ID=$CLASSROOM_ID"
  CLASSROOM_IDS="$CLASSROOM_IDS $CLASSROOM_ID"
else
  echo "FAIL: Failed to create classroom"
  exit 1
fi

echo -e "\n=== Testing CRUD Teachers ==="
echo "Creating teacher..."
TEACHER=$(curl -s -X POST "$API_URL/teachers" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Teacher","login":"test_teacher","password":"test123","role":"teacher"}')
TEACHER_ID=$(echo $TEACHER | jq -r '.id')
if [ "$TEACHER_ID" != "null" ] && [ -n "$TEACHER_ID" ]; then
  echo "PASS: Created teacher with ID=$TEACHER_ID"
  TEACHER_IDS="$TEACHER_IDS $TEACHER_ID"
else
  echo "FAIL: Failed to create teacher"
  exit 1
fi

echo -e "\n=== Testing Schedule Creation ==="
echo "Creating schedule entry..."
ENTRY=$(curl -s -X POST "$API_URL/schedule/entries" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"group_id\":$GROUP_ID,\"date\":\"2024-06-15\",\"subject_id\":$SUBJECT_ID,\"teacher_id\":$TEACHER_ID,\"classroom_id\":$CLASSROOM_ID,\"pair_number\":1}")
ENTRY_ID=$(echo $ENTRY | jq -r '.entry.id // empty')
if [ -n "$ENTRY_ID" ] && [ "$ENTRY_ID" != "null" ] && [ "$ENTRY_ID" -gt 0 ]; then
  echo "PASS: Created schedule entry with ID=$ENTRY_ID"
else
  echo "FAIL: Failed to create schedule entry"
  echo $ENTRY | jq .
  exit 1
fi

echo -e "\n=== Testing Schedule Conflict Detection ==="
echo "Trying to create conflicting entry (same pair, same classroom, same teacher)..."
CONFLICT=$(curl -s -X POST "$API_URL/schedule/entries" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"group_id\":$GROUP_ID,\"date\":\"2024-06-15\",\"subject_id\":$SUBJECT_ID,\"teacher_id\":$TEACHER_ID,\"classroom_id\":$CLASSROOM_ID,\"pair_number\":1}")
ERROR=$(echo $CONFLICT | jq -r '.error // empty')
if [ "$ERROR" = "conflict detected" ]; then
  echo "PASS: Conflict detected correctly"
  DETAILS=$(echo $CONFLICT | jq '.details | length')
  echo "  Found $DETAILS conflict warnings"
else
  echo "FAIL: Conflict not detected"
  echo $CONFLICT | jq .
  exit 1
fi

echo -e "\n=== Testing Different Pair (No Conflict) ==="
echo "Creating entry with different pair..."
ENTRY2=$(curl -s -X POST "$API_URL/schedule/entries" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"group_id\":$GROUP_ID,\"date\":\"2024-06-15\",\"subject_id\":$SUBJECT_ID,\"teacher_id\":$TEACHER_ID,\"classroom_id\":$CLASSROOM_ID,\"pair_number\":2}")
ENTRY2_ID=$(echo $ENTRY2 | jq -r '.entry.id // empty')
if [ -n "$ENTRY2_ID" ] && [ "$ENTRY2_ID" != "null" ] && [ "$ENTRY2_ID" -gt 0 ]; then
  echo "PASS: Created second entry with ID=$ENTRY2_ID (different pair)"
else
  echo "FAIL: Failed to create entry with different pair"
  echo $ENTRY2 | jq .
  exit 1
fi

echo -e "\n=== Testing Get Schedule ==="
SCHEDULE=$(curl -s "$API_URL/schedule?group_id=$GROUP_ID" \
  -H "Authorization: Bearer $TOKEN")
PAIRS=$(echo $SCHEDULE | jq '.[0].pairs | length')
if [ "$PAIRS" -ge 2 ]; then
  echo "PASS: Schedule has $PAIRS pairs"
else
  echo "FAIL: Expected at least 2 pairs, got $PAIRS"
  exit 1
fi

echo -e "\n=== Testing Schedule Entry Update ==="
if [ -n "$ENTRY_ID" ] && [ "$ENTRY_ID" != "null" ]; then
  echo "Updating entry $ENTRY_ID to pair 3..."
  UPDATED=$(curl -s -X PUT "$API_URL/schedule/entries/$ENTRY_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"group_id\":$GROUP_ID,\"date\":\"2024-06-15\",\"subject_id\":$SUBJECT_ID,\"teacher_id\":$TEACHER_ID,\"classroom_id\":$CLASSROOM_ID,\"pair_number\":3}")
  UPDATED_PAIR=$(echo $UPDATED | jq -r '.entry.pair_number // empty')
  if [ "$UPDATED_PAIR" = "3" ]; then
    echo "PASS: Entry updated to pair 3"
  else
    echo "FAIL: Entry update failed"
    echo $UPDATED | jq .
    exit 1
  fi
fi

echo -e "\n=== Testing Schedule Entry Delete ==="
if [ -n "$ENTRY2_ID" ] && [ "$ENTRY2_ID" != "null" ]; then
  echo "Deleting entry $ENTRY2_ID..."
  DELETE_RESP=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$API_URL/schedule/entries/$ENTRY2_ID" \
    -H "Authorization: Bearer $TOKEN")
  if [ "$DELETE_RESP" = "204" ]; then
    echo "PASS: Entry deleted (HTTP 204)"
  else
    echo "FAIL: Entry delete failed (HTTP $DELETE_RESP)"
    exit 1
  fi
fi

echo -e "\n=== Testing Unauthorized Access ==="
RESPONSE=$(curl -s -w "%{http_code}" -o /dev/null "$API_URL/teachers")
if [ "$RESPONSE" = "401" ]; then
  echo "PASS: Unauthorized access blocked (HTTP 401)"
else
  echo "FAIL: Expected 401, got $RESPONSE"
  exit 1
fi

echo -e "\n=== Testing CORS ==="
CORS=$(curl -s -X OPTIONS "$API_URL/teachers" \
  -H "Origin: http://example.com" \
  -H "Access-Control-Request-Method: GET" \
  -w "%{http_code}" -o /dev/null)
if [ "$CORS" = "204" ] || [ "$CORS" = "200" ]; then
  echo "PASS: CORS preflight handled (HTTP $CORS)"
else
  echo "WARN: CORS preflight returned HTTP $CORS"
fi

echo -e "\n=== Testing Invalid Login ==="
INVALID=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"login":"admin","password":"wrongpassword"}')
ERROR_MSG=$(echo $INVALID | jq -r '.error // empty')
if [ "$ERROR_MSG" = "invalid credentials" ]; then
  echo "PASS: Invalid login rejected"
else
  echo "FAIL: Invalid login not rejected properly"
  exit 1
fi

echo -e "\n========================================="
echo "ALL TESTS PASSED!"
echo "========================================="

exit 0
