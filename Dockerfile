FROM docker.io/alpine

ADD ./build/covid-away /covid-away/covid-away

CMD /covid-away/covid-away