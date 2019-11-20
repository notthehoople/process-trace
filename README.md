# process-trace
Process the output from bpftrace
## Description
Takes the output from bpftrace and processes it ready for sending to Graphite
### wrapper
The wrapper script that runs bpftrace and sends output to process-trace which processes the data delivered, then uploads to Graphite
### process-trace.service
systemd service configuration. Set to keep the bpftrace and process-trace service running and delivering data non-interactively.
To install:
- cp process-trace.service /etc/systemd/system
- chmod 644 /etc/systemd/system/process-trace.service
- systemctl daemon-reload
- systemctl status process-trace
- systemctl start process-trace
- systemctl enable process-trace    # Assuming the service should restart after a reboot
### process-trace.go
Processing script. At a high level the code:
- Looks for "@usecs:" in the data coming in. One of these should appear every second from bpftrace
- Then deals with all the data that follows @usecs:
- Every 15 times "@usecs:" is seen, the data is formatted ready for Graphite and then sent. NOTE that currently the process will attempt to open port 2003 on the localhost as the IP address/hostname of the Graphite server isn't known
- Every 15 minutes the network connection to Graphite is closed and a new connection is opened. This is to avoid opening and closing many network ports rapidly and potentially starving the server we're running on of resources.