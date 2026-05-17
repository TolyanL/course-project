#!/bin/bash

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
ADMIN_LOGIN="admin"
ADMIN_PASSWORD="admin123"

cleanup_all() {
  docker exec -i schedule_db psql -U postgres -d schedule_db -c "
        DELETE FROM schedule_entries;
        DELETE FROM schedules;
        DELETE FROM teachers WHERE login != 'admin';
        DELETE FROM subjects;
        DELETE FROM classrooms;
        DELETE FROM groups;
    " >/dev/null 2>&1 || true
}

cleanup_all
echo "=== Starting API Tests ==="
echo ""

TOKEN=$(curl -s "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"login\":\"$ADMIN_LOGIN\",\"password\":\"$ADMIN_PASSWORD\"}" |
  jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "ERROR: Login failed"
  exit 1
fi
echo "✓ Login successful"

GROUP_RESP=$(curl -s -X POST "$BASE_URL/api/groups" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"TestGroup"}')
GROUP_ID=$(echo "$GROUP_RESP" | jq -r '.id')
echo "✓ Created group (ID: $GROUP_ID)"

SUBJ_RESP=$(curl -s -X POST "$BASE_URL/api/subjects" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"TestSubject"}')
SUBJ_ID=$(echo "$SUBJ_RESP" | jq -r '.id')
echo "✓ Created subject (ID: $SUBJ_ID)"

CLASS_RESP=$(curl -s -X POST "$BASE_URL/api/classrooms" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"number":"202"}')
CLASS_ID=$(echo "$CLASS_RESP" | jq -r '.id')
echo "✓ Created classroom (ID: $CLASS_ID)"

TEACH_RESP=$(curl -s -X POST "$BASE_URL/api/teachers" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Teacher","login":"testteacher","password":"pass123","role":"teacher"}')
TEACH_ID=$(echo "$TEACH_RESP" | jq -r '.id')
echo "✓ Created teacher (ID: $TEACH_ID)"

ENTRY_RESP=$(curl -s -X POST "$BASE_URL/api/schedule/entries" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"group_id\":$GROUP_ID,\"date\":\"2024-02-01\",\"subject_id\":$SUBJ_ID,\"teacher_id\":$TEACH_ID,\"classroom_id\":$CLASS_ID,\"pair_number\":1}")
ENTRY_ID=$(echo "$ENTRY_RESP" | jq -r '.entry.id')
echo "✓ Created schedule entry (ID: $ENTRY_ID)"

curl -s "$BASE_URL/api/schedule?group_id=$GROUP_ID" \
  -H "Authorization: Bearer $TOKEN" | jq -e 'length > 0' >/dev/null &&
  echo "✓ Get schedule works" || echo "✗ Get schedule failed"

CONFLICT=$(curl -s -X POST "$BASE_URL/api/schedule/entries" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"group_id\":$GROUP_ID,\"date\":\"2024-02-01\",\"subject_id\":$SUBJ_ID,\"teacher_id\":$TEACH_ID,\"classroom_id\":$CLASS_ID,\"pair_number\":1}")
echo "$CONFLICT" | jq -e '.error == "conflict detected"' >/dev/null &&
  echo "✓ Conflict detection works" || echo "✗ Conflict detection failed"

UNAUTH=$(curl -s "$BASE_URL/api/teachers")
echo "$UNAUTH" | jq -e '.error' >/dev/null &&
  echo "✓ Unauthorized access blocked" || echo "✗ Auth middleware issue"

echo ""
echo "=== Delete tests (FK constraints) ==="

DELETE_CONFLICT=$(curl -s -X DELETE "$BASE_URL/api/teachers/$TEACH_ID" \
  -H "Authorization: Bearer $TOKEN")
echo "$DELETE_CONFLICT" | jq -e '.error | contains("cannot delete")' >/dev/null &&
  echo "✓ Teacher FK constraint returns proper error" || echo "✗ Teacher delete failed"

DELETE_CONFLICT=$(curl -s -X DELETE "$BASE_URL/api/subjects/$SUBJ_ID" \
  -H "Authorization: Bearer $TOKEN")
echo "$DELETE_CONFLICT" | jq -e '.error | contains("cannot delete")' >/dev/null &&
  echo "✓ Subject FK constraint returns proper error" || echo "✗ Subject delete failed"

DELETE_CONFLICT=$(curl -s -X DELETE "$BASE_URL/api/classrooms/$CLASS_ID" \
  -H "Authorization: Bearer $TOKEN")
echo "$DELETE_CONFLICT" | jq -e '.error | contains("cannot delete")' >/dev/null &&
  echo "✓ Classroom FK constraint returns proper error" || echo "✗ Classroom delete failed"

echo ""
echo "=== Cleanup ==="
cleanup_all

sleep 0.5

GROUPS=$(curl -s "$BASE_URL/api/groups" -H "Authorization: Bearer $TOKEN" 2>/dev/null | jq -r 'if type == "array" then length else 0 end' 2>/dev/null || echo "0")
if [ "$GROUPS" = "0" ]; then
  echo "✓ Cleanup successful"
else
  echo "✓ Cleanup done (DB state: $GROUPS groups)"
fi

echo ""
echo "=== All tests completed! ==="

