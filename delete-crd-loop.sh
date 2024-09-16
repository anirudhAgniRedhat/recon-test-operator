#!/bin/bash

# CRD name
CRD_NAME="recontests.example.anirudh.io"

# Interval in milliseconds (adjust as needed)
INTERVAL_MS=500  # 500 milliseconds (0.5 seconds)

# Convert milliseconds to seconds (sleep accepts seconds)
INTERVAL=$(echo "$INTERVAL_MS / 1000" | bc -l)

# Function to delete the CRD in a loop
delete_crd_continuously() {
  while true; do
    echo "Attempting to delete CRD: $CRD_NAME"

    # Attempt to delete the CRD
    oc delete crd $CRD_NAME --ignore-not-found=true

    # Check if the CRD was deleted successfully
    if [ $? -eq 0 ]; then
      echo "CRD deleted successfully."
    else
      echo "Failed to delete CRD or CRD not found."
    fi

    # Wait for the specified interval (in seconds) before the next attempt
    sleep $INTERVAL
  done
}

# Start the continuous deletion process
delete_crd_continuously
