#!/bin/bash

# Interval in milliseconds (adjust as needed)
INTERVAL_MS=10  # 100 milliseconds

# Convert milliseconds to seconds (sleep accepts seconds)
INTERVAL=$(echo "$INTERVAL_MS / 1000" | bc -l)

# Function to delete CRDs with specific naming pattern
delete_reconciler_crds() {
  while true; do
    echo "Searching for CRDs to delete..."

    # Find and delete CRDs matching the pattern recontests*.example.anirudh.io
    oc get crd | grep "recontests[0-9]*.example.anirudh.io" | awk '{print $1}' | while read -r crd; do
      echo "Attempting to delete CRD: $crd"

      # Attempt to delete the CRD
      oc delete crd "$crd" --ignore-not-found=true

      if [ $? -eq 0 ]; then
        echo "CRD $crd deleted successfully."
      else
        echo "Failed to delete CRD $crd."
      fi
    done

    # Wait for the specified interval (in seconds) before the next attempt
    sleep $INTERVAL
  done
}

# Handle script interruption gracefully
trap 'echo "Stopping CRD deletion script..."; exit 0' SIGINT SIGTERM

# Start the continuous deletion process
delete_reconciler_crds