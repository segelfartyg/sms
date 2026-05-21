# SMS - Segel Management System

<img width="200" height="200" alt="bild" src="https://github.com/user-attachments/assets/c9af1f75-cd7b-4e9d-a583-c58348aafe8d" />

## Introduction

Have you ever wanted to show off your homelab to other people? I have. Personally I find it very difficult to stay up to date with my devices and software, and also to provide the current state of my homelab to other people in an understandable way. Often I just talk about it, which leads to me always forgetting something :P This is a way for me to do avoid both of those things: The Segel Management System.

This is basically just a framework you can make use of, to provide up to date technical details about your homelab environment, both to people and LLMs in accessible formats. The fundamental thing about SMS is that you should configure it once, then it keeps track of things for you. However you dont have to use it this way, but I will :)

## Capabilities

- List deployments in a Kubernetes cluster
- Multi platform CLI tool for reporting technical details using fastfetch

## Arcitechture

- Warehouse: all your data is saved here <- this is the big thing
- Backend: Provides warehouse data to clients
- Clients (web, markdown, csv): Shows warehouse data to humans
