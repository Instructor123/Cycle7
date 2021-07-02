# Cycle 7
## Project Description
The goal this cycle was to create a server/client environment that would automatically setup and configure a Wireguard VPN to enable communication. This server/client combo only establishes the VPN connection, requiring users to actually use it in whatever way/means they want.

## Reason
I spent some time playing with Wireguard in a previous course and really enjoyed it, but became frustrated with the tool's interface and documentation. I wanted to create something that utilized Wireguard but kept as many of the command details hidden from the user as possible. 

An additional artifact from this work is a Golang module for wireguard. It currently implements minimal functionality, but it separates the needed functionality of wireguard from using the client/servers themselves.

## Future Work
Adding more features to the wireguard module to reflect everything that the wireguard program can do is a must. Exposing the whole tool would offer more flexibility, and really allow for a wide range of features to be explored.

Also adding the ability for different clients to communicate as peers would be really cool; essentially creating a distributed mesh network (to use the buzz words.) This would require a decent amount of work to extend the server and client roles, in addition to adding more wireguard features. A key feature would also be interface management. Currenly users must manually delete the interface (wg0) when they are done using it.
