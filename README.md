# SMS - Segel Management System

<img width="200" height="200" alt="bild" src="https://github.com/user-attachments/assets/c9af1f75-cd7b-4e9d-a583-c58348aafe8d" />

## Introduction
Have you ever wanted to show off your homelabbing to other people? A problem I face when doing this is how you do it in an understandable way. This is a way for me to do it: The Segel Management System. 
It's basically a framework you can install and make use of, in order to show technical details about your homelab environments to people and LLMs in accessible formats.

## Capabilities
- List deployments in a kubernetes cluster

## Arcitechture
- Warehouse: all your data is saved here <- this is the big thing
- Backend: Provides warehouse data to clients
- Clients (web, markdown, csv): Shows warehouse data to humans 
