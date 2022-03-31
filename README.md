# honeypot

This is a golang honeypot that implements the gliderlabs ssh package.

The aim of the honeypot is to trick an "attacker" into logging into an
ssh server hosted on port 2222.

Having done that the honeypot logs that information, reports and error, and sends a push notification to the user's phone.
