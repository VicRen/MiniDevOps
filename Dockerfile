FROM docker.io/alpine

ADD ./build/miniDevOps /minidevops/miniDevOps

CMD /minidevops/miniDevOps