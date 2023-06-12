#!/bin/bash

echo "Creating tasks"
curl -X POST -H "Content-Type: application/json" -d '{"text": "Task 1"}' http://localhost:8001/tasks
curl -X POST -H "Content-Type: application/json" -d '{"text": "Task 2"}' http://localhost:8001/tasks

echo "Listing tasks"
curl http://localhost:8001/tasks

echo "Getting a single task"
curl http://localhost:8001/tasks/1

echo "Updating a task"
curl -X PUT -H "Content-Type: application/json" -d '{"text": "Updated Task 1", "completed": true}' http://localhost:8001/tasks/1

echo "Deleting a task"
curl -X DELETE http://localhost:8001/tasks/1


