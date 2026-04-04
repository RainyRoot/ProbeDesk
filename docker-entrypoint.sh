#!/bin/sh
Xvfb :99 -screen 0 1280x800x24 -nolisten tcp &
XVFB_PID=$!
export DISPLAY=:99

# Wait for Xvfb to be ready
for i in $(seq 1 10); do
    xdpyinfo -display :99 >/dev/null 2>&1 && break
    sleep 0.2
done

exec probedesk "$@"
