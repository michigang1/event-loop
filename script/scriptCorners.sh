#!/usr/bin/env bash
curl -d "reset
  figure 1 1
  figure 0 0
  figure 0 1
  figure 1 0
  figure 0.5 0.5
  update" http://localhost:17000