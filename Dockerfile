FROM docker.io/alpine

ADD ./build/miniDevOps /covid-away/miniDevOps

CMD /covid-away/miniDevOps