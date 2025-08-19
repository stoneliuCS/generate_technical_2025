#!/bin/bash

if [ -z "$1" ] || [ "$1" != "--email" ] || [ -z "$2" ]; then
    echo "Usage: ./candidate --email candidate@example.com"
    exit 1
fi

EMAIL="$2"

echo "Looking up stats for: $EMAIL"
echo "================================"

SUMMARY=$(PGPASSWORD="" psql \
    -h \
    -p 5432 \
    -U \
    -d postgres \
    -t \
    -A \
    -c "
    SELECT 
    m.email,
        COALESCE(MIN(CASE WHEN s.score != -1 THEN s.score END), 0) as best_score,
        COALESCE(COUNT(s.id), 0) as attempts,
        m.created_at::date as registered
    FROM members m 
    LEFT JOIN scores s ON m.id = s.user_id 
    WHERE m.email = '$EMAIL'
    GROUP BY m.id, m.email, m.created_at;
    ")

if [ -z "$SUMMARY" ]; then
    echo "Candidate not found"
    exit 1
fi

IFS='|' read -r email best_score attempts avg_score registered <<< "$SUMMARY"

echo "Email: $email"
echo "Best Score: $best_score"
echo "Total Attempts: $attempts"
echo "Registered: $registered"
echo ""
echo "All Scores:"
echo "----------"

PGPASSWORD="" psql \
    -h \
    -p 5432 \
    -U \
    -d postgres \
    -c "
    SELECT 
        CASE 
            WHEN s.score = -1 THEN 'FAILED (-1)'
            ELSE s.score::text
        END as score,
        s.created_at::timestamp as submitted_at,
        s.challenge_type
    FROM members m 
    JOIN scores s ON m.id = s.user_id 
    WHERE m.email = '$EMAIL'
    ORDER BY s.created_at DESC;
    "

FRONTEND_STATS=$(PGPASSWORD="" psql \
    -h \
    -p 5432 \
    -U \
    -d postgres \
    -t \
    -A \
    -c "
    SELECT
        COALESCE(COUNT(f.id), 0) as request_count,
        MIN(f.timestamp)::date as first_request,
        MAX(f.timestamp)::date as last_request
    FROM members m
    LEFT JOIN frontend_usages f ON m.id::uuid = f.user_id
    WHERE m.email = '$EMAIL';
    ")

IFS='|' read -r frontend_requests first_frontend last_frontend <<< "$FRONTEND_STATS"

echo "Frontend Usage:"
echo "---------------"
echo "Total Requests: $frontend_requests"
if [ "$frontend_requests" -gt 0 ]; then
    echo "First Request: $first_frontend"
    echo "Last Request: $last_frontend"
else
    echo "No frontend requests made"
fi