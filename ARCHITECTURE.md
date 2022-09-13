# Architecture

This documents describes the high-level architecture of goGoPing.

At the highest level this project is a port of open source C application commonly shipped in Linux and BSD systems. Really it is meant as a fun side project, but has the added benefit of added concurrency to the application.

## Code Map

The code is organized by application activities and by the packages required for those applications.

/cli

Contains a wrapper around an instance of golang.org/x/term and helper functions to support user cli activities

/net

Contains code that performs networking. This includes IMCP communication, opening sockets, and DNS resolution.