#!/bin/sh
eval "$(ssh-agent)"
ssh-add /home/iomz/.ssh/id_ed25519
/home/iomz/go/bin/autolog
